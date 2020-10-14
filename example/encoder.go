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
		EncoderConfig: &zapcore.EncoderConfig{
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
