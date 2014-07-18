package gin

import (
	"crypto/tls"
	"errors"
	"github.com/gin-gonic/gin/render"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"math"
	"net"
	"net/http"
	"path"
	"sync"
	"time"
)

const (
	AbortIndex   = math.MaxInt8 / 2
	MIMEJSON     = "application/json"
	MIMEHTML     = "text/html"
	MIMEXML      = "application/xml"
	MIMEXML2     = "text/xml"
	MIMEPlain    = "text/plain"
	MIMEPOSTForm = "application/x-www-form-urlencoded"
)

type (
	HandlerFunc func(*Context)

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
		HTMLRender   render.Render
		cache        sync.Pool
		finalNoRoute []HandlerFunc
		noRoute      []HandlerFunc
		router       *httprouter.Router
		listener     *stoppableListener
	}
)

func (engine *Engine) handle404(w http.ResponseWriter, req *http.Request) {
	c := engine.createContext(w, req, nil, engine.finalNoRoute)
	c.Writer.setStatus(404)
	c.Next()
	if !c.Writer.Written() {
		c.Data(404, MIMEPlain, []byte("404 page not found"))
	}
	engine.cache.Put(c)
}

// Returns a new blank Engine instance without any middleware attached.
// The most basic configuration
func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{nil, "/", nil, engine}
	engine.router = httprouter.New()
	engine.router.NotFound = engine.handle404
	engine.cache.New = func() interface{} {
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
	templ := template.Must(template.ParseGlob(pattern))
	engine.SetHTMLTemplate(templ)
}

func (engine *Engine) LoadHTMLFiles(files ...string) {
	templ := template.Must(template.ParseFiles(files...))
	engine.SetHTMLTemplate(templ)
}

func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
	engine.HTMLRender = render.HTMLRender{
		Template: templ,
	}
}

// Adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.finalNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) Use(middlewares ...HandlerFunc) {
	engine.RouterGroup.Use(middlewares...)
	engine.finalNoRoute = engine.combineHandlers(engine.noRoute)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(w, req)
}

func (engine *Engine) Run(addr string) {
	server := &http.Server{Addr: addr, Handler: engine}
	err := engine.listenAndServe(server)
	if err != nil {
		if err != stoppedError {
			panic(err)
		}
	}
}

func (engine *Engine) RunTLS(addr string, cert string, key string) {
	server := &http.Server{Addr: addr, Handler: engine}
	err := engine.listenAndServeTLS(server, cert, key)
	if err != nil {
		if err != stoppedError {
			panic(err)
		}
	}
}

func (engine *Engine) Stop() {
	engine.listener.Stop()
}

// Inlined from net/http source so we can inject our own listener.
func (engine *Engine) listenAndServe(srv *http.Server) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	engine.listener, err = newStoppableListener(tcpKeepAliveListener{ln.(*net.TCPListener)})
	if err != nil {
		return err
	}
	return srv.Serve(engine.listener)
}

// Inlined from net/http source so we can inject our own listener.
func (engine *Engine) listenAndServeTLS(srv *http.Server, certFile string, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	engine.listener, err = newStoppableListener(tlsListener)
	if err != nil {
		return err
	}
	return srv.Serve(engine.listener)
}

/************************************/
/******** KEEP ALIVE LISTENER *******/
/************************************/

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

/************************************/
/******** STOPPABLE LISTENER ********/
/************************************/

var stoppedError = errors.New("Webserver is being stopped")

type stoppableListener struct {
	tcpKeepAliveListener          //Wrapped listener
	stop                 chan int //Channel used only to indicate listener should shutdown
}

func newStoppableListener(l net.Listener) (*stoppableListener, error) {
	tcpL, ok := l.(tcpKeepAliveListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &stoppableListener{}
	retval.tcpKeepAliveListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

func (sl *stoppableListener) Accept() (net.Conn, error) {
	for {
		//Wait up to one second for a new connection
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.tcpKeepAliveListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, stoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *stoppableListener) Stop() {
	close(sl.stop)
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

// Adds middlewares to the group, see example code in github.
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := group.pathFor(component)

	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		parent:   group,
		prefix:   prefix,
		engine:   group.engine,
	}
}

func (group *RouterGroup) pathFor(p string) string {
	joined := path.Join(group.prefix, p)
	// Append a '/' if the last component had one, but only if it's not there already
	if len(p) > 0 && p[len(p)-1] == '/' && joined[len(joined)-1] != '/' {
		return joined + "/"
	}
	return joined
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
	p = group.pathFor(p)
	handlers = group.combineHandlers(handlers)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := group.engine.createContext(w, req, params, handlers)
		c.Next()
		group.engine.cache.Put(c)
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

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (group *RouterGroup) Static(p, root string) {
	prefix := group.pathFor(p)
	p = path.Join(p, "/*filepath")
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	group.GET(p, func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
	group.HEAD(p, func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(group.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, group.Handlers...)
	h = append(h, handlers...)
	return h
}
