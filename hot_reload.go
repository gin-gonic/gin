//go:build !windows

// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const hotReloadListenerEnv = "GIN_LISTENER_FD"

// RunWithHotReload runs the engine and enables zero-downtime hot reload via
// SIGHUP. On SIGHUP, a child process inherits the listening socket and begins
// serving immediately while the parent drains in-flight requests (up to 30s)
// and exits. The child handles subsequent SIGHUPs the same way.
//
// Send SIGINT or SIGTERM for a clean shutdown without spawning a replacement.
//
// Note: hot reload re-executes the same binary. Rebuilding must be handled
// externally (e.g. with make or a file watcher) before sending SIGHUP.
func (engine *Engine) RunWithHotReload(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.")
	}
	engine.updateRouteTrees()

	if fdStr := os.Getenv(hotReloadListenerEnv); fdStr != "" {
		fd, parseErr := strconv.Atoi(fdStr)
		if parseErr != nil {
			return fmt.Errorf("gin: invalid %s=%q: %w", hotReloadListenerEnv, fdStr, parseErr)
		}
		return engine.runInherited(fd)
	}

	address := resolveAddress(addr)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	debugPrint("Listening and serving HTTP on %s (hot reload enabled — send SIGHUP to reload)\n", address)
	return engine.serveWithSignals(ln)
}

// runInherited is the child-process entry point: it reconstructs the listener
// from an fd inherited via ExtraFiles and hands off to serveWithSignals.
func (engine *Engine) runInherited(fd int) error {
	f := os.NewFile(uintptr(fd), "gin-listener")
	ln, err := net.FileListener(f)
	f.Close() // net.FileListener dups the fd; our copy is no longer needed
	if err != nil {
		return fmt.Errorf("gin: could not create listener from fd %d: %w", fd, err)
	}
	defer ln.Close()
	debugPrint("Listening and serving HTTP on inherited socket (hot reload enabled — send SIGHUP to reload)\n")
	return engine.serveWithSignals(ln)
}

// serveWithSignals starts the HTTP server on ln and blocks until a signal
// arrives. SIGHUP forks a child then drains and exits; SIGINT/SIGTERM drain
// and exit without spawning a replacement.
func (engine *Engine) serveWithSignals(ln net.Listener) error {
	srv := &http.Server{Handler: engine.Handler()}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go srv.Serve(ln) //nolint:errcheck

	for sig := range sigCh {
		switch sig {
		case syscall.SIGHUP:
			debugPrint("received SIGHUP — forking child for zero-downtime reload\n")
			if err := spawnChild(ln); err != nil {
				debugPrint("hot reload fork failed: %v\n", err)
				continue
			}
			// Give the child a moment to call Accept before we stop.
			time.Sleep(100 * time.Millisecond)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := srv.Shutdown(ctx)
			cancel()
			return err
		case syscall.SIGINT, syscall.SIGTERM:
			debugPrint("received %v — shutting down gracefully\n", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := srv.Shutdown(ctx)
			cancel()
			return err
		}
	}
	return nil
}

// spawnChild forks a new instance of the current binary, passing the listening
// socket as fd 3 via ExtraFiles and advertising it through GIN_LISTENER_FD.
func spawnChild(ln net.Listener) error {
	tcpLn, ok := ln.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("gin: hot reload requires a TCP listener, got %T", ln)
	}

	f, err := tcpLn.File()
	if err != nil {
		return fmt.Errorf("gin: could not duplicate listener fd: %w", err)
	}
	defer f.Close()

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("gin: could not resolve executable path: %w", err)
	}

	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=3", hotReloadListenerEnv))
	cmd.ExtraFiles = []*os.File{f} // ExtraFiles[0] becomes fd 3 in the child

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("gin: failed to start child process: %w", err)
	}
	go cmd.Wait() //nolint:errcheck — best-effort zombie reap before parent exits
	return nil
}
