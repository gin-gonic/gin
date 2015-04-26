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
	assert.Equal(t, Default("GET", ""), GETForm)
	assert.Equal(t, Default("GET", MIMEJSON), GETForm)

	assert.Equal(t, Default("POST", MIMEJSON), JSON)
	assert.Equal(t, Default("PUT", MIMEJSON), JSON)

	assert.Equal(t, Default("POST", MIMEXML), XML)
	assert.Equal(t, Default("PUT", MIMEXML2), XML)

	assert.Equal(t, Default("POST", MIMEPOSTForm), POSTForm)
	assert.Equal(t, Default("DELETE", MIMEPOSTForm), POSTForm)
}

func TestBindingJSON(t *testing.T) {
	testBinding(t,
		JSON, MIMEJSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingPOSTForm(t *testing.T) {
	testBinding(t,
		POSTForm, MIMEPOSTForm, "post_form",
		"/", "/",
		"foo=bar", "bar=foo")
}

func TestBindingGETForm(t *testing.T) {
	testBinding(t,
		GETForm, MIMEPOSTForm, "get_form",
		"/?foo=bar", "/?bar=foo",
		"", "")
}

func TestBindingXML(t *testing.T) {
	testBinding(t,
		XML, MIMEXML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func testBinding(t *testing.T, b Binding, contentType, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStruct{}
	req := requestWithBody(path, body, contentType)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")

	obj = FooStruct{}
	req = requestWithBody(badPath, badBody, contentType)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(path, body, contentType string) (req *http.Request) {
	req, _ = http.NewRequest("POST", path, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", contentType)
	return req
}
