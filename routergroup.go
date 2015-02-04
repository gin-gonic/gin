// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
)

// Used internally to configure router, a RouterGroup is associated with a prefix
// and an array of handlers (middlewares)
type RouterGroup struct {
	Handlers     []HandlerFunc
	absolutePath string
	engine       *Engine
}

// Adds middlewares to the group, see example code in github.
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers:     group.combineHandlers(handlers),
		absolutePath: group.calculateAbsolutePath(relativePath),
		engine:       group.engine,
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
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers []HandlerFunc) {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers[nuHandlers-1])
		debugPrint("%-5s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	group.engine.router.Handle(httpMethod, absolutePath, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		context := group.engine.createContext(w, req, params, handlers)
		context.Next()
		context.Writer.WriteHeaderNow()
		group.engine.reuseContext(context)
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	group.Handle("POST", relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	group.Handle("GET", relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	group.Handle("DELETE", relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) {
	group.Handle("PATCH", relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	group.Handle("PUT", relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	group.Handle("OPTIONS", relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) {
	group.Handle("HEAD", relativePath, handlers)
}

// LINK is a shortcut for router.Handle("LINK", path, handle)
func (group *RouterGroup) LINK(relativePath string, handlers ...HandlerFunc) {
	group.Handle("LINK", relativePath, handlers)
}

// UNLINK is a shortcut for router.Handle("UNLINK", path, handle)
func (group *RouterGroup) UNLINK(relativePath string, handlers ...HandlerFunc) {
	group.Handle("UNLINK", relativePath, handlers)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (group *RouterGroup) Static(relativePath, root string) {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handler := group.createStaticHandler(absolutePath, root)
	absolutePath = path.Join(absolutePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(absolutePath, handler)
	group.HEAD(absolutePath, handler)
}

func (group *RouterGroup) createStaticHandler(absolutePath, root string) func(*Context) {
	fileServer := http.StripPrefix(absolutePath, http.FileServer(http.Dir(root)))
	return func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, group.Handlers...)
	return append(mergedHandlers, handlers...)
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return group.absolutePath
	}
	absolutePath := path.Join(group.absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(absolutePath) != '/'
	if appendSlash {
		return absolutePath + "/"
	}
	return absolutePath
}
