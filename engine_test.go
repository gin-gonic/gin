package gin

import (
	"testing"
)

func TestRoutesConcurrent(t *testing.T) {
	r := New()

	done := make(chan bool)

	// Concurrently read routes
	go func() {
		_ = r.Routes()
		done <- true
	}()

	// Register a route at the same time
	r.GET("/", func(c *Context) { c.String(200, "OK") })

	<-done
}
