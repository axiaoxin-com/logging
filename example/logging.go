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
