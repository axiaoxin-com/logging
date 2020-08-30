package main

import (
	"fmt"
	"os"

	"github.com/axiaoxin-com/logging"
	"github.com/gin-gonic/gin"
)

func main() {

	// set sentry dsn, error level's log will report to sentry automatically
	os.Setenv(logging.SentryDSNEnvKey, "http://your-sentry-dsn")
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	// you can custom the config or use logging.GinLogger() by default config
	conf := logging.GinLoggerConfig{
		Formatter: func(m logging.GinLogMsg) string {
			return fmt.Sprintf("%s use %s request %s, handler %s use %f seconds to respond it with %d at %v",
				m.ClientIP, m.Method, m.RequestURI, m.HandlerName, m.Latency, m.StatusCode, m.Timestamp)
		},
		SkipPaths:              []string{},
		DisableDetails:         false,
		DetailsWithContextKeys: true,
		DetailsWithBody:        true,
		TraceIDFunc:            func(c *gin.Context) string { return "fake-uuid" },
	}
	app.Use(logging.GinLoggerWithConfig(conf))
	app.POST("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	app.Run(":8888")
}
