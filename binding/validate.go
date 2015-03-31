// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"reflect"
	"strings"
)

func Validate(obj interface{}) error {
	return validate(obj, "{{ROOT}}")
}

func validate(obj interface{}, parent string) error {
	typ, val := inspectObject(obj)
	switch typ.Kind() {
	case reflect.Struct:
		return validateStruct(typ, val, parent)

	case reflect.Slice:
		return validateSlice(typ, val, parent)

	default:
		return errors.New("The object is not a slice or struct.")
	}
}

func inspectObject(obj interface{}) (typ reflect.Type, val reflect.Value) {
	typ = reflect.TypeOf(obj)
	val = reflect.ValueOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	return
}

func validateSlice(typ reflect.Type, val reflect.Value, parent string) error {
	if typ.Elem().Kind() == reflect.Struct {
		for i := 0; i < val.Len(); i++ {
			itemValue := val.Index(i).Interface()
			if err := validate(itemValue, parent); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateStruct(typ reflect.Type, val reflect.Value, parent string) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Allow ignored and unexported fields in the struct
		// TODO should include  || field.Tag.Get("form") == "-"
		if len(field.PkgPath) > 0 {
			continue
		}

		fieldValue := val.Field(i).Interface()
		requiredField := strings.Index(field.Tag.Get("binding"), "required") > -1

		if requiredField {
			zero := reflect.Zero(field.Type).Interface()
			if reflect.DeepEqual(zero, fieldValue) {
				return errors.New("Required " + field.Name + " in " + parent)
			}
		}
		fieldType := field.Type.Kind()
		if fieldType == reflect.Struct || fieldType == reflect.Slice {
			if err := validate(fieldValue, field.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
