package render

import (
	"fmt"
	"net/http"
)

type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

func (r Redirect) Write(w http.ResponseWriter) error {
	if r.Code < 300 || r.Code > 308 {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}
