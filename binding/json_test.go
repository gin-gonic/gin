// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `json:"foo"`
	}
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO"}`), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}

func TestJSONBindingBindBodyMap(t *testing.T) {
	s := make(map[string]string)
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO","hello":"world"}`), &s)
	require.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Equal(t, "FOO", s["foo"])
	assert.Equal(t, "world", s["hello"])
}

func TestJSONBindingBindEmpty(t *testing.T) {
	var s struct {
		Foo string `binding:"required"`
	}
	err := jsonBinding{}.BindBody([]byte(""), &s)
	require.True(t, strings.Contains(err.Error(), "empty body"))
}
