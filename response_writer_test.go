// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO
// func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// func (w *responseWriter) CloseNotify() <-chan bool {
// func (w *responseWriter) Flush() {

var _ ResponseWriter = &responseWriter{}
var _ http.ResponseWriter = &responseWriter{}
var _ http.ResponseWriter = ResponseWriter(&responseWriter{})
var _ http.Hijacker = ResponseWriter(&responseWriter{})
var _ http.Flusher = ResponseWriter(&responseWriter{})
var _ http.CloseNotifier = ResponseWriter(&responseWriter{})

func init() {
	SetMode(TestMode)
}

func TestResponseWriterReset(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	var w ResponseWriter = writer

	writer.reset(testWritter)
	assert.Equal(t, -1, writer.size)
	assert.Equal(t, 200, writer.status)
	assert.Equal(t, testWritter, writer.ResponseWriter)
	assert.Equal(t, -1, w.Size())
	assert.Equal(t, 200, w.Status())
	assert.False(t, w.Written())
}

func TestResponseWriterWriteHeader(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWritter)
	w := ResponseWriter(writer)

	w.WriteHeader(300)
	assert.False(t, w.Written())
	assert.Equal(t, 300, w.Status())
	assert.NotEqual(t, testWritter.Code, 300)

	w.WriteHeader(-1)
	assert.Equal(t, 300, w.Status())
}

func TestResponseWriterWriteHeadersNow(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWritter)
	w := ResponseWriter(writer)

	w.WriteHeader(300)
	w.WriteHeaderNow()

	assert.True(t, w.Written())
	assert.Equal(t, 0, w.Size())
	assert.Equal(t, 300, testWritter.Code)

	writer.size = 10
	w.WriteHeaderNow()
	assert.Equal(t, 10, w.Size())
}

func TestResponseWriterWrite(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWritter)
	w := ResponseWriter(writer)

	n, err := w.Write([]byte("hola"))
	assert.Equal(t, 4, n)
	assert.Equal(t, 4, w.Size())
	assert.Equal(t, 200, w.Status())
	assert.Equal(t, 200, testWritter.Code)
	assert.Equal(t, "hola", testWritter.Body.String())
	assert.NoError(t, err)

	n, err = w.Write([]byte(" adios"))
	assert.Equal(t, 6, n)
	assert.Equal(t, 10, w.Size())
	assert.Equal(t, "hola adios", testWritter.Body.String())
	assert.NoError(t, err)
}

func TestResponseWriterHijack(t *testing.T) {
	testWritter := httptest.NewRecorder()
	writer := &responseWriter{}
	writer.reset(testWritter)
	w := ResponseWriter(writer)

	assert.Panics(t, func() {
		w.Hijack()
	})
	assert.True(t, w.Written())

	assert.Panics(t, func() {
		w.CloseNotify()
	})

	w.Flush()
}
