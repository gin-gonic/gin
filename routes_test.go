// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func testRouteOK(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	r.Handle(method, "/test", []HandlerFunc{func(c *Context) {
		passed = true
	}})
	// RUN
	w := performRequest(r, method, "/test")

	// TEST
	assert.True(t, passed)
	assert.Equal(t, w.Code, http.StatusOK)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK(method string, t *testing.T) {
	// SETUP
	passed := false
	router := New()
	router.Handle(method, "/test_2", []HandlerFunc{func(c *Context) {
		passed = true
	}})

	// RUN
	w := performRequest(router, method, "/test")

	// TEST
	assert.False(t, passed)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK2(method string, t *testing.T) {
	// SETUP
	passed := false
	router := New()
	var methodRoute string
	if method == "POST" {
		methodRoute = "GET"
	} else {
		methodRoute = "POST"
	}
	router.Handle(methodRoute, "/test", []HandlerFunc{func(c *Context) {
		passed = true
	}})

	// RUN
	w := performRequest(router, method, "/test")

	// TEST
	assert.False(t, passed)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestRouterGroupRouteOK(t *testing.T) {
	testRouteOK("POST", t)
	testRouteOK("DELETE", t)
	testRouteOK("PATCH", t)
	testRouteOK("PUT", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("HEAD", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func TestRouteNotOK(t *testing.T) {
	testRouteNotOK("POST", t)
	testRouteNotOK("DELETE", t)
	testRouteNotOK("PATCH", t)
	testRouteNotOK("PUT", t)
	testRouteNotOK("OPTIONS", t)
	testRouteNotOK("HEAD", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func TestRouteNotOK2(t *testing.T) {
	testRouteNotOK2("POST", t)
	testRouteNotOK2("DELETE", t)
	testRouteNotOK2("PATCH", t)
	testRouteNotOK2("PUT", t)
	testRouteNotOK2("OPTIONS", t)
	testRouteNotOK2("HEAD", t)
}

// TestHandleStaticFile - ensure the static file handles properly
func TestHandleStaticFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	filePath := path.Join("/", path.Base(f.Name()))
	f.WriteString("Gin Web Framework")
	f.Close()

	// SETUP gin
	r := New()
	r.Static("./", testRoot)

	// RUN
	w := performRequest(r, "GET", filePath)

	// TEST
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestHandleStaticDir(t *testing.T) {
	// SETUP
	r := New()
	r.Static("/", "./")

	// RUN
	w := performRequest(r, "GET", "/")

	// TEST
	bodyAsString := w.Body.String()
	assert.Equal(t, w.Code, 200)
	assert.NotEmpty(t, bodyAsString)
	assert.Contains(t, bodyAsString, "gin.go")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
func TestHandleHeadToDir(t *testing.T) {
	// SETUP
	router := New()
	router.Static("/", "./")

	// RUN
	w := performRequest(router, "HEAD", "/")

	// TEST
	bodyAsString := w.Body.String()
	assert.Equal(t, w.Code, 200)
	assert.NotEmpty(t, bodyAsString)
	assert.Contains(t, bodyAsString, "gin.go")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

func TestContextGeneralCase(t *testing.T) {
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
func TestContextNextOrder(t *testing.T) {
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
func TestAbortHandlersChain(t *testing.T) {
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

func TestAbortHandlersChainAndNext(t *testing.T) {
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

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestContextParamsByName(t *testing.T) {
	name := ""
	lastName := ""
	router := New()
	router.GET("/test/:name/:last_name", func(c *Context) {
		name = c.Params.ByName("name")
		lastName = c.Params.ByName("last_name")
	})
	// RUN
	w := performRequest(router, "GET", "/test/john/smith")

	// TEST
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, name, "john")
	assert.Equal(t, lastName, "smith")
}

// TestFailHandlersChain - ensure that Fail interrupt used middlewares in fifo order as
// as well as Abort
func TestFailHandlersChain(t *testing.T) {
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
