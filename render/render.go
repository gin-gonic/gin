// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

type (
	Render interface {
		Render(http.ResponseWriter, int, ...interface{}) error
	}

	jsonRender struct{}

	indentedJSON struct{}

	xmlRender struct{}

	plainTextRender struct{}

	htmlPlainRender struct{}

	redirectRender struct{}

	HTMLRender struct {
		Template *template.Template
	}
)

var (
	JSON         = jsonRender{}
	IndentedJSON = indentedJSON{}
	XML          = xmlRender{}
	HTMLPlain    = htmlPlainRender{}
	Plain        = plainTextRender{}
	Redirect     = redirectRender{}
)

func (_ redirectRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	req := data[0].(*http.Request)
	location := data[1].(string)
	http.Redirect(w, req, location, code)
	return nil
}

func (_ jsonRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "application/json")
	return json.NewEncoder(w).Encode(data[0])
}

func (_ indentedJSON) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "application/json")
	jsonData, err := json.MarshalIndent(data[0], "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}

func (_ xmlRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "application/xml")
	return xml.NewEncoder(w).Encode(data[0])
}

func (_ plainTextRender) Render(w http.ResponseWriter, code int, data ...interface{}) (err error) {
	WriteHeader(w, code, "text/plain")
	format := data[0].(string)
	args := data[1].([]interface{})
	if len(args) > 0 {
		_, err = fmt.Fprintf(w, format, args...)
	} else {
		_, err = w.Write([]byte(format))
	}
	return
}

func (_ htmlPlainRender) Render(w http.ResponseWriter, code int, data ...interface{}) (err error) {
	WriteHeader(w, code, "text/html")
	format := data[0].(string)
	args := data[1].([]interface{})
	if len(args) > 0 {
		_, err = fmt.Fprintf(w, format, args...)
	} else {
		_, err = w.Write([]byte(format))
	}
	return
}

func (html HTMLRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "text/html")
	file := data[0].(string)
	args := data[1]
	return html.Template.ExecuteTemplate(w, file, args)
}

func WriteHeader(w http.ResponseWriter, code int, contentType string) {
	contentType = joinStrings(contentType, "; charset=utf-8")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func joinStrings(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := 0
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	n = 0
	for _, s := range a {
		n += copy(b[n:], s)
	}
	return string(b)
}
