// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicClean(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	password := "my-super-secret-password"
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(c *Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery",
		header{
			Key:   "Host",
			Value: "www.google.com",
		},
		header{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", password),
		},
		header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check the buffer does not have the secret key
	assert.NotContains(t, buffer.String(), password)
}

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
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	SetMode(TestMode)
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

func TestCustomRecoveryWithWriter(t *testing.T) {
	errBuffer := new(bytes.Buffer)
	buffer := new(bytes.Buffer)
	router := New()
	handleRecovery := func(c *Context, err interface{}) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(CustomRecoveryWithWriter(buffer, handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oupps, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}

func TestCustomRecovery(t *testing.T) {
	errBuffer := new(bytes.Buffer)
	buffer := new(bytes.Buffer)
	router := New()
	DefaultErrorWriter = buffer
	handleRecovery := func(c *Context, err interface{}) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(CustomRecovery(handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oupps, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}

func TestRecoveryWithWriterWithCustomRecovery(t *testing.T) {
	errBuffer := new(bytes.Buffer)
	buffer := new(bytes.Buffer)
	router := New()
	DefaultErrorWriter = buffer
	handleRecovery := func(c *Context, err interface{}) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(RecoveryWithWriter(DefaultErrorWriter, handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oupps, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}
