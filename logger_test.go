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
	sc, err := NewSentryClient(dsn, true)
	if err != nil {
		// t.Error("TestNewLogger NewSentryClient err", err)
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
		AtomicLevelServer: AtomicLevelServerOption{Addr: ":1903"},
	}
	logger, err := NewLogger(options)
	if err != nil {
		t.Error(err)
	}
	logger.Debug("TestNewLogger Debug")
	logger.Error("TestNewLogger Error")

	// TEST HTTP Level
	// query level
	url := "http://localhost:1903"
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

func TestCloneLogger(t *testing.T) {
	nlogger := CloneLogger("cloned")
	if reflect.DeepEqual(nlogger, logger) {
		t.Error("CloneLogger should not be default logger")
	}
	if &nlogger == &logger {
		t.Error("CloneLogger should not be default logger")
	}
}

func TestSetLevel(t *testing.T) {
	logger.Debug("TestChangeLevel raw debug level")
	t.Log("current level:", atomicLevel.Level())
	atomicLevel.SetLevel(zap.InfoLevel)
	t.Log("new level:", atomicLevel.Level())
	logger.Debug("TestChangeLevel raw debug level should not be logged")
	// reset
	atomicLevel.SetLevel(zap.DebugLevel)
}

func TestTextLevel(t *testing.T) {
	level := TextLevel()
	if level != "debug" {
		t.Error(level)
	}

}
