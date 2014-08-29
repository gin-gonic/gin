// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strconv"
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
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return Validate(obj)
}

func mapForm(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	formStruct := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		if inputFieldName := typeField.Tag.Get("form"); inputFieldName != "" {
			structField := formStruct.Field(i)
			if !structField.CanSet() {
				continue
			}

			inputValue, exists := form[inputFieldName]
			if !exists {
				continue
			}
			numElems := len(inputValue)
			if structField.Kind() == reflect.Slice && numElems > 0 {
				sliceOf := structField.Type().Elem().Kind()
				slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
				for i := 0; i < numElems; i++ {
					if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
						return err
					}
				}
				formStruct.Elem().Field(i).Set(slice)
			} else {
				if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return err
		} else {
			structField.SetInt(int64(intVal))
		}
	case reflect.Bool:
		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
	return nil
}

// Don't pass in pointers to bind to. Can lead to bugs. See:
// https://github.com/codegangsta/martini-contrib/issues/40
// https://github.com/codegangsta/martini-contrib/pull/34#issuecomment-29683659
func ensureNotPointer(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as binding models")
	}
}

func Validate(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// Allow ignored fields in the struct
			if field.Tag.Get("form") == "-" {
				continue
			}

			fieldValue := val.Field(i).Interface()
			zero := reflect.Zero(field.Type).Interface()

			if strings.Index(field.Tag.Get("binding"), "required") > -1 {
				fieldType := field.Type.Kind()
				if fieldType == reflect.Struct {
					err := Validate(fieldValue)
					if err != nil {
						return err
					}
				} else if reflect.DeepEqual(zero, fieldValue) {
					return errors.New("Required " + field.Name)
				} else if fieldType == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
					err := Validate(fieldValue)
					if err != nil {
						return err
					}
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fieldValue := val.Index(i).Interface()
			err := Validate(fieldValue)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	return nil
}
