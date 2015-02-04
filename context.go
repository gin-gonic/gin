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
	"net"
	"net/http"
	"strings"
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
	accepted  []string
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (engine *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.Params = params
	c.handlers = handlers
	c.Keys = nil
	c.index = -1
	c.accepted = nil
	c.Errors = c.Errors[0:0]
	return c
}

func (engine *Engine) reuseContext(c *Context) {
	engine.pool.Put(c)
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

// Forces the system to do not continue calling the pending handlers in the chain.
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
func (c *Context) Get(key string) (interface{}, error) {
	if c.Keys != nil {
		value, ok := c.Keys[key]
		if ok {
			return value, nil
		}
	}
	return nil, errors.New("Key does not exist.")
}

// MustGet returns the value for the given key or panics if the value doesn't exist.
func (c *Context) MustGet(key string) interface{} {
	value, err := c.Get(key)
	if err != nil || value == nil {
		log.Panicf("Key %s doesn't exist", value)
	}
	return value
}

func ipInMasks(ip net.IP, masks []interface{}) bool {
	for _, proxy := range masks {
		var mask *net.IPNet
		var err error

		switch t := proxy.(type) {
		case string:
			if _, mask, err = net.ParseCIDR(t); err != nil {
				panic(err)
			}
		case net.IP:
			mask = &net.IPNet{IP: t, Mask: net.CIDRMask(len(t)*8, len(t)*8)}
		case net.IPNet:
			mask = &t
		}

		if mask.Contains(ip) {
			return true
		}
	}

	return false
}

// the ForwardedFor middleware unwraps the X-Forwarded-For headers, be careful to only use this
// middleware if you've got servers in front of this server. The list with (known) proxies and
// local ips are being filtered out of the forwarded for list, giving the last not local ip being
// the real client ip.
func ForwardedFor(proxies ...interface{}) HandlerFunc {
	if len(proxies) == 0 {
		// default to local ips
		var reservedLocalIps = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

		proxies = make([]interface{}, len(reservedLocalIps))

		for i, v := range reservedLocalIps {
			proxies[i] = v
		}
	}

	return func(c *Context) {
		// the X-Forwarded-For header contains an array with left most the client ip, then
		// comma separated, all proxies the request passed. The last proxy appears
		// as the remote address of the request. Returning the client
		// ip to comply with default RemoteAddr response.

		// check if remoteaddr is local ip or in list of defined proxies
		remoteIp := net.ParseIP(strings.Split(c.Request.RemoteAddr, ":")[0])

		if !ipInMasks(remoteIp, proxies) {
			return
		}

		if forwardedFor := c.Request.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")

			for i := len(parts) - 1; i >= 0; i-- {
				part := parts[i]

				ip := net.ParseIP(strings.TrimSpace(part))

				if ipInMasks(ip, proxies) {
					continue
				}

				// returning remote addr conform the original remote addr format
				c.Request.RemoteAddr = ip.String() + ":0"

				// remove forwarded for address
				c.Request.Header.Set("X-Forwarded-For", "")
				return
			}
		}
	}
}

func (c *Context) ClientIP() string {
	return c.Request.RemoteAddr
}

/************************************/
/********* PARSING REQUEST **********/
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
	case MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		if len(config.HTMLPath) == 0 {
			panic("negotiate config is wrong. html path is needed")
		}
		c.HTML(code, config.HTMLPath, data)

	case MIMEXML:
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
	if c.accepted == nil {
		c.accepted = parseAccept(c.Request.Header.Get("Accept"))
	}
	if len(c.accepted) == 0 {
		return offered[0]

	} else {
		for _, accepted := range c.accepted {
			for _, offert := range offered {
				if accepted == offert {
					return offert
				}
			}
		}
		return ""
	}
}

func (c *Context) SetAccepted(formats ...string) {
	c.accepted = formats
}
