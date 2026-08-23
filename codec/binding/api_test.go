package bindingcodec

import (
	"testing"
)

// TestAPIInitialization tests that the API is properly initialized
func TestAPIInitialization(t *testing.T) {
	if API == nil {
		t.Fatal("API should not be nil after initialization")
	}
}
