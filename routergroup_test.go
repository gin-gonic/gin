// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
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
	assert.Equal(t, "/hola", group.BasePath())
	assert.Equal(t, router, group.engine)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, "/hola/manu", group2.BasePath())
	assert.Equal(t, router, group2.engine)
}

func TestRouterGroupBasicHandle(t *testing.T) {
	performRequestInGroup(t, http.MethodGet)
	performRequestInGroup(t, http.MethodPost)
	performRequestInGroup(t, http.MethodPut)
	performRequestInGroup(t, http.MethodPatch)
	performRequestInGroup(t, http.MethodDelete)
	performRequestInGroup(t, http.MethodHead)
	performRequestInGroup(t, http.MethodOptions)
}

func performRequestInGroup(t *testing.T, method string) {
	router := New()
	v1 := router.Group("v1", func(c *Context) {})
	assert.Equal(t, "/v1", v1.BasePath())

	login := v1.Group("/login/", func(c *Context) {}, func(c *Context) {})
	assert.Equal(t, "/v1/login/", login.BasePath())

	handler := func(c *Context) {
		c.String(http.StatusBadRequest, "the method was %s and index %d", c.Request.Method, c.index)
	}

	switch method {
	case http.MethodGet:
		v1.GET("/test", handler)
		login.GET("/test", handler)
	case http.MethodPost:
		v1.POST("/test", handler)
		login.POST("/test", handler)
	case http.MethodPut:
		v1.PUT("/test", handler)
		login.PUT("/test", handler)
	case http.MethodPatch:
		v1.PATCH("/test", handler)
		login.PATCH("/test", handler)
	case http.MethodDelete:
		v1.DELETE("/test", handler)
		login.DELETE("/test", handler)
	case http.MethodHead:
		v1.HEAD("/test", handler)
		login.HEAD("/test", handler)
	case http.MethodOptions:
		v1.OPTIONS("/test", handler)
		login.OPTIONS("/test", handler)
	default:
		panic("unknown method")
	}

	w := performRequest(router, method, "/v1/login/test")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "the method was "+method+" and index 3", w.Body.String())

	w = performRequest(router, method, "/v1/test")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "the method was "+method+" and index 1", w.Body.String())
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
		router.Handle(http.MethodGet, "/")
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

func testRoutesInterface(t *testing.T, r IRoutes) {
	handler := func(c *Context) {}
	assert.Equal(t, r, r.Use(handler))

	assert.Equal(t, r, r.Handle(http.MethodGet, "/handler", handler))
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
