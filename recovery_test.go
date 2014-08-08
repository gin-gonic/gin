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
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Disable panic logs for testing
	log.SetOutput(bytes.NewBuffer(nil))

	r := Default()
	r.GET("/", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})

	r.ServeHTTP(w, req)

	// restore logging
	log.SetOutput(os.Stderr)

	if w.Code != 500 {
		t.Errorf("Response code should be Internal Server Error, was: %s", w.Code)
	}
	bodyAsString := w.Body.String()

	//	fixme:
	if bodyAsString != "" {
		t.Errorf("Response body should be empty, was  %s", bodyAsString)
	}
	//fixme:
	if len(w.HeaderMap) != 0 {
		t.Errorf("No headers should be provided, was %s", w.HeaderMap)
	}

}
