package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	/* zap Debug */
	logging.Debug(nil, "Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.32688+08:00","logger":"root.ctxLogger","msg":"Debug message","pid":27907,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	/* zap sugared logger Debug */
	logging.Debugs(nil, "Debugs message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// {"level":"DEBUG","time":"2020-04-12T02:56:39.327239+08:00","logger":"root.ctxLogger","msg":"Debugs message123 false [1 2 3] map[i:1 s:s]","pid":27907}

	/* zap sugared logger Debugf */
	logging.Debugf(nil, "Debugf message, %s", "ok")
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327287+08:00","logger":"root.ctxLogger","msg":"Debugf message, ok","pid":27907}
	/* zap sugared logger Debugw */
	logging.Debugw(nil, "Debugw message", "name", "axiaoxin", "age", 18)
	//{"level":"DEBUG","time":"2020-04-12T02:56:39.327301+08:00","logger":"root.ctxLogger","msg":"Debugw message","pid":27907,"name":"axiaoxin","age":18}

	/* with context */
	c := logging.Context(context.Background(), logging.DefaultLogger(), "trace-id-123")
	logging.Debug(c, "Debug with trace id")
	// {"level":"DEBUG","time":"2020-04-14T16:16:29.404008+08:00","logger":"root","msg":"Debug with trace id","pid":44559,"traceID":"trace-id-123"}

	/* extra fields */
	logging.Debug(c, "extra fields demo", logging.ExtraField("k1", "v1", "k2", 2, "k3", true))
	// {"level":"DEBUG","time":"2020-04-14T23:50:05.056916+08:00","logger":"root","msg":"extra fields demo","pid":98214,"traceID":"trace-id-123","extra":{"k1":"v1","k2":2,"k3":true}}
}
