// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

// Delims represents a set of Left and Right delimiters for HTML template rendering.
type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(string, interface{}) Render
	// Add new template files to this instance with glob
	ParseGlob(pattern string)
	// Add new template files to this instance
	ParseFiles(files ...string)
}

// HTMLProduction contains template reference and its delims.
type HTMLProduction struct {
	Template *template.Template
}

// HTMLDebug contains template delims and pattern and function with file list.
type HTMLDebug struct {
	Files   []string
	Globs   []string
	Delims  Delims
	FuncMap template.FuncMap
}

// HTML contains template reference and its name with given interface object.
type HTML struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

var htmlContentType = []string{"text/html; charset=utf-8"}

// Instance (HTMLProduction) returns an HTML instance which it realizes Render interface.
func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

// Add new template files to this instance (HTMLProduction) with glob
func (r *HTMLProduction) ParseGlob(pattern string) {
	template.Must(r.Template.ParseGlob(pattern))
}

// Add new template files to this instance (HTMLProduction) with glob
func (r *HTMLProduction) ParseFiles(files ...string) {
	template.Must(r.Template.ParseFiles(files...))
}

// Instance (HTMLDebug) returns an HTML instance which it realizes Render interface.
func (r HTMLDebug) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

// Add new template files to this instance (HTMLProduction) with glob
func (r *HTMLDebug) ParseGlob(pattern string) {
	r.Globs = append(r.Globs, pattern)
}

// Add new template files to this instance (HTMLDebug) with glob
func (r *HTMLDebug) ParseFiles(files ...string) {
	r.Files = append(r.Files, files...)
}

func (r HTMLDebug) loadTemplate() *template.Template {
	if r.FuncMap == nil {
		r.FuncMap = template.FuncMap{}
	}
	if len(r.Files) == 0 && len(r.Globs) == 0 {
		panic("the HTML debug render was created without files or glob pattern")
	}
	t := template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap)
	if len(r.Files) > 0 {
		template.Must(t.ParseFiles(r.Files...))
	}
	for _, glob := range r.Globs {
		template.Must(t.ParseGlob(glob))
	}
	return t
}

// Render (HTML) executes template and writes its result with custom ContentType for response.
func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
