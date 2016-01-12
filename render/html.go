// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/gin-gonic/gin/template"
)

// @modified by henry, 2016.1.12
var HTMLTemplate = template.New("gin")

type (
	HTMLRender interface {
		Instance(string, interface{}) Render
	}

	HTMLProduction struct {
		Template *template.Template
	}

	HTMLDebug struct {
		Files []string
		Glob  string
	}

	HTML struct {
		Template *template.Template
		Name     string
		Data     interface{}
	}
)

var htmlContentType = []string{"text/html; charset=utf-8"}

func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

func (r HTMLDebug) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

// @modified by henry, 2016.1.12
func (r HTMLDebug) loadTemplate() *template.Template {
	if len(r.Files) > 0 {
		return template.Must(template.Must(HTMLTemplate.Clone()).ParseFiles(r.Files...))
	}
	if len(r.Glob) > 0 {
		return template.Must(template.Must(HTMLTemplate.Clone()).ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

func (r HTML) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)
	if len(r.Name) == 0 {
		return r.Template.Execute(w, r.Data)
	} else {
		return r.Template.ExecuteTemplate(w, r.Name, r.Data)
	}
}
