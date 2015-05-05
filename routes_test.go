// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
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

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestRouteParamsByName(t *testing.T) {
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

// TestHandleStaticFile - ensure the static file handles properly
func TestRouteStaticFile(t *testing.T) {
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
func TestRouteStaticDir(t *testing.T) {
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
func TestRouteHeadToDir(t *testing.T) {
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
