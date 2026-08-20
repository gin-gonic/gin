// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package yaml

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	Foo string `yaml:"foo"`
	Bar string `yaml:"bar" binding:"required"`
}

func TestRender(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Render(c, http.StatusOK, payload{Foo: "foo", Bar: "bar"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/yaml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "foo: foo")
	assert.Contains(t, w.Body.String(), "bar: bar")
}

func TestBinding(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("foo: foo\nbar: bar\n"))
	var out payload
	require.NoError(t, Binding.Bind(req, &out))
	assert.Equal(t, "foo", out.Foo)
	assert.Equal(t, "bar", out.Bar)
	assert.Equal(t, "yaml", Binding.Name())
}

func TestBindingValidation(t *testing.T) {
	var out payload
	err := Binding.BindBody([]byte("foo: foo\n"), &out) // Bar missing
	require.Error(t, err)
}

func TestRegistration(t *testing.T) {
	for _, ct := range []string{MIMEYAML, MIMEYAML2} {
		assert.Equal(t, Binding, binding.Default(http.MethodPost, ct), ct)
		r, ok := render.Negotiate(ct, payload{Foo: "x"})
		assert.True(t, ok, ct)
		assert.NotNil(t, r)
	}
}
