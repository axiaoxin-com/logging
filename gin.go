package logging

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// prometheus namespace
	promNamespace = "logging"
	// gin prometheus labels
	promGinLabels = []string{
		"status_code",
		"path",
		"method",
	}
	promGinReqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "req_count",
			Help:      "gin server request count",
		}, promGinLabels,
	)
	promGinReqLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: promNamespace,
			Name:      "req_latency",
			Help:      "gin server request latency in seconds",
		}, promGinLabels,
	)

	// 默认慢请求时间 3s
	defaultGinSlowThreshold = time.Second * 3
)

// GetGinTraceIDFromHeader 从 gin 的 request header 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(string(TraceIDKeyname))
}

// GetGinTraceIDFromQueryString 从 gin 的 querystring 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(string(TraceIDKeyname))
}

// GetGinTraceIDFromPostForm 从 gin 的 postform 中获取 key 为 TraceIDKeyname 的值作为 traceid
func GetGinTraceIDFromPostForm(c *gin.Context) string {
	return c.PostForm(string(TraceIDKeyname))
}

// GinLogDetails gin 日志中间件记录的信息
type GinLogDetails struct {
	// 请求处理完成时间
	Timestamp time.Time `json:"timestamp"`
	// 请求方法
	Method string `json:"method"`
	// 请求 Path
	Path string `json:"path"`
	// 请求 RawQuery
	Query string `json:"query"`
	// http 协议版本
	Proto string `json:"proto"`
	// 请求内容长度
	ContentLength int `json:"content_length"`
	// 请求的 host host:port
	Host string `json:"host"`
	// 请求 remote addr  host:port
	RemoteAddr string `json:"remote_addr"`
	// uri
	RequestURI string `json:"request_uri"`
	// referer
	Referer string `json:"referer"`
	// user agent
	UserAgent string `json:"user_agent"`
	// 真实客户端 ip
	ClientIP string `json:"client_ip"`
	// content type
	ContentType string `json:"content_type"`
	// handler name
	HandlerName string `json:"handler_name"`
	// http 状态码
	StatusCode int `json:"status_code"`
	// 响应 body 字节数
	BodySize int `json:"body_size"`
	// 请求处理耗时 (秒)
	Latency float64 `json:"latency"`
	// Context 中的 Keys
	ContextKeys map[string]interface{} `json:"context_keys,omitempty"`
	// http request header
	RequestHeader http.Header `json:"request_header,omitempty"`
	// http Request Form
	RequestForm url.Values `json:"request_form,omitempty"`
	// 请求 body
	RequestBody interface{} `json:"request_body,omitempty"`
	// 响应 Body
	ResponseBody interface{} `json:"response_body,omitempty"`
}

// GinLoggerConfig GinLogger 支持的配置项字段定义
type GinLoggerConfig struct {
	// Optional. Default value is logging.defaultGinLogFormatter
	Formatter func(context.Context, GinLogDetails) string
	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
	// SkipPathRegexps skip path by regexp
	SkipPathRegexps []string
	// TraceIDFunc 获取或生成 trace id 的函数
	// Optional.
	TraceIDFunc func(context.Context) string
	// InitFieldsFunc 获取 logger 初始字段方法 key 为字段名 value 为字段值
	InitFieldsFunc func(context.Context) map[string]interface{}
	// 是否使用详细模式打印日志，记录更多字段信息
	// Optional.
	EnableDetails bool

	// 以下选项开启后对性能有影响，适用于接口调试，慎用。

	// 是否打印 context keys
	// Optional.
	EnableContextKeys bool
	// 是否打印请求头信息
	// Optional.
	EnableRequestHeader bool
	// 是否打印请求form信息
	// Optional.
	EnableRequestForm bool
	// 是否打印请求体信息
	// Optional.
	EnableRequestBody bool
	// 是否打印响应体信息
	// Optional.
	EnableResponseBody bool

	// 慢请求时间阈值 请求处理时间超过该值则使用 Error 级别打印日志
	SlowThreshold time.Duration
}

// GinLogger 以默认配置生成 gin 的 Logger 中间件
func GinLogger() gin.HandlerFunc {
	return GinLoggerWithConfig(GinLoggerConfig{})
}

// gin 访问日志中 msg 字段的输出格式
func defaultGinLogFormatter(c context.Context, m GinLogDetails) string {
	_, shortHandlerName := path.Split(m.HandlerName)
	msg := fmt.Sprintf("%s|%s|%s%s|%s|%d|%f",
		m.ClientIP,
		m.Method,
		m.Host,
		m.RequestURI,
		shortHandlerName,
		m.StatusCode,
		m.Latency,
	)
	return msg
}

