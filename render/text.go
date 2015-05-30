// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"io"
	"net/http"
)

type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (r String) Write(w http.ResponseWriter) error {
	header := w.Header()
	if _, exist := header["Content-Type"]; !exist {
		header["Content-Type"] = plainContentType
	}
	if len(r.Data) > 0 {
		fmt.Fprintf(w, r.Format, r.Data...)
	} else {
		io.WriteString(w, r.Format)
	}
	return nil
}
