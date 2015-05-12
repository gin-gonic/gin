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

func TestLastChar(t *testing.T) {
	assert.Equal(t, lastChar("hola"), uint8('a'))
	assert.Equal(t, lastChar("adios"), uint8('s'))
	assert.Panics(t, func() { lastChar("") })
}

func TestParseAccept(t *testing.T) {
	parts := parseAccept("text/html , application/xhtml+xml,application/xml;q=0.9,  */* ;q=0.8")
	assert.Len(t, parts, 4)
	assert.Equal(t, parts[0], "text/html")
	assert.Equal(t, parts[1], "application/xhtml+xml")
	assert.Equal(t, parts[2], "application/xml")
	assert.Equal(t, parts[3], "*/*")
}

func TestChooseData(t *testing.T) {
	A := "a"
	B := "b"
	assert.Equal(t, chooseData(A, B), A)
	assert.Equal(t, chooseData(nil, B), B)
	assert.Panics(t, func() { chooseData(nil, nil) })
}

func TestFilterFlags(t *testing.T) {
	result := filterFlags("text/html ")
	assert.Equal(t, result, "text/html")

	result = filterFlags("text/html;")
	assert.Equal(t, result, "text/html")
}

func TestFunctionName(t *testing.T) {
	assert.Equal(t, nameOfFunction(somefunction), "github.com/gin-gonic/gin.somefunction")
}

func somefunction() {
	// this empty function is used by TestFunctionName()
}

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, joinPaths("", ""), "")
	assert.Equal(t, joinPaths("", "/"), "/")
	assert.Equal(t, joinPaths("/a", ""), "/a")
	assert.Equal(t, joinPaths("/a/", ""), "/a/")
	assert.Equal(t, joinPaths("/a/", "/"), "/a/")
	assert.Equal(t, joinPaths("/a", "/"), "/a/")
	assert.Equal(t, joinPaths("/a", "/hola"), "/a/hola")
	assert.Equal(t, joinPaths("/a/", "/hola"), "/a/hola")
	assert.Equal(t, joinPaths("/a/", "/hola/"), "/a/hola/")
	assert.Equal(t, joinPaths("/a/", "/hola//"), "/a/hola/")
}
