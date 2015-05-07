package render

import (
	"encoding/xml"
	"net/http"
)

type xmlRender struct{}

func (_ xmlRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	return WriteXML(w, code, data[0])
}

func WriteXML(w http.ResponseWriter, code int, data interface{}) error {
	WriteHeader(w, code, "application/xml")
	return xml.NewEncoder(w).Encode(data)
}
