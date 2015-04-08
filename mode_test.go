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

func TestSetMode(t *testing.T) {
	SetMode(DebugMode)
	assert.Equal(t, ginMode, debugCode)
	assert.Equal(t, Mode(), DebugMode)

	SetMode(ReleaseMode)
	assert.Equal(t, ginMode, releaseCode)
	assert.Equal(t, Mode(), ReleaseMode)

	SetMode(TestMode)
	assert.Equal(t, ginMode, testCode)
	assert.Equal(t, Mode(), TestMode)

	assert.Panics(t, func() { SetMode("unknown") })
}
