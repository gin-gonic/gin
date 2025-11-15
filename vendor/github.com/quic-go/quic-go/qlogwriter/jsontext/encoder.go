// Package jsontext provides a fast JSON encoder providing only the necessary features
// for qlog encoding. No efforts are made to add any features beyond qlog's requirements.
//
// The API aims to be compatible with the standard library's encoding/json/jsontext package.
package jsontext

import (
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

type kind uint8

const (
	kindString kind = iota
	kindInt
	kindUint
	kindFloat
	kindBool
	kindNull
	kindObjectStart
	kindObjectEnd
	kindArrayStart
	kindArrayEnd
)

// Token represents a JSON token.
type Token struct {
	kind kind
	str  string
	i64  int64
	u64  uint64
	f64  float64
	b    bool
}

// String creates a string token.
func String(s string) Token {
	return Token{kind: kindString, str: s}
}

// Int creates an int token.
func Int(i int64) Token {
	return Token{kind: kindInt, i64: i}
}

// Uint creates a uint token.
func Uint(u uint64) Token {
	return Token{kind: kindUint, u64: u}
}

// Float creates a float token.
func Float(f float64) Token {
	return Token{kind: kindFloat, f64: f}
}

// Bool creates a bool token.
func Bool(b bool) Token {
	return Token{kind: kindBool, b: b}
}

// Null is a null token.
var Null Token = Token{kind: kindNull}

// BeginObject is the begin object token.
var BeginObject Token = Token{kind: kindObjectStart}

// EndObject is the end object token.
var EndObject Token = Token{kind: kindObjectEnd}

// BeginArray is the begin array token.
var BeginArray Token = Token{kind: kindArrayStart}

// EndArray is the end array token.
var EndArray Token = Token{kind: kindArrayEnd}

// True is a true token.
var True Token = Bool(true)

// False is a false token.
var False Token = Bool(false)

var hexDigits = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

var (
	commaByte       = []byte(",")
	quoteByte       = []byte(`"`)
	colonByte       = []byte(":")
	trueByte        = []byte("true")
	falseByte       = []byte("false")
	nullByte        = []byte("null")
	openObjectByte  = []byte("{")
	closeObjectByte = []byte("}")
	openArrayByte   = []byte("[")
	closeArrayByte  = []byte("]")
	newlineByte     = []byte("\n")
	escapeQuote     = []byte(`\"`)
	escapeBackslash = []byte(`\\`)
	escapeBackspace = []byte(`\b`)
	escapeFormfeed  = []byte(`\f`)
	escapeNewline   = []byte(`\n`)
	escapeCarriage  = []byte(`\r`)
	escapeTab       = []byte(`\t`)
	escapeUnicode   = []byte(`\u00`)
)

type context struct {
	isObject   bool
	needsComma bool
	expectKey  bool
}

// Encoder encodes JSON to an io.Writer.
type Encoder struct {
	w     io.Writer
	buf   [64]byte // scratch buffer for number formatting
	stack []context
}

// NewEncoder creates a new Encoder.
func NewEncoder(w io.Writer) *Encoder {
	stack := make([]context, 0, 8)
	stack = append(stack, context{isObject: false, needsComma: false, expectKey: false})
	return &Encoder{
		w:     w,
		stack: stack,
	}
}

// WriteToken writes a token to the encoder.
func (e *Encoder) WriteToken(t Token) error {
	if len(e.stack) == 0 {
		return fmt.Errorf("empty stack")
	}
	curr := &e.stack[len(e.stack)-1]
	isClosing := t.kind == kindObjectEnd || t.kind == kindArrayEnd
	if !isClosing && curr.needsComma {
		if _, err := e.w.Write(commaByte); err != nil {
			return err
		}
		curr.needsComma = false
	}
	var err error
	switch t.kind {
	case kindString:
		data := stringToBytes(t.str)
		needsEscape := false
		for _, b := range data {
			if b == '"' || b == '\\' || b < 0x20 {
				needsEscape = true
				break
			}
		}
		if !needsEscape {
			if _, err = e.w.Write(quoteByte); err != nil {
				return err
			}
			if _, err = e.w.Write(data); err != nil {
				return err
			}
			if _, err = e.w.Write(quoteByte); err != nil {
				return err
			}
		} else {
			if _, err = e.w.Write(quoteByte); err != nil {
				return err
			}
			for i := 0; i < len(t.str); i++ {
				c := t.str[i]
				switch c {
				case '"':
					if _, err = e.w.Write(escapeQuote); err != nil {
						return err
					}
				case '\\':
					if _, err = e.w.Write(escapeBackslash); err != nil {
						return err
					}
				case '\b':
					if _, err = e.w.Write(escapeBackspace); err != nil {
						return err
					}
				case '\f':
					if _, err = e.w.Write(escapeFormfeed); err != nil {
						return err
					}
				case '\n':
					if _, err = e.w.Write(escapeNewline); err != nil {
						return err
					}
				case '\r':
					if _, err = e.w.Write(escapeCarriage); err != nil {
						return err
					}
				case '\t':
					if _, err = e.w.Write(escapeTab); err != nil {
						return err
					}
				default:
					if c < 0x20 {
						if _, err = e.w.Write(escapeUnicode); err != nil {
							return err
						}
						if _, err = e.w.Write([]byte{hexDigits[c>>4], hexDigits[c&0xf]}); err != nil {
							return err
						}
					} else {
						if _, err = e.w.Write([]byte{c}); err != nil {
							return err
						}
					}
				}
			}
			if _, err = e.w.Write(quoteByte); err != nil {
				return err
			}
		}
		if curr.isObject {
			if curr.expectKey {
				// key
				if _, err = e.w.Write(colonByte); err != nil {
					return err
				}
				curr.expectKey = false
				return nil // do not call afterValue for keys
			} else {
				// value
				e.afterValue()
			}
		} else {
			e.afterValue()
		}
	case kindInt:
		b := strconv.AppendInt(e.buf[:0], t.i64, 10)
		if _, err = e.w.Write(b); err != nil {
			return err
		}
		e.afterValue()
	case kindUint:
		b := strconv.AppendUint(e.buf[:0], t.u64, 10)
		if _, err = e.w.Write(b); err != nil {
			return err
		}
		e.afterValue()
	case kindFloat:
		b := strconv.AppendFloat(e.buf[:0], t.f64, 'g', -1, 64)
		if _, err = e.w.Write(b); err != nil {
			return err
		}
		e.afterValue()
	case kindBool:
		if t.b {
			if _, err = e.w.Write(trueByte); err != nil {
				return err
			}
		} else {
			if _, err = e.w.Write(falseByte); err != nil {
				return err
			}
		}
		e.afterValue()
	case kindNull:
		if _, err = e.w.Write(nullByte); err != nil {
			return err
		}
		e.afterValue()
	case kindObjectStart:
		if _, err = e.w.Write(openObjectByte); err != nil {
			return err
		}
		e.stack = append(e.stack, context{isObject: true, needsComma: false, expectKey: true})
		return nil
	case kindObjectEnd:
		if _, err = e.w.Write(closeObjectByte); err != nil {
			return err
		}
		e.stack = e.stack[:len(e.stack)-1]
		e.afterValue()
		if len(e.stack) == 1 {
			if _, err = e.w.Write(newlineByte); err != nil {
				return err
			}
		}
		return nil
	case kindArrayStart:
		if _, err = e.w.Write(openArrayByte); err != nil {
			return err
		}
		e.stack = append(e.stack, context{isObject: false, needsComma: false, expectKey: false})
		return nil
	case kindArrayEnd:
		if _, err = e.w.Write(closeArrayByte); err != nil {
			return err
		}
		e.stack = e.stack[:len(e.stack)-1]
		e.afterValue()
		if len(e.stack) == 1 {
			if _, err = e.w.Write(newlineByte); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown token kind")
	}
	return err
}

// afterValue updates the state after encoding a value
func (e *Encoder) afterValue() {
	if len(e.stack) > 1 {
		curr := &e.stack[len(e.stack)-1]
		curr.needsComma = true
		if curr.isObject {
			curr.expectKey = true
		}
	}
}

func stringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
