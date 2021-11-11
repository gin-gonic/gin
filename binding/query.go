// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"context"
	"net/http"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (b queryBinding) Bind(req *http.Request, obj interface{}) error {
	return b.BindContext(context.Background(), req, obj)
}

func (queryBinding) BindContext(ctx context.Context, req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}
