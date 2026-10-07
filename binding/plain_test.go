// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// errReadCloser simulates a ReadCloser whose Read returns a fixed error.
type errReadCloser struct{ err error }

func (e *errReadCloser) Read(p []byte) (int, error) { return 0, e.err }
func (e *errReadCloser) Close() error               { return nil }

func TestDecodePlain_String_Success(t *testing.T) {
	t.Parallel()
	var s string
	if err := (plainBinding{}).BindBody([]byte("hello world"), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", s)
	}
}

func TestDecodePlain_ByteSlice_Success(t *testing.T) {
	t.Parallel()
	in := []byte{1, 2, 3, 4}
	var b []byte
	if err := (plainBinding{}).BindBody(in, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, in) {
		t.Fatalf("expected %v, got %v", in, b)
	}
}

func TestPlainBind_UsesHTTPRequestBody(t *testing.T) {
	t.Parallel()
	var s string
	req := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte("reqbody")))}
	if err := (plainBinding{}).Bind(req, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "reqbody" {
		t.Fatalf("expected %q, got %q", "reqbody", s)
	}
}

func TestDecodePlain_NilObj_NoPanic(t *testing.T) {
	// Passing nil obj should be a no-op and return nil error.
	if err := (plainBinding{}).BindBody([]byte("x"), nil); err != nil {
		t.Fatalf("expected nil error for nil obj, got %v", err)
	}

	// Passing a nil pointer (e.g., *string == nil) should also return nil error.
	var ps *string = nil
	if err := (plainBinding{}).BindBody([]byte("x"), ps); err != nil {
		t.Fatalf("expected nil error for nil pointer obj, got %v", err)
	}
}

func TestDecodePlain_UnsupportedType_Error(t *testing.T) {
	var x int
	err := (plainBinding{}).BindBody([]byte("x"), &x)
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Fatalf("expected error to contain 'unknown type', got %v", err)
	}
}

func TestPlainBind_ReadError(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("read fail")
	req := &http.Request{Body: &errReadCloser{err: sentinel}}
	var s string
	err := (plainBinding{}).Bind(req, &s)
	if err == nil {
		t.Fatalf("expected read error, got nil")
	}
	if err != sentinel {
		t.Fatalf("expected sentinel error %v, got %v", sentinel, err)
	}
}
