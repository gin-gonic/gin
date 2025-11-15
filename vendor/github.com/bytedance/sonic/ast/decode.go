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
	"encoding/base64"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/internal/utils"
	"github.com/bytedance/sonic/unquote"
)


var bytesNull   = []byte("null")

const (
    strNull   = "null"
    bytesTrue   = "true"
    bytesFalse  = "false"
    bytesObject = "{}"
    bytesArray  = "[]"
)

//go:nocheckptr
func skipBlank(src string, pos int) int {
    se := uintptr(rt.IndexChar(src, len(src)))
    sp := uintptr(rt.IndexChar(src, pos))

    for sp < se {
        if !utils.IsSpace(*(*byte)(unsafe.Pointer(sp))) {
            break
        }
        sp += 1
    }
    if sp >= se {
        return -int(types.ERR_EOF)
    }
    runtime.KeepAlive(src)
    return int(sp - uintptr(rt.IndexChar(src, 0)))
}

func decodeNull(src string, pos int) (ret int) {
    ret = pos + 4
    if ret > len(src) {
        return -int(types.ERR_EOF)
    }
    if src[pos:ret] == strNull {
        return ret
    } else {
        return -int(types.ERR_INVALID_CHAR)
    }
}

func decodeTrue(src string, pos int) (ret int) {
    ret = pos + 4
    if ret > len(src) {
        return -int(types.ERR_EOF)
    }
    if src[pos:ret] == bytesTrue {
        return ret
    } else {
        return -int(types.ERR_INVALID_CHAR)
    }

}

func decodeFalse(src string, pos int) (ret int) {
    ret = pos + 5
    if ret > len(src) {
        return -int(types.ERR_EOF)
    }
    if src[pos:ret] == bytesFalse {
        return ret
    }
    return -int(types.ERR_INVALID_CHAR)
}

//go:nocheckptr
func decodeString(src string, pos int) (ret int, v string) {
    ret, ep := skipString(src, pos)
    if ep == -1 {
        (*rt.GoString)(unsafe.Pointer(&v)).Ptr = rt.IndexChar(src, pos+1)
        (*rt.GoString)(unsafe.Pointer(&v)).Len = ret - pos - 2
        return ret, v
    }

    result, err := unquote.String(src[pos:ret])
    if err != 0 {
        return -int(types.ERR_INVALID_CHAR), ""
    }

    runtime.KeepAlive(src)
    return ret, result
}

func decodeBinary(src string, pos int) (ret int, v []byte) {
    var vv string
    ret, vv = decodeString(src, pos)
    if ret < 0 {
        return ret, nil
    }
    var err error
    v, err = base64.StdEncoding.DecodeString(vv)
    if err != nil {
        return -int(types.ERR_INVALID_CHAR), nil
    }
    return ret, v
}

func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

//go:nocheckptr
func decodeInt64(src string, pos int) (ret int, v int64, err error) {
    sp := uintptr(rt.IndexChar(src, pos))
    ss := uintptr(sp)
    se := uintptr(rt.IndexChar(src, len(src)))
    if uintptr(sp) >= se {
        return -int(types.ERR_EOF), 0, nil
    }

    if c := *(*byte)(unsafe.Pointer(sp)); c == '-' {
        sp += 1
    }
    if sp == se {
        return -int(types.ERR_EOF), 0, nil
    }

    for ; sp < se; sp += uintptr(1) {
        if !isDigit(*(*byte)(unsafe.Pointer(sp))) {
            break
        }
    }

    if sp < se {
        if c := *(*byte)(unsafe.Pointer(sp)); c == '.' || c == 'e' || c == 'E' {
            return -int(types.ERR_INVALID_NUMBER_FMT), 0, nil
        }
    }

    var vv string
    ret = int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr))
    (*rt.GoString)(unsafe.Pointer(&vv)).Ptr = unsafe.Pointer(ss)
    (*rt.GoString)(unsafe.Pointer(&vv)).Len = ret - pos

    v, err = strconv.ParseInt(vv, 10, 64)
    if err != nil {
        //NOTICE: allow overflow here
        if err.(*strconv.NumError).Err == strconv.ErrRange {
            return ret, 0, err
        }
        return -int(types.ERR_INVALID_CHAR), 0, err
    }

    runtime.KeepAlive(src)
    return ret, v, nil
}

func isNumberChars(c byte) bool {
    return (c >= '0' && c <= '9') || c == '+' || c == '-' || c == 'e' || c == 'E' || c == '.'
}

