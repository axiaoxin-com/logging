package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	/* 克隆一个带有初始字段的默认 logger */
	// 初始字段可以不传，克隆的 logger 名称会是 logging.subname ，该 logger 打印的日志都会带上传入的字段
	cloneDefaultLogger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	cloneDefaultLogger.Debug("CloneDefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548271","logger":"logging.subname","msg":"CloneDefaultLogger","pid":68701,"str_field":"field_value"}

	/* 为 clone logger 设置 sentry core */
	// logging 内部默认的 logger 不支持 sentry 上报，可以通过以下方法设置 sentry
	// 创建 sentry 客户端
	sentryClient, _ := logging.NewSentryClient("YOUR_SENTRY_DSN", false)
	// 设置 sentry ，使用该 logger 打印 Error 及其以上级别的日志事件将会自动上报到 Sentry
	cloneDefaultLogger = logging.SentryAttach(cloneDefaultLogger, sentryClient)

	/* 使用 Options 创建 logger */
	// 可以直接使用空 Options 创建默认配置项的 logger
	// 不支持 sentry 和 http 动态修改日志级别，日志输出到 stderr
	emptyOptionsLogger, _ := logging.NewLogger(logging.Options{})
	emptyOptionsLogger.Debug("emptyOptionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548323","logger":"logging","caller":"example/logger.go:main:48","msg":"emptyOptionsLogger","pid":68701}

	// 配置 Options 创建 logger
	// 日志级别定义在外层，便于代码内部可以动态修改日志级别
	level := logging.TextLevelMap["debug"]
	options := logging.Options{
		Name:              "logging",                      // logger 名称
		Level:             level,                          // zap 的 AtomicLevel ， logger 日志级别
		Format:            "json",                         // 日志输出格式为 json
		OutputPaths:       []string{"stderr"},             // 日志输出位置为 stderr
		InitialFields:     logging.DefaultInitialFields(), // DefaultInitialFields 初始 logger 带有 pid 字段
		DisableCaller:     false,                          // 是否打印调用的代码行位置
		DisableStacktrace: false,                          // 错误日志是否打印调用栈信息
		SentryClient:      sentryClient,                   // sentry 客户端
		AtomicLevelAddr:   ":8080",                        // http 动态修改日志级别的端口地址，不设置则不开启 http 服务
	}
	optionsLogger, _ := logging.NewLogger(options)
	optionsLogger.Debug("optionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548363","logger":"logging","caller":"example/logger.go:main:67","msg":"optionsLogger","pid":68701}

	/* 从 context.Context 或*gin.Context 中获取或创建 logger */
	ctx := context.Background()
	ctxLogger := logging.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxLogger.Debug("ctxLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548414","logger":"logging.ctx_logger","msg":"ctxLogger","pid":68701,"field1":"xxx"}
}
