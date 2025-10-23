// Copyright 2025 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

package binding

import "net/http"

// Content-Type MIME of msgpack.
const (
	MIMEMSGPACK  = "application/x-msgpack"
	MIMEMSGPACK2 = "application/msgpack"
)

// MsgPack implement the BindingBody interface.
var MsgPack BindingBody = msgpackBinding{}

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack

	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEYAML, MIMEYAML2:
		return YAML
	case MIMETOML:
		return TOML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}
