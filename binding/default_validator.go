// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"fmt"
	"sync"
	"strings"
	"log"

	"github.com/go-playground/validator/v10"
	// "github.com/ahmetb/go-linq"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

type sliceValidateError []error 

func (err sliceValidateError) Error() string {
	var errMsgs []string
	for i, e := range err {
		if e == nil {
			continue
		}
		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, err.Error()))
	}
	return strings.Join(errMsgs, "\n")
}

var _ ValidatorImp = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) Validate(obj interface{}) error {
	log.Println("called")
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	log.Printf("valueType: %v, %#v\n", valueType, valueType)
	if valueType == reflect.Ptr {
		value = value.Elem()
		valueType = value.Kind()
	}

	switch valueType {	
	case reflect.Struct:
		log.Printf("goto validateStruct: %v, %#v\n", obj, obj)
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)
		for i := 0; i < count; i++ {
			log.Println("called inside slice")
			if err := v.Validate(value.Index(i)); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		log.Println(validateRet)
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet		
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	v.lazyinit()
	return v.validate.Struct(obj)
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
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
