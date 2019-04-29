// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"

	"github.com/gin-gonic/gin/render/common"
)

var (
	_ common.Render = String{}
	_ common.Render = Redirect{}
	_ common.Render = Data{}
	_ common.Render = HTML{}
	_ HTMLRender    = HTMLDebug{}
	_ HTMLRender    = HTMLProduction{}
	_ common.Render = Reader{}
)

//YAML return the render for yaml if loaded
func YAML(obj interface{}) common.Render {
	return retRender("YAML", obj, nil)
}

//XML return the render for xml if loaded
func XML(obj interface{}) common.Render {
	return retRender("XML", obj, nil)
}

//ProtoBuf return the render for ProtoBuf if loaded
func ProtoBuf(obj interface{}) common.Render {
	return retRender("ProtoBuf", obj, nil)
}

//MsgPack return the render for MsgPack if loaded
func MsgPack(obj interface{}) common.Render {
	return retRender("MsgPack", obj, nil)
}

//JSON return the render for JSON if loaded
func JSON(obj interface{}) common.Render {
	return retRender("JSON", obj, nil)
}

//IndentedJSON return the render for IndentedJSON if loaded
func IndentedJSON(obj interface{}) common.Render {
	return retRender("IndentedJSON", obj, nil)
}

//SecureJSON return the render for SecureJSON if loaded
func SecureJSON(prefix string, obj interface{}) common.Render {
	return retRender("SecureJSON", obj, map[string]string{
		"Prefix": prefix,
	})
}

//JsonpJSON return the render for JsonpJSON if loaded
func JsonpJSON(callback string, obj interface{}) common.Render {
	return retRender("JsonpJSON", obj, map[string]string{
		"Callback": callback,
	})
}

//AsciiJSON return the render for AsciiJSON if loaded
func AsciiJSON(obj interface{}) common.Render {
	return retRender("AsciiJSON", obj, nil)
}

//Search for a render
func retRender(rID string, obj interface{}, opts map[string]string) common.Render {
	r, ok := common.List[rID]
	if !ok {
		panic(fmt.Sprintf("Undefined render %s", rID))
	}
	return r(obj, opts)
}
