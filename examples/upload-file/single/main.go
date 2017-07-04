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

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		src, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("file open err: %s", err.Error()))
			return
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Create file err: %s", err.Error()))
			return
		}
		defer dst.Close()

		// Copy
		io.Copy(dst, src)

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email))
	})
	router.Run(":8080")
}
