package gin

import (
	"log"
	"time"
)

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
