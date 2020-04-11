# logging

zap logger wrapped with sentry core


## install

```
go get -u github/axiaoxin-com/logging
```

## usage


### custom logger

```
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
```

### logging global funcs

```
package main

import (
	"github/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	// zap Debug
	logging.Debug("Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.32688+08:00","logger":"root","caller":"logging/global.go:120","msg":"Debug message","pid":27907,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	// zap sugared logger Debug
	logging.SDebug("SDebug message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.327239+08:00","logger":"root","caller":"logging/global.go:10","msg":"SDebug message123 false [1 2 3] map[i:1 s:s]","pid":27907}

	// zap sugared logger Debugf
	logging.SDebugf("SDebugf message, %s", "ok")
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327287+08:00","logger":"root","caller":"logging/global.go:47","msg":"SDebugf message, ok","pid":27907}

	// zap sugared logger Debugw
	logging.SDebugw("SDebug message", "name", "axiaoxin", "age", 18)
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327301+08:00","logger":"root","caller":"logging/global.go:84","msg":"SDebug message","pid":27907,"name":"axiaoxin","age":18}

}
```
