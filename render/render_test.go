// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"errors"
	"html/template"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/internal/json"
	testdata "github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// TODO unit tests
// test errors

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo":  "bar",
		"html": "<b>",
	}

	(JSON{data}).WriteContentType(w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	err := (JSON{data}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	require.Error(t, (JSON{data}).Render(w))
}

func TestRenderIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
		"bar": "foo",
	}

	err := (IndentedJSON{data}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\"\n}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderIndentedJSONPanics(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (IndentedJSON{data}).Render(w)
	require.Error(t, err)
}

func TestRenderSecureJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
	}

	(SecureJSON{"while(1);", data}).WriteContentType(w1)
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	err1 := (SecureJSON{"while(1);", data}).Render(w1)

	require.NoError(t, err1)
	assert.Equal(t, "{\"foo\":\"bar\"}", w1.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]any{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	err2 := (SecureJSON{"while(1);", datas}).Render(w2)
	require.NoError(t, err2)
	assert.Equal(t, "while(1);[{\"foo\":\"bar\"},{\"bar\":\"foo\"}]", w2.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w2.Header().Get("Content-Type"))
}

func TestRenderSecureJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (SecureJSON{"while(1);", data}).Render(w)
	require.Error(t, err)
}

func TestRenderJsonpJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
	}

	(JsonpJSON{"x", data}).WriteContentType(w1)
	assert.Equal(t, "application/javascript; charset=utf-8", w1.Header().Get("Content-Type"))

	err1 := (JsonpJSON{"x", data}).Render(w1)

	require.NoError(t, err1)
	assert.Equal(t, "x({\"foo\":\"bar\"});", w1.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	datas := []map[string]any{{
		"foo": "bar",
	}, {
		"bar": "foo",
	}}

	err2 := (JsonpJSON{"x", datas}).Render(w2)
	require.NoError(t, err2)
	assert.Equal(t, "x([{\"foo\":\"bar\"},{\"bar\":\"foo\"}]);", w2.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w2.Header().Get("Content-Type"))
}

type errorWriter struct {
	bufString string
	*httptest.ResponseRecorder
}

var _ http.ResponseWriter = (*errorWriter)(nil)

func (w *errorWriter) Write(buf []byte) (int, error) {
	if string(buf) == w.bufString {
		return 0, errors.New(`write "` + w.bufString + `" error`)
	}
	return w.ResponseRecorder.Write(buf)
}

func TestRenderJsonpJSONError(t *testing.T) {
	ew := &errorWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}

	jsonpJSON := JsonpJSON{
		Callback: "foo",
		Data: map[string]string{
			"foo": "bar",
		},
	}

	cb := template.JSEscapeString(jsonpJSON.Callback)
	ew.bufString = cb
	err := jsonpJSON.Render(ew) // error was returned while writing callback
	assert.Equal(t, `write "`+cb+`" error`, err.Error())

	ew.bufString = `(`
	err = jsonpJSON.Render(ew)
	assert.Equal(t, `write "`+`(`+`" error`, err.Error())

	data, _ := json.Marshal(jsonpJSON.Data) // error was returned while writing data
	ew.bufString = string(data)
	err = jsonpJSON.Render(ew)
	assert.Equal(t, `write "`+string(data)+`" error`, err.Error())

	ew.bufString = `);`
	err = jsonpJSON.Render(ew)
	assert.Equal(t, `write "`+`);`+`" error`, err.Error())
}

func TestRenderJsonpJSONError2(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
	}
	(JsonpJSON{"", data}).WriteContentType(w)
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))

	e := (JsonpJSON{"", data}).Render(w)
	require.NoError(t, e)

	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderJsonpJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	err := (JsonpJSON{"x", data}).Render(w)
	require.Error(t, err)
}

func TestRenderAsciiJSON(t *testing.T) {
	w1 := httptest.NewRecorder()
	data1 := map[string]any{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	err := (AsciiJSON{data1}).Render(w1)

	require.NoError(t, err)
	assert.Equal(t, "{\"lang\":\"GO\\u8bed\\u8a00\",\"tag\":\"\\u003cbr\\u003e\"}", w1.Body.String())
	assert.Equal(t, "application/json", w1.Header().Get("Content-Type"))

	w2 := httptest.NewRecorder()
	data2 := 3.1415926

	err = (AsciiJSON{data2}).Render(w2)
	require.NoError(t, err)
	assert.Equal(t, "3.1415926", w2.Body.String())
}

func TestRenderAsciiJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	// json: unsupported type: chan int
	require.Error(t, (AsciiJSON{data}).Render(w))
}

func TestRenderPureJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo":  "bar",
		"html": "<b>",
	}
	err := (PureJSON{data}).Render(w)
	require.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"<b>\"}\n", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

type xmlmap map[string]any

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
	(YAML{data}).WriteContentType(w)
	assert.Equal(t, "application/yaml; charset=utf-8", w.Header().Get("Content-Type"))

	err := (YAML{data}).Render(w)
	require.NoError(t, err)
	assert.Equal(t, "|4-\n    a : Easy!\n    b:\n    \tc: 2\n    \td: [3, 4]\n    \t\n", w.Body.String())
	assert.Equal(t, "application/yaml; charset=utf-8", w.Header().Get("Content-Type"))
}

type fail struct{}

// Hook MarshalYAML
func (ft *fail) MarshalYAML() (any, error) {
	return nil, errors.New("fail")
}

