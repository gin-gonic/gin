// +build !amd64,!arm64 go1.26 !go1.17 arm64,!go1.20

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
    `bytes`
    `encoding/json`
    `io`
    `reflect`

    `github.com/bytedance/sonic/option`
)

const apiKind = UseStdJSON

type frozenConfig struct {
    Config
}

// Froze convert the Config to API
func (cfg Config) Froze() API {
    api := &frozenConfig{Config: cfg}
    return api
}

func (cfg frozenConfig) marshalOptions(val interface{}, prefix, indent string) ([]byte, error) {
    w := bytes.NewBuffer([]byte{})
    enc := json.NewEncoder(w)
    enc.SetEscapeHTML(cfg.EscapeHTML)
    enc.SetIndent(prefix, indent)
    err := enc.Encode(val)
	out := w.Bytes()

	// json.Encoder always appends '\n' after encoding,
	// which is not same with json.Marshal()
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out, err
}

// Marshal is implemented by sonic
func (cfg frozenConfig) Marshal(val interface{}) ([]byte, error) {
    if !cfg.EscapeHTML {
        return cfg.marshalOptions(val, "", "")
    }
    return json.Marshal(val)
}

// MarshalToString is implemented by sonic
func (cfg frozenConfig) MarshalToString(val interface{}) (string, error) {
    out, err := cfg.Marshal(val)
    return string(out), err
}

// MarshalIndent is implemented by sonic
func (cfg frozenConfig) MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
    if !cfg.EscapeHTML {
        return cfg.marshalOptions(val, prefix, indent)
    }
    return json.MarshalIndent(val, prefix, indent)
}

// UnmarshalFromString is implemented by sonic
func (cfg frozenConfig) UnmarshalFromString(buf string, val interface{}) error {
    r := bytes.NewBufferString(buf)
    dec := json.NewDecoder(r)
    if cfg.UseNumber {
        dec.UseNumber()
    }
    if cfg.DisallowUnknownFields {
        dec.DisallowUnknownFields()
    }
    err := dec.Decode(val)
    if err != nil {
        return err
    }

    // check the trailing chars
    offset := dec.InputOffset()
    if t, err := dec.Token(); !(t == nil && err == io.EOF) {
        return &json.SyntaxError{ Offset: offset}
    }
    return nil
}

// Unmarshal is implemented by sonic
func (cfg frozenConfig) Unmarshal(buf []byte, val interface{}) error {
    return cfg.UnmarshalFromString(string(buf), val)
}

// NewEncoder is implemented by sonic
func (cfg frozenConfig) NewEncoder(writer io.Writer) Encoder {
    enc := json.NewEncoder(writer)
    if !cfg.EscapeHTML {
        enc.SetEscapeHTML(cfg.EscapeHTML)
    }
    return enc
}

// NewDecoder is implemented by sonic
func (cfg frozenConfig) NewDecoder(reader io.Reader) Decoder {
    dec := json.NewDecoder(reader)
    if cfg.UseNumber {
        dec.UseNumber()
    }
    if cfg.DisallowUnknownFields {
        dec.DisallowUnknownFields()
    }
    return dec
}

// Valid is implemented by sonic
func (cfg frozenConfig) Valid(data []byte) bool {
    return json.Valid(data)
}

// Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
// order to reduce the first-hit latency at **amd64** Arch.
// Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
// a compile option to set the depth of recursive compile for the nested struct type.
// * This is the none implement for !amd64.
// It will be useful for someone who develop with !amd64 arch,like Mac M1.
func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    return nil
}

