// +build go1.7

package gin

import (
	"context"
	"net/http"
)

const ParamsKey = "_gin-gonic/gin/paramskey"

// WithParams is a helper function to add Params in native context
// Returns a http request
func WithParams(r *http.Request, params Params) *http.Request {
	ctx := context.WithValue(r.Context(), ParamsKey, params)
	return r.WithContext(ctx)
}

// GetParams is a helper function to get Params in native context
// Returns a Gin Params
func GetParams(r *http.Request) Params {
	if params := r.Context().Value(ParamsKey); params != nil {
		return params.(Params)
	}
	return nil
}
