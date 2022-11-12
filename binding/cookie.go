// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type cookieBinding struct{}

func (cookieBinding) Name() string {
	return "cookie"
}

func (c cookieBinding) Bind(req *http.Request, obj any) error {
	cookies := make(map[string][]string, len(req.Cookies()))
	for _, cookie := range req.Cookies() {
		cookies[cookie.Name] = append(cookies[cookie.Name], cookie.Value)
	}

	return mapFormByTag(obj, cookies, c.Name())
}
