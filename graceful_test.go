// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineShutdown(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	// Start server in goroutine
	go func() {
		err := router.Run(":18080")
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()
	time.Sleep(100 * time.Millisecond) // Wait for server start

	// Verify server is running
	resp, err := http.Get("http://localhost:18080/")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = router.Shutdown(ctx)
	require.NoError(t, err)

	// Wait a moment for server to fully stop
	time.Sleep(50 * time.Millisecond)

	// Verify server is stopped
	_, err = http.Get("http://localhost:18080/")
	require.Error(t, err)
}

func TestEngineShutdownBeforeStart(t *testing.T) {
	router := New()

	// Shutdown before starting should not error
	err := router.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestEngineShutdownTLS(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	// Start TLS server in goroutine
	go func() {
		err := router.RunTLS(":18443", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem")
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()
	time.Sleep(100 * time.Millisecond) // Wait for server start

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := router.Shutdown(ctx)
	require.NoError(t, err)
}

func TestEngineShutdownWithActiveRequest(t *testing.T) {
	router := New()

	requestStarted := make(chan struct{})
	requestDone := make(chan struct{})

	router.GET("/slow", func(c *Context) {
		close(requestStarted)
		time.Sleep(500 * time.Millisecond) // Simulate slow request
		c.String(http.StatusOK, "done")
		close(requestDone)
	})

	// Start server
	go func() {
		_ = router.Run(":18081")
	}()
	time.Sleep(100 * time.Millisecond)

	// Start slow request
	go func() {
		resp, err := http.Get("http://localhost:18081/slow")
		if err == nil {
			resp.Body.Close()
		}
	}()

	// Wait for request to start
	<-requestStarted

	// Initiate shutdown while request is in progress
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- router.Shutdown(ctx)
	}()

	// Verify request completes before shutdown finishes
	select {
	case <-requestDone:
		// Request completed - this is expected
	case err := <-shutdownDone:
		t.Errorf("Shutdown completed before request finished: %v", err)
	}

	// Wait for shutdown to complete
	err := <-shutdownDone
	require.NoError(t, err)
}

func TestRunWithShutdown(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- router.RunWithShutdown(":18082", 5*time.Second)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:18082/")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Send shutdown signal to self
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = p.Signal(syscall.SIGINT)
	require.NoError(t, err)

	// Wait for shutdown to complete
	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Shutdown timed out")
	}
}

func TestRunWithShutdownConfig(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	config := ShutdownConfig{
		Timeout: 5 * time.Second,
		Signals: []os.Signal{syscall.SIGUSR1},
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- router.RunWithShutdownConfig(":18083", config)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:18083/")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Send custom signal
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = p.Signal(syscall.SIGUSR1)
	require.NoError(t, err)

	// Wait for shutdown to complete
	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Shutdown timed out")
	}
}

func TestRunWithShutdownConfigDefaults(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	// Test with zero values to check defaults are applied
	config := ShutdownConfig{}

	errCh := make(chan error, 1)
	go func() {
		errCh <- router.RunWithShutdownConfig(":18084", config)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:18084/")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGINT (default signal)
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = p.Signal(syscall.SIGINT)
	require.NoError(t, err)

	// Wait for shutdown to complete
	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(15 * time.Second):
		t.Fatal("Shutdown timed out")
	}
}

func TestRunWithShutdownServerError(t *testing.T) {
	router := New()

	// Start a server on the same port first
	listener, err := net.Listen("tcp", ":18085")
	require.NoError(t, err)
	defer listener.Close()

	// Try to run on the same port - should fail
	errCh := make(chan error, 1)
	go func() {
		errCh <- router.RunWithShutdown(":18085", 5*time.Second)
	}()

	// Should get an error because port is already in use
	select {
	case err := <-errCh:
		require.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Expected error but timed out")
	}
}
