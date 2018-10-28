// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "github.com/gin-gonic/gin/internal"

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) BindUri(p internal.Params, obj interface{}) error {
	if err := mapUri(obj, p); err != nil {
		return err
	}
	return validate(obj)
}
