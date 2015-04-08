package main

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	router := gin.Default()
	router.HTMLRender = newPongoRender()

	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Gin meets pongo2 !",
			"name":  c.Input.Get("name"),
		})
	})
	router.Run(":8080")
}

type pongoRender struct {
	cache map[string]*pongo2.Template
}

func newPongoRender() *pongoRender {
	return &pongoRender{map[string]*pongo2.Template{}}
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
	render.WriteHeader(w, code, "text/html")
	return t.ExecuteWriter(ctx, w)
}
