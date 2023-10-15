// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"flag"
	"os"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv(EnvGinMode, TestMode)
}

func TestSetMode(t *testing.T) {
	assert.Equal(t, testCode, ginMode)
	assert.Equal(t, TestMode, Mode())
	os.Unsetenv(EnvGinMode)

	err := SetMode("")
	assert.NoError(t, err)
	assert.Equal(t, testCode, ginMode)
	assert.Equal(t, TestMode, Mode())

	tmp := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	err = SetMode("")
	assert.NoError(t, err)
	assert.Equal(t, debugCode, ginMode)
	assert.Equal(t, DebugMode, Mode())
	flag.CommandLine = tmp

	err = SetMode(DebugMode)
	assert.NoError(t, err)
	assert.Equal(t, debugCode, ginMode)
	assert.Equal(t, DebugMode, Mode())

	err = SetMode(ReleaseMode)
	assert.NoError(t, err)
	assert.Equal(t, releaseCode, ginMode)
	assert.Equal(t, ReleaseMode, Mode())

	err = SetMode(TestMode)
	assert.NoError(t, err)
	assert.Equal(t, testCode, ginMode)
	assert.Equal(t, TestMode, Mode())

	err = SetMode("unknown")
	assert.Error(t, err)
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
