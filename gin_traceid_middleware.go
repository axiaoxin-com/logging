package logging

import "github.com/gin-gonic/gin"

// GinTraceIDMiddleware is a gin middleware for gen a trace id in context
func GinTraceIDMiddleware(traceIDKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 header 获取 trace id，不存在则新生成
		traceID := c.Request.Header.Get(traceIDKey)
		if traceID == "" {
			traceID = CtxTraceID(c)
		}
		// 设置 trace id 到 header 中
		c.Writer.Header().Set(traceIDKey, traceID)
		// 设置 trace id 和 ctxLogger 和 context 中
		Context(c, traceID)

		c.Next()
	}
}