//go:nocheckptr
func decodeFloat64(src string, pos int) (ret int, v float64, err error) {
    sp := uintptr(rt.IndexChar(src, pos))
    ss := uintptr(sp)
    se := uintptr(rt.IndexChar(src, len(src)))
    if uintptr(sp) >= se {
        return -int(types.ERR_EOF), 0, nil
    }

    if c := *(*byte)(unsafe.Pointer(sp)); c == '-' {
        sp += 1
    }
    if sp == se {
        return -int(types.ERR_EOF), 0, nil
    }

    for ; sp < se; sp += uintptr(1) {
        if !isNumberChars(*(*byte)(unsafe.Pointer(sp))) {
            break
        }
    }

    var vv string
    ret = int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr))
    (*rt.GoString)(unsafe.Pointer(&vv)).Ptr = unsafe.Pointer(ss)
    (*rt.GoString)(unsafe.Pointer(&vv)).Len = ret - pos

    v, err = strconv.ParseFloat(vv, 64)
    if err != nil {
        //NOTICE: allow overflow here
        if err.(*strconv.NumError).Err == strconv.ErrRange {
            return ret, 0, err
        }
        return -int(types.ERR_INVALID_CHAR), 0, err
    }

    runtime.KeepAlive(src)
    return ret, v, nil
}

func decodeValue(src string, pos int, skipnum bool) (ret int, v types.JsonState) {
    pos = skipBlank(src, pos)
    if pos < 0 {
        return pos, types.JsonState{Vt: types.ValueType(pos)}
    }
    switch c := src[pos]; c {
    case 'n':
        ret = decodeNull(src, pos)
        if ret < 0 {
            return ret, types.JsonState{Vt: types.ValueType(ret)}
        }
        return ret, types.JsonState{Vt: types.V_NULL}
    case '"':
        var ep int
        ret, ep = skipString(src, pos)
        if ret < 0 {
            return ret, types.JsonState{Vt: types.ValueType(ret)}
        }
        return ret, types.JsonState{Vt: types.V_STRING, Iv: int64(pos + 1), Ep: ep}
    case '{':
        return pos + 1, types.JsonState{Vt: types.V_OBJECT}
    case '[':
        return pos + 1, types.JsonState{Vt: types.V_ARRAY}
    case 't':
        ret = decodeTrue(src, pos)
        if ret < 0 {
            return ret, types.JsonState{Vt: types.ValueType(ret)}
        }
        return ret, types.JsonState{Vt: types.V_TRUE}
    case 'f':
        ret = decodeFalse(src, pos)
        if ret < 0 {
            return ret, types.JsonState{Vt: types.ValueType(ret)}
        }
        return ret, types.JsonState{Vt: types.V_FALSE}
    case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        if skipnum {
            ret = skipNumber(src, pos)
            if ret >= 0 {
                return ret, types.JsonState{Vt: types.V_DOUBLE, Iv: 0, Ep: pos}
            } else {
                return ret, types.JsonState{Vt: types.ValueType(ret)}
            }
        } else {
            var iv int64
            ret, iv, _ = decodeInt64(src, pos)
            if ret >= 0 {
                return ret, types.JsonState{Vt: types.V_INTEGER, Iv: iv, Ep: pos}
            } else if ret != -int(types.ERR_INVALID_NUMBER_FMT) {
                return ret, types.JsonState{Vt: types.ValueType(ret)}
            }
            var fv float64
            ret, fv, _ = decodeFloat64(src, pos)
            if ret >= 0 {
                return ret, types.JsonState{Vt: types.V_DOUBLE, Dv: fv, Ep: pos}
            } else {
                return ret, types.JsonState{Vt: types.ValueType(ret)}
            }
        }
        
    default:
        return -int(types.ERR_INVALID_CHAR), types.JsonState{Vt:-types.ValueType(types.ERR_INVALID_CHAR)}
    }
}

//go:nocheckptr
func skipNumber(src string, pos int) (ret int) {
    return utils.SkipNumber(src, pos)
}

//go:nocheckptr
func skipString(src string, pos int) (ret int, ep int) {
    if pos+1 >= len(src) {
        return -int(types.ERR_EOF), -1
    }

    sp := uintptr(rt.IndexChar(src, pos))
    se := uintptr(rt.IndexChar(src, len(src)))

    // not start with quote
    if *(*byte)(unsafe.Pointer(sp)) != '"' {
        return -int(types.ERR_INVALID_CHAR), -1
    }
    sp += 1

    ep = -1
    for sp < se {
        c := *(*byte)(unsafe.Pointer(sp))
        if c == '\\' {
            if ep == -1 {
                ep = int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr))
            }
            sp += 2
            continue
        }
        sp += 1
        if c == '"' {
            return int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr)), ep
        }
    }

    runtime.KeepAlive(src)
    // not found the closed quote until EOF
    return -int(types.ERR_EOF), -1
}

//go:nocheckptr
func skipPair(src string, pos int, lchar byte, rchar byte) (ret int) {
    if pos+1 >= len(src) {
        return -int(types.ERR_EOF)
    }

    sp := uintptr(rt.IndexChar(src, pos))
    se := uintptr(rt.IndexChar(src, len(src)))

    if *(*byte)(unsafe.Pointer(sp)) != lchar {
        return -int(types.ERR_INVALID_CHAR)
    }

    sp += 1
    nbrace := 1
    inquote := false

    for sp < se {
        c := *(*byte)(unsafe.Pointer(sp))
        if c == '\\' {
            sp += 2
            continue
        } else if c == '"' {
            inquote = !inquote
        } else if c == lchar {
            if !inquote {
                nbrace += 1
            }
        } else if c == rchar {
            if !inquote {
                nbrace -= 1
                if nbrace == 0 {
                    sp += 1
                    break
                }
            }
        }
        sp += 1
    }

    if nbrace != 0 {
        return -int(types.ERR_INVALID_CHAR)
    }

    runtime.KeepAlive(src)
    return int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr))
}

