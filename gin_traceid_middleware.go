package logging

import (
	"github.com/gin-gonic/gin"
)

// TraceIDFunc 生成 trace id 的回调函数类型
type TraceIDFunc func(*gin.Context) string

// GinTraceIDMiddleware is a gin middleware for gen a trace id in context
func GinTraceIDMiddleware(traceIDFunc TraceIDFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用传入的回调方法获取 traceid ，不存在则新生成
		traceID := traceIDFunc(c)
		if traceID == "" {
			traceID = CtxTraceID(c)
		}
		// 设置 trace id 到 header 中
		c.Writer.Header().Set(TraceIDKey, traceID)
		// 设置 trace id 和 ctxLogger 到 context 中
		Context(c, DefaultLogger(), traceID)

		c.Next()
	}
}

// GetTraceIDFromHeader 从 request header 中获取 key 为 TraceIDKey 的值作为 traceid
func GetTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(TraceIDKey)
}

// GetTraceIDFromQueryString 从 querystring 中获取 key 为 TraceIDKey 的值作为 traceid
func GetTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(TraceIDKey)
}

// GetTraceIDFromPostForm 从 postform 中获取 key 为 TraceIDKey 的值作为 traceid
func GetTraceIDFromPostForm(c *gin.Context) string {
	return c.PostForm(TraceIDKey)
}
