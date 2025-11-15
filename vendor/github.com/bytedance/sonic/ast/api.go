//go:build (amd64 && go1.17 && !go1.26) || (arm64 && go1.20 && !go1.26)
// +build amd64,go1.17,!go1.26 arm64,go1.20,!go1.26

/*
 * Copyright 2022 ByteDance Inc.
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

package ast

import (
	"runtime"
	"unsafe"

    "github.com/bytedance/sonic/encoder"
    "github.com/bytedance/sonic/internal/encoder/alg"
    "github.com/bytedance/sonic/internal/native"
    "github.com/bytedance/sonic/internal/native/types"
    "github.com/bytedance/sonic/internal/rt"
    "github.com/bytedance/sonic/utf8"
)

var typeByte = rt.UnpackEface(byte(0)).Type

func quote(buf *[]byte, val string) {
    *buf = alg.Quote(*buf, val, false)
}

func (self *Parser) decodeValue() (val types.JsonState) {
    sv := (*rt.GoString)(unsafe.Pointer(&self.s))
    flag := types.F_USE_NUMBER
    if self.dbuf != nil {
        flag = 0
        val.Dbuf = self.dbuf
        val.Dcap = types.MaxDigitNums
    }
    self.p = native.Value(sv.Ptr, sv.Len, self.p, &val, uint64(flag))
    return
}

func (self *Parser) skip() (int, types.ParsingError) {
    fsm := types.NewStateMachine()
    start := native.SkipOne(&self.s, &self.p, fsm, 0)
    types.FreeStateMachine(fsm)

    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func (self *Node) encodeInterface(buf *[]byte) error {
    //WARN: NOT compatible with json.Encoder
    return encoder.EncodeInto(buf, self.packAny(), encoder.NoEncoderNewline)
}

func (self *Parser) skipFast() (int, types.ParsingError) {
    start := native.SkipOneFast(&self.s, &self.p)
    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func (self *Parser) getByPath(validate bool, path ...interface{}) (int, types.ParsingError) {
    var fsm *types.StateMachine
    if validate {
        fsm = types.NewStateMachine()
    }
    start := native.GetByPath(&self.s, &self.p, &path, fsm)
    if validate {
        types.FreeStateMachine(fsm)
    }
    runtime.KeepAlive(path)
    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func validate_utf8(str string) bool {
    return utf8.ValidateString(str)
}
