// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

var intBitSize = map[reflect.Kind]int{
	reflect.Int:   0,
	reflect.Int8:  8,
	reflect.Int16: 16,
	reflect.Int32: 32,
	reflect.Int64: 64,
}

var uintBitSize = map[reflect.Kind]int{
	reflect.Uint:   0,
	reflect.Uint8:  8,
	reflect.Uint16: 16,
	reflect.Uint32: 32,
	reflect.Uint64: 64,
}

var floatBitSize = map[reflect.Kind]int{
	reflect.Float32: 32,
	reflect.Float64: 64,
}

var durationType = reflect.TypeOf(time.Duration(0))
var timeType = reflect.TypeOf(time.Time{})

func setTime(value string, val reflect.Value) (err error) {
	var t time.Time
	if t, err = time.ParseInLocation(time.RFC3339, value, time.Local); err != nil {
		return err
	}
	val.Set(reflect.ValueOf(t))
	return
}

func parseBaseTypeVar(value string, ptr reflect.Value) (err error) {
	val := ptr.Elem()
	switch val.Kind() {
	// bool
	case reflect.Bool:
		return setBoolField(value, val)
		// string
	case reflect.String:
		val.SetString(value)
		return
	case reflect.Int64:
		if val.Type() == durationType {
			return setTimeDuration(value, val, reflect.StructField{})
		}

	case reflect.Struct:
		if val.Type() == timeType {
			return setTime(value, val)
		}
	}

	// int, int8, int16, int32, int64
	if bs, ok := intBitSize[val.Kind()]; ok {
		return setIntField(value, bs, val)
	}

	// uint, uint8, uint16, uint32, uint64
	if bs, ok := uintBitSize[val.Kind()]; ok {
		return setUintField(value, bs, val)
	}

	// float32 float64
	if bs, ok := floatBitSize[val.Kind()]; ok {
		return setFloatField(value, bs, val)
	}

	return nil
}

func setSlice2(values []string, ptr reflect.Value) error {
	slice := reflect.MakeSlice(ptr.Elem().Type(), len(values), len(values))
	ptr.Elem().Set(slice)
	if err := setArray2(values, ptr); err != nil {
		return err
	}

	ptr.Elem().Set(slice)
	return nil

}

func setArray2(values []string, ptr reflect.Value) error {
	if ptr.Elem().Len() != len(values) {
		return fmt.Errorf("Unequal length:%d:%d", ptr.Elem().Len(), len(values))
	}

	for i, v := range values {
		if err := parseBaseTypeVar(v, ptr.Elem().Index(i).Addr()); err != nil {
			return err
		}
	}
	return nil
}

// slice, array
// base type
func parseTypeVar(ptr reflect.Value, values []string) error {
	switch ptr.Elem().Kind() {
	case reflect.Slice:
		return setSlice2(values, ptr)
	case reflect.Array:
		return setArray2(values, ptr)
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Struct, reflect.UnsafePointer:

		if ptr.Elem().Type() == timeType {
			return parseBaseTypeVar(values[0], ptr)
		}
		return errors.New("Unsupported type")
	default:
		return parseBaseTypeVar(values[0], ptr)
	}
}

func SetValue(ptr, defaultVal reflect.Value, values []string, ok bool) error {
	if ok {
		return parseTypeVar(ptr, values)
	}

	ptr.Elem().Set(defaultVal)
	return nil
}
