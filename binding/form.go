// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"net/http"
)

const defaultMemory = 32 << 20

type (
	formBinding          struct{}
	formPostBinding      struct{}
	formMultipartBinding struct{}
)

func (formBinding) Name() string {
	return "form"
}

func (b formBinding) Bind(req *http.Request, obj any) error {
	return b.BindWithMemory(req, obj, defaultMemory)
}

// BindWithMemory binds form data using the given multipart memory threshold.
func (formBinding) BindWithMemory(req *http.Request, obj any, maxMemory int64) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(maxMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
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

func (b formMultipartBinding) Bind(req *http.Request, obj any) error {
	return b.BindWithMemory(req, obj, defaultMemory)
}

// BindWithMemory binds multipart form data using the given memory threshold.
func (formMultipartBinding) BindWithMemory(req *http.Request, obj any, maxMemory int64) error {
	if err := req.ParseMultipartForm(maxMemory); err != nil {
		return err
	}
	if err := mappingByPtr(obj, (*multipartRequest)(req), "form"); err != nil {
		return err
	}

	return validate(obj)
}
