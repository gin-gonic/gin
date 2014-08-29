// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"log"
	"net/http"
)

type (
	ResponseWriter interface {
		http.ResponseWriter
		Status() int
		Written() bool
		WriteHeaderNow()
	}

	responseWriter struct {
		http.ResponseWriter
		status  int
		written bool
	}
)

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.status = 200
	w.written = false
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 {
		w.status = code
		if w.written {
			log.Println("[GIN] WARNING. Headers were already written!")
		}
	}
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.written {
		w.written = true
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Written() bool {
	return w.written
}
