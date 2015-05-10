// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO unit tests
// test errors

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo": "bar",
	}

	err := JSON.Render(w, 201, data)
	WriteJSON(w2, 201, data)

	assert.Equal(t, w, w2)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "{\"foo\":\"bar\"}\n")
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}

func TestRenderIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	err := IndentedJSON.Render(w, 202, map[string]interface{}{
		"foo": "bar",
		"bar": "foo",
	})

	assert.NoError(t, err)
	assert.Equal(t, w.Code, 202)
	assert.Equal(t, w.Body.String(), "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\"\n}")
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
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
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}

func TestRenderXML(t *testing.T) {
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	data := xmlmap{
		"foo": "bar",
	}

	err := XML.Render(w, 200, data)
	WriteXML(w2, 200, data)

	assert.Equal(t, w, w2)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "<map><foo>bar</foo></map>")
	assert.Equal(t, w.Header().Get("Content-Type"), "application/xml; charset=utf-8")
}

func TestRenderRedirect(t *testing.T) {
	// TODO
}

func TestRenderData(t *testing.T) {
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	data := []byte("#!PNG some raw data")

	err := Data.Render(w, 400, "image/png", data)
	WriteData(w2, 400, "image/png", data)

	assert.Equal(t, w, w2)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "#!PNG some raw data")
	assert.Equal(t, w.Header().Get("Content-Type"), "image/png")
}

func TestRenderPlain(t *testing.T) {
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	err := Plain.Render(w, 400, "hola %s %d", []interface{}{"manu", 2})
	WritePlainText(w2, 400, "hola %s %d", []interface{}{"manu", 2})

	assert.Equal(t, w, w2)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, 400)
	assert.Equal(t, w.Body.String(), "hola manu 2")
	assert.Equal(t, w.Header().Get("Content-Type"), "text/plain; charset=utf-8")
}

func TestRenderPlainHTML(t *testing.T) {
	w := httptest.NewRecorder()
	err := HTMLPlain.Render(w, 401, "hola %s %d", []interface{}{"manu", 2})

	assert.NoError(t, err)
	assert.Equal(t, w.Code, 401)
	assert.Equal(t, w.Body.String(), "hola manu 2")
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}

func TestRenderHTMLTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
	htmlRender := HTMLRender{Template: templ}
	err := htmlRender.Render(w, 402, "t", map[string]interface{}{
		"name": "alexandernyquist",
	})

	assert.NoError(t, err)
	assert.Equal(t, w.Code, 402)
	assert.Equal(t, w.Body.String(), "Hello alexandernyquist")
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}
