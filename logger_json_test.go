package gin

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var buffer *bytes.Buffer

func init() {
	SetMode(TestMode)
	buffer = new(bytes.Buffer)
}

func TestJsonLogger(t *testing.T) {
	router := New()
	router.Use(JsonLogger(JsonLoggerConfig{
		Output:    buffer,
		IsConsole: false,
		LogColor:  false,
	}))
	router.GET("/example", func(c *Context) {})
	router.POST("/example", func(c *Context) {})
	router.PUT("/example", func(c *Context) {})
	router.DELETE("/example", func(c *Context) {})
	router.PATCH("/example", func(c *Context) {})
	router.HEAD("/example", func(c *Context) {})
	router.OPTIONS("/example", func(c *Context) {})

	performRequest(router, "GET", "/example?a=100")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// I wrote these first (extending the above) but then realized they are more
	// like integration tests because they test the whole logging process rather
	// than individual functions.  Im not sure where these should go.
	buffer.Reset()
	performRequest(router, "POST", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestJsonLoggerWithConfig(t *testing.T) {
	router := New()
	router.Use(JsonLoggerWithConfig(JsonLoggerConfig{Output: buffer}))
	router.GET("/example", func(c *Context) {})
	router.POST("/example", func(c *Context) {})
	router.PUT("/example", func(c *Context) {})
	router.DELETE("/example", func(c *Context) {})
	router.PATCH("/example", func(c *Context) {})
	router.HEAD("/example", func(c *Context) {})
	router.OPTIONS("/example", func(c *Context) {})

	performRequest(router, "GET", "/example?a=100")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// I wrote these first (extending the above) but then realized they are more
	// like integration tests because they test the whole logging process rather
	// than individual functions.  Im not sure where these should go.
	buffer.Reset()
	performRequest(router, "POST", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestJsonLoggerConfig_SetFilePath2FileName(t *testing.T) {
	f, _ := os.Create("gin.log")
	conf := &JsonLoggerConfig{Output: f}
	conf.SetFilePath2FileName()
	if conf.logDir != "./" || conf.logName != "gin.log" {
		t.Error("SetFilePath2FileName is failing")
	}
}

func TestJsonLoggerConfig_Rename2File(t *testing.T) {
	conf := JsonLoggerConfig{logFilePath: "gin.log"}
	fileName := conf.Rename2File()
	conf.logFilePath = fileName
	if !conf.IsExist() {
		t.Error("Rename2File is failing")
	}
}

func TestJsonLoggerConfig_DeleteLogFile(t *testing.T) {
	router := New()
	router.Use(JsonLoggerWithConfig(JsonLoggerConfig{Output: buffer}))
	router.GET("/example", func(c *Context) {})
	_, _ = os.Create("gin.log.2018-01-01 01:01:01")
	time.Sleep(time.Second)
	conf := JsonLoggerConfig{LogExpDays: 30, logName: "gin.log", logDir: "./"}
	conf.DeleteLogFile()
}

func TestJsonLoggerConfig_CheckFileSize(t *testing.T) {
	conf := &JsonLoggerConfig{logFilePath: "./logger_json_test.go"}
	if conf.CheckFileSize() == 0 {
		t.Error("CheckFileSize is failing")
	}
}

func TestJsonLoggerConfig_InitLogConfig(t *testing.T) {
	conf := &JsonLoggerConfig{}
	conf.InitLogConfig()
	if conf.LogExpDays != 30 ||
		conf.logLimitNums != 1024*1024*1024 ||
		conf.LogLevel != 0 {
		t.Error("InitLogConfig is failing")
	}
}

func TestJsonLoggerConfig_CheckLogExpDays(t *testing.T) {
	conf := &JsonLoggerConfig{}
	conf.CheckLogExpDays()
	if conf.LogExpDays != 30 {
		t.Error("CheckLogExpDays is failing")
	}
}

func TestJsonLoggerConfig_SetLoglevel(t *testing.T) {
	conf := &JsonLoggerConfig{LogLevel: -2}
	logger = &log.Logger
	conf.SetLoglevel()
	if conf.LogLevel != 0 {
		t.Error("SetLoglevel is failing")
	}
}

func TestJsonLoggerConfig_CheckLogWriteSize(t *testing.T) {
	conf := &JsonLoggerConfig{}
	conf.CheckLogWriteSize()
	if conf.LogWriteSize != 1000 {
		t.Error("SetLogWriteSize is failing")
	}
}

func TestJsonLoggerConfig_SetLogFileSize(t *testing.T) {
	conf := &JsonLoggerConfig{LogLimitSize: "1G"}
	conf.SetLogFileSize()
	if conf.logLimitNums != 1024*1024*1024 {
		t.Error("TestJsonLoggerConfig is failing")
	}

	conf = &JsonLoggerConfig{LogLimitSize: "512MB"}
	conf.SetLogFileSize()
	if conf.logLimitNums != 512*1024*1024 {
		t.Error("TestJsonLoggerConfig is failing")
	}
}

func TestJsonLoggerConfig_IsExist(t *testing.T) {
	conf := &JsonLoggerConfig{logFilePath: "logger_json_test.go"}
	if !conf.IsExist() {
		t.Error("logger_json_test.go is not exist")
	}
}

func TestCreateUuid(t *testing.T) {
	data := JsonLoggerConfig{Caller: true}
	if CreateUuid(data) == "" {
		t.Error("CreateUuid is failing")
	}
}
