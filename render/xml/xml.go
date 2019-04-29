// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin/render/common"
)

func init() {
	common.List["XML"] = NewXML
}

// XML contains the given interface object.
type XML struct {
	Data interface{}
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render (XML) encodes the given interface object and writes data with custom ContentType.
func (r XML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

// WriteContentType (XML) writes XML ContentType for response.
func (r XML) WriteContentType(w http.ResponseWriter) {
	common.WriteContentType(w, xmlContentType)
}

//NewXML build a new xml render
func NewXML(obj interface{}, _ map[string]string) common.Render {
	return XML{Data: obj}
}
