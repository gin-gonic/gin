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

package prim

import (
	"encoding"
	"encoding/json"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/resolver"
	"github.com/bytedance/sonic/internal/rt"
)

func Compact(p *[]byte, v []byte) error {
	buf := vars.NewBuffer()
	err := json.Compact(buf, v)

	/* check for errors */
	if err != nil {
		return err
	}

	/* add to result */
	v = buf.Bytes()
	*p = append(*p, v...)

	/* return the buffer into pool */
	vars.FreeBuffer(buf)
	return nil
}

func EncodeNil(rb *[]byte) error {
	*rb = append(*rb, 'n', 'u', 'l', 'l')
	return nil
}

// func Make_EncodeTypedPointer(computor func(*rt.GoType, ...interface{}) (interface{}, error)) func(*[]byte, *rt.GoType, *unsafe.Pointer, *vars.Stack, uint64) error {
// 	return func(buf *[]byte, vt *rt.GoType, vp *unsafe.Pointer, sb *vars.Stack, fv uint64) error {
// 		if vt == nil {
// 			return EncodeNil(buf)
// 		} else if fn, err := vars.FindOrCompile(vt, (fv&(1<<BitPointerValue)) != 0, computor); err != nil {
// 			return err
// 		} else if vt.Indirect() {
// 			err := fn(buf, *vp, sb, fv)
// 			return err
// 		} else {
// 			err := fn(buf, unsafe.Pointer(vp), sb, fv)
// 			return err
// 		}
// 	}
// }

func EncodeJsonMarshaler(buf *[]byte, val json.Marshaler, opt uint64) error {
	if ret, err := val.MarshalJSON(); err != nil {
		return err
	} else {
		if opt&(1<<alg.BitCompactMarshaler) != 0 {
			return Compact(buf, ret)
		}
		if opt&(1<<alg.BitNoValidateJSONMarshaler) == 0 {
			if ok, s := alg.Valid(ret); !ok {
				return vars.Error_marshaler(ret, s)
			}
		}
		*buf = append(*buf, ret...)
		return nil
	}
}

func EncodeTextMarshaler(buf *[]byte, val encoding.TextMarshaler, opt uint64) error {
	if ret, err := val.MarshalText(); err != nil {
		return err
	} else {
		if opt&(1<<alg.BitNoQuoteTextMarshaler) != 0 {
			*buf = append(*buf, ret...)
			return nil
		}
		*buf = alg.Quote(*buf, rt.Mem2Str(ret), false)
		return nil
	}
}

func IsZero(val unsafe.Pointer, fv *resolver.FieldMeta) bool {
	rv := reflect.NewAt(fv.Type, val).Elem()
	b1 := fv.IsZero == nil && rv.IsZero()
	b2 := fv.IsZero != nil && fv.IsZero(rv)
	return  b1 || b2
}
