// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
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
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
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
	Foo interface{} `json:"foo" binding:"required"`
}

type FooStructDisallowUnknownFields struct {
	Foo interface{} `json:"foo" binding:"required"`
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
	MapFoo map[string]interface{} `form:"map_foo"`
}

type FooStructForIgnoreFormTag struct {
	Foo *string `form:"-"`
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

	assert.Equal(t, FormMultipart, Default("POST", MIMEMultipartPOSTForm))
	assert.Equal(t, FormMultipart, Default("PUT", MIMEMultipartPOSTForm))

	assert.Equal(t, ProtoBuf, Default("POST", MIMEPROTOBUF))
	assert.Equal(t, ProtoBuf, Default("PUT", MIMEPROTOBUF))

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

func TestBindingJSONDisallowUnknownFields(t *testing.T) {
	testBodyBindingDisallowUnknownFields(t, JSON,
		"/", "/",
		`{"foo": "bar"}`, `{"foo": "bar", "what": "this"}`)
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

func createFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createDefaultFormPostRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMap(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo={\"bar\":123}"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormPostRequestForMapFail(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", bytes.NewBufferString("map_foo=hello"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormFilesMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	assert.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file", "form.go")
	assert.NoError(t, err1)
	_, err = io.Copy(fw, f)
	assert.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormFilesMultipartRequestFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))

	f, err := os.Open("form.go")
	assert.NoError(t, err)
	defer f.Close()
	fw, err1 := mw.CreateFormFile("file_foo", "form_foo.go")
	assert.NoError(t, err1)
	_, err = io.Copy(fw, f)
	assert.NoError(t, err)

	req, err2 := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err2)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)

	return req
}

func createFormMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("foo", "bar"))
	assert.NoError(t, mw.WriteField("bar", "foo"))
	req, err := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMap(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("map_foo", "{\"bar\":123, \"name\":\"thinkerou\", \"pai\": 3.14}"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func createFormMultipartRequestForMapFail(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	assert.NoError(t, mw.SetBoundary(boundary))
	assert.NoError(t, mw.WriteField("map_foo", "3.14"))
	req, err := http.NewRequest("POST", "/?map_foo=getfoo", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest(t)
	var obj FooBarStruct
	assert.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "form-urlencoded", FormPost.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingDefaultValueFormPost(t *testing.T) {
	req := createDefaultFormPostRequest(t)
	var obj FooDefaultBarStruct
	assert.NoError(t, FormPost.Bind(req, &obj))

	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "hello", obj.Bar)
}

func TestBindingFormPostForMap(t *testing.T) {
	req := createFormPostRequestForMap(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), obj.MapFoo["bar"].(float64))
}

func TestBindingFormPostForMapFail(t *testing.T) {
	req := createFormPostRequestForMapFail(t)
	var obj FooStructForMapType
	err := FormPost.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormFilesMultipart(t *testing.T) {
	req := createFormFilesMultipartRequest(t)
	var obj FooBarFileStruct
	err := FormMultipart.Bind(req, &obj)
	assert.NoError(t, err)

	// file from os
	f, _ := os.Open("form.go")
	defer f.Close()
	fileActual, _ := ioutil.ReadAll(f)

	// file from multipart
	mf, _ := obj.File.Open()
	defer mf.Close()
	fileExpect, _ := ioutil.ReadAll(mf)

	assert.Equal(t, FormMultipart.Name(), "multipart/form-data")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, fileExpect, fileActual)
}

func TestBindingFormFilesMultipartFail(t *testing.T) {
	req := createFormFilesMultipartRequestFail(t)
	var obj FooBarFileFailStruct
	err := FormMultipart.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest(t)
	var obj FooBarStruct
	assert.NoError(t, FormMultipart.Bind(req, &obj))

	assert.Equal(t, "multipart/form-data", FormMultipart.Name())
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func TestBindingFormMultipartForMap(t *testing.T) {
	req := createFormMultipartRequestForMap(t)
	var obj FooStructForMapType
	err := FormMultipart.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), obj.MapFoo["bar"].(float64))
	assert.Equal(t, "thinkerou", obj.MapFoo["name"].(string))
	assert.Equal(t, float64(3.14), obj.MapFoo["pai"].(float64))
}

func TestBindingFormMultipartForMapFail(t *testing.T) {
	req := createFormMultipartRequestForMapFail(t)
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

func TestRequiredSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := JSON.Bind(req, &obj)
	assert.NoError(t, err)
}

func TestRequiredFails(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"foo" binding:"required"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"boen": 0}`)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
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
	assert.NoError(t, h.Bind(req, &theader))
	assert.Equal(t, 1000, theader.Limit)

	req = requestWithBody("GET", "/", "")
	req.Header.Add("fail", `{fail:fail}`)

	type failStruct struct {
		Fail map[string]interface{} `header:"fail"`
	}

	err := h.Bind(req, &failStruct{})
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
	assert.NoError(t, Uri.BindUri(m, &tag))
	assert.Equal(t, tag.Name, expectedName)
	assert.Equal(t, tag.S.Age, expectedAge)
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
	assert.NoError(t, err)
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

func TestFormBindingMultipartFail(t *testing.T) {
	obj := FooBarStruct{}
	req, err := http.NewRequest("POST", "/", strings.NewReader("foo=bar"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+";boundary=testboundary")
	_, err = req.MultipartReader()
	assert.NoError(t, err)
	err = Form.Bind(req, &obj)
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
	assert.Equal(t, int64(1562400033000000123), obj.CreateTime.UnixNano())
	assert.Equal(t, int64(1562400033), obj.UnixTime.Unix())

	obj = FooBarStructForTimeType{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
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
	assert.Error(t, err)

	obj = FooStructForTimeTypeNotUnixFormat{}
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

func testFormBindingIgnoreField(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, "form", b.Name())

	obj := FooStructForIgnoreFormTag{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)

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
		assert.NoError(t, err)
		assert.Equal(t, float64(123), obj.MapFoo["bar"].(float64))
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

func testBodyBindingDisallowUnknownFields(t *testing.T, b Binding, path, badPath, body, badBody string) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	obj := FooStructDisallowUnknownFields{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStructDisallowUnknownFields{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "what")
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

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
