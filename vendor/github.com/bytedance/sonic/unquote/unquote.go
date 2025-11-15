//go:build (amd64 && go1.17 && !go1.26) || (arm64 && go1.20 && !go1.26)
// +build amd64,go1.17,!go1.26 arm64,go1.20,!go1.26


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

package unquote

import (
    `unsafe`
    `runtime`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)

// String unescapes an escaped string (not including `"` at beginning and end)
// It validates invalid UTF8 and replace with `\ufffd`
func String(s string) (ret string, err types.ParsingError) {
    mm := make([]byte, 0, len(s))
    err = intoBytesUnsafe(s, &mm, true)
    ret = rt.Mem2Str(mm)
    return
}

// IntoBytes is same with String besides it output result into a buffer m
func IntoBytes(s string, m *[]byte) types.ParsingError {
    if cap(*m) < len(s) {
        return types.ERR_EOF
    } else {
        return intoBytesUnsafe(s, m, true)
    }
}

// String unescapes an escaped string (not including `"` at beginning and end)
//   - replace enables replacing invalid utf8 escaped char with `\uffd`
func _String(s string, replace bool) (ret string, err error) {
    mm := make([]byte, 0, len(s))
    err = intoBytesUnsafe(s, &mm, replace)
    ret = rt.Mem2Str(mm)
    return
}

func intoBytesUnsafe(s string, m *[]byte, replace bool) types.ParsingError {
    pos := -1
    slv := (*rt.GoSlice)(unsafe.Pointer(m))
    str := (*rt.GoString)(unsafe.Pointer(&s))

    flags := uint64(0)
    if replace {
        /* unquote as the default configuration, replace invalid unicode with \ufffd */
        flags |= types.F_UNICODE_REPLACE
    }

    ret := native.Unquote(str.Ptr, str.Len, slv.Ptr, &pos, flags)

    /* check for errors */
    if ret < 0 {
        return types.ParsingError(-ret)
    }

    /* update the length */
    slv.Len = ret
    runtime.KeepAlive(s)
    return 0
}



