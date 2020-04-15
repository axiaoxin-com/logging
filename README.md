# logging

logging 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的方法，除 zap 的 DPanic、DPanicf 方法外所有日志打印方法开箱即用
- 提供多种快速创建 logger 的方法
- 支持在使用 Error 及其以上级别打印日志时自动将该事件上报到 **Sentry**
- 支持从 context.Context/gin.Context 中创建、获取带有 **Trace ID** 的 logger
- 提供 gin 中 Trace ID 的中间件，支持自定义方法获取 Trace ID
- 支持服务内部函数方式和外部 HTTP 方式**动态调整日志级别**，无需修改配置、重启服务
- 支持自定义 logger EncoderConfig 字段名
- 支持将日志保存到文件并自动 rotate

logging 只提供 zap 使用时的常用方法汇总，不是对 zap 进行二次开发，拒绝过度封装。

## 安装

```
go get -u github.com/axiaoxin-com/logging
```

## 开箱即用

logging 提供的开箱即用方法都是使用自身默认 logger 克隆出的 CtxLogger 实际执行的，
在 logging 被 import 时，会生成内部使用的默认 logger，
默认 logger 使用 JSON 格式打印日志内容到 stderr ，
不带 Sentry 上报功能，
可通过 HTTP 调用 `curl -XPUT "http://localhost:1903" -d '{"level": "info"}'` 动态修改日志级别，
默认带有初始字段 pid 打印进程 ID

开箱即用的方法第一个参数为 context.Context, 可以传入 gin.Context，会尝试从其中获取 Trace ID 进行日志打印，无需 Trace ID 可以直接传 nil

**示例**

```golang
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
```
## 快速创建你的 Logger

logging 提供多种方式快速获取一个 logger 来打印日志

**示例0**：快速获取默认 logger

```golang
import "github.com/axiaoxin-com/logging"

// 使用 logging 提供的方法获取默认 logger
logger := logging.DefaultLogger()
// 使用 logging 提供的方法获取默认 slogger
slogger := logging.DefaultSLogger()
```

**示例1**：为默认 logger 设置 sentry Error以上日志自动上报错误事件

```golang
import "github.com/axiaoxin-com/logging"

// 创建一个默认 logger
logger := logging.DefaultLogger()

// logging 内部默认的 logger 不支持 sentry 上报，可以通过以下方法设置 sentry
// 创建 sentry 客户端
sentryClient, _ := logging.GetSentryClientByDSN("YOUR_SENTRY_DSN", false)
// 设置 sentry，使用该 logger 打印 Error 及其以上级别的日志事件将会自动上报到 Sentry
logger = logging.SentryAttach(logger, sentryClient)
```

**示例2**: 使用 NewLogger 方法创建一个默认配置的 logger （不支持 sentry 和 http 动态修改日志级别，日志输出到stderr）

```golang
import "github.com/axiaoxin-com/logging"

logger, _ := logging.NewLogger(logging.Options{})
```

**示例3**: 创建有配置项的 logger （支持 sentry 和 http 动态修改日志级别）

```golang
import "github.com/axiaoxin-com/logging"

// sentry client for reporting Error to Sentry
sc, _ := logging.GetSentryClientByDSN("your_sentry_dsn", true)

// atomic level for changing it dynamicly
level := logging.TextLevelMap["debug"]

options := logging.Options{
    Name:              "your_logger_name",              // logger 名称
    Level:             level,                           // zap 的 AtomicLevel，logger 日志级别
    Format:            "json",                          // 日志输出格式为 json
    OutputPaths:       []string{"stderr"},              // 日志输出位置为 stderr
    InitialFields:     logging.DefaultInitialFields(),  // DefaultInitialFields 初始 logger 带有 pid 字段
    DisableCaller:     false,                           // 是否打印调用的代码行位置
    DisableStacktrace: false,                           // 错误日志是否打印调用栈信息
    SentryClient:      sc,                              // sentry 客户端
    AtomicLevelAddr:   ":8080",                         // http 动态修改日志级别的端口地址，不设置则不开启 http 服务
}

// new logger with options
logger, _ := logging.NewLogger(options)
```

**示例4**: 快速克隆一个默认 logger，并添加初始字段

```golang
import "github.com/axiaoxin-com/logging"

logger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))

logger.Debug("CloneDefaultLogger Debug")

// {"level":"DEBUG","time":"2020-04-13T00:20:36.614438+08:00","logger":"root.subname","msg":"CloneDefaultLogger Debug","pid":54273,"str_field":"field_value"}
```

