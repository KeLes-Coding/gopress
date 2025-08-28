// package logger 封装了日志系统的初始化和配置。
// 目标是提供一个全局的、配置好的日志记录器(Logger)，
// 让其他包可以方便地使用它来记录日志。
package logger

import (
	"os"
	"path/filepath"

	"github.com/KeLes-Coding/gopress/internal/config" // 导入我们自己的 config 包
	"go.uber.org/zap"                                 // 导入 zap 核心库
	"go.uber.org/zap/zapcore"                         // 导入 zapcore，用于更精细的配置
	"gopkg.in/natefinch/lumberjack.v2"                // 导入 lumberjack 用于日志切割
)

// L 是一个全局的 zap Logger 指针。
// 项目的其他部分将通过 logger.L 来访问这个全局的日志记录器实例。
var L *zap.Logger

// Init 函数负责根据配置文件初始化全局的 Logger。
func Init() error {
	// getLogWriter 函数返回一个 io.Writer，它决定了日志将被写入到哪里。
	// 这里我们配置它同时写入文件和控制台。
	writeSyncer := getLogWriter(
		config.Conf.Log.Filename,
		config.Conf.Log.MaxSize,
		config.Conf.Log.MaxBackups,
		config.Conf.Log.MaxAge,
		config.Conf.Log.Compress,
	)

	// getEncoder 函数返回一个 zapcore.Encoder，它决定了日志的格式（例如 JSON 或 Console）。
	encoder := getEncoder()

	// level 是一个 zapcore.Level 类型的变量，用于表示日志级别。
	level := new(zapcore.Level)
	// UnmarshalText 会将字符串形式的日志级别（如 "debug"）解析为 zapcore.Level 类型。
	if err := level.UnmarshalText([]byte(config.Conf.Log.Level)); err != nil {
		return err
	}

	// zapcore.NewCore 创建一个 zap 的核心(Core)。
	// Core 需要三个基本组件：
	// 1. Encoder: 决定日志的格式。
	// 2. WriteSyncer: 决定日志的输出位置。
	// 3. LevelEnabler: 决定哪些级别的日志应该被记录。
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// zap.New 根据创建的 Core 构建 Logger。
	// zap.AddCaller() 是一个选项，它会在日志中添加调用者的文件名和行号，非常有助于调试。
	L = zap.New(core, zap.AddCaller())

	// zap.ReplaceGlobals(L) 将我们自定义的 Logger 替换掉 zap 库的全局 Logger。
	// 这样，我们就可以在项目中使用 zap.S() 和 zap.L() 来获取全局的 SugaredLogger 和 Logger。
	zap.ReplaceGlobals(L)

	return nil
}

// getEncoder 函数用于创建并配置 zapcore.Encoder。
func getEncoder() zapcore.Encoder {
	// NewProductionEncoderConfig 是 zap 提供的一个生产环境的推荐编码器配置。
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间戳的格式为 ISO8601 (例如 "2023-10-27T10:00:00.000Z")。
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 将日志级别字符串转换为大写（例如 "INFO", "DEBUG"）。
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 根据配置文件中的 format 选项，返回不同的编码器。
	if config.Conf.Log.Format == "json" {
		// NewJSONEncoder 会以 JSON 格式输出日志。
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	// NewConsoleEncoder 会以更适合人类阅读的格式（例如 "INFO logger/logger.go:42 message"）输出日志。
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter 函数配置日志的输出目标，集成了 lumberjack 来实现日志切割。
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	// 确保日志文件所在的目录存在，如果不存在则创建。
	logDir := filepath.Dir(filename)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// os.MkdirAll 会创建所有必要的父目录。
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	// 使用 lumberjack 库来创建日志切割器。
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件的位置
		MaxSize:    maxSize,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		Compress:   compress,  // 是否压缩/归档旧文件
	}

	// zapcore.NewMultiWriteSyncer 可以将日志同时输出到多个目标。
	// 这里我们将日志同时输出到标准输出（控制台）和配置好的日志文件。
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
