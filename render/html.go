package render

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type (
	HTMLRender struct {
		Template *template.Template
	}

	htmlPlainRender struct{}

	HTMLDebugRender struct {
		Files []string
		Glob  string
	}
)

func (html HTMLRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html; charset=utf-8")
	file := data[0].(string)
	args := data[1]
	return html.Template.ExecuteTemplate(w, file, args)
}

func (r *HTMLDebugRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html; charset=utf-8")
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

func (_ htmlPlainRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	format := data[0].(string)
	values := data[1].([]interface{})
	WriteHTMLString(w, code, format, values)
	return nil
}

func WriteHTMLString(w http.ResponseWriter, code int, format string, values []interface{}) {
	writeHeader(w, code, "text/html; charset=utf-8")
	if len(values) > 0 {
		fmt.Fprintf(w, format, values...)
	} else {
		w.Write([]byte(format))
	}
}
