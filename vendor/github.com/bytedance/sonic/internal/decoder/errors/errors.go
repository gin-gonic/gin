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

package errors

import (
    `encoding/json`
    `errors`
    `fmt`
    `reflect`
    `strconv`
    `strings`

    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)

type SyntaxError struct {
    Pos  int
    Src  string
    Code types.ParsingError
    Msg  string
}

func (self SyntaxError) Error() string {
    return fmt.Sprintf("%q", self.Description())
}

func (self SyntaxError) Description() string {
    return "Syntax error " + self.description()
}

func (self SyntaxError) description() string {
    /* check for empty source */
    if self.Src == "" {
        return fmt.Sprintf("no sources available, the input json is empty: %#v", self)
    }

    p, x, q, y := calcBounds(len(self.Src), self.Pos)

    /* compose the error description */
    return fmt.Sprintf(
        "at index %d: %s\n\n\t%s\n\t%s^%s\n",
        self.Pos,
        self.Message(),
        self.Src[p:q],
        strings.Repeat(".", x),
        strings.Repeat(".", y),
    )
}

func calcBounds(size int, pos int) (lbound int, lwidth int, rbound int, rwidth int) {
    if pos >= size || pos < 0 {
        return 0, 0, size, 0
    }

    i := 16
    lbound = pos - i
    rbound = pos + i

    /* prevent slicing before the beginning */
    if lbound < 0 {
        lbound, rbound, i = 0, rbound - lbound, i + lbound
    }

    /* prevent slicing beyond the end */
    if n := size; rbound > n {
        n = rbound - n
        rbound = size

        /* move the left bound if possible */
        if lbound > n {
            i += n
            lbound -= n
        }
    }

    /* left and right length */
    lwidth = clamp_zero(i)
    rwidth = clamp_zero(rbound - lbound - i - 1)

    return
}

func (self SyntaxError) Message() string {
    if self.Msg == "" {
        return self.Code.Message()
    }
    return self.Msg
}

func clamp_zero(v int) int {
    if v < 0 {
        return 0
    } else {
        return v
    }
}

/** JIT Error Helpers **/

var StackOverflow = &json.UnsupportedValueError {
    Str   : "Value nesting too deep",
    Value : reflect.ValueOf("..."),
}

func ErrorWrap(src string, pos int, code types.ParsingError) error {
    return *error_wrap_heap(src, pos, code)
}

//go:noinline
func error_wrap_heap(src string, pos int, code types.ParsingError) *SyntaxError {
    return &SyntaxError {
        Pos  : pos,
        Src  : src,
        Code : code,
    }
}

func ErrorType(vt *rt.GoType) error {
    return &json.UnmarshalTypeError{Type: vt.Pack()}
}

type MismatchTypeError struct {
    Pos  int
    Src  string
    Type reflect.Type
}

func swithchJSONType (src string, pos int) string {
    var val string
    switch src[pos] {
        case 'f': fallthrough
        case 't': val = "bool"
        case '"': val = "string"
        case '{': val = "object"
        case '[': val = "array"
        case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': val = "number"        
    }
    return val
}

func (self MismatchTypeError) Error() string {
    se := SyntaxError {
        Pos  : self.Pos,
        Src  : self.Src,
        Code : types.ERR_MISMATCH,
    }
    return fmt.Sprintf("Mismatch type %s with value %s %q", self.Type.String(), swithchJSONType(self.Src, self.Pos), se.description())
}

func (self MismatchTypeError) Description() string {
    se := SyntaxError {
        Pos  : self.Pos,
        Src  : self.Src,
        Code : types.ERR_MISMATCH,
    }
    return fmt.Sprintf("Mismatch type %s with value %s %s", self.Type.String(), swithchJSONType(self.Src, self.Pos), se.description())
}

func ErrorMismatch(src string, pos int, vt *rt.GoType) error {
    return &MismatchTypeError {
        Pos  : pos,
        Src  : src,
        Type : vt.Pack(),
    }
}

func ErrorField(name string) error {
    return errors.New("json: unknown field " + strconv.Quote(name))
}

func ErrorValue(value string, vtype reflect.Type) error {
    return &json.UnmarshalTypeError {
        Type  : vtype,
        Value : value,
    }
}
