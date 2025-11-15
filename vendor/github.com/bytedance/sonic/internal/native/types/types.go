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

package types

import (
    `fmt`
    `sync`
    `unsafe`
)

type ValueType = int64
type ParsingError uint
type SearchingError uint

// NOTE: !NOT MODIFIED ONLY.
// This definitions are followed in native/types.h.

const BufPaddingSize int     = 64

const (
    V_EOF     ValueType = 1
    V_NULL    ValueType = 2
    V_TRUE    ValueType = 3
    V_FALSE   ValueType = 4
    V_ARRAY   ValueType = 5
    V_OBJECT  ValueType = 6
    V_STRING  ValueType = 7
    V_DOUBLE  ValueType = 8
    V_INTEGER ValueType = 9
    _         ValueType = 10    // V_KEY_SEP
    _         ValueType = 11    // V_ELEM_SEP
    _         ValueType = 12    // V_ARRAY_END
    _         ValueType = 13    // V_OBJECT_END
    V_MAX
)

const (
    // for native.Unquote() flags
    B_DOUBLE_UNQUOTE  = 0
    B_UNICODE_REPLACE = 1

    // for native.Value() flags
    B_USE_NUMBER      = 1
    B_VALIDATE_STRING = 5
    B_ALLOW_CONTROL   = 31

    // for native.SkipOne() flags
    B_NO_VALIDATE_JSON= 6
)

const (
    F_DOUBLE_UNQUOTE  = 1 << B_DOUBLE_UNQUOTE
    F_UNICODE_REPLACE = 1 << B_UNICODE_REPLACE

    F_USE_NUMBER      = 1 << B_USE_NUMBER
    F_VALIDATE_STRING = 1 << B_VALIDATE_STRING
    F_ALLOW_CONTROL   = 1 << B_ALLOW_CONTROL
)

const (
    MAX_RECURSE = 4096
)

const (
    SPACE_MASK = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')
)

const (
    ERR_EOF                ParsingError = 1
    ERR_INVALID_CHAR       ParsingError = 2
    ERR_INVALID_ESCAPE     ParsingError = 3
    ERR_INVALID_UNICODE    ParsingError = 4
    ERR_INTEGER_OVERFLOW   ParsingError = 5
    ERR_INVALID_NUMBER_FMT ParsingError = 6
    ERR_RECURSE_EXCEED_MAX ParsingError = 7
    ERR_FLOAT_INFINITY     ParsingError = 8
    ERR_MISMATCH           ParsingError = 9
    ERR_INVALID_UTF8       ParsingError = 10

    // error code used in ast
    ERR_NOT_FOUND          ParsingError = 33
    ERR_UNSUPPORT_TYPE     ParsingError = 34
)

var _ParsingErrors = []string{
    0                      : "ok",
    ERR_EOF                : "eof",
    ERR_INVALID_CHAR       : "invalid char",
    ERR_INVALID_ESCAPE     : "invalid escape char",
    ERR_INVALID_UNICODE    : "invalid unicode escape",
    ERR_INTEGER_OVERFLOW   : "integer overflow",
    ERR_INVALID_NUMBER_FMT : "invalid number format",
    ERR_RECURSE_EXCEED_MAX : "recursion exceeded max depth",
    ERR_FLOAT_INFINITY     : "float number is infinity",
    ERR_MISMATCH           : "mismatched type with value",
    ERR_INVALID_UTF8       : "invalid UTF8",
}

func (self ParsingError) Error() string {
    return "json: error when parsing input: " + self.Message()
}

func (self ParsingError) Message() string {
    if int(self) < len(_ParsingErrors) {
        return _ParsingErrors[self]
    } else {
        return fmt.Sprintf("unknown error %d", self)
    }
}

type JsonState struct {
    Vt ValueType
    Dv   float64
    Iv   int64
    Ep   int
    Dbuf *byte
    Dcap int
}

type StateMachine struct {
    Sp int
    Vt [MAX_RECURSE]int
}

var stackPool = sync.Pool{
    New: func()interface{}{
        return &StateMachine{}
    },
}

func NewStateMachine() *StateMachine {
    return stackPool.Get().(*StateMachine)
}

func FreeStateMachine(fsm *StateMachine) {
    stackPool.Put(fsm)
}

const MaxDigitNums = 800

var digitPool = sync.Pool{
    New: func() interface{} {
        return (*byte)(unsafe.Pointer(&[MaxDigitNums]byte{}))
    },
}

func NewDbuf() *byte {
    return digitPool.Get().(*byte)
}

func FreeDbuf(p *byte) {
    digitPool.Put(p)
}