初始字段可以不传，克隆的 logger 名称会是 root.subname，该 logger 打印的日志都会带上传入的字段

**示例5**: 快速克隆一个默认 sugared logger，并添加初始字段

```golang
import "github.com/axiaoxin-com/logging"

logger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))

logger.Debug("CloneDefaultSLogger Debug")

// {"level":"DEBUG","time":"2020-04-13T00:24:41.629175+08:00","logger":"root.subname","msg":"CloneDefaultSLogger Debug","pid":73087,"foo":123,"str_field":"field_value"}
```

初始字段可以不传，克隆的 sugared logger 名称会是 root.subname，添加的初始字段则该 logger 打印的日志都会带上传入的字段

**示例6**: 创建一个 CtxLogger

```golang
/* new context logger */
ctx := context.Background()
ctxlogger := logging.CtxLogger(ctx, zap.String("fie    ld1", "xxx"))
ctxlogger.Debug("ctxlogger debug")
// {"level":"DEBUG","time":"2020-04-13T14:52:29.00566+08:00","logger":"root.ctxLogger","msg":"ctxlogger debug","pid":53998,"field1":"xxx"}
 ```


## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。


**示例1**: 普通函数中打印打印带 Trace ID 的日志

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
	// 生成一个 trace id，如果 context 是 gin.Context，会尝试从其中获取，否则尝试从 context.Context 获取，获取不到则新生成
	traceID := logging.CtxTraceID(ctx)
	// 设置 trace id 到 context 中， 会尝试同时设置到 gin.Context 中
	ctx = logging.Context(ctx, logging.CtxLogger(ctx), traceID)
	// 从 context 中获取 logger，会尝试从 gin.Context 中获取，context 中没有 logger 则克隆默认 logger 作为 context logger
	ctxlogger := logging.CtxLogger(ctx)
	// log with trace id
	ctxlogger.Debug("ctxlogger with trace id debug")
	logging.Debug(ctx, "global debug with ctx")
	// Output:
	// {"level":"DEBUG","time":"2020-04-14T16:32:36.565279+08:00","logger":"root.ctxLogger","msg":"ctxlogger with trace id debug","pid":17930,"traceID":"logging-bqana93ipt34c2lc9lgg"}
	// {"level":"DEBUG","time":"2020-04-14T16:32:36.565394+08:00","logger":"root.ctxLogger","msg":"global debug with ctx","pid":17930,"traceID":"logging-bqana93ipt34c2lc9lgg"}
}
```

**示例2**: gin 中打印带 Trace ID 的日志

```golang
package main

import (
	"context"
	"github.com/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	// 使用CtxLogger打印带trace id的日志
	logging.CtxLogger(c).Info("func1 will call func2")
	func2(c)
	// 使用logging全局方法打印带trace id的日志
	logging.Info(c, "func2 is called")
}

func func2(c context.Context) {
	logging.CtxLogger(c).Info("func2 will call func3")
	func3(c)
	logging.Info(c, "func3 is called")
}

func func3(c context.Context) {
	logging.CtxLogger(c).Info("func3 be called")
}

