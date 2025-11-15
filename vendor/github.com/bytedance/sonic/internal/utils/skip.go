
package utils

import (
    `runtime`
    `unsafe`

    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)

func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

//go:nocheckptr
func SkipNumber(src string, pos int) (ret int) {
    sp := uintptr(rt.IndexChar(src, pos))
    se := uintptr(rt.IndexChar(src, len(src)))
    if uintptr(sp) >= se {
        return -int(types.ERR_EOF)
    }

    if c := *(*byte)(unsafe.Pointer(sp)); c == '-' {
        sp += 1
    }
    ss := sp

    var pointer bool
    var exponent bool
    var lastIsDigit bool
    var nextNeedDigit = true

    for ; sp < se; sp += uintptr(1) {
        c := *(*byte)(unsafe.Pointer(sp))
        if isDigit(c) {
            lastIsDigit = true
            nextNeedDigit = false
            continue
        } else if nextNeedDigit {
            return -int(types.ERR_INVALID_CHAR)
        } else if c == '.' {
            if !lastIsDigit || pointer || exponent || sp == ss {
                return -int(types.ERR_INVALID_CHAR)
            }
            pointer = true
            lastIsDigit = false
            nextNeedDigit = true
            continue
        } else if c == 'e' || c == 'E' {
            if !lastIsDigit || exponent {
                return -int(types.ERR_INVALID_CHAR)
            }
            if sp == se-1 {
                return -int(types.ERR_EOF)
            }
            exponent = true
            lastIsDigit = false
            nextNeedDigit = false
            continue
        } else if c == '-' || c == '+' {
            if prev := *(*byte)(unsafe.Pointer(sp - 1)); prev != 'e' && prev != 'E' {
                return -int(types.ERR_INVALID_CHAR)
            }
            lastIsDigit = false
            nextNeedDigit = true
            continue
        } else {
            break
        }
    }

    if nextNeedDigit {
        return -int(types.ERR_EOF)
    }

    runtime.KeepAlive(src)
    return int(uintptr(sp) - uintptr((*rt.GoString)(unsafe.Pointer(&src)).Ptr))
}

// Hack: this is used for both checking space and cause friendly compile errors in 32-bit arch.
const _Sonic_Not_Support_32Bit_Arch__Checking_32Bit_Arch_Here = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')


func IsSpace(c byte) bool {
    return (int(1<<c) & _Sonic_Not_Support_32Bit_Arch__Checking_32Bit_Arch_Here) != 0
}
