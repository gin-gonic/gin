package render

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

type (
	Render interface {
		Render(http.ResponseWriter, int, ...interface{}) error
	}

	// JSON binding
	jsonRender struct{}

	// XML binding
	xmlRender struct{}

	// Plain text
	plainRender struct{}

	// form binding
	HTMLRender struct {
		Template *template.Template
	}
)

var (
	JSON  = jsonRender{}
	XML   = xmlRender{}
	Plain = plainRender{}
)

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func (_ jsonRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data[0])
}

func (_ xmlRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "application/xml")
	encoder := xml.NewEncoder(w)
	return encoder.Encode(data[0])
}

func (html HTMLRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/html")
	file := data[0].(string)
	obj := data[1]
	return html.Template.ExecuteTemplate(w, file, obj)
}

func (_ plainRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "text/plain")
	format := data[0].(string)
	args := data[1].([]interface{})
	var err error
	if len(args) > 0 {
		_, err = w.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = w.Write([]byte(format))
	}
	return err
}
