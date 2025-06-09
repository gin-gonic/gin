// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"
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

	err = jsonBinding{}.BindBody([]byte(`{"foo": "FOO}`), &s)
	assert.Equal(t, ErrInvalidJSON, err)
}

func TestJSONBindingBindBodyMap(t *testing.T) {
	s := make(map[string]string)
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO","hello":"world"}`), &s)
	require.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Equal(t, "FOO", s["foo"])
	assert.Equal(t, "world", s["hello"])

	err = jsonBinding{}.BindBody([]byte(`{"foo": "FOO","hello":"world}`), &s)
	assert.Equal(t, ErrInvalidJSON, err)
}

func TestTestJSONBindingBind(t *testing.T) {
	var s struct {
		Foo string `json:"foo"`
	}
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"foo":"FOO"}`))
	err := jsonBinding{}.Bind(req, &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)

	req, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"foo":"FOO}`))
	err = jsonBinding{}.Bind(req, &s)
	assert.Equal(t, ErrInvalidJSON, err)
}
