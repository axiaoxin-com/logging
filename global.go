// 开箱即用的方法

package logging

import "go.uber.org/zap"

// Debugs 使用logging默认的slogger记录debug级别的日志
// logging.Debugs("abc", 123)
func Debugs(args ...interface{}) {
	defer slogger.Sync()
	slogger.Debug(args...)
}

// Infos 使用logging默认的slogger记录info级别的日志
func Infos(args ...interface{}) {
	defer slogger.Sync()
	slogger.Info(args...)
}

// Warns 使用logging默认的slogger记录warn级别的日志
func Warns(args ...interface{}) {
	defer slogger.Sync()
	slogger.Warn(args...)
}

// Errors 使用logging默认的slogger记录Error级别的日志
func Errors(args ...interface{}) {
	defer slogger.Sync()
	slogger.Error(args...)
}

// Panics 使用logging默认的slogger记录Panic级别的日志
func Panics(args ...interface{}) {
	defer slogger.Sync()
	slogger.Panic(args...)
}

// Fatals 使用logging默认的slogger记录Fatal级别的日志
func Fatals(args ...interface{}) {
	defer slogger.Sync()
	slogger.Fatal(args...)
}

// Debugf 使用logging默认的slogger模板字符串记录debug级别的日志
// logging.Debugf("str:%s", "abd")
func Debugf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Debugf(template, args...)
}

// Infof 使用logging默认的slogger模板字符串记录info级别的日志
func Infof(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Infof(template, args...)
}

// Warnf 使用logging默认的slogger模板字符串记录warn级别的日志
func Warnf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Warnf(template, args...)
}

// Errorf 使用logging默认的slogger模板字符串记录error级别的日志
func Errorf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Errorf(template, args...)
}

// Panicf 使用logging默认的slogger模板字符串记录panic级别的日志
func Panicf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Panicf(template, args...)
}

// Fatalf 使用logging默认的slogger模板字符串记录fatal级别的日志
func Fatalf(template string, args ...interface{}) {
	defer slogger.Sync()
	slogger.Fatalf(template, args...)
}

// Debugw 使用logging默认的sloggerkv记录debug级别的日志
// logging.Debugw("msg", "k1", "v1", "k2", "v2")
func Debugw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Debugw(msg, keysAndValues...)
}

// Infow 使用logging默认的sloggerkv记录info级别的日志
func Infow(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Infow(msg, keysAndValues...)
}

// Warnw 使用logging默认的slogger kv记录warn级别的日志
func Warnw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Warnw(msg, keysAndValues...)
}

// Errorw 使用logging默认的slogger kv记录error级别的日志
func Errorw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Errorw(msg, keysAndValues...)
}

// Panicw 使用logging默认的slogger kv记录panic级别的日志
func Panicw(msg string, keysAndValues ...interface{}) {
	defer slogger.Sync()
	slogger.Panicw(msg, keysAndValues...)
}

// Fatalw 使用logging默认的slogger kv记录fatal级别的日志
func Fatalw(msg string, keysAndValues ...interface{}) {
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
