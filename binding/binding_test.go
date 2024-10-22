// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type appkey struct {
	Appkey string `json:"appkey" form:"appkey"`
}

type QueryTest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	appkey
}

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required,max=32"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar" xml:"bar" binding:"required"`
}

type FooBarFileStruct struct {
	FooBarStruct
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FooBarFileFailStruct struct {
	FooBarStruct
	File *multipart.FileHeader `invalid_name:"file" binding:"required"`
	// for unexport test
	data *multipart.FileHeader `form:"data" binding:"required"`
}

type FooDefaultBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar,default=hello" xml:"bar" binding:"required"`
}

type FooStructUseNumber struct {
	Foo any `json:"foo" binding:"required"`
}

type FooStructDisallowUnknownFields struct {
	Foo any `json:"foo" binding:"required"`
}

type FooBarStructForTimeType struct {
	TimeFoo    time.Time `form:"time_foo" time_format:"2006-01-02" time_utc:"1" time_location:"Asia/Chongqing"`
	TimeBar    time.Time `form:"time_bar" time_format:"2006-01-02" time_utc:"1"`
	CreateTime time.Time `form:"createTime" time_format:"unixNano"`
	UnixTime   time.Time `form:"unixTime" time_format:"unix"`
}

type FooStructForTimeTypeNotUnixFormat struct {
	CreateTime time.Time `form:"createTime" time_format:"unixNano"`
	UnixTime   time.Time `form:"unixTime" time_format:"unix"`
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
	MapFoo map[string]any `form:"map_foo"`
}

type FooStructForIgnoreFormTag struct {
	Foo *string `form:"-"`
}

type InvalidNameType struct {
	TestName string `invalid_name:"test_name"`
}

type InvalidNameMapType struct {
	TestName struct {
		MapFoo map[string]any `form:"map_foo"`
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
	SliceMapFoo []map[string]any `form:"slice_map_foo"`
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
}

type FooStructForStringPtrType struct {
	PtrFoo *string `form:"ptr_foo"`
	PtrBar *string `form:"ptr_bar" binding:"required"`
}

type FooStructForMapPtrType struct {
	PtrBar *map[string]any `form:"ptr_bar"`
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

	assert.Equal(t, FormMultipart, Default("POST", MIMEMultipartPOSTForm))
	assert.Equal(t, FormMultipart, Default("PUT", MIMEMultipartPOSTForm))

	assert.Equal(t, ProtoBuf, Default("POST", MIMEPROTOBUF))
	assert.Equal(t, ProtoBuf, Default("PUT", MIMEPROTOBUF))

	assert.Equal(t, YAML, Default("POST", MIMEYAML))
	assert.Equal(t, YAML, Default("PUT", MIMEYAML))
	assert.Equal(t, YAML, Default("POST", MIMEYAML2))
	assert.Equal(t, YAML, Default("PUT", MIMEYAML2))

	assert.Equal(t, TOML, Default("POST", MIMETOML))
	assert.Equal(t, TOML, Default("PUT", MIMETOML))
}

func TestBindingJSONNilBody(t *testing.T) {
	var obj FooStruct
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := JSON.Bind(req, &obj)
	require.Error(t, err)
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		JSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingJSONSlice(t *testing.T) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	testBodyBindingSlice(t, JSON, "json", "/", "/", `[]`, ``)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": ""}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": 123}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"bar": 123}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": "123456789012345678901234567890123"}]`)
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

func TestBindingJSONDisallowUnknownFields(t *testing.T) {
	testBodyBindingDisallowUnknownFields(t, JSON,
		"/", "/",
		`{"foo": "bar"}`, `{"foo": "bar", "what": "this"}`)
}

