// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"errors"
	"html/template"
	"net/http"
)

type HTMLDebugRender struct {
	Files []string
	Glob  string
}

func (r *HTMLDebugRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "text/html")
	file := data[0].(string)
	obj := data[1]

	if t, err := r.loadTemplate(); err == nil {
		return t.ExecuteTemplate(w, file, obj)
	} else {
		return err
	}
}

func (r *HTMLDebugRender) loadTemplate() (*template.Template, error) {
	if len(r.Files) > 0 {
		return template.ParseFiles(r.Files...)
	}
	if len(r.Glob) > 0 {
		return template.ParseGlob(r.Glob)
	}
	return nil, errors.New("the HTML debug render was created without files or glob pattern")
}
