// 默认的 logging 全局开箱即用的方法（如： logging.Debug , logging.Debugf 等）都是使用默认 logger 执行的，
// 默认 logger 不支持 Sentry 和输出日志到文件，可以通过创建一个新的 logger，
// 再使用 ReplaceDefaultLogger 方法替换默认 logger 为新的 logger 来解决。

package main

import (
	"os"

	"github.com/axiaoxin-com/logging"
)

func main() {
	// 默认使用全局方法不会保存到文件和上报 Sentry
	logging.Error(nil, "default logger no sentry and file")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661457","logger":"root.ctxLogger","msg":"default logger no sentry and file","pid":73847}

	// 创建一个支持 sentry 和 lumberjack 的 logger
	sentryClient, _ := logging.NewSentryClient(os.Getenv("dsn"), true)
	options := logging.Options{
		Name:           "replacedLogger",
		OutputPaths:    []string{"stderr", "lumberjack:"},
		LumberjackSink: logging.NewLumberjackSink("lumberjack", "/tmp/replace.log", 1, 1, 10, true, true),
		SentryClient:   sentryClient,
	}
	logger, _ := logging.NewLogger(options)
	// 替换默认 logger
	resetLogger := logging.ReplaceDefaultLogger(logger)

	// 全局方法将使用新的 logger，上报 sentry 并输出到文件
	logging.Error(nil, "ReplaceDefaultLogger")
	// Output并保存到文件:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661927","logger":"replacedLogger.ctxLogger","caller":"logging/global.go:Error:166","msg":"ReplaceDefaultLogger","pid":73847,"stacktrace":"github.com/axiaoxin-com/logging.Error\n\t/Users/ashin/go/src/logging/global.go:166\nmain.main\n\t/Users/ashin/go/src/logging/example/replace.go:30\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 重置默认 logger
	resetLogger()

	// 全局方法将恢复使用原始的 logger，不再上报 sentry 和输出到文件
	logging.Error(nil, "ResetDefaultLogger")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.742995","logger":"root.ctxLogger","msg":"ResetDefaultLogger","pid":73847}
}
