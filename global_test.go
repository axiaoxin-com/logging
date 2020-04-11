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

	SDebug("SDebug")
	SInfo("SInfo")
	SWarn("SWarn")
	SError("SError")

	SDebugf("%s", "SDebugf")
	SInfof("%s", "SInfof")
	SWarnf("%s", "SWarnf")
	SErrorf("%s", "SErrorf")

	SDebugw("SDebugw", "k1", "v1")
	SInfow("SInfow", "k1", "v1")
	SWarnw("SWarnw", "k1", "v1")
	SErrorw("SErrorw", "k1", "v1")
}
