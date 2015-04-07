// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/julienschmidt/httprouter"
)

const AbortIndex = math.MaxInt8 / 2

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	Engine    *Engine
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   httprouter.Params
	Input    inputHolder
	handlers []HandlerFunc
	index    int8

	Keys     map[string]interface{}
	Errors   errorMsgs
	Accepted []string
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.Keys = nil
	c.index = -1
	c.Accepted = nil
	c.Errors = c.Errors[0:0]
}

func (c *Context) Copy() *Context {
	var cp Context = *c
	cp.index = AbortIndex
	cp.handlers = nil
	return &cp
}

/************************************/
/*************** FLOW ***************/
/************************************/

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

// Forces the system to not continue calling the pending handlers in the chain.
func (c *Context) Abort() {
	c.index = AbortIndex
}

// Same than AbortWithStatus() but also writes the specified response status code.
// For example, the first handler checks if the request is authorized. If it's not, context.AbortWithStatus(401) should be called.
func (c *Context) AbortWithStatus(code int) {
	c.Writer.WriteHeader(code)
	c.Abort()
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Fail is the same as Abort plus an error message.
// Calling `context.Fail(500, err)` is equivalent to:
// ```
// context.Error("Operation aborted", err)
// context.AbortWithStatus(500)
// ```
func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.AbortWithStatus(code)
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
	nuErrors := len(c.Errors)
	if nuErrors > 0 {
		return errors.New(c.Errors[nuErrors-1].Err)
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
func (c *Context) Get(key string) (value interface{}, ok bool) {
	if c.Keys != nil {
		value, ok = c.Keys[key]
	}
	return
}

// MustGet returns the value for the given key or panics if the value doesn't exist.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	} else {
		log.Panicf("Key %s does not exist", key)
	}
	return nil
}

/************************************/
/********* PARSING REQUEST **********/
/************************************/

func (c *Context) ClientIP() string {
	clientIP := c.Request.Header.Get("X-Real-IP")
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = c.Request.Header.Get("X-Forwarded-For")
	clientIP = strings.Split(clientIP, ",")[0]
	if len(clientIP) > 0 {
		return clientIP
	}
	return c.Request.RemoteAddr
}

func (c *Context) ContentType() string {
	return filterFlags(c.Request.Header.Get("Content-Type"))
}

// This function checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// "application/json" --> JSON binding
// "application/xml"  --> XML binding
// else --> returns an error
// if Parses the request's body as JSON if Content-Type == "application/json"  using JSON or XML  as a JSON input. It decodes the json payload into the struct specified as a pointer.Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(obj interface{}) bool {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

func (c *Context) BindWith(obj interface{}, b binding.Binding) bool {
	if err := b.Bind(c.Request, obj); err != nil {
		c.Fail(400, err)
		return false
	}
	return true
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

func (c *Context) Render(code int, render render.Render, obj ...interface{}) {
	if err := render.Render(c.Writer, code, obj...); err != nil {
		c.ErrorTyped(err, ErrorTypeInternal, obj)
		c.AbortWithStatus(500)
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

// Writes the given string into the response body and sets the Content-Type to "text/html" without template.
func (c *Context) HTMLString(code int, format string, values ...interface{}) {
	c.Render(code, render.HTMLPlain, format, values)
}

// Returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	if code >= 300 && code <= 308 {
		c.Render(code, render.Redirect, c.Request, location)
	} else {
		log.Panicf("Cannot redirect with status code %d", code)
	}
}

// Writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	c.Writer.WriteHeader(code)
	c.Writer.Write(data)
}

// Writes the specified file into the body stream
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

/************************************/
/******** CONTENT NEGOTIATION *******/
/************************************/

type Negotiate struct {
	Offered  []string
	HTMLPath string
	HTMLData interface{}
	JSONData interface{}
	XMLData  interface{}
	Data     interface{}
}

func (c *Context) Negotiate(code int, config Negotiate) {
	switch c.NegotiateFormat(config.Offered...) {
	case binding.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case binding.MIMEHTML:
		if len(config.HTMLPath) == 0 {
			log.Panic("negotiate config is wrong. html path is needed")
		}
		data := chooseData(config.HTMLData, config.Data)
		c.HTML(code, config.HTMLPath, data)

	case binding.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		c.XML(code, data)

	default:
		c.Fail(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server"))
	}
}

func (c *Context) NegotiateFormat(offered ...string) string {
	if len(offered) == 0 {
		log.Panic("you must provide at least one offer")
	}
	if c.Accepted == nil {
		c.Accepted = parseAccept(c.Request.Header.Get("Accept"))
	}
	if len(c.Accepted) == 0 {
		return offered[0]
	}
	for _, accepted := range c.Accepted {
		for _, offert := range offered {
			if accepted == offert {
				return offert
			}
		}
	}
	return ""
}

func (c *Context) SetAccepted(formats ...string) {
	c.Accepted = formats
}
