package gin

import(
	"testing"
	"html/template"
	"net/http"
	"net/http/httptest"
)

// TestRouterGroupGETRouteOK tests that GET route is correctly invoked.
func TestRouterGroupGETRouteOK(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	passed := false

	r := Default()
	r.GET("/test", func (c *Context) {
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
	r.GET("/test", func (c *Context) {
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
	r.POST("/test", func (c *Context) {
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
	r.DELETE("/test", func (c *Context) {
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
	r.PATCH("/test", func (c *Context) {
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
	r.PUT("/test", func (c *Context) {
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
	r.OPTIONS("/test", func (c *Context) {
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
	r.HEAD("/test", func (c *Context) {
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

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestContextParamsByName(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/alexandernyquist", nil)
	w := httptest.NewRecorder()
	name := ""

	r := Default()
	r.GET("/test/:name", func (c *Context) {
		name = c.Params.ByName("name")
	})

	r.ServeHTTP(w, req)

	if name != "alexandernyquist" {
		t.Errorf("Url parameter was not correctly parsed. Should be alexandernyquist, was %s.", name)
	}
}

// TestContextSetGet tests that a parameter is set correctly on the
// current context and can be retrieved using Get.
func TestContextSetGet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.GET("/test", func (c *Context) {
		// Key should be lazily created
		if c.Keys != nil {
			t.Error("Keys should be nil")
		}

		// Set
		c.Set("foo", "bar")

		if v := c.Get("foo"); v != "bar" {
			t.Errorf("Value should be bar, was %s", v)
		}
	})

	r.ServeHTTP(w, req)
}

// TestContextJSON tests that the response is serialized as JSON
// and Content-Type is set to application/json
func TestContextJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.GET("/test", func (c *Context) {
		c.JSON(200, H{"foo": "bar"})
	})

	r.ServeHTTP(w, req)

	if w.Body.String() != "{\"foo\":\"bar\"}\n" {
		t.Errorf("Response should be {\"foo\":\"bar\"}, was: %s", w.Body.String())
	}

	if w.HeaderMap.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type should be application/json, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestContextHTML tests that the response executes the templates
// and responds with Content-Type set to text/html
func TestContextHTML(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.HTMLTemplates = template.Must(template.New("t").Parse(`Hello {{.Name}}`))

	type TestData struct { Name string }

	r.GET("/test", func (c *Context) {
		c.HTML(200, "t", TestData{"alexandernyquist"})
	})

	r.ServeHTTP(w, req)

	if w.Body.String() != "Hello alexandernyquist" {
		t.Errorf("Response should be Hello alexandernyquist, was: %s", w.Body.String())
	}

	if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Errorf("Content-Type should be text/html, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

// TestContextString tests that the response is returned
// with Content-Type set to text/plain
func TestContextString(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r := Default()
	r.GET("/test", func (c *Context) {
		c.String(200, "test")
	})

	r.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Errorf("Response should be test, was: %s", w.Body.String())
	}

	if w.HeaderMap.Get("Content-Type") != "text/plain" {
		t.Errorf("Content-Type should be text/plain, was %s", w.HeaderMap.Get("Content-Type"))
	}
}