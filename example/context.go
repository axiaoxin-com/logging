package main

import (
	"context"

	"github.com/axiaoxin-com/logging"
)

/* context logger with trace id */
func main() {
	// 初始化一个 context
	ctx := context.Background()
	// 生成一个 trace id，如果 context 是 gin.Context，会尝试从其中获取，否则尝试从 context.Context 获取，获取不到则新生成
	traceID := logging.CtxTraceID(ctx)
	// 设置 trace id 到 context 和 logger 中， 会尝试同时设置到 gin.Context 中
	ctx = logging.Context(ctx, logging.CtxLogger(ctx), traceID)
	// 从 context 中获取 logger，会尝试从 gin.Context 中获取，context 中没有 logger 则克隆默认 logger 作为 context logger
	ctxlogger := logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 19:12:45.263227","logger":"root.ctxLogger","msg":"ctxlogger with trace id debug","pid":17044,"traceID":"logging-bqbeobbipt345502logg"}

	logging.Debug(ctx, "global debug with ctx")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 19:12:45.263333","logger":"root.ctxLogger","msg":"global debug with ctx","pid":17044,"traceID":"logging-bqbeobbipt345502logg"}
}
