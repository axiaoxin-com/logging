package logging

import (
	"github.com/gin-gonic/gin"
)

// TraceIDFunc 生成trace id的回调函数类型
type TraceIDFunc func(*gin.Context) string

// GinTraceIDMiddleware is a gin middleware for gen a trace id in context
func GinTraceIDMiddleware(traceIDFunc TraceIDFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用传入的回调方法获取traceid，不存在则新生成
		traceID := traceIDFunc(c)
		if traceID == "" {
			traceID = CtxTraceID(c)
		}
		// 设置 trace id 到 header 中
		c.Writer.Header().Set(TraceIDKey, traceID)
		// 设置 trace id 和 ctxLogger 和 context 中
		Context(c, traceID)

		c.Next()
	}
}

// GetTraceIDFromHeader 从request header中获取key为TraceIDKey的值作为traceid
func GetTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(TraceIDKey)
}

// GetTraceIDFromQueryString 从querystring中获取key为TraceIDKey的值作为traceid
func GetTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(TraceIDKey)
}

// GetTraceIDFromPostForm 从postform中获取key为TraceIDKey的值作为traceid
func GetTraceIDFromPostForm(c *gin.Context) string {
	return c.PostForm(TraceIDKey)
}
