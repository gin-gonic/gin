// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net"
	"sync"
)

// limitedListener wraps a net.Listener and limits the number of concurrent connections
// using a buffered channel as a semaphore.
type limitedListener struct {
	net.Listener
	sem chan struct{}
}

// Accept accepts a new connection. If the connection limit has been reached,
// the new connection is immediately closed.
func (l *limitedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	select {
	case l.sem <- struct{}{}:
		return &limitedConn{Conn: conn, sem: l.sem}, nil
	default:
		conn.Close()
		return nil, nil
	}
}

// limitedConn wraps a net.Conn and releases the semaphore slot on Close.
type limitedConn struct {
	net.Conn
	sem       chan struct{}
	closeOnce sync.Once
}

// Close closes the connection and releases the semaphore slot.
func (c *limitedConn) Close() error {
	c.closeOnce.Do(func() { <-c.sem })
	return c.Conn.Close()
}