func main() {
	r := gin.Default()

    // 使用默认的回调方法从Header中获取Key为traceID的值作为trace id
    // 可以自定义方法
    r.Use(logging.GinTraceIDMiddleware(logging.GetTraceIDFromHeader))

	r.GET("/ping", func(c *gin.Context) {
		logging.CtxLogger(c).Info("ping ping pong pong")
		func1(c)
		c.String(200, "pong")
	})

	r.Run(":8080")
}
```

请求日志：
```json
{"level":"INFO","time":"2020-04-14T16:35:36.151951+08:00","logger":"root.ctxLogger","msg":"ping ping pong pong","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
{"level":"INFO","time":"2020-04-14T16:35:36.15217+08:00","logger":"root.ctxLogger","msg":"func1 will call func2","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
{"level":"INFO","time":"2020-04-14T16:35:36.152178+08:00","logger":"root.ctxLogger","msg":"func2 will call func3","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
{"level":"INFO","time":"2020-04-14T16:35:36.152183+08:00","logger":"root.ctxLogger","msg":"func3 be called","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
{"level":"INFO","time":"2020-04-14T16:35:36.152189+08:00","logger":"root.ctxLogger","msg":"func3 is called","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
{"level":"INFO","time":"2020-04-14T16:35:36.152197+08:00","logger":"root.ctxLogger","msg":"func2 is called","pid":30881,"traceID":"logging-bqanbm3ipt37h899lbu0"}
```

请求响应头中也包含 Trace ID：

```curl
curl localhost:8080/ping -v
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /ping HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Traceid: logging-bqanbm3ipt37h899lbu0
< Date: Tue, 14 Apr 2020 08:35:36 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
pong* Closing connection 0
```

请求时如果指定Header `-H "traceID: x-y-z"`，demo将使用该值作为trace id

## 动态修改 logger 日志级别

logging 可以在代码中对 AtomicLevel 调用 SetLevel 动态修改日志级别，也可以通过请求 HTTP 接口修改。
默认 logger 使用 `:1903` 运行 HTTP 服务来接收请求修改日志级别。实际使用中日志级别通常写在配置文件中，
可以通过监听配置文件的修改来动态调用 SetLevel 方法。

**示例**:

```golang
package main

import (
    "fmt"
    "github.com/axiaoxin-com/logging"
    "io/ioutil"
    "net/http"
    "strings"

    "go.uber.org/zap"
)

// level 全局变量，便于动态修改，初始化为 Debug 级别
var level zap.AtomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)

func main() {
    /* change log level on fly */

    // 创建指定Level的logger，并开启http服务
    options := logging.Options{
        Level:           level,
        AtomicLevelAddr: ":2012",
    }
    logger, _ := logging.NewLogger(options)
    logger.Debug("Debug level msg", zap.Any("current level", level.Level()))

    /* 函数内部修改 */
    // 使用SetLevel动态修改logger 日志级别为error
    // 实际应用中可以监听配置文件中日志级别配置项的变化动态调用该函数
    level.SetLevel(zap.ErrorLevel)
    // Info 级别将不会被打印
    logger.Info("Info level msg will not be logged")
    // 只会打印error以上
    logger.Error("Error level msg", zap.Any("current level", level.Level()))

    // Output:
    // {"level":"DEBUG","time":"2020-04-13T19:34:46.12339+08:00","logger":"root","caller":"example/atomiclevel.go:18","msg":"Debug level msg","pid":21546,"current level":"debug"}
    // {"level":"ERROR","time":"2020-04-13T19:34:46.123555+08:00","logger":"root","caller":"example/atomiclevel.go:26","msg":"Error Level msg","pid":21546,"current level":"error"}

    /* 在外部通过HTTP接口修改 */
    // 通过HTTP方式动态修改当前的error level为debug level
    // 查询当前 level
    url := "http://localhost" + options.AtomicLevelAddr
    resp, _ := http.Get(url)
    content, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    fmt.Println("currentlevel:", string(content))
    logger.Info("Info level will not be logged")

    // 修改level为debug
    c := &http.Client{}
    req, _ := http.NewRequest("PUT", url, strings.NewReader(`{"level": "debug"}`))
    resp, _ = c.Do(req)
    content, _ = ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    fmt.Println("newlevel:", string(content))

    logger.Debug("level is changed on fly!")

    // Output:
    // currentlevel: {"level":"error"}
    //
    // newlevel: {"level":"debug"}
    //
    // {"level":"DEBUG","time":"2020-04-13T20:04:25.694969+08:00","logger":"root","caller":"example/atomiclevel.go:56","msg":"level is changed on fly!","pid":55317}
}
```

## 自定义 logger EncoderConfig 字段名

**示例**：

```
import "github.com/axiaoxin-com/logging"

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
```


## 日志保存到文件并自动 rotate

使用 lumberjack 将日志保存到文件并 rotate ，采用 zap 的 RegisterSink 方法和 Config.OutputPaths 字段添加自定义的日志输出的方式来使用 lumberjack。


**示例**

```golang
package main

import (
    "github.com/axiaoxin-com/logging"
)

// Options 传入 LumberjacSink，并在 OutputPaths 中添加对应 scheme 就能将日志保存到文件并自动 rotate
func main() {
    /* 使用logger将日志输出到x.log */
    // 创建一个lumberjack的sink，scheme 为 lumberjack，日志文件为 /tmp/x.log , 保存 7 天，保留 10 份文件，文件大小超过 100M，使用压缩备份，压缩文件名使用 localtime
    sink := logging.NewLumberjackSink("lumberjack", "/tmp/x.log", 7, 10, 100, true, true)
    // 创建logger时，设置该sink，OutputPaths 设置为对应 scheme
    options := logging.Options{
        LumberjackSink: sink,
        // 使用 sink 中设置的url scheme 即 lumberjack: 或 lumberjack://
        OutputPaths: []string{"lumberjack:"},
    }
    // 创建logger
    logger, _ := logging.NewLogger(options)
    // 日志将打到x.log文件中
    logger.Debug("xxx")
}
```
