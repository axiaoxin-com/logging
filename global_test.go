package logging

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestGlobal(t *testing.T) {
	c, _ := NewCtxLogger(context.Background(), CloneLogger("test"), "xxx-yyy-zzz")
	Debug(nil, "Debug nil", zap.Int("k", 1))
	Debug(c, "Debug", zap.Int("k", 1))
	Info(c, "Info")
	Warn(c, "Warn")
	Error(c, "Error")

	Debugs(c, "Debugs")
	Infos(c, "Infos")
	Warns(c, "Warns")
	Errors(c, "Errors")

	Debugf(c, "%s", "Debugf")
	Infof(c, "%s", "Infof")
	Warnf(c, "%s", "Warnf")
	Errorf(c, "%s", "Errorf")

	Debugw(c, "Debugw", "k1", "v1")
	Infow(c, "Infow", "k1", "v1")
	Warnw(c, "Warnw", "k1", "v1")
	Errorw(c, "Errorw", "k1", "v1")

	f := ExtraField("k1", "v1", "k2", 2)
	t.Logf("%#v", f)
}
