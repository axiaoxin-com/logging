# logging

logging 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的方法，除 zap 的 DPanic 、 DPanicf 方法外所有日志打印方法开箱即用
- 提供多种快速创建 logger 的方法
- 集成 **Sentry**，设置 DSN 后可直接使用 Sentry ，支持在使用 Error 及其以上级别打印日志时自动将该事件上报到 Sentry
- 支持从 context.Context/gin.Context 中创建、获取带有 **Trace ID** 的 logger
- 提供 gin 中 Trace ID 的中间件，支持自定义方法获取 Trace ID
- 支持服务内部函数方式和外部 HTTP 方式**动态调整日志级别**，无需修改配置、重启服务
- 支持自定义 logger Encoder 配置
- 支持将日志保存到文件并自动 rotate
- 支持 Gorm 日志打印 Trace ID

logging 只提供 zap 使用时的常用方法汇总，不是对 zap 进行二次开发，拒绝过度封装。

## 安装

```
go get -u github.com/axiaoxin-com/logging
```

## 开箱即用

logging 提供的开箱即用方法都是使用自身默认 logger 克隆出的 CtxLogger 实际执行的。
在 logging 被 import 时，会生成内部使用的默认 logger 。
默认 logger 使用 JSON 格式打印日志内容到 stderr 。
默认不带 Sentry 上报功能，可以通过设置环境变量或者替换 logger 方法支持。
默认 logger 可通过代码内部动态修改日志级别， 默认不支持 HTTP 方式动态修改日志级别，需要指定端口创建新的 logger 来支持。
默认带有初始字段 pid 打印进程 ID 。

开箱即用的方法第一个参数为 context.Context, 可以传入 gin.Context ，会尝试从其中获取 Trace ID 进行日志打印，无需 Trace ID 可以直接传 nil

**示例**

```golang
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
	// 全局方法使用的默认 logger 在默认情况下不支持 sentry 上报，通过配置环境变量 SENTRY_DSN 后自动支持
	logging.Error(nil, "dsn env")

	// 如果环境变量配置了 sentry dsn ，则会创建一个默认 sentry client 并初始化 sentry ，可以通过 DefaultSentryClient 获取原始的 sentry client
	if logging.DefaultSentryClient() != nil {
		// 如果已经初始化过 sentry ，则可以使用 sentry hub 直接上报数据到 sentry
		sentry.CaptureMessage("hello sentry hub msg!")
		sentry.Flush(2 * time.Second)
	}

	// 配置了 sentry 后，可以通过全局的方法上报 sentry
	// 封装了上面包含 Flush 方法的示例写法
	logging.SentryCaptureMessage("Hello sentry")

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
```

全局开箱即用的方法默认不支持 sentry 自动上报 Error 级别的事件，有两种方式可以使其支持：

1. 通过设置系统环境变量 `SENTRY_DSN` 和 `SENTRY_DEBUG` 来实现自动上报。

2. 也可以通过替换默认 logger 来实现让全局方法支持 Error 以上级别自动上报，以下示例：

```golang
// 默认的 logging 全局开箱即用的方法（如： logging.Debug , logging.Debugf 等）都是使用默认 logger 执行的，
// 默认 logger 不支持 Sentry 和输出日志到文件，可以通过创建一个新的 logger ，
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
	sentryClient, _ := logging.GetSentryClientByDSN(os.Getenv("dsn"), true)
	options := logging.Options{
		Name:           "replacedLogger",
		OutputPaths:    []string{"stderr", "lumberjack:"},
		LumberjackSink: logging.NewLumberjackSink("lumberjack", "/tmp/replace.log", 1, 1, 10, true, true),
		SentryClient:   sentryClient,
	}
	logger, _ := logging.NewLogger(options)
	// 替换默认 logger
	resetLogger := logging.ReplaceDefaultLogger(logger)

	// 全局方法将使用新的 logger ，上报 sentry 并输出到文件
	logging.Error(nil, "ReplaceDefaultLogger")
	// Output 并保存到文件:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661927","logger":"replacedLogger.ctxLogger","caller":"logging/global.go:Error:166","msg":"ReplaceDefaultLogger","pid":73847,"stacktrace":"github.com/axiaoxin-com/logging.Error\n\t/Users/ashin/go/src/logging/global.go:166\nmain.main\n\t/Users/ashin/go/src/logging/example/replace.go:30\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 重置默认 logger
	resetLogger()

	// 全局方法将恢复使用原始的 logger ，不再上报 sentry 和输出到文件
	logging.Error(nil, "ResetDefaultLogger")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.742995","logger":"root.ctxLogger","msg":"ResetDefaultLogger","pid":73847}
}
```

