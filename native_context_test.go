// +build go1.7

package gin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParamsWithWrap(t *testing.T) {
	router := New()
	router.GET("/hello/:name", WrapH(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/hello/gopher", req.URL.Path)
		assert.Equal(t, "gopher", GetParams(req).ByName("name"))
	})))

	router.GET("/hello2/:name", WrapF(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/hello2/gopher", req.URL.Path)
		assert.Equal(t, "gopher", GetParams(req).ByName("name"))
	}))

	router.GET("/hello", WrapH(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/hello", req.URL.Path)
		assert.Equal(t, "", GetParams(req).ByName("name"))
	})))

	w := performRequest(router, "GET", "/hello/gopher")
	assert.Equal(t, 200, w.Code)

	w = performRequest(router, "GET", "/hello2/gopher")
	assert.Equal(t, 200, w.Code)

	w = performRequest(router, "GET", "/hello")
	assert.Equal(t, 200, w.Code)
}

func TestGetParamsWithRequest(t *testing.T) {
	req := &http.Request{}
	assert.Equal(t, "", GetParams(req).ByName("name"))
}
