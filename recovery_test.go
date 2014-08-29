// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	// SETUP
	log.SetOutput(bytes.NewBuffer(nil)) // Disable panic logs for testing
	r := New()
	r.Use(Recovery())
	r.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	w := PerformRequest(r, "GET", "/recovery")

	// restore logging
	log.SetOutput(os.Stderr)

	if w.Code != 500 {
		t.Errorf("Response code should be Internal Server Error, was: %s", w.Code)
	}
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	// SETUP
	log.SetOutput(bytes.NewBuffer(nil))
	r := New()
	r.Use(Recovery())
	r.GET("/recovery", func(c *Context) {
		c.Abort(400)
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	w := PerformRequest(r, "GET", "/recovery")

	// restore logging
	log.SetOutput(os.Stderr)

	// TEST
	if w.Code != 500 {
		t.Errorf("Response code should be Bad request, was: %s", w.Code)
	}
}
