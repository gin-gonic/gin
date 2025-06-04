package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/examples/url/ini"
)

func main() {
	fmt.Println("Hello Gin")
	router := gin.Default()
	ini.UrlInit(router)
	err := router.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
