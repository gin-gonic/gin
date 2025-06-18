// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import "net/http"

// CreateTestContext returns a fresh Engine and a Context associated with it.
// This is useful for tests that need to set up a new Gin engine instance
// along with a context, for example, to test middleware that doesn't depend on
// specific routes. The ResponseWriter `w` is used to initialize the context's writer.
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
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	c = r.allocateContext(r.maxParams)
	c.reset()
	c.writermem.reset(w)
	return
}
