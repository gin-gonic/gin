package gin

import (
	"net/http"
)

func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext()
	c.reset()
	c.writermem.reset(w)
	return
}
