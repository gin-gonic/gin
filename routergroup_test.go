// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestRouterGroupBasic(t *testing.T) {
	router := New()
	group := router.Group("/hola", func(c *Context) {})
	group.Use(func(c *Context) {})

	assert.Len(t, group.Handlers, 2)
	assert.Equal(t, group.BasePath, "/hola")
	assert.Equal(t, group.engine, router)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, group2.BasePath, "/hola/manu")
	assert.Equal(t, group2.engine, router)
}

func TestRouterGroupBasicHandle(t *testing.T) {
	performRequestInGroup(t, "GET")
	performRequestInGroup(t, "POST")
	performRequestInGroup(t, "PUT")
	performRequestInGroup(t, "PATCH")
	performRequestInGroup(t, "DELETE")
	performRequestInGroup(t, "HEAD")
	performRequestInGroup(t, "OPTIONS")
}

func performRequestInGroup(t *testing.T, method string) {
	router := New()
	v1 := router.Group("v1", func(c *Context) {})
	assert.Equal(t, v1.BasePath, "/v1")

	login := v1.Group("/login/", func(c *Context) {}, func(c *Context) {})
	assert.Equal(t, login.BasePath, "/v1/login/")

	handler := func(c *Context) {
		c.String(400, "the method was %s and index %d", c.Request.Method, c.index)
	}

	switch method {
	case "GET":
		v1.GET("/test", handler)
		login.GET("/test", handler)
	case "POST":
		v1.POST("/test", handler)
		login.POST("/test", handler)
	case "PUT":
		v1.PUT("/test", handler)
		login.PUT("/test", handler)
	case "PATCH":
		v1.PATCH("/test", handler)
		login.PATCH("/test", handler)
	case "DELETE":
		v1.DELETE("/test", handler)
		login.DELETE("/test", handler)
	case "HEAD":
		v1.HEAD("/test", handler)
		login.HEAD("/test", handler)
	case "OPTIONS":
		v1.OPTIONS("/test", handler)
		login.OPTIONS("/test", handler)
	default:
		panic("unknown method")
	}

	w := performRequest(router, method, "/v1/login/test")
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "the method was "+method+" and index 3")

	w = performRequest(router, method, "/v1/test")
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "the method was "+method+" and index 1")
}

func TestRouterGroupInvalidStatic(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Static("/path/:param", "/")
	})

	assert.Panics(t, func() {
		router.Static("/path/*param", "/")
	})
}

func TestRouterGroupInvalidStaticFile(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.StaticFile("/path/:param", "favicon.ico")
	})

	assert.Panics(t, func() {
		router.StaticFile("/path/*param", "favicon.ico")
	})
}

func TestRouterGroupTooManyHandlers(t *testing.T) {
	router := New()
	handlers1 := make([]HandlerFunc, 40)
	router.Use(handlers1...)

	handlers2 := make([]HandlerFunc, 26)
	assert.Panics(t, func() {
		router.Use(handlers2...)
	})
	assert.Panics(t, func() {
		router.GET("/", handlers2...)
	})
}

func TestRouterGroupBadMethod(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Handle("get", "/")
	})
	assert.Panics(t, func() {
		router.Handle(" GET", "/")
	})
	assert.Panics(t, func() {
		router.Handle("GET ", "/")
	})
	assert.Panics(t, func() {
		router.Handle("", "/")
	})
	assert.Panics(t, func() {
		router.Handle("PO ST", "/")
	})
	assert.Panics(t, func() {
		router.Handle("1GET", "/")
	})
	assert.Panics(t, func() {
		router.Handle("PATCh", "/")
	})
}

func TestRouterGroupPipeline(t *testing.T) {
	router := New()
	testRoutesInterface(t, router)

	v1 := router.Group("/v1")
	testRoutesInterface(t, v1)
}

func testRoutesInterface(t *testing.T, r RoutesInterface) {
	handler := func(c *Context) {}
	assert.Equal(t, r, r.Use(handler))

	assert.Equal(t, r, r.Handle("GET", "/handler", handler))
	assert.Equal(t, r, r.Any("/any", handler))
	assert.Equal(t, r, r.GET("/", handler))
	assert.Equal(t, r, r.POST("/", handler))
	assert.Equal(t, r, r.DELETE("/", handler))
	assert.Equal(t, r, r.PATCH("/", handler))
	assert.Equal(t, r, r.PUT("/", handler))
	assert.Equal(t, r, r.OPTIONS("/", handler))
	assert.Equal(t, r, r.HEAD("/", handler))

	assert.Equal(t, r, r.StaticFile("/file", "."))
	assert.Equal(t, r, r.Static("/static", "."))
	assert.Equal(t, r, r.StaticFS("/static2", Dir(".", false)))
}
