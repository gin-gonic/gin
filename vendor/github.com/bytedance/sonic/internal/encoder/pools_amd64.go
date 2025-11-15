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

package encoder

import (
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/encoder/x86"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/option"
)

func ForceUseJit() {
	x86.SetCompiler(makeEncoderX86)
	pretouchType = pretouchTypeX86
	encodeTypedPointer = x86.EncodeTypedPointer
	vars.UseVM = false
}

func init() {
	if vars.UseVM {
		ForceUseVM()
	} else {
		ForceUseJit()
	}
}

var _KeepAlive struct {
	rb    *[]byte
	vp    unsafe.Pointer
	sb    *vars.Stack
	fv    uint64
	err   error
	frame [x86.FP_offs]byte
}

func makeEncoderX86(vt *rt.GoType, ex ...interface{}) (interface{}, error) {
	pp, err := NewCompiler().Compile(vt.Pack(), ex[0].(bool))
	if err != nil {
		return nil, err
	}
	as := x86.NewAssembler(pp)
	as.Name = vt.String()
	return as.Load(), nil
}

func pretouchTypeX86(_vt reflect.Type, opts option.CompileOptions, v uint8) (map[reflect.Type]uint8, error) {
	/* compile function */
	compiler := NewCompiler().apply(opts)
	encoder := func(vt *rt.GoType, ex ...interface{}) (interface{}, error) {
		pp, err := compiler.Compile(vt.Pack(), ex[0].(bool))
		if err != nil {
			return nil, err
		}
		as := x86.NewAssembler(pp)
		as.Name = vt.String()
		return as.Load(), nil
	}

	/* find or compile */
	vt := rt.UnpackType(_vt)
	if val := vars.GetProgram(vt); val != nil {
		return nil, nil
	} else if _, err := vars.ComputeProgram(vt, encoder, v == 1); err == nil {
		return compiler.rec, nil
	} else {
		return nil, err
	}
}
