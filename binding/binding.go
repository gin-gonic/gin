package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

type (
	Binding interface {
		Bind(*http.Request, interface{}) error
	}

	// JSON binding
	jsonBinding struct{}

	// XML binding
	xmlBinding struct{}

	// // form binding
	formBinding struct{}
)

var (
	JSON = jsonBinding{}
	XML  = xmlBinding{}
	Form = formBinding{} // todo
)

func (_ jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err == nil {
		return Validate(obj)
	} else {
		return err
	}
}

func (_ xmlBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err == nil {
		return Validate(obj)
	} else {
		return err
	}
}

func (_ formBinding) Bind(req *http.Request, obj interface{}) error {
	return nil
}

func Validate(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue)) {
			if err := Validate(fieldValue); err != nil {
				return err
			}
		}

		if strings.Index(field.Tag.Get("binding"), "required") > -1 {
			if reflect.DeepEqual(zero, fieldValue) {
				name := field.Name
				if j := field.Tag.Get("json"); j != "" {
					name = j
				} else if f := field.Tag.Get("form"); f != "" {
					name = f
				}
				return errors.New("Required " + name)
			}
		}
	}
	return nil
}
