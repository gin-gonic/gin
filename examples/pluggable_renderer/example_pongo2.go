package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pongoRender struct {
	cache map[string]*pongo2.Template
}

func newPongoRender() *pongoRender {
	return &pongoRender{map[string]*pongo2.Template{}}
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	if code >= 0 {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
	}
}

func (p *pongoRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	file := data[0].(string)
	ctx := data[1].(pongo2.Context)
	var t *pongo2.Template

	if tmpl, ok := p.cache[file]; ok {
		t = tmpl
	} else {
		tmpl, err := pongo2.FromFile(file)
		if err != nil {
			return err
		}
		p.cache[file] = tmpl
		t = tmpl
	}
	writeHeader(w, code, "text/html")
	return t.ExecuteWriter(ctx, w)
}

func main() {
	r := gin.Default()
	r.HTMLRender = newPongoRender()

	r.GET("/index", func(c *gin.Context) {
		name := c.Request.FormValue("name")
		ctx := pongo2.Context{
			"title": "Gin meets pongo2 !",
			"name":  name,
		}
		c.HTML(200, "index.html", ctx)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
