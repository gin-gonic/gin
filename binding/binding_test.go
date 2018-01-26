// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

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

type FooStructUseNumber struct {
	Foo interface{} `json:"foo" binding:"required"`
}

type FooBarStructForTimeType struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2006-01-02" time_utc:"1" time_location:"Asia/Chongqing"`
	TimeBar time.Time `form:"time_bar" time_format:"2006-01-02" time_utc:"1"`
}

type FooStructForTimeTypeNotFormat struct {
	TimeFoo time.Time `form:"time_foo"`
}

type FooStructForTimeTypeFailFormat struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2017-11-15"`
}

type FooStructForTimeTypeFailLocation struct {
	TimeFoo time.Time `form:"time_foo" time_format:"2006-01-02" time_location:"/asia/chongqing"`
}

type FooStructForMapType struct {
	// Unknown type: not support map
	MapFoo map[string]interface{} `form:"map_foo"`
}

type InvalidNameType struct {
	TestName string `invalid_name:"test_name"`
}

type InvalidNameMapType struct {
	TestName struct {
		MapFoo map[string]interface{} `form:"map_foo"`
	}
}

type FooStructForSliceType struct {
	SliceFoo []int `form:"slice_foo"`
}

type FooStructForSliceMapType struct {
	// Unknown type: not support map
	SliceMapFoo []map[string]interface{} `form:"slice_map_foo"`
}

type FooBarStructForIntType struct {
	IntFoo int `form:"int_foo"`
	IntBar int `form:"int_bar" binding:"required"`
}

type FooBarStructForInt8Type struct {
	Int8Foo int8 `form:"int8_foo"`
	Int8Bar int8 `form:"int8_bar" binding:"required"`
}

type FooBarStructForInt16Type struct {
	Int16Foo int16 `form:"int16_foo"`
	Int16Bar int16 `form:"int16_bar" binding:"required"`
}

type FooBarStructForInt32Type struct {
	Int32Foo int32 `form:"int32_foo"`
	Int32Bar int32 `form:"int32_bar" binding:"required"`
}

type FooBarStructForInt64Type struct {
	Int64Foo int64 `form:"int64_foo"`
	Int64Bar int64 `form:"int64_bar" binding:"required"`
}

type FooBarStructForUintType struct {
	UintFoo uint `form:"uint_foo"`
	UintBar uint `form:"uint_bar" binding:"required"`
}

type FooBarStructForUint8Type struct {
	Uint8Foo uint8 `form:"uint8_foo"`
	Uint8Bar uint8 `form:"uint8_bar" binding:"required"`
}

type FooBarStructForUint16Type struct {
	Uint16Foo uint16 `form:"uint16_foo"`
	Uint16Bar uint16 `form:"uint16_bar" binding:"required"`
}

type FooBarStructForUint32Type struct {
	Uint32Foo uint32 `form:"uint32_foo"`
	Uint32Bar uint32 `form:"uint32_bar" binding:"required"`
}

type FooBarStructForUint64Type struct {
	Uint64Foo uint64 `form:"uint64_foo"`
	Uint64Bar uint64 `form:"uint64_bar" binding:"required"`
}

type FooBarStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
	BoolBar bool `form:"bool_bar" binding:"required"`
}

type FooBarStructForFloat32Type struct {
	Float32Foo float32 `form:"float32_foo"`
	Float32Bar float32 `form:"float32_bar" binding:"required"`
}