func skipValueFast(src string, pos int) (ret int, start int) {
    pos = skipBlank(src, pos)
    if pos < 0 {
        return pos, -1
    }
    switch c := src[pos]; c {
    case 'n':
        ret = decodeNull(src, pos)
    case '"':
        ret, _ = skipString(src, pos)
    case '{':
        ret = skipPair(src, pos, '{', '}')
    case '[':
        ret = skipPair(src, pos, '[', ']')
    case 't':
        ret = decodeTrue(src, pos)
    case 'f':
        ret = decodeFalse(src, pos)
    case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        ret = skipNumber(src, pos)
    default:
        ret = -int(types.ERR_INVALID_CHAR)
    }
    return ret, pos
}

func skipValue(src string, pos int) (ret int, start int) {
    pos = skipBlank(src, pos)
    if pos < 0 {
        return pos, -1
    }
    switch c := src[pos]; c {
    case 'n':
        ret = decodeNull(src, pos)
    case '"':
        ret, _ = skipString(src, pos)
    case '{':
        ret, _ = skipObject(src, pos)
    case '[':
        ret, _ = skipArray(src, pos)
    case 't':
        ret = decodeTrue(src, pos)
    case 'f':
        ret = decodeFalse(src, pos)
    case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        ret = skipNumber(src, pos)
    default:
        ret = -int(types.ERR_INVALID_CHAR)
    }
    return ret, pos
}

func skipObject(src string, pos int) (ret int, start int) {
    start = skipBlank(src, pos)
    if start < 0 {
        return start, -1
    }

    if src[start] != '{' {
        return -int(types.ERR_INVALID_CHAR), -1
    }

    pos = start + 1
    pos = skipBlank(src, pos)
    if pos < 0 {
        return pos, -1
    }
    if src[pos] == '}' {
        return pos + 1, start
    }

    for {
        pos, _ = skipString(src, pos)
        if pos < 0 {
            return pos, -1
        }

        pos = skipBlank(src, pos)
        if pos < 0 {
            return pos, -1
        }
        if src[pos] != ':' {
            return -int(types.ERR_INVALID_CHAR), -1
        }

        pos++
        pos, _ = skipValue(src, pos)
        if pos < 0 {
            return pos, -1
        }

        pos = skipBlank(src, pos)
        if pos < 0 {
            return pos, -1
        }
        if src[pos] == '}' {
            return pos + 1, start
        }
        if src[pos] != ',' {
            return -int(types.ERR_INVALID_CHAR), -1
        }

        pos++
        pos = skipBlank(src, pos)
        if pos < 0 {
            return pos, -1
        }

    }
}

func skipArray(src string, pos int) (ret int, start int) {
    start = skipBlank(src, pos)
    if start < 0 {
        return start, -1
    }

    if src[start] != '[' {
        return -int(types.ERR_INVALID_CHAR), -1
    }

    pos = start + 1
    pos = skipBlank(src, pos)
    if pos < 0 {
        return pos, -1
    }
    if src[pos] == ']' {
        return pos + 1, start
    }

    for {
        pos, _ = skipValue(src, pos)
        if pos < 0 {
            return pos, -1
        }

        pos = skipBlank(src, pos)
        if pos < 0 {
            return pos, -1
        }
        if src[pos] == ']' {
            return pos + 1, start
        }
        if src[pos] != ',' {
            return -int(types.ERR_INVALID_CHAR), -1
        }
        pos++
    }
}

// DecodeString decodes a JSON string from pos and return golang string.
//   - needEsc indicates if to unescaped escaping chars
//   - hasEsc tells if the returned string has escaping chars
//   - validStr enables validating UTF8 charset
//
func _DecodeString(src string, pos int, needEsc bool, validStr bool) (v string, ret int, hasEsc bool) {
    p := NewParserObj(src)
    p.p = pos
    switch val := p.decodeValue(); val.Vt {
    case types.V_STRING:
        str := p.s[val.Iv : p.p-1]
        if validStr && !validate_utf8(str) {
           return "", -int(types.ERR_INVALID_UTF8), false
        }
        /* fast path: no escape sequence */
        if val.Ep == -1 {
            return str, p.p, false
        } else if !needEsc {
            return str, p.p, true
        }
        /* unquote the string */
        out, err := unquote.String(str)
        /* check for errors */
        if err != 0 {
            return "", -int(err), true
        } else {
            return out, p.p, true
        }
    default:
        return "", -int(_ERR_UNSUPPORT_TYPE), false
    }
}
