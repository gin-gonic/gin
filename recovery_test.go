// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package gin

import (
	"bytes"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), "TestPanicInHandler")
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	router := New()
	router.Use(RecoveryWithWriter(nil))
	router.GET("/recovery", func(c *Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSource(t *testing.T) {
	bs := source(nil, 0)
	assert.Equal(t, []byte("???"), bs)

	in := [][]byte{
		[]byte("Hello world."),
		[]byte("Hi, gin.."),
	}
	bs = source(in, 10)
	assert.Equal(t, []byte("???"), bs)

	bs = source(in, 1)
	assert.Equal(t, []byte("Hello world."), bs)
}

func TestFunction(t *testing.T) {
	bs := function(1)
	assert.Equal(t, []byte("???"), bs)
}

// TestPanicWithBrokenPipe asserts that recovery specifically handles
// writing responses to broken pipes
func TestPanicWithBrokenPipe(t *testing.T) {
	const expectCode = 204

	expectMsgs := map[syscall.Errno]string{
		syscall.EPIPE:      "broken pipe",
		syscall.ECONNRESET: "connection reset by peer",
	}

	for errno, expectMsg := range expectMsgs {
		t.Run(expectMsg, func(t *testing.T) {

			var buf bytes.Buffer

			router := New()
			router.Use(RecoveryWithWriter(&buf))
			router.GET("/recovery", func(c *Context) {
				// Start writing response
				c.Header("X-Test", "Value")
				c.Status(expectCode)

				// Oops. Client connection closed
				e := &net.OpError{Err: &os.SyscallError{Err: errno}}
				panic(e)
			})
			// RUN
			w := performRequest(router, "GET", "/recovery")
			// TEST
			assert.Equal(t, expectCode, w.Code)
			assert.Contains(t, strings.ToLower(buf.String()), expectMsg)
		})
	}
}
