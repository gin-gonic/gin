// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value
		}
	}
	return ""
}

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
