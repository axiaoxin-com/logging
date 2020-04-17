package main

import (
	"github.com/axiaoxin-com/logging"
)

// Options 传入 LumberjacSink ，并在 OutputPaths 中添加对应 scheme 就能将日志保存到文件并自动 rotate
func main() {
	// scheme 为 lumberjack ，日志文件为 /tmp/x.log , 保存 7 天，保留 10 份文件，文件大小超过 100M ，使用压缩备份，压缩文件名使用 localtime
	sink := logging.NewLumberjackSink("lumberjack", "/tmp/x.log", 7, 10, 100, true, true)
	options := logging.Options{
		LumberjackSink: sink,
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack:"},
	}
	logger, _ := logging.NewLogger(options)
	logger.Debug("xxx")

	sink2 := logging.NewLumberjackSink("lumberjack2", "/tmp/x2.log", 7, 10, 100, true, true)
	options2 := logging.Options{
		LumberjackSink: sink2,
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack2:"},
	}
	logger2, _ := logging.NewLogger(options2)
	logger2.Debug("yyy")
}