## 快速获取、创建你的 Logger

logging 提供多种方式快速获取一个 logger 来打印日志

**示例**：

```golang
package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	/* 获取默认 logger */
	defaultLogger := logging.DefaultLogger()
	defaultLogger.Debug("DefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548141","logger":"root","msg":"DefaultLogger","pid":68701}

	/* 为默认 logger 设置 sentry core */
	// logging 内部默认的 logger 不支持 sentry 上报，可以通过以下方法设置 sentry
	// 创建 sentry 客户端
	sentryClient, _ := logging.GetSentryClientByDSN("YOUR_SENTRY_DSN", false)
	// 设置 sentry ，使用该 logger 打印 Error 及其以上级别的日志事件将会自动上报到 Sentry
	defaultLogger = logging.SentryAttach(defaultLogger, sentryClient)

	/* 克隆一个带有初始字段的默认 logger */
	// 初始字段可以不传，克隆的 logger 名称会是 root.subname ，该 logger 打印的日志都会带上传入的字段
	cloneDefaultLogger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	cloneDefaultLogger.Debug("CloneDefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548271","logger":"root.subname","msg":"CloneDefaultLogger","pid":68701,"str_field":"field_value"}

	/* 使用 Options 创建 logger */
	// 可以直接使用空 Options 创建默认配置项的 logger
	// 不支持 sentry 和 http 动态修改日志级别，日志输出到 stderr
	emptyOptionsLogger, _ := logging.NewLogger(logging.Options{})
	emptyOptionsLogger.Debug("emptyOptionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548323","logger":"root","caller":"example/logger.go:main:48","msg":"emptyOptionsLogger","pid":68701}

	// 配置 Options 创建 logger
	// 日志级别定义在外层，便于代码内部可以动态修改日志级别
	level := logging.TextLevelMap["debug"]
	options := logging.Options{
		Name:              "root",                         // logger 名称
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
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548363","logger":"root","caller":"example/logger.go:main:67","msg":"optionsLogger","pid":68701}

	/* 从 context.Context 或*gin.Context 中获取或创建 logger */
	ctx := context.Background()
	ctxLogger := logging.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxLogger.Debug("ctxLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548414","logger":"root.ctxLogger","msg":"ctxLogger","pid":68701,"field1":"xxx"}
}
```

## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。


**示例 1**: 普通函数中打印打印带 Trace ID 的日志

```golang
package main

import (
	"context"

	"github.com/axiaoxin-com/logging"
)

/* context logger with trace id */
func main() {
	// 初始化一个 context
	ctx := context.Background()
	// 生成一个 trace id ，如果 context 是 gin.Context ，会尝试从其中获取，否则尝试从 context.Context 获取，获取不到则新生成
	traceID := logging.CtxTraceID(ctx)
	// 设置 trace id 到 context 和 logger 中， 会尝试同时设置到 gin.Context 中
	ctx = logging.Context(ctx, logging.CtxLogger(ctx), traceID)
	// 从 context 中获取 logger ，会尝试从 gin.Context 中获取， context 中没有 logger 则克隆默认 logger 作为 context logger
	ctxlogger := logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 19:12:45.263227","logger":"root.ctxLogger","msg":"ctxlogger with trace id debug","pid":17044,"traceID":"logging-bqbeobbipt345502logg"}

	logging.Debug(ctx, "global debug with ctx")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 19:12:45.263333","logger":"root.ctxLogger","msg":"global debug with ctx","pid":17044,"traceID":"logging-bqbeobbipt345502logg"}
}
```

**示例 2**: gin 中打印带 Trace ID 的日志

