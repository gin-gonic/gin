// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimitedListenerUnderLimit(t *testing.T) {
	router := New(WithMaxConns(10))
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	server := httptest.NewServer(router.Handler())
	defer server.Close()

	// Should be able to make requests under limit
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestLimitedListenerAtLimit(t *testing.T) {
	// Create a server that holds connections
	var activeConns atomic.Int32
	var wg sync.WaitGroup
	block := make(chan struct{})

	router := New(WithMaxConns(2))
	router.GET("/", func(c *Context) {
		activeConns.Add(1)
		wg.Add(1)
		<-block // Block until test releases
		activeConns.Add(-1)
		wg.Done()
		c.String(http.StatusOK, "ok")
	})

	server := httptest.NewServer(router.Handler())
	defer server.Close()

	// Start 2 requests that will block
	for i := 0; i < 2; i++ {
		go func() {
			resp, err := http.Get(server.URL)
			if err == nil {
				resp.Body.Close()
			}
		}()
	}

	// Wait for both connections to be active
	require.Eventually(t, func() bool {
		return activeConns.Load() == 2
	}, 2*time.Second, 10*time.Millisecond)

	// Third request should be rejected immediately
	client := &http.Client{Timeout: 500 * time.Millisecond}
	_, err := client.Get(server.URL)
	// Connection should be rejected
	require.Error(t, err, "expected connection to be rejected")

	// Release the blocked connections
	close(block)
	wg.Wait()
}

func TestLimitedConnRelease(t *testing.T) {
	block := make(chan struct{})

	router := New(WithMaxConns(1))
	router.GET("/", func(c *Context) {
		<-block
		c.String(http.StatusOK, "ok")
	})

	server := httptest.NewServer(router.Handler())
	defer server.Close()

	// Start one blocking request
	var firstDone atomic.Bool
	go func() {
		resp, err := http.Get(server.URL)
		if err == nil {
			resp.Body.Close()
			firstDone.Store(true)
		}
	}()

	// Give the first request time to start
	time.Sleep(100 * time.Millisecond)

	// Second request should fail
	client := &http.Client{Timeout: 200 * time.Millisecond}
	_, err := client.Get(server.URL)
	require.Error(t, err, "expected connection to be rejected when limit reached")

	// Release the first connection
	close(block)

	// Eventually first request should complete
	require.Eventually(t, firstDone.Load, 2*time.Second, 10*time.Millisecond)
}

func TestLimitedListenerZeroLimit(t *testing.T) {
	router := New(WithMaxConns(0))
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	server := httptest.NewServer(router.Handler())
	defer server.Close()

	// Should allow unlimited connections (zero means no limit)
	for i := 0; i < 5; i++ {
		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}
}

func TestWrapListenerNoLimit(t *testing.T) {
	engine := New()
	assert.Equal(t, int64(0), engine.MaxConns)

	// Create a dummy listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	// Should return original listener when no limit
	wrapped := engine.wrapListener(ln)
	assert.Equal(t, ln, wrapped)
}

func TestWrapListenerWithLimit(t *testing.T) {
	engine := New(WithMaxConns(5))

	// Create a dummy listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	// Should return wrapped listener
	wrapped := engine.wrapListener(ln)
	assert.NotNil(t, wrapped)
	_, ok := wrapped.(*limitedListener)
	assert.True(t, ok)
}
