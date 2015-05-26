// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

const MAX_MEMORY = 1 * 1024 * 1024

type multipartFormBinding struct{}

func (_ multipartFormBinding) Name() string {
	return "multipart form"
}

func (_ multipartFormBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(MAX_MEMORY); err != nil {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return Validate(obj)
}
