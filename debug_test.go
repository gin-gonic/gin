// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO
// func debugRoute(httpMethod, absolutePath string, handlers []HandlerFunc) {
// func debugPrint(format string, values ...interface{}) {

func TestIsDebugging(t *testing.T) {
	SetMode(DebugMode)
	assert.True(t, IsDebugging())
	SetMode(ReleaseMode)
	assert.False(t, IsDebugging())
	SetMode(TestMode)
	assert.False(t, IsDebugging())
}
