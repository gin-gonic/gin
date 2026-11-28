// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-contrib/sse"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareGeneralCase(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
	})
	router.GET("/", func(c *Context) {
		signature += "D"
	})
	router.NoRoute(func(c *Context) {
		signature += " X "
	})
	router.NoMethod(func(c *Context) {
		signature += " XX "
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestMiddlewareNoRoute(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	router.NoRoute(func(c *Context) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *Context) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	router.NoMethod(func(c *Context) {
		signature += " X "
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "ACEGHFDB", signature)
}

func TestMiddlewareNoMethodEnabled(t *testing.T) {
	signature := ""
	router := New()
	router.HandleMethodNotAllowed = true
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.Next()
		signature += "D"
	})
	router.NoMethod(func(c *Context) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *Context) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	router.NoRoute(func(c *Context) {
		signature += " X "
	})
	router.POST("/", func(c *Context) {
		signature += " XX "
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "ACEGHFDB", signature)
}

func TestMiddlewareNoMethodDisabled(t *testing.T) {
	signature := ""
	router := New()

	// NoMethod disabled
	router.HandleMethodNotAllowed = false

	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.Next()
		signature += "D"
	})
	router.NoMethod(func(c *Context) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *Context) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	router.NoRoute(func(c *Context) {
		signature += " X "
	})
	router.POST("/", func(c *Context) {
		signature += " XX "
	})

	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "AC X DB", signature)
}

func TestMiddlewareAbort(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.AbortWithStatus(http.StatusUnauthorized)
		c.Next()
		signature += "D"
	})
	router.GET("/", func(c *Context) {
		signature += " X "
		c.Next()
		signature += " XX "
	})

	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "ACD", signature)
}

func TestMiddlewareAbortHandlersChainAndNext(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		c.AbortWithStatus(http.StatusGone)
		signature += "B"
	})
	router.GET("/", func(c *Context) {
		signature += "C"
		c.Next()
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusGone, w.Code)
	assert.Equal(t, "ACB", signature)
}

// TestMiddlewareFailHandlersChain - ensure that Fail interrupt used middleware in fifo order as
// as well as Abort
func TestMiddlewareFailHandlersChain(t *testing.T) {
	// SETUP
	signature := ""
	router := New()
	router.Use(func(context *Context) {
		signature += "A"
		context.AbortWithError(http.StatusInternalServerError, errors.New("foo")) //nolint: errcheck
	})
	router.Use(func(context *Context) {
		signature += "B"
		context.Next()
		signature += "C"
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "A", signature)
}

func TestMiddlewareWrite(t *testing.T) {
	router := New()
	router.Use(func(c *Context) {
		c.String(http.StatusBadRequest, "hola\n")
	})
	router.Use(func(c *Context) {
		c.XML(http.StatusBadRequest, H{"foo": "bar"})
	})
	router.Use(func(c *Context) {
		c.JSON(http.StatusBadRequest, H{"foo": "bar"})
	})
	router.GET("/", func(c *Context) {
		c.JSON(http.StatusBadRequest, H{"foo": "bar"})
	}, func(c *Context) {
		c.Render(http.StatusBadRequest, sse.Event{
			Event: "test",
			Data:  "message",
		})
	})

	w := PerformRequest(router, http.MethodGet, "/")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, strings.ReplaceAll("hola\n<map><foo>bar</foo></map>{\"foo\":\"bar\"}{\"foo\":\"bar\"}event:test\ndata:message\n\n", " ", ""), strings.ReplaceAll(w.Body.String(), " ", ""))
}

// captureDebugOutput runs fn with gin switched to debug mode and returns
// everything debugPrint writes to DefaultWriter during that window.
func captureDebugOutput(t *testing.T, fn func()) string {
	t.Helper()
	oldMode, oldWriter := Mode(), DefaultWriter
	SetMode(DebugMode)
	buf := new(bytes.Buffer)
	DefaultWriter = buf
	defer func() {
		DefaultWriter = oldWriter
		SetMode(oldMode)
	}()
	fn()
	return buf.String()
}

func TestMultipleResponseWritesWarning(t *testing.T) {
	out := captureDebugOutput(t, func() {
		router := New()
		router.GET("/", func(c *Context) {
			c.String(http.StatusOK, "first\n")
			c.String(http.StatusOK, "second\n")
		})

		w := PerformRequest(router, http.MethodGet, "/")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "first\nsecond\n", w.Body.String())
	})

	// The second render must warn exactly once; the first one must not.
	assert.Equal(t, 1, strings.Count(out, "Response body has already been written"))
}

func TestNoWarningOnHeaderFlushBeforeRender(t *testing.T) {
	out := captureDebugOutput(t, func() {
		router := New()
		router.GET("/", func(c *Context) {
			c.Header("Content-Type", "text/plain")
			c.Writer.Flush() // header-only flush, no body bytes written yet
			c.Data(http.StatusOK, "text/plain", []byte("hello"))
		})

		w := PerformRequest(router, http.MethodGet, "/")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "hello", w.Body.String())
	})

	assert.NotContains(t, out, "Response body has already been written")
}

func TestNoWarningOnSSEStream(t *testing.T) {
	out := captureDebugOutput(t, func() {
		router := New()
		router.GET("/sse", func(c *Context) {
			c.Header("Content-Type", "text/event-stream")
			c.Writer.Flush()
			c.SSEvent("message", "hello")
			c.SSEvent("message", "world")
			c.SSEvent("message", "bye")
		})

		w := PerformRequest(router, http.MethodGet, "/sse")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "data:hello")
		assert.Contains(t, w.Body.String(), "data:bye")
	})

	// Streaming events are legitimate repeated writes; they must stay silent.
	assert.NotContains(t, out, "Response body has already been written")
}
