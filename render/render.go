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

	// JSON binding
	jsonRender struct{}

	// XML binding
	xmlRender struct{}

	// Plain text
	plainRender struct{}

	// Redirects
	redirectRender struct{}

	// Redirects
	htmlDebugRender struct{}

	// form binding
	HTMLRender struct {
		Template *template.Template
	}
)

var (
	JSON      = jsonRender{}
	XML       = xmlRender{}
	Plain     = plainRender{}
	Redirect  = redirectRender{}
	HTMLDebug = htmlDebugRender{}
)

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func (_ jsonRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data[0])
}

func (_ redirectRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	w.Header().Set("Location", data[0].(string))
	w.WriteHeader(code)
	return nil
}

func (_ xmlRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "application/xml")
	encoder := xml.NewEncoder(w)
	return encoder.Encode(data[0])
}

func (_ plainRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/plain")
	format := data[0].(string)
	args := data[1].([]interface{})
	var err error
	if len(args) > 0 {
		_, err = w.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = w.Write([]byte(format))
	}
	return err
}

func (_ htmlDebugRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html")
	file := data[0].(string)
	obj := data[1]
	t, err := template.ParseFiles(file)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, file, obj)
}

func (html HTMLRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html")
	file := data[0].(string)
	obj := data[1]
	return html.Template.ExecuteTemplate(w, file, obj)
}
