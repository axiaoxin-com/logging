package logging

import (
	"testing"
)

func TestInit(t *testing.T) {
	if Logger == nil {
		t.Error("Logger is nil")
	}
	if SLogger == nil {
		t.Error("SLogger is nil")
	}
}

func TestNewLoggerNoParam(t *testing.T) {
	logger, err := NewLogger(Options{})
	if err != nil {
		t.Error(err)
	}
	if logger == nil {
		t.Error("return a nil logger")
	}
	logger.Debug("TestNewLoggerNoParam Debug")
}

func TestNewLogger(t *testing.T) {
	dsn := "sentrydsn"
	sc, err := GetSentryClientByDSN(dsn, true)
	if err != nil {
		t.Error("TestNewLogger GetSentryClientByDSN err", err)
	}
	options := Options{
		Name:              "tlogger",
		Level:             "debug",
		Format:            "json",
		OutputPaths:       []string{"stderr"},
		InitialFields:     map[string]interface{}{"service_name": "testing"},
		DisableCaller:     false,
		DisableStacktrace: false,
		SentryClient:      sc,
	}
	logger, err := NewLogger(options)
	if err != nil {
		t.Error(err)
	}
	logger.Debug("TestNewLogger Debug")
	logger.Error("TestNewLogger Error")
}

func TestCloneLogger(t *testing.T) {
	logger := CloneLogger("cloned")
	logger.Info("TestCloneLogger Info")
}

func TestCloneSLogger(t *testing.T) {
	logger := CloneSLogger("cloned-slogger")
	logger.Info("TestCloneSLogger Info")
}
