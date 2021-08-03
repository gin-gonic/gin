// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"
)

const defaultMemory = 32 << 20

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (b formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := b.BindOnly(req, obj); err != nil {
		return err
	}
	return validate(obj)
}

func (b formBinding) BindOnly(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && err != http.ErrNotMultipart {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return nil
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (b formPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := b.BindOnly(req, obj); err != nil {
		return err
	}
	return validate(obj)
}

func (b formPostBinding) BindOnly(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return nil
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (b formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := b.BindOnly(req, obj); err != nil {
		return err
	}
	return validate(obj)
}

func (b formMultipartBinding) BindOnly(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	if err := mappingByPtr(obj, (*multipartRequest)(req), "form"); err != nil {
		return err
	}

	return nil
}
