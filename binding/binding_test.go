// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin/binding/example"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar" xml:"bar" binding:"required"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, Default("GET", ""), Form)
	assert.Equal(t, Default("GET", MIMEJSON), Form)

	assert.Equal(t, Default("POST", MIMEJSON), JSON)
	assert.Equal(t, Default("PUT", MIMEJSON), JSON)

	assert.Equal(t, Default("POST", MIMEXML), XML)
	assert.Equal(t, Default("PUT", MIMEXML2), XML)

	assert.Equal(t, Default("POST", MIMEPOSTForm), Form)
	assert.Equal(t, Default("PUT", MIMEPOSTForm), Form)

	assert.Equal(t, Default("POST", MIMEMultipartPOSTForm), Form)
	assert.Equal(t, Default("PUT", MIMEMultipartPOSTForm), Form)

	assert.Equal(t, Default("POST", MIMEPROTOBUF), ProtoBuf)
	assert.Equal(t, Default("PUT", MIMEPROTOBUF), ProtoBuf)

	assert.Equal(t, Default("POST", MIMEMSGPACK), MsgPack)
	assert.Equal(t, Default("PUT", MIMEMSGPACK2), MsgPack)
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
		"foo=bar&bar=foo", "bar2=foo")
}

func TestBindingForm2(t *testing.T) {
	testFormBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"", "")
}

func TestBindingXML(t *testing.T) {
	testBodyBinding(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func createFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	mw.SetBoundary(boundary)
	mw.WriteField("foo", "bar")
	mw.WriteField("bar", "foo")
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest()
	var obj FooBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest()
	var obj FooBarStruct
	FormMultipart.Bind(req, &obj)

	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestBindingProtoBuf(t *testing.T) {
	test := &example.Test{
		Label: proto.String("yes"),
	}
	data, _ := proto.Marshal(test)

	testProtoBodyBinding(t,
		ProtoBuf, "protobuf",
		"/", "/",
		string(data), string(data[1:]))
}

func TestBindingMsgPack(t *testing.T) {
	test := FooStruct{
		Foo: "bar",
	}

	h := new(codec.MsgpackHandle)
	assert.NotNil(t, h)
	buf := bytes.NewBuffer([]byte{})
	assert.NotNil(t, buf)
	err := codec.NewEncoder(buf, h).Encode(test)
	assert.NoError(t, err)

	data := buf.Bytes()

	testMsgPackBodyBinding(t,
		MsgPack, "msgpack",
		"/", "/",
		string(data), string(data[1:]))
}

func TestValidationFails(t *testing.T) {
	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func TestValidationDisabled(t *testing.T) {
	backup := Validator
	Validator = nil
	defer func() { Validator = backup }()

	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	assert.NoError(t, err)
}

func TestExistsSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"exists"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := JSON.Bind(req, &obj)
	assert.NoError(t, err)
}

func TestExistsFails(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"foo" binding:"exists"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")

	obj = FooBarStruct{}
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

func testProtoBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := example.Test{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, *obj.Label, "yes")

	obj = example.Test{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err = ProtoBuf.Bind(req, &obj)
	assert.Error(t, err)
}

func testMsgPackBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", MIMEMSGPACK)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", MIMEMSGPACK)
	err = MsgPack.Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
