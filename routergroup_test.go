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
	assert.Equal(t, group.absolutePath, "/hola")
	assert.Equal(t, group.engine, router)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, group2.absolutePath, "/hola/manu")
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
	performRequestInGroup(t, "LINK")
	performRequestInGroup(t, "UNLINK")

}

func performRequestInGroup(t *testing.T, method string) {
	router := New()
	v1 := router.Group("v1", func(c *Context) {})
	assert.Equal(t, v1.absolutePath, "/v1")

	login := v1.Group("/login/", func(c *Context) {}, func(c *Context) {})
	assert.Equal(t, login.absolutePath, "/v1/login/")

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
	case "LINK":
		v1.LINK("/test", handler)
		login.LINK("/test", handler)
	case "UNLINK":
		v1.UNLINK("/test", handler)
		login.UNLINK("/test", handler)
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
