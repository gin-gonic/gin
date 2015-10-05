// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginS

import (
	"html/template"
	"net/http"
	"sync"

	. "github.com/gin-gonic/gin"
)

var once sync.Once
var internalEngine *Engine

func engine() *Engine {
	once.Do(func() {
		internalEngine = Default()
	})
	return internalEngine
}

func LoadHTMLGlob(pattern string) {
	engine().LoadHTMLGlob(pattern)
}

func LoadHTMLFiles(files ...string) {
	engine().LoadHTMLFiles(files...)
}

func SetHTMLTemplate(templ *template.Template) {
	engine().SetHTMLTemplate(templ)
}

// Adds handlers for NoRoute. It return a 404 code by default.
func NoRoute(handlers ...HandlerFunc) {
	engine().NoRoute(handlers...)
}

// Sets the handlers called when... TODO
func NoMethod(handlers ...HandlerFunc) {
	engine().NoMethod(handlers...)
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return engine().Group(relativePath, handlers...)
}

func Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().Handle(httpMethod, relativePath, handlers...)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().POST(relativePath, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().GET(relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().DELETE(relativePath, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().PATCH(relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().PUT(relativePath, handlers...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().OPTIONS(relativePath, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().HEAD(relativePath, handlers...)
}

func Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine().Any(relativePath, handlers...)
}

func StaticFile(relativePath, filepath string) IRoutes {
	return engine().StaticFile(relativePath, filepath)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func Static(relativePath, root string) IRoutes {
	return engine().Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	return engine().StaticFS(relativePath, fs)
}

// Attachs a global middleware to the router. ie. the middlewares attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func Use(middlewares ...HandlerFunc) IRoutes {
	return engine().Use(middlewares...)
}

// The router is attached to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func Run(addr ...string) (err error) {
	return engine().Run(addr...)
}

// The router is attached to a http.Server and starts listening and serving HTTPS requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func RunTLS(addr string, certFile string, keyFile string) (err error) {
	return engine().RunTLS(addr, certFile, keyFile)
}

// The router is attached to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func RunUnix(file string) (err error) {
	return engine().RunUnix(file)
}
