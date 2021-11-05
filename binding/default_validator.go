// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validatorTags = make(map[reflect.Type]string)

// RegisterValidatorTag registers a validator tag against a number of types.
// This allows defining validation for custom slice, array, and map types. For example:
//     type CustomMap map[int]string
//     ...
//     binding.RegisterValidatorTag("gt=0", CustomMap{})
//
// Do not use the "dive" tag (unless in conjunction with "keys"/"endkeys").
// Slice/array/map elements are validated independently.
//
// This function will not have any effect is binding.Validator has been replaced.
//
// NOTE: This function is not thread-safe. It is intended that these all be registered prior to any validation.
func RegisterValidatorTag(tag string, types ...interface{}) {
	for _, typ := range types {
		t := reflect.TypeOf(typ)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Slice && t.Kind() != reflect.Array && t.Kind() != reflect.Map {
			panic("validator tags can be registered only for slices, arrays, and maps")
		}
		validatorTags[t] = tag
	}
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

// SliceFieldError is returned for invalid slice or array elements.
// It extends validator.FieldError with the index of the failing element.
type SliceFieldError interface {
	validator.FieldError
	Index() int
}

type sliceFieldError struct {
	validator.FieldError
	index int
}

func (fe sliceFieldError) Index() int {
	return fe.index
}

func (fe sliceFieldError) Error() string {
	return fmt.Sprintf("[%d]: %s", fe.index, fe.FieldError.Error())
}

func (fe sliceFieldError) Unwrap() error {
	return fe.FieldError
}

// MapFieldError is returned for invalid map values.
// It extends validator.FieldError with the key of the failing value.
type MapFieldError interface {
	validator.FieldError
	Key() interface{}
}

type mapFieldError struct {
	validator.FieldError
	key interface{}
}

func (fe mapFieldError) Key() interface{} {
	return fe.key
}

func (fe mapFieldError) Error() string {
	return fmt.Sprintf("[%v]: %s", fe.key, fe.FieldError.Error())
}

func (fe mapFieldError) Unwrap() error {
	return fe.FieldError
}

var _ StructValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but validates only structs, pointers, slices, arrays, and maps.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		var errs validator.ValidationErrors

		if tag, ok := validatorTags[value.Type()]; ok {
			if err := v.validateVar(obj, tag); err != nil {
				errs = append(errs, err.(validator.ValidationErrors)...) // nolint: errorlint
			}
		}

		count := value.Len()
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				for _, fieldError := range err.(validator.ValidationErrors) { // nolint: errorlint
					errs = append(errs, sliceFieldError{fieldError, i})
				}
			}
		}

		if len(errs) > 0 {
			return errs
		}
		return nil
	case reflect.Map:
		var errs validator.ValidationErrors

		if tag, ok := validatorTags[value.Type()]; ok {
			if err := v.validateVar(obj, tag); err != nil {
				errs = append(errs, err.(validator.ValidationErrors)...) // nolint: errorlint
			}
		}

		for _, key := range value.MapKeys() {
			if err := v.ValidateStruct(value.MapIndex(key).Interface()); err != nil {
				for _, fieldError := range err.(validator.ValidationErrors) { // nolint: errorlint
					errs = append(errs, mapFieldError{fieldError, key.Interface()})
				}
			}
		}
		if len(errs) > 0 {
			return errs
		}
		return nil
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	v.lazyinit()
	return v.validate.Struct(obj)
}

// validateStruct receives slice, array, and map types
func (v *defaultValidator) validateVar(obj interface{}, tag string) error {
	v.lazyinit()
	return v.validate.Var(obj, tag)
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://pkg.go.dev/github.com/go-playground/validator/v10
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}
