package main

import (
	"context"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/getsentry/sentry-go"

	"go.uber.org/zap"
)

func main() {
	/* Error sentry dsn env */
	// 全局方法使用的默认logger在默认情况下不支持sentry上报，通过配置环境变量LoggingSentryDSN后自动支持
	logging.Error(nil, "dsn env")
	// 如果环境变量配置了sentry dsn，则会创建一个默认sentry client并初始化sentry，可以通过DefaultSentryClient 获取原始的sentry client
	if logging.DefaultSentryClient() != nil {
		// 如果已经初始化过sentry，则可以使用sentry hub直接上报数据到sentry
		sentry.CaptureMessage("hello sentry hub msg!")
		// waiting for report complete
		time.Sleep(2 * time.Second)
	}

	/* zap Debug */
	logging.Debug(nil, "Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991006","logger":"root.ctxLogger","msg":"Debug message","pid":45713,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	/* zap sugared logger Debug */
	logging.Debugs(nil, "Debugs message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991239","logger":"root.ctxLogger","msg":"Debugs message123 false [1 2 3] map[i:1 s:s]","pid":45713}

	/* zap sugared logger Debugf */
	logging.Debugf(nil, "Debugf message, %s", "ok")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991268","logger":"root.ctxLogger","msg":"Debugf message, ok","pid":45713}

	/* zap sugared logger Debugw */
	logging.Debugw(nil, "Debugw message", "name", "axiaoxin", "age", 18)
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991277","logger":"root.ctxLogger","msg":"Debugw message","pid":45713,"name":"axiaoxin","age":18}

	/* with context */
	c := logging.Context(context.Background(), logging.DefaultLogger(), "trace-id-123")
	logging.Debug(c, "Debug with trace id")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991314","logger":"root","msg":"Debug with trace id","pid":45713,"traceID":"trace-id-123"}

	/* extra fields */
	logging.Debug(c, "extra fields demo", logging.ExtraField("k1", "v1", "k2", 2, "k3", true))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991348","logger":"root","msg":"extra fields demo","pid":45713,"traceID":"trace-id-123","extra":{"k1":"v1","k2":2,"k3":true}}
}
