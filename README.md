# logging

logging 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的全部方法，所有日志打印方法开箱即用
- 提供多种快速创建 logger 的方法
- 支持在使用 Error 及其以上级别打印日志时自动将该事件上报到 Sentry
- 支持从 context.Context/gin.Context 中创建、获取带有 Trace ID 的 logger


## 安装

```
go get -u github/axiaoxin-com/logging
```

## 开箱即用

logging 提供的开箱即用方法都是使用自身默认 logger 和 sugared logger 打印，
默认 logger 使用 JSON 格式打印日志内容到 stderr ，
不带 Sentry 上报功能，
可通过 HTTP 调用 `curl -XPUT "http://localhost:1903" -d '{"level": "info"}'` 动态修改日志级别，
默认带有初始字段 pid 打印进程 ID

**示例**

```golang
import "github.com/axiaoxin-com/logging"

// 等同于 *zap.Logger 的 Debug
logging.Debug("Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
// {"level":"DEBUG","time":"2020-04-12T02:56:39.32688+08:00","logger":"root","caller":"logging/global.go:120","msg":"Debug message","pid":27907,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

// 等同于 *zap.Logger.Sugar().Debug
logging.Debugs("SDebug message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
// {"level":"DEBUG","time":"2020-04-12T02:56:39.327239+08:00","logger":"root","caller":"logging/global.go:10","msg":"SDebug message123 false [1 2 3] map[i:1 s:s]","pid":27907}

// 等同于 *zap.Logger.Sugar().Debugf
logging.Debugf("SDebugf message, %s", "ok")
//{"level":"DEBUG","time":"2020-04-12T02:56:39.327287+08:00","logger":"root","caller":"logging/global.go:47","msg":"SDebugf message, ok","pid":27907}

// 等同于 *zap.Logger.Sugar().Debugw
logging.Debugw("SDebug message", "name", "axiaoxin", "age", 18)
//{"level":"DEBUG","time":"2020-04-12T02:56:39.327301+08:00","logger":"root","caller":"logging/global.go:84","msg":"SDebug message","pid":27907,"name":"axiaoxin","age":18}
```
## 快速创建你的 Logger

logging 提供多种方式快速获取一个 logger 来打印日志


**示例1**：创建一个 logging 自身使用的默认 logger，并设置 sentry 上报错误

```golang
import "github.com/axiaoxin-com/logging"

// 创建一个默认 logger
logger := logging.DefaultLogger()

// logging 内部默认的 logger 不支持 sentry 上报，可以通过以下方法设置 sentry
// 创建 sentry 客户端
sentryClient, _ := logging.GetSentryClientByDSN("YOUR_SENTRY_DSN")
// 设置 sentry，使用该 logger 打印 Error 及其以上级别的日志事件将会自动上报到 Sentry
logger = logging.SentryAttach(logger, sentryClient)
```

**示例2**: 使用 NewLogger 方法创建一个默认配置的 logger （不支持 sentry 和 http 动态修改日志级别）

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

// {"level":"DEBUG","time":"2020-04-13T00:20:36.614438+08:00","logger":"root.subname","caller":"example/logging.go:27","msg":"CloneDefaultLogger Debug","pid":54273,"str_field":"field_value"}
```

初始字段可以不传，克隆的 logger 名称会是 root.subname，该 logger 打印的日志都会带上传入的字段

**示例5**: 快速克隆一个默认 sugared logger，并添加初始字段

```golang
import "github.com/axiaoxin-com/logging"

logger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))

logger.Debug("CloneDefaultSLogger Debug")

// {"level":"DEBUG","time":"2020-04-13T00:24:41.629175+08:00","logger":"root.subname","caller":"example/logging.go:32","msg":"CloneDefaultSLogger Debug","pid":73087,"foo":123,"str_field":"field_value"}
```

初始字段可以不传，克隆的 sugared logger 名称会是 root.subname，添加的初始字段则该 logger 打印的日志都会带上传入的字段


## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。


**示例1**: 普通函数中使用 CtxLogger

```golang
import "github.com/axiaoxin-com/logging"

// 初始化一个 context
ctx := context.Background()
// 生成一个 trace id，如果 context 是 gin.Context，会尝试从其中获取，否则尝试从 context.Context 获取，获取不到则新生成
traceID := logging.CtxTraceID(ctx)
// 设置 trace id 到 context 中， 会尝试同时设置到 gin.Context 中
ctx = logging.Context(ctx, traceID)
// 从 context 中获取 logger，会尝试从 gin.Context 中获取，context 中没有 logger 则克隆默认 logger 作为 context logger
ctxlogger = logging.CtxLogger(ctx)
// log with trace id
ctxlogger.Debug("ctxlogger with trace id debug")

// Output:
// {"level":"DEBUG","time":"2020-04-13T01:34:19.697443+08:00","logger":"root","caller":"logging/global.go:120","msg":"no logger in context, clone the default logger as ctxLogger","pid":88649}
// {"level":"DEBUG","time":"2020-04-13T01:34:19.697453+08:00","logger":"root.ctxLogger","caller":"example/main.go:51","msg":"ctxlogger with trace id debug","pid":88649,"traceID":"logging-bq9l26ript35kicii5tg"}
```

**示例2**: gin 使用 CtxLogger 打印带 Trace ID 的日志

```golang
package main

import (
	"context"
	"github/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	logging.CtxLogger(c).Info("func1 will call func2")
	func2(c)
}

func func2(c context.Context) {
	logging.CtxLogger(c).Info("func2 will call func3")
	func3(c)
}

func func3(c context.Context) {
	logging.CtxLogger(c).Info("func3 be called")
}

func main() {
	r := gin.Default()

	r.Use(logging.GinTraceIDMiddleware(logging.TraceIDKey))

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
{"level":"DEBUG","time":"2020-04-13T03:22:34.86741+08:00","logger":"root","caller":"logging/global.go:120","msg":"context dose not exist trace id key, generate a new trace id","pid":82451}
{"level":"DEBUG","time":"2020-04-13T03:22:34.867633+08:00","logger":"root","caller":"logging/global.go:120","msg":"no logger in context, clone the default logger as ctxLogger","pid":82451}
{"level":"INFO","time":"2020-04-13T03:22:34.867649+08:00","logger":"root.ctxLogger","caller":"example/gin.go:30","msg":"ping ping pong pong","pid":82451,"traceID":"logging-bq9mkujipt3444tmd1vg"}
{"level":"INFO","time":"2020-04-13T03:22:34.86766+08:00","logger":"root.ctxLogger","caller":"example/gin.go:11","msg":"func1 will call func2","pid":82451,"traceID":"logging-bq9mkujipt3444tmd1vg"}
{"level":"INFO","time":"2020-04-13T03:22:34.867668+08:00","logger":"root.ctxLogger","caller":"example/gin.go:16","msg":"func2 will call func3","pid":82451,"traceID":"logging-bq9mkujipt3444tmd1vg"}
{"level":"INFO","time":"2020-04-13T03:22:34.867674+08:00","logger":"root.ctxLogger","caller":"example/gin.go:21","msg":"func3 be called","pid":82451,"traceID":"logging-bq9mkujipt3444tmd1vg"}
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
< Traceid: logging-bq9mkujipt3444tmd1vg
< Date: Sun, 12 Apr 2020 19:22:34 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
pong* Closing connection 0
```
