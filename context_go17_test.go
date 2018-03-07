// +build go1.7

package gin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type interceptedWriter struct {
	ResponseWriter
	b *bytes.Buffer
}

func (i interceptedWriter) WriteHeader(code int) {
	i.Header().Del("X-Test")
	i.ResponseWriter.WriteHeader(code)
}
func TestInterceptedHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := CreateTestContext(w)

	r.Use(func(c *Context) {
		i := interceptedWriter{
			ResponseWriter: c.Writer,
			b:              bytes.NewBuffer(nil),
		}
		c.Writer = i
		c.Next()
		c.Header("X-Test", "overridden")
		c.Writer = i.ResponseWriter
	})
	r.GET("/", func(c *Context) {
		c.Header("X-Test", "original")
		c.Header("X-Test-2", "present")
		c.String(http.StatusOK, "hello world")
	})
	c.Request = httptest.NewRequest("GET", "/", nil)
	r.HandleContext(c)
	// Result() has headers frozen when WriteHeaderNow() has been called
	// Compared to this time, this is when the response headers will be flushed
	// As response is flushed on c.String, the Header cannot be set by the first
	// middleware. Assert this
	assert.Equal(t, "", w.Result().Header.Get("X-Test"))
	assert.Equal(t, "present", w.Result().Header.Get("X-Test-2"))
}
