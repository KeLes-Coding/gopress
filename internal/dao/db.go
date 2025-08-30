// package dao (Data Access Object) 专门用于与数据库进行交互。
// 所有数据库相关的操作，如连接、CRUD，都应该封装在这一层。
package dao

import (
	"fmt"
	"time"

	"github.com/KeLes-Coding/gopress/internal/config" // 导入配置
	"github.com/KeLes-Coding/gopress/internal/logger" // 导入日志
	"github.com/KeLes-Coding/gopress/internal/model"  // 导入数据模型
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// _db 是一个包内私有的全局变量，用于持有数据库连接实例。
// 使用下划线前缀 `_` 是一种约定，表示这个变量不应该被包外直接访问。
var _db *gorm.DB

// InitMySQL 函数负责初始化到 MySQL 数据库的连接。
func InitMySQL() (err error) {
	// 从已加载的配置中获取 MySQL 的配置信息。
	c := config.Conf.MySQL
	// 构建 DSN (Data Source Name) 字符串，这是连接数据库所必需的信息。
	// 格式: "user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// parseTime=True: 告诉驱动需要将数据库中的时间类型(如 DATETIME)解析为 time.Time。
	// loc=Local: 设置时区为本地时区。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)

	// 根据我们应用的日志级别，来决定 GORM 的日志级别。
	var gormLogLevel gormlogger.LogLevel
	switch config.Conf.Log.Level {
	case "debug":
		gormLogLevel = gormlogger.Info // GORM 的 Info 级别会打印所有 SQL，适合调试
	case "info":
		gormLogLevel = gormlogger.Warn
	case "warn":
		gormLogLevel = gormlogger.Warn
	case "error":
		gormLogLevel = gormlogger.Error
	default:
		gormLogLevel = gormlogger.Info // 默认级别
	}

	// 实例化我们自定义的 GORM 日志记录器，并设置其日志级别。
	newLogger := logger.NewGormLogger(logger.L).LogMode(gormLogLevel)

	// gorm.Open 尝试打开一个数据库连接。
	// 它需要一个 gorm.Dialector (这里是 mysql.Open(dsn)) 和 gorm.Config。
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, // 将 GORM 的 logger 设置为我们自定义的 logger 实例。
	})
	if err != nil {
		// 如果连接失败，记录致命错误并返回。
		logger.L.Error("Failed to connect to MySQL", zap.Error(err))
		return
	}

	// GORM 使用 database/sql 包来维护连接池。
	// 我们可以通过 _db.DB() 获取底层的 *sql.DB 对象，并对其进行配置。
	sqlDB, err := _db.DB()
	if err != nil {
		logger.L.Error("Failed to get underlying sql.DB", zap.Error(err))
		return
	}
	// SetMaxIdleConns 设置连接池中的最大空闲连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置数据库的最大打开连接数。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置连接可以被重用的最长时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.L.Info("MySQL connection successful")
	return nil
}

// GetDB 函数提供了一个公共的访问方式来获取数据库连接实例。
// 其他包（如 service 层）应该通过这个函数来获取 *gorm.DB 对象，而不是直接访问 _db 变量。
func GetDB() *gorm.DB {
	return _db
}

// AutoMigrateTables 函数使用 GORM 的 AutoMigrate 功能来自动创建或更新数据库表结构。
// 这在开发阶段非常方便，可以确保数据库表结构与代码中的模型定义保持同步。
func AutoMigrateTables() error {
	// AutoMigrate 会检查传入的模型的表是否存在，如果不存在则创建。
	// 如果表已存在，它会检查并添加缺失的字段、索引等，但不会删除或修改现有的列。
	err := _db.AutoMigrate(
		&model.User{}, // 传入需要迁移的模型结构体的指针
		&model.Category{},
		&model.Tag{},
		&model.Post{},
		// &model.Post{},
		// &model.Category{},
	)
	if err != nil {
		logger.L.Error("Failed to auto-migrate tables", zap.Error(err))
		return err
	}
	logger.L.Info("Tables auto-migrated successfully")
	return nil
}
