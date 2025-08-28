// package main 表明这是一个可执行程序。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KeLes-Coding/gopress/internal/api/middleware"
	"github.com/KeLes-Coding/gopress/internal/config"
	"github.com/KeLes-Coding/gopress/internal/dao"
	"github.com/KeLes-Coding/gopress/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// main 函数是整个程序的入口点。
func main() {
	// --- 1. 初始化配置 ---
	// 程序启动的第一步是加载配置文件。
	// 如果配置加载失败，程序将无法正常运行，因此直接 panic。
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize config: %v", err))
	}

	// --- 2. 初始化日志 ---
	// 在配置加载后，立即初始化日志系统。
	// 这样，后续的所有启动步骤都可以使用日志记录器。
	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	// 使用 defer 语句确保在 main 函数退出时，将日志缓冲区的日志同步到文件中，防止日志丢失。
	defer func() {
		_ = logger.L.Sync()
	}()

	// --- 3. 初始化 MySQL 连接 ---
	if err := dao.InitMySQL(); err != nil {
		// 如果数据库连接失败，这是一个致命错误，程序无法继续。
		logger.L.Fatal("Failed to initialize MySQL", zap.Error(err))
	}

	// --- 4. 自动迁移数据表 ---
	// 在开发环境中，自动迁移表结构非常方便。
	if err := dao.AutoMigrateTables(); err != nil {
		logger.L.Fatal("Failed to auto-migrate tables", zap.Error(err))
	}

	// --- 5. 设置 Gin 模式并创建引擎 ---
	gin.SetMode(config.Conf.Server.Mode)
	// gin.New() 创建一个不带任何默认中间件的纯净的 Gin 引擎。
	// 这给了我们完全的控制权来决定使用哪些中间件。
	r := gin.New()
	// 注册我们自定义的日志中间件和 Gin 官方的 Recovery 中间件。
	// Recovery 中间件可以在发生 panic 时捕获它，并返回一个 500 错误，防止整个程序崩溃。
	r.Use(middleware.GinLogger(logger.L), gin.Recovery())

	// --- 6. 注册路由 ---
	r.GET("/ping", func(c *gin.Context) {
		logger.L.Info("Received ping request", zap.String("client_ip", c.ClientIP()))
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// --- 7. 启动服务并实现优雅关停 (Graceful Shutdown) ---
	// 创建一个 http.Server 实例，这样我们可以更好地控制服务的启动和关闭。
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler: r, // 将 Gin 引擎作为请求处理器
	}

	// 开启一个新的 goroutine 来启动服务，这样主 goroutine 就不会被阻塞。
	go func() {
		logger.L.Info("Server is starting...",
			zap.Int("port", config.Conf.Server.Port),
			zap.Int("pid", os.Getpid()),
		)
		// srv.ListenAndServe() 会阻塞，直到服务被关闭或发生错误。
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 如果错误不是 http.ErrServerClosed (这是我们主动关闭服务时产生的)，
			// 那么就是一个意外的致命错误。
			logger.L.Fatal("Failed to listen and serve", zap.Error(err))
		}
	}()

	// 主 goroutine 会在这里等待中断信号。
	// 我们创建一个 channel 来接收 os.Signal。
	quit := make(chan os.Signal, 1)
	// signal.Notify 会将指定的信号（这里是 SIGINT 和 SIGTERM）转发到 quit channel。
	// SIGINT: 用户按下 Ctrl+C。
	// SIGTERM: 系统（如 Docker, Kubernetes）发出的标准终止信号。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 程序会在这里阻塞，直到接收到一个信号。

	logger.L.Info("Shutting down server...")

	// 创建一个带有超时的 context，用于通知 srv.Shutdown 我们最多等待 5 秒。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保在退出时释放 context 相关的资源。

	// srv.Shutdown 会尝试优雅地关闭服务器。
	// 它会停止接收新的请求，并等待当前正在处理的请求完成，直到超时。
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.L.Info("Server exiting.")
}
