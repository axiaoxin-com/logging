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
	"log"
	"net/http"
	"reflect"
	"strings"
	"syscall"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// logger default global zap Logger with pid field
	logger *zap.Logger
	// slogger default global zap sugared Logger with pid field
	slogger *zap.SugaredLogger
	// defaultOutPaths zap日志默认输出位置
	defaultOutPaths = []string{"stderr"}
	// defaultInitialFields 默认初始字段为进程id
	defaultInitialFields = map[string]interface{}{
		"pid": syscall.Getpid(),
	}
	// defaultLoggerName 默认logger name为root
	defaultLoggerName = "root"
	// defaultLevel 默认日志级别为debug
	defaultLevel           = zap.NewAtomicLevelAt(zap.DebugLevel)
	defaultAtomicLevelAddr = ":1903"
	// defaultEncoderConfig 默认的日志字段名配置
	defaultEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// TextLevelMap string level mapping zap AtomicLevel
	TextLevelMap = map[string]zap.AtomicLevel{
		"debug":  zap.NewAtomicLevelAt(zap.DebugLevel),
		"info":   zap.NewAtomicLevelAt(zap.InfoLevel),
		"warn":   zap.NewAtomicLevelAt(zap.WarnLevel),
		"error":  zap.NewAtomicLevelAt(zap.ErrorLevel),
		"dpanic": zap.NewAtomicLevelAt(zap.DPanicLevel),
		"panic":  zap.NewAtomicLevelAt(zap.PanicLevel),
		"fatal":  zap.NewAtomicLevelAt(zap.FatalLevel),
	}
)

// Options new logger options
type Options struct {
	Name              string                 // logger 名称
	Level             zap.AtomicLevel        // 日志级别
	Format            string                 // 日志格式
	OutputPaths       []string               // 日志输出位置
	InitialFields     map[string]interface{} // 日志初始字段
	DisableCaller     bool                   // 是否关闭打印caller
	DisableStacktrace bool                   // 是否关闭打印stackstrace
	SentryClient      *sentry.Client         // sentry客户端
	AtomicLevelAddr   string                 // http动态修改日志级别的地址，传空不启用
	EncoderConfig     zapcore.EncoderConfig  // 配置日志地点key的名称
	LumberjackSink    *LumberjackSink        // lumberjack sink支持日志文件rotate
}

// init the global default logger
func init() {
	options := Options{
		Name:              defaultLoggerName,
		Level:             defaultLevel,
		Format:            "json",
		OutputPaths:       defaultOutPaths,
		InitialFields:     defaultInitialFields,
		DisableCaller:     true,
		DisableStacktrace: true,
		SentryClient:      nil,
		AtomicLevelAddr:   defaultAtomicLevelAddr,
		EncoderConfig:     defaultEncoderConfig,
		LumberjackSink:    nil,
	}
	var err error
	logger, err = NewLogger(options)
	if err != nil {
		log.Println(err)
	} else {
		slogger = logger.Sugar()
	}
}

// NewLogger return a zap Logger instance
func NewLogger(options Options) (*zap.Logger, error) {
	cfg := zap.Config{}
	// 设置日志级别
	if options.Level == (zap.AtomicLevel{}) {
		cfg.Level = defaultLevel
	} else {
		cfg.Level = options.Level
	}
	// 设置encoding 默认为json
	if strings.ToLower(options.Format) == "console" {
		cfg.Encoding = "console"
	} else {
		cfg.Encoding = "json"
	}
	// 设置output 没有传参默认全部输出到stderr
	if len(options.OutputPaths) == 0 {
		cfg.OutputPaths = defaultOutPaths
		cfg.ErrorOutputPaths = defaultOutPaths
	} else {
		cfg.OutputPaths = options.OutputPaths
		cfg.ErrorOutputPaths = options.OutputPaths
	}
	// 设置InitialFields 没有传参使用默认字段
	if len(options.InitialFields) == 0 {
		cfg.InitialFields = defaultInitialFields
	} else {
		cfg.InitialFields = options.InitialFields
	}
	// 设置disablecaller
	cfg.DisableCaller = options.DisableCaller
	// 设置disablestacktrace
	cfg.DisableStacktrace = options.DisableStacktrace

	// 设置encoderConfig
	if reflect.DeepEqual(options.EncoderConfig, zapcore.EncoderConfig{}) {
		cfg.EncoderConfig = defaultEncoderConfig
	} else {
		cfg.EncoderConfig = options.EncoderConfig
	}

	// Sampling实现了日志的流控功能，或者叫采样配置，主要有两个配置参数，Initial和Thereafter，实现的效果是在1s的时间单位内，如果某个日志级别下同样内容的日志输出数量超过了Initial的数量，那么超过之后，每隔Thereafter的数量，才会再输出一次。是一个对日志输出的保护功能。
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	// 注册 lumberjack sink，支持Outputs指定为文件时可以使用lumberjack对日志文件自动rotate
	if options.LumberjackSink != nil {
		if err := RegisterLumberjackSink(options.LumberjackSink); err != nil {
			Error(nil, "RegisterSink error", zap.Error(err))
		}
	}

	// 生成logger
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	// 如果传了sentryclient则设置sentrycore
	if options.SentryClient != nil {
		logger = SentryAttach(logger, options.SentryClient)
	}

	// 设置logger名字，没有传参使用默认名字
	if options.Name != "" {
		logger = logger.Named(options.Name)
	} else {
		logger = logger.Named(defaultLoggerName)
	}
	if options.AtomicLevelAddr != "" {
		go func() {
			// curl -X GET localhost:1903
			// curl -X PUT localhost:1903 -d '{"level":"info"}'
			levelServer := http.NewServeMux()
			levelServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				cfg.Level.ServeHTTP(w, r)
				logger.Debug("Handled AtomicLevel HTTP request", zap.String("method", r.Method))
			})
			if err := http.ListenAndServe(options.AtomicLevelAddr, levelServer); err != nil {
				Error(nil, "logging NewLogger levelServer ListenAndServe error", zap.Error(err))
			}
		}()
	}
	return logger, nil
}

// DefaultLogger return the global logger
func DefaultLogger() *zap.Logger {
	copy := *logger
	clogger := &copy
	return clogger
}

// DefaultSLogger return the global slogger
func DefaultSLogger() *zap.SugaredLogger {
	copy := *slogger
	cslogger := &copy
	return cslogger
}

// CloneDefaultLogger return the global logger copy which add a new name
func CloneDefaultLogger(name string, fields ...zap.Field) *zap.Logger {
	copy := *logger
	clogger := &copy
	clogger = clogger.Named(name)
	if len(fields) > 0 {
		clogger = clogger.With(fields...)
	}
	return clogger
}

// CloneDefaultSLogger return the global slogger copy which add a new name
func CloneDefaultSLogger(name string, args ...interface{}) *zap.SugaredLogger {
	copy := *slogger
	cslogger := &copy
	cslogger = cslogger.Named(name)
	if len(args) > 0 {
		cslogger = cslogger.With(args...)
	}
	return cslogger
}

// AttachCore add a core to zap logger
func AttachCore(l *zap.Logger, c zapcore.Core) *zap.Logger {
	return l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, c)
	}))
}

// DefaultInitialFields return defaultInitialFields
func DefaultInitialFields() map[string]interface{} {
	return defaultInitialFields
}

// DefaultAtomicLevelAddr return defaultAtomicLevelAddr
func DefaultAtomicLevelAddr() string {
	return defaultAtomicLevelAddr
}
