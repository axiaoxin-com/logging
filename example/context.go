package main

import (
	"context"

	"github.com/axiaoxin-com/logging"
)

/* context logger with trace id */
func main() {
	// 初始化一个 context
	ctx := context.Background()
	// 从 context 中获取 logger ，会尝试从 gin.Context 中获取， context 中没有 logger 则克隆默认 logger 作为 context logger
	// context 中无 trace id 会默认生成一个新的 trace id
	ctxlogger := logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
	// Output:
	// {"level":"DEBUG","time":"2020-06-10 20:30:48.588416","logger":"root.ctxLogger","msg":"ctxlogger with trace id debug","pid":3242,"traceID":"logging-brgd4u3ipt30pamqff80"}

	// 设置 一个指定的 trace id 和 logger 到 context 中， 会尝试同时设置到 gin.Context 中
	traceID := "this-is-a-trace-id"
	ctx = logging.Context(ctx, logging.DefaultLogger(), traceID)
	logging.Debug(ctx, "global debug with ctx")
	// Output:
	// {"level":"DEBUG","time":"2020-06-10 20:30:48.588510","logger":"root","msg":"global debug with ctx","pid":3242,"traceID":"this-is-a-trace-id"}

	ctxlogger2 := logging.CtxLogger(ctx)
	ctxlogger2.Debug("ctxlogger2 with special trace id")
	// Output:
	// {"level":"DEBUG","time":"2020-06-10 20:30:48.588521","logger":"root","msg":"ctxlogger2 with special trace id","pid":3242,"traceID":"this-is-a-trace-id"}
}