func defaultGinTraceIDFunc(c context.Context) (traceID string) {
	if c == nil {
		c = context.Background()
	}

	if gc, ok := c.(*gin.Context); ok {

		traceID = GetGinTraceIDFromHeader(gc)
		if traceID != "" {
			return
		}
		traceID = GetGinTraceIDFromPostForm(gc)
		if traceID != "" {
			return
		}
		traceID = GetGinTraceIDFromQueryString(gc)
		if traceID != "" {
			return
		}
	}
	traceID = CtxTraceID(c)
	return
}

// GinLoggerWithConfig 根据配置信息生成 gin 的 Logger 中间件
// 中间件会记录访问信息，根据状态码确定日志级别， 500 以上为 Error ， 400-500 默认为 Warn ， 400 以下默认为 Info
// api 请求进来的 context 的函数无需在其中打印 err ，使用 c.Error(err)会在请求完成时自动打印 error
// context 中有 error 则日志忽略返回码始终使用 error 级别
func GinLoggerWithConfig(conf GinLoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultGinLogFormatter
	}
	getTraceID := conf.TraceIDFunc
	if getTraceID == nil {
		getTraceID = defaultGinTraceIDFunc
	}

	var skip map[string]struct{}
	if length := len(conf.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range conf.SkipPaths {
			skip[path] = struct{}{}
		}
	}
	var skipRegexps []*regexp.Regexp
	for _, p := range conf.SkipPathRegexps {
		if r, err := regexp.Compile(p); err != nil {
			Error(nil, "skip path regexps compile "+p+" error:"+err.Error())
		} else {
			skipRegexps = append(skipRegexps, r)
		}
	}

	if conf.SlowThreshold.Seconds() <= 0 {
		conf.SlowThreshold = defaultGinSlowThreshold
	}

	return func(c *gin.Context) {
		traceID := getTraceID(c)
		// 设置 trace id 到 request header 中
		c.Request.Header.Set(string(TraceIDKeyname), traceID)
		// 设置 trace id 到 response header 中
		c.Writer.Header().Set(string(TraceIDKeyname), traceID)
		// 设置 trace id 和 ctxLogger 到 context 中
		ginLogger := CloneLogger("gin")
		if conf.InitFieldsFunc != nil {
			for k, v := range conf.InitFieldsFunc(c) {
				ginLogger = ginLogger.With(zap.Any(k, v))
			}
		}
		_, ctxLogger := NewCtxLogger(c, ginLogger, traceID)

		start := time.Now()

		// 获取请求信息
		details := GinLogDetails{
			Method:        c.Request.Method,
			Path:          c.Request.URL.Path,
			Query:         c.Request.URL.RawQuery,
			Proto:         c.Request.Proto,
			ContentLength: int(c.Request.ContentLength),
			Host:          c.Request.Host,
			RemoteAddr:    c.Request.RemoteAddr,
			RequestURI:    c.Request.RequestURI,
			Referer:       c.Request.Referer(),
			UserAgent:     c.Request.UserAgent(),

			ClientIP:    c.ClientIP(),
			ContentType: c.ContentType(),
			HandlerName: c.HandlerName(),
		}

		// 获取并保存原始请求 body
		if conf.EnableRequestBody {
			body := GetGinRequestBody(c)
			if err := jsoniter.Unmarshal(body, &details.RequestBody); err != nil {
				details.RequestBody = string(body)
			}
		}
		// 获取并保存原始请求 form
		if conf.EnableRequestForm {
			details.RequestForm = c.Request.Form
		}
		// 获取并保存原始请求 header
		if conf.EnableRequestHeader {
			details.RequestHeader = c.Request.Header
		}
		rspBodyWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		if conf.EnableResponseBody {
			// 开启记录响应 body 时，保存 body 到 rspBodyWriter.body 中
			c.Writer = rspBodyWriter
		}

		defer func() {
			// 获取响应信息
			details.StatusCode = c.Writer.Status()
			details.BodySize = c.Writer.Size()
			details.Timestamp = time.Now()
			details.Latency = details.Timestamp.Sub(start).Seconds()

			// 创建 logger
			accessLogger := ctxLogger.Named("access_logger").With(
				zap.String("client_ip", details.ClientIP),
				zap.String("method", details.Method),
				zap.String("path", details.Path),
				zap.String("host", details.Host),
				zap.Int("status_code", details.StatusCode),
				zap.Float64("latency", details.Latency),
			)
			// handler 中使用 c.Error(err) 后，会打印到 context_errors 字段中
			if len(c.Errors) > 0 {
				accessLogger = accessLogger.With(zap.String("context_errors", c.Errors.String()))
			}
			// 判断是否打印 context keys
			if conf.EnableContextKeys {
				details.ContextKeys = c.Keys
				accessLogger = accessLogger.With(zap.Any("context_keys", details.ContextKeys))
			}
			// 判断是否打印请求 header
			if conf.EnableRequestHeader {
				accessLogger = accessLogger.With(zap.Any("request_header", details.RequestHeader))
			}
			// 判断是否打印请求 form
			if conf.EnableRequestForm {
				accessLogger = accessLogger.With(zap.Any("request_form", details.RequestForm))
			}
			// 判断是否打印请求 body
			if conf.EnableRequestBody {
				accessLogger = accessLogger.With(zap.Any("request_body", details.RequestBody))
			}
			// 判断是否打印响应 body
			if conf.EnableResponseBody {
				if err := jsoniter.Unmarshal(rspBodyWriter.body.Bytes(), &details.ResponseBody); err != nil {
					details.ResponseBody = rspBodyWriter.body.String()
				}
				accessLogger = accessLogger.With(zap.Any("response_body", details.ResponseBody))
			}

			// details logger 可以打印更多字段
			detailsLogger := accessLogger.Named("details").With(
				zap.String("query", details.Query),
				zap.String("proto", details.Proto),
				zap.Int("content_length", details.ContentLength),
				zap.String("remote_addr", details.RemoteAddr),
				zap.String("request_uri", details.RequestURI),
				zap.String("referer", details.Referer),
				zap.String("user_agent", details.UserAgent),
				zap.String("content_type", details.ContentType),
				zap.Int("body_size", details.BodySize),
				zap.String("handler_name", details.HandlerName),
			)

			logger := accessLogger
			// 是否打印 details 字段
			if conf.EnableDetails {
				logger = detailsLogger
			}

			// 打印访问日志，根据状态码确定日志打印级别
			log := logger.Info
			if details.StatusCode >= http.StatusInternalServerError {
				// 500+ 始终打印带 details 的 error 级别日志
				errLogger := detailsLogger.Named("err")
				// 无视配置开关，打印全部能搜集的信息
				if len(details.ContextKeys) == 0 {
					errLogger = errLogger.With(zap.Any("context_keys", c.Keys))
				}
				if len(details.RequestHeader) == 0 {
					errLogger = errLogger.With(zap.Any("request_header", c.Request.Header))
				}
				if len(details.RequestForm) == 0 {
					errLogger = errLogger.With(zap.Any("request_form", c.Request.Form))
				}
				if details.RequestBody == nil {
					errLogger = errLogger.With(zap.String("request_body", string(GetGinRequestBody(c))))
				}
				if details.ResponseBody == nil {
					errLogger = errLogger.With(zap.String("response_body", rspBodyWriter.body.String()))
				}
				log = errLogger.Error
			} else if details.StatusCode >= http.StatusBadRequest {
				// 400+ 默认使用 warn 级别。如果有 errors 则使用 error 级别
				log = logger.Warn
				if len(c.Errors) > 0 {
					log = logger.Error
				}
			} else if len(c.Errors) > 0 {
				log = logger.Error
			}

			skipLog := false
			if _, exists := skip[details.Path]; exists {
				skipLog = true
			} else {
				for _, p := range skipRegexps {
					if p.MatchString(details.Path) {
						skipLog = true
						break
					}
				}
			}
			if !skipLog {
				// 慢请求使用 Warn 记录
				if details.Latency > conf.SlowThreshold.Seconds() {
					logger.Warn(
						formatter(c, details)+" hit slow request.",
						zap.Float64("slow_threshold", conf.SlowThreshold.Seconds()),
					)
				} else {
					log(formatter(c, details))
				}

				// update prometheus info
				labels := []string{fmt.Sprint(details.StatusCode), details.Path, details.Method}
				promGinReqCount.WithLabelValues(labels...).Inc()
				promGinReqLatency.WithLabelValues(labels...).Observe(details.Latency)
			}
		}()

		c.Next()
	}
}

// GetGinRequestBody 获取请求 body
func GetGinRequestBody(c *gin.Context) []byte {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		} else {
			requestBody = body
			// body 被 read 、 bind 之后会被置空，需要重置
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
	}
	return requestBody
}

// 用于记录响应 body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 覆盖 ResponseWriter 接口的 Write 方法，将 body 保存到 responseBodyWriter.body 中
func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
