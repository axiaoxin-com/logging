package logging

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	if logger == nil {
		t.Error("logger is nil")
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
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	options := Options{
		Name:              "tlogger",
		Level:             level,
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

func TestSetLevel(t *testing.T) {
	logger.Debug("TestChangeLevel raw debug level")
	t.Log("current level:", defaultLoggerLevel.Level())
	defaultLoggerLevel.SetLevel(zap.InfoLevel)
	t.Log("new level:", defaultLoggerLevel.Level())
	logger.Debug("TestChangeLevel raw debug level should not be logged")
	// reset
	defaultLoggerLevel.SetLevel(zap.DebugLevel)
}

func TestHTTPSetLevel(t *testing.T) {
	// query level
	url := "http://localhost" + defaultAtomicLevelAddr
	logger.Debug("TestChangeLevel raw debug level")
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	t.Log("current level:", string(content))

	// set level
	c := &http.Client{}
	req, _ := http.NewRequest("PUT", url, strings.NewReader(`{"level": "info"}`))
	resp, err = c.Do(req)
	if err != nil {
		t.Error(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	t.Log("current level:", string(content))

	logger.Debug("TestChangeLevel raw debug level should not be logged")
}
