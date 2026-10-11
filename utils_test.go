// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	SetMode(TestMode)
}

func BenchmarkParseAccept(b *testing.B) {
	for b.Loop() {
		parseAccept("text/html , application/xhtml+xml,application/xml;q=0.9,  */* ;q=0.8")
	}
}

type testStruct struct {
	T *testing.T
}

func (t *testStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	assert.Equal(t.T, http.MethodPost, req.Method)
	assert.Equal(t.T, "/path", req.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "hello")
}

func TestWrap(t *testing.T) {
	router := New()
	router.POST("/path", WrapH(&testStruct{t}))
	router.GET("/path2", WrapF(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/path2", req.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "hola!")
	}))

	w := PerformRequest(router, http.MethodPost, "/path")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "hello", w.Body.String())

	w = PerformRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, http.StatusBadRequest, w.Code)
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
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple type", "text/html", "text/html"},
		{"type with params", "text/html; charset=utf-8", "text/html"},
		{"type with trailing semicolon", "text/html;", "text/html"},

		{"leading whitespace", " text/html", "text/html"},
		{"trailing whitespace", "text/html ", "text/html"},
		{"leading and trailing whitespace", " text/html ", "text/html"},
		{"space before semicolon", "text/html ; charset=utf-8", "text/html"},
		{"no space before semicolon", "text/html;charset=utf-8", "text/html"},

		{"uppercase type", "TEXT/HTML", "text/html"},
		{"mixed case type", "Text/Html", "text/html"},
		{"uppercase with params", "TEXT/HTML; CHARSET=UTF-8", "text/html"},

		{"valid base with invalid param", `text/html; charset="`, "text/html"},
		{"valid base with malformed param value", "text/html; charset=\xff", "text/html"},
		{"valid base with duplicate params", "application/xml;charset=utf-8;charset=utf-8", "application/xml"},

		{"completely invalid", "\xff", "\xff"},
		{"invalid with semicolon", "; charset=utf-8", ""},

		{"application json", "application/json", "application/json"},
		{"application xml", "application/xml", "application/xml"},
		{"multipart form data", "multipart/form-data; boundary=something", "multipart/form-data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterFlags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFunctionName(t *testing.T) {
	assert.Regexp(t, `^(.*/vendor/)?github.com/gin-gonic/gin.somefunction$`, nameOfFunction(somefunction))
}

func somefunction() {
	// this empty function is used by TestFunctionName()
}

func TestJoinPaths(t *testing.T) {
	assert.Empty(t, joinPaths("", ""))
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
	PerformRequest(router, http.MethodGet, "/?foo=hola&bar=10")
	assert.True(t, called)
	assert.Equal(t, "hola", value.Foo)
	assert.Equal(t, 10, value.Bar)

	called = false
	PerformRequest(router, http.MethodGet, "/?foo=hola&bar=1")
	assert.False(t, called)

	assert.Panics(t, func() {
		Bind(&bindTestStruct{})
	})
}

func TestMarshalXMLforH(t *testing.T) {
	h := H{
		"": "test",
	}
	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	var x xml.StartElement
	e := h.MarshalXML(enc, x)
	assert.Error(t, e)
}

func TestMarshalXMLforHSuccess(t *testing.T) {
	h := H{
		"key1": "value1",
		"key2": 123,
	}
	data, err := xml.Marshal(h)
	require.NoError(t, err)
	assert.Contains(t, string(data), "<key1>value1</key1>")
	assert.Contains(t, string(data), "<key2>123</key2>")
}

func TestIsASCII(t *testing.T) {
	assert.True(t, isASCII("test"))
	assert.False(t, isASCII("🧡💛💚💙💜"))
}

func TestSafeInt8(t *testing.T) {
	assert.Equal(t, int8(100), safeInt8(100))
	assert.Equal(t, int8(math.MaxInt8), safeInt8(int(math.MaxInt8)+123))
}

func TestSafeUint16(t *testing.T) {
	assert.Equal(t, uint16(100), safeUint16(100))
	assert.Equal(t, uint16(math.MaxUint16), safeUint16(int(math.MaxUint16)+123))
}
