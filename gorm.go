// 当前版本的 gorm 暂时不能通过在日志回调中获取 context ，无法直接添加 trace id 到 sql 日志中
// gorm v2 将会支持 context 特性(https://github.com/jinzhu/gorm/issues/2886),
// 当前 v2 还没有 release ，这里提供一种临时解决方案：
// 在每次即将进行 gorm 的操作时，都手动设置一个带有当前 trace id 的 logger 作为 gorm 的 logger 来打印日志

package logging

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

const (
	// GormLoggerName gorm logger 名称
	GormLoggerName = "gorm"
)

// GormLogger 使用 zap 来打印 gorm 的日志
// 初始化时在内部的 logger 中添加 trace id 可以追踪 sql 执行记录
type GormLogger struct {
	logger *zap.Logger
}

// Print 实现 gorm 定义的 logger 接口
// 使用 logger 的 debug 级别打印日志
func (g GormLogger) Print(values ...interface{}) {
	if values[0] == "sql" {
		g.logger.Debug(values[3].(string),
			zap.Any("vars", values[4]),
			zap.Int64("rowsAffected", values[5].(int64)),
			zap.Float64("duration", values[2].(time.Duration).Seconds()),
		)
	} else {
		g.logger.Sugar().Debug(values)
	}
}

// NewGormLogger 返回带 zap logger 的 GormLogger
func NewGormLogger(logger *zap.Logger) GormLogger {
	logger = logger.Named(GormLoggerName).WithOptions(zap.AddCallerSkip(7))
	return GormLogger{
		logger: logger,
	}
}

// CtxGormLogger 从 context 中获取 zap logger 来创建 GormLogger
func CtxGormLogger(c context.Context) GormLogger {
	logger := CtxLogger(c).Named(GormLoggerName).WithOptions(zap.AddCallerSkip(7))
	return GormLogger{
		logger: logger,
	}
}

// GormDBWithCtxLogger 为 gorm DB 设置 context 中的 logger ，返回带有新 logger 的 db 对象
func GormDBWithCtxLogger(c context.Context, db *gorm.DB) *gorm.DB {
	dbCopy := *db
	logger := CtxGormLogger(c)
	dbCopy.SetLogger(logger)
	return &dbCopy
}
