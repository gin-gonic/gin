// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"net/http"
	"time"
)

// CreateTestContext returns a fresh Engine and a Context associated with it.
// This is useful for tests that need to set up a new Gin engine instance
// along with a context, for example, to test middleware that doesn't depend on
// specific routes. The ResponseWriter `w` is used to initialize the context's writer.
//
// Unlike Engine.ServeHTTP, CreateTestContext does not flush the response at the
// end of the handler. The status code set by Context.Status is buffered until
// the response body is written (Context.Writer.Write/WriteString) or
// Context.Writer.WriteHeaderNow is called. To observe the status code in the
// ResponseWriter after a handler that only calls Context.Status, call
// c.Writer.WriteHeaderNow() explicitly.
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext(0)
	c.reset()
	c.writermem.reset(w)
	return
}

// CreateTestContextOnly returns a fresh Context associated with the provided Engine `r`.
// This is useful for tests that operate on an existing, possibly pre-configured,
// Gin engine instance and need a new context for it.
// The ResponseWriter `w` is used to initialize the context's writer.
// The context is allocated with the `maxParams` setting from the provided engine.
//
// Unlike Engine.ServeHTTP, CreateTestContextOnly does not flush the response at the
// end of the handler. See CreateTestContext for details on observing the status code.
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	c = r.allocateContext(r.maxParams)
	c.reset()
	c.writermem.reset(w)
	return
}

// waitForServerReady waits for a server to be ready by making HTTP requests
// with exponential backoff. This is more reliable than time.Sleep() for testing.
func waitForServerReady(url string, maxAttempts int) error {
	client := &http.Client{
		Timeout: 100 * time.Millisecond,
	}

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}

		// Exponential backoff: 10ms, 20ms, 40ms, 80ms, 160ms...
		backoff := min(time.Duration(10*(1<<uint(i)))*time.Millisecond, 500*time.Millisecond)
		time.Sleep(backoff)
	}

	return fmt.Errorf("server at %s did not become ready after %d attempts", url, maxAttempts)
}
