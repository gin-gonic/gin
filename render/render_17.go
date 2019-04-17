// +build go1.7

package render

import "github.com/gin-gonic/gin/render/common"

//PureJSON return the render for AsciiJSON if loaded
func PureJSON(obj interface{}) common.Render {
	return retRender("PureJSON", obj, nil)
}