type FooBarStructForFloat64Type struct {
	Float64Foo float64 `form:"float64_foo"`
	Float64Bar float64 `form:"float64_bar" binding:"required"`
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

func TestBindingJSONUseNumber(t *testing.T) {
	testBodyBindingUseNumber(t,
		JSON, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
}

func TestBindingJSONUseNumber2(t *testing.T) {
	testBodyBindingUseNumber2(t,
		JSON, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
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

func TestBindingFormForTime(t *testing.T) {
	testFormBindingForTime(t, "POST",
		"/", "/",
		"time_foo=2017-11-15&time_bar=", "bar2=foo")
	testFormBindingForTimeNotFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
	testFormBindingForTimeFailFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
	testFormBindingForTimeFailLocation(t, "POST",
		"/", "/",
		"time_foo=2017-11-15", "bar2=foo")
}

func TestBindingFormForTime2(t *testing.T) {
	testFormBindingForTime(t, "GET",
		"/?time_foo=2017-11-15&time_bar=", "/?bar2=foo",
		"", "")
	testFormBindingForTimeNotFormat(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
	testFormBindingForTimeFailFormat(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
	testFormBindingForTimeFailLocation(t, "GET",
		"/?time_foo=2017-11-15", "/?bar2=foo",
		"", "")
}

func TestBindingFormInvalidName(t *testing.T) {
	testFormBindingInvalidName(t, "POST",
		"/", "/",
		"test_name=bar", "bar2=foo")
}

func TestBindingFormInvalidName2(t *testing.T) {
	testFormBindingInvalidName2(t, "POST",
		"/", "/",
		"map_foo=bar", "bar2=foo")
}

func TestBindingFormForType(t *testing.T) {
	testFormBindingForType(t, "POST",
		"/", "/",
		"map_foo=", "bar2=1", "Map")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_foo=1&slice_foo=2", "bar2=1&bar2=2", "Slice")

	testFormBindingForType(t, "GET",
		"/?slice_foo=1&slice_foo=2", "/?bar2=1&bar2=2",
		"", "", "Slice")

	testFormBindingForType(t, "POST",
		"/", "/",
		"slice_map_foo=1&slice_map_foo=2", "bar2=1&bar2=2", "SliceMap")

	testFormBindingForType(t, "GET",
		"/?slice_map_foo=1&slice_map_foo=2", "/?bar2=1&bar2=2",
		"", "", "SliceMap")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int_foo=&int_bar=-12", "bar2=-123", "Int")

	testFormBindingForType(t, "GET",
		"/?int_foo=&int_bar=-12", "/?bar2=-123",
		"", "", "Int")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int8_foo=&int8_bar=-12", "bar2=-123", "Int8")

	testFormBindingForType(t, "GET",
		"/?int8_foo=&int8_bar=-12", "/?bar2=-123",
		"", "", "Int8")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int16_foo=&int16_bar=-12", "bar2=-123", "Int16")

	testFormBindingForType(t, "GET",
		"/?int16_foo=&int16_bar=-12", "/?bar2=-123",
		"", "", "Int16")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int32_foo=&int32_bar=-12", "bar2=-123", "Int32")

	testFormBindingForType(t, "GET",
		"/?int32_foo=&int32_bar=-12", "/?bar2=-123",
		"", "", "Int32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"int64_foo=&int64_bar=-12", "bar2=-123", "Int64")

	testFormBindingForType(t, "GET",
		"/?int64_foo=&int64_bar=-12", "/?bar2=-123",
		"", "", "Int64")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint_foo=&uint_bar=12", "bar2=123", "Uint")

	testFormBindingForType(t, "GET",
		"/?uint_foo=&uint_bar=12", "/?bar2=123",
		"", "", "Uint")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint8_foo=&uint8_bar=12", "bar2=123", "Uint8")

	testFormBindingForType(t, "GET",
		"/?uint8_foo=&uint8_bar=12", "/?bar2=123",
		"", "", "Uint8")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint16_foo=&uint16_bar=12", "bar2=123", "Uint16")

	testFormBindingForType(t, "GET",
		"/?uint16_foo=&uint16_bar=12", "/?bar2=123",
		"", "", "Uint16")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint32_foo=&uint32_bar=12", "bar2=123", "Uint32")

	testFormBindingForType(t, "GET",
		"/?uint32_foo=&uint32_bar=12", "/?bar2=123",
		"", "", "Uint32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"uint64_foo=&uint64_bar=12", "bar2=123", "Uint64")

	testFormBindingForType(t, "GET",
		"/?uint64_foo=&uint64_bar=12", "/?bar2=123",
		"", "", "Uint64")

	testFormBindingForType(t, "POST",
		"/", "/",
		"bool_foo=&bool_bar=true", "bar2=true", "Bool")

	testFormBindingForType(t, "GET",
		"/?bool_foo=&bool_bar=true", "/?bar2=true",
		"", "", "Bool")

	testFormBindingForType(t, "POST",
		"/", "/",
		"float32_foo=&float32_bar=-12.34", "bar2=12.3", "Float32")

	testFormBindingForType(t, "GET",
		"/?float32_foo=&float32_bar=-12.34", "/?bar2=12.3",
		"", "", "Float32")

	testFormBindingForType(t, "POST",
		"/", "/",
		"float64_foo=&float64_bar=-12.34", "bar2=12.3", "Float64")

	testFormBindingForType(t, "GET",
		"/?float64_foo=&float64_bar=-12.34", "/?bar2=12.3",
		"", "", "Float64")
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST",
		"/?foo=bar&bar=foo", "/",
		"foo=unused", "bar2=foo")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"foo=unused", "")
}

func TestBindingQueryFail(t *testing.T) {
	testQueryBindingFail(t, "POST",
		"/?map_foo=", "/",
		"map_foo=unused", "bar2=foo")
}

func TestBindingQueryFail2(t *testing.T) {
	testQueryBindingFail(t, "GET",
		"/?map_foo=", "/?bar2=foo",
		"map_foo=unused", "")
}

func TestBindingXML(t *testing.T) {
	testBodyBinding(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func TestBindingXMLFail(t *testing.T) {
	testBodyBindingFail(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar<foo></map>", "<map><bar>foo</bar></map>")
}

func createFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestFail() *http.Request {
	req, _ := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo=bar"))
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

func createFormMultipartRequestFail() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	mw.SetBoundary(boundary)
	mw.WriteField("map_foo", "bar")
	req, _ := http.NewRequest("POST", "/?map_foo=getfoo", body)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest()
	var obj FooBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, FormPost.Name(), "form-urlencoded")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestBindingFormPostFail(t *testing.T) {
	req := createFormPostRequestFail()
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest()
	var obj FooBarStruct
	FormMultipart.Bind(req, &obj)

	assert.Equal(t, FormMultipart.Name(), "multipart/form-data")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestBindingFormMultipartFail(t *testing.T) {
	req := createFormMultipartRequestFail()
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.Error(t, err)
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

func TestBindingProtoBufFail(t *testing.T) {
	test := &example.Test{
		Label: proto.String("yes"),
	}
	data, _ := proto.Marshal(test)

	testProtoBodyBindingFail(t,
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

func TestFormBindingFail(t *testing.T) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormPostBindingFail(t *testing.T) {
	b := FormPost
	assert.Equal(t, b.Name(), "form-urlencoded")

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormMultipartBindingFail(t *testing.T) {
	b := FormMultipart
	assert.Equal(t, b.Name(), "multipart/form-data")

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTime(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooBarStructForTimeType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)

	assert.NoError(t, err)
	assert.Equal(t, obj.TimeFoo.Unix(), int64(1510675200))
	assert.Equal(t, obj.TimeFoo.Location().String(), "Asia/Chongqing")
	assert.Equal(t, obj.TimeBar.Unix(), int64(-62135596800))
	assert.Equal(t, obj.TimeBar.Location().String(), "UTC")

	obj = FooBarStructForTimeType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeNotFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooStructForTimeTypeNotFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = FooStructForTimeTypeNotFormat{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooStructForTimeTypeFailFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = FooStructForTimeTypeFailFormat{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeFailLocation(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooStructForTimeTypeFailLocation{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = FooStructForTimeTypeFailLocation{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingInvalidName(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := InvalidNameType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.TestName, "")

	obj = InvalidNameType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingInvalidName2(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := InvalidNameMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = InvalidNameMapType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForType(t *testing.T, method, path, badPath, body, badBody string, typ string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	switch typ {
	case "Int":
		obj := FooBarStructForIntType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.IntFoo, int(0))
		assert.Equal(t, obj.IntBar, int(-12))

		obj = FooBarStructForIntType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int8":
		obj := FooBarStructForInt8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Int8Foo, int8(0))
		assert.Equal(t, obj.Int8Bar, int8(-12))

		obj = FooBarStructForInt8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int16":
		obj := FooBarStructForInt16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Int16Foo, int16(0))
		assert.Equal(t, obj.Int16Bar, int16(-12))

		obj = FooBarStructForInt16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int32":
		obj := FooBarStructForInt32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Int32Foo, int32(0))
		assert.Equal(t, obj.Int32Bar, int32(-12))

		obj = FooBarStructForInt32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int64":
		obj := FooBarStructForInt64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Int64Foo, int64(0))
		assert.Equal(t, obj.Int64Bar, int64(-12))

		obj = FooBarStructForInt64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint":
		obj := FooBarStructForUintType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.UintFoo, uint(0x0))
		assert.Equal(t, obj.UintBar, uint(0xc))

		obj = FooBarStructForUintType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint8":
		obj := FooBarStructForUint8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Uint8Foo, uint8(0x0))
		assert.Equal(t, obj.Uint8Bar, uint8(0xc))

		obj = FooBarStructForUint8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint16":
		obj := FooBarStructForUint16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Uint16Foo, uint16(0x0))
		assert.Equal(t, obj.Uint16Bar, uint16(0xc))

		obj = FooBarStructForUint16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint32":
		obj := FooBarStructForUint32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Uint32Foo, uint32(0x0))
		assert.Equal(t, obj.Uint32Bar, uint32(0xc))

		obj = FooBarStructForUint32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint64":
		obj := FooBarStructForUint64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Uint64Foo, uint64(0x0))
		assert.Equal(t, obj.Uint64Bar, uint64(0xc))

		obj = FooBarStructForUint64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Float32":
		obj := FooBarStructForFloat32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Float32Foo, float32(0.0))
		assert.Equal(t, obj.Float32Bar, float32(-12.34))

		obj = FooBarStructForFloat32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Float64":
		obj := FooBarStructForFloat64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.Float64Foo, float64(0.0))
		assert.Equal(t, obj.Float64Bar, float64(-12.34))

		obj = FooBarStructForFloat64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Bool":
		obj := FooBarStructForBoolType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.BoolFoo, false)
		assert.Equal(t, obj.BoolBar, true)

		obj = FooBarStructForBoolType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Slice":
		obj := FooStructForSliceType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, obj.SliceFoo, []int{1, 2})

		obj = FooStructForSliceType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Map":
		obj := FooStructForMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "SliceMap":
		obj := FooStructForSliceMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	}
}

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, b.Name(), "query")

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func testQueryBindingFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, b.Name(), "query")

	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
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

func testBodyBindingUseNumber(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	assert.NoError(t, e)
	assert.Equal(t, v, int64(123))

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingUseNumber2(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.Equal(t, obj.Foo, float64(123))

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
	assert.Equal(t, obj.Foo, "")

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

type hook struct{}

func (h hook) Read([]byte) (int, error) {
	return 0, errors.New("error")
}

func testProtoBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := example.Test{}
	req := requestWithBody("POST", path, body)

	req.Body = ioutil.NopCloser(&hook{})
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	assert.Error(t, err)

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
