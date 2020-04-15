package main

import (
	"context"

	"github.com/axiaoxin-com/logging"

	"github.com/gin-gonic/gin"
)

func func1(c context.Context) {
	// 使用CtxLogger打印带trace id的日志
	logging.CtxLogger(c).Info("func1 begin")
	func2(c)
	// 使用logging全局方法打印带trace id的日志
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

	// 使用中间件注册获取trace id
	// 使用默认的回调方法从Header中获取Key为traceID的值作为trace id
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

请求响应头中也包含 Trace ID, 请求时如果指定Header `-H "traceID: x-y-z"`，demo将使用该值作为trace id

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