func TestBindingJSONStringMap(t *testing.T) {
	testBodyBindingStringMap(t, JSON,
		"/", "/",
		`{"foo": "bar", "hello": "world"}`, `{"num": 2}`)
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

func TestBindingFormEmbeddedStruct(t *testing.T) {
	testFormBindingEmbeddedStruct(t, "POST",
		"/", "/",
		"page=1&size=2&appkey=test-appkey", "bar2=foo")
}

func TestBindingFormEmbeddedStruct2(t *testing.T) {
	testFormBindingEmbeddedStruct(t, "GET",
		"/?page=1&size=2&appkey=test-appkey", "/?bar2=foo",
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
		"time_foo=2017-11-15&time_bar=&createTime=1562400033000000123&unixTime=1562400033", "bar2=foo")
	testFormBindingForTimeNotUnixFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15&createTime=bad&unixTime=bad", "bar2=foo")
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
		"/?time_foo=2017-11-15&time_bar=&createTime=1562400033000000123&unixTime=1562400033", "/?bar2=foo",
		"", "")
	testFormBindingForTimeNotUnixFormat(t, "POST",
		"/", "/",
		"time_foo=2017-11-15&createTime=bad&unixTime=bad", "bar2=foo")
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

func TestFormBindingIgnoreField(t *testing.T) {
	testFormBindingIgnoreField(t, "POST",
		"/", "/",
		"-=bar", "")
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
		"map_foo={\"bar\":123}", "map_foo=1", "Map")

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

func TestBindingFormStringMap(t *testing.T) {
	testBodyBindingStringMap(t, Form,
		"/", "",
		`foo=bar&hello=world`, "")
	// Should pick the last value
	testBodyBindingStringMap(t, Form,
		"/", "",
		`foo=something&foo=bar&hello=world`, "")
}

func TestBindingFormStringSliceMap(t *testing.T) {
	obj := make(map[string][]string)
	req := requestWithBody("POST", "/", "foo=something&foo=bar&hello=world")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err := Form.Bind(req, &obj)
	require.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	target := map[string][]string{
		"foo":   {"something", "bar"},
		"hello": {"world"},
	}
	assert.True(t, reflect.DeepEqual(obj, target))

	objInvalid := make(map[string][]int)
	req = requestWithBody("POST", "/", "foo=something&foo=bar&hello=world")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	err = Form.Bind(req, &objInvalid)
	require.Error(t, err)
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

func TestBindingQueryStringMap(t *testing.T) {
	b := Query

	obj := make(map[string]string)
	req := requestWithBody("GET", "/?foo=bar&hello=world", "")
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "bar", obj["foo"])
	assert.Equal(t, "world", obj["hello"])

	obj = make(map[string]string)
	req = requestWithBody("GET", "/?foo=bar&foo=2&hello=world", "") // should pick last
	err = b.Bind(req, &obj)
	require.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "2", obj["foo"])
	assert.Equal(t, "world", obj["hello"])
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

func TestBindingTOML(t *testing.T) {
	testBodyBinding(t,
		TOML, "toml",
		"/", "/",
		`foo="bar"`, `bar="foo"`)
}

func TestBindingTOMLFail(t *testing.T) {
	testBodyBindingFail(t,
		TOML, "toml",
		"/", "/",
		`foo=\n"bar"`, `bar="foo"`)
}

func TestBindingYAML(t *testing.T) {
	testBodyBinding(t,
		YAML, "yaml",
		"/", "/",
		`foo: bar`, `bar: foo`)
}

func TestBindingYAMLStringMap(t *testing.T) {
	// YAML is a superset of JSON, so the test below is JSON (to avoid newlines)
	testBodyBindingStringMap(t, YAML,
		"/", "/",
		`{"foo": "bar", "hello": "world"}`, `{"nested": {"foo": "bar"}}`)
}

func TestBindingYAMLFail(t *testing.T) {
	testBodyBindingFail(t,
		YAML, "yaml",
		"/", "/",
		`foo:\nbar`, `bar: foo`)
}

func createFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createDefaultFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMap(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo={\"bar\":123}"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMapFail(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo=hello"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormFilesMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	require.NoError(t, mw.SetBoundary(boundary))
	require.NoError(t, mw.WriteField("foo", "bar"))
	require.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	require.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file", "form.go")
	require.NoError(t, err1)
	_, err = io.Copy(fw, f)
	require.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	require.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormFilesMultipartRequestFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	require.NoError(t, mw.SetBoundary(boundary))
	require.NoError(t, mw.WriteField("foo", "bar"))
	require.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	require.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file_foo", "form_foo.go")
	require.NoError(t, err1)
	_, err = io.Copy(fw, f)
	require.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	require.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	require.NoError(t, mw.SetBoundary(boundary))
	require.NoError(t, mw.WriteField("foo", "bar"))
	require.NoError(t, mw.WriteField("bar", "foo"))
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMap(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	require.NoError(t, mw.SetBoundary(boundary))
	require.NoError(t, mw.WriteField("map_foo", "{\"bar\":123, \"name\":\"thinkerou\", \"pai\": 3.14}"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMapFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	require.NoError(t, mw.SetBoundary(boundary))
	require.NoError(t, mw.WriteField("map_foo", "3.14"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest(t)
	var obj FooBarStruct
	require.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingDefaultValueFormPost(t *testing.T) {
	req := createDefaultFormPostRequest(t)
	var obj FooDefaultBarStruct
	require.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
}

func TestBindingFormPostForMap(t *testing.T) {
	req := createFormPostRequestForMap(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	require.NoError(t, err)
	assert.InDelta(t, float64(123), obj.MapFoo["bar"].(float64), 0.01)
}

func TestBindingFormPostForMapFail(t *testing.T) {
	req := createFormPostRequestForMapFail(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	require.Error(t, err)
}

func TestBindingFormFilesMultipart(t *testing.T) {
	req := createFormFilesMultipartRequest(t)
	var obj FooBarFileStruct
	err := FormMultipart.Bind(req, &obj)
	require.NoError(t, err)

	// file from os
	f, _ := os.Open("form.go")
	defer f.Close()
	fileActual, _ := io.ReadAll(f)

	// file from multipart
	mf, _ := obj.File.Open()
	defer mf.Close()
	fileExpect, _ := io.ReadAll(mf)

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
	assert.Equal(t, fileExpect, fileActual)
}

func TestBindingFormFilesMultipartFail(t *testing.T) {
	req := createFormFilesMultipartRequestFail(t)
	var obj FooBarFileFailStruct
	err := FormMultipart.Bind(req, &obj)
	require.Error(t, err)
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest(t)
	var obj FooBarStruct
	require.NoError(t, FormMultipart.Bind(req, &obj))

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingFormMultipartForMap(t *testing.T) {
	req := createFormMultipartRequestForMap(t)
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	require.NoError(t, err)
	assert.InDelta(t, float64(123), obj.MapFoo["bar"].(float64), 0.01)
	assert.Equal(t, "thinkerou", obj.MapFoo["name"].(string))
	assert.InDelta(t, float64(3.14), obj.MapFoo["pai"].(float64), 0.01)
}

func TestBindingFormMultipartForMapFail(t *testing.T) {
	req := createFormMultipartRequestForMapFail(t)
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	require.Error(t, err)
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

func TestValidationFails(t *testing.T) {
	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	require.Error(t, err)
}

func TestValidationDisabled(t *testing.T) {
	backup := Validator
	Validator = nil
	defer func() { Validator = backup }()

	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	require.NoError(t, err)
}

func TestRequiredSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := JSON.Bind(req, &obj)
	require.NoError(t, err)
}

func TestRequiredFails(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"foo" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := JSON.Bind(req, &obj)
	require.Error(t, err)
}

func TestHeaderBinding(t *testing.T) {
	h := Header
	assert.Equal(t, "header", h.Name())

	type tHeader struct {
		Limit int `header:"limit"`
	}

	var theader tHeader
	req := requestWithBody("GET", "/", "")
	req.Header.Add("limit", "1000")
	require.NoError(t, h.Bind(req, &theader))
	assert.Equal(t, 1000, theader.Limit)

	req = requestWithBody("GET", "/", "")
	req.Header.Add("fail", `{fail:fail}`)

	type failStruct struct {
		Fail map[string]any `header:"fail"`
	}

	err := h.Bind(req, &failStruct{})
	require.Error(t, err)
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
	require.NoError(t, b.BindUri(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]any `uri:"name"`
	}
	var not NotSupportStruct
	require.Error(t, b.BindUri(m, &not))
	assert.Equal(t, map[string]any(nil), not.Name)
}

func TestUriInnerBinding(t *testing.T) {
	type Tag struct {
		Name string `uri:"name"`
		S    struct {
			Age int `uri:"age"`
		}
	}

	expectedName := "mike"
	expectedAge := 25

	m := map[string][]string{
		"name": {expectedName},
		"age":  {strconv.Itoa(expectedAge)},
	}

	var tag Tag
	require.NoError(t, Uri.BindUri(m, &tag))
	assert.Equal(t, expectedName, tag.Name)
	assert.Equal(t, expectedAge, tag.S.Age)
}

func testFormBindingEmbeddedStruct(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := QueryTest{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, 1, obj.Page)
	assert.Equal(t, 2, obj.Size)
	assert.Equal(t, "test-appkey", obj.Appkey)
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
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)

	obj = FooBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
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
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)

	obj = FooDefaultBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func TestFormBindingFail(t *testing.T) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	require.Error(t, err)
}

func TestFormBindingMultipartFail(t *testing.T) {
	obj := FooBarStruct{}
	req, err := http.NewRequest("POST", "/", strings.NewReader("foo=bar"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+";boundary=testboundary")
	_, err = req.MultipartReader()
	require.NoError(t, err)
	err = Form.Bind(req, &obj)
	require.Error(t, err)
}

func TestFormPostBindingFail(t *testing.T) {
	b := FormPost
	assert.Equal(t, "form-urlencoded", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	require.Error(t, err)
}

func TestFormMultipartBindingFail(t *testing.T) {
	b := FormMultipart
	assert.Equal(t, "multipart/form-data", b.Name())

	obj := FooBarStruct{}
	req, _ := http.NewRequest("POST", "/", nil)
	err := b.Bind(req, &obj)
	require.Error(t, err)
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

	require.NoError(t, err)
	assert.Equal(t, int64(1510675200), obj.TimeFoo.Unix())
	assert.Equal(t, "Asia/Chongqing", obj.TimeFoo.Location().String())
	assert.Equal(t, int64(-62135596800), obj.TimeBar.Unix())
	assert.Equal(t, "UTC", obj.TimeBar.Location().String())
	assert.Equal(t, int64(1562400033000000123), obj.CreateTime.UnixNano())
	assert.Equal(t, int64(1562400033), obj.UnixTime.Unix())

	obj = FooBarStructForTimeType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testFormBindingForTimeNotUnixFormat(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForTimeTypeNotUnixFormat{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	require.Error(t, err)

	obj = FooStructForTimeTypeNotUnixFormat{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
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
	require.Error(t, err)

	obj = FooStructForTimeTypeNotFormat{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
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
	require.Error(t, err)

	obj = FooStructForTimeTypeFailFormat{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
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
	require.Error(t, err)

	obj = FooStructForTimeTypeFailLocation{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testFormBindingIgnoreField(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForIgnoreFormTag{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	require.NoError(t, err)

	assert.Nil(t, obj.Foo)
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
	require.NoError(t, err)
	assert.Equal(t, "", obj.TestName)

	obj = InvalidNameType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
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
	require.Error(t, err)

	obj = InvalidNameMapType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testFormBindingForType(t *testing.T, method, path, badPath, body, badBody string, typ string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	switch typ {
	case "Slice":
		obj := FooStructForSliceType{}
		err := b.Bind(req, &obj)
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2}, obj.SliceFoo)

		obj = FooStructForSliceType{}
		req = requestWithBody(method, badPath, badBody)
		err = JSON.Bind(req, &obj)
		require.Error(t, err)
	case "Struct":
		obj := FooStructForStructType{}
		err := b.Bind(req, &obj)
		require.NoError(t, err)
		assert.Equal(t,
			struct {
				Idx int "form:\"idx\""
			}{Idx: 123},
			obj.StructFoo)
	case "StructPointer":
		obj := FooStructForStructPointerType{}
		err := b.Bind(req, &obj)
		require.NoError(t, err)
		assert.Equal(t,
			struct {
				Name string "form:\"name\""
			}{Name: "thinkerou"},
			*obj.StructPointerFoo)
	case "Map":
		obj := FooStructForMapType{}
		err := b.Bind(req, &obj)
		require.NoError(t, err)
		assert.InDelta(t, float64(123), obj.MapFoo["bar"].(float64), 0.01)
	case "SliceMap":
		obj := FooStructForSliceMapType{}
		err := b.Bind(req, &obj)
		require.Error(t, err)
	case "Ptr":
		obj := FooStructForStringPtrType{}
		err := b.Bind(req, &obj)
		require.NoError(t, err)
		assert.Nil(t, obj.PtrFoo)
		assert.Equal(t, "test", *obj.PtrBar)

		obj = FooStructForStringPtrType{}
		obj.PtrBar = new(string)
		err = b.Bind(req, &obj)
		require.NoError(t, err)
		assert.Equal(t, "test", *obj.PtrBar)

		objErr := FooStructForMapPtrType{}
		err = b.Bind(req, &objErr)
		require.Error(t, err)

		obj = FooStructForStringPtrType{}
		req = requestWithBody(method, badPath, badBody)
		err = b.Bind(req, &obj)
		require.Error(t, err)
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
	require.NoError(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testBodyBindingSlice(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	var obj1 []FooStruct
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj1)
	require.NoError(t, err)

	var obj2 []FooStruct
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj2)
	require.Error(t, err)
}

func testBodyBindingStringMap(t *testing.T, b Binding, path, badPath, body, badBody string) {
	obj := make(map[string]string)
	req := requestWithBody("POST", path, body)
	if b.Name() == "form" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "bar", obj["foo"])
	assert.Equal(t, "world", obj["hello"])

	if badPath != "" && badBody != "" {
		obj = make(map[string]string)
		req = requestWithBody("POST", badPath, badBody)
		err = b.Bind(req, &obj)
		require.Error(t, err)
	}

	objInt := make(map[string]int)
	req = requestWithBody("POST", path, body)
	err = b.Bind(req, &objInt)
	require.Error(t, err)
}

func testBodyBindingUseNumber(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	require.NoError(t, e)
	assert.Equal(t, int64(123), v)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testBodyBindingUseNumber2(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.InDelta(t, float64(123), obj.Foo, 0.01)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testBodyBindingDisallowUnknownFields(t *testing.T, b Binding, path, badPath, body, badBody string) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	obj := FooStructDisallowUnknownFields{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStructDisallowUnknownFields{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "what")
}

func testBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	require.Error(t, err)
	assert.Equal(t, "", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	require.Error(t, err)
}

func testProtoBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := protoexample.Test{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	assert.Equal(t, "yes", *obj.Label)

	obj = protoexample.Test{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err = ProtoBuf.Bind(req, &obj)
	require.Error(t, err)
}

type hook struct{}

func (h hook) Read([]byte) (int, error) {
	return 0, errors.New("error")
}

type failRead struct{}

func (f *failRead) Read(b []byte) (n int, err error) {
	return 0, errors.New("my fail")
}

func (f *failRead) Close() error {
	return nil
}

func TestPlainBinding(t *testing.T) {
	p := Plain
	assert.Equal(t, "plain", p.Name())

	var s string
	req := requestWithBody("POST", "/", "test string")
	require.NoError(t, p.Bind(req, &s))
	assert.Equal(t, "test string", s)

	var bs []byte
	req = requestWithBody("POST", "/", "test []byte")
	require.NoError(t, p.Bind(req, &bs))
	assert.Equal(t, bs, []byte("test []byte"))

	var i int
	req = requestWithBody("POST", "/", "test fail")
	require.Error(t, p.Bind(req, &i))

	req = requestWithBody("POST", "/", "")
	req.Body = &failRead{}
	require.Error(t, p.Bind(req, &s))

	req = requestWithBody("POST", "/", "")
	require.NoError(t, p.Bind(req, nil))

	var ptr *string
	req = requestWithBody("POST", "/", "")
	require.NoError(t, p.Bind(req, ptr))
}

func testProtoBodyBindingFail(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := protoexample.Test{}
	req := requestWithBody("POST", path, body)

	req.Body = io.NopCloser(&hook{})
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err := b.Bind(req, &obj)
	require.Error(t, err)

	invalidobj := FooStruct{}
	req.Body = io.NopCloser(strings.NewReader(`{"msg":"hello"}`))
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err = b.Bind(req, &invalidobj)
	require.Error(t, err)
	assert.Equal(t, "obj is not ProtoMessage", err.Error())

	obj = protoexample.Test{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", MIMEPROTOBUF)
	err = ProtoBuf.Bind(req, &obj)
	require.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
