package logging

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.JSON(200, "world")
}

func TestGinLogger(t *testing.T) {
	app := gin.New()
	app.Use(GinLogger())
	app.GET("/hello", hello)
	go app.Run()
	time.Sleep(100 * time.Millisecond)

	_, err := http.Get("http://localhost:8080/hello?k=v")
	if err != nil {
		t.Fatal(err)
	}
}
