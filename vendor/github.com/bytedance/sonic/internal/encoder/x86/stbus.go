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

package x86

import (
	"unsafe"
	_ "unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/prim"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/loader"
	_ "github.com/cloudwego/base64x"
)

var compiler func(*rt.GoType, ... interface{}) (interface{}, error)

func SetCompiler(c func(*rt.GoType, ... interface{}) (interface{}, error)) {
	compiler = c
}

func ptoenc(p loader.Function) vars.Encoder {
    return *(*vars.Encoder)(unsafe.Pointer(&p))
}

func EncodeTypedPointer(buf *[]byte, vt *rt.GoType, vp *unsafe.Pointer, sb *vars.Stack, fv uint64) error {
	if vt == nil {
		return prim.EncodeNil(buf)
	} else if fn, err := vars.FindOrCompile(vt, (fv&(1<<alg.BitPointerValue)) != 0, compiler); err != nil {
		return err
	} else if vt.Indirect() {
		return	fn.(vars.Encoder)(buf, *vp, sb, fv)
	} else {
		return fn.(vars.Encoder)(buf, unsafe.Pointer(vp), sb, fv)
	}
}

