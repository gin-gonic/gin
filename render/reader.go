// Copyright 2018 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"io"
	"net/http"
	"strconv"
)

// Reader contains the IO reader and its length, and custom ContentType and other headers.
type Reader struct {
	ContentType   string
	ContentLength int64
	Reader        io.Reader
	Headers       map[string]string
}

// Render (Reader) writes data with custom ContentType and headers.
func (r Reader) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if r.ContentLength >= 0 {
		if r.Headers == nil {
			r.Headers = map[string]string{}
		}
		r.Headers["Content-Length"] = strconv.FormatInt(r.ContentLength, 10)
	}
	r.writeHeaders(w, r.Headers)
	_, err = io.Copy(w, r.Reader)
	return
}

// WriteContentType (Reader) writes custom ContentType.
func (r Reader) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{r.ContentType})
}

// writeHeaders writes custom Header.
func (r Reader) writeHeaders(w http.ResponseWriter, headers map[string]string) {
	header := w.Header()
	for k, v := range headers {
		if header.Get(k) == "" {
			header.Set(k, v)
		}
	}
}

type ReaderStream struct {
	ContentType string
	Reader      io.Reader
	Headers     map[string]string
}

// Render (ReaderStream) writes data with custom ContentType and headers.
func (r ReaderStream) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	r.writeHeaders(w, r.Headers)
	_, err = io.Copy(w, r.Reader)
	return
}

func (r ReaderStream) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{r.ContentType})
}

func (r ReaderStream) writeHeaders(w http.ResponseWriter, headers map[string]string) {
	header := w.Header()
	for k, v := range headers {
		if val := header[k]; len(val) == 0 {
			header[k] = []string{v}
		}
	}
}
