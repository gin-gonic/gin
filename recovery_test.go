package gin

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/recovery", nil)
	w := httptest.NewRecorder()

	// Disable panic logs for testing
	log.SetOutput(bytes.NewBuffer(nil))

	r := Default()
	r.GET("/recovery", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})

	r.ServeHTTP(w, req)

	// restore logging
	log.SetOutput(os.Stderr)

	if w.Code != 500 {
		t.Errorf("Response code should be Internal Server Error, was: %s", w.Code)
	}
	bodyAsString := w.Body.String()

	//fixme: no message provided?
	if bodyAsString != "" {
		t.Errorf("Response body should be empty, was  %s", bodyAsString)
	}
	//fixme:
	if len(w.HeaderMap) != 0 {
		t.Errorf("No headers should be provided, was %s", w.HeaderMap)
	}

}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	req, _ := http.NewRequest("GET", "/recovery", nil)
	w := httptest.NewRecorder()

	// Disable panic logs for testing
	log.SetOutput(bytes.NewBuffer(nil))

	r := Default()
	r.GET("/recovery", func(c *Context) {
		c.Abort(400)
		panic("Oupps, Houston, we have a problem")
	})

	r.ServeHTTP(w, req)

	// restore logging
	log.SetOutput(os.Stderr)

	// fixme: why not 500?
	if w.Code != 400 {
		t.Errorf("Response code should be Bad request, was: %s", w.Code)
	}
	bodyAsString := w.Body.String()

	//fixme: no message provided?
	if bodyAsString != "" {
		t.Errorf("Response body should be empty, was  %s", bodyAsString)
	}
	//fixme:
	if len(w.HeaderMap) != 0 {
		t.Errorf("No headers should be provided, was %s", w.HeaderMap)
	}

}
