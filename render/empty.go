package render

import (
	"fmt"
	"net/http"
)

type EmptyRenderFactory struct{}

type EmptyRender struct{}

func init() {
	Register(EmptyRenderType, EmptyRenderFactory{})
}

// Instance apply opts to build a new EmptyRender instance
func (EmptyRenderFactory) Instance() RenderRecycler {
	return &EmptyRender{}
}

// Render writes data with custom ContentType.
func (*EmptyRender) Render(http.ResponseWriter) error {
	return fmt.Errorf("empty render,you need register one first")
}

// WriteContentType writes custom ContentType.
func (*EmptyRender) WriteContentType(w http.ResponseWriter) {
	// Empty
}

// Setup set data and opts
func (*EmptyRender) Setup(data interface{}, opts ...interface{}) {
	// Empty
}

// Reset clean data and opts
func (*EmptyRender) Reset() {
	// Empty
}
