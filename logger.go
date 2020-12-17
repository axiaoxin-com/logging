// Package logging 简单封装了在日常使用 zap 打日志时的常用方法。
//
// 提供快速使用 zap 打印日志的全部方法，所有日志打印方法开箱即用
//
// 提供多种快速创建 logger 的方法
//
// 支持在使用 Error 及其以上级别打印日志时自动将该事件上报到 Sentry
//
// 支持从 context.Context/gin.Context 中创建、获取带有 Trace ID 的 logger
package logging

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global zap Logger with pid field
	logger *zap.Logger
	// 默认 sentry client
	sentryClient *sentry.Client
	// outPaths zap 日志默认输出位置
	outPaths = []string{"stdout"}
	// initialFields 默认初始字段为进程 id
	initialFields = map[string]interface{}{
		"pid":       syscall.Getpid(),
		"server_ip": ServerIP(),
	}
	// loggerName 默认 logger name 为 logging
	loggerName = "logging"
	// atomicLevel 默认 logger atomic level 级别默认为 debug
	atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	// EncoderConfig 默认的日志字段名配置
	EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   CallerEncoder,
	}
	// 读写锁
	rwMutex sync.RWMutex
)

// AtomicLevelServerOption AtomicLevel server 相关配置
type AtomicLevelServerOption struct {
	Addr     string // http 动态修改日志级别服务运行地址
	Path     string // 设置 url path ，可选
	Username string // 请求时设置 basic auth 认证的用户名，可选
	Password string // 请求时设置 basic auth 认证的密码，可选，与 username 同时存在才开启 basic auth
}

// Options new logger options
type Options struct {
	Name              string                  // logger 名称
	Level             string                  // 日志级别 debug, info, warn, error dpanic, panic, fatal
	Format            string                  // 日志格式
	OutputPaths       []string                // 日志输出位置
	InitialFields     map[string]interface{}  // 日志初始字段
	DisableCaller     bool                    // 是否关闭打印 caller
	DisableStacktrace bool                    // 是否关闭打印 stackstrace
	SentryClient      *sentry.Client          // sentry 客户端
	EncoderConfig     *zapcore.EncoderConfig  // 配置日志字段 key 的名称
	LumberjackSink    *LumberjackSink         // lumberjack sink 支持日志文件 rotate
	AtomicLevelServer AtomicLevelServerOption // AtomicLevel server 相关配置
}

const (
	// SentryDSNEnvKey 引入包时默认创建 logger 将尝试从该环境变量名中获取 sentry dsn
	SentryDSNEnvKey = "SENTRY_DSN"
	// SentryDebugEnvKey 尝试从该环境变量中获取 sentry 是否开启 debug 模式
	SentryDebugEnvKey = "SENTRY_DEBUG"
	// AtomicLevelAddrEnvKey 初始化时尝试获取该环境变量用于设置动态修改日志级别的 http 服务运行地址
	AtomicLevelAddrEnvKey = "ATOMIC_LEVEL_ADDR"
)

// init the global logger
func init() {
	var err error
	// 尝试从环境变量获取 sentry dsn
	if dsn := os.Getenv(SentryDSNEnvKey); dsn != "" {
		debugStr := os.Getenv(SentryDebugEnvKey)
		debug := false
		if strings.ToLower(debugStr) != "" {
			debug = true
		}
		sentryClient, err = NewSentryClient(dsn, debug)
		if err != nil {
			log.Println(err)
		}
	}

	options := Options{
		Name:              loggerName,
		Level:             "debug",
		Format:            "json",
		OutputPaths:       outPaths,
		InitialFields:     initialFields,
		DisableCaller:     false,
		DisableStacktrace: true,
		SentryClient:      sentryClient,
		AtomicLevelServer: AtomicLevelServerOption{
			Addr: os.Getenv(AtomicLevelAddrEnvKey),
		},
		EncoderConfig:  &EncoderConfig,
		LumberjackSink: nil,
	}
	logger, err = NewLogger(options)
	if err != nil {
		log.Println(err)
	}
}

