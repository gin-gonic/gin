// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "github.com/gin-gonic/gin/binding/common"

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) BindUri(m map[string][]string, obj interface{}) error {
	if err := mapUri(obj, m); err != nil {
		return err
	}
	return common.Validate(obj)
}
