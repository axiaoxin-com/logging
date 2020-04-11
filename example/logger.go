package main

import (
	"github/axiaoxin-com/logging"
)

func main() {
	// get a sentry client as core
	dsn := "sentrydsn"
	sc, _ := logging.GetSentryClientByDSN(dsn, true)

	// set logger options
	options := logging.Options{
		Name:              logging.DefaultLoggerName,
		Level:             "debug",
		Format:            "json",
		OutputPaths:       []string{"stderr"},
		InitialFields:     logging.DefaultInitialFields,
		DisableCaller:     false,
		DisableStacktrace: false,
		SentryClient:      sc,
	}

	// new logger
	logger, _ := logging.NewLogger(options)
	logger.Debug("Such as logging.Debug")
	// {"level":"DEBUG","time":"2020-04-12T03:10:09.220667+08:00","logger":"root","caller":"example/logger.go:26","msg":"Such as logging.Debug","pid":69775}
	logger.Error("Such as logging.Error")
	// {"level":"ERROR","time":"2020-04-12T03:10:09.220913+08:00","logger":"root","caller":"example/logger.go:27","msg":"Such as logging.Error","pid":69775,"stacktrace":"main.main\n\t/Users/ashin/go/src/logging/example/logger.go:27\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// sugared logger
	sugaredLogger := logger.Sugar()
	sugaredLogger.Debug("Such as logging.SDebug")
	// {"level":"DEBUG","time":"2020-04-12T03:10:09.220968+08:00","logger":"root","caller":"example/logger.go:31","msg":"Such as logging.SDebug","pid":69775}

	// clone logging default Logger
	clonedLogger := logging.CloneLogger("subname")
	clonedLogger.Debug("I have a new logger name")
	// {"level":"DEBUG","time":"2020-04-12T03:10:09.220994+08:00","logger":"root.subname","caller":"example/logger.go:35","msg":"I have a new logger name","pid":69775}
}
