
package consts

import (
    `github.com/bytedance/sonic/internal/native/types`
)


const (
    F_use_int64       = 0
    F_disable_urc     = 2
    F_disable_unknown = 3
    F_copy_string     = 4

    F_use_number      = types.B_USE_NUMBER
    F_validate_string = types.B_VALIDATE_STRING
    F_allow_control   = types.B_ALLOW_CONTROL
    F_no_validate_json = types.B_NO_VALIDATE_JSON
    F_case_sensitive = 7
)

type Options uint64

const (
    OptionUseInt64         Options = 1 << F_use_int64
    OptionUseNumber        Options = 1 << F_use_number
    OptionUseUnicodeErrors Options = 1 << F_disable_urc
    OptionDisableUnknown   Options = 1 << F_disable_unknown
    OptionCopyString       Options = 1 << F_copy_string
    OptionValidateString   Options = 1 << F_validate_string
    OptionNoValidateJSON   Options = 1 << F_no_validate_json
    OptionCaseSensitive    Options = 1 << F_case_sensitive
)

const (
	MaxStack = 4096
)