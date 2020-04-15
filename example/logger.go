package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	/* 获取默认logger */
	defaultLogger := logging.DefaultLogger()
	defaultLogger.Debug("DefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548141","logger":"root","msg":"DefaultLogger","pid":68701}

	/* 为默认logger设置sentry core */
	// logging 内部默认的 logger 不支持 sentry 上报，可以通过以下方法设置 sentry
	// 创建 sentry 客户端
	sentryClient, _ := logging.GetSentryClientByDSN("YOUR_SENTRY_DSN", false)
	// 设置 sentry，使用该 logger 打印 Error 及其以上级别的日志事件将会自动上报到 Sentry
	defaultLogger = logging.SentryAttach(defaultLogger, sentryClient)

	/* 获取默认sugared logger */
	defaultSLogger := logging.DefaultSLogger()
	defaultSLogger.Debug("DefaultSLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548263","logger":"root","msg":"DefaultSLogger","pid":68701}

	/* 克隆一个带有初始字段的默认logger */
	// 初始字段可以不传，克隆的 logger 名称会是 root.subname，该 logger 打印的日志都会带上传入的字段
	cloneDefaultLogger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	cloneDefaultLogger.Debug("CloneDefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548271","logger":"root.subname","msg":"CloneDefaultLogger","pid":68701,"str_field":"field_value"}

	/* 克隆一个带有初始字段的默认 sugared logger */
	cloneDefaultSLogger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))
	cloneDefaultSLogger.Debug("CloneDefaultSLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548283","logger":"root.subname","msg":"CloneDefaultSLogger","pid":68701,"foo":123,"str_field":"field_value"}

	/* 使用Options创建logger */
	// 可以直接使用空Options创建默认配置项的logger
	// 不支持 sentry 和 http 动态修改日志级别，日志输出到stderr
	emptyOptionsLogger, _ := logging.NewLogger(logging.Options{})
	emptyOptionsLogger.Debug("emptyOptionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548323","logger":"root","caller":"example/logger.go:main:48","msg":"emptyOptionsLogger","pid":68701}

	// 配置Options创建logger
	// 日志级别定义在外层，便于代码内部可以动态修改日志级别
	level := logging.TextLevelMap["debug"]
	options := logging.Options{
		Name:              "root",                         // logger 名称
		Level:             level,                          // zap 的 AtomicLevel，logger 日志级别
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
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548363","logger":"root","caller":"example/logger.go:main:67","msg":"optionsLogger","pid":68701}

	/* 从context.Context或*gin.Context中获取或创建logger */
	ctx := context.Background()
	ctxLogger := logging.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxLogger.Debug("ctxLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548414","logger":"root.ctxLogger","msg":"ctxLogger","pid":68701,"field1":"xxx"}
}
