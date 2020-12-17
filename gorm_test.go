package logging

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCtxGormLogger(t *testing.T) {
	logger := NewGormLogger(zapcore.InfoLevel, zap.DebugLevel, 5*time.Second)
	if logger == (GormLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
}
