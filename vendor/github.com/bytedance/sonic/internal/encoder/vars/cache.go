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
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

type Encoder func(
	rb *[]byte,
	vp unsafe.Pointer,
	sb *Stack,
	fv uint64,
) error

func FindOrCompile(vt *rt.GoType, pv bool, compiler func(*rt.GoType, ... interface{}) (interface{}, error)) (interface{}, error) {
	if val := programCache.Get(vt); val != nil {
		return val, nil
	} else if ret, err := programCache.Compute(vt, compiler, pv); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
}

func GetProgram(vt *rt.GoType) (interface{}) {
	return programCache.Get(vt)
}

func ComputeProgram(vt *rt.GoType, compute func(*rt.GoType, ... interface{}) (interface{}, error), pv bool) (interface{}, error) {
	return programCache.Compute(vt, compute, pv)
}