package gin

import (
	"github.com/gin-gonic/gin/binding"
)

// DEPRECATED, use Bind() instead.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) EnsureBody(item interface{}) bool {
	return c.Bind(item)
}

// DEPRECATED use bindings directly
// Parses the body content as a JSON input. It decodes the json payload into the struct specified as a pointer.
func (c *Context) ParseBody(item interface{}) error {
	return binding.JSON.Bind(c.Req, item)
}
