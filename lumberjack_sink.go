// 使用zap.RegisterSink函数和Config.OutputPaths字段添加自定义日志目标。
// RegisterSink将URL方案映射到Sink构造函数，OutputPaths配置日志目的地（编码为URL）。
// *lumberjack.Logger已经实现了几乎所有的zap.Sink接口。只缺少Sync方法。

package logging

import (
	"net/url"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

const (
	// DefaultLogFilename 默认日志文件名
	DefaultLogFilename = "/tmp/logging.log"
)

// LumberjackSink 将日志输出到lumberjack进行rotate
type LumberjackSink struct {
	*lumberjack.Logger
	Scheme string
}

// Sync lumberjack Logger默认已实现Sink的其他方法，这里实现Sync后就成为一个Sink对象
func (LumberjackSink) Sync() error {
	return nil
}

// RegisterLumberjackSink 注册lumberjack sink
// 在OutputPaths中指定输出为sink.Scheme://log_filename即可使用
// path url中不指定日志文件名则使用默认的名称
// 一个scheme只能对应一个文件名，相同的scheme注册无效，会全部写入同一个文件
func RegisterLumberjackSink(sink *LumberjackSink) error {
	err := zap.RegisterSink(sink.Scheme, func(*url.URL) (zap.Sink, error) {
		if sink.Filename == "" {
			sink.Filename = DefaultLogFilename
		}
		return sink, nil
	})
	return err
}

// NewLumberjackSink 创建LumberjackSink对象
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
