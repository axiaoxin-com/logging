// gorm v2

package logging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

const (
	// GormLoggerName gorm logger 名称
	GormLoggerName = "gorm"
)

// GormLogger 使用 zap 来打印 gorm 的日志
// 初始化时在内部的 logger 中添加 trace id 可以追踪 sql 执行记录
type GormLogger struct {
	logger *zap.Logger
	// 日志级别
	logLevel zapcore.Level
	// 指定慢查询时间
	slowThreshold time.Duration
}

var gormLogLevelMap = map[gormlogger.LogLevel]zapcore.Level{
	gormlogger.Info:  zap.InfoLevel,
	gormlogger.Warn:  zap.WarnLevel,
	gormlogger.Error: zap.ErrorLevel,
}

// LogMode 实现 gorm logger 接口方法
func (g GormLogger) LogMode(gormLogLevel gormlogger.LogLevel) gormlogger.Interface {
	zaplevel, exists := gormLogLevelMap[gormLogLevel]
	if !exists {
		zaplevel = zap.DebugLevel
	}
	newlogger := g
	newlogger.logLevel = zaplevel
	return &newlogger
}

// Info 实现 gorm logger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gorm logger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gorm logger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace 实现 gorm logger 接口方法
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	now := time.Now()
	latency := now.Sub(begin).Seconds()
	sql, rows := fc()
	sql = strings.Replace(sql, "\t", " ", -1)
	sql = strings.Replace(sql, "\n", " ", -1)
	switch {
	case err != nil:
		CtxLogger(ctx).Error("gorm trace [error]", zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		CtxLogger(ctx).Warn("gorm trace [slow]", zap.String("threshold", fmt.Sprintf("%v", g.slowThreshold)), zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	case g.logLevel <= zap.InfoLevel:
		CtxLogger(ctx).Info("gorm trace", zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	}
}

// NewGormLogger 返回带 zap logger 的 GormLogger
func NewGormLogger(logger *zap.Logger, logLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	logger = logger.Named(GormLoggerName).WithOptions(zap.AddCallerSkip(7))
	return GormLogger{
		logger:        logger,
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}

// CtxGormLogger 从 context 中获取 zap logger 来创建 GormLogger
func CtxGormLogger(c context.Context, logLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	logger := CtxLogger(c).Named(GormLoggerName).WithOptions(zap.AddCallerSkip(7))
	return GormLogger{
		logger:        logger,
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}
