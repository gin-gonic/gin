// +build !go1.7

package gin

import "net/http"

// WithParams is a helper function to add Params in native context
// Returns a http request
func WithParams(r *http.Request, params Params) *http.Request {
	return r
}

// GetParams is a helper function to get Params in native context
// Returns a Gin Params
func GetParams(r *http.Request) Params {
	return nil
}
