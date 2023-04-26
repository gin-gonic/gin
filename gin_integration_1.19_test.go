// Copyright 2023 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build go1.19

package gin

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunQUIC(t *testing.T) {
	router := New()
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

		assert.NoError(t, router.RunQUIC(":8443", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	}()

	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.RunQUIC(":8443", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	testRequest(t, "https://localhost:8443/example")
}
