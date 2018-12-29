// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"encoding/xml"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"

	testdata "github.com/gin-gonic/gin/testdata/protoexample"
)

// TODO unit tests
// test errors

func TestRenderEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	data := "empty data"
	r := Default(EmptyRenderType)
	r.Setup(data)
	err := r.Render(w)
	Recycle(EmptyRenderType, r)
	assert.EqualError(t, err, "empty render,you need register one first")
}

// test not registered render type
func TestRenderUnknown(t *testing.T) {
	r := Default(unknownRenderType)
	_, ok := r.(*EmptyRender)
	assert.True(t, ok)
}

func TestRenderMsgPack(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	r := Default(MsgPackRenderType)
	r.Setup(data)
	r.WriteContentType(w)
	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data)
	err := r.Render(w)
	Recycle(MsgPackRenderType, r)
	assert.NoError(t, err)

	h := new(codec.MsgpackHandle)
	assert.NotNil(t, h)
	buf := bytes.NewBuffer([]byte{})
	assert.NotNil(t, buf)
	err = codec.NewEncoder(buf, h).Encode(data)

	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), string(buf.Bytes()))
	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo":  "bar",
		"html": "<b>",
	}

	r := Default(JSONRenderType)
	r.Setup(data)
	r.WriteContentType(w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data)
	err := r.Render(w)
	Recycle(JSONRenderType, r)

	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJSONPanics(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	assert.Panics(t, func() {
		r := Default(JSONRenderType)
		r.Setup(data)
		r.Render(w)
		Recycle(JSONRenderType, r)
	})
}

func TestRenderIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
		"bar": "foo",
	}

	r := Default(IntendedJSONRenderType)
	r.Setup(data)
	err := r.Render(w)
	Recycle(IntendedJSONRenderType, r)

	assert.NoError(t, err)
	assert.Equal(t, "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\"\n}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderIndentedJSONPanics(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	r := Default(IntendedJSONRenderType)
	r.Setup(data)
	err := r.Render(w)
	Recycle(IntendedJSONRenderType, r)

	assert.Error(t, err)
}

func TestRenderSecureJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	r := Default(SecureJSONRenderType)
	r.Setup(data, "while(1);")
	r.WriteContentType(w1)
	Recycle(SecureJSONRenderType, r)

	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	r1 := Default(SecureJSONRenderType)
	r1.Setup(data, "while(1);")
	err1 := r1.Render(w1)
	Recycle(SecureJSONRenderType, r1)

	assert.NoError(t, err1)
	assert.Equal(t, "{\"foo\":\"bar\"}", w1.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]interface{}{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	r2 := Default(SecureJSONRenderType)
	r2.Setup(datas, "while(1);")
	err2 := r2.Render(w2)
	Recycle(SecureJSONRenderType, r2)

	assert.NoError(t, err2)
	assert.Equal(t, "while(1);[{\"foo\":\"bar\"},{\"bar\":\"foo\"}]", w2.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w2.Header().Get("Content-Type"))
}

func TestRenderSecureJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	r := Default(SecureJSONRenderType)
	r.Setup(data, "while(1);")
	err := r.Render(w)
	Recycle(SecureJSONRenderType, r)

	assert.Error(t, err)
}

func TestRenderJsonpJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	r := Default(JsonpJSONRenderType)
	r.Setup(data, "x")
	r.WriteContentType(w)
	Recycle(JsonpJSONRenderType, r)
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))

	w1 := httptest.NewRecorder()
	r1 := Default(JsonpJSONRenderType)
	r1.Setup(data, "x")
	err1 := r1.Render(w1)
	// Recycle(JsonpJSONRenderType, r1)

	assert.NoError(t, err1)
	assert.Equal(t, "x({\"foo\":\"bar\"})", w1.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]interface{}{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	r2 := Default(JsonpJSONRenderType)
	r2.Setup(datas, "x")
	err2 := r2.Render(w2)
	Recycle(JsonpJSONRenderType, r2)
	assert.NoError(t, err2)
	assert.Equal(t, "x([{\"foo\":\"bar\"},{\"bar\":\"foo\"}])", w2.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w2.Header().Get("Content-Type"))
}

func TestRenderJsonpJSONError2(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}
	r := Default(JsonpJSONRenderType)
	r.Setup(data, "")
	r.WriteContentType(w)
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data, "")
	e := r.Render(w)
	Recycle(JsonpJSONRenderType, r)
	assert.NoError(t, e)

	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJsonpJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	r := Default(JsonpJSONRenderType)
	r.Setup(data, "x")
	// json: unsupported type: chan int
	err := r.Render(w)
	Recycle(JsonpJSONRenderType, r)
	assert.Error(t, err)
}

func TestRenderAsciiJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data1 := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	r := Default(AsciiJSONRenderType)
	r.Setup(data1)
	err := r.Render(w1)

	assert.NoError(t, err)
	assert.Equal(t, "{\"lang\":\"GO\\u8bed\\u8a00\",\"tag\":\"\\u003cbr\\u003e\"}", w1.Body.String())
	assert.Equal(t, "application/json", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	data2 := float64(3.1415926)

	r.Reset()
	r.Setup(data2)
	err = r.Render(w2)
	Recycle(AsciiJSONRenderType, r)
	assert.NoError(t, err)
	assert.Equal(t, "3.1415926", w2.Body.String())
}

func TestRenderAsciiJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	r := Default(AsciiJSONRenderType)
	r.Setup(data)
	// json: unsupported type: chan int
	assert.Error(t, r.Render(w))
	Recycle(AsciiJSONRenderType, r)
}

