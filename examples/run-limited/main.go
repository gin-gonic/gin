package main

import "github.com/gin-gonic/gin"

const maxConnections = 10
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.RunLimited(maxConnections, ":80") // listen and serve on 0.0.0.0:8080
}
