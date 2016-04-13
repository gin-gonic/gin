package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.AddHTMLTemplate("/", "layouts/default.html", "views/blocks1.html")
	router.AddHTMLTemplate("/blocks2/", "layouts/default.html", "views/blocks2.html")

	router.GET("/", render)
	router.GET("/blocks2/", render)
	router.Run(":8080")

}

func render(c *gin.Context) {
	c.HTML(200, c.Request.URL.Path, gin.H{})
}
