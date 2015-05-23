package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type NgHTML struct {
	Template map[string]*pongo2.Template
	Name     string
	Data     interface{}
}

func (n NgHTML) Write(w http.ResponseWriter) error {
	file := n.Name
	ctx := n.Data.(pongo2.Context)

	var t *pongo2.Template

	if tmpl, ok := n.Template[file]; ok {
		t = tmpl
	} else {
		tmpl, err := pongo2.FromFile(file)
		if err != nil {
			return err
		}
		n.Template[file] = tmpl
		t = tmpl
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteWriter(ctx, w)
}

type NgPongRender struct {
	Template map[string]*pongo2.Template
}

func (n *NgPongRender) Instance(name string, data interface{}) render.Render {
	return NgHTML{
		Template: n.Template,
		Name:     name,
		Data:     data,
	}
}
func NewNgPongRender() *NgPongRender {
	return &NgPongRender{Template: map[string]*pongo2.Template{}}
}
