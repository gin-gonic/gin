// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

// Factory builds a Render for the given data. Optional format subpackages
// (github.com/gin-gonic/gin/render/<format>) register one per content type so
// that Context.Negotiate can produce their Render without the core importing
// the underlying codec library.
type Factory func(data any) Render

// registry maps a content type to its Factory. It is only written from init()
// functions before main runs, so it needs no synchronization.
var registry = map[string]Factory{}

// Register associates a Factory with a content type. It is intended to be
// called from an init() function in a format subpackage.
func Register(contentType string, factory Factory) {
	registry[contentType] = factory
}

// Negotiate returns a Render for the content type and data when a Factory has
// been registered for it (i.e. the matching format subpackage was imported).
// The boolean reports whether a Factory was found.
func Negotiate(contentType string, data any) (Render, bool) {
	factory, ok := registry[contentType]
	if !ok {
		return nil, false
	}
	return factory(data), true
}
