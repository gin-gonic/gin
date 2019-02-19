// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer))
	router.GET("/example", func(c *Context) {})
	router.POST("/example", func(c *Context) {})
	router.PUT("/example", func(c *Context) {})
	router.DELETE("/example", func(c *Context) {})
	router.PATCH("/example", func(c *Context) {})
	router.HEAD("/example", func(c *Context) {})
	router.OPTIONS("/example", func(c *Context) {})

	performRequest(router, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// I wrote these first (extending the above) but then realized they are more
	// like integration tests because they test the whole logging process rather
	// than individual functions.  Im not sure where these should go.
	buffer.Reset()
	performRequest(router, "POST", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithConfig(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer}))
	router.GET("/example", func(c *Context) {})
	router.POST("/example", func(c *Context) {})
	router.PUT("/example", func(c *Context) {})
	router.DELETE("/example", func(c *Context) {})
	router.PATCH("/example", func(c *Context) {})
	router.HEAD("/example", func(c *Context) {})
	router.OPTIONS("/example", func(c *Context) {})

	performRequest(router, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// I wrote these first (extending the above) but then realized they are more
	// like integration tests because they test the whole logging process rather
	// than individual functions.  Im not sure where these should go.
	buffer.Reset()
	performRequest(router, "POST", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithFormatter(t *testing.T) {
	buffer := new(bytes.Buffer)

	d := DefaultWriter
	DefaultWriter = buffer
	defer func() {
		DefaultWriter = d
	}()

	router := New()
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %s\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.GET("/example", func(c *Context) {})
	performRequest(router, "GET", "/example?a=100")

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")
}

func TestLoggerWithConfigFormatting(t *testing.T) {
	var gotParam LogFormatterParams
	buffer := new(bytes.Buffer)

	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
		Formatter: func(param LogFormatterParams) string {
			// for assert test
			gotParam = param

			return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %s\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		},
	}))
	router.GET("/example", func(c *Context) {
		// set dummy ClientIP
		c.Request.Header.Set("X-Forwarded-For", "20.20.20.20")
	})
	performRequest(router, "GET", "/example?a=100")

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// LogFormatterParams test
	assert.NotNil(t, gotParam.Request)
	assert.NotEmpty(t, gotParam.TimeStamp)
	assert.Equal(t, 200, gotParam.StatusCode)
	assert.NotEmpty(t, gotParam.Latency)
	assert.Equal(t, "20.20.20.20", gotParam.ClientIP)
	assert.Equal(t, "GET", gotParam.Method)
	assert.Equal(t, "/example?a=100", gotParam.Path)
	assert.Empty(t, gotParam.ErrorMessage)

}

func TestDefaultLogFormatter(t *testing.T) {
	timeStamp := time.Unix(1544173902, 0).UTC()

	termFalseParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Second * 5,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		IsTerm:       false,
	}

	termTrueParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Second * 5,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		IsTerm:       true,
	}

	assert.Equal(t, "[GIN] 2018/12/07 - 09:11:42 | 200 |            5s |     20.20.20.20 | GET      /\n", defaultLogFormatter(termFalseParam))

	assert.Equal(t, "[GIN] 2018/12/07 - 09:11:42 |\x1b[97;42m 200 \x1b[0m|            5s |     20.20.20.20 |\x1b[97;44m GET     \x1b[0m /\n", defaultLogFormatter(termTrueParam))
}

func TestColorForMethod(t *testing.T) {
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 52, 109}), colorForMethod("GET"), "get should be blue")
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 54, 109}), colorForMethod("POST"), "post should be cyan")
	assert.Equal(t, string([]byte{27, 91, 57, 48, 59, 52, 51, 109}), colorForMethod("PUT"), "put should be yellow")
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 49, 109}), colorForMethod("DELETE"), "delete should be red")
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 50, 109}), colorForMethod("PATCH"), "patch should be green")
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 53, 109}), colorForMethod("HEAD"), "head should be magenta")
	assert.Equal(t, string([]byte{27, 91, 57, 48, 59, 52, 55, 109}), colorForMethod("OPTIONS"), "options should be white")
	assert.Equal(t, string([]byte{27, 91, 48, 109}), colorForMethod("TRACE"), "trace is not defined and should be the reset color")
}

func TestColorForStatus(t *testing.T) {
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 50, 109}), colorForStatus(http.StatusOK), "2xx should be green")
	assert.Equal(t, string([]byte{27, 91, 57, 48, 59, 52, 55, 109}), colorForStatus(http.StatusMovedPermanently), "3xx should be white")
	assert.Equal(t, string([]byte{27, 91, 57, 48, 59, 52, 51, 109}), colorForStatus(http.StatusNotFound), "4xx should be yellow")
	assert.Equal(t, string([]byte{27, 91, 57, 55, 59, 52, 49, 109}), colorForStatus(2), "other things should be red")
}

func TestErrorLogger(t *testing.T) {
	router := New()
	router.Use(ErrorLogger())
	router.GET("/error", func(c *Context) {
		c.Error(errors.New("this is an error")) // nolint: errcheck
	})
	router.GET("/abort", func(c *Context) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("no authorized")) // nolint: errcheck
	})
	router.GET("/print", func(c *Context) {
		c.Error(errors.New("this is an error")) // nolint: errcheck
		c.String(http.StatusInternalServerError, "hola!")
	})

	w := performRequest(router, "GET", "/error")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"error\":\"this is an error\"}", w.Body.String())

	w = performRequest(router, "GET", "/abort")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"error\":\"no authorized\"}", w.Body.String())

	w = performRequest(router, "GET", "/print")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "hola!{\"error\":\"this is an error\"}", w.Body.String())
}

func TestLoggerWithWriterSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer, "/skipped"))
	router.GET("/logged", func(c *Context) {})
	router.GET("/skipped", func(c *Context) {})

	performRequest(router, "GET", "/logged")
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	performRequest(router, "GET", "/skipped")
	assert.Contains(t, buffer.String(), "")
}

func TestLoggerWithConfigSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output:    buffer,
		SkipPaths: []string{"/skipped"},
	}))
	router.GET("/logged", func(c *Context) {})
	router.GET("/skipped", func(c *Context) {})

	performRequest(router, "GET", "/logged")
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	performRequest(router, "GET", "/skipped")
	assert.Contains(t, buffer.String(), "")
}

func TestDisableConsoleColor(t *testing.T) {
	New()
	assert.False(t, disableColor)
	DisableConsoleColor()
	assert.True(t, disableColor)
}
