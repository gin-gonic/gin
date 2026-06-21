// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package msgpack

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugorji/go/codec"
)

type payload struct {
	Foo string `msgpack:"foo" json:"foo"`
	Bar string `msgpack:"bar" json:"bar" binding:"required"`
}

func decode(t *testing.T, b []byte) payload {
	t.Helper()
	var out payload
	var mh codec.MsgpackHandle
	require.NoError(t, codec.NewDecoder(bytes.NewReader(b), &mh).Decode(&out))
	return out
}

func TestRender(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Render(c, http.StatusOK, payload{Foo: "foo", Bar: "bar"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))
	got := decode(t, w.Body.Bytes())
	assert.Equal(t, "foo", got.Foo)
	assert.Equal(t, "bar", got.Bar)
}

func TestBinding(t *testing.T) {
	var mh codec.MsgpackHandle
	var buf bytes.Buffer
	require.NoError(t, codec.NewEncoder(&buf, &mh).Encode(payload{Foo: "foo", Bar: "bar"}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(buf.Bytes()))
	var out payload
	require.NoError(t, Binding.Bind(req, &out))
	assert.Equal(t, "foo", out.Foo)
	assert.Equal(t, "bar", out.Bar)
	assert.Equal(t, "msgpack", Binding.Name())
}

func TestBindingValidation(t *testing.T) {
	var mh codec.MsgpackHandle
	var buf bytes.Buffer
	require.NoError(t, codec.NewEncoder(&buf, &mh).Encode(payload{Foo: "foo"})) // Bar missing

	var out payload
	err := Binding.BindBody(buf.Bytes(), &out)
	require.Error(t, err) // binding:"required" on Bar must fire
}

func TestRegistration(t *testing.T) {
	// init() must have wired both content types into the registries.
	for _, ct := range []string{MIMEMSGPACK, MIMEMSGPACK2} {
		assert.Equal(t, Binding, binding.Default(http.MethodPost, ct), ct)
		r, ok := render.Negotiate(ct, payload{Foo: "x"})
		assert.True(t, ok, ct)
		assert.NotNil(t, r)
	}
}
