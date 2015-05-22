// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cachedDebugLogger *log.Logger = nil

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

func setup(w io.Writer) {
	SetMode(DebugMode)
	if cachedDebugLogger == nil {
		cachedDebugLogger = debugLogger
		debugLogger = log.New(w, debugLogger.Prefix(), 0)
	} else {
		panic("setup failed")
	}
}

func teardown() {
	SetMode(TestMode)
	if cachedDebugLogger != nil {
		debugLogger = cachedDebugLogger
		cachedDebugLogger = nil
	} else {
		panic("teardown failed")
	}
}
