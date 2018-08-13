# CHANGELOG

### Gin 1.3.0

- [NEW] Add [`func (*Context) QueryMap`](https://godoc.org/github.com/gin-gonic/gin#Context.QueryMap), [`func (*Context) GetQueryMap`](https://godoc.org/github.com/gin-gonic/gin#Context.GetQueryMap), [`func (*Context) PostFormMap`](https://godoc.org/github.com/gin-gonic/gin#Context.PostFormMap) and [`func (*Context) GetPostFormMap`](https://godoc.org/github.com/gin-gonic/gin#Context.GetPostFormMap) to support `type map[string]string` as query string or form parameters, see [#1383](https://github.com/gin-gonic/gin/pull/1383)
- [NEW] Add [`func (*Context) AsciiJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.AsciiJSON), see [#1358](https://github.com/gin-gonic/gin/pull/1358)
- [NEW] Add `Pusher()` in [`type ResponseWriter`](https://godoc.org/github.com/gin-gonic/gin#ResponseWriter) for supporting http2 push, see [#1273](https://github.com/gin-gonic/gin/pull/1273)
- [NEW] Add [`func (*Context) DataFromReader`](https://godoc.org/github.com/gin-gonic/gin#Context.DataFromReader) for serving dynamic data, see [#1304](https://github.com/gin-gonic/gin/pull/1304)
- [NEW] Add [`func (*Context) ShouldBindBodyWith`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindBodyWith) allowing to call binding multiple times, see [#1341](https://github.com/gin-gonic/gin/pull/1341)
- [NEW] Support pointers in form binding, see [#1336](https://github.com/gin-gonic/gin/pull/1336)
- [NEW] Add [`func (*Context) JSONP`](https://godoc.org/github.com/gin-gonic/gin#Context.JSONP), see [#1333](https://github.com/gin-gonic/gin/pull/1333)
- [NEW] Support default value in form binding, see [#1138](https://github.com/gin-gonic/gin/pull/1138)
- [NEW] Expose validator engine in [`type StructValidator`](https://godoc.org/github.com/gin-gonic/gin/binding#StructValidator), see [#1277](https://github.com/gin-gonic/gin/pull/1277)
- [NEW] Add [`func (*Context) ShouldBind`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBind), [`func (*Context) ShouldBindQuery`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindQuery) and [`func (*Context) ShouldBindJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindJSON), see [#1047](https://github.com/gin-gonic/gin/pull/1047)
- [NEW] Add support for `time.Time` location in form binding, see [#1117](https://github.com/gin-gonic/gin/pull/1117)
- [NEW] Add [`func (*Context) BindQuery`](https://godoc.org/github.com/gin-gonic/gin#Context.BindQuery), see [#1029](https://github.com/gin-gonic/gin/pull/1029)
- [NEW] Make [jsonite](https://github.com/json-iterator/go) optional with build tags, see [#1026](https://github.com/gin-gonic/gin/pull/1026)
- [NEW] Show query string in logger, see [#999](https://github.com/gin-gonic/gin/pull/999)
- [NEW] Add [`func (*Context) SecureJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.SecureJSON), see [#987](https://github.com/gin-gonic/gin/pull/987) and [#993](https://github.com/gin-gonic/gin/pull/993)
- [DEPRECATE] `func (*Context) GetCookie` for [`func (*Context) Cookie`](https://godoc.org/github.com/gin-gonic/gin#Context.Cookie)
- [FIX] Don't display color tags if [`func DisableConsoleColor`](https://godoc.org/github.com/gin-gonic/gin#DisableConsoleColor) called, see [#1072](https://github.com/gin-gonic/gin/pull/1072)
- [FIX] Gin Mode `""` when calling [`func Mode`](https://godoc.org/github.com/gin-gonic/gin#Mode) now returns `const DebugMode`, see [#1250](https://github.com/gin-gonic/gin/pull/1250)
- [FIX] `Flush()` now doesn't overwrite `responseWriter` status code, see [#1460](https://github.com/gin-gonic/gin/pull/1460)

### Gin 1.2.0

- [NEW] Switch from godeps to govendor
- [NEW] Add support for Let's Encrypt via gin-gonic/autotls
- [NEW] Improve README examples and add extra at examples folder
- [NEW] Improved support with App Engine
- [NEW] Add custom template delimiters, see #860
- [NEW] Add Template Func Maps, see #962
- [NEW] Add \*context.Handler(), see #928
- [NEW] Add \*context.GetRawData()
- [NEW] Add \*context.GetHeader() (request)
- [NEW] Add \*context.AbortWithStatusJSON() (JSON content type)
- [NEW] Add \*context.Keys type cast helpers
- [NEW] Add \*context.ShouldBindWith()
- [NEW] Add \*context.MustBindWith()
- [NEW] Add \*engine.SetFuncMap()
- [DEPRECATE] On next release: \*context.BindWith(), see #855
- [FIX] Refactor render
- [FIX] Reworked tests
- [FIX] logger now supports cygwin
- [FIX] Use X-Forwarded-For before X-Real-Ip
- [FIX] time.Time binding (#904)

### Gin 1.1.4

- [NEW] Support google appengine for IsTerminal func

### Gin 1.1.3

- [FIX] Reverted Logger: skip ANSI color commands

### Gin 1.1

- [NEW] Implement QueryArray and PostArray methods 
- [NEW] Refactor GetQuery and GetPostForm 
- [NEW] Add contribution guide 
- [FIX] Corrected typos in README
- [FIX] Removed additional Iota  
- [FIX] Changed imports to gopkg instead of github in README (#733) 
- [FIX] Logger: skip ANSI color commands if output is not a tty

### Gin 1.0rc2 (...)

- [PERFORMANCE] Fast path for writing Content-Type.
- [PERFORMANCE] Much faster 404 routing
- [PERFORMANCE] Allocation optimizations
- [PERFORMANCE] Faster root tree lookup
- [PERFORMANCE] Zero overhead, String() and JSON() rendering.
- [PERFORMANCE] Faster ClientIP parsing
- [PERFORMANCE] Much faster SSE implementation
- [NEW] Benchmarks suite
- [NEW] Bind validation can be disabled and replaced with custom validators.
- [NEW] More flexible HTML render
- [NEW] Multipart and PostForm bindings
- [NEW] Adds method to return all the registered routes
- [NEW] Context.HandlerName() returns the main handler's name
- [NEW] Adds Error.IsType() helper
- [FIX] Binding multipart form
- [FIX] Integration tests
- [FIX] Crash when binding non struct object in Context.
- [FIX] RunTLS() implementation
- [FIX] Logger() unit tests
- [FIX] Adds SetHTMLTemplate() warning
- [FIX] Context.IsAborted()
- [FIX] More unit tests
- [FIX] JSON, XML, HTML renders accept custom content-types
- [FIX] gin.AbortIndex is unexported
- [FIX] Better approach to avoid directory listing in StaticFS()
- [FIX] Context.ClientIP() always returns the IP with trimmed spaces.
- [FIX] Better warning when running in debug mode.
- [FIX] Google App Engine integration. debugPrint does not use os.Stdout
- [FIX] Fixes integer overflow in error type
- [FIX] Error implements the json.Marshaller interface
- [FIX] MIT license in every file


### Gin 1.0rc1 (May 22, 2015)

- [PERFORMANCE] Zero allocation router
- [PERFORMANCE] Faster JSON, XML and text rendering
- [PERFORMANCE] Custom hand optimized HttpRouter for Gin
- [PERFORMANCE] Misc code optimizations. Inlining, tail call optimizations
- [NEW] Built-in support for golang.org/x/net/context
- [NEW] Any(path, handler). Create a route that matches any path
- [NEW] Refactored rendering pipeline (faster and static typeded)
- [NEW] Refactored errors API
- [NEW] IndentedJSON() prints pretty JSON
- [NEW] Added gin.DefaultWriter
- [NEW] UNIX socket support
- [NEW] RouterGroup.BasePath is exposed
- [NEW] JSON validation using go-validate-yourself (very powerful options)
- [NEW] Completed suite of unit tests
- [NEW] HTTP streaming with c.Stream()
- [NEW] StaticFile() creates a router for serving just one file.
- [NEW] StaticFS() has an option to disable directory listing.
- [NEW] StaticFS() for serving static files through virtual filesystems
- [NEW] Server-Sent Events native support
- [NEW] WrapF() and WrapH() helpers for wrapping http.HandlerFunc and http.Handler
- [NEW] Added LoggerWithWriter() middleware
- [NEW] Added RecoveryWithWriter() middleware
- [NEW] Added DefaultPostFormValue()
- [NEW] Added DefaultFormValue()
- [NEW] Added DefaultParamValue()
- [FIX] BasicAuth() when using custom realm
- [FIX] Bug when serving static files in nested routing group
- [FIX] Redirect using built-in http.Redirect()
- [FIX] Logger when printing the requested path
- [FIX] Documentation typos
- [FIX] Context.Engine renamed to Context.engine
- [FIX] Better debugging messages
- [FIX] ErrorLogger
- [FIX] Debug HTTP render
- [FIX] Refactored binding and render modules
- [FIX] Refactored Context initialization
- [FIX] Refactored BasicAuth()
- [FIX] NoMethod/NoRoute handlers
- [FIX] Hijacking http
- [FIX] Better support for Google App Engine (using log instead of fmt)


### Gin 0.6 (Mar 9, 2015)

- [NEW] Support multipart/form-data
- [NEW] NoMethod handler
- [NEW] Validate sub structures
- [NEW] Support for HTTP Realm Auth
- [FIX] Unsigned integers in binding
- [FIX] Improve color logger


### Gin 0.5 (Feb 7, 2015)

- [NEW] Content Negotiation
- [FIX] Solved security bug that allow a client to spoof ip
- [FIX] Fix unexported/ignored fields in binding


### Gin 0.4 (Aug 21, 2014)

- [NEW] Development mode
- [NEW] Unit tests
- [NEW] Add Content.Redirect()
- [FIX] Deferring WriteHeader()
- [FIX] Improved documentation for model binding


### Gin 0.3 (Jul 18, 2014)

- [PERFORMANCE] Normal log and error log are printed in the same call.
- [PERFORMANCE] Improve performance of NoRouter()
- [PERFORMANCE] Improve context's memory locality, reduce CPU cache faults.
- [NEW] Flexible rendering API
- [NEW] Add Context.File()
- [NEW] Add shorcut RunTLS() for http.ListenAndServeTLS
- [FIX] Rename NotFound404() to NoRoute()
- [FIX] Errors in context are purged
- [FIX] Adds HEAD method in Static file serving
- [FIX] Refactors Static() file serving
- [FIX] Using keyed initialization to fix app-engine integration
- [FIX] Can't unmarshal JSON array, #63
- [FIX] Renaming Context.Req to Context.Request
- [FIX] Check application/x-www-form-urlencoded when parsing form


### Gin 0.2b (Jul 08, 2014)
- [PERFORMANCE] Using sync.Pool to allocatio/gc overhead
- [NEW] Travis CI integration
- [NEW] Completely new logger
- [NEW] New API for serving static files. gin.Static()
- [NEW] gin.H() can be serialized into XML
- [NEW] Typed errors. Errors can be typed. Internet/external/custom.
- [NEW] Support for Godeps
- [NEW] Travis/Godocs badges in README
- [NEW] New Bind() and BindWith() methods for parsing request body.
- [NEW] Add Content.Copy()
- [NEW] Add context.LastError()
- [NEW] Add shorcut for OPTIONS HTTP method
- [FIX] Tons of README fixes
- [FIX] Header is written before body
- [FIX] BasicAuth() and changes API a little bit
- [FIX] Recovery() middleware only prints panics
- [FIX] Context.Get() does not panic anymore. Use MustGet() instead.
- [FIX] Multiple http.WriteHeader() in NotFound handlers
- [FIX] Engine.Run() panics if http server can't be setted up
- [FIX] Crash when route path doesn't start with '/'
- [FIX] Do not update header when status code is negative
- [FIX] Setting response headers before calling WriteHeader in context.String()
- [FIX] Add MIT license
- [FIX] Changes behaviour of ErrorLogger() and Logger()
