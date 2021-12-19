// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"net/http"
)

const defaultMemory = 32 << 20

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return validate(obj)
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	if err := mappingByPtr(obj, (*multipartRequest)(req), "form"); err != nil {
		return err
	}

	return validate(obj)
}
