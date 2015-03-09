// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/gin-gonic/gin/render"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"math"
	"net/http"
	"sync"
)

const (
	AbortIndex            = math.MaxInt8 / 2
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type (
	HandlerFunc func(*Context)

	// Represents the web framework, it wraps the blazing fast httprouter multiplexer and a list of global middlewares.
	Engine struct {
		*RouterGroup
		HTMLRender         render.Render
		Default404Body     []byte
		Default405Body     []byte
		pool               sync.Pool
		allNoRouteNoMethod []HandlerFunc
		noRoute            []HandlerFunc
		noMethod           []HandlerFunc
		router             *httprouter.Router
	}
)

// Returns a new blank Engine instance without any middleware attached.
// The most basic configuration
func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		Handlers:     nil,
		absolutePath: "/",
		engine:       engine,
	}
	engine.router = httprouter.New()
	engine.Default404Body = []byte("404 page not found")
	engine.Default405Body = []byte("405 method not allowed")
	engine.router.NotFound = engine.handle404
	engine.router.MethodNotAllowed = engine.handle405
	engine.pool.New = func() interface{} {
		c := &Context{Engine: engine}
		c.Writer = &c.writermem
		return c
	}
	return engine
}

// Returns a Engine instance with the Logger and Recovery already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	if IsDebugging() {
		render.HTMLDebug.AddGlob(pattern)
		engine.HTMLRender = render.HTMLDebug
	} else {
		templ := template.Must(template.ParseGlob(pattern))
		engine.SetHTMLTemplate(templ)
	}
}

func (engine *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		render.HTMLDebug.AddFiles(files...)
		engine.HTMLRender = render.HTMLDebug
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
	engine.allNoRouteNoMethod = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoRouteNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) handle404(w http.ResponseWriter, req *http.Request) {
	c := engine.createContext(w, req, nil, engine.allNoRouteNoMethod)
	// set 404 by default, useful for logging
	c.Writer.WriteHeader(404)
	c.Next()
	if !c.Writer.Written() {
		if c.Writer.Status() == 404 {
			c.Data(-1, MIMEPlain, engine.Default404Body)
		} else {
			c.Writer.WriteHeaderNow()
		}
	}
	engine.reuseContext(c)
}

func (engine *Engine) handle405(w http.ResponseWriter, req *http.Request) {
	c := engine.createContext(w, req, nil, engine.allNoRouteNoMethod)
	// set 405 by default, useful for logging
	c.Writer.WriteHeader(405)
	c.Next()
	if !c.Writer.Written() {
		if c.Writer.Status() == 405 {
			c.Data(-1, MIMEPlain, engine.Default405Body)
		} else {
			c.Writer.WriteHeaderNow()
		}
	}
	engine.reuseContext(c)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	engine.router.ServeHTTP(writer, request)
}

func (engine *Engine) Run(addr string) error {
	debugPrint("Listening and serving HTTP on %s\n", addr)
	if err := http.ListenAndServe(addr, engine); err != nil {
		return err
	}
	return nil
}

func (engine *Engine) RunTLS(addr string, cert string, key string) error {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	if err := http.ListenAndServeTLS(addr, cert, key, engine); err != nil {
		return err
	}
	return nil
}
