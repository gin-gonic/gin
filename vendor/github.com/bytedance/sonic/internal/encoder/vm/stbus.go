/**
 * Copyright 2024 ByteDance Inc.
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

package vm

import (
	"unsafe"
	_ "unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/ir"
	"github.com/bytedance/sonic/internal/encoder/prim"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
)

func EncodeTypedPointer(buf *[]byte, vt *rt.GoType, vp *unsafe.Pointer, sb *vars.Stack, fv uint64) error {
	if vt == nil {
		return prim.EncodeNil(buf)
	} else if pp, err := vars.FindOrCompile(vt, (fv&(1<<alg.BitPointerValue)) != 0, compiler); err != nil {
		return err
	} else if vt.Indirect() {
		return Execute(buf, *vp, sb, fv, pp.(*ir.Program))
	} else {
		return Execute(buf, unsafe.Pointer(vp), sb, fv, pp.(*ir.Program))
	}
}

var compiler func(*rt.GoType, ... interface{}) (interface{}, error)

func SetCompiler(c func(*rt.GoType, ... interface{}) (interface{}, error)) {
	compiler = c
}
