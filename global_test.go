package logging

import (
	"testing"

	"go.uber.org/zap"
)

func TestGlobal(t *testing.T) {
	Debug("Debug", zap.Int("k", 1))
	Info("Info")
	Warn("Warn")
	Error("Error")

	Debugs("SDebug")
	Infos("SInfo")
	Warns("SWarn")
	Errors("SError")

	Debugf("%s", "SDebugf")
	Infof("%s", "SInfof")
	Warnf("%s", "SWarnf")
	Errorf("%s", "SErrorf")

	Debugw("SDebugw", "k1", "v1")
	Infow("SInfow", "k1", "v1")
	Warnw("SWarnw", "k1", "v1")
	Errorw("SErrorw", "k1", "v1")
}
