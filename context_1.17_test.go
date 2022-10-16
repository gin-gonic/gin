// Copyright 2021 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build go1.17
// +build go1.17

package gin

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type interceptedWriter struct {
	ResponseWriter
	b *bytes.Buffer
}

func (i interceptedWriter) WriteHeader(code int) {
	i.Header().Del("X-Test")
	i.ResponseWriter.WriteHeader(code)
}

func TestContextFormFileFailed17(t *testing.T) {
	if !isGo117OrGo118() {
		return
	}
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	defer func(mw *multipart.Writer) {
		err := mw.Close()
		if err != nil {
			assert.Error(t, err)
		}
	}(mw)
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	c.engine.MaxMultipartMemory = 8 << 20
	assert.Panics(t, func() {
		f, err := c.FormFile("file")
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}

func TestInterceptedHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := CreateTestContext(w)

	r.Use(func(c *Context) {
		i := interceptedWriter{
			ResponseWriter: c.Writer,
			b:              bytes.NewBuffer(nil),
		}
		c.Writer = i
		c.Next()
		c.Header("X-Test", "overridden")
		c.Writer = i.ResponseWriter
	})
	r.GET("/", func(c *Context) {
		c.Header("X-Test", "original")
		c.Header("X-Test-2", "present")
		c.String(http.StatusOK, "hello world")
	})
	c.Request = httptest.NewRequest("GET", "/", nil)
	r.HandleContext(c)
	// Result() has headers frozen when WriteHeaderNow() has been called
	// Compared to this time, this is when the response headers will be flushed
	// As response is flushed on c.String, the Header cannot be set by the first
	// middleware. Assert this
	assert.Equal(t, "", w.Result().Header.Get("X-Test"))
	assert.Equal(t, "present", w.Result().Header.Get("X-Test-2"))
}

func isGo117OrGo118() bool {
	version := strings.Split(runtime.Version()[2:], ".")
	if len(version) >= 2 {
		x := version[0]
		y := version[1]
		if x == "1" && (y == "17" || y == "18") {
			return true
		}
	}
	return false
}
