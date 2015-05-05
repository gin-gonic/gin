// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	err := JSON.Render(w, 201, map[string]interface{}{
		"foo": "bar",
	})

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

func TestRenderPlain(t *testing.T) {
	w := httptest.NewRecorder()
	err := Plain.Render(w, 400, "hola %s %d", []interface{}{"manu", 2})

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

func TestRenderJoinStrings(t *testing.T) {
	assert.Equal(t, joinStrings("a", "BB", "c"), "aBBc")
	assert.Equal(t, joinStrings("a", "", "c"), "ac")
	assert.Equal(t, joinStrings("text/html", "; charset=utf-8"), "text/html; charset=utf-8")
}
