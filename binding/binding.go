// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"fmt"

	"github.com/gin-gonic/gin/binding/common"
)

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	Uri           = uriBinding{}
	FormMultipart = formMultipartBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) common.Binding {
	if method == "GET" {
		return Form
	}
	switch contentType {
	case common.MIMEMultipartPOSTForm:
		return FormMultipart
	default:
		b, ok := common.List[contentType]
		if !ok {
			return Form //Default to Form
		}
		return b
	}
}

//YAML return the binding for yaml if loaded
func YAML() common.BindingBody {
	return retBinding(common.MIMEYAML)
}

//JSON return the binding for json if loaded
func JSON() common.BindingBody {
	return retBinding(common.MIMEJSON)
}

//XML return the binding for xml if loaded
func XML() common.BindingBody {
	return retBinding(common.MIMEXML)
}

//ProtoBuf return the binding for ProtoBuf if loaded
func ProtoBuf() common.BindingBody {
	return retBinding(common.MIMEPROTOBUF)
}

//MsgPack return the binding for MsgPack if loaded
func MsgPack() common.BindingBody {
	return retBinding(common.MIMEMSGPACK)
}

//retBinding Search for a render and panic on not found
func retBinding(contentType string) common.BindingBody {
	b, ok := common.List[contentType]
	if !ok {
		panic(fmt.Sprintf("Undefined binding %s", contentType))
	}
	return b
}
