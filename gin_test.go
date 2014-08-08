package gin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

// TestRouterGroupGETRouteOK tests that GET route is correctly invoked.
func TestRouterGroupGETRouteOK(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.GET("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("GET route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupGETNoRootExistsRouteOK tests that a GET requse to root is correctly
// handled (404) when no root route exists.
func TestRouterGroupGETNoRootExistsRouteOK(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.GET("/test", func(c *Context) {
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		// If this fails, it's because httprouter needs to be updated to at least f78f58a0db
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
}

// TestRouterGroupPOSTRouteOK tests that POST route is correctly invoked.
func TestRouterGroupPOSTRouteOK(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.POST("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("POST route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupDELETERouteOK tests that DELETE route is correctly invoked.
func TestRouterGroupDELETERouteOK(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.DELETE("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("DELETE route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupPATCHRouteOK tests that PATCH route is correctly invoked.
func TestRouterGroupPATCHRouteOK(t *testing.T) {
	req, _ := http.NewRequest("PATCH", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.PATCH("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("PATCH route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupPUTRouteOK tests that PUT route is correctly invoked.
func TestRouterGroupPUTRouteOK(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.PUT("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("PUT route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupOPTIONSRouteOK tests that OPTIONS route is correctly invoked.
func TestRouterGroupOPTIONSRouteOK(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.OPTIONS("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("OPTIONS route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroupHEADRouteOK tests that HEAD route is correctly invoked.
func TestRouterGroupHEADRouteOK(t *testing.T) {
	req, _ := http.NewRequest("HEAD", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.HEAD("/test", func(c *Context) {
		passed = true
	})

	r.ServeHTTP(w, req)

	if passed == false {
		t.Errorf("HEAD route handler was not invoked.")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

// TestRouterGroup404 tests that 404 is returned for a route that does not exist.
func TestEngine404(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Response code should be %v, was %d", http.StatusNotFound, w.Code)
	}
}

// TestHandleStaticFile - ensure the static file handles properly
func TestHandleStaticFile(t *testing.T) {

	testRoot, _ := os.Getwd()

	f, err := ioutil.TempFile(testRoot, "")
	defer os.Remove(f.Name())

	if err != nil {
		t.Error(err)
	}

	filePath := path.Join("/", path.Base(f.Name()))
	req, _ := http.NewRequest("GET", filePath, nil)

	f.WriteString("Gin Web Framework")
	f.Close()

	w := httptest.NewRecorder()

	r := Default()
	r.Static("./", testRoot)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}

	if w.Body.String() != "Gin Web Framework" {
		t.Errorf("Response should be test, was: %s", w.Body.String())
	}

	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestHandleStaticDir(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()

	r := Default()
	r.Static("/", "./")

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}

	bodyAsString := w.Body.String()

	if len(bodyAsString) == 0 {
		t.Errorf("Got empty body instead of file tree")
	}

	if !strings.Contains(bodyAsString, "gin.go") {
		t.Errorf("Can't find:`gin.go` in file tree: %s", bodyAsString)
	}

	if w.HeaderMap.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
func TestHandleHeadToDir(t *testing.T) {

	req, _ := http.NewRequest("HEAD", "/", nil)

	w := httptest.NewRecorder()

	r := Default()
	r.Static("/", "./")

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}

	bodyAsString := w.Body.String()

	if len(bodyAsString) == 0 {
		t.Errorf("Got empty body instead of file tree")
	}
	if !strings.Contains(bodyAsString, "gin.go") {
		t.Errorf("Can't find:`gin.go` in file tree: %s", bodyAsString)
	}

	if w.HeaderMap.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}
