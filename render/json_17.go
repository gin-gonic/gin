// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package render

import (
	"net/http"

	"github.com/gin-gonic/gin/internal/json"
)

type PureJSON struct {
	Data interface{}
}

func (r PureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(r.Data)
}

func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
