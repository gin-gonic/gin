package json

import (
	"io"

	"github.com/goccy/go-json/internal/decoder"
	"github.com/goccy/go-json/internal/encoder"
)

type EncodeOption = encoder.Option
type EncodeOptionFunc func(*EncodeOption)

// UnorderedMap doesn't sort when encoding map type.
func UnorderedMap() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.UnorderedMapOption
	}
}

// DisableHTMLEscape disables escaping of HTML characters ( '&', '<', '>' ) when encoding string.
func DisableHTMLEscape() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag &= ^encoder.HTMLEscapeOption
	}
}

// DisableNormalizeUTF8
// By default, when encoding string, UTF8 characters in the range of 0x80 - 0xFF are processed by applying \ufffd for invalid code and escaping for \u2028 and \u2029.
// This option disables this behaviour. You can expect faster speeds by applying this option, but be careful.
// encoding/json implements here: https://github.com/golang/go/blob/6178d25fc0b28724b1b5aec2b1b74fc06d9294c7/src/encoding/json/encode.go#L1067-L1093.
func DisableNormalizeUTF8() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag &= ^encoder.NormalizeUTF8Option
	}
}

// Debug outputs debug information when panic occurs during encoding.
func Debug() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.DebugOption
	}
}

// DebugWith sets the destination to write debug messages.
func DebugWith(w io.Writer) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.DebugOut = w
	}
}

// DebugDOT sets the destination to write opcodes graph.
func DebugDOT(w io.WriteCloser) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.DebugDOTOut = w
	}
}

// Colorize add an identifier for coloring to the string of the encoded result.
func Colorize(scheme *ColorScheme) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.ColorizeOption
		opt.ColorScheme = scheme
	}
}

type DecodeOption = decoder.Option
type DecodeOptionFunc func(*DecodeOption)

// DecodeFieldPriorityFirstWin
// in the default behavior, go-json, like encoding/json,
// will reflect the result of the last evaluation when a field with the same name exists.
// This option allow you to change this behavior.
// this option reflects the result of the first evaluation if a field with the same name exists.
// This behavior has a performance advantage as it allows the subsequent strings to be skipped if all fields have been evaluated.
func DecodeFieldPriorityFirstWin() DecodeOptionFunc {
	return func(opt *DecodeOption) {
		opt.Flags |= decoder.FirstWinOption
	}
}
