// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type formBinding struct{}

func (_ formBinding) Name() string {
	return "query"
}

func (_ formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return Validate(obj)
}
