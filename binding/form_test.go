// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormBindings_Names(t *testing.T) {
	assert.Equal(t, "form", Form.Name())
	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
}

func TestFormBinding_BasicPost(t *testing.T) {
	b := Form
	obj := FooBarStruct{}
	req := requestWithBody(http.MethodPost, "/", "foo=bar&bar=foo")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestFormPostBinding_Basic(t *testing.T) {
	b := FormPost
	obj := FooBarStruct{}
	req := requestWithBody(http.MethodPost, "/", "foo=bar&bar=foo")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestFormMultipartBinding_BasicMultipart(t *testing.T) {
	// reuse helper that exists in binding_test.go
	req := createFormMultipartRequest(t)
	var obj FooBarStruct
	err := FormMultipart.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

// ensure that a non-multipart POST with form content still binds with Form (ParseMultipartForm should be ignored)
func TestFormBinding_NonMultipartPostIgnoresMultipartError(t *testing.T) {
	b := Form
	obj := FooBarStruct{}
	req := requestWithBody(http.MethodPost, "/", "foo=bar&bar=foo")
	// Intentionally do not set multipart content-type so ParseMultipartForm returns ErrNotMultipart and is ignored
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}
