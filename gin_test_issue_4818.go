package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMethodNotAllowedSkippedNodesPanic(t *testing.T) {
	SetMode(ReleaseMode)
	router := New()
	router.HandleMethodNotAllowed = true

	h := func(c *Context) {}
	router.OPTIONS("/:p0/:p1/a/:p2", h)
	router.GET("/:p0/:p1/a/:p2", h)
	router.PATCH("/b/:p0/:p1/c", h)
	router.DELETE("/b/:p0/:p1/d/:p3", h)
	router.GET("/b/:p0/:p1/e/f", h)
	router.POST("/b/:p0/:p1/g/:p4/h", h)
	router.OPTIONS("/b/:p0/:p1/g/:p4/h", h)
	router.DELETE("/b/cache", h)
	router.GET("/b/clients/:p1/g", h)
	router.POST("/b/clients/:p1/g", h)
	router.PATCH("/b/clients/:p1/g/:p4", h)
	router.OPTIONS("/b/clients/:p1/g/:p4", h)

	req := httptest.NewRequest(http.MethodPost, "/b/clients/42", nil)
	w := httptest.NewRecorder()

	// This should return 405 Method Not Allowed, not panic
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
