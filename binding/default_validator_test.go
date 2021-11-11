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

func TestRegisterValidatorTag(t *testing.T) {
	type CustomSlice []struct {
		A string
	}
	type CustomArray [10]struct {
		A string
	}
	type CustomMap map[string]struct {
		A string
	}
	type CustomStruct struct {
		A string
	}
	type CustomInt int

	// only slice, array, and map types are accepted
	RegisterValidatorTag("gt=0", CustomSlice{})
	RegisterValidatorTag("gt=0", &CustomSlice{})
	RegisterValidatorTag("gt=0", CustomArray{})
	RegisterValidatorTag("gt=0", &CustomArray{})
	RegisterValidatorTag("gt=0", CustomMap{})
	RegisterValidatorTag("gt=0", &CustomMap{})
	assert.Panics(t, func() { RegisterValidatorTag("gt=0", CustomStruct{}) })
	assert.Panics(t, func() { RegisterValidatorTag("gt=0", &CustomStruct{}) })
	assert.Panics(t, func() { var i CustomInt; RegisterValidatorTag("gt=0", i) })
	assert.Panics(t, func() { var i CustomInt; RegisterValidatorTag("gt=0", &i) })
}

func TestValidatorTagsSlice(t *testing.T) {
	type CustomSlice []struct {
		A string `binding:"max=8"`
	}

	var (
		invalidSlice    = CustomSlice{{"12345678"}}
		invalidVal      = CustomSlice{{"123456789"}, {"abcdefgh"}}
		validSlice      = CustomSlice{{"12345678"}, {"abcdefgh"}}
		invalidSliceVal = CustomSlice{{"123456789"}}
	)

	v := &defaultValidator{}

	// no tags registered for the slice itself yet, so only elements are validated
	assert.NoError(t, v.ValidateStruct(invalidSlice))
	assert.Error(t, v.ValidateStruct(invalidVal))
	assert.NoError(t, v.ValidateStruct(validSlice))
	assert.NoError(t, v.ValidateStruct(&invalidSlice))
	assert.Error(t, v.ValidateStruct(&invalidVal))
	assert.NoError(t, v.ValidateStruct(&validSlice))

	err := v.ValidateStruct(invalidSliceVal)
	assert.Error(t, err)
	assert.Len(t, err, 1) // only value error

	RegisterValidatorTag("gt=1", CustomSlice{})

	assert.Error(t, v.ValidateStruct(invalidSlice))
	assert.Error(t, v.ValidateStruct(invalidVal))
	assert.NoError(t, v.ValidateStruct(validSlice))
	assert.Error(t, v.ValidateStruct(&invalidSlice))
	assert.Error(t, v.ValidateStruct(&invalidVal))
	assert.NoError(t, v.ValidateStruct(&validSlice))

	err = v.ValidateStruct(invalidSliceVal)
	assert.Error(t, err)
	assert.Len(t, err, 2) // both slice length and value error
}

func TestValidatorTagsMap(t *testing.T) {
	type CustomMap map[string]struct {
		B int `binding:"gt=0"`
	}

	var (
		invalidMap    = CustomMap{"12345678": {1}}
		invalidKey    = CustomMap{"123456789": {1}, "abcdefgh": {2}}
		invalidVal    = CustomMap{"12345678": {0}, "abcdefgh": {2}}
		invalidMapVal = CustomMap{"12345678": {0}}
		validMap      = CustomMap{"12345678": {1}, "abcdefgh": {2}}
	)

	v := &defaultValidator{}

	// no tags registered for the map itself yet, so only values are validated
	assert.NoError(t, v.ValidateStruct(invalidMap))
	assert.NoError(t, v.ValidateStruct(invalidKey))
	assert.Error(t, v.ValidateStruct(invalidVal))
	assert.NoError(t, v.ValidateStruct(validMap))
	assert.NoError(t, v.ValidateStruct(&invalidMap))
	assert.NoError(t, v.ValidateStruct(&invalidKey))
	assert.Error(t, v.ValidateStruct(&invalidVal))
	assert.NoError(t, v.ValidateStruct(&validMap))

	err := v.ValidateStruct(invalidMapVal)
	assert.Error(t, err)
	assert.Len(t, err, 1) // only value error

	RegisterValidatorTag("gt=1,dive,keys,max=8,endkeys", CustomMap{})

	assert.Error(t, v.ValidateStruct(invalidMap))
	assert.Error(t, v.ValidateStruct(invalidKey))
	assert.Error(t, v.ValidateStruct(invalidVal))
	assert.NoError(t, v.ValidateStruct(validMap))
	assert.Error(t, v.ValidateStruct(&invalidMap))
	assert.Error(t, v.ValidateStruct(&invalidKey))
	assert.Error(t, v.ValidateStruct(&invalidVal))
	assert.NoError(t, v.ValidateStruct(&validMap))

	err = v.ValidateStruct(invalidMapVal)
	assert.Error(t, err)
	assert.Len(t, err, 2) // both map size and value errors
}
