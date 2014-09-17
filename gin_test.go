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
	"strings"
	"testing"
)

func init() {
	SetMode(TestMode)
}

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteOK(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	r.Handle(method, "/test", []HandlerFunc{func(c *Context) {
		passed = true
	}})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == false {
		t.Errorf(method + " route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
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
func testRouteNotOK(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	r.Handle(method, "/test_2", []HandlerFunc{func(c *Context) {
		passed = true
	}})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == true {
		t.Errorf(method + " route handler was invoked, when it should not")
	}
	if w.Code != http.StatusNotFound {
		// If this fails, it's because httprouter needs to be updated to at least f78f58a0db
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
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
func testRouteNotOK2(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	var methodRoute string
	if method == "POST" {
		methodRoute = "GET"
	} else {
		methodRoute = "POST"
	}
	r.Handle(methodRoute, "/test", []HandlerFunc{func(c *Context) {
		passed = true
	}})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == true {
		t.Errorf(method + " route handler was invoked, when it should not")
	}
	if w.Code != http.StatusNotFound {
		// If this fails, it's because httprouter needs to be updated to at least f78f58a0db
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
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
	w := PerformRequest(r, "GET", filePath)

	// TEST
	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}
	if w.Body.String() != "Gin Web Framework" {
		t.Errorf("Response should be test, was: %s", w.Body.String())
	}
	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestHandleStaticDir(t *testing.T) {
	// SETUP
	r := New()
	r.Static("/", "./")

	// RUN
	w := PerformRequest(r, "GET", "/")

	// TEST
	bodyAsString := w.Body.String()
	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}
	if len(bodyAsString) == 0 {
		t.Errorf("Got empty body instead of file tree")
	}
	if !strings.Contains(bodyAsString, "gin.go") {
		t.Errorf("Can't find:`gin.go` in file tree: %s", bodyAsString)
	}
	if w.HeaderMap.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
func TestHandleHeadToDir(t *testing.T) {
	// SETUP
	r := New()
	r.Static("/", "./")

	// RUN
	w := PerformRequest(r, "HEAD", "/")

	// TEST
	bodyAsString := w.Body.String()
	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}
	if len(bodyAsString) == 0 {
		t.Errorf("Got empty body instead of file tree")
	}
	if !strings.Contains(bodyAsString, "gin.go") {
		t.Errorf("Can't find:`gin.go` in file tree: %s", bodyAsString)
	}
	if w.HeaderMap.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}
