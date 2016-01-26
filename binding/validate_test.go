// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testInterface interface {
	String() string
}

type substruct_noValidation struct {
	I_String string
	I_Int    int
}

type mapNoValidationSub map[string]substruct_noValidation

type struct_noValidation_values struct {
	substruct_noValidation

	Boolean bool

	Uinteger   uint
	Integer    int
	Integer8   int8
	Integer16  int16
	Integer32  int32
	Integer64  int64
	Uinteger8  uint8
	Uinteger16 uint16
	Uinteger32 uint32
	Uinteger64 uint64

	Float32 float32
	Float64 float64

	String string

	Date time.Time

	Struct        substruct_noValidation
	InlinedStruct struct {
		String  []string
		Integer int
	}

	IntSlice           []int
	IntPointerSlice    []*int
	StructPointerSlice []*substruct_noValidation
	StructSlice        []substruct_noValidation
	InterfaceSlice     []testInterface

	UniversalInterface interface{}
	CustomInterface    testInterface

	FloatMap  map[string]float32
	StructMap mapNoValidationSub
}

func createNoValidation_values() struct_noValidation_values {
	integer := 1
	s := struct_noValidation_values{
		Boolean:            true,
		Uinteger:           1 << 29,
		Integer:            -10000,
		Integer8:           120,
		Integer16:          -20000,
		Integer32:          1 << 29,
		Integer64:          1 << 61,
		Uinteger8:          250,
		Uinteger16:         50000,
		Uinteger32:         1 << 31,
		Uinteger64:         1 << 62,
		Float32:            123.456,
		Float64:            123.456789,
		String:             "text",
		Date:               time.Time{},
		CustomInterface:    &bytes.Buffer{},
		Struct:             substruct_noValidation{},
		IntSlice:           []int{-3, -2, 1, 0, 1, 2, 3},
		IntPointerSlice:    []*int{&integer},
		StructSlice:        []substruct_noValidation{},
		UniversalInterface: 1.2,
		FloatMap: map[string]float32{
			"foo": 1.23,
			"bar": 232.323,
		},
		StructMap: mapNoValidationSub{
			"foo": substruct_noValidation{},
			"bar": substruct_noValidation{},
		},
		// StructPointerSlice []noValidationSub
		// InterfaceSlice     []testInterface
	}
	s.InlinedStruct.Integer = 1000
	s.InlinedStruct.String = []string{"first", "second"}
	s.I_String = "substring"
	s.I_Int = 987654
	return s
}

func TestValidateNoValidationValues(t *testing.T) {
	origin := createNoValidation_values()
	test := createNoValidation_values()
	empty := struct_noValidation_values{}

	assert.Nil(t, validate(test))
	assert.Nil(t, validate(&test))
	assert.Nil(t, validate(empty))
	assert.Nil(t, validate(&empty))

	assert.Equal(t, origin, test)
}

type struct_noValidation_pointer struct {
	substruct_noValidation

	Boolean bool

	Uinteger   *uint
	Integer    *int
	Integer8   *int8
	Integer16  *int16
	Integer32  *int32
	Integer64  *int64
	Uinteger8  *uint8
	Uinteger16 *uint16
	Uinteger32 *uint32
	Uinteger64 *uint64

	Float32 *float32
	Float64 *float64

	String *string

	Date *time.Time

	Struct *substruct_noValidation

	IntSlice           *[]int
	IntPointerSlice    *[]*int
	StructPointerSlice *[]*substruct_noValidation
	StructSlice        *[]substruct_noValidation
	InterfaceSlice     *[]testInterface

	FloatMap  *map[string]float32
	StructMap *mapNoValidationSub
}

func TestValidateNoValidationPointers(t *testing.T) {
	//origin := createNoValidation_values()
	//test := createNoValidation_values()
	empty := struct_noValidation_pointer{}

	//assert.Nil(t, validate(test))
	//assert.Nil(t, validate(&test))
	assert.Nil(t, validate(empty))
	assert.Nil(t, validate(&empty))

	//assert.Equal(t, origin, test)
}

type Object map[string]interface{}

func TestValidatePrimitives(t *testing.T) {
	obj := Object{"foo": "bar", "bar": 1}
	assert.NoError(t, validate(obj))
	assert.NoError(t, validate(&obj))
	assert.Equal(t, obj, Object{"foo": "bar", "bar": 1})

	obj2 := []Object{{"foo": "bar", "bar": 1}, {"foo": "bar", "bar": 1}}
	assert.NoError(t, validate(obj2))
	assert.NoError(t, validate(&obj2))

	nu := 10
	assert.NoError(t, validate(nu))
	assert.NoError(t, validate(&nu))
	assert.Equal(t, nu, 10)

	str := "value"
	assert.NoError(t, validate(str))
	assert.NoError(t, validate(&str))
	assert.Equal(t, str, "value")
}
