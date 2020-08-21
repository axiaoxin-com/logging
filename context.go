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

var (
	// CtxLoggerName define the ctx logger name
	CtxLoggerName Ctxkey = "ctxlogger"
	// TraceIDKeyname define the trace id keyname
	TraceIDKeyname Ctxkey = "traceid"
	// TraceIDPrefix set the prefix when gen a trace id
	TraceIDPrefix string = "logging-"
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLogger *zap.Logger
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(string(CtxLoggerName))
	} else {
		ctxLoggerItf = c.Value(CtxLoggerName)
	}

	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		ctxLogger = CloneDefaultLogger(string(CtxLoggerName))
	}

	// try to get trace id from ctx
	traceID := CtxTraceID(c)
	// then set trace id into ctxlogger
	ctxLogger = ctxLogger.With(zap.String(string(TraceIDKeyname), traceID))

	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxTraceID get trace id from context
func CtxTraceID(c context.Context) string {
	if c == nil {
		c = context.Background()
	}
	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(string(TraceIDKeyname)); traceID != "" {
			return traceID
		}
	}
	// get from go context
	traceIDItf := c.Value(TraceIDKeyname)
	if traceIDItf != nil {
		return traceIDItf.(string)
	}
	// return default value
	return TraceIDPrefix + xid.New().String()
}

// Context set trace id and logger into context.Context and gin.Context
func Context(c context.Context, logger *zap.Logger, traceID string) context.Context {
	if c == nil {
		c = context.Background()
	}
	if gc, ok := c.(*gin.Context); ok {
		// set ctxlogger in gin.Context
		gc.Set(string(CtxLoggerName), logger)
		// set traceID in gin.Context
		gc.Set(string(TraceIDKeyname), traceID)
	}
	// set ctxlogger in context.Context
	c = context.WithValue(c, CtxLoggerName, logger)
	// set traceID in context.Context
	c = context.WithValue(c, TraceIDKeyname, traceID)
	return c
}
