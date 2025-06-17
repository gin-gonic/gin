package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/examples/url/ini"
)

func main() {
	fmt.Println("h")
	fmt.Println("Hello Gin")
	// Force log's color
	gin.ForceConsoleColor()
	// disable log`s color
	//gin.DisableConsoleColor()
	router := gin.Default()
	ini.UrlInit(router)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second, // ReadTimeout
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
