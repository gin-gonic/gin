// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package yaml

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding/common"
	"gopkg.in/yaml.v2"
)

func init() {
	common.List[common.MIMEYAML] = yamlBinding{}
}

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (yamlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeYAML(req.Body, obj)
}

func (yamlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeYAML(bytes.NewReader(body), obj)
}

func decodeYAML(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return common.Validate(obj)
}
