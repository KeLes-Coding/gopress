// package logger ...
package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // 导入 GORM 的 logger 接口定义
)

// GormLogger 是一个自定义的结构体，它将实现 gorm/logger.Interface 接口。
// 这样做的目的是为了将 GORM 产生的日志，通过我们自己的 zap 日志系统来输出，
// 从而统一项目的日志格式和管理方式。
type GormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel // GORM 的日志级别
}

// NewGormLogger 创建一个新的 GormLogger 实例。
func NewGormLogger(zapLogger *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger: zapLogger,
		LogLevel:  logger.Info, // 默认的 GORM 日志级别
	}
}

// LogMode 是 gorm/logger.Interface 接口要求实现的方法。
// 它返回一个新的 logger 实例，但日志级别被设置为指定的 level。
// 这允许 GORM 在不同场景下动态改变日志级别。
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 方法用于记录 Info 级别的日志。
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 方法用于记录 Warn 级别的日志。
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 方法用于记录 Error 级别的日志。
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 方法是 GORM 日志的核心，它在每条 SQL 执行前后被调用，用于记录 SQL 语句、执行时间和错误。
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 如果日志级别设置为 Silent，则不记录任何 trace 信息。
	if l.LogLevel <= logger.Silent {
		return
	}

	// 计算 SQL 执行耗时。
	elapsed := time.Since(begin)
	// fc() 函数会返回 SQL 语句和受影响的行数。
	sql, rows := fc()

	// 错误处理逻辑：
	// 如果存在错误，并且这个错误不是 "记录未找到" (gorm.ErrRecordNotFound)，
	// 那么我们就以 Error 级别记录这条日志。
	// gorm.ErrRecordNotFound 通常是业务逻辑的一部分，不应被视为系统级错误。
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.ZapLogger.Error("GORM Trace",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
		return
	}

	// 慢查询日志逻辑：
	// 定义一个慢查询阈值，例如 200 毫秒。
	slowThreshold := 200 * time.Millisecond
	// 如果执行时间超过了阈值，就以 Warn 级别记录日志，提醒开发者注意优化。
	if elapsed > slowThreshold {
		l.ZapLogger.Warn("GORM Trace (Slow Query)",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
		return
	}

	// 正常 SQL 查询逻辑：
	// 如果日志级别允许（Info 或更高），就以 Debug 级别记录所有正常的 SQL 执行。
	// 使用 Debug 级别可以避免在生产环境中输出过多的 SQL 日志，只在需要时开启。
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Debug("GORM Trace",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
