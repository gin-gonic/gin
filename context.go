// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/net/context"
)

const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
)

const AbortIndex = math.MaxInt8 / 2

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
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	context.Context
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8

	Engine   *Engine
	Keys     map[string]interface{}
	Errors   errorMsgs
	Accepted []string
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
}

func (c *Context) Copy() *Context {
	var cp Context = *c
	cp.writermem.ResponseWriter = nil
	cp.Writer = &cp.writermem
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

func (c *Context) IsAborted() bool {
	return c.index == AbortIndex
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

func (c *Context) ErrorTyped(err error, typ int, meta interface{}) {
	c.Errors = append(c.Errors, errorMsg{
		Error: err,
		Flags: typ,
		Meta:  meta,
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
		return c.Errors[nuErrors-1].Error
	}
	return nil
}

/************************************/
/************ INPUT DATA ************/
/************************************/

/** Shortcut for c.Request.FormValue(key) */
func (c *Context) FormValue(key string) (va string) {
	va, _ = c.formValue(key)
	return
}

/** Shortcut for c.Request.PostFormValue(key) */
func (c *Context) PostFormValue(key string) (va string) {
	va, _ = c.postFormValue(key)
	return
}

/** Shortcut for c.Params.ByName(key) */
func (c *Context) ParamValue(key string) (va string) {
	va, _ = c.paramValue(key)
	return
}

func (c *Context) DefaultPostFormValue(key, defaultValue string) string {
	if va, ok := c.postFormValue(key); ok {
		return va
	}
	return defaultValue
}

func (c *Context) DefaultFormValue(key, defaultValue string) string {
	if va, ok := c.formValue(key); ok {
		return va
	}
	return defaultValue
}

func (c *Context) DefaultParamValue(key, defaultValue string) string {
	if va, ok := c.paramValue(key); ok {
		return va
	}
	return defaultValue
}

func (c *Context) paramValue(key string) (string, bool) {
	return c.Params.Get(key)
}

func (c *Context) formValue(key string) (string, bool) {
	req := c.Request
	req.ParseForm()
	if values, ok := req.Form[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (c *Context) postFormValue(key string) (string, bool) {
	req := c.Request
	req.ParseForm()
	if values, ok := req.PostForm[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Sets a new pair key/value just for the specified context.
// It also lazy initializes the hashmap.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get returns the value for the given key or an error if the key does not exist.
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.Keys != nil {
		value, exists = c.Keys[key]
	}
	return
}

// MustGet returns the value for the given key or panics if the value doesn't exist.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	} else {
		panic("Key \"" + key + "\" does not exist")
	}
}

func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return c.Context.Value(key)
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
		return strings.TrimSpace(clientIP)
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

func (c *Context) renderingError(err error, meta ...interface{}) {
	c.ErrorTyped(err, ErrorTypeInternal, meta)
	c.AbortWithStatus(500)
}

func (c *Context) Render(code int, render render.Render, obj ...interface{}) {
	if err := render.Render(c.Writer, code, obj...); err != nil {
		c.renderingError(err, obj)
	}
}

// Renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj interface{}) {
	c.Render(code, c.Engine.HTMLRender, name, obj)
}

func (c *Context) IndentedJSON(code int, obj interface{}) {
	if err := render.WriteIndentedJSON(c.Writer, code, obj); err != nil {
		c.renderingError(err, obj)
	}
}

// Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	if err := render.WriteJSON(c.Writer, code, obj); err != nil {
		c.renderingError(err, obj)
	}
}

// Serializes the given struct as XML into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	if err := render.WriteXML(c.Writer, code, obj); err != nil {
		c.renderingError(err, obj)
	}
}

// Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) String(code int, format string, values ...interface{}) {
	render.WritePlainText(c.Writer, code, format, values)
}

// Writes the given string into the response body and sets the Content-Type to "text/html" without template.
func (c *Context) HTMLString(code int, format string, values ...interface{}) {
	render.WriteHTMLString(c.Writer, code, format, values)
}

// Returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	render.WriteRedirect(c.Writer, code, c.Request, location)
}

// Writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	render.WriteData(c.Writer, code, contentType, data)
}

// Writes the specified file into the body stream
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) SSEvent(name string, message interface{}) {
	render.WriteSSEvent(c.Writer, name, message)
}

func (c *Context) Header(code int, headers map[string]string) {
	if len(headers) > 0 {
		header := c.Writer.Header()
		for key, value := range headers {
			header.Set(key, value)
		}
	}
	c.Writer.WriteHeader(code)
}

func (c *Context) Stream(step func(w io.Writer) bool) {
	w := c.Writer
	clientGone := w.CloseNotify()
	for {
		select {
		case <-clientGone:
			return
		default:
			keepopen := step(w)
			w.Flush()
			if !keepopen {
				return
			}
		}
	}
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
			panic("negotiate config is wrong. html path is needed")
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
		panic("you must provide at least one offer")
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
