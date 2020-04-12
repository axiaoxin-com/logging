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
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	Context(c, "1234")
	_, exists := c.Get(CtxLoggerKey)
	if !exists {
		t.Fatal("set ctxLogger failed")
	}

	if tid := CtxTraceID(c); tid != "1234" {
		t.Fatal("invalid tid", tid)
	}
}

func TestContext(t *testing.T) {
	c := context.Background()

	c = Context(c, "1234")
	if tid := CtxTraceID(c); tid != "1234" {
		t.Fatal("invalid tid", c, tid)
	}
}

func TestGinCtxLoggerEmpty(t *testing.T) {
	c := &gin.Context{}

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

func TestGinCtxLoggerEmptyField(t *testing.T) {
	c := &gin.Context{}

	logger := CtxLogger(c, zap.String("field1", "1"))
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx but with field")
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

	Context(c, "rid")
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
	logger.Info("this is a logger from default logger")
}

func TestCtxLoggerDefaultLogger(t *testing.T) {
	c := context.Background()

	Context(c, "rid")
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
}

func TestGinCtxLoggerDefaultLoggerWithField(t *testing.T) {
	c := &gin.Context{}

	Context(c, "rid")
	ctxlogger := CtxLogger(c, zap.String("myfield", "xxx"))
	if ctxlogger == nil {
		t.Fatal("context also must should return a logger")
	}
	if ctxlogger == logger {
		t.Fatal("with field will get a logger")
	}
	logger.Info("this is a logger from default logger with field")
}

func TestCtxLoggerDefaultLoggerWithField(t *testing.T) {
	c := context.Background()

	Context(c, "rid")
	ctxlogger := CtxLogger(c, zap.String("myfield", "xxx"))
	if ctxlogger == nil {
		t.Fatal("context also must should return a logger")
	}
	if ctxlogger == logger {
		t.Fatal("with field will get a logger")
	}
	logger.Info("this is a logger from default logger with field")
}

func TestGinCtxTraceID(t *testing.T) {
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	if CtxTraceID(c) == "" {
		t.Fatal("context should return default value")
	}
	c.Set(TraceIDKey, "IAMAREQUESTID")
	if CtxTraceID(c) != "IAMAREQUESTID" {
		t.Fatal("context should return set value")
	}
}

func TestCtxTraceID(t *testing.T) {
	c := context.Background()
	if CtxTraceID(c) == "" {
		t.Fatal("context should return default value")
	}
	c = context.WithValue(c, TraceIDKey, "IAMAREQUESTID")
	if CtxTraceID(c) != "IAMAREQUESTID" {
		t.Fatal("context should return set value")
	}
}
