// context 中不能使用 global 中的方法打印日志， global 会调用 context 的方法，会陷入循环

package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// Ctxkey context key 类型
type Ctxkey string

var (
	// CtxLoggerName define the ctx logger name
	CtxLoggerName Ctxkey = "ctx_logger"
	// TraceIDKeyname define the trace id keyname
	TraceIDKeyname Ctxkey = "trace_id"
	// TraceIDPrefix set the prefix when gen a trace id
	TraceIDPrefix = "logging_"
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(string(CtxLoggerName))
	} else {
		ctxLoggerItf = c.Value(CtxLoggerName)
	}

	var ctxLogger *zap.Logger
	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		_, ctxLogger = NewCtxLogger(c, CloneLogger(string(CtxLoggerName)), CtxTraceID(c))
	}

	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxTraceID get trace id from context
// Modify TraceIDPrefix change change the prefix
func CtxTraceID(c context.Context) string {
	if c == nil {
		c = context.Background()
	}
	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(string(TraceIDKeyname)); traceID != "" {
			return traceID
		}
		if traceID := gc.Query(string(TraceIDKeyname)); traceID != "" {
			return traceID
		}
		if traceID := jsoniter.Get(GetGinRequestBody(gc), string(TraceIDKeyname)).ToString(); traceID != "" {
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

// NewCtxLogger return a context with logger and trace id and a logger with trace id
func NewCtxLogger(c context.Context, logger *zap.Logger, traceID string) (context.Context, *zap.Logger) {
	if c == nil {
		c = context.Background()
	}
	if traceID == "" {
		traceID = CtxTraceID(c)
	}
	ctxLogger := logger.With(zap.String(string(TraceIDKeyname), traceID))
	if gc, ok := c.(*gin.Context); ok {
		// set ctxlogger in gin.Context
		gc.Set(string(CtxLoggerName), ctxLogger)
		// set traceID in gin.Context
		gc.Set(string(TraceIDKeyname), traceID)
	}
	// set ctxlogger in context.Context
	c = context.WithValue(c, CtxLoggerName, ctxLogger)
	// set traceID in context.Context
	c = context.WithValue(c, TraceIDKeyname, traceID)
	return c, ctxLogger
}
