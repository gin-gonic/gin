package render

import (
	"fmt"
	"net/http"
)

type String struct {
	Format string
	Data   []interface{}
}

func (r String) Write(w http.ResponseWriter) error {
	header := w.Header()
	if _, exist := header["Content-Type"]; !exist {
		header.Set("Content-Type", "text/plain; charset=utf-8")
	}
	if len(r.Data) > 0 {
		fmt.Fprintf(w, r.Format, r.Data...)
	} else {
		w.Write([]byte(r.Format))
	}
	return nil
}
