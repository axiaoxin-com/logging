package logging

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestGinContext(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	NewCtxLogger(c, CloneLogger("test"), "1234")
	_, exists := c.Get(string(CtxLoggerName))
	if !exists {
		t.Fatal("set ctxLogger failed")
	}

	if tid := CtxTraceID(c); tid != "1234" {
		t.Fatal("invalid tid", tid)
	}
}

func TestContext(t *testing.T) {
	c := context.Background()

	c, _ = NewCtxLogger(c, CloneLogger("test"), "1234")
	if tid := CtxTraceID(c); tid != "1234" {
		t.Fatal("invalid tid", c, tid)
	}
}

func TestGinCtxLoggerEmpty(t *testing.T) {
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx")
}

func TestCtxLoggerEmpty(t *testing.T) {
	c := context.Background()

	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx")
}

func TestCtxLoggerEmptyField(t *testing.T) {
	c := context.Background()

	logger := CtxLogger(c, zap.String("field1", "1"))
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx but with field")
}

func TestGinCtxLoggerDefaultLogger(t *testing.T) {
	c := &gin.Context{}

	NewCtxLogger(c, CloneLogger("test"), "rid")
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
	logger.Info("this is a logger from default logger")
}

func TestGinCtxLoggerDefaultLoggerWithField(t *testing.T) {
	ginctx := &gin.Context{}

	NewCtxLogger(ginctx, CloneLogger("test"), "rid")
	ginCtxlogger := CtxLogger(ginctx, zap.String("myfield", "xxx"))
	if ginCtxlogger == nil {
		t.Fatal("gin context logger is nil")
	}
	ginCtxlogger.Info("this is a logger from default logger with field")
}

func TestCtxLoggerDefaultLogger(t *testing.T) {
	c := context.Background()

	NewCtxLogger(c, CloneLogger("test"), "rid")
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
}

func TestCtxLoggerDefaultLoggerWithField(t *testing.T) {
	c := context.Background()

	NewCtxLogger(c, CtxLogger(c), "rid-xx")
	ctxlogger := CtxLogger(c, zap.String("field", "xx"))
	if ctxlogger == nil {
		t.Fatal("context logger is nil")
	}
	ctxlogger.Info("this is a logger from default logger with field")
}

func TestGinCtxTraceID(t *testing.T) {
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	if CtxTraceID(c) == "" {
		t.Fatal("context should return default value")
	}
	c.Set(string(TraceIDKeyname), "IAMAREQUESTID")
	if CtxTraceID(c) != "IAMAREQUESTID" {
		t.Fatal("context should return set value")
	}
}

func TestCtxTraceID(t *testing.T) {
	c := context.Background()
	if CtxTraceID(c) == "" {
		t.Fatal("context should return default value")
	}
	c = context.WithValue(c, TraceIDKeyname, "IAMAREQUESTID")
	if CtxTraceID(c) != "IAMAREQUESTID" {
		t.Fatal("context should return set value")
	}
}
