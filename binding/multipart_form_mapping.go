// Copyright 2019 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"
)

type multipartRequest http.Request

var _ setter = (*multipartRequest)(nil)

var (
	// ErrMultiFileHeader multipart.FileHeader invalid
	ErrMultiFileHeader = errors.New("unsupported field type for multipart.FileHeader")

	// ErrMultiFileHeaderLenInvalid array for []*multipart.FileHeader len invalid
	ErrMultiFileHeaderLenInvalid = errors.New("unsupported len of array for []*multipart.FileHeader")
)

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *multipartRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (bool, error) {
	if files := r.MultipartForm.File[key]; len(files) != 0 {
		return setByMultipartFormFile(value, field, files)
	}

	return setByForm(value, field, r.MultipartForm.Value, key, opt)
}

func setByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(), len(files), len(files))
		isSet, err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSet {
			return isSet, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, ErrMultiFileHeader
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	if value.Len() != len(files) {
		return false, ErrMultiFileHeaderLenInvalid
	}
	for i := range files {
		set, err := setByMultipartFormFile(value.Index(i), field, files[i:i+1])
		if err != nil || !set {
			return set, err
		}
	}
	return true, nil
}
