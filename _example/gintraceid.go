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
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 使用 GinLogger 中间件记录访问日志和生成 trace id
	r.Use(logging.GinLogger())

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

{"level":"ERROR","time":"2020-08-28 09:21:57.675677","logger":"logging.gin","caller":"example/gin.go:func1:36","msg":"ping ping pong pong","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.675892","logger":"logging.gin","caller":"example/gin.go:func1:13","msg":"func1 begin","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.675907","logger":"logging.gin","caller":"example/gin.go:func2:20","msg":"func2 begin","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.675917","logger":"logging.gin","caller":"example/gin.go:func3:26","msg":"in func3","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.675933","logger":"logging.gin","caller":"example/gin.go:func2:22","msg":"func2 end","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.675957","logger":"logging.gin","caller":"example/gin.go:func1:16","msg":"func1 end","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0"}
{"level":"INFO","time":"2020-08-28 09:21:57.676170","logger":"logging.gin.access_logger","caller":"logging/gin.go:func1:229","msg":"2020-08-28 09:21:57.675988|127.0.0.1|GET|localhost:8080/ping|main.main.func1|200|0.000329","pid":61374,"server_ip":"10.64.35.43","trace_id":"logging_bt45od98d3bevfj7lni0","details":{"req_time":"2020-08-28T09:21:57.675988+08:00","method":"GET","path":"/ping","query":"","proto":"HTTP/1.1","content_length":0,"host":"localhost:8080","remote_addr":"127.0.0.1:57316","request_uri":"/ping","referer":"","user_agent":"curl/7.64.1","client_ip":"127.0.0.1","content_type":"","handler_name":"main.main.func1","status_code":200,"body_size":4,"latency":0.000328505,"context_keys":{"ctx_logger":{},"trace_id":"logging_bt45od98d3bevfj7lni0"}}}

请求响应头中也包含 Trace ID, 请求时如果指定 Header `-H "trace_id: x-y-z"`， demo 将使用该值作为 trace id

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
< Trace_id: logging-bqbeq9ript38cuae9nb0
< Content-Length: 4
<
* Connection #0 to host localhost left intact
pong* Closing connection 0
*/
