// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"path"
	"strings"
)

// Used internally to configure router, a RouterGroup is associated with a prefix
// and an array of handlers (middlewares)
type RouterGroup struct {
	Handlers HandlersChain
	BasePath string
	engine   *Engine
}

// Adds middlewares to the group, see example code in github.
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		BasePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

// Handle registers a new request handle and middlewares with the given path and method.
// The last handler should be the real handler, the other ones should be middlewares that can and should be shared among different routes.
// See the example code in github.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
}

func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	group.handle(httpMethod, relativePath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	group.handle("POST", relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	group.handle("GET", relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	group.handle("DELETE", relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) {
	group.handle("PATCH", relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	group.handle("PUT", relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	group.handle("OPTIONS", relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) {
	group.handle("HEAD", relativePath, handlers)
}

func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) {
	// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE
	group.handle("GET", relativePath, handlers)
	group.handle("POST", relativePath, handlers)
	group.handle("PUT", relativePath, handlers)
	group.handle("PATCH", relativePath, handlers)
	group.handle("HEAD", relativePath, handlers)
	group.handle("OPTIONS", relativePath, handlers)
	group.handle("DELETE", relativePath, handlers)
	group.handle("CONNECT", relativePath, handlers)
	group.handle("TRACE", relativePath, handlers)
}

func (group *RouterGroup) StaticFile(relativePath, filepath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *Context) {
		c.File(filepath)
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (group *RouterGroup) Static(relativePath, root string) {
	group.StaticFS(relativePath, http.Dir(root), false)
}

func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem, listDirectory bool) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs, listDirectory)
	relativePath = path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem, listDirectory bool) HandlerFunc {
	absolutePath := group.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		if !listDirectory && lastChar(c.Request.URL.Path) == '/' {
			http.NotFound(c.Writer, c.Request)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.BasePath, relativePath)
}
