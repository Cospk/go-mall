package dao

import (
	"context"
	"errors"
	"github.com/Cospk/go-mall/pkg/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

// GormLogger 实现 gormLogger.Interface，并使用自定义 Logger 记录日志
type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 500 * time.Millisecond, // 超过500毫秒算慢查询, 如果需要按环境定制化, 这个放到配置文件中
	}
}

func (l *GormLogger) LogMode(lev gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.NewLogger(ctx).Info(msg, "data", data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.NewLogger(ctx).Error(msg, "data", data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.NewLogger(ctx).Error(msg, "data", data)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	duration := time.Since(begin).Milliseconds()
	// 获取 SQL 语句和返回条数
	sql, rows := fc()
	// Gorm 错误时记录错误日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.NewLogger(ctx).Error("SQL ERROR", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
	// 慢查询日志
	if duration > l.SlowThreshold.Milliseconds() {
		logger.NewLogger(ctx).Warn("SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", duration)
	} else {
		logger.NewLogger(ctx).Debug("SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
}
