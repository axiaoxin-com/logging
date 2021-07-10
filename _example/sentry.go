package main

import (
	"os"
	"time"

	"github.com/axiaoxin-com/logging"
)

func main() {
	debug := true
	sentryClient, _ := logging.NewSentryClient(os.Getenv("dsn"), debug)
	logger := logging.CloneLogger("sentry")
	logger = logging.SentryAttach(logger, sentryClient)
	logger.Error("hello sentry!")
	time.Sleep(2 * time.Second)
	// Output:
	// [Sentry] 2020/04/15 14:27:40 Integration installed: ContextifyFrames
	// [Sentry] 2020/04/15 14:27:40 Integration installed: Environment
	// [Sentry] 2020/04/15 14:27:40 Integration installed: Modules
	// [Sentry] 2020/04/15 14:27:40 Integration installed: IgnoreErrors
	// [Sentry] 2020/04/15 14:27:40 Integration installed: ContextifyFrames
	// [Sentry] 2020/04/15 14:27:40 Integration installed: Environment
	// [Sentry] 2020/04/15 14:27:40 Integration installed: Modules
	// [Sentry] 2020/04/15 14:27:40 Integration installed: IgnoreErrors
	// {"level":"ERROR","time":"2020-04-15T14:27:40.716083+08:00","logger":"logging.sentry","msg":"hello sentry!","pid":47015}
	// [Sentry] 2020/04/15 14:27:40 ModuleIntegration wasn't able to extract modules: module integration failed
	// [Sentry] 2020/04/15 14:27:40 Sending error event [xxx] to host project: id
}
