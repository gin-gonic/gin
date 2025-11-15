package yaml

import (
	"context"
	"io"
	"reflect"

	"github.com/goccy/go-yaml/ast"
)

// DecodeOption functional option type for Decoder
type DecodeOption func(d *Decoder) error

// ReferenceReaders pass to Decoder that reference to anchor defined by passed readers
func ReferenceReaders(readers ...io.Reader) DecodeOption {
	return func(d *Decoder) error {
		d.referenceReaders = append(d.referenceReaders, readers...)
		return nil
	}
}

// ReferenceFiles pass to Decoder that reference to anchor defined by passed files
func ReferenceFiles(files ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceFiles = files
		return nil
	}
}

// ReferenceDirs pass to Decoder that reference to anchor defined by files under the passed dirs
func ReferenceDirs(dirs ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceDirs = dirs
		return nil
	}
}

// RecursiveDir search yaml file recursively from passed dirs by ReferenceDirs option
func RecursiveDir(isRecursive bool) DecodeOption {
	return func(d *Decoder) error {
		d.isRecursiveDir = isRecursive
		return nil
	}
}

// Validator set StructValidator instance to Decoder
func Validator(v StructValidator) DecodeOption {
	return func(d *Decoder) error {
		d.validator = v
		return nil
	}
}

// Strict enable DisallowUnknownField
func Strict() DecodeOption {
	return func(d *Decoder) error {
		d.disallowUnknownField = true
		return nil
	}
}

// DisallowUnknownField causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func DisallowUnknownField() DecodeOption {
	return func(d *Decoder) error {
		d.disallowUnknownField = true
		return nil
	}
}

// AllowDuplicateMapKey ignore syntax error when mapping keys that are duplicates.
func AllowDuplicateMapKey() DecodeOption {
	return func(d *Decoder) error {
		d.allowDuplicateMapKey = true
		return nil
	}
}

// UseOrderedMap can be interpreted as a map,
// and uses MapSlice ( ordered map ) aggressively if there is no type specification
func UseOrderedMap() DecodeOption {
	return func(d *Decoder) error {
		d.useOrderedMap = true
		return nil
	}
}

// UseJSONUnmarshaler if neither `BytesUnmarshaler` nor `InterfaceUnmarshaler` is implemented
// and `UnmashalJSON([]byte)error` is implemented, convert the argument from `YAML` to `JSON` and then call it.
func UseJSONUnmarshaler() DecodeOption {
	return func(d *Decoder) error {
		d.useJSONUnmarshaler = true
		return nil
	}
}

// CustomUnmarshaler overrides any decoding process for the type specified in generics.
//
// NOTE: If RegisterCustomUnmarshaler and CustomUnmarshaler of DecodeOption are specified for the same type,
// the CustomUnmarshaler specified in DecodeOption takes precedence.
func CustomUnmarshaler[T any](unmarshaler func(*T, []byte) error) DecodeOption {
	return func(d *Decoder) error {
		var typ *T
		d.customUnmarshalerMap[reflect.TypeOf(typ)] = func(ctx context.Context, v interface{}, b []byte) error {
			return unmarshaler(v.(*T), b)
		}
		return nil
	}
}

// CustomUnmarshalerContext overrides any decoding process for the type specified in generics.
// Similar to CustomUnmarshaler, but allows passing a context to the unmarshaler function.
func CustomUnmarshalerContext[T any](unmarshaler func(context.Context, *T, []byte) error) DecodeOption {
	return func(d *Decoder) error {
		var typ *T
		d.customUnmarshalerMap[reflect.TypeOf(typ)] = func(ctx context.Context, v interface{}, b []byte) error {
			return unmarshaler(ctx, v.(*T), b)
		}
		return nil
	}
}

// EncodeOption functional option type for Encoder
type EncodeOption func(e *Encoder) error

// Indent change indent number
func Indent(spaces int) EncodeOption {
	return func(e *Encoder) error {
		e.indentNum = spaces
		return nil
	}
}

// IndentSequence causes sequence values to be indented the same value as Indent
func IndentSequence(indent bool) EncodeOption {
	return func(e *Encoder) error {
		e.indentSequence = indent
		return nil
	}
}

// UseSingleQuote determines if single or double quotes should be preferred for strings.
func UseSingleQuote(sq bool) EncodeOption {
	return func(e *Encoder) error {
		e.singleQuote = sq
		return nil
	}
}

// Flow encoding by flow style
func Flow(isFlowStyle bool) EncodeOption {
	return func(e *Encoder) error {
		e.isFlowStyle = isFlowStyle
		return nil
	}
}

// WithSmartAnchor when multiple map values share the same pointer,
// an anchor is automatically assigned to the first occurrence, and aliases are used for subsequent elements.
// The map key name is used as the anchor name by default.
// If key names conflict, a suffix is automatically added to avoid collisions.
// This is an experimental feature and cannot be used simultaneously with anchor tags.
func WithSmartAnchor() EncodeOption {
	return func(e *Encoder) error {
		e.enableSmartAnchor = true
		return nil
	}
}

