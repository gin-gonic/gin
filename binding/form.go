// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"
	"strings"
)

const defaultMemory = 32 * 1024 * 1024

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(defaultMemory)
	if err := mapForm(obj, resetForm(req.Form)); err != nil {
		return err
	}
	return validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(req *http.Request, obj interface{}) error {
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

func (formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	if err := mapForm(obj, req.MultipartForm.Value); err != nil {
		return err
	}
	return validate(obj)
}

func resetForm(f map[string][]string) map[string][]string {
	dicts := make(map[string][]string)
	for k, v := range f {
		lists := strings.Split(k, "&")
		if len(lists) == 1 {
			dicts[k] = v
			continue
		}
		for _, vv := range lists {
			p := strings.Split(vv, "=")
			dicts[p[0]] = append(dicts[p[0]], p[1])
		}
	}
	return dicts
}
