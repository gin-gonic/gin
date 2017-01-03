package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["files"]

		for _, file := range files {
			// Source
			src, _ := file.Open()
			defer src.Close()

			// Destination
			dst, _ := os.Create(file.Filename)
			defer dst.Close()

			// Copy
			io.Copy(dst, src)
		}

		c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email))
	})
	router.Run(":8080")
}
