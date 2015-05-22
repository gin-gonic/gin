// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	ErrorTypeBind    = 1 << 30
	ErrorTypeRender  = 1 << 29
	ErrorTypePrivate = 1 << 0
	ErrorTypePublic  = 1 << 1

	ErrorTypeAny = 0xfffffff
	ErrorTypeNu  = 2
)

type (
	Error struct {
		Err  error       `json:"error"`
		Type int         `json:"-"`
		Meta interface{} `json:"meta"`
	}

	errorMsgs []*Error
)

var _ error = &Error{}

func (msg *Error) SetType(flags int) *Error {
	msg.Type = flags
	return msg
}

func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

func (msg *Error) JSON() interface{} {
	json := H{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				json[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			json["meta"] = msg.Meta
		}
	}
	if _, ok := json["error"]; !ok {
		json["error"] = msg.Error()
	}
	return json
}

func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

func (msg *Error) Error() string {
	return msg.Err.Error()
}

func (a errorMsgs) ByType(typ int) errorMsgs {
	if len(a) == 0 {
		return a
	}
	result := make(errorMsgs, 0, len(a))
	for _, msg := range a {
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}
	return result
}

func (a errorMsgs) Last() *Error {
	length := len(a)
	if length == 0 {
		return nil
	}
	return a[length-1]
}

func (a errorMsgs) Errors() []string {
	if len(a) == 0 {
		return []string{}
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

func (a errorMsgs) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

func (a errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", (i + 1), msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
