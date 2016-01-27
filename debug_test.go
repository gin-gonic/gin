// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO
// func debugRoute(httpMethod, absolutePath string, handlers HandlersChain) {
// func debugPrint(format string, values ...interface{}) {

func TestIsDebugging(t *testing.T) {
	SetMode(DebugMode)
	assert.True(t, IsDebugging())
	SetMode(ReleaseMode)
	assert.False(t, IsDebugging())
	SetMode(TestMode)
	assert.False(t, IsDebugging())
}

func TestDebugPrint(t *testing.T) {
	var w bytes.Buffer
	setup(&w)
	defer teardown()

	SetMode(ReleaseMode)
	debugPrint("DEBUG this!")
	SetMode(TestMode)
	debugPrint("DEBUG this!")
	assert.Empty(t, w.String())

	SetMode(DebugMode)
	debugPrint("these are %d %s\n", 2, "error messages")
	assert.Equal(t, w.String(), "[GIN-debug] these are 2 error messages\n")
}

func TestDebugPrintError(t *testing.T) {
	var w bytes.Buffer
	setup(&w)
	defer teardown()

	SetMode(DebugMode)
	debugPrintError(nil)
	assert.Empty(t, w.String())

	debugPrintError(errors.New("this is an error"))
	assert.Equal(t, w.String(), "[GIN-debug] [ERROR] this is an error\n")
}

func TestDebugPrintRoutes(t *testing.T) {
	var w bytes.Buffer
	setup(&w)
	defer teardown()

	debugPrintRoute("GET", "/path/to/route/:param", HandlersChain{func(c *Context) {}, handlerNameTest})
	assert.Regexp(t, `^\[GIN-debug\] GET    /path/to/route/:param     --> (.*/vendor/)?github.com/gin-gonic/gin.handlerNameTest \(2 handlers\)\n$`, w.String())
}

func setup(w io.Writer) {
	SetMode(DebugMode)
	log.SetOutput(w)
}

func teardown() {
	SetMode(TestMode)
	log.SetOutput(os.Stdout)
}
