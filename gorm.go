// gorm v2

package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

var (
	// GormLoggerName gorm logger 名称
	GormLoggerName = "gorm"
	// GormLoggerCallerSkip caller skip
	GormLoggerCallerSkip = 3
)

// GormLogger 使用 zap 来打印 gorm 的日志
// 初始化时在内部的 logger 中添加 trace id 可以追踪 sql 执行记录
type GormLogger struct {
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

// CtxLogger 创建打印日志的 ctxlogger
func (g GormLogger) CtxLogger(ctx context.Context) *zap.Logger {
	return CtxLogger(ctx).Named(GormLoggerName).WithOptions(zap.AddCallerSkip(GormLoggerCallerSkip))
}

// Info 实现 gorm logger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		g.CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gorm logger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		g.CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gorm logger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		g.CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace 实现 gorm logger 接口方法
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	now := time.Now()
	latency := now.Sub(begin).Seconds()
	sql, rows := fc()
	sql = goutils.RemoveDuplicateWhitespace(sql, true)
	switch {
	case err != nil:
		g.CtxLogger(ctx).Error("gorm sql error trace", zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		g.CtxLogger(ctx).Warn("gorm sql slow trace", zap.String("threshold", fmt.Sprintf("%v", g.slowThreshold)), zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	case g.logLevel <= zap.InfoLevel:
		g.CtxLogger(ctx).Info("gorm sql trace", zap.Float64("latency", latency), zap.String("sql", sql), zap.Int64("rows", rows))
	}
}

// NewGormLogger 返回带 zap logger 的 GormLogger
func NewGormLogger(logLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	return GormLogger{
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}
