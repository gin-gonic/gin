// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")

type (
	HandlerFunc func(*Context)

	// Represents the web framework, it wraps the blazing fast httprouter multiplexer and a list of global middlewares.
	Engine struct {
		*RouterGroup
		HTMLRender  render.Render
		pool        sync.Pool
		allNoRoute  []HandlerFunc
		allNoMethod []HandlerFunc
		noRoute     []HandlerFunc
		noMethod    []HandlerFunc
		trees       map[string]*node

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool

		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool

		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool
	}
)

// Returns a new blank Engine instance without any middleware attached.
// The most basic configuration
func New() *Engine {
	engine := &Engine{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		trees: make(map[string]*node),
	}
	engine.RouterGroup = &RouterGroup{
		Handlers:     nil,
		absolutePath: "/",
		engine:       engine,
	}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

// Returns a Engine instance with the Logger and Recovery already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (engine *Engine) allocateContext() (context *Context) {
	context = &Context{Engine: engine}
	context.Writer = &context.writermem
	context.Input = inputHolder{context: context}
	return
}

func (engine *Engine) createContext(w http.ResponseWriter, req *http.Request) *Context {
	c := engine.pool.Get().(*Context)
	c.reset()
	c.writermem.reset(w)
	c.Request = req
	return c
}

func (engine *Engine) reuseContext(c *Context) {
	engine.pool.Put(c)
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	if IsDebugging() {
		r := &render.HTMLDebugRender{Glob: pattern}
		engine.HTMLRender = r
	} else {
		templ := template.Must(template.ParseGlob(pattern))
		engine.SetHTMLTemplate(templ)
	}
}

func (engine *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		r := &render.HTMLDebugRender{Files: files}
		engine.HTMLRender = r
	} else {
		templ := template.Must(template.ParseFiles(files...))
		engine.SetHTMLTemplate(templ)
	}
}

func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
	engine.HTMLRender = render.HTMLRender{
		Template: templ,
	}
}

// Adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

func (engine *Engine) Use(middlewares ...HandlerFunc) {
	engine.RouterGroup.Use(middlewares...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) handle404(c *Context) {
	// set 404 by default, useful for logging
	c.handlers = engine.allNoRoute
	c.Writer.WriteHeader(404)
	c.Next()
	if !c.Writer.Written() {
		if c.Writer.Status() == 404 {
			c.Data(-1, binding.MIMEPlain, default404Body)
		} else {
			c.Writer.WriteHeaderNow()
		}
	}
}

func (engine *Engine) handle405(c *Context) {
	// set 405 by default, useful for logging
	c.handlers = engine.allNoMethod
	c.Writer.WriteHeader(405)
	c.Next()
	if !c.Writer.Written() {
		if c.Writer.Status() == 405 {
			c.Data(-1, binding.MIMEPlain, default405Body)
		} else {
			c.Writer.WriteHeaderNow()
		}
	}
}

func (engine *Engine) handle(method, path string, handlers []HandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}

	//methodCode := codeForHTTPMethod(method)
	root := engine.trees[method]
	if root == nil {
		root = new(node)
		engine.trees[method] = root
	}
	root.addRoute(path, handlers)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.createContext(w, req)
	//methodCode := codeForHTTPMethod(req.Method)
	if root := engine.trees[req.Method]; root != nil {
		path := req.URL.Path
		if handlers, params, _ := root.getValue(path, c.Params); handlers != nil {
			c.handlers = handlers
			c.Params = params
			c.Next()
			engine.reuseContext(c)
			return
		}
	}

	// Handle 404
	engine.handle404(c)
	engine.reuseContext(c)
}

func (engine *Engine) Run(addr string) error {
	debugPrint("Listening and serving HTTP on %s\n", addr)
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) RunTLS(addr string, cert string, key string) error {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	return http.ListenAndServeTLS(addr, cert, key, engine)
}
