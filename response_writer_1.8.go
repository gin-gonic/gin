// +build go1.8

// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
)

// ResponseWriter ...
type ResponseWriter interface {
	responseWriterBase
	// get the http.Pusher for server push
	Pusher() http.Pusher
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
