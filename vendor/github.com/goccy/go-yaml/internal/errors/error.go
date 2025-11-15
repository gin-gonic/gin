package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
)

var (
	As  = errors.As
	Is  = errors.Is
	New = errors.New
)

const (
	defaultFormatColor   = false
	defaultIncludeSource = true
)

type Error interface {
	error
	GetToken() *token.Token
	GetMessage() string
	FormatError(bool, bool) string
}

var (
	_ Error = new(SyntaxError)
	_ Error = new(TypeError)
	_ Error = new(OverflowError)
	_ Error = new(DuplicateKeyError)
	_ Error = new(UnknownFieldError)
	_ Error = new(UnexpectedNodeTypeError)
)

type SyntaxError struct {
	Message string
	Token   *token.Token
}

type TypeError struct {
	DstType         reflect.Type
	SrcType         reflect.Type
	StructFieldName *string
	Token           *token.Token
}

type OverflowError struct {
	DstType reflect.Type
	SrcNum  string
	Token   *token.Token
}

type DuplicateKeyError struct {
	Message string
	Token   *token.Token
}

type UnknownFieldError struct {
	Message string
	Token   *token.Token
}

type UnexpectedNodeTypeError struct {
	Actual   ast.NodeType
	Expected ast.NodeType
	Token    *token.Token
}

// ErrSyntax create syntax error instance with message and token
func ErrSyntax(msg string, tk *token.Token) *SyntaxError {
	return &SyntaxError{
		Message: msg,
		Token:   tk,
	}
}

// ErrOverflow creates an overflow error instance with message and a token.
func ErrOverflow(dstType reflect.Type, num string, tk *token.Token) *OverflowError {
	return &OverflowError{
		DstType: dstType,
		SrcNum:  num,
		Token:   tk,
	}
}

// ErrTypeMismatch cerates an type mismatch error instance with token.
func ErrTypeMismatch(dstType, srcType reflect.Type, token *token.Token) *TypeError {
	return &TypeError{
		DstType: dstType,
		SrcType: srcType,
		Token:   token,
	}
}

// ErrDuplicateKey creates an duplicate key error instance with token.
func ErrDuplicateKey(msg string, tk *token.Token) *DuplicateKeyError {
	return &DuplicateKeyError{
		Message: msg,
		Token:   tk,
	}
}

// ErrUnknownField creates an unknown field error instance with token.
func ErrUnknownField(msg string, tk *token.Token) *UnknownFieldError {
	return &UnknownFieldError{
		Message: msg,
		Token:   tk,
	}
}

func ErrUnexpectedNodeType(actual, expected ast.NodeType, tk *token.Token) *UnexpectedNodeTypeError {
	return &UnexpectedNodeTypeError{
		Actual:   actual,
		Expected: expected,
		Token:    tk,
	}
}

func (e *SyntaxError) GetMessage() string {
	return e.Message
}

func (e *SyntaxError) GetToken() *token.Token {
	return e.Token
}

func (e *SyntaxError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *SyntaxError) FormatError(colored, inclSource bool) string {
	return FormatError(e.Message, e.Token, colored, inclSource)
}

func (e *OverflowError) GetMessage() string {
	return e.msg()
}

func (e *OverflowError) GetToken() *token.Token {
	return e.Token
}

func (e *OverflowError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *OverflowError) FormatError(colored, inclSource bool) string {
	return FormatError(e.msg(), e.Token, colored, inclSource)
}

func (e *OverflowError) msg() string {
	return fmt.Sprintf("cannot unmarshal %s into Go value of type %s ( overflow )", e.SrcNum, e.DstType)
}

func (e *TypeError) msg() string {
	if e.StructFieldName != nil {
		return fmt.Sprintf("cannot unmarshal %s into Go struct field %s of type %s", e.SrcType, *e.StructFieldName, e.DstType)
	}
	return fmt.Sprintf("cannot unmarshal %s into Go value of type %s", e.SrcType, e.DstType)
}

func (e *TypeError) GetMessage() string {
	return e.msg()
}

func (e *TypeError) GetToken() *token.Token {
	return e.Token
}

func (e *TypeError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *TypeError) FormatError(colored, inclSource bool) string {
	return FormatError(e.msg(), e.Token, colored, inclSource)
}

func (e *DuplicateKeyError) GetMessage() string {
	return e.Message
}

func (e *DuplicateKeyError) GetToken() *token.Token {
	return e.Token
}

func (e *DuplicateKeyError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *DuplicateKeyError) FormatError(colored, inclSource bool) string {
	return FormatError(e.Message, e.Token, colored, inclSource)
}

func (e *UnknownFieldError) GetMessage() string {
	return e.Message
}

func (e *UnknownFieldError) GetToken() *token.Token {
	return e.Token
}

func (e *UnknownFieldError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *UnknownFieldError) FormatError(colored, inclSource bool) string {
	return FormatError(e.Message, e.Token, colored, inclSource)
}

func (e *UnexpectedNodeTypeError) GetMessage() string {
	return e.msg()
}

func (e *UnexpectedNodeTypeError) GetToken() *token.Token {
	return e.Token
}

func (e *UnexpectedNodeTypeError) Error() string {
	return e.FormatError(defaultFormatColor, defaultIncludeSource)
}

func (e *UnexpectedNodeTypeError) FormatError(colored, inclSource bool) string {
	return FormatError(e.msg(), e.Token, colored, inclSource)
}

func (e *UnexpectedNodeTypeError) msg() string {
	return fmt.Sprintf("%s was used where %s is expected", e.Actual.YAMLName(), e.Expected.YAMLName())
}

func FormatError(errMsg string, token *token.Token, colored, inclSource bool) string {
	var pp printer.Printer
	if token == nil {
		return pp.PrintErrorMessage(errMsg, colored)
	}
	pos := fmt.Sprintf("[%d:%d] ", token.Position.Line, token.Position.Column)
	msg := pp.PrintErrorMessage(fmt.Sprintf("%s%s", pos, errMsg), colored)
	if inclSource {
		msg += "\n" + pp.PrintErrorToken(token, colored)
	}
	return msg
}
