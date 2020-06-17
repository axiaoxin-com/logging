package logging

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestCtxGormLogger(t *testing.T) {
	logger := CtxGormLogger(nil, zapcore.InfoLevel)
	if logger == (GormLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
}
