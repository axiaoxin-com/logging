package main

import (
	"context"
	"github/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	// zap Debug
	logging.Debug("Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.32688+08:00","logger":"root","caller":"logging/global.go:120","msg":"Debug message","pid":27907,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	// zap sugared logger Debug
	logging.Debugs("SDebug message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.327239+08:00","logger":"root","caller":"logging/global.go:10","msg":"SDebug message123 false [1 2 3] map[i:1 s:s]","pid":27907}

	// zap sugared logger Debugf
	logging.Debugf("SDebugf message, %s", "ok")
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327287+08:00","logger":"root","caller":"logging/global.go:47","msg":"SDebugf message, ok","pid":27907}

	// zap sugared logger Debugw
	logging.Debugw("SDebug message", "name", "axiaoxin", "age", 18)
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327301+08:00","logger":"root","caller":"logging/global.go:84","msg":"SDebug message","pid":27907,"name":"axiaoxin","age":18}

	// clone default logger
	logger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	logger.Debug("CloneDefaultLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:20:36.614438+08:00","logger":"root.subname","caller":"example/logging.go:27","msg":"CloneDefaultLogger Debug","pid":54273,"str_field":"field_value"}

	// clone default sugared logger
	slogger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))
	slogger.Debug("CloneDefaultSLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:24:41.629175+08:00","logger":"root.subname","caller":"example/logging.go:32","msg":"CloneDefaultSLogger Debug","pid":73087,"foo":123,"str_field":"field_value"}

	// new context logger
	ctx := context.Background()
	ctxlogger := logging.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxlogger.Debug("ctxlogger debug")

	// context logger with request id
	ctx = context.Background()
	// get a trace id
	traceID := logging.CtxTraceID(ctx)
	// set request id in context
	ctx = logging.Context(ctx, traceID)
	// get a logger with trace id
	ctxlogger = logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
}
