// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO
// func (engine *Engine) LoadHTMLGlob(pattern string) {
// func (engine *Engine) LoadHTMLFiles(files ...string) {
// func (engine *Engine) Run(addr string) error {
// func (engine *Engine) RunTLS(addr string, cert string, key string) error {

func init() {
	SetMode(TestMode)
}

func TestCreateEngine(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.BasePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)

	assert.Panics(t, func() { router.addRoute("", "/", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "a", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "/", HandlersChain{}) })
}

func TestCreateDefaultRouter(t *testing.T) {
	router := Default()
	assert.Len(t, router.Handlers, 2)
}

func TestNoRouteWithoutGlobalHandlers(t *testing.T) {
	middleware0 := func(c *Context) {}
	middleware1 := func(c *Context) {}

	router := New()

	router.NoRoute(middleware0)
	assert.Nil(t, router.Handlers)
	assert.Len(t, router.noRoute, 1)
	assert.Len(t, router.allNoRoute, 1)
	assert.Equal(t, router.noRoute[0], middleware0)
	assert.Equal(t, router.allNoRoute[0], middleware0)

	router.NoRoute(middleware1, middleware0)
	assert.Len(t, router.noRoute, 2)
	assert.Len(t, router.allNoRoute, 2)
	assert.Equal(t, router.noRoute[0], middleware1)
	assert.Equal(t, router.allNoRoute[0], middleware1)
	assert.Equal(t, router.noRoute[1], middleware0)
	assert.Equal(t, router.allNoRoute[1], middleware0)
}

func TestNoRouteWithGlobalHandlers(t *testing.T) {
	middleware0 := func(c *Context) {}
	middleware1 := func(c *Context) {}
	middleware2 := func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoRoute(middleware0)
	assert.Len(t, router.allNoRoute, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noRoute, 1)

	assert.Equal(t, router.Handlers[0], middleware2)
	assert.Equal(t, router.noRoute[0], middleware0)
	assert.Equal(t, router.allNoRoute[0], middleware2)
	assert.Equal(t, router.allNoRoute[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoRoute, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noRoute, 1)

	assert.Equal(t, router.Handlers[0], middleware2)
	assert.Equal(t, router.Handlers[1], middleware1)
	assert.Equal(t, router.noRoute[0], middleware0)
	assert.Equal(t, router.allNoRoute[0], middleware2)
	assert.Equal(t, router.allNoRoute[1], middleware1)
	assert.Equal(t, router.allNoRoute[2], middleware0)
}

func TestNoMethodWithoutGlobalHandlers(t *testing.T) {
	middleware0 := func(c *Context) {}
	middleware1 := func(c *Context) {}

	router := New()

	router.NoMethod(middleware0)
	assert.Empty(t, router.Handlers)
	assert.Len(t, router.noMethod, 1)
	assert.Len(t, router.allNoMethod, 1)
	assert.Equal(t, router.noMethod[0], middleware0)
	assert.Equal(t, router.allNoMethod[0], middleware0)

	router.NoMethod(middleware1, middleware0)
	assert.Len(t, router.noMethod, 2)
	assert.Len(t, router.allNoMethod, 2)
	assert.Equal(t, router.noMethod[0], middleware1)
	assert.Equal(t, router.allNoMethod[0], middleware1)
	assert.Equal(t, router.noMethod[1], middleware0)
	assert.Equal(t, router.allNoMethod[1], middleware0)
}

func TestRebuild404Handlers(t *testing.T) {

}

func TestNoMethodWithGlobalHandlers(t *testing.T) {
	middleware0 := func(c *Context) {}
	middleware1 := func(c *Context) {}
	middleware2 := func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoMethod(middleware0)
	assert.Len(t, router.allNoMethod, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noMethod, 1)

	assert.Equal(t, router.Handlers[0], middleware2)
	assert.Equal(t, router.noMethod[0], middleware0)
	assert.Equal(t, router.allNoMethod[0], middleware2)
	assert.Equal(t, router.allNoMethod[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoMethod, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noMethod, 1)

	assert.Equal(t, router.Handlers[0], middleware2)
	assert.Equal(t, router.Handlers[1], middleware1)
	assert.Equal(t, router.noMethod[0], middleware0)
	assert.Equal(t, router.allNoMethod[0], middleware2)
	assert.Equal(t, router.allNoMethod[1], middleware1)
	assert.Equal(t, router.allNoMethod[2], middleware0)
}
