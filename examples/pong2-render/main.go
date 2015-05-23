package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.HTMLRender = NewNgPongRender()
	pages := r.Group("/")
	{
		pages.GET("/", func(c *gin.Context) {
			ctx := pongo2.Context{
				"hello": "Hello Home Pages",
			}
			c.HTML(200, "templates/sites/home/index.html", ctx)
		})
		pages.GET("/user", func(c *gin.Context) {
			ctx := pongo2.Context{
				"hello": "Hello User Pages",
			}
			c.HTML(200, "templates/sites/user/index.html", ctx)
		})
	}
	admins := r.Group("/admin")
	{
		admins.GET("/", func(c *gin.Context) {
			ctx := pongo2.Context{
				"hello": "Hello Admin Dashboard Pages",
			}
			c.HTML(200, "templates/admins/dashboard/index.html", ctx)
		})
	}
	r.Run(":8080")
}
