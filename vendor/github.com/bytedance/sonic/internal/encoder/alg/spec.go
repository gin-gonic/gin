//go:build (amd64 && go1.16 && !go1.26) || (arm64 && go1.20 && !go1.26)
// +build amd64,go1.16,!go1.26 arm64,go1.20,!go1.26

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

package alg

import (
	"runtime"
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
)

// Valid validates json and returns first non-blank character position,
// if it is only one valid json value.
// Otherwise returns invalid character position using start.
//
// Note: it does not check for the invalid UTF-8 characters.
func Valid(data []byte) (ok bool, start int) {
    n := len(data)
    if n == 0 {
        return false, -1
    }
    s := rt.Mem2Str(data)
    p := 0
    m := types.NewStateMachine()
    ret := native.ValidateOne(&s, &p, m, 0)
    types.FreeStateMachine(m)

    if ret < 0 {
        return false, p-1
    }

    /* check for trailing spaces */
    for ;p < n; p++ {
        if (types.SPACE_MASK & (1 << data[p])) == 0 {
            return false, p
        }
    }

    return true, ret
}

var typeByte = rt.UnpackEface(byte(0)).Type

func Quote(buf []byte, val string, double bool) []byte {
	if len(val) == 0 {
		if double {
			return append(buf, `"\"\""`...)
		}
		return append(buf, `""`...)
	}

	if double {
		buf = append(buf, `"\"`...)
	} else {
		buf = append(buf, `"`...)
	}
	sp := rt.IndexChar(val, 0)
	nb := len(val)

	buf = rt.GuardSlice2(buf, nb+1)
	b := (*rt.GoSlice)(unsafe.Pointer(&buf))

	// input buffer
	for nb > 0 {
		// output buffer
		dp := unsafe.Pointer(uintptr(b.Ptr) + uintptr(b.Len))
		dn := b.Cap - b.Len
		// call native.Quote, dn is byte count it outputs
		opts := uint64(0)
		if double {
			opts = types.F_DOUBLE_UNQUOTE
		}
		ret := native.Quote(sp, nb, dp, &dn, opts)
		// update *buf length
		b.Len += dn

		// no need more output
		if ret >= 0 {
			break
		}

		// double buf size
		*b = rt.GrowSlice(typeByte, *b, b.Cap*2)
		// ret is the complement of consumed input
		ret = ^ret
		// update input buffer
		nb -= ret
		if nb > 0 {
			sp = unsafe.Pointer(uintptr(sp) + uintptr(ret))
		}
	}

	runtime.KeepAlive(buf)
	runtime.KeepAlive(sp)
	if double {
		buf = append(buf, `\""`...)
	} else {
		buf = append(buf, `"`...)
	}

	return buf
}

func HtmlEscape(dst []byte, src []byte) []byte {
	var sidx int

	dst = append(dst, src[:0]...) // avoid check nil dst
	sbuf := (*rt.GoSlice)(unsafe.Pointer(&src))
	dbuf := (*rt.GoSlice)(unsafe.Pointer(&dst))

	/* grow dst if it is shorter */
	if cap(dst)-len(dst) < len(src)+types.BufPaddingSize {
		cap := len(src)*3/2 + types.BufPaddingSize
		*dbuf = rt.GrowSlice(typeByte, *dbuf, cap)
	}

	for sidx < sbuf.Len {
		sp := rt.Add(sbuf.Ptr, uintptr(sidx))
		dp := rt.Add(dbuf.Ptr, uintptr(dbuf.Len))

		sn := sbuf.Len - sidx
		dn := dbuf.Cap - dbuf.Len
		nb := native.HTMLEscape(sp, sn, dp, &dn)

		/* check for errors */
		if dbuf.Len += dn; nb >= 0 {
			break
		}

		/* not enough space, grow the slice and try again */
		sidx += ^nb
		*dbuf = rt.GrowSlice(typeByte, *dbuf, dbuf.Cap*2)
	}
	return dst
}

func F64toa(buf []byte, v float64) ([]byte) {
	if v == 0 {
		return append(buf, '0')
	}
	buf = rt.GuardSlice2(buf, 64)
	ret := native.F64toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

func F32toa(buf []byte, v float32) ([]byte) {
	if v == 0 {
		return append(buf, '0')
	}
	buf = rt.GuardSlice2(buf, 64)
	ret := native.F32toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

func I64toa(buf []byte, v int64) ([]byte) {
	return strconv.AppendInt(buf, v, 10)
}

func U64toa(buf []byte, v uint64) ([]byte) {
	return strconv.AppendUint(buf, v, 10)
}
