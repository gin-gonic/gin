package gin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"
)

const (
	AbortIndex = math.MaxInt8 / 2
	MIMEJSON   = "application/json"
	MIMEHTML   = "text/html"
	MIMEXML    = "application/xml"
	MIMEXML2   = "text/xml"
	MIMEPlain  = "text/plain"
)

type (
	HandlerFunc func(*Context)

	H map[string]interface{}

	// Used internally to collect errors that occurred during an http request.
	ErrorMsg struct {
		Err  string      `json:"error"`
		Meta interface{} `json:"meta"`
	}

	ErrorMsgs []ErrorMsg

	Config struct {
		CacheSize    int
		Preallocated int
	}

	// Context is the most important part of gin. It allows us to pass variables between middleware,
	// manage the flow, validate the JSON of a request and render a JSON response for example.
	Context struct {
		Req      *http.Request
		Writer   ResponseWriter
		Keys     map[string]interface{}
		Errors   ErrorMsgs
		Params   httprouter.Params
		Engine   *Engine
		handlers []HandlerFunc
		index    int8
	}

	// Used internally to configure router, a RouterGroup is associated with a prefix
	// and an array of handlers (middlewares)
	RouterGroup struct {
		Handlers []HandlerFunc
		prefix   string
		parent   *RouterGroup
		engine   *Engine
	}

	// Represents the web framework, it wraps the blazing fast httprouter multiplexer and a list of global middlewares.
	Engine struct {
		*RouterGroup
		HTMLTemplates *template.Template
		cache         chan *Context
		handlers404   []HandlerFunc
		router        *httprouter.Router
	}
)

var (
	DefaultConfig = Config{
		CacheSize:    1024,
		Preallocated: 512,
	}
)

func (a ErrorMsgs) String() string {
	var buffer bytes.Buffer
	for i, msg := range a {
		text := fmt.Sprintf("Error #%02d: %s \n     Meta: %v\n", (i + 1), msg.Err, msg.Meta)
		buffer.WriteString(text)
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func NewWithConfig(config Config) *Engine {
	if config.CacheSize < 2 {
		panic("CacheSize must be at least 2")
	}
	if config.Preallocated > config.CacheSize {
		panic("Preallocated must be less or equal to CacheSize")
	}
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{nil, "/", nil, engine}
	engine.router = httprouter.New()
	engine.router.NotFound = engine.handle404
	engine.cache = make(chan *Context, config.CacheSize)

	// Fill it with empty contexts
	for i := 0; i < config.Preallocated; i++ {
		engine.cache <- &Context{Engine: engine, Writer: &responseWriter{}}
	}
	return engine
}

// Returns a new blank Engine instance without any middleware attached.
// The most basic configuration
func New() *Engine {
	return NewWithConfig(DefaultConfig)
}

// Returns a Engine instance with the Logger and Recovery already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (engine *Engine) LoadHTMLTemplates(pattern string) {
	engine.HTMLTemplates = template.Must(template.ParseGlob(pattern))
}

// Adds handlers for NotFound. It return a 404 code by default.
func (engine *Engine) NotFound404(handlers ...HandlerFunc) {
	engine.handlers404 = handlers
}

func (engine *Engine) CacheStress() float32 {
	return 1.0 - float32(len(engine.cache))/float32(cap(engine.cache))
}

func (engine *Engine) handle404(w http.ResponseWriter, req *http.Request) {
	handlers := engine.combineHandlers(engine.handlers404)
	c := engine.createContext(w, req, nil, handlers)
	c.Writer.setStatus(404)
	c.Next()
	if !c.Writer.Written() {
		c.String(404, "404 page not found")
	}
	engine.reuseContext(c)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (engine *Engine) ServeFiles(path string, root http.FileSystem) {
	engine.router.ServeFiles(path, root)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(w, req)
}

func (engine *Engine) Run(addr string) {
	if err := http.ListenAndServe(addr, engine); err != nil {
		panic(err)
	}
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

func (engine *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	select {
	case c := <-engine.cache:
		c.Writer.reset(w)
		c.Req = req
		c.Params = params
		c.handlers = handlers
		c.Keys = nil
		c.index = -1
		return c
	default:
		return &Context{
			Writer:   &responseWriter{w, -1, false},
			Req:      req,
			Params:   params,
			handlers: handlers,
			index:    -1,
			Engine:   engine,
		}
	}
}

func (engine *Engine) reuseContext(c *Context) {
	select {
	case engine.cache <- c:
	default:
	}
}

// Adds middlewares to the group, see example code in github.
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		parent:   group,
		prefix:   prefix,
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
func (group *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	p = path.Join(group.prefix, p)
	handlers = group.combineHandlers(handlers)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := group.engine.createContext(w, req, params, handlers)
		c.Next()
		group.engine.reuseContext(c)
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.Handle("POST", path, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	group.Handle("GET", path, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	group.Handle("DELETE", path, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	group.Handle("PATCH", path, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	group.Handle("PUT", path, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (group *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) {
	group.Handle("OPTIONS", path, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (group *RouterGroup) HEAD(path string, handlers ...HandlerFunc) {
	group.Handle("HEAD", path, handlers)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(group.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, group.Handlers...)
	h = append(h, handlers...)
	return h
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

// Attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together, print a log, or append it in the HTTP response.
func (c *Context) Error(err error, meta interface{}) {
	c.Errors = append(c.Errors, ErrorMsg{
		Err:  err.Error(),
		Meta: meta,
	})
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

func filterFlags(content string) string {
	for i, a := range content {
		if a == ' ' || a == ';' {
			return content[:i]
		}
	}
	return content
}

// This function checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// "application/json" --> JSON binding
// "application/xml"  --> XML binding
// else --> returns an error
// if Parses the request's body as JSON if Content-Type == "application/json"  using JSON or XML  as a JSON input. It decodes the json payload into the struct specified as a pointer.Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(obj interface{}) bool {
	var b binding.Binding
	ctype := filterFlags(c.Req.Header.Get("Content-Type"))
	switch {
	case c.Req.Method == "GET":
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
	if err := b.Bind(c.Req, obj); err != nil {
		c.Fail(400, err)
		return false
	}
	return true
}

// Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", MIMEJSON)
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.Error(err, obj)
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Serializes the given struct as XML into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", MIMEXML)
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.Error(err, obj)
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Writer.Header().Set("Content-Type", MIMEHTML)
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	if err := c.Engine.HTMLTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Error(err, map[string]interface{}{
			"name": name,
			"data": data,
		})
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) String(code int, msg string) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	c.Writer.Header().Set("Content-Type", MIMEPlain)
	c.Writer.Write([]byte(msg))
}

// Writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, data []byte) {
	c.Writer.WriteHeader(code)
	c.Writer.Write(data)
}
