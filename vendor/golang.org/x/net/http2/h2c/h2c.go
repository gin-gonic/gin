// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package h2c implements the unencrypted "h2c" form of HTTP/2.
//
// The h2c protocol is the non-TLS version of HTTP/2 which is not available from
// net/http or golang.org/x/net/http2.
package h2c

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2"
)

var (
	http2VerboseLogs bool
)

func init() {
	e := os.Getenv("GODEBUG")
	if strings.Contains(e, "http2debug=1") || strings.Contains(e, "http2debug=2") {
		http2VerboseLogs = true
	}
}

// h2cHandler is a Handler which implements h2c by hijacking the HTTP/1 traffic
// that should be h2c traffic. There are two ways to begin a h2c connection
// (RFC 7540 Section 3.2 and 3.4): (1) Starting with Prior Knowledge - this
// works by starting an h2c connection with a string of bytes that is valid
// HTTP/1, but unlikely to occur in practice and (2) Upgrading from HTTP/1 to
// h2c - this works by using the HTTP/1 Upgrade header to request an upgrade to
// h2c. When either of those situations occur we hijack the HTTP/1 connection,
// convert it to an HTTP/2 connection and pass the net.Conn to http2.ServeConn.
type h2cHandler struct {
	Handler http.Handler
	s       *http2.Server
}

// NewHandler returns an http.Handler that wraps h, intercepting any h2c
// traffic. If a request is an h2c connection, it's hijacked and redirected to
// s.ServeConn. Otherwise the returned Handler just forwards requests to h. This
// works because h2c is designed to be parseable as valid HTTP/1, but ignored by
// any HTTP server that does not handle h2c. Therefore we leverage the HTTP/1
// compatible parts of the Go http library to parse and recognize h2c requests.
// Once a request is recognized as h2c, we hijack the connection and convert it
// to an HTTP/2 connection which is understandable to s.ServeConn. (s.ServeConn
// understands HTTP/2 except for the h2c part of it.)
//
// The first request on an h2c connection is read entirely into memory before
// the Handler is called. To limit the memory consumed by this request, wrap
// the result of NewHandler in an http.MaxBytesHandler.
func NewHandler(h http.Handler, s *http2.Server) http.Handler {
	return &h2cHandler{
		Handler: h,
		s:       s,
	}
}

// extractServer extracts existing http.Server instance from http.Request or create an empty http.Server
func extractServer(r *http.Request) *http.Server {
	server, ok := r.Context().Value(http.ServerContextKey).(*http.Server)
	if ok {
		return server
	}
	return new(http.Server)
}

// ServeHTTP implement the h2c support that is enabled by h2c.GetH2CHandler.
func (s h2cHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle h2c with prior knowledge (RFC 7540 Section 3.4)
	if r.Method == "PRI" && len(r.Header) == 0 && r.URL.Path == "*" && r.Proto == "HTTP/2.0" {
		if http2VerboseLogs {
			log.Print("h2c: attempting h2c with prior knowledge.")
		}
		conn, err := initH2CWithPriorKnowledge(w)
		if err != nil {
			if http2VerboseLogs {
				log.Printf("h2c: error h2c with prior knowledge: %v", err)
			}
			return
		}
		defer conn.Close()
		s.s.ServeConn(conn, &http2.ServeConnOpts{
			Context:          r.Context(),
			BaseConfig:       extractServer(r),
			Handler:          s.Handler,
			SawClientPreface: true,
		})
		return
	}
	// Handle Upgrade to h2c (RFC 7540 Section 3.2)
	if isH2CUpgrade(r.Header) {
		conn, settings, err := h2cUpgrade(w, r)
		if err != nil {
			if http2VerboseLogs {
				log.Printf("h2c: error h2c upgrade: %v", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		s.s.ServeConn(conn, &http2.ServeConnOpts{
			Context:        r.Context(),
			BaseConfig:     extractServer(r),
			Handler:        s.Handler,
			UpgradeRequest: r,
			Settings:       settings,
		})
		return
	}
	s.Handler.ServeHTTP(w, r)
	return
}

// initH2CWithPriorKnowledge implements creating a h2c connection with prior
// knowledge (Section 3.4) and creates a net.Conn suitable for http2.ServeConn.
// All we have to do is look for the client preface that is suppose to be part
// of the body, and reforward the client preface on the net.Conn this function
// creates.
func initH2CWithPriorKnowledge(w http.ResponseWriter) (net.Conn, error) {
	rc := http.NewResponseController(w)
	conn, rw, err := rc.Hijack()
	if err != nil {
		return nil, err
	}

	const expectedBody = "SM\r\n\r\n"

	buf := make([]byte, len(expectedBody))
	n, err := io.ReadFull(rw, buf)
	if err != nil {
		return nil, fmt.Errorf("h2c: error reading client preface: %s", err)
	}

	if string(buf[:n]) == expectedBody {
		return newBufConn(conn, rw), nil
	}

	conn.Close()
	return nil, errors.New("h2c: invalid client preface")
}

// h2cUpgrade establishes a h2c connection using the HTTP/1 upgrade (Section 3.2).
func h2cUpgrade(w http.ResponseWriter, r *http.Request) (_ net.Conn, settings []byte, err error) {
	settings, err = getH2Settings(r.Header)
	if err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	rc := http.NewResponseController(w)
	conn, rw, err := rc.Hijack()
	if err != nil {
		return nil, nil, err
	}

	rw.Write([]byte("HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Upgrade: h2c\r\n\r\n"))
	return newBufConn(conn, rw), settings, nil
}

// isH2CUpgrade returns true if the header properly request an upgrade to h2c
// as specified by Section 3.2.
func isH2CUpgrade(h http.Header) bool {
	return httpguts.HeaderValuesContainsToken(h[textproto.CanonicalMIMEHeaderKey("Upgrade")], "h2c") &&
		httpguts.HeaderValuesContainsToken(h[textproto.CanonicalMIMEHeaderKey("Connection")], "HTTP2-Settings")
}

// getH2Settings returns the settings in the HTTP2-Settings header.
func getH2Settings(h http.Header) ([]byte, error) {
	vals, ok := h[textproto.CanonicalMIMEHeaderKey("HTTP2-Settings")]
	if !ok {
		return nil, errors.New("missing HTTP2-Settings header")
	}
	if len(vals) != 1 {
		return nil, fmt.Errorf("expected 1 HTTP2-Settings. Got: %v", vals)
	}
	settings, err := base64.RawURLEncoding.DecodeString(vals[0])
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func newBufConn(conn net.Conn, rw *bufio.ReadWriter) net.Conn {
	rw.Flush()
	if rw.Reader.Buffered() == 0 {
		// If there's no buffered data to be read,
		// we can just discard the bufio.ReadWriter.
		return conn
	}
	return &bufConn{conn, rw.Reader}
}

// bufConn wraps a net.Conn, but reads drain the bufio.Reader first.
type bufConn struct {
	net.Conn
	*bufio.Reader
}

func (c *bufConn) Read(p []byte) (int, error) {
	if c.Reader == nil {
		return c.Conn.Read(p)
	}
	n := c.Reader.Buffered()
	if n == 0 {
		c.Reader = nil
		return c.Conn.Read(p)
	}
	if n < len(p) {
		p = p[:n]
	}
	return c.Reader.Read(p)
}