// NewLogger return a zap Logger instance
func NewLogger(options Options) (*zap.Logger, error) {
	cfg := zap.Config{}
	// 设置日志级别
	lvl := strings.ToLower(options.Level)
	if _, exists := AtomicLevelMap[lvl]; !exists {
		cfg.Level = atomicLevel
	} else {
		cfg.Level = AtomicLevelMap[lvl]
		atomicLevel = cfg.Level
	}
	// 设置 encoding 默认为 json
	if strings.ToLower(options.Format) == "console" {
		cfg.Encoding = "console"
	} else {
		cfg.Encoding = "json"
	}
	// 设置 output 没有传参默认全部输出到 stderr
	if len(options.OutputPaths) == 0 {
		cfg.OutputPaths = outPaths
		cfg.ErrorOutputPaths = outPaths
	} else {
		cfg.OutputPaths = options.OutputPaths
		cfg.ErrorOutputPaths = options.OutputPaths
	}
	// 设置 InitialFields 没有传参使用默认字段
	// 传了就添加到现有的初始化字段中
	if len(options.InitialFields) > 0 {
		for k, v := range options.InitialFields {
			initialFields[k] = v
		}
	}
	cfg.InitialFields = initialFields
	// 设置 disablecaller
	cfg.DisableCaller = options.DisableCaller
	// 设置 disablestacktrace
	cfg.DisableStacktrace = options.DisableStacktrace

	// 设置 encoderConfig
	if options.EncoderConfig == nil {
		cfg.EncoderConfig = EncoderConfig
	} else {
		cfg.EncoderConfig = *options.EncoderConfig
	}

	// Sampling 实现了日志的流控功能，或者叫采样配置，主要有两个配置参数， Initial 和 Thereafter ，实现的效果是在 1s 的时间单位内，如果某个日志级别下同样内容的日志输出数量超过了 Initial 的数量，那么超过之后，每隔 Thereafter 的数量，才会再输出一次。是一个对日志输出的保护功能。
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	// 注册 lumberjack sink ，支持 Outputs 指定为文件时可以使用 lumberjack 对日志文件自动 rotate
	if options.LumberjackSink != nil {
		if err := RegisterLumberjackSink(options.LumberjackSink); err != nil {
			Error(nil, "RegisterSink error", zap.Error(err))
		}
	}

	// 生成 logger
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	// 如果传了 sentryclient 则设置 sentrycore
	if options.SentryClient != nil {
		logger = SentryAttach(logger, options.SentryClient)
	}

	// 设置 logger 名字，没有传参使用默认名字
	if options.Name != "" {
		logger = logger.Named(options.Name)
	} else {
		logger = logger.Named(loggerName)
	}
	if options.AtomicLevelServer.Addr != "" {
		runAtomicLevelServer(cfg.Level, options.AtomicLevelServer)
	}
	return logger, nil
}

// CloneLogger return the global logger copy which add a new name
func CloneLogger(name string, fields ...zap.Field) *zap.Logger {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	copy := *logger
	clogger := &copy
	clogger = clogger.Named(name)
	if len(fields) > 0 {
		clogger = clogger.With(fields...)
	}
	return clogger
}

// AttachCore add a core to zap logger
func AttachCore(l *zap.Logger, c zapcore.Core) *zap.Logger {
	return l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, c)
	}))
}

// ReplaceLogger 替换默认的全局 logger 为传入的新 logger
// 返回函数，调用它可以恢复全局 logger 为上一次的 logger
func ReplaceLogger(newLogger *zap.Logger) func() {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	// 备份原始 logger 以便恢复
	prevLogger := logger
	// 替换为新 logger
	logger = newLogger
	return func() { ReplaceLogger(prevLogger) }
}

// TextLevel 返回默认 logger 的 字符串 level
func TextLevel() string {
	b, _ := atomicLevel.MarshalText()
	return string(b)
}

// SetLevel 使用字符串级别设置默认 logger 的 atomic level
func SetLevel(lvl string) {
	Warn(nil, "Set logging atomicLevel "+lvl)
	atomicLevel.UnmarshalText([]byte(strings.ToLower(lvl)))
}

// SentryClient 返回默认 sentry client
func SentryClient() *sentry.Client {
	return sentryClient
}

// ServerIP 获取当前 IP
func ServerIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// 运行 atomic level server
func runAtomicLevelServer(atomicLevel zap.AtomicLevel, options AtomicLevelServerOption) {
	go func() {
		// curl -X GET http://host:port
		// curl -X PUT http://host:port -d '{"level":"info"}'
		Debug(nil, "Running AtomicLevel HTTP server on "+options.Addr+options.Path)
		urlPath := "/"
		if options.Path != "" {
			urlPath = options.Path
		}

		levelServer := http.NewServeMux()
		levelServer.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
			msg := fmt.Sprintf("%s %s the logger atomic level", r.RemoteAddr, r.Method)
			_, logger := NewCtxLogger(r.Context(), CloneLogger("atomiclevel"), r.Header.Get(string(TraceIDKeyname)))
			if r.Method == http.MethodPut {
				b, _ := ioutil.ReadAll(r.Body)
				msg += " to " + string(b)
				r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
				logger.Warn(msg)
			} else {
				logger.Info(msg)
			}
			if options.Username != "" && options.Password != "" {
				if _, _, ok := r.BasicAuth(); !ok {
					http.Error(w, "need to basic auth", http.StatusUnauthorized)
					return
				}
			}
			atomicLevel.ServeHTTP(w, r)
		})
		if err := http.ListenAndServe(options.Addr, levelServer); err != nil {
			Error(nil, "logging NewLogger levelServer ListenAndServe error:"+err.Error())
		}
	}()
}
