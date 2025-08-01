// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO
// func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// func (w *responseWriter) CloseNotify() <-chan bool {
// func (w *responseWriter) Flush() {

var (
	_ ResponseWriter      = &responseWriter{}
	_ http.ResponseWriter = &responseWriter{}
	_ http.ResponseWriter = ResponseWriter(&responseWriter{})
	_ http.Hijacker       = ResponseWriter(&responseWriter{})
	_ http.Flusher        = ResponseWriter(&responseWriter{})
	_ http.CloseNotifier  = ResponseWriter(&responseWriter{})
)

func init() {
	SetMode(TestMode)
}

func TestResponseWriterUnwrap(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: testWriter}
	assert.Same(t, testWriter, writer.Unwrap())
}

func TestResponseWriterReset(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	var w ResponseWriter = writer

	writer.reset(testWriter)
	assert.Equal(t, -1, writer.size)
	assert.Equal(t, http.StatusOK, writer.status)
	assert.Equal(t, testWriter, writer.ResponseWriter)
	assert.Equal(t, -1, w.Size())
	assert.Equal(t, http.StatusOK, w.Status())
	assert.False(t, w.Written())
}

func TestResponseWriterWriteHeader(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWriter)
	w := ResponseWriter(writer)

	w.WriteHeader(http.StatusMultipleChoices)
	assert.False(t, w.Written())
	assert.Equal(t, http.StatusMultipleChoices, w.Status())
	assert.NotEqual(t, http.StatusMultipleChoices, testWriter.Code)

	w.WriteHeader(-1)
	assert.Equal(t, http.StatusMultipleChoices, w.Status())
}

func TestResponseWriterWriteHeadersNow(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWriter)
	w := ResponseWriter(writer)

	w.WriteHeader(http.StatusMultipleChoices)
	w.WriteHeaderNow()

	assert.True(t, w.Written())
	assert.Equal(t, 0, w.Size())
	assert.Equal(t, http.StatusMultipleChoices, testWriter.Code)

	writer.size = 10
	w.WriteHeaderNow()
	assert.Equal(t, 10, w.Size())
}

func TestResponseWriterWrite(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWriter)
	w := ResponseWriter(writer)

	n, err := w.Write([]byte("hola"))
	assert.Equal(t, 4, n)
	assert.Equal(t, 4, w.Size())
	assert.Equal(t, http.StatusOK, w.Status())
	assert.Equal(t, http.StatusOK, testWriter.Code)
	assert.Equal(t, "hola", testWriter.Body.String())
	require.NoError(t, err)

	n, err = w.Write([]byte(" adios"))
	assert.Equal(t, 6, n)
	assert.Equal(t, 10, w.Size())
	assert.Equal(t, "hola adios", testWriter.Body.String())
	require.NoError(t, err)
}

func TestResponseWriterHijack(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWriter)
	w := ResponseWriter(writer)

	assert.Panics(t, func() {
		_, _, err := w.Hijack()
		require.NoError(t, err)
	})
	assert.True(t, w.Written())

	assert.Panics(t, func() {
		w.CloseNotify()
	})

	w.Flush()
}

type mockHijacker struct {
	*httptest.ResponseRecorder
	hijacked bool
}

// Hijack implements the http.Hijacker interface. It just records that it was called.
func (m *mockHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	m.hijacked = true
	return nil, nil, nil
}

func TestResponseWriterHijackAfterWrite(t *testing.T) {
	tests := []struct {
		name                      string
		action                    func(w ResponseWriter) error // Action to perform before hijacking
		expectWrittenBeforeHijack bool
		expectHijackSuccess       bool
		expectWrittenAfterHijack  bool
		expectError               error
	}{
		{
			name:                      "hijack before write should succeed",
			action:                    func(w ResponseWriter) error { return nil },
			expectWrittenBeforeHijack: false,
			expectHijackSuccess:       true,
			expectWrittenAfterHijack:  true, // Hijack itself marks the writer as written
			expectError:               nil,
		},
		{
			name: "hijack after write should fail",
			action: func(w ResponseWriter) error {
				_, err := w.Write([]byte("test"))
				return err
			},
			expectWrittenBeforeHijack: true,
			expectHijackSuccess:       false,
			expectWrittenAfterHijack:  true,
			expectError:               errHijackAlreadyWritten,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hijacker := &mockHijacker{ResponseRecorder: httptest.NewRecorder()}
			writer := &responseWriter{}
			writer.reset(hijacker)
			w := ResponseWriter(writer)

			// Check initial state
			assert.False(t, w.Written(), "should not be written initially")

			// Perform pre-hijack action
			require.NoError(t, tc.action(w), "unexpected error during pre-hijack action")

			// Check state before hijacking
			assert.Equal(t, tc.expectWrittenBeforeHijack, w.Written(), "unexpected w.Written() state before hijack")

			// Attempt to hijack
			_, _, hijackErr := w.Hijack()

			// Check results
			require.ErrorIs(t, hijackErr, tc.expectError, "unexpected error from Hijack()")
			assert.Equal(t, tc.expectHijackSuccess, hijacker.hijacked, "unexpected hijacker.hijacked state")
			assert.Equal(t, tc.expectWrittenAfterHijack, w.Written(), "unexpected w.Written() state after hijack")
		})
	}
}

func TestResponseWriterFlush(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{}
		writer.reset(w)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Flush()
	}))
	defer testServer.Close()

	// should return 500
	resp, err := http.Get(testServer.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestResponseWriterStatusCode(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWriter)
	w := ResponseWriter(writer)

	w.WriteHeader(http.StatusOK)
	w.WriteHeaderNow()

	assert.Equal(t, http.StatusOK, w.Status())
	assert.True(t, w.Written())

	w.WriteHeader(http.StatusUnauthorized)

	// status must be 200 although we tried to change it
	assert.Equal(t, http.StatusOK, w.Status())
}

// mockPusherResponseWriter is an http.ResponseWriter that implements http.Pusher.
type mockPusherResponseWriter struct {
	http.ResponseWriter
}

func (m *mockPusherResponseWriter) Push(target string, opts *http.PushOptions) error {
	return nil
}

// nonPusherResponseWriter is an http.ResponseWriter that does not implement http.Pusher.
type nonPusherResponseWriter struct {
	http.ResponseWriter
}

func TestPusherWithPusher(t *testing.T) {
	rw := &mockPusherResponseWriter{}
	w := &responseWriter{ResponseWriter: rw}

	pusher := w.Pusher()
	assert.NotNil(t, pusher, "Expected pusher to be non-nil")
}

func TestPusherWithoutPusher(t *testing.T) {
	rw := &nonPusherResponseWriter{}
	w := &responseWriter{ResponseWriter: rw}

	pusher := w.Pusher()
	assert.Nil(t, pusher, "Expected pusher to be nil")
}
