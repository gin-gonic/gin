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

	"github.com/gin-gonic/gin/testdata/protoexample"
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

type FooDefaultBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar,default=hello" xml:"bar" binding:"required"`
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

type FooStructForStructType struct {
	StructFoo struct {
		Idx int `form:"idx"`
	}
}

type FooStructForStructPointerType struct {
	StructPointerFoo *struct {
		Name string `form:"name"`
	}
}

type FooStructForSliceMapType struct {
	// Unknown type: not support map
	SliceMapFoo []map[string]interface{} `form:"slice_map_foo"`
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
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

type FooStructForStringPtrType struct {
	PtrFoo *string `form:"ptr_foo"`
	PtrBar *string `form:"ptr_bar" binding:"required"`
}

type FooStructForMapPtrType struct {
	PtrBar *map[string]interface{} `form:"ptr_bar"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, Form, Default("GET", ""))
	assert.Equal(t, Form, Default("GET", MIMEJSON))

	assert.Equal(t, JSON, Default("POST", MIMEJSON))
	assert.Equal(t, JSON, Default("PUT", MIMEJSON))

	assert.Equal(t, XML, Default("POST", MIMEXML))
	assert.Equal(t, XML, Default("PUT", MIMEXML2))

	assert.Equal(t, Form, Default("POST", MIMEPOSTForm))
	assert.Equal(t, Form, Default("PUT", MIMEPOSTForm))

	assert.Equal(t, Form, Default("POST", MIMEMultipartPOSTForm))
	assert.Equal(t, Form, Default("PUT", MIMEMultipartPOSTForm))

	assert.Equal(t, ProtoBuf, Default("POST", MIMEPROTOBUF))
	assert.Equal(t, ProtoBuf, Default("PUT", MIMEPROTOBUF))

	assert.Equal(t, MsgPack, Default("POST", MIMEMSGPACK))
	assert.Equal(t, MsgPack, Default("PUT", MIMEMSGPACK2))

	assert.Equal(t, YAML, Default("POST", MIMEYAML))
	assert.Equal(t, YAML, Default("PUT", MIMEYAML))
}

func TestBindingJSONNilBody(t *testing.T) {
	var obj FooStruct
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
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

func TestBindingFormDefaultValue(t *testing.T) {
	testFormBindingDefaultValue(t, "POST",
		"/", "/",
		"foo=bar", "bar2=foo")
}

func TestBindingFormDefaultValue2(t *testing.T) {
	testFormBindingDefaultValue(t, "GET",
		"/?foo=bar", "/?bar2=foo",
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

	testFormBindingForType(t, "POST",
		"/", "/",
		"ptr_bar=test", "bar2=test", "Ptr")

	testFormBindingForType(t, "GET",
		"/?ptr_bar=test", "/?bar2=test",
		"", "", "Ptr")

	testFormBindingForType(t, "POST",
		"/", "/",
		"idx=123", "id1=1", "Struct")

	testFormBindingForType(t, "GET",
		"/?idx=123", "/?id1=1",
		"", "", "Struct")

	testFormBindingForType(t, "POST",
		"/", "/",
		"name=thinkerou", "name1=ou", "StructPointer")

	testFormBindingForType(t, "GET",
		"/?name=thinkerou", "/?name1=ou",
		"", "", "StructPointer")
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

func TestBindingQueryBoolFail(t *testing.T) {
	testQueryBindingBoolFail(t, "GET",
		"/?bool_foo=fasl", "/?bar2=foo",
		"bool_foo=unused", "")
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

func TestBindingYAML(t *testing.T) {
	testBodyBinding(t,
		YAML, "yaml",
		"/", "/",
		`foo: bar`, `bar: foo`)
}

func TestBindingYAMLFail(t *testing.T) {
	testBodyBindingFail(t,
		YAML, "yaml",
		"/", "/",
		`foo:\nbar`, `bar: foo`)
}

func createFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createDefaultFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar"))
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

	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingDefaultValueFormPost(t *testing.T) {
	req := createDefaultFormPostRequest()
	var obj FooDefaultBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
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

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingFormMultipartFail(t *testing.T) {
	req := createFormMultipartRequestFail()
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingProtoBuf(t *testing.T) {
	test := &protoexample.Test{
		Label: proto.String("yes"),
	}
	data, _ := proto.Marshal(test)

	testProtoBodyBinding(t,
		ProtoBuf, "protobuf",
		"/", "/",
		string(data), string(data[1:]))
}

func TestBindingProtoBufFail(t *testing.T) {
	test := &protoexample.Test{
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

func TestUriBinding(t *testing.T) {
	b := Uri
	assert.Equal(t, "uri", b.Name())

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NoError(t, b.BindUri(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]interface{} `uri:"name"`
	}
	var not NotSupportStruct
	assert.Error(t, b.BindUri(m, &not))
	assert.Equal(t, map[string]interface{}(nil), not.Name)
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)

	obj = FooBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingDefaultValue(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooDefaultBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)

	obj = FooDefaultBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormBindingFail(t *testing.T) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormPostBindingFail(t *testing.T) {
	b := FormPost
	assert.Equal(t, "form-urlencoded", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func TestFormMultipartBindingFail(t *testing.T) {
	b := FormMultipart
	assert.Equal(t, "multipart/form-data", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTime(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStructForTimeType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)

	assert.NoError(t, err)
	assert.Equal(t, int64(1510675200), obj.TimeFoo.Unix())
	assert.Equal(t, "Asia/Chongqing", obj.TimeFoo.Location().String())
	assert.Equal(t, int64(-62135596800), obj.TimeBar.Unix())
	assert.Equal(t, "UTC", obj.TimeBar.Location().String())

	obj = FooBarStructForTimeType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingForTimeNotFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

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
	assert.Equal(t, "form", b.Name())

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
	assert.Equal(t, "form", b.Name())

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
	assert.Equal(t, "form", b.Name())

	obj := InvalidNameType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "", obj.TestName)

	obj = InvalidNameType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testFormBindingInvalidName2(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

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
	assert.Equal(t, "form", b.Name())

	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	switch typ {
	case "Int":
		obj := FooBarStructForIntType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int(0), obj.IntFoo)
		assert.Equal(t, int(-12), obj.IntBar)

		obj = FooBarStructForIntType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int8":
		obj := FooBarStructForInt8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int8(0), obj.Int8Foo)
		assert.Equal(t, int8(-12), obj.Int8Bar)

		obj = FooBarStructForInt8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int16":
		obj := FooBarStructForInt16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int16(0), obj.Int16Foo)
		assert.Equal(t, int16(-12), obj.Int16Bar)

		obj = FooBarStructForInt16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int32":
		obj := FooBarStructForInt32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int32(0), obj.Int32Foo)
		assert.Equal(t, int32(-12), obj.Int32Bar)

		obj = FooBarStructForInt32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Int64":
		obj := FooBarStructForInt64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), obj.Int64Foo)
		assert.Equal(t, int64(-12), obj.Int64Bar)

		obj = FooBarStructForInt64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint":
		obj := FooBarStructForUintType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint(0x0), obj.UintFoo)
		assert.Equal(t, uint(0xc), obj.UintBar)

		obj = FooBarStructForUintType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint8":
		obj := FooBarStructForUint8Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x0), obj.Uint8Foo)
		assert.Equal(t, uint8(0xc), obj.Uint8Bar)

		obj = FooBarStructForUint8Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint16":
		obj := FooBarStructForUint16Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0x0), obj.Uint16Foo)
		assert.Equal(t, uint16(0xc), obj.Uint16Bar)

		obj = FooBarStructForUint16Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint32":
		obj := FooBarStructForUint32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x0), obj.Uint32Foo)
		assert.Equal(t, uint32(0xc), obj.Uint32Bar)

		obj = FooBarStructForUint32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Uint64":
		obj := FooBarStructForUint64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0x0), obj.Uint64Foo)
		assert.Equal(t, uint64(0xc), obj.Uint64Bar)

		obj = FooBarStructForUint64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Float32":
		obj := FooBarStructForFloat32Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, float32(0.0), obj.Float32Foo)
		assert.Equal(t, float32(-12.34), obj.Float32Bar)

		obj = FooBarStructForFloat32Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Float64":
		obj := FooBarStructForFloat64Type{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, float64(0.0), obj.Float64Foo)
		assert.Equal(t, float64(-12.34), obj.Float64Bar)

		obj = FooBarStructForFloat64Type{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Bool":
		obj := FooBarStructForBoolType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.False(t, obj.BoolFoo)
		assert.True(t, obj.BoolBar)

		obj = FooBarStructForBoolType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Slice":
		obj := FooStructForSliceType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, obj.SliceFoo)

		obj = FooStructForSliceType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		assert.Error(t, err)
	case "Struct":
		obj := FooStructForStructType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t,
			struct {
				Idx int "form:\"idx\""
			}(struct {
				Idx int "form:\"idx\""
			}{Idx: 123}),
			obj.StructFoo)
	case "StructPointer":
		obj := FooStructForStructPointerType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t,
			struct {
				Name string "form:\"name\""
			}(struct {
				Name string "form:\"name\""
			}{Name: "thinkerou"}),
			*obj.StructPointerFoo)
	case "Map":
		obj := FooStructForMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "SliceMap":
		obj := FooStructForSliceMapType{}
		err := b.Bind(req, &obj)
		assert.Error(t, err)
	case "Ptr":
		obj := FooStructForStringPtrType{}
		err := b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Nil(t, obj.PtrFoo)
		assert.Equal(t, "test", *obj.PtrBar)

		obj = FooStructForStringPtrType{}
		obj.PtrBar = new(string)
		err = b.Bind(req, &obj)
		assert.NoError(t, err)
		assert.Equal(t, "test", *obj.PtrBar)

		objErr := FooStructForMapPtrType{}
		err = b.Bind(req, &objErr)
		assert.Error(t, err)

		obj = FooStructForStringPtrType{}
		req = requestWithBody(method, badPath, badBody)
		err = b.Bind(req, &obj)
		assert.Error(t, err)
	}
}

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testQueryBindingFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testQueryBindingBoolFail(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, "query", b.Name())

	obj := FooStructForBoolType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingUseNumber(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	assert.NoError(t, e)
	assert.Equal(t, int64(123), v)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingUseNumber2(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.Equal(t, float64(123), obj.Foo)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.Error(t, err)
	assert.Equal(t, "", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testProtoBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := protoexample.Test{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "yes", *obj.Label)

	obj = protoexample.Test{}
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
	assert.Equal(t, name, b.Name())

	obj := protoexample.Test{}
	req := requestWithBody("POST", path, body)

	req.Body = ioutil.NopCloser(&hook{})
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	assert.Error(t, err)

	obj = protoexample.Test{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err = ProtoBuf.Bind(req, &obj)
	assert.Error(t, err)
}

func testMsgPackBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", MIMEMSGPACK)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

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

func TestCanSet(t *testing.T) {
	type CanSetStruct struct {
		lowerStart string `form:"lower"`
	}

	var c CanSetStruct
	assert.Nil(t, mapForm(&c, nil))
}
