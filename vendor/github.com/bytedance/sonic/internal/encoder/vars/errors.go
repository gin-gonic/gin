/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vars

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

var ERR_too_deep = &json.UnsupportedValueError {
    Str   : "Value nesting too deep",
    Value : reflect.ValueOf("..."),
}

var ERR_nan_or_infinite = &json.UnsupportedValueError {
    Str   : "NaN or ±Infinite",
    Value : reflect.ValueOf("NaN or ±Infinite"),
}

func Error_type(vtype reflect.Type) error {
    return &json.UnsupportedTypeError{Type: vtype}
}

func Error_number(number json.Number) error {
    return &json.UnsupportedValueError {
        Str   : "invalid number literal: " + strconv.Quote(string(number)),
        Value : reflect.ValueOf(number),
    }
}

func Error_unsuppoted(typ *rt.GoType) error {
	return &json.UnsupportedTypeError{Type: typ.Pack() }
}

func Error_marshaler(ret []byte, pos int) error {
    return fmt.Errorf("invalid Marshaler output json syntax at %d: %q", pos, ret)
}

const (
    PanicNilPointerOfNonEmptyString int = 1 + iota
)

func GoPanic(code int, val unsafe.Pointer, buf string) {
    sb := strings.Builder{}
    switch(code){
    case PanicNilPointerOfNonEmptyString:
        sb.WriteString(fmt.Sprintf("val: %#v has nil pointer while its length is not zero!\nThis is a nil pointer exception (NPE) problem. There might be a data race issue. It is recommended to execute the tests related to the code with the `-race` compile flag to detect the problem.\n", (*rt.GoString)(val)))
    default:
        sb.WriteString("encoder error: ")
        sb.WriteString(strconv.Itoa(code))
        sb.WriteString("\n")
    }
    sb.WriteString("JSON: ")
    if len(buf) > maxJSONLength {
        sb.WriteString(buf[len(buf)-maxJSONLength:])
    } else {
        sb.WriteString(buf)
    }
    panic(sb.String())
}

var maxJSONLength = 1024

func init() {
    if v := os.Getenv("SONIC_PANIC_MAX_JSON_LENGTH"); v != "" {
        if i, err := strconv.Atoi(v); err == nil {
            maxJSONLength = i
        }
    }
}
