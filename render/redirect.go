package render

import (
	"fmt"
	"net/http"
)

type redirectRender struct{}

func (_ redirectRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	req := data[0].(*http.Request)
	location := data[1].(string)
	WriteRedirect(w, code, req, location)
	return nil
}

func WriteRedirect(w http.ResponseWriter, code int, req *http.Request, location string) {
	if code < 300 || code > 308 {
		panic(fmt.Sprintf("Cannot redirect with status code %d", code))
	}
	http.Redirect(w, req, location, code)
}
