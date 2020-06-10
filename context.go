// context 中不能使用 global 中的方法打印日志， global 会调用 context 的方法，会陷入循环

package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// Ctxkey context key 类型
type Ctxkey string

const (
	// CtxLoggerName define the ctx logger name
	CtxLoggerName = "ctxLogger"
	// TraceIDKey define the trace id keyname
	TraceIDKey = "traceID"

	// CtxkeyCtxLoggerName context keyname
	CtxkeyCtxLoggerName Ctxkey = CtxLoggerName
	// CtxkeyTraceID context keyname
	CtxkeyTraceID Ctxkey = TraceIDKey
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLogger *zap.Logger
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(CtxLoggerName)
	} else {
		ctxLoggerItf = c.Value(CtxkeyCtxLoggerName)
	}

	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		ctxLogger = CloneDefaultLogger(CtxLoggerName)
	}

	// try to get trace id from ctx
	traceID := CtxTraceID(c)
	// then set trace id into ctxlogger
	ctxLogger = ctxLogger.With(zap.String(TraceIDKey, traceID))

	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxTraceID get trace id from context
func CtxTraceID(c context.Context) string {
	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(TraceIDKey); traceID != "" {
			return traceID
		}
	}
	// get from go context
	traceIDItf := c.Value(CtxkeyTraceID)
	if traceIDItf != nil {
		return traceIDItf.(string)
	}
	// return default value
	return "logging-" + xid.New().String()
}

// Context set trace id and logger into context.Context and gin.Context
func Context(c context.Context, logger *zap.Logger, traceID string) context.Context {
	if gc, ok := c.(*gin.Context); ok {
		// set ctxlogger in gin.Context
		gc.Set(CtxLoggerName, logger)
		// set traceID in gin.Context
		gc.Set(TraceIDKey, traceID)
	}
	// set ctxlogger in context.Context
	c = context.WithValue(c, CtxkeyCtxLoggerName, logger)
	// set traceID in context.Context
	c = context.WithValue(c, CtxkeyTraceID, traceID)
	return c
}
