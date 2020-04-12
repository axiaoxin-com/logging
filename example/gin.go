package main

import (
	"context"
	"github/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	logging.CtxLogger(c).Info("func1 will call func2")
	func2(c)
}

func func2(c context.Context) {
	logging.CtxLogger(c).Info("func2 will call func3")
	func3(c)
}

func func3(c context.Context) {
	logging.CtxLogger(c).Info("func3 be called")
}

func main() {
	r := gin.Default()

	r.Use(logging.GinTraceIDMiddleware(logging.TraceIDKey))

	r.GET("/ping", func(c *gin.Context) {
		logging.CtxLogger(c).Info("ping ping pong pong")
		func1(c)
		c.String(200, "pong")
	})

	r.Run(":8080")
}
