// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (b yamlBinding) Bind(req *http.Request, obj interface{}) error {
	if err := b.BindOnly(req, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (yamlBinding) BindOnly(req *http.Request, obj interface{}) error {
	return decodeYAML(req.Body, obj)
}

func (b yamlBinding) BindBody(body []byte, obj interface{}) error {
	if err := b.BindBodyOnly(body, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (yamlBinding) BindBodyOnly(body []byte, obj interface{}) error {
	return decodeYAML(bytes.NewReader(body), obj)
}

func decodeYAML(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(obj)
}
