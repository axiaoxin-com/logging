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
		Debug("no logger in context, clone the default logger as ctxLogger")
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
			Debug("get trace id from gin.Context")
			return traceID
		}
	}
	// get from go context
	traceIDItf := c.Value(TraceIDKey)
	if traceIDItf != nil {
		Debug("get trace id from context.Context")
		return traceIDItf.(string)
	}
	// return default value
	Debug("context dose not exist trace id key, generate a new trace id")
	return "logging-" + xid.New().String()
}

// Context set ctxlogger and trace id into context.Context and gin.Context
func Context(c context.Context, traceID string) context.Context {
	// set trace id in ctxlogger
	logger := CtxLogger(c, zap.String(TraceIDKey, traceID))

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
