package main

import (
	"context"
	"github/axiaoxin-com/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	/* zap Debug */
	logging.Debug("Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.32688+08:00","logger":"root","msg":"Debug message","pid":27907,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	/* zap sugared logger Debug */
	logging.Debugs("SDebug message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.327239+08:00","logger":"root","msg":"SDebug message123 false [1 2 3] map[i:1 s:s]","pid":27907}

	/* zap sugared logger Debugf */
	logging.Debugf("SDebugf message, %s", "ok")
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327287+08:00","logger":"root","msg":"SDebugf message, ok","pid":27907}

	/* zap sugared logger Debugw */
	logging.Debugw("SDebug message", "name", "axiaoxin", "age", 18)
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327301+08:00","logger":"root","msg":"SDebug message","pid":27907,"name":"axiaoxin","age":18}

	/* clone default logger */
	logger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	logger.Debug("CloneDefaultLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:20:36.614438+08:00","logger":"root.subname","msg":"CloneDefaultLogger Debug","pid":54273,"str_field":"field_value"}

	/* clone default sugared logger */
	slogger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))
	slogger.Debug("CloneDefaultSLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:24:41.629175+08:00","logger":"root.subname","msg":"CloneDefaultSLogger Debug","pid":73087,"foo":123,"str_field":"field_value"}

	/* new context logger */
	ctx := context.Background()
	ctxlogger := logging.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxlogger.Debug("ctxlogger debug")
	// {"level":"DEBUG","time":"2020-04-13T14:52:29.00566+08:00","logger":"root.ctxLogger","msg":"ctxlogger debug","pid":53998,"field1":"xxx"}

	/* context logger with request id */
	ctx = context.Background()
	// get a trace id
	traceID := logging.CtxTraceID(ctx)
	// set request id in context
	ctx = logging.Context(ctx, traceID)
	// get a logger with trace id
	ctxlogger = logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
	// {"level":"DEBUG","time":"2020-04-13T14:52:29.005685+08:00","logger":"root.ctxLogger","msg":"ctxlogger with trace id debug","pid":53998,"traceID":"logging-bqa0obbipt3d5rj6gus0"}

	/* custom logger encoder key name */
	options := logging.Options{
		Name: "apiserver",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "LogTime",
			LevelKey:       "LogLevel",
			NameKey:        "ServiceName",
			CallerKey:      "LogLine",
			MessageKey:     "Message",
			StacktraceKey:  "Stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, _ = logging.NewLogger(options)
	logger.Debug("EncoderConfig Debug", zap.Reflect("Tags", map[string]interface{}{
		"Status":     "200 OK",
		"StatusCode": 200,
		"Latency":    0.075,
	}))
	// {"LogLevel":"DEBUG","LogTime":"2020-04-13T14:51:39.478605+08:00","ServiceName":"apiserver","LogLine":"example/logging.go:72","Message":"EncoderConfig Debug","pid":50014,"Tags":{"Latency":0.075,"Status":"200 OK","StatusCode":200}}
}
