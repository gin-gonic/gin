package gin

import (
	"log"
	"time"
)

func ErrorLogger() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if len(c.Errors) > 0 {
				log.Println(c.Errors)
				c.JSON(-1, c.Errors)
			}
		}()
		c.Next()
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {

		// Start timer
		t := time.Now()

		// Process request
		c.Next()

		// Calculate resolution time
		log.Printf("%s in %v", c.Req.RequestURI, time.Since(t))
	}
}
