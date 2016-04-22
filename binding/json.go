// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"errors"

	"net/http"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req.Body == nil {
		return errors.New("can not bind a request with a nil body")
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
