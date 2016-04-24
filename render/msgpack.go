// Copyright 2016 Andida Syahendar.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"net/http"
)

type Msgpack struct {
	Data interface{}
}

var msgpackContentType = []string{"application/x-msgpack; charset=utf-8"}

func (r Msgpack) Render(w http.ResponseWriter) error {
	writeContentType(w, msgpackContentType)

	return msgpack.NewEncoder(w).Encode(r.Data)
}
