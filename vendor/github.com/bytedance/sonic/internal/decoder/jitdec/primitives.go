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

package jitdec

import (
    `encoding`
    `encoding/json`
    `unsafe`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/rt`
)

func decodeTypedPointer(s string, i int, vt *rt.GoType, vp unsafe.Pointer, sb *_Stack, fv uint64) (int, error) {
    if fn, err := findOrCompile(vt); err != nil {
        return 0, err
    } else {
        rt.MoreStack(_FP_size + _VD_size + native.MaxFrameSize)
        ret, err := fn(s, i, vp, sb, fv, "", nil)
        return ret, err
    }
}

func decodeJsonUnmarshaler(vv interface{}, s string) error {
    return vv.(json.Unmarshaler).UnmarshalJSON(rt.Str2Mem(s))
}

// used to distinguish between MismatchQuoted and other MismatchedTyped errors, see issue #670 and #716
type MismatchQuotedError struct {}

func (*MismatchQuotedError) Error() string {
    return "mismatch quoted"
}

func decodeJsonUnmarshalerQuoted(vv interface{}, s string) error {
    if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
        return &MismatchQuotedError{}
    }
    return vv.(json.Unmarshaler).UnmarshalJSON(rt.Str2Mem(s[1:len(s)-1]))
}

func decodeTextUnmarshaler(vv interface{}, s string) error {
    return vv.(encoding.TextUnmarshaler).UnmarshalText(rt.Str2Mem(s))
}
