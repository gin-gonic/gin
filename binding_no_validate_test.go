// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

var (
	_ binding.BindingNoValidate = binding.JSON
	_ binding.BindingNoValidate = binding.XML
	_ binding.BindingNoValidate = binding.Form
	_ binding.BindingNoValidate = binding.Query
	_ binding.BindingNoValidate = binding.FormPost
	_ binding.BindingNoValidate = binding.FormMultipart
	_ binding.BindingNoValidate = binding.ProtoBuf
	_ binding.BindingNoValidate = binding.YAML
	_ binding.BindingNoValidate = binding.Header
	_ binding.BindingNoValidate = binding.Plain
	_ binding.BindingNoValidate = binding.TOML
	_ binding.BindingNoValidate = binding.BSON

	_ binding.BindingBodyNoValidate = binding.JSON
	_ binding.BindingBodyNoValidate = binding.XML
	_ binding.BindingBodyNoValidate = binding.ProtoBuf
	_ binding.BindingBodyNoValidate = binding.YAML
	_ binding.BindingBodyNoValidate = binding.Plain
	_ binding.BindingBodyNoValidate = binding.TOML
	_ binding.BindingBodyNoValidate = binding.BSON

	_ binding.BindingUriNoValidate = binding.Uri
)

type deferredBindingRequest struct {
	ID    string `uri:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Token string `header:"X-Token" binding:"required"`
}

func newNoValidateContext(method, target, body, contentType string) *Context {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(method, target, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", contentType)
	return c
}

func TestShouldBindNoValidateCombinesSources(t *testing.T) {
	c := newNoValidateContext(http.MethodPost, "/users/42", `{"name":"Ada"}`, binding.MIMEJSON)
	c.Params = Params{{Key: "id", Value: "42"}}
	c.Request.Header.Set("X-Token", "secret")

	validator := binding.Validator
	var got deferredBindingRequest
	require.NoError(t, c.ShouldBindUriNoValidate(&got))
	require.NoError(t, c.ShouldBindJSONNoValidate(&got))
	require.NoError(t, c.ShouldBindHeaderNoValidate(&got))
	require.Equal(t, deferredBindingRequest{ID: "42", Name: "Ada", Token: "secret"}, got)
	require.NoError(t, binding.Validator.ValidateStruct(&got))
	require.Same(t, validator, binding.Validator)
}

func TestShouldBindNoValidateLeavesValidationExplicit(t *testing.T) {
	withoutValidation := newNoValidateContext(http.MethodPost, "/", `{"name":"Ada"}`, binding.MIMEJSON)
	var partiallyBound deferredBindingRequest
	require.NoError(t, withoutValidation.ShouldBindNoValidate(&partiallyBound))
	require.Error(t, binding.Validator.ValidateStruct(&partiallyBound))

	withValidation := newNoValidateContext(http.MethodPost, "/", `{"name":"Ada"}`, binding.MIMEJSON)
	var validated deferredBindingRequest
	require.Error(t, withValidation.ShouldBind(&validated))
}

func TestShouldBindNoValidateStillReturnsDecodeErrors(t *testing.T) {
	c := newNoValidateContext(http.MethodPost, "/", `{"name":`, binding.MIMEJSON)
	var got deferredBindingRequest
	require.Error(t, c.ShouldBindJSONNoValidate(&got))
}

func TestShouldBindQueryNoValidateLeavesValidationExplicit(t *testing.T) {
	type queryRequest struct {
		Page int `form:"page" binding:"required,min=1"`
		Size int `form:"size" binding:"required,min=1"`
	}

	c := newNoValidateContext(http.MethodGet, "/?page=3", "", binding.MIMEPOSTForm)
	var got queryRequest
	require.NoError(t, c.ShouldBindQueryNoValidate(&got))
	require.Equal(t, 3, got.Page)
	require.Zero(t, got.Size)
	require.Error(t, binding.Validator.ValidateStruct(&got))
}

func TestShouldBindBodyWithNoValidateReusesBody(t *testing.T) {
	c := newNoValidateContext(http.MethodPost, "/", `{"name":"Ada"}`, binding.MIMEJSON)

	var first, second deferredBindingRequest
	require.NoError(t, c.ShouldBindBodyWithJSONNoValidate(&first))
	require.NoError(t, c.ShouldBindBodyWithJSONNoValidate(&second))
	require.Equal(t, first, second)
	require.Equal(t, "Ada", second.Name)
}

type failingNoValidateBody struct{}

func (failingNoValidateBody) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (failingNoValidateBody) Close() error { return nil }

func TestShouldBindBodyWithNoValidateReturnsReadError(t *testing.T) {
	c := newNoValidateContext(http.MethodPost, "/", "", binding.MIMEJSON)
	c.Request.Body = failingNoValidateBody{}

	var got deferredBindingRequest
	require.EqualError(t, c.ShouldBindBodyWithJSONNoValidate(&got), "read failed")
}

func TestShouldBindNoValidateFormatShortcuts(t *testing.T) {
	type formatRequest struct {
		Name     string `xml:"name" yaml:"name" toml:"name"`
		Required string `xml:"required" yaml:"required" toml:"required" binding:"required"`
	}

	tests := []struct {
		name        string
		body        string
		contentType string
		bind        func(*Context, any) error
		bindBody    func(*Context, any) error
		bodyBinding binding.BindingBodyNoValidate
	}{
		{"XML", `<request><name>Ada</name></request>`, binding.MIMEXML, (*Context).ShouldBindXMLNoValidate, (*Context).ShouldBindBodyWithXMLNoValidate, binding.XML},
		{"YAML", "name: Ada\n", binding.MIMEYAML, (*Context).ShouldBindYAMLNoValidate, (*Context).ShouldBindBodyWithYAMLNoValidate, binding.YAML},
		{"TOML", `name = "Ada"`, binding.MIMETOML, (*Context).ShouldBindTOMLNoValidate, (*Context).ShouldBindBodyWithTOMLNoValidate, binding.TOML},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newNoValidateContext(http.MethodPost, "/", tt.body, tt.contentType)
			var direct formatRequest
			require.NoError(t, tt.bind(c, &direct))
			require.Equal(t, "Ada", direct.Name)
			require.Error(t, binding.Validator.ValidateStruct(&direct))

			c = newNoValidateContext(http.MethodPost, "/", tt.body, tt.contentType)
			var first, second formatRequest
			require.NoError(t, tt.bindBody(c, &first))
			require.NoError(t, c.ShouldBindBodyWithNoValidate(&second, tt.bodyBinding))
			require.Equal(t, first, second)
		})
	}
}

func TestShouldBindNoValidatePlainShortcuts(t *testing.T) {
	c := newNoValidateContext(http.MethodPost, "/", "hello", binding.MIMEPlain)
	var direct string
	require.NoError(t, c.ShouldBindPlainNoValidate(&direct))
	require.Equal(t, "hello", direct)

	c = newNoValidateContext(http.MethodPost, "/", "hello", binding.MIMEPlain)
	var first, second string
	require.NoError(t, c.ShouldBindBodyWithPlainNoValidate(&first))
	require.NoError(t, c.ShouldBindBodyWithNoValidate(&second, binding.Plain))
	require.Equal(t, first, second)
}

type customNoValidateBinding struct{}

func (customNoValidateBinding) Name() string { return "custom" }

func (customNoValidateBinding) Bind(*http.Request, any) error {
	return errors.New("validated binding path called")
}

func (customNoValidateBinding) BindNoValidate(_ *http.Request, obj any) error {
	value, ok := obj.(*string)
	if !ok {
		return errors.New("expected *string")
	}
	*value = "bound"
	return nil
}

func TestShouldBindWithNoValidateUsesCustomBinder(t *testing.T) {
	c := newNoValidateContext(http.MethodGet, "/", "", "")
	var got string
	require.NoError(t, c.ShouldBindWithNoValidate(&got, customNoValidateBinding{}))
	require.Equal(t, "bound", got)
}
