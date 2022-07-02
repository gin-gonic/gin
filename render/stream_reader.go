// Copyright 2018 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"
)

// StreamReader contains the IO reader and its length, and custom ContentType and other headers.
type StreamReader struct {
	Reader
}

type writerFlusher struct {
	http.ResponseWriter
}

func (w *writerFlusher) Write(buf []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(buf)
	if err == nil {
		w.ResponseWriter.(http.Flusher).Flush()
	}
	return
}

// Render (StreamReader) writes data with custom ContentType and headers.
func (r StreamReader) Render(w http.ResponseWriter) (err error) {
	return r.Reader.Render(&writerFlusher{w})
}
