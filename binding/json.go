// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	stdjson "encoding/json"

	"github.com/gin-gonic/gin/codec/json"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// any as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj any) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Fast path: standard decode when the struct has no time_format tags.
	if !hasTimeFormatTags(obj) {
		decoder := json.API.NewDecoder(bytes.NewReader(body))
		if EnableDecoderUseNumber {
			decoder.UseNumber()
		}
		if EnableDecoderDisallowUnknownFields {
			decoder.DisallowUnknownFields()
		}
		if err := decoder.Decode(obj); err != nil {
			return err
		}
		return validate(obj)
	}

	// Slow path: decode field-by-field so that time.Time fields with a
	// time_format tag are parsed by setTimeField instead of the standard
	// JSON decoder (which only understands RFC3339).
	if err := decodeJSONWithTimeFormat(body, obj); err != nil {
		return err
	}
	return validate(obj)
}

// ---------------------------------------------------------------------------
// time_format tag support
// ---------------------------------------------------------------------------

// jsonKey returns the JSON object key for a struct field, respecting the
// "json" struct tag.  It returns the field name when no tag is present.
func jsonKey(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("json")
	if !ok {
		return field.Name
	}
	if idx := strings.IndexByte(tag, ','); idx >= 0 {
		tag = tag[:idx]
	}
	if tag == "" || tag == "-" {
		return field.Name
	}
	return tag
}

// hasTimeFormatTags returns true when the value (or any nested struct reachable
// through exported fields) contains at least one time.Time field with a
// non-empty time_format struct tag.
func hasTimeFormatTags(obj any) bool {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return hasTimeFormatTagsValue(rv)
}

func hasTimeFormatTagsValue(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Tag.Get("time_format") != "" {
			return true
		}
		ft := field.Type
		if ft == reflect.TypeOf(time.Time{}) {
			continue
		}
		// Recurse into nested structs and pointer-to-structs.
		switch ft.Kind() {
		case reflect.Struct:
			if hasTimeFormatTagsValue(reflect.New(ft).Elem()) {
				return true
			}
		case reflect.Ptr:
			if ft.Elem().Kind() == reflect.Struct &&
				hasTimeFormatTagsValue(reflect.New(ft.Elem()).Elem()) {
				return true
			}
		}
	}
	return false
}

// decodeJSONWithTimeFormat decodes a JSON body into a struct, handling
// time.Time fields that carry a time_format tag by parsing them with
// setTimeField instead of the standard JSON decoder.
func decodeJSONWithTimeFormat(body []byte, obj any) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return decodeJSONValue(body, rv)
}

func decodeJSONValue(body []byte, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		// Unmarshal the raw body into a map of raw JSON messages so we can
		// process each field individually.
		var rawMap map[string]stdjson.RawMessage
		if err := stdjson.Unmarshal(body, &rawMap); err != nil {
			return err
		}
		t := rv.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			key := jsonKey(field)
			raw, ok := rawMap[key]
			if !ok {
				continue
			}
			fv := rv.Field(i)

			if err := decodeField(raw, field, fv); err != nil {
				return err
			}
		}
		return nil

	case reflect.Slice, reflect.Array:
		var items []stdjson.RawMessage
		if err := stdjson.Unmarshal(body, &items); err != nil {
			return err
		}
		for i := 0; i < rv.Len() && i < len(items); i++ {
			if err := decodeJSONValue(items[i], rv.Index(i)); err != nil {
				return err
			}
		}
		return nil

	default:
		return json.API.Unmarshal(body, rv.Addr().Interface())
	}
}

func decodeField(raw stdjson.RawMessage, field reflect.StructField, fv reflect.Value) error {
	ft := field.Type

	// time.Time with time_format tag → use setTimeField.
	if ft == reflect.TypeOf(time.Time{}) && field.Tag.Get("time_format") != "" {
		rawStr, err := rawToStr(raw)
		if err != nil {
			return nil // skip unparseable values
		}
		return setTimeField(rawStr, field, fv)
	}

	// *time.Time with time_format tag → allocate, parse, set pointer.
	if ft.Kind() == reflect.Ptr && ft.Elem() == reflect.TypeOf(time.Time{}) &&
		field.Tag.Get("time_format") != "" {
		rawStr, err := rawToStr(raw)
		if err != nil {
			return nil
		}
		newTime := reflect.New(reflect.TypeOf(time.Time{})).Elem()
		if err := setTimeField(rawStr, field, newTime); err != nil {
			return err
		}
		fv.Set(newTime.Addr())
		return nil
	}

	// Nested struct → recurse.
	if ft.Kind() == reflect.Struct && ft != reflect.TypeOf(time.Time{}) {
		return decodeJSONValue(raw, fv)
	}

	// Pointer to nested struct → allocate if nil, then recurse.
	if ft.Kind() == reflect.Ptr && ft.Elem().Kind() == reflect.Struct && ft.Elem() != reflect.TypeOf(time.Time{}) {
		if fv.IsNil() {
			fv.Set(reflect.New(ft.Elem()))
		}
		return decodeJSONValue(raw, fv.Elem())
	}

	// All other types: standard JSON unmarshal.
	return json.API.Unmarshal(raw, fv.Addr().Interface())
}

// rawToStr extracts a string value from a raw JSON message, handling both
// quoted strings and numeric values (unix timestamps).
func rawToStr(raw stdjson.RawMessage) (string, error) {
	var s string
	if err := stdjson.Unmarshal(raw, &s); err == nil {
		return s, nil
	}
	var n int64
	if err := stdjson.Unmarshal(raw, &n); err == nil {
		return strconv.FormatInt(n, 10), nil
	}
	var f float64
	if err := stdjson.Unmarshal(raw, &f); err == nil {
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	}
	return "", errors.New("not a string or number")
}