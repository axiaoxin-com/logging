package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axiaoxin-com/logging"

	"go.uber.org/zap"
)

func main() {
	/* change log level on fly */

	// 创建指定 Level 的 logger ，并开启 http 服务
	options := logging.Options{
		Format: "json",
		Level:  "debug",
		AtomicLevelServer: &logging.AtomicLevelServerOption{
			Addr:     ":8999",
			Path:     "/level",
			Username: "admin",
			Password: "admin",
		},
		DisableStacktrace: true,
	}
	logger, _ := logging.NewLogger(options)
	// 替换 logging 默认 logger
	logging.ReplaceLogger(logger)
	logging.Debug(nil, "Debug level msg", zap.Any("current level", logging.TextLevel()))

	// 使用 SetLevel 动态修改 logger 日志级别为 error
	// 实际应用中可以监听配置文件中日志级别配置项的变化动态调用该函数
	logging.SetLevel("error")
	// Info 级别将不会被打印
	logging.Info(nil, "--> [FAIL] Info level msg will not be logged")
	// 只会打印 error 以上
	logging.Error(nil, "Error level msg", zap.Any("current level", logging.TextLevel()))

	// 通过 HTTP 方式动态修改当前的 error level 为 info level
	url := "http://localhost" + options.AtomicLevelServer.Addr + options.AtomicLevelServer.Path
	c := &http.Client{}
	req, _ := http.NewRequest("PUT", url, strings.NewReader(`{"level": "info"}`))
	req.SetBasicAuth("admin", "admin")
	resp, _ := c.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(content))

	logging.Debug(nil, "--> [FAILe] debug level will not be logger")

	/* 修改默认 logger 日志级别 */
	logging.Info(nil, "level change to info success")
}
