package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgutz/ansi"
)

var (
	white  = ansi.ColorCode("white+h:black")
	red    = ansi.ColorCode("red+h:black")
	green  = ansi.ColorCode("green+h:black")
	yellow = ansi.ColorCode("yellow+h:black")
	blue   = ansi.ColorCode("blue+h:black")
	reset  = ansi.ColorCode("reset")
)

//
// Example of an extended ansi-colored logger using the
// ctx.Writer.Status() function
func logger(c *gin.Context) {
	start := time.Now()

	// save the IP of the requester
	requester := c.Req.Header.Get("X-Real-IP")

	// if the requester-header is empty, check the forwarded-header
	if requester == "" {
		requester = c.Req.Header.Get("X-Forwarded-For")
	}

	// if the requester is still empty, use the hard-coded address from the socket
	if requester == "" {
		requester = c.Req.RemoteAddr
	}

	// ... finally, log the fact we got a request
	log.Printf("<-- %16s | %6s | %s\n", requester, c.Req.Method, c.Req.URL.Path)

	c.Next()

	var color string
	if code := c.Writer.Status(); code >= 200 && code <= 299 {
		color = green
	} else if code >= 300 && code <= 399 {
		color = white
	} else if code >= 400 && code <= 499 {
		color = yellow
	} else {
		color = red
	}

	log.Printf("--> %s%16s | %6d | %s | %s%s\n",
		color,
		requester, c.Writer.Status(), time.Since(start), c.Req.URL.Path,
		reset,
	)
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger)
	// or modify func.logger to return a handler and use:
	// r.Use(logger())

	// Ping test
	r.GET("/:code", func(c *gin.Context) {
		asInt, err := strconv.ParseInt(c.Params.ByName("code"), 10, 32)
		if err != nil {
			c.String(400, err.Error())
		} else {
			c.String(int(asInt), c.Params.ByName("code"))
		}
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}