func TestRenderYAMLFail(t *testing.T) {
	w := httptest.NewRecorder()
	err := (YAML{&fail{}}).Render(w)
	require.Error(t, err)
}

func TestRenderTOML(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo":  "bar",
		"html": "<b>",
	}
	(TOML{data}).WriteContentType(w)
	assert.Equal(t, "application/toml; charset=utf-8", w.Header().Get("Content-Type"))

	err := (TOML{data}).Render(w)
	require.NoError(t, err)
	assert.Equal(t, "foo = 'bar'\nhtml = '<b>'\n", w.Body.String())
	assert.Equal(t, "application/toml; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderTOMLFail(t *testing.T) {
	w := httptest.NewRecorder()
	err := (TOML{net.IPv4bcast}).Render(w)
	require.Error(t, err)
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

	(ProtoBuf{data}).WriteContentType(w)
	protoData, err := proto.Marshal(data)
	require.NoError(t, err)
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))

	err = (ProtoBuf{data}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, string(protoData), w.Body.String())
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))
}

func TestRenderProtoBufFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := &testdata.Test{}
	err := (ProtoBuf{data}).Render(w)
	require.Error(t, err)
}

func TestRenderXML(t *testing.T) {
	w := httptest.NewRecorder()
	data := xmlmap{
		"foo": "bar",
	}

	(XML{data}).WriteContentType(w)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))

	err := (XML{data}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, "<map><foo>bar</foo></map>", w.Body.String())
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/test-redirect", nil)
	require.NoError(t, err)

	data1 := Redirect{
		Code:     http.StatusMovedPermanently,
		Request:  req,
		Location: "/new/location",
	}

	w := httptest.NewRecorder()
	err = data1.Render(w)
	require.NoError(t, err)

	data2 := Redirect{
		Code:     http.StatusOK,
		Request:  req,
		Location: "/new/location",
	}

	w = httptest.NewRecorder()
	assert.PanicsWithValue(t, "Cannot redirect with status code 200", func() {
		err := data2.Render(w)
		require.NoError(t, err)
	})

	data3 := Redirect{
		Code:     http.StatusCreated,
		Request:  req,
		Location: "/new/location",
	}

	w = httptest.NewRecorder()
	err = data3.Render(w)
	require.NoError(t, err)

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

	require.NoError(t, err)
	assert.Equal(t, "#!PNG some raw data", w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
}

func TestRenderString(t *testing.T) {
	w := httptest.NewRecorder()

	(String{
		Format: "hello %s %d",
		Data:   []any{},
	}).WriteContentType(w)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	err := (String{
		Format: "hola %s %d",
		Data:   []any{"manu", 2},
	}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, "hola manu 2", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderStringLenZero(t *testing.T) {
	w := httptest.NewRecorder()

	err := (String{
		Format: "hola %s %d",
		Data:   []any{},
	}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, "hola %s %d", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("t", map[string]any{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)

	require.NoError(t, err)
	assert.Equal(t, "Hello alexandernyquist", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLTemplateEmptyName(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("").Parse(`Hello {{.name}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("", map[string]any{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)

	require.NoError(t, err)
	assert.Equal(t, "Hello alexandernyquist", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugFiles(t *testing.T) {
	w := httptest.NewRecorder()
	htmlRender := HTMLDebug{
		Files:   []string{"../testdata/template/hello.tmpl"},
		Glob:    "",
		Delims:  Delims{Left: "{[{", Right: "}]}"},
		FuncMap: nil,
	}
	instance := htmlRender.Instance("hello.tmpl", map[string]any{
		"name": "thinkerou",
	})

	err := instance.Render(w)

	require.NoError(t, err)
	assert.Equal(t, "<h1>Hello thinkerou</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugGlob(t *testing.T) {
	w := httptest.NewRecorder()
	htmlRender := HTMLDebug{
		Files:   nil,
		Glob:    "../testdata/template/hello*",
		Delims:  Delims{Left: "{[{", Right: "}]}"},
		FuncMap: nil,
	}
	instance := htmlRender.Instance("hello.tmpl", map[string]any{
		"name": "thinkerou",
	})

	err := instance.Render(w)

	require.NoError(t, err)
	assert.Equal(t, "<h1>Hello thinkerou</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugPanics(t *testing.T) {
	htmlRender := HTMLDebug{
		Files:   nil,
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
	headers["x-request-id"] = "requestId"

	err := (Reader{
		ContentLength: int64(len(body)),
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(len(body)), w.Header().Get("Content-Length"))
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
	assert.Equal(t, headers["x-request-id"], w.Header().Get("x-request-id"))
}

func TestRenderReaderNoContentLength(t *testing.T) {
	w := httptest.NewRecorder()

	body := "#!PNG some raw data"
	headers := make(map[string]string)
	headers["Content-Disposition"] = `attachment; filename="filename.png"`
	headers["x-request-id"] = "requestId"

	err := (Reader{
		ContentLength: -1,
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	require.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.NotContains(t, "Content-Length", w.Header())
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
	assert.Equal(t, headers["x-request-id"], w.Header().Get("x-request-id"))
}

func TestRenderWriteError(t *testing.T) {
	data := []interface{}{"value1", "value2"}
	prefix := "my-prefix:"
	r := SecureJSON{Data: data, Prefix: prefix}
	ew := &errorWriter{
		bufString:        prefix,
		ResponseRecorder: httptest.NewRecorder(),
	}
	err := r.Render(ew)
	require.Error(t, err)
	assert.Equal(t, `write "my-prefix:" error`, err.Error())
}
