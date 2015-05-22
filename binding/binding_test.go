// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FooStruct struct {
	Foo string `json:"foo" form:"foo" xml:"foo" binding:"required"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, Default("GET", ""), Form)
	assert.Equal(t, Default("GET", MIMEJSON), Form)

	assert.Equal(t, Default("POST", MIMEJSON), JSON)
	assert.Equal(t, Default("PUT", MIMEJSON), JSON)

	assert.Equal(t, Default("POST", MIMEXML), XML)
	assert.Equal(t, Default("PUT", MIMEXML2), XML)

	assert.Equal(t, Default("POST", MIMEPOSTForm), Form)
	assert.Equal(t, Default("DELETE", MIMEPOSTForm), Form)
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		JSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingForm(t *testing.T) {
	testFormBinding(t, "POST",
		"/", "/",
		"foo=bar", "bar=foo")
}

func TestBindingForm2(t *testing.T) {
	testFormBinding(t, "GET",
		"/?foo=bar", "/?bar=foo",
		"", "")
}

func TestBindingXML(t *testing.T) {
	testBodyBinding(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "query")

	obj := FooStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")

	obj = FooStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
