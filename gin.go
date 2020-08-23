package logging

import (
	"github.com/gin-gonic/gin"
)

// GinTraceIDFunc 在 gin 的 context 中 生成 trace id 的回调函数类型
type GinTraceIDFunc func(*gin.Context) string

// GinTraceID is a gin middleware for gen a trace id in context
func GinTraceID(traceIDFuncs ...GinTraceIDFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用传入的回调方法获取 traceid ，获取失败则新生成
		traceID := ""
		for _, traceIDFunc := range traceIDFuncs {
			traceID = traceIDFunc(c)
		}
		if traceID == "" {
			traceID = CtxTraceID(c)
		}
		// 设置 trace id 到 header 中
		c.Writer.Header().Set(string(TraceIDKeyname), traceID)
		// 设置 trace id 和 ctxLogger 到 context 中
		Context(c, DefaultLogger(), traceID)

		c.Next()
	}
}

// GetGinTraceIDFromHeader 从 gin 的 request header 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(string(TraceIDKeyname))
}

// GetGinTraceIDFromQueryString 从 gin 的 querystring 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(string(TraceIDKeyname))
}

// GetGinTraceIDFromPostForm 从 gin 的 postform 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromPostForm(c *gin.Context) string {
	return c.PostForm(string(TraceIDKeyname))
}
