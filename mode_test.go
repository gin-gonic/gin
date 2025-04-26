// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"os"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv(EnvGinMode, TestMode)
}

func TestSetMode(t *testing.T) {
	assert.Equal(t, int32(testCode), atomic.LoadInt32(&ginMode))
	assert.Equal(t, TestMode, Mode())
	os.Unsetenv(EnvGinMode)

	SetMode("")
	assert.Equal(t, int32(testCode), atomic.LoadInt32(&ginMode))
	assert.Equal(t, TestMode, Mode())

	SetMode(DebugMode)
	assert.Equal(t, int32(debugCode), atomic.LoadInt32(&ginMode))
	assert.Equal(t, DebugMode, Mode())

	SetMode(ReleaseMode)
	assert.Equal(t, int32(releaseCode), atomic.LoadInt32(&ginMode))
	assert.Equal(t, ReleaseMode, Mode())

	SetMode(TestMode)
	assert.Equal(t, int32(testCode), atomic.LoadInt32(&ginMode))
	assert.Equal(t, TestMode, Mode())

	assert.Panics(t, func() { SetMode("unknown") })
}

func TestDisableBindValidation(t *testing.T) {
	v := binding.Validator
	assert.NotNil(t, binding.Validator)
	DisableBindValidation()
	assert.Nil(t, binding.Validator)
	binding.Validator = v
}

func TestEnableJsonDecoderUseNumber(t *testing.T) {
	assert.False(t, binding.EnableDecoderUseNumber)
	EnableJsonDecoderUseNumber()
	assert.True(t, binding.EnableDecoderUseNumber)
}

func TestEnableJsonDecoderDisallowUnknownFields(t *testing.T) {
	assert.False(t, binding.EnableDecoderDisallowUnknownFields)
	EnableJsonDecoderDisallowUnknownFields()
	assert.True(t, binding.EnableDecoderDisallowUnknownFields)
}
