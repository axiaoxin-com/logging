// context中不能使用global中的方法打印日志，global会调用context的方法，会陷入循环

package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

const (
	// CtxLoggerKey define the logger keyname which in context
	CtxLoggerKey = "ctxLogger"
	// TraceIDKey define the trace id keyname
	TraceIDKey = "traceID"
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLogger *zap.Logger
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(CtxLoggerKey)
	} else {
		ctxLoggerItf = c.Value(CtxLoggerKey)
	}

	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		ctxLogger = CloneDefaultLogger(CtxLoggerKey)
	}
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
	traceIDItf := c.Value(TraceIDKey)
	if traceIDItf != nil {
		return traceIDItf.(string)
	}
	// return default value
	return "logging-" + xid.New().String()
}

// Context set trace id in logger,then set trace id and logger into context.Context and gin.Context
func Context(c context.Context, logger *zap.Logger, traceID string) context.Context {
	// set trace id in ctxlogger
	logger = logger.With(zap.String(TraceIDKey, traceID))
	if gc, ok := c.(*gin.Context); ok {
		// set ctxlogger in gin.Context
		gc.Set(CtxLoggerKey, logger)
		// set traceID in gin.Context
		gc.Set(TraceIDKey, traceID)
	}
	// set ctxlogger in context.Context
	c = context.WithValue(c, CtxLoggerKey, logger)
	// set traceID in context.Context
	c = context.WithValue(c, TraceIDKey, traceID)
	return c
}
