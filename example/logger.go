package main

import (
	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	/* clone default logger */
	logger := logging.CloneDefaultLogger("subname", zap.String("str_field", "field_value"))
	logger.Debug("CloneDefaultLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:20:36.614438+08:00","logger":"root.subname","msg":"CloneDefaultLogger Debug","pid":54273,"str_field":"field_value"}

	/* clone default sugared logger */
	slogger := logging.CloneDefaultSLogger("subname", "foo", 123, zap.String("str_field", "field_value"))
	slogger.Debug("CloneDefaultSLogger Debug")
	// {"level":"DEBUG","time":"2020-04-13T00:24:41.629175+08:00","logger":"root.subname","msg":"CloneDefaultSLogger Debug","pid":73087,"foo":123,"str_field":"field_value"}

	/* custom logger encoder key name */
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
}
