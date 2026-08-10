// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
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

// Test the fix for https://github.com/gin-gonic/gin/issues/4189
func TestMiddlewareNoMethodSkipped(t *testing.T) {
	signature := ""
	router := New()
	router.HandleMethodNotAllowed = true
	router.SkipMethodNotAllowedMiddleware = true
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
	assert.Equal(t, http.MethodPost, w.Header().Get("Allow"))
	assert.Equal(t, "EGHF", signature)
}

// The flag is read when the request is served, so it takes effect even when it
// is set after Use() and NoMethod() have already built the handlers chain.
func TestMiddlewareNoMethodSkippedSetAfterRegistration(t *testing.T) {
	signature := ""
	router := New()
	router.HandleMethodNotAllowed = true
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.NoMethod(func(c *Context) {
		signature += "E"
		c.Next()
		signature += "F"
	})
	router.POST("/", func(c *Context) {
		signature += " XX "
	})
	router.SkipMethodNotAllowedMiddleware = true

	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "EF", signature)
}

func TestMiddlewareNoMethodSkippedWithoutNoMethodHandlers(t *testing.T) {
	signature := ""
	router := New()
	router.HandleMethodNotAllowed = true
	router.SkipMethodNotAllowedMiddleware = true
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.POST("/", func(c *Context) {
		signature += " XX "
	})

	// RUN
	w := PerformRequest(router, http.MethodGet, "/")

	// TEST
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "405 method not allowed", w.Body.String())
	assert.Empty(t, signature)
}

// The flag only covers 405 responses, requests falling through to NoRoute must
// keep running the global middleware.
func TestMiddlewareNoMethodSkippedDoesNotAffectNoRoute(t *testing.T) {
	signature := ""
	router := New()
	router.HandleMethodNotAllowed = true
	router.SkipMethodNotAllowedMiddleware = true
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
		signature += " E "
	})
	router.NoRoute(func(c *Context) {
		signature += " X "
	})
	router.POST("/", func(c *Context) {
		signature += " XX "
	})

	// RUN
	w := PerformRequest(router, http.MethodGet, "/not-registered")

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
