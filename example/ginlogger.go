package main

import (
	"context"
	"fmt"

	"github.com/axiaoxin-com/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	// you can custom the config or use logging.GinLogger() by default config
	conf := logging.GinLoggerConfig{
		Formatter: func(m logging.GinLogDetails) string {
			return fmt.Sprintf("%s use %s request %s, handler %s use %f seconds to respond it with %d at %v",
				m.ClientIP, m.Method, m.RequestURI, m.HandlerName, m.Latency, m.StatusCode, m.Timestamp)
		},
		SkipPaths:     []string{},
		EnableDetails: false,
		TraceIDFunc:   func(context.Context) string { return "fake-uuid" },
	}
	app.Use(logging.GinLoggerWithConfig(conf))
	app.POST("/ping", func(c *gin.Context) {
		// panic("xx")
		c.JSON(200, string(logging.GetGinRequestBody(c)))
	})
	app.Run(":8888")
}
