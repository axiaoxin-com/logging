package logging

import "testing"

func TestCtxGormLogger(t *testing.T) {
	logger := CtxGormLogger(nil)
	if logger == (GormLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
}
