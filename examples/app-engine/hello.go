package hello

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func init() {
	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.GET("/", func(c *gin.Context){
		c.String(200, "Hello World!")
	})
	r.GET("/ping/", func(c *gin.Context){
		c.String(200, "pong")
	})

	// Handle all requests using net/http
	http.Handle("/", r)
}