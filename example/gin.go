package main

import (
	"context"
	"github/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	// 使用CtxLogger打印带trace id的日志
	logging.CtxLogger(c).Info("func1 will call func2")
	func2(c)
	// 使用logging全局方法打印带trace id的日志
	logging.Info(c, "func2 is called")
}

func func2(c context.Context) {
	logging.CtxLogger(c).Info("func2 will call func3")
	func3(c)
	logging.Info(c, "func3 is called")
}

func func3(c context.Context) {
	logging.CtxLogger(c).Info("func3 be called")
}

func main() {
	r := gin.Default()

	// 使用默认的回调方法从Header中获取Key为traceID的值作为trace id
	// 可以自定义方法
	r.Use(logging.GinTraceIDMiddleware(logging.GetTraceIDFromHeader))

	r.GET("/ping", func(c *gin.Context) {
		logging.CtxLogger(c).Info("ping ping pong pong")
		func1(c)
		c.String(200, "pong")
	})

	r.Run(":8080")
}
