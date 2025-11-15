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

package sonic

import (
    `io`
    `reflect`

    `github.com/bytedance/sonic/decoder`
    `github.com/bytedance/sonic/encoder`
    `github.com/bytedance/sonic/option`
    `github.com/bytedance/sonic/internal/rt`
)

const apiKind = UseSonicJSON

type frozenConfig struct {
    Config
    encoderOpts encoder.Options
    decoderOpts decoder.Options
}

// Froze convert the Config to API
func (cfg Config) Froze() API {
    api := &frozenConfig{Config: cfg}

    // configure encoder options:
    if cfg.EscapeHTML {
        api.encoderOpts |= encoder.EscapeHTML
    }
    if cfg.SortMapKeys {
        api.encoderOpts |= encoder.SortMapKeys
    }
    if cfg.CompactMarshaler {
        api.encoderOpts |= encoder.CompactMarshaler
    }
    if cfg.NoQuoteTextMarshaler {
        api.encoderOpts |= encoder.NoQuoteTextMarshaler
    }
    if cfg.NoNullSliceOrMap {
        api.encoderOpts |= encoder.NoNullSliceOrMap
    }
    if cfg.ValidateString {
        api.encoderOpts |= encoder.ValidateString
    }
    if cfg.NoValidateJSONMarshaler {
        api.encoderOpts |= encoder.NoValidateJSONMarshaler
    }
    if cfg.NoEncoderNewline {
        api.encoderOpts |= encoder.NoEncoderNewline
    }
    if cfg.EncodeNullForInfOrNan {
        api.encoderOpts |= encoder.EncodeNullForInfOrNan
    }

    // configure decoder options:
    if cfg.NoValidateJSONSkip {
        api.decoderOpts |= decoder.OptionNoValidateJSON
    }
    if cfg.UseInt64 {
        api.decoderOpts |= decoder.OptionUseInt64
    }
    if cfg.UseNumber {
        api.decoderOpts |= decoder.OptionUseNumber
    }
    if cfg.DisallowUnknownFields {
        api.decoderOpts |= decoder.OptionDisableUnknown
    }
    if cfg.CopyString {
        api.decoderOpts |= decoder.OptionCopyString
    }
    if cfg.ValidateString {
        api.decoderOpts |= decoder.OptionValidateString
    }
    if cfg.CaseSensitive {
        api.decoderOpts |= decoder.OptionCaseSensitive
    }
    return api
}

// Marshal is implemented by sonic
func (cfg frozenConfig) Marshal(val interface{}) ([]byte, error) {
    return encoder.Encode(val, cfg.encoderOpts)
}

// MarshalToString is implemented by sonic
func (cfg frozenConfig) MarshalToString(val interface{}) (string, error) {
    buf, err := encoder.Encode(val, cfg.encoderOpts)
    return rt.Mem2Str(buf), err
}

// MarshalIndent is implemented by sonic
func (cfg frozenConfig) MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
    return encoder.EncodeIndented(val, prefix, indent, cfg.encoderOpts)
}

// UnmarshalFromString is implemented by sonic
func (cfg frozenConfig) UnmarshalFromString(buf string, val interface{}) error {
    dec := decoder.NewDecoder(buf)
    dec.SetOptions(cfg.decoderOpts)
    err := dec.Decode(val)

    /* check for errors */
    if err != nil {
        return err
    }

    return dec.CheckTrailings()
}

// Unmarshal is implemented by sonic
func (cfg frozenConfig) Unmarshal(buf []byte, val interface{}) error {
    return cfg.UnmarshalFromString(string(buf), val)
}

// NewEncoder is implemented by sonic
func (cfg frozenConfig) NewEncoder(writer io.Writer) Encoder {
    enc := encoder.NewStreamEncoder(writer)
    enc.Opts = cfg.encoderOpts
    return enc
}

// NewDecoder is implemented by sonic
func (cfg frozenConfig) NewDecoder(reader io.Reader) Decoder {
    dec := decoder.NewStreamDecoder(reader)
    dec.SetOptions(cfg.decoderOpts)
    return dec
}

// Valid is implemented by sonic
func (cfg frozenConfig) Valid(data []byte) bool {
    ok, _ := encoder.Valid(data)
    return ok
}

// Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
// order to reduce the first-hit latency.
//
// Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
// a compile option to set the depth of recursive compile for the nested struct type.
func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    if err := encoder.Pretouch(vt, opts...); err != nil {
        return err
    } 
    if err := decoder.Pretouch(vt, opts...); err != nil {
        return err
    }
    // to pretouch the corresponding pointer type as well
    if vt.Kind() == reflect.Ptr {
        vt = vt.Elem()
    } else {
        vt = reflect.PtrTo(vt)
    }
    if err := encoder.Pretouch(vt, opts...); err != nil {
        return err
    } 
    if err := decoder.Pretouch(vt, opts...); err != nil {
        return err
    }
    return nil
}
