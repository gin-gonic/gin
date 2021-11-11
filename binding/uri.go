// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "context"

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (b uriBinding) BindUri(m map[string][]string, obj interface{}) error {
	return b.BindUriContext(context.Background(), m, obj)
}

func (uriBinding) BindUriContext(ctx context.Context, m map[string][]string, obj interface{}) error {
	if err := mapURI(obj, m); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}
