// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"testing"

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

	performRequest(router, "GET", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")

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

func TestColorForMethod(t *testing.T) {
	assert.Equal(t, colorForMethod("GET"), string([]byte{27, 91, 57, 55, 59, 52, 52, 109}), "get should be blue")
	assert.Equal(t, colorForMethod("POST"), string([]byte{27, 91, 57, 55, 59, 52, 54, 109}), "post should be cyan")
	assert.Equal(t, colorForMethod("PUT"), string([]byte{27, 91, 57, 55, 59, 52, 51, 109}), "put should be yellow")
	assert.Equal(t, colorForMethod("DELETE"), string([]byte{27, 91, 57, 55, 59, 52, 49, 109}), "delete should be red")
	assert.Equal(t, colorForMethod("PATCH"), string([]byte{27, 91, 57, 55, 59, 52, 50, 109}), "patch should be green")
	assert.Equal(t, colorForMethod("HEAD"), string([]byte{27, 91, 57, 55, 59, 52, 53, 109}), "head should be magenta")
	assert.Equal(t, colorForMethod("OPTIONS"), string([]byte{27, 91, 57, 48, 59, 52, 55, 109}), "options should be white")
	assert.Equal(t, colorForMethod("TRACE"), string([]byte{27, 91, 48, 109}), "trace is not defined and should be the reset color")
}

func TestColorForStatus(t *testing.T) {
	assert.Equal(t, colorForStatus(200), string([]byte{27, 91, 57, 55, 59, 52, 50, 109}), "2xx should be green")
	assert.Equal(t, colorForStatus(301), string([]byte{27, 91, 57, 48, 59, 52, 55, 109}), "3xx should be white")
	assert.Equal(t, colorForStatus(404), string([]byte{27, 91, 57, 55, 59, 52, 51, 109}), "4xx should be yellow")
	assert.Equal(t, colorForStatus(2), string([]byte{27, 91, 57, 55, 59, 52, 49, 109}), "other things should be red")
}

func TestErrorLogger(t *testing.T) {
	router := New()
	router.Use(ErrorLogger())
	router.GET("/error", func(c *Context) {
		c.Error(errors.New("this is an error"))
	})
	router.GET("/abort", func(c *Context) {
		c.AbortWithError(401, errors.New("no authorized"))
	})
	router.GET("/print", func(c *Context) {
		c.Error(errors.New("this is an error"))
		c.String(500, "hola!")
	})

	w := performRequest(router, "GET", "/error")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"error\":\"this is an error\"}", w.Body.String())

	w = performRequest(router, "GET", "/abort")
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"error\":\"no authorized\"}", w.Body.String())

	w = performRequest(router, "GET", "/print")
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "hola!{\"error\":\"this is an error\"}", w.Body.String())
}

func TestSkippingPaths(t *testing.T) {
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

func TestDisableConsoleColor(t *testing.T) {
	New()
	assert.False(t, disableColor)
	DisableConsoleColor()
	assert.True(t, disableColor)
}
