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
		log.Printf("[%d] %s in %v", c.Writer.Status(), c.Req.RequestURI, time.Since(t))
	}
}