// UseLiteralStyleIfMultiline causes encoding multiline strings with a literal syntax,
// no matter what characters they include
func UseLiteralStyleIfMultiline(useLiteralStyleIfMultiline bool) EncodeOption {
	return func(e *Encoder) error {
		e.useLiteralStyleIfMultiline = useLiteralStyleIfMultiline
		return nil
	}
}

// JSON encode in JSON format
func JSON() EncodeOption {
	return func(e *Encoder) error {
		e.isJSONStyle = true
		e.isFlowStyle = true
		return nil
	}
}

// MarshalAnchor call back if encoder find an anchor during encoding
func MarshalAnchor(callback func(*ast.AnchorNode, interface{}) error) EncodeOption {
	return func(e *Encoder) error {
		e.anchorCallback = callback
		return nil
	}
}

// UseJSONMarshaler if neither `BytesMarshaler` nor `InterfaceMarshaler`
// nor `encoding.TextMarshaler` is implemented and `MarshalJSON()([]byte, error)` is implemented,
// call `MarshalJSON` to convert the returned `JSON` to `YAML` for processing.
func UseJSONMarshaler() EncodeOption {
	return func(e *Encoder) error {
		e.useJSONMarshaler = true
		return nil
	}
}

// CustomMarshaler overrides any encoding process for the type specified in generics.
//
// NOTE: If type T implements MarshalYAML for pointer receiver, the type specified in CustomMarshaler must be *T.
// If RegisterCustomMarshaler and CustomMarshaler of EncodeOption are specified for the same type,
// the CustomMarshaler specified in EncodeOption takes precedence.
func CustomMarshaler[T any](marshaler func(T) ([]byte, error)) EncodeOption {
	return func(e *Encoder) error {
		var typ T
		e.customMarshalerMap[reflect.TypeOf(typ)] = func(ctx context.Context, v interface{}) ([]byte, error) {
			return marshaler(v.(T))
		}
		return nil
	}
}

// CustomMarshalerContext overrides any encoding process for the type specified in generics.
// Similar to CustomMarshaler, but allows passing a context to the marshaler function.
func CustomMarshalerContext[T any](marshaler func(context.Context, T) ([]byte, error)) EncodeOption {
	return func(e *Encoder) error {
		var typ T
		e.customMarshalerMap[reflect.TypeOf(typ)] = func(ctx context.Context, v interface{}) ([]byte, error) {
			return marshaler(ctx, v.(T))
		}
		return nil
	}
}

// AutoInt automatically converts floating-point numbers to integers when the fractional part is zero.
// For example, a value of 1.0 will be encoded as 1.
func AutoInt() EncodeOption {
	return func(e *Encoder) error {
		e.autoInt = true
		return nil
	}
}

// OmitEmpty behaves in the same way as the interpretation of the omitempty tag in the encoding/json library.
// set on all the fields.
// In the current implementation, the omitempty tag is not implemented in the same way as encoding/json,
// so please specify this option if you expect the same behavior.
func OmitEmpty() EncodeOption {
	return func(e *Encoder) error {
		e.omitEmpty = true
		return nil
	}
}

// OmitZero forces the encoder to assume an `omitzero` struct tag is
// set on all the fields. See `Marshal` commentary for the `omitzero` tag logic.
func OmitZero() EncodeOption {
	return func(e *Encoder) error {
		e.omitZero = true
		return nil
	}
}

// CommentPosition type of the position for comment.
type CommentPosition int

const (
	CommentHeadPosition CommentPosition = CommentPosition(iota)
	CommentLinePosition
	CommentFootPosition
)

func (p CommentPosition) String() string {
	switch p {
	case CommentHeadPosition:
		return "Head"
	case CommentLinePosition:
		return "Line"
	case CommentFootPosition:
		return "Foot"
	default:
		return ""
	}
}

// LineComment create a one-line comment for CommentMap.
func LineComment(text string) *Comment {
	return &Comment{
		Texts:    []string{text},
		Position: CommentLinePosition,
	}
}

// HeadComment create a multiline comment for CommentMap.
func HeadComment(texts ...string) *Comment {
	return &Comment{
		Texts:    texts,
		Position: CommentHeadPosition,
	}
}

// FootComment create a multiline comment for CommentMap.
func FootComment(texts ...string) *Comment {
	return &Comment{
		Texts:    texts,
		Position: CommentFootPosition,
	}
}

// Comment raw data for comment.
type Comment struct {
	Texts    []string
	Position CommentPosition
}

// CommentMap map of the position of the comment and the comment information.
type CommentMap map[string][]*Comment

// WithComment add a comment using the location and text information given in the CommentMap.
func WithComment(cm CommentMap) EncodeOption {
	return func(e *Encoder) error {
		commentMap := map[*Path][]*Comment{}
		for k, v := range cm {
			path, err := PathString(k)
			if err != nil {
				return err
			}
			commentMap[path] = v
		}
		e.commentMap = commentMap
		return nil
	}
}

// CommentToMap apply the position and content of comments in a YAML document to a CommentMap.
func CommentToMap(cm CommentMap) DecodeOption {
	return func(d *Decoder) error {
		if cm == nil {
			return ErrInvalidCommentMapValue
		}
		d.toCommentMap = cm
		return nil
	}
}
