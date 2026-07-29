// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	ginjson "github.com/gin-gonic/gin/codec/json"
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

// jsonSource implements setter for JSON binding, allowing time_format
// and time_utc tags to work with JSON binding (ShouldBindJSON).
// For time.Time fields with time_format tags, it uses setTimeField
// to parse the string value. For all other fields, it uses standard
// JSON unmarshaling.
type jsonSource map[string]json.RawMessage

var _ setter = jsonSource(nil)

// timeTimeType is the reflect.Type for time.Time, used for comparison.
var timeTimeType = reflect.TypeOf(time.Time{})

// TrySet tries to set a value from a JSON source.
func (j jsonSource) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSet bool, err error) {
	raw, ok := j[key]
	if !ok {
		return false, nil
	}

	// For time.Time fields with a time_format tag, use setTimeField
	// to parse the string value using the custom format.
	if value.Type() == timeTimeType && field.Tag.Get("time_format") != "" {
		// Unmarshal the raw JSON value as a string
		var s string
		if err := ginjson.API.Unmarshal(raw, &s); err != nil {
			// If the raw value is a number (unix timestamp), try that too
			var n int64
			if err2 := ginjson.API.Unmarshal(raw, &n); err2 != nil {
				return false, err
			}
			// Convert number to string for setTimeField
			s = strconv.FormatInt(n, 10)
		}
		return true, setTimeField(s, field, value)
	}

	// For all other types, use standard JSON unmarshaling.
	if err := ginjson.API.Unmarshal(raw, value.Addr().Interface()); err != nil {
		return false, err
	}
	return true, nil
}

func decodeJSON(r io.Reader, obj any) error {
	// Read the full body first so we can retry if needed
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty JSON body")
	}

	// Fast path: try the original decoder approach first (preserves UseNumber,
	// DisallowUnknownFields, and works with any JSON codec backend).
	decoder := ginjson.API.NewDecoder(bytes.NewReader(body))
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		// If the struct has time_format tags, the decode may have failed because
		// a time.Time field used a non-RFC3339 format. Fall back to the jsonSource
		// approach which respects the time_format tag.
		if hasTimeFormatTag(obj) {
			if err2 := decodeJSONWithTimeFormat(body, obj); err2 != nil {
				return err // return the original decode error, which is more descriptive
			}
			return validate(obj)
		}
		return err
	}

	return validate(obj)
}

// decodeJSONWithTimeFormat decodes JSON using mappingByPtr with a jsonSource,
// which respects time_format tags on time.Time fields.
func decodeJSONWithTimeFormat(body []byte, obj any) error {
	var rawMap map[string]json.RawMessage
	if err := ginjson.API.Unmarshal(body, &rawMap); err != nil {
		return err
	}

	if err := mappingByPtr(obj, jsonSource(rawMap), "json"); err != nil {
		return err
	}

	return nil
}

// hasTimeFormatTag checks if the given struct (or any nested struct) has
// any time.Time fields with a time_format tag.
func hasTimeFormatTag(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return hasTimeFormatTagRecursive(v.Type())
}

func hasTimeFormatTagRecursive(t reflect.Type) bool {
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			if f.Type == timeTimeType {
				if f.Tag.Get("time_format") != "" {
					return true
				}
			} else if f.Anonymous {
				if hasTimeFormatTagRecursive(f.Type) {
					return true
				}
			}
		}
	}
	return false
}