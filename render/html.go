package render

import (
	"html/template"
	"net/http"
)

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
func (r HTMLDebug) loadTemplate() *template.Template {
	if len(r.Files) > 0 {
		return template.Must(template.ParseFiles(r.Files...))
	}
	if len(r.Glob) > 0 {
		return template.Must(template.ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

func (r HTML) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}
