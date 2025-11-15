package encoder

import (
	"context"
	"io"
)

type OptionFlag uint8

const (
	HTMLEscapeOption OptionFlag = 1 << iota
	IndentOption
	UnorderedMapOption
	DebugOption
	ColorizeOption
	ContextOption
	NormalizeUTF8Option
	FieldQueryOption
)

type Option struct {
	Flag        OptionFlag
	ColorScheme *ColorScheme
	Context     context.Context
	DebugOut    io.Writer
	DebugDOTOut io.WriteCloser
}

type EncodeFormat struct {
	Header string
	Footer string
}

type EncodeFormatScheme struct {
	Int       EncodeFormat
	Uint      EncodeFormat
	Float     EncodeFormat
	Bool      EncodeFormat
	String    EncodeFormat
	Binary    EncodeFormat
	ObjectKey EncodeFormat
	Null      EncodeFormat
}

type (
	ColorScheme = EncodeFormatScheme
	ColorFormat = EncodeFormat
)
