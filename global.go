// Package logging provides ...
package logging

import "go.uber.org/zap"

// SDebug 使用logging默认的slogger记录debug级别的日志
// logging.Debug("abc", 123)
func SDebug(args ...interface{}) {
	defer slogger.Sync()
	slogger.Debug(args...)
}

// SInfo 使用logging默认的slogger记录info级别的日志
func SInfo(args ...interface{}) {
	defer slogger.Sync()
	slogger.Info(args...)
}

// SWarn 使用logging默认的slogger记录warn级别的日志
func SWarn(args ...interface{}) {
	defer slogger.Sync()
	slogger.Warn(args...)
}

// SError 使用logging默认的slogger记录Error级别的日志
func SError(args ...interface{}) {
	defer slogger.Sync()
	slogger.Error(args...)
}

// SPanic 使用logging默认的slogger记录Panic级别的日志
func SPanic(args ...interface{}) {
	defer slogger.Sync()
	slogger.Panic(args...)
}

// SFatal 使用logging默认的slogger记录Fatal级别的日志
func SFatal(args ...interface{}) {
	defer slogger.Sync()
	slogger.Fatal(args...)
}

// SDebugf 使用logging默认的slogger模板字符串记录debug级别的日志
// logging.Debugf("str:%s", "abd")
func SDebugf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Debugf(template, args...)
}

// SInfof 使用logging默认的slogger模板字符串记录info级别的日志
func SInfof(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Infof(template, args...)
}

// SWarnf 使用logging默认的slogger模板字符串记录warn级别的日志
func SWarnf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Warnf(template, args...)
}

// SErrorf 使用logging默认的slogger模板字符串记录error级别的日志
func SErrorf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Errorf(template, args...)
}

// SPanicf 使用logging默认的slogger模板字符串记录panic级别的日志
func SPanicf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Panicf(template, args...)
}

// SFatalf 使用logging默认的slogger模板字符串记录fatal级别的日志
func SFatalf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Fatalf(template, args...)
}

// SDebugw 使用logging默认的sloggerkv记录debug级别的日志
// logging.Debugw("msg", "k1", "v1", "k2", "v2")
func SDebugw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Debugw(msg, keysAndValues...)
}

// SInfow 使用logging默认的sloggerkv记录info级别的日志
func SInfow(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Infow(msg, keysAndValues...)
}

// SWarnw 使用logging默认的slogger kv记录warn级别的日志
func SWarnw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Warnw(msg, keysAndValues...)
}

// SErrorw 使用logging默认的slogger kv记录error级别的日志
func SErrorw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Errorw(msg, keysAndValues...)
}

// SPanicw 使用logging默认的slogger kv记录panic级别的日志
func SPanicw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Panicw(msg, keysAndValues...)
}

// SFatalw 使用logging默认的slogger kv记录fatal级别的日志
func SFatalw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Fatalw(msg, keysAndValues...)
}

// Debug 使用logging默认的logger记录debug级别的日志
func Debug(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Debug(msg, fields...)
}

// Info 使用logging默认的logger记录info级别的日志
func Info(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Info(msg, fields...)
}

// Warn 使用logging默认的logger记录warn级别的日志
func Warn(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Warn(msg, fields...)
}

// Error 使用logging默认的logger记录error级别的日志
func Error(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Error(msg, fields...)
}

// Panic 使用logging默认的logger记录panic级别的日志
func Panic(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Panic(msg, fields...)
}

// Fatal 使用logging默认的logger记录fatal级别的日志
func Fatal(msg string, fields ...zap.Field) {
	defer logger.Sync()
	logger.Fatal(msg, fields...)
}
