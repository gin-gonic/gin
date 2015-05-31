// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type struct1 struct {
	Value float64 `binding:"required"`
}

type struct2 struct {
	RequiredValue string `binding:"required"`
	Value         float64
}

type struct3 struct {
	Integer    int
	String     string
	BasicSlice []int
	Boolean    bool

	RequiredInteger       int       `binding:"required"`
	RequiredString        string    `binding:"required"`
	RequiredAnotherStruct struct1   `binding:"required"`
	RequiredBasicSlice    []int     `binding:"required"`
	RequiredComplexSlice  []struct2 `binding:"required"`
	RequiredBoolean       bool      `binding:"required"`
}

func createStruct() struct3 {
	return struct3{
		RequiredInteger:       2,
		RequiredString:        "hello",
		RequiredAnotherStruct: struct1{1.5},
		RequiredBasicSlice:    []int{1, 2, 3, 4},
		RequiredComplexSlice: []struct2{
			{RequiredValue: "A"},
			{RequiredValue: "B"},
		},
		RequiredBoolean: true,
	}
}

func TestValidateGoodObject(t *testing.T) {
	test := createStruct()
	assert.Nil(t, validate(&test))
}

type Object map[string]interface{}
type MyObjects []Object

func TestValidateSlice(t *testing.T) {
	var obj MyObjects
	var obj2 Object
	var nu = 10

	assert.NoError(t, validate(obj))
	assert.NoError(t, validate(&obj))
	assert.NoError(t, validate(obj2))
	assert.NoError(t, validate(&obj2))
	assert.NoError(t, validate(nu))
	assert.NoError(t, validate(&nu))
}
