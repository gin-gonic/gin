// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"
)

type (
	CustomRenderFunc func(http.ResponseWriter, string, interface{}) error

	CustomRender interface {
		Instance(string, interface{}) Render
	}

	CustomRenderProduction struct {
		RenderFunc    CustomRenderFunc
	}

	Custom struct {
		RenderFunc CustomRenderFunc
		Name       string
		Data       interface{}
	}
)

func (r CustomRenderProduction) Instance(name string, data interface{}) Render {
	return Custom{
		RenderFunc: r.RenderFunc,
		Name: name,
		Data: data,
	}
}

func (r Custom) Render(w http.ResponseWriter) error {
	return r.RenderFunc(w, r.Name, r.Data)
}
