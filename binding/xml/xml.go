// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding/common"
)

func init() {
	xml := xmlBinding{}
	common.List[common.MIMEXML] = xml
	common.List[common.MIMEXML2] = xml
}

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeXML(req.Body, obj)
}

func (xmlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body), obj)
}
func decodeXML(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return common.Validate(obj)
}