```golang
package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	// 使用 CtxLogger 打印带 trace id 的日志
	logging.CtxLogger(c).Info("func1 begin")
	func2(c)
	// 使用 logging 全局方法打印带 trace id 的日志
	logging.Info(c, "func1 end")
}

func func2(c context.Context) {
	logging.CtxLogger(c).Info("func2 begin")
	func3(c)
	logging.Info(c, "func2 end")
}

func func3(c context.Context) {
	logging.CtxLogger(c).Info("in func3")
}

func main() {
	r := gin.Default()

	// 使用中间件注册获取 trace id
	// 使用默认的回调方法从 Header 中获取 Key 为 traceID 的值作为 trace id
	// 可以自定义方法
	r.Use(logging.GinTraceIDMiddleware(logging.GetTraceIDFromHeader))

	r.GET("/ping", func(c *gin.Context) {
		logging.Error(c, "ping ping pong pong")
		// 模拟内部函数调用中打日志
		func1(c)
		c.String(200, "pong")
	})

	r.Run(":8080")
}

/*
日志输出

{"level":"ERROR","time":"2020-04-15 19:16:55.739465","logger":"root.ctxLogger","msg":"ping ping pong pong","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}
{"level":"INFO","time":"2020-04-15 19:16:55.739504","logger":"root.ctxLogger","msg":"func1 begin","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}
{"level":"INFO","time":"2020-04-15 19:16:55.739510","logger":"root.ctxLogger","msg":"func2 begin","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}
{"level":"INFO","time":"2020-04-15 19:16:55.739530","logger":"root.ctxLogger","msg":"in func3","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}
{"level":"INFO","time":"2020-04-15 19:16:55.739534","logger":"root.ctxLogger","msg":"func2 end","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}
{"level":"INFO","time":"2020-04-15 19:16:55.739540","logger":"root.ctxLogger","msg":"func1 end","pid":34425,"traceID":"logging-bqbeq9ript38cuae9nb0"}

请求响应头中也包含 Trace ID, 请求时如果指定 Header `-H "traceID: x-y-z"`， demo 将使用该值作为 trace id

curl
curl localhost:8080/ping -v
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /ping HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Traceid: logging-bqbeq9ript38cuae9nb0
< Content-Length: 4
<
* Connection #0 to host localhost left intact
pong* Closing connection 0
*/
```

## 动态修改 logger 日志级别

logging 可以在代码中对 AtomicLevel 调用 SetLevel 动态修改日志级别，也可以通过请求 HTTP 接口修改。
默认 logger 使用 `:1903` 运行 HTTP 服务来接收请求修改日志级别。实际使用中日志级别通常写在配置文件中，
可以通过监听配置文件的修改来动态调用 SetLevel 方法。

**示例**:

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

// level 全局变量，便于动态修改，初始化为 Debug 级别
var level zap.AtomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)

func main() {
	/* change log level on fly */

	// 创建指定 Level 的 logger ，并开启 http 服务
	options := logging.Options{
		Level:           level,
		AtomicLevelAddr: ":2012",
	}
	logger, _ := logging.NewLogger(options)
	logger.Debug("Debug level msg", zap.Any("current level", level.Level()))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:03:17.799767","logger":"root","caller":"example/atomiclevel.go:main:26","msg":"Debug level msg","pid":6088,"current level":"debug"}

	// 使用 SetLevel 动态修改 logger 日志级别为 error
	// 实际应用中可以监听配置文件中日志级别配置项的变化动态调用该函数
	level.SetLevel(zap.ErrorLevel)
	// Info 级别将不会被打印
	logger.Info("Info level msg will not be logged")
	// 只会打印 error 以上
	logger.Error("Error level msg", zap.Any("current level", level.Level()))
	// Output:
	// {"level":"ERROR","time":"2020-04-15 18:03:17.799999","logger":"root","caller":"example/atomiclevel.go:main:34","msg":"Error level msg","pid":6088,"current level":"error","stacktrace":"main.main\n\t/Users/ashin/go/src/logging/example/atomiclevel.go:34\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 通过 HTTP 方式动态修改当前的 error level 为 debug level
	// 查询当前 level
	url := "http://localhost" + options.AtomicLevelAddr
	resp, _ := http.Get(url)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("currentlevel:", string(content))
	// Output: currentlevel: {"level":"error"}

	logger.Info("Info level will not be logged")

	// 修改 level 为 debug
	c := &http.Client{}
	req, _ := http.NewRequest("PUT", url, strings.NewReader(`{"level": "debug"}`))
	resp, _ = c.Do(req)
	content, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("newlevel:", string(content))
	// Output: newlevel: {"level":"debug"}

	logger.Debug("level is changed on fly!")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:03:17.805293","logger":"root","caller":"example/atomiclevel.go:main:57","msg":"level is changed on fly!","pid":6088}

	/* 修改默认 logger 日志级别 */
	logging.Info(nil, "default logger level")
	// 修改前 Output:
	// {"level":"INFO","time":"2020-04-16 13:33:50.178265","logger":"root.ctxLogger","msg":"default logger level","pid":45311}

	// 获取默认 logger 的 level
	defaultLoggerLevel := logging.DefaultLoggerLevel()
	// 修改 level 为 error
	defaultLoggerLevel.SetLevel(zap.ErrorLevel)

	// info 将不会打印
	logging.Info(nil, "info level will not be print")
	logging.Error(nil, "new level")
	// Output:
	// {"level":"ERROR","time":"2020-04-16 13:33:50.178273","logger":"root.ctxLogger","msg":"new level","pid":45311}

}
```

## 自定义 logger Encoder 配置

**示例**：

```golang
package main

