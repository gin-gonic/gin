// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

type TemplateStorage struct {
	Storage map[string]*template.Template
}

func (t TemplateStorage) Instance(name string, data interface{}) Render {
	return HTMLWithBlock{
		Template: t.Storage[name],
		Name:     name,
		Data:     data,
	}
}

type HTMLWithBlock struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

func (r HTMLWithBlock) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)
	return r.Template.Execute(w, r.Data)
}
