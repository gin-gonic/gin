package gin

import (
	"fmt"
	"log"
	"time"
)

func ErrorLogger() HandlerFunc {
	return func(c *Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// -1 status code = do not change current one
			c.JSON(-1, c.Errors)
		}
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
		if len(c.Errors) > 0 {
			fmt.Println(c.Errors)
		}
	}
}
