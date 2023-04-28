// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"context"
	"net/http"
)

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext(0)
	c.reset()
	c.writermem.reset(w)
	return
}

// CreateTestContextOnly returns a fresh context base on the engine for testing purposes
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	c = r.allocateContext(r.maxParams)
	c.reset()
	c.writermem.reset(w)
	return
}

// CreateTestContextOnly returns a fresh context and its closer
func CreateTestContextWithCloser(w http.ResponseWriter) (c *Context, closeClient context.CancelFunc) {
	r := New()
	c = r.allocateContext(0)
	c.reset()
	c.writermem.reset(w)
	ctx, closeClient := context.WithCancel(context.Background())
	var req http.Request
	c.Request = req.WithContext(ctx)
	return c, closeClient
}
