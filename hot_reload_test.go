//go:build !windows

// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWithHotReload_InvalidFD(t *testing.T) {
	t.Setenv(hotReloadListenerEnv, "not-a-number")
	err := New().RunWithHotReload()
	assert.ErrorContains(t, err, "invalid")
}

// TestRunWithHotReload_GracefulShutdown starts the engine via RunWithHotReload
// and verifies it shuts down cleanly on SIGTERM. signal.Notify inside
// serveWithSignals captures SIGTERM before the default handler fires, so the
// test process is not terminated.
func TestRunWithHotReload_GracefulShutdown(t *testing.T) {
	// Reserve a free port then release it; there is a small TOCTOU window.
	ln0, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := ln0.Addr().String()
	ln0.Close()

	engine := New()
	engine.GET("/ping", func(c *Context) { c.String(200, "pong") })

	errCh := make(chan error, 1)
	go func() { errCh <- engine.RunWithHotReload(addr) }()

	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 5*time.Second, 10*time.Millisecond, "server never became reachable")

	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, proc.Signal(syscall.SIGTERM))

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("server did not shut down within 10s")
	}
}

// TestRunWithHotReload_InheritedListener exercises the child-process path by
// pre-opening a TCP socket, duplicating its fd, and advertising it via the
// environment variable that runInherited reads.
func TestRunWithHotReload_InheritedListener(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	addr := ln.Addr().String()

	// tcpLn.File() dups the underlying fd; we then dup again so that
	// runInherited's f.Close() doesn't affect our reference.
	tcpLn := ln.(*net.TCPListener)
	f, err := tcpLn.File()
	require.NoError(t, err)
	defer f.Close()

	dupFD, err := syscall.Dup(int(f.Fd()))
	require.NoError(t, err)
	// dupFD is now owned by RunWithHotReload; do not close it here.

	t.Setenv(hotReloadListenerEnv, fmt.Sprintf("%d", dupFD))

	engine := New()
	engine.GET("/ping", func(c *Context) { c.String(200, "pong") })

	errCh := make(chan error, 1)
	go func() { errCh <- engine.RunWithHotReload() }()

	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 5*time.Second, 10*time.Millisecond, "inherited server never became reachable")

	proc, _ := os.FindProcess(os.Getpid())
	proc.Signal(syscall.SIGTERM)

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("inherited server did not shut down within 10s")
	}
}
