// Package logging provides ...
package logging

import "go.uber.org/zap"

// SDebug 使用logging默认的SLogger记录debug级别的日志
// logging.Debug("abc", 123)
func SDebug(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Debug(args...)
}

// SInfo 使用logging默认的SLogger记录info级别的日志
func SInfo(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Info(args...)
}

// SWarn 使用logging默认的SLogger记录warn级别的日志
func SWarn(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Warn(args...)
}

// SError 使用logging默认的SLogger记录Error级别的日志
func SError(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Error(args...)
}

// SPanic 使用logging默认的SLogger记录Panic级别的日志
func SPanic(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Panic(args...)
}

// SFatal 使用logging默认的SLogger记录Fatal级别的日志
func SFatal(args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Fatal(args...)
}

// SDebugf 使用logging默认的SLogger模板字符串记录debug级别的日志
// logging.Debugf("str:%s", "abd")
func SDebugf(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Debugf(template, args...)
}

// SInfof 使用logging默认的SLogger模板字符串记录info级别的日志
func SInfof(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Infof(template, args...)
}

// SWarnf 使用logging默认的SLogger模板字符串记录warn级别的日志
func SWarnf(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Warnf(template, args...)
}

// SErrorf 使用logging默认的SLogger模板字符串记录error级别的日志
func SErrorf(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Errorf(template, args...)
}

// SPanicf 使用logging默认的SLogger模板字符串记录panic级别的日志
func SPanicf(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Panicf(template, args...)
}

// SFatalf 使用logging默认的SLogger模板字符串记录fatal级别的日志
func SFatalf(template string, args ...interface{}) {
	defer SLogger.Sync()
	SLogger.Fatalf(template, args...)
}

// SDebugw 使用logging默认的SLoggerkv记录debug级别的日志
// logging.Debugw("msg", "k1", "v1", "k2", "v2")
func SDebugw(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Debugw(msg, keysAndValues...)
}

// SInfow 使用logging默认的SLoggerkv记录info级别的日志
func SInfow(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Infow(msg, keysAndValues...)
}

// SWarnw 使用logging默认的SLogger kv记录warn级别的日志
func SWarnw(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Warnw(msg, keysAndValues...)
}

// SErrorw 使用logging默认的SLogger kv记录error级别的日志
func SErrorw(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Errorw(msg, keysAndValues...)
}

// SPanicw 使用logging默认的SLogger kv记录panic级别的日志
func SPanicw(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Panicw(msg, keysAndValues...)
}

// SFatalw 使用logging默认的SLogger kv记录fatal级别的日志
func SFatalw(msg string, keysAndValues ...interface{}) {
	defer SLogger.Sync()
	SLogger.Fatalw(msg, keysAndValues...)
}

// Debug 使用logging默认的Logger记录debug级别的日志
func Debug(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Debug(msg, fields...)
}

// Info 使用logging默认的Logger记录info级别的日志
func Info(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Info(msg, fields...)
}

// Warn 使用logging默认的Logger记录warn级别的日志
func Warn(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Warn(msg, fields...)
}

// Error 使用logging默认的Logger记录error级别的日志
func Error(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Error(msg, fields...)
}

// Panic 使用logging默认的Logger记录panic级别的日志
func Panic(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Panic(msg, fields...)
}

// Fatal 使用logging默认的Logger记录fatal级别的日志
func Fatal(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Fatal(msg, fields...)
}
