package json

import (
	"github.com/goccy/go-json/internal/errors"
)

// Before Go 1.2, an InvalidUTF8Error was returned by Marshal when
// attempting to encode a string value with invalid UTF-8 sequences.
// As of Go 1.2, Marshal instead coerces the string to valid UTF-8 by
// replacing invalid bytes with the Unicode replacement rune U+FFFD.
//
// Deprecated: No longer used; kept for compatibility.
type InvalidUTF8Error = errors.InvalidUTF8Error

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError = errors.InvalidUnmarshalError

// A MarshalerError represents an error from calling a MarshalJSON or MarshalText method.
type MarshalerError = errors.MarshalerError

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError = errors.SyntaxError

// An UnmarshalFieldError describes a JSON object key that
// led to an unexported (and therefore unwritable) struct field.
//
// Deprecated: No longer used; kept for compatibility.
type UnmarshalFieldError = errors.UnmarshalFieldError

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError = errors.UnmarshalTypeError

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError = errors.UnsupportedTypeError

type UnsupportedValueError = errors.UnsupportedValueError

type PathError = errors.PathError
