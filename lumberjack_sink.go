// 使用 zap.RegisterSink 函数和 Config.OutputPaths 字段添加自定义日志目标。
// RegisterSink 将 URL 方案映射到 Sink 构造函数， OutputPaths 配置日志目的地（编码为 URL ）。
// *lumberjack.Logger 已经实现了几乎所有的 zap.Sink 接口。只缺少 Sync 方法。

package logging

import (
	"net/url"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

const (
	// LogFilename 默认日志文件名
	LogFilename = "/tmp/logging.log"
)

// LumberjackSink 将日志输出到 lumberjack 进行 rotate
type LumberjackSink struct {
	*lumberjack.Logger
	Scheme string
}

// Sync lumberjack Logger 默认已实现 Sink 的其他方法，这里实现 Sync 后就成为一个 Sink 对象
func (LumberjackSink) Sync() error {
	return nil
}

// RegisterLumberjackSink 注册 lumberjack sink
// 在 OutputPaths 中指定输出为 sink.Scheme://log_filename 即可使用
// path url 中不指定日志文件名则使用默认的名称
// 一个 scheme 只能对应一个文件名，相同的 scheme 注册无效，会全部写入同一个文件
func RegisterLumberjackSink(sink *LumberjackSink) error {
	err := zap.RegisterSink(sink.Scheme, func(*url.URL) (zap.Sink, error) {
		if sink.Filename == "" {
			sink.Filename = LogFilename
		}
		return sink, nil
	})
	return err
}

// NewLumberjackSink 创建 LumberjackSink 对象
func NewLumberjackSink(scheme, filename string, maxAge, maxBackups, maxSize int, compress, localtime bool) *LumberjackSink {
	return &LumberjackSink{
		Logger: &lumberjack.Logger{
			Filename:   filename,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			MaxSize:    maxSize,
			Compress:   compress,
			LocalTime:  localtime,
		},
		Scheme: scheme,
	}
}
