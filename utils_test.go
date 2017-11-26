// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

type testStruct struct {
	T *testing.T
}

func (t *testStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	assert.Equal(t.T, "POST", req.Method)
	assert.Equal(t.T, "/path", req.URL.Path)
	w.WriteHeader(500)
	fmt.Fprint(w, "hello")
}

func TestWrap(t *testing.T) {
	router := New()
	router.POST("/path", WrapH(&testStruct{t}))
	router.GET("/path2", WrapF(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/path2", req.URL.Path)
		w.WriteHeader(400)
		fmt.Fprint(w, "hola!")
	}))

	w := performRequest(router, "POST", "/path")
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "hello", w.Body.String())

	w = performRequest(router, "GET", "/path2")
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "hola!", w.Body.String())
}

func TestLastChar(t *testing.T) {
	assert.Equal(t, uint8('a'), lastChar("hola"))
	assert.Equal(t, uint8('s'), lastChar("adios"))
	assert.Panics(t, func() { lastChar("") })
}

func TestParseAccept(t *testing.T) {
	parts := parseAccept("text/html , application/xhtml+xml,application/xml;q=0.9,  */* ;q=0.8")
	assert.Len(t, parts, 4)
	assert.Equal(t, "text/html", parts[0])
	assert.Equal(t, "application/xhtml+xml", parts[1])
	assert.Equal(t, "application/xml", parts[2])
	assert.Equal(t, "*/*", parts[3])
}

func TestChooseData(t *testing.T) {
	A := "a"
	B := "b"
	assert.Equal(t, A, chooseData(A, B))
	assert.Equal(t, B, chooseData(nil, B))
	assert.Panics(t, func() { chooseData(nil, nil) })
}

func TestFilterFlags(t *testing.T) {
	result := filterFlags("text/html ")
	assert.Equal(t, "text/html", result)

	result = filterFlags("text/html;")
	assert.Equal(t, "text/html", result)
}

func TestFunctionName(t *testing.T) {
	assert.Regexp(t, `^(.*/vendor/)?github.com/gin-gonic/gin.somefunction$`, nameOfFunction(somefunction))
}

func somefunction() {
	// this empty function is used by TestFunctionName()
}

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, "", joinPaths("", ""))
	assert.Equal(t, "/", joinPaths("", "/"))
	assert.Equal(t, "/a", joinPaths("/a", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", "/"))
	assert.Equal(t, "/a/", joinPaths("/a", "/"))
	assert.Equal(t, "/a/hola", joinPaths("/a", "/hola"))
	assert.Equal(t, "/a/hola", joinPaths("/a/", "/hola"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola/"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola//"))
}

type bindTestStruct struct {
	Foo string `form:"foo" binding:"required"`
	Bar int    `form:"bar" binding:"min=4"`
}

func TestBindMiddleware(t *testing.T) {
	var value *bindTestStruct
	var called bool
	router := New()
	router.GET("/", Bind(bindTestStruct{}), func(c *Context) {
		called = true
		value = c.MustGet(BindKey).(*bindTestStruct)
	})
	performRequest(router, "GET", "/?foo=hola&bar=10")
	assert.True(t, called)
	assert.Equal(t, "hola", value.Foo)
	assert.Equal(t, 10, value.Bar)

	called = false
	performRequest(router, "GET", "/?foo=hola&bar=1")
	assert.False(t, called)

	assert.Panics(t, func() {
		Bind(&bindTestStruct{})
	})
}
