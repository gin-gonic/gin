// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type postFormBinding struct{}

func (_ postFormBinding) Name() string {
	return "post_form"
}

func (_ postFormBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	if err := _validator.ValidateStruct(obj); err != nil {
		return error(err)
	}
	return nil
}
