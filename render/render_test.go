// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.JSONEq(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
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
	assert.JSONEq(t, "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\"\n}", w.Body.String())
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
	assert.JSONEq(t, "{\"foo\":\"bar\"}", w1.Body.String())
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
	bufString    string
	ErrThreshold int // 1-based threshold. If 1, errors on the 1st Write call.
	writeCount   int
	*httptest.ResponseRecorder
}

var _ http.ResponseWriter = (*errorWriter)(nil)

func (w *errorWriter) Header() http.Header {
	if w.ResponseRecorder == nil {
		w.ResponseRecorder = httptest.NewRecorder()
	}
	return w.ResponseRecorder.Header()
}

func (w *errorWriter) WriteHeader(statusCode int) {
	if w.ResponseRecorder == nil {
		w.ResponseRecorder = httptest.NewRecorder()
	}
	w.ResponseRecorder.WriteHeader(statusCode)
}

func (w *errorWriter) Write(buf []byte) (int, error) {
	w.writeCount++
	if (w.bufString != "" && string(buf) == w.bufString) || (w.ErrThreshold > 0 && w.writeCount >= w.ErrThreshold) {
		return 0, errors.New(`write error`)
	}
	if w.ResponseRecorder == nil {
		w.ResponseRecorder = httptest.NewRecorder()
	}
	return w.ResponseRecorder.Write(buf)
}

func (w *errorWriter) reset() {
	w.writeCount = 0
	w.ResponseRecorder = httptest.NewRecorder()
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

	// error was returned while writing callback
	ew.reset()
	ew.ErrThreshold = 1
	err := jsonpJSON.Render(ew)
	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())

	// error was returned while writing "("
	ew.reset()
	ew.ErrThreshold = 2
	err = jsonpJSON.Render(ew)
	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())

	// error was returned while writing data
	ew.reset()
	ew.ErrThreshold = 3
	err = jsonpJSON.Render(ew)
	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())

	// error was returned while writing ");"
	ew.reset()
	ew.ErrThreshold = 4
	err = jsonpJSON.Render(ew)
	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())
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

	assert.JSONEq(t, "{\"foo\":\"bar\"}", w.Body.String())
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
	assert.JSONEq(t, "{\"lang\":\"GO\\u8bed\\u8a00\",\"tag\":\"\\u003cbr\\u003e\"}", w1.Body.String())
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
	assert.JSONEq(t, "{\"foo\":\"bar\",\"html\":\"<b>\"}\n", w.Body.String())
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

func TestRenderXMLError(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	err := (XML{data}).Render(w)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "xml: unsupported type: chan int")
}

func TestRenderPDF(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("%Test pdf content")

	pdf := PDF{data}

	pdf.WriteContentType(w)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))

	err := pdf.Render(w)
	require.NoError(t, err)

	assert.Equal(t, data, w.Body.Bytes())
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
}

func TestRenderRedirect(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/test-redirect", nil)
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
	assert.Equal(t, "19", w.Header().Get("Content-Length"))
}

func TestRenderDataContentLength(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		assert.NoError(t, err)

		data := Data{
			ContentType: "application/octet-stream",
			Data:        make([]byte, size),
		}
		assert.NoError(t, data.Render(w))
	}))
	t.Cleanup(srv.Close)

	for _, size := range []int{0, 1, 100, 100_000} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			resp, err := http.Get(srv.URL + "?size=" + strconv.Itoa(size))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, "application/octet-stream", resp.Header.Get("Content-Type"))
			assert.Equal(t, strconv.Itoa(size), resp.Header.Get("Content-Length"))

			actual, err := io.Copy(io.Discard, resp.Body)
			require.NoError(t, err)
			assert.EqualValues(t, size, actual)
		})
	}
}

func TestRenderDataError(t *testing.T) {
	ew := &errorWriter{
		ErrThreshold:     1,
		ResponseRecorder: httptest.NewRecorder(),
	}
	data := []byte("#!PNG some raw data")

	err := (Data{
		ContentType: "image/png",
		Data:        data,
	}).Render(ew)

	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())
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
		Files:      []string{"../testdata/template/hello.tmpl"},
		Glob:       "",
		FileSystem: nil,
		Patterns:   nil,
		Delims:     Delims{Left: "{[{", Right: "}]}"},
		FuncMap:    nil,
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
		Files:      nil,
		Glob:       "../testdata/template/hello*",
		FileSystem: nil,
		Patterns:   nil,
		Delims:     Delims{Left: "{[{", Right: "}]}"},
		FuncMap:    nil,
	}
	instance := htmlRender.Instance("hello.tmpl", map[string]any{
		"name": "thinkerou",
	})

	err := instance.Render(w)

	require.NoError(t, err)
	assert.Equal(t, "<h1>Hello thinkerou</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRenderHTMLDebugFS(t *testing.T) {
	w := httptest.NewRecorder()
	htmlRender := HTMLDebug{
		Files:      nil,
		Glob:       "",
		FileSystem: http.Dir("../testdata/template"),
		Patterns:   []string{"hello.tmpl"},
		Delims:     Delims{Left: "{[{", Right: "}]}"},
		FuncMap:    nil,
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
		Files:      nil,
		Glob:       "",
		FileSystem: nil,
		Patterns:   nil,
		Delims:     Delims{"{{", "}}"},
		FuncMap:    nil,
	}
	assert.Panics(t, func() { htmlRender.Instance("", nil) })
}

func TestRenderHTMLTemplateError(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("t").Parse(`Hello {{if .name}}{{.name.DoesNotExist}}{{end}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("t", map[string]any{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)
	require.Error(t, err)
}

func TestRenderHTMLTemplateExecuteError(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("t").Parse(`Hello {{.name.invalid}}`))

	htmlRender := HTMLProduction{Template: templ}
	instance := htmlRender.Instance("t", map[string]any{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)
	require.Error(t, err)
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
	data := []any{"value1", "value2"}
	prefix := "my-prefix:"
	r := SecureJSON{Data: data, Prefix: prefix}
	ew := &errorWriter{
		ErrThreshold:     1,
		ResponseRecorder: httptest.NewRecorder(),
	}
	err := r.Render(ew)
	require.Error(t, err)
	assert.Equal(t, "write error", err.Error())
}
