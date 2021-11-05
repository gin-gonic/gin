// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestSliceFieldError(t *testing.T) {
	var fe validator.FieldError = dummyFieldError{msg: "test error"}

	var err SliceFieldError = sliceFieldError{fe, 10}
	assert.Equal(t, 10, err.Index())
	assert.Equal(t, "[10]: test error", err.Error())
	assert.Equal(t, fe, errors.Unwrap(err))
}

func TestMapFieldError(t *testing.T) {
	var fe validator.FieldError = dummyFieldError{msg: "test error"}

	var err MapFieldError = mapFieldError{fe, "test key"}
	assert.Equal(t, "test key", err.Key())
	assert.Equal(t, "[test key]: test error", err.Error())
	assert.Equal(t, fe, errors.Unwrap(err))

	err = mapFieldError{fe, 123}
	assert.Equal(t, 123, err.Key())
	assert.Equal(t, "[123]: test error", err.Error())
	assert.Equal(t, fe, errors.Unwrap(err))
}

type dummyFieldError struct {
	validator.FieldError
	msg string
}

func (fe dummyFieldError) Error() string {
	return fe.msg
}

func TestDefaultValidator(t *testing.T) {
	type exampleStruct struct {
		A string `binding:"max=8"`
		B int    `binding:"gt=0"`
	}
	tests := []struct {
		name    string
		v       *defaultValidator
		obj     interface{}
		wantErr bool
	}{
		{"validate nil obj", &defaultValidator{}, nil, false},
		{"validate int obj", &defaultValidator{}, 3, false},
		{"validate struct failed-1", &defaultValidator{}, exampleStruct{A: "123456789", B: 1}, true},
		{"validate struct failed-2", &defaultValidator{}, exampleStruct{A: "12345678", B: 0}, true},
		{"validate struct passed", &defaultValidator{}, exampleStruct{A: "12345678", B: 1}, false},
		{"validate *struct failed-1", &defaultValidator{}, &exampleStruct{A: "123456789", B: 1}, true},
		{"validate *struct failed-2", &defaultValidator{}, &exampleStruct{A: "12345678", B: 0}, true},
		{"validate *struct passed", &defaultValidator{}, &exampleStruct{A: "12345678", B: 1}, false},
		{"validate []struct failed-1", &defaultValidator{}, []exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []struct failed-2", &defaultValidator{}, []exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []struct passed", &defaultValidator{}, []exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate []*struct failed-1", &defaultValidator{}, []*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []*struct failed-2", &defaultValidator{}, []*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []*struct passed", &defaultValidator{}, []*exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]struct failed-1", &defaultValidator{}, &[]exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]struct failed-2", &defaultValidator{}, &[]exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]struct passed", &defaultValidator{}, &[]exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]*struct failed-1", &defaultValidator{}, &[]*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]*struct failed-2", &defaultValidator{}, &[]*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]*struct passed", &defaultValidator{}, &[]*exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate map[string]struct failed-1", &defaultValidator{}, map[string]exampleStruct{"x": {A: "123456789", B: 1}}, true},
		{"validate map[string]struct failed-2", &defaultValidator{}, map[string]exampleStruct{"x": {A: "12345678", B: 0}}, true},
		{"validate map[string]struct passed", &defaultValidator{}, map[string]exampleStruct{"x": {A: "12345678", B: 1}}, false},
		{"validate map[string]*struct failed-1", &defaultValidator{}, map[string]*exampleStruct{"x": {A: "123456789", B: 1}}, true},
		{"validate map[string]*struct failed-2", &defaultValidator{}, map[string]*exampleStruct{"x": {A: "12345678", B: 0}}, true},
		{"validate map[string]*struct passed", &defaultValidator{}, map[string]*exampleStruct{"x": {A: "12345678", B: 1}}, false},
		{"validate *map[string]struct failed-1", &defaultValidator{}, &map[string]exampleStruct{"x": {A: "123456789", B: 1}}, true},
		{"validate *map[string]struct failed-2", &defaultValidator{}, &map[string]exampleStruct{"x": {A: "12345678", B: 0}}, true},
		{"validate *map[string]struct passed", &defaultValidator{}, &map[string]exampleStruct{"x": {A: "12345678", B: 1}}, false},
		{"validate *map[string]*struct failed-1", &defaultValidator{}, &map[string]*exampleStruct{"x": {A: "123456789", B: 1}}, true},
		{"validate *map[string]*struct failed-2", &defaultValidator{}, &map[string]*exampleStruct{"x": {A: "12345678", B: 0}}, true},
		{"validate *map[string]*struct passed", &defaultValidator{}, &map[string]*exampleStruct{"x": {A: "12345678", B: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.ValidateStruct(tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("defaultValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
