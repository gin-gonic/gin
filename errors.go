// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/internal/json"
)

// ErrorType is an unsigned 64-bit error code as defined in the gin spec.
type ErrorType uint64

const (
	ErrorTypeBind      ErrorType = 1 << 63
	ErrorTypeRender    ErrorType = 1 << 62
	ErrorTypePrivate   ErrorType = 1 << 0
	ErrorTypePublic    ErrorType = 1 << 1
	ErrorTypeAny       ErrorType = 1<<64 - 1
	ErrorTypeNu        = 2
)

// Error represents an error's specification.
type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

// errorMsgs is a slice of errors.
type errorMsgs []*Error

// SetType sets the error's type.
func (err *Error) SetType(flags ErrorType) *Error {
	err.Type = flags
	return err
}

// SetMeta sets the error's meta data.
func (err *Error) SetMeta(data interface{}) *Error {
	err.Meta = data
	return err
}

// JSON creates a properly formatted JSON.
func (err *Error) JSON() interface{} {
	jsonData := map[string]interface{}{
		"error": err.Err.Error(),
	}

	if err.Meta != nil {
		value := reflect.ValueOf(err.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return err.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = err.Meta
		}
	}

	return jsonData
}

// MarshalJSON implements the json.Marshaller interface.
func (err *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(err.JSON())
}

// Error implements the error interface.
func (err Error) Error() string {
	return fmt.Sprintf("error: %v", err.Err)
}

// IsType checks if the error has a specific type.
func (err *Error) IsType(flags ErrorType) bool {
	return (err.Type & flags) > 0
}

// Unwrap returns the wrapped error.
func (err *Error) Unwrap() error {
	return err.Err
}

// ByType returns a readonly copy filtered by the error type.
func (errs errorMsgs) ByType(typ ErrorType) errorMsgs {
	if len(errs) == 0 || typ == ErrorTypeAny {
		return errs
	}

	var result errorMsgs
	for _, err := range errs {
		if err.IsType(typ) {
			result = append(result, err)
		}
	}
	return result
}

// Last returns the last error in the slice.
func (errs errorMsgs) Last() *Error {
	if length := len(errs); length > 0 {
		return errs[length-1]
	}
	return nil
}

// Errors returns an array with all the error messages.
func (errs errorMsgs) Errors() []string {
	errorStrings := make([]string, len(errs))
	for i, err := range errs {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

// JSON creates a JSON representation of the error slice.
func (errs errorMsgs) JSON() interface{} {
	switch length := len(errs); length {
	case 0:
		return nil
	case 1:
		return errs.Last().JSON()
	default:
		jsonData := make([]interface{}, length)
		for i, err := range errs {
			jsonData[i] = err.JSON()
		}
		return jsonData
	}
}

// MarshalJSON implements the json.Marshaller interface for the error slice.
func (errs errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(errs.JSON())
}

// String returns a formatted string representation of the error slice.
func (errs errorMsgs) String() string {
	var buffer strings.Builder
	for i, err := range errs {
		fmt.Fprintf(&buffer, "Error #%02d: %v\n", i+1, err)
		if err.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", err.Meta)
		}
	}
	return buffer.String()
}
