// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"

	"testing"

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
		signature += "X"
	})
	router.NoMethod(func(c *Context) {
		signature += "X"
	})
	// RUN
	w := performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, signature, "ACDB")
}

// TestBadAbortHandlersChain - ensure that Abort after switch context will not interrupt pending handlers
func TestMiddlewareNextOrder(t *testing.T) {
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
	// RUN
	w := performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, w.Code, 404)
	assert.Equal(t, signature, "ACEGHFDB")
}

// TestAbortHandlersChain - ensure that Abort interrupt used middlewares in fifo order
func TestMiddlewareAbortHandlersChain(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.AbortWithStatus(409)
		c.Next()
		signature += "D"
	})
	router.GET("/", func(c *Context) {
		signature += "D"
		c.Next()
		signature += "E"
	})

	// RUN
	w := performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, w.Code, 409)
	assert.Equal(t, signature, "ACD")
}

func TestMiddlewareAbortHandlersChainAndNext(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.AbortWithStatus(410)
		c.Next()
		signature += "B"

	})
	router.GET("/", func(c *Context) {
		signature += "C"
		c.Next()
	})
	// RUN
	w := performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, w.Code, 410)
	assert.Equal(t, signature, "AB")
}

// TestFailHandlersChain - ensure that Fail interrupt used middlewares in fifo order as
// as well as Abort
func TestMiddlewareFailHandlersChain(t *testing.T) {
	// SETUP
	signature := ""
	router := New()
	router.Use(func(context *Context) {
		signature += "A"
		context.Fail(500, errors.New("foo"))
	})
	router.Use(func(context *Context) {
		signature += "B"
		context.Next()
		signature += "C"
	})
	// RUN
	w := performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, w.Code, 500)
	assert.Equal(t, signature, "A")
}
