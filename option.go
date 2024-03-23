package gin

import "net/http"

// OptionFunc defines the function to change the default configuration
type OptionFunc func(*Engine)

// Use attaches a global middleware to the router
func Use(middleware ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.Use(middleware...)
	}
}

// GET is a shortcut for RouterGroup.Handle("GET", path, handle)
func GET(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.GET(path, handlers...)
	}
}

// POST is a shortcut for RouterGroup.Handle("POST", path, handle)
func POST(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.POST(path, handlers...)
	}
}

// PUT is a shortcut for RouterGroup.Handle("PUT", path, handle)
func PUT(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.PUT(path, handlers...)
	}
}

// DELETE is a shortcut for RouterGroup.Handle("DELETE", path, handle)
func DELETE(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.DELETE(path, handlers...)
	}
}

// PATCH is a shortcut for RouterGroup.Handle("PATCH", path, handle)
func PATCH(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.PATCH(path, handlers...)
	}
}

// HEAD is a shortcut for RouterGroup.Handle("HEAD", path, handle)
func HEAD(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.HEAD(path, handlers...)
	}
}

// OPTIONS is a shortcut for RouterGroup.Handle("OPTIONS", path, handle)
func OPTIONS(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.OPTIONS(path, handlers...)
	}
}

// Any is a shortcut for RouterGroup.Handle("GET", path, handle)
func Any(path string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.Any(path, handlers...)
	}
}

// Group is used to create a new router group. You should add all the routes that have common middlewares or the same path prefix
func Group(path string, groupFunc func(*RouterGroup), handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		groupFunc(
			e.Group(path, handlers...),
		)
	}
}

// Route is a shortcut for RouterGroup.Handle
func Route(httpMethod, relativePath string, handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.Handle(httpMethod, relativePath, handlers...)
	}
}

// StaticFS returns a middleware that serves static files in the given file system
func StaticFS(path string, fs http.FileSystem) OptionFunc {
	return func(e *Engine) {
		e.StaticFS(path, fs)
	}
}

// StaticFile returns a middleware that serves a single file
func StaticFile(path, file string) OptionFunc {
	return func(e *Engine) {
		e.StaticFile(path, file)
	}
}

// Static returns a middleware that serves static files from a directory
func Static(path, root string) OptionFunc {
	return func(e *Engine) {
		e.Static(path, root)
	}
}

// NoRoute is a global handler for no matching routes
func NoRoute(handlers ...HandlerFunc) OptionFunc {
	return func(e *Engine) {
		e.NoRoute(handlers...)
	}
}
