package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data interface{}
}

func (r XML) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	return xml.NewEncoder(w).Encode(r.Data)
}
