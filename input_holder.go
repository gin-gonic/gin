// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

type inputHolder struct {
	context *Context
}

func (i inputHolder) FromGET(key string) (va string) {
	va, _ = i.fromGET(key)
	return
}

func (i inputHolder) FromPOST(key string) (va string) {
	va, _ = i.fromPOST(key)
	return
}

func (i inputHolder) Get(key string) string {
	if value, exists := i.fromPOST(key); exists {
		return value
	}
	if value, exists := i.fromGET(key); exists {
		return value
	}
	return ""
}

func (i inputHolder) fromGET(key string) (string, bool) {
	req := i.context.Request
	req.ParseForm()
	if values, ok := req.Form[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (i inputHolder) fromPOST(key string) (string, bool) {
	req := i.context.Request
	req.ParseForm()
	if values, ok := req.PostForm[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
