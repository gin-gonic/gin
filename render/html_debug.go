// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

type HTMLDebugRender struct {
	files []string
	globs []string
}

func (r *HTMLDebugRender) AddGlob(pattern string) {
	r.globs = append(r.globs, pattern)
}

func (r *HTMLDebugRender) AddFiles(files ...string) {
	r.files = append(r.files, files...)
}

func (r *HTMLDebugRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(w, code, "text/html")
	file := data[0].(string)
	obj := data[1]

	if t, err := r.newTemplate(); err == nil {
		return t.ExecuteTemplate(w, file, obj)
	} else {
		return err
	}
}

func (r *HTMLDebugRender) newTemplate() (*template.Template, error) {
	t := template.New("")
	if len(r.files) > 0 {
		if _, err := t.ParseFiles(r.files...); err != nil {
			return nil, err
		}
	}
	for _, glob := range r.globs {
		if _, err := t.ParseGlob(glob); err != nil {
			return nil, err
		}
	}
	return t, nil
}
