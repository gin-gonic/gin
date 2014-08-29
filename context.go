// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const (
	ErrorTypeInternal = 1 << iota
	ErrorTypeExternal = 1 << iota
	ErrorTypeAll      = 0xffffffff
)

// Used internally to collect errors that occurred during an http request.
type errorMsg struct {
	Err  string      `json:"error"`
	Type uint32      `json:"-"`
	Meta interface{} `json:"meta"`
}

type errorMsgs []errorMsg

func (a errorMsgs) ByType(typ uint32) errorMsgs {
	if len(a) == 0 {
		return a
	}
	result := make(errorMsgs, 0, len(a))
	for _, msg := range a {
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}
	return result
}

func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		text := fmt.Sprintf("Error #%02d: %s \n     Meta: %v\n", (i + 1), msg.Err, msg.Meta)
		buffer.WriteString(text)
	}
	return buffer.String()
}

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter
	Keys      map[string]interface{}
	Errors    errorMsgs
	Params    httprouter.Params
	Engine    *Engine
	handlers  []HandlerFunc
	index     int8
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

func (engine *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	c := engine.cache.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.Params = params
	c.handlers = handlers
	c.Keys = nil
	c.index = -1
	c.Errors = c.Errors[0:0]
	return c
}

/************************************/
/****** FLOW AND ERROR MANAGEMENT****/
/************************************/

func (c *Context) Copy() *Context {
	var cp Context = *c
	cp.index = AbortIndex
	cp.handlers = nil
	return &cp
}

// Next should be used only in the middlewares.
// It executes the pending handlers in the chain inside the calling handler.
// See example in github.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Forces the system to do not continue calling the pending handlers.
// For example, the first handler checks if the request is authorized. If it's not, context.Abort(401) should be called.
// The rest of pending handlers would never be called for that request.
func (c *Context) Abort(code int) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	c.index = AbortIndex
}

// Fail is the same as Abort plus an error message.
// Calling `context.Fail(500, err)` is equivalent to:
// ```
// context.Error("Operation aborted", err)
// context.Abort(500)
// ```
func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.Abort(code)
}

func (c *Context) ErrorTyped(err error, typ uint32, meta interface{}) {
	c.Errors = append(c.Errors, errorMsg{
		Err:  err.Error(),
		Type: typ,
		Meta: meta,
	})
}

// Attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together, print a log, or append it in the HTTP response.
func (c *Context) Error(err error, meta interface{}) {
	c.ErrorTyped(err, ErrorTypeExternal, meta)
}

func (c *Context) LastError() error {
	s := len(c.Errors)
	if s > 0 {
		return errors.New(c.Errors[s-1].Err)
	} else {
		return nil
	}
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Sets a new pair key/value just for the specified context.
// It also lazy initializes the hashmap.
func (c *Context) Set(key string, item interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = item
}

// Get returns the value for the given key or an error if the key does not exist.
func (c *Context) Get(key string) (interface{}, error) {
	if c.Keys != nil {
		item, ok := c.Keys[key]
		if ok {
			return item, nil
		}
	}
	return nil, errors.New("Key does not exist.")
}

// MustGet returns the value for the given key or panics if the value doesn't exist.
func (c *Context) MustGet(key string) interface{} {
	value, err := c.Get(key)
	if err != nil || value == nil {
		log.Panicf("Key %s doesn't exist", key)
	}
	return value
}

/************************************/
/******** ENCOGING MANAGEMENT********/
/************************************/

// This function checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// "application/json" --> JSON binding
// "application/xml"  --> XML binding
// else --> returns an error
// if Parses the request's body as JSON if Content-Type == "application/json"  using JSON or XML  as a JSON input. It decodes the json payload into the struct specified as a pointer.Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(obj interface{}) bool {
	var b binding.Binding
	ctype := filterFlags(c.Request.Header.Get("Content-Type"))
	switch {
	case c.Request.Method == "GET" || ctype == MIMEPOSTForm:
		b = binding.Form
	case ctype == MIMEJSON:
		b = binding.JSON
	case ctype == MIMEXML || ctype == MIMEXML2:
		b = binding.XML
	default:
		c.Fail(400, errors.New("unknown content-type: "+ctype))
		return false
	}
	return c.BindWith(obj, b)
}

func (c *Context) BindWith(obj interface{}, b binding.Binding) bool {
	if err := b.Bind(c.Request, obj); err != nil {
		c.Fail(400, err)
		return false
	}
	return true
}

func (c *Context) Render(code int, render render.Render, obj ...interface{}) {
	if err := render.Render(c.Writer, code, obj...); err != nil {
		c.ErrorTyped(err, ErrorTypeInternal, obj)
		c.Abort(500)
	}
}

// Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON, obj)
}

// Serializes the given struct as XML into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	c.Render(code, render.XML, obj)
}

// Renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj interface{}) {
	c.Render(code, c.Engine.HTMLRender, name, obj)
}

// Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.Plain, format, values)
}

// Returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	if code >= 300 && code <= 308 {
		c.Render(code, render.Redirect, location)
	} else {
		panic(fmt.Sprintf("Cannot send a redirect with status code %d", code))
	}
}

// Writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	c.Writer.Write(data)
}

// Writes the specified file into the body stream
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}
