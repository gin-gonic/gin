// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownConfig holds configuration for graceful shutdown.
type ShutdownConfig struct {
	// Timeout is the maximum duration to wait for active connections to finish.
	// Default: 10 seconds
	Timeout time.Duration

	// Signals are the OS signals that will trigger shutdown.
	// Default: SIGINT, SIGTERM
	Signals []os.Signal
}

// RunWithShutdown starts the HTTP server and handles graceful shutdown on SIGINT/SIGTERM.
// It blocks until the server is shut down.
// The timeout parameter specifies the maximum duration to wait for active connections to finish.
func (engine *Engine) RunWithShutdown(addr string, timeout time.Duration) error {
	return engine.RunWithShutdownConfig(addr, ShutdownConfig{
		Timeout: timeout,
		Signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	})
}

// RunWithShutdownConfig starts the HTTP server with custom shutdown configuration.
// It blocks until the server is shut down.
func (engine *Engine) RunWithShutdownConfig(addr string, config ShutdownConfig) error {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if len(config.Signals) == 0 {
		config.Signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	ctx, stop := signal.NotifyContext(context.Background(), config.Signals...)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := engine.Run(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	return engine.Shutdown(shutdownCtx)
}
