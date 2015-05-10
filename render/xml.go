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
	writeHeader(w, code, "application/xml; charset=utf-8")
	return xml.NewEncoder(w).Encode(data)
}