type xmlmap map[string]interface{}

// Allows type H to be used with xml.Marshal
func (h xmlmap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func TestRenderYAML(t *testing.T) {
	w := httptest.NewRecorder()
	data := `
a : Easy!
b:
	c: 2
	d: [3, 4]
	`
	r := Default(YAMLRenderType)
	r.Setup(data)
	r.WriteContentType(w)
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data)
	err := r.Render(w)
	Recycle(YAMLRenderType, r)
	assert.NoError(t, err)
	assert.Equal(t, "\"\\na : Easy!\\nb:\\n\\tc: 2\\n\\td: [3, 4]\\n\\t\"\n", w.Body.String())
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))
}

type fail struct{}

// Hook MarshalYAML
func (ft *fail) MarshalYAML() (interface{}, error) {
	return nil, errors.New("fail")
}

func TestRenderYAMLFail(t *testing.T) {
	w := httptest.NewRecorder()
	r := Default(YAMLRenderType)
	r.Setup(&fail{})
	err := r.Render(w)
	Recycle(YAMLRenderType, r)
	assert.Error(t, err)
}

// test Protobuf rendering
func TestRenderProtoBuf(t *testing.T) {
	w := httptest.NewRecorder()
	reps := []int64{int64(1), int64(2)}
	label := "test"
	data := &testdata.Test{
		Label: &label,
		Reps:  reps,
	}

	r := Default(ProtoBufRenderType)
	r.Setup(data)
	r.WriteContentType(w)
	protoData, err := proto.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data)
	err = r.Render(w)
	Recycle(ProtoBufRenderType, r)

	assert.NoError(t, err)
	assert.Equal(t, string(protoData), w.Body.String())
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))
}

func TestRenderProtoBufFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := &testdata.Test{}
	r := Default(ProtoBufRenderType)
	r.Setup(data)
	err := r.Render(w)
	Recycle(ProtoBufRenderType, r)
	assert.Error(t, err)
}

func TestRenderXML(t *testing.T) {
	w := httptest.NewRecorder()
	data := xmlmap{
		"foo": "bar",
	}

	r := Default(XMLRenderType)
	r.Setup(data)
	r.WriteContentType(w)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))

	r.Reset()
	r.Setup(data)
	err := r.Render(w)
	Recycle(XMLRenderType, r)

	assert.NoError(t, err)
	assert.Equal(t, "<map><foo>bar</foo></map>", w.Body.String())
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/test-redirect", nil)
	assert.NoError(t, err)

	data1 := Redirect{
		Code:     http.StatusMovedPermanently,
		Request:  req,
		Location: "/new/location",
	}

	w := httptest.NewRecorder()
	err = data1.Render(w)
	assert.NoError(t, err)

	data2 := Redirect{
		Code:     http.StatusOK,
		Request:  req,
		Location: "/new/location",
	}

	w = httptest.NewRecorder()
	assert.Panics(t, func() { data2.Render(w) })

	// only improve coverage
	data2.WriteContentType(w)
}

func TestRenderData(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("#!PNG some raw data")

	err := (Data{
		ContentType: "image/png",
		Data:        data,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "#!PNG some raw data", w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
}

func TestRenderString(t *testing.T) {
	w := httptest.NewRecorder()

	(String{
		Format: "hello %s %d",
		Data:   []interface{}{},
	}).WriteContentType(w)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	err := (String{
		Format: "hola %s %d",
		Data:   []interface{}{"manu", 2},
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "hola manu 2", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderStringLenZero(t *testing.T) {
	w := httptest.NewRecorder()

	err := (String{
		Format: "hola %s %d",
		Data:   []interface{}{},
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "hola %s %d", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("t", map[string]interface{}{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "Hello alexandernyquist", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLTemplateEmptyName(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("").Parse(`Hello {{.name}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("", map[string]interface{}{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "Hello alexandernyquist", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugFiles(t *testing.T) {
	w := httptest.NewRecorder()
	htmlRender := HTMLDebug{Files: []string{"../testdata/template/hello.tmpl"},
		Glob:    "",
		Delims:  Delims{Left: "{[{", Right: "}]}"},
		FuncMap: nil,
	}
	instance := htmlRender.Instance("hello.tmpl", map[string]interface{}{
		"name": "thinkerou",
	})

	err := instance.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "<h1>Hello thinkerou</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugGlob(t *testing.T) {
	w := httptest.NewRecorder()
	htmlRender := HTMLDebug{Files: nil,
		Glob:    "../testdata/template/hello*",
		Delims:  Delims{Left: "{[{", Right: "}]}"},
		FuncMap: nil,
	}
	instance := htmlRender.Instance("hello.tmpl", map[string]interface{}{
		"name": "thinkerou",
	})

	err := instance.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "<h1>Hello thinkerou</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugPanics(t *testing.T) {
	htmlRender := HTMLDebug{Files: nil,
		Glob:    "",
		Delims:  Delims{"{{", "}}"},
		FuncMap: nil,
	}
	assert.Panics(t, func() { htmlRender.Instance("", nil) })
}

func TestRenderReader(t *testing.T) {
	w := httptest.NewRecorder()

	body := "#!PNG some raw data"
	headers := make(map[string]string)
	headers["Content-Disposition"] = `attachment; filename="filename.png"`

	err := (Reader{
		ContentLength: int64(len(body)),
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(len(body)), w.Header().Get("Content-Length"))
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
}
