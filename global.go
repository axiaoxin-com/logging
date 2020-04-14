// 开箱即用的方法
// 第一个参数为context，会尝试从其中获取带 trace id 的logger进行打印

package logging

import (
	"context"

	"go.uber.org/zap"
)

// Debugs 尝试从Context中获取带trace id的sugared logger来记录debug级别的日志
// logging.Debugs(nil, "abc", 123)
func Debugs(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Debug(args...)
}

// Infos 尝试从Context中获取带trace id的sugared logger来记录info级别的日志
func Infos(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Info(args...)
}

// Warns 尝试从Context中获取带trace id的sugared logger来记录warn级别的日志
func Warns(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Warn(args...)
}

// Errors 尝试从Context中获取带trace id的sugared logger来记录Error级别的日志
func Errors(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Error(args...)
}

// Panics 尝试从Context中获取带trace id的sugared logger来记录Panic级别的日志
func Panics(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Panic(args...)
}

// Fatals 尝试从Context中获取带trace id的sugared logger来记录Fatal级别的日志
func Fatals(c context.Context, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Fatal(args...)
}

// Debugf 尝试从Context中获取带trace id的sugared logger来模板字符串记录debug级别的日志
// logging.Debugf(nil, "str:%s", "abd")
func Debugf(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Debugf(template, args...)
}

// Infof 尝试从Context中获取带trace id的sugared logger来模板字符串记录info级别的日志
func Infof(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Infof(template, args...)
}

// Warnf 尝试从Context中获取带trace id的sugared logger来模板字符串记录warn级别的日志
func Warnf(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Warnf(template, args...)
}

// Errorf 尝试从Context中获取带trace id的sugared logger来模板字符串记录error级别的日志
func Errorf(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Errorf(template, args...)
}

// Panicf 尝试从Context中获取带trace id的sugared logger来模板字符串记录panic级别的日志
func Panicf(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Panicf(template, args...)
}

// Fatalf 尝试从Context中获取带trace id的sugared logger来模板字符串记录fatal级别的日志
func Fatalf(c context.Context, template string, args ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Fatalf(template, args...)
}

// Debugw 尝试从Context中获取带trace id的sugared logger来kv记录debug级别的日志
// logging.Debugw(nil, "msg", "k1", "v1", "k2", "v2")
func Debugw(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Debugw(msg, keysAndValues...)
}

// Infow 尝试从Context中获取带trace id的sugared logger来kv记录info级别的日志
func Infow(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Infow(msg, keysAndValues...)
}

// Warnw 尝试从Context中获取带trace id的sugared logger来 kv记录warn级别的日志
func Warnw(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Warnw(msg, keysAndValues...)
}

// Errorw 尝试从Context中获取带trace id的sugared logger来 kv记录error级别的日志
func Errorw(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Errorw(msg, keysAndValues...)
}

// Panicw 尝试从Context中获取带trace id的sugared logger来 kv记录panic级别的日志
func Panicw(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Panicw(msg, keysAndValues...)
}

// Fatalw 尝试从Context中获取带trace id的sugared logger来 kv记录fatal级别的日志
func Fatalw(c context.Context, msg string, keysAndValues ...interface{}) {
	slogger := CtxLogger(c).Sugar()
	defer slogger.Sync()
	slogger.Fatalw(msg, keysAndValues...)
}

// Debug 尝试从Context中获取带trace id的logger记录debug级别的日志
func Debug(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Debug(msg, fields...)
}

// Info 尝试从Context中获取带trace id的logger记录info级别的日志
func Info(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Info(msg, fields...)
}

// Warn 尝试从Context中获取带trace id的logger记录warn级别的日志
func Warn(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Warn(msg, fields...)
}

// Error 尝试从Context中获取带trace id的logger记录error级别的日志
func Error(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Error(msg, fields...)
}

// Panic 尝试从Context中获取带trace id的logger记录panic级别的日志
func Panic(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Panic(msg, fields...)
}

// Fatal 尝试从Context中获取带trace id的logger记录fatal级别的日志
func Fatal(c context.Context, msg string, fields ...zap.Field) {
	logger := CtxLogger(c)
	defer logger.Sync()
	logger.Fatal(msg, fields...)
}
