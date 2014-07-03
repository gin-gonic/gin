package binding

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type (
	Binding interface {
		Bind(io.Reader, interface{}) error
	}

	// JSON binding
	jsonBinding struct{}

	// JSON binding
	xmlBinding struct{}
)

var (
	JSON = jsonBinding{}
	XML  = xmlBinding{}
)

func (_ jsonBinding) Bind(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&obj)
}

func (_ xmlBinding) Bind(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	return decoder.Decode(&obj)
}