import (
	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	/* custom logger encoder */
	options := logging.Options{
		Name: "apiserver",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Logger",
			CallerKey:      "Caller",
			MessageKey:     "Message",
			StacktraceKey:  "Stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     logging.TimeEncoder, // 使用 logging 的 time 格式
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   logging.CallerEncoder, // 使用 logging 的 caller 格式
		},
		DisableCaller: false,
	}
	logger, _ := logging.NewLogger(options)
	logger.Debug("EncoderConfig Debug", zap.Reflect("Tags", map[string]interface{}{
		"Status":     "200 OK",
		"StatusCode": 200,
		"Latency":    0.075,
	}))
	// Output:
	// {"Level":"DEBUG","Time":"2020-04-15 19:23:44.373302","Logger":"apiserver","Caller":"example/encoder.go:main:30","Message":"EncoderConfig Debug","pid":66937,"Tags":{"Latency":0.075,"Status":"200 OK","StatusCode":200}}
}
```


## 日志保存到文件并自动 rotate

使用 lumberjack 将日志保存到文件并 rotate ，采用 zap 的 RegisterSink 方法和 Config.OutputPaths 字段添加自定义的日志输出的方式来使用 lumberjack 。


**示例**

```golang
package main

import (
	"github.com/axiaoxin-com/logging"
)

// Options 传入 LumberjacSink ，并在 OutputPaths 中添加对应 scheme 就能将日志保存到文件并自动 rotate
func main() {
	// scheme 为 lumberjack ，日志文件为 /tmp/x.log , 保存 7 天，保留 10 份文件，文件大小超过 100M ，使用压缩备份，压缩文件名使用 localtime
	sink := logging.NewLumberjackSink("lumberjack", "/tmp/x.log", 7, 10, 100, true, true)
	options := logging.Options{
		LumberjackSink: sink,
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack:"},
	}
	logger, _ := logging.NewLogger(options)
	logger.Debug("xxx")

	sink2 := logging.NewLumberjackSink("lumberjack2", "/tmp/x2.log", 7, 10, 100, true, true)
	options2 := logging.Options{
		LumberjackSink: sink2,
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack2:"},
	}
	logger2, _ := logging.NewLogger(options2)
	logger2.Debug("yyy")
}
```

## 支持 Gorm 日志打印 Trace ID


在每一次使用 gorm 进行 db 操作前，调用 GormDBWithCtxLogger 来设置替换 gorm DB 对象的默认 logger 并生成新的 DB 对象，之后使用新的 DB 对象来操作 gorm 即可。

示例：

```golang
package main

import (
	"context"
	"os"
	"sync"

	"github.com/axiaoxin-com/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Product test model
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var (
	db  *gorm.DB
	err error
	wg  sync.WaitGroup
)

func init() {
	// Create gorm db instance
	db, err = gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})
	// Enable Logger, show detailed log
	db.LogMode(true)

}

// G 模拟一次请求处理
func G(traceID string) {

	// Mock a context with a trace id and logger
	ctx := logging.Context(context.Background(), logging.DefaultLogger(), traceID)

	// 打印带 trace id 的 gorm 日志
	// 必须先对 db 对象设置带有 trace id 的 ctxlogger 作为 sql 日志打印的 logger
	// 后续的 gorm 操作使用新的 db 对象即可
	db := logging.GormDBWithCtxLogger(ctx, db)
	// Create
	db.Create(&Product{Code: traceID, Price: 1000})
	wg.Done()
}

func main() {
	// defer clear
	defer db.Close()
	defer os.Remove("./sqlite3.db")

	// 模拟并发
	wg.Add(4)
	go G("g1")
	go G("g2")
	go G("g3")
	go G("g4")
	wg.Wait()
}

// log:
// {"level":"DEBUG","time":"2020-04-21 17:08:44.449254","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g4","vars":["2020-04-21T17:08:44.448622+08:00","2020-04-21T17:08:44.448622+08:00",null,"g4",1000],"rowsAffected":1,"duration":0.000613636}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.452657","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g2","vars":["2020-04-21T17:08:44.44919+08:00","2020-04-21T17:08:44.44919+08:00",null,"g2",1000],"rowsAffected":1,"duration":0.0034358}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.458721","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g1","vars":["2020-04-21T17:08:44.44946+08:00","2020-04-21T17:08:44.44946+08:00",null,"g1",1000],"rowsAffected":1,"duration":0.009227084}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.471094","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g3","vars":["2020-04-21T17:08:44.449387+08:00","2020-04-21T17:08:44.449387+08:00",null,"g3",1000],"rowsAffected":1,"duration":0.021678226}
```
