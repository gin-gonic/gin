package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	r := gin.Default()

	gin.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Listen and Server in 0.0.0.0:443
	r.RunAutoTLS("example1.com", "example2.com")
}
