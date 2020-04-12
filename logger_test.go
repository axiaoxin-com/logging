package logging

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	if logger == nil {
		t.Error("logger is nil")
	}
	if slogger == nil {
		t.Error("slogger is nil")
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
		// t.Error("TestNewLogger GetSentryClientByDSN err", err)
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

func TestCloneDefaultLogger(t *testing.T) {
	nlogger := CloneDefaultLogger("cloned")
	if reflect.DeepEqual(nlogger, logger) {
		t.Error("CloneDefaultLogger should not be default logger")
	}
	if &nlogger == &logger {
		t.Error("CloneDefaultLogger should not be default logger")
	}
}

func TestCloneDefaultSLogger(t *testing.T) {
	nlogger := CloneDefaultSLogger("cloned-slogger")
	nlogger.Info("TestCloneDefaultSLogger Info")
	if reflect.DeepEqual(nlogger, logger) {
		t.Error("CloneDefaultLogger should not be default logger")
	}
	if &nlogger == &logger {
		t.Error("CloneDefaultLogger should not be default logger")
	}
}
