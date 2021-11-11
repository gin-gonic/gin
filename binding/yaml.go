// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (b yamlBinding) Bind(req *http.Request, obj interface{}) error {
	return b.BindContext(context.Background(), req, obj)
}

func (yamlBinding) BindContext(ctx context.Context, req *http.Request, obj interface{}) error {
	return decodeYAML(ctx, req.Body, obj)
}

func (b yamlBinding) BindBody(body []byte, obj interface{}) error {
	return b.BindBodyContext(context.Background(), body, obj)
}

func (yamlBinding) BindBodyContext(ctx context.Context, body []byte, obj interface{}) error {
	return decodeYAML(ctx, bytes.NewReader(body), obj)
}

func decodeYAML(ctx context.Context, r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}
