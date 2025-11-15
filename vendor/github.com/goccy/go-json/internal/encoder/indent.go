package encoder

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json/internal/errors"
)

func takeIndentSrcRuntimeContext(src []byte) (*RuntimeContext, []byte) {
	ctx := TakeRuntimeContext()
	buf := ctx.Buf[:0]
	buf = append(append(buf, src...), nul)
	ctx.Buf = buf
	return ctx, buf
}

func Indent(buf *bytes.Buffer, src []byte, prefix, indentStr string) error {
	if len(src) == 0 {
		return errors.ErrUnexpectedEndOfJSON("", 0)
	}

	srcCtx, srcBuf := takeIndentSrcRuntimeContext(src)
	dstCtx := TakeRuntimeContext()
	dst := dstCtx.Buf[:0]

	dst, err := indentAndWrite(buf, dst, srcBuf, prefix, indentStr)
	if err != nil {
		ReleaseRuntimeContext(srcCtx)
		ReleaseRuntimeContext(dstCtx)
		return err
	}
	dstCtx.Buf = dst
	ReleaseRuntimeContext(srcCtx)
	ReleaseRuntimeContext(dstCtx)
	return nil
}

func indentAndWrite(buf *bytes.Buffer, dst []byte, src []byte, prefix, indentStr string) ([]byte, error) {
	dst, err := doIndent(dst, src, prefix, indentStr, false)
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(dst); err != nil {
		return nil, err
	}
	return dst, nil
}

func doIndent(dst, src []byte, prefix, indentStr string, escape bool) ([]byte, error) {
	buf, cursor, err := indentValue(dst, src, 0, 0, []byte(prefix), []byte(indentStr), escape)
	if err != nil {
		return nil, err
	}
	if err := validateEndBuf(src, cursor); err != nil {
		return nil, err
	}
	return buf, nil
}

func indentValue(
	dst []byte,
	src []byte,
	indentNum int,
	cursor int64,
	prefix []byte,
	indentBytes []byte,
	escape bool) ([]byte, int64, error) {
	for {
		switch src[cursor] {
		case ' ', '\t', '\n', '\r':
			cursor++
			continue
		case '{':
			return indentObject(dst, src, indentNum, cursor, prefix, indentBytes, escape)
		case '}':
			return nil, 0, errors.ErrSyntax("unexpected character '}'", cursor)
		case '[':
			return indentArray(dst, src, indentNum, cursor, prefix, indentBytes, escape)
		case ']':
			return nil, 0, errors.ErrSyntax("unexpected character ']'", cursor)
		case '"':
			return compactString(dst, src, cursor, escape)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return compactNumber(dst, src, cursor)
		case 't':
			return compactTrue(dst, src, cursor)
		case 'f':
			return compactFalse(dst, src, cursor)
		case 'n':
			return compactNull(dst, src, cursor)
		default:
			return nil, 0, errors.ErrSyntax(fmt.Sprintf("unexpected character '%c'", src[cursor]), cursor)
		}
	}
}

func indentObject(
	dst []byte,
	src []byte,
	indentNum int,
	cursor int64,
	prefix []byte,
	indentBytes []byte,
	escape bool) ([]byte, int64, error) {
	if src[cursor] == '{' {
		dst = append(dst, '{')
	} else {
		return nil, 0, errors.ErrExpected("expected { character for object value", cursor)
	}
	cursor = skipWhiteSpace(src, cursor+1)
	if src[cursor] == '}' {
		dst = append(dst, '}')
		return dst, cursor + 1, nil
	}
	indentNum++
	var err error
	for {
		dst = append(append(dst, '\n'), prefix...)
		for i := 0; i < indentNum; i++ {
			dst = append(dst, indentBytes...)
		}
		cursor = skipWhiteSpace(src, cursor)
		dst, cursor, err = compactString(dst, src, cursor, escape)
		if err != nil {
			return nil, 0, err
		}
		cursor = skipWhiteSpace(src, cursor)
		if src[cursor] != ':' {
			return nil, 0, errors.ErrSyntax(
				fmt.Sprintf("invalid character '%c' after object key", src[cursor]),
				cursor+1,
			)
		}
		dst = append(dst, ':', ' ')
		dst, cursor, err = indentValue(dst, src, indentNum, cursor+1, prefix, indentBytes, escape)
		if err != nil {
			return nil, 0, err
		}
		cursor = skipWhiteSpace(src, cursor)
		switch src[cursor] {
		case '}':
			dst = append(append(dst, '\n'), prefix...)
			for i := 0; i < indentNum-1; i++ {
				dst = append(dst, indentBytes...)
			}
			dst = append(dst, '}')
			cursor++
			return dst, cursor, nil
		case ',':
			dst = append(dst, ',')
		default:
			return nil, 0, errors.ErrSyntax(
				fmt.Sprintf("invalid character '%c' after object key:value pair", src[cursor]),
				cursor+1,
			)
		}
		cursor++
	}
}

func indentArray(
	dst []byte,
	src []byte,
	indentNum int,
	cursor int64,
	prefix []byte,
	indentBytes []byte,
	escape bool) ([]byte, int64, error) {
	if src[cursor] == '[' {
		dst = append(dst, '[')
	} else {
		return nil, 0, errors.ErrExpected("expected [ character for array value", cursor)
	}
	cursor = skipWhiteSpace(src, cursor+1)
	if src[cursor] == ']' {
		dst = append(dst, ']')
		return dst, cursor + 1, nil
	}
	indentNum++
	var err error
	for {
		dst = append(append(dst, '\n'), prefix...)
		for i := 0; i < indentNum; i++ {
			dst = append(dst, indentBytes...)
		}
		dst, cursor, err = indentValue(dst, src, indentNum, cursor, prefix, indentBytes, escape)
		if err != nil {
			return nil, 0, err
		}
		cursor = skipWhiteSpace(src, cursor)
		switch src[cursor] {
		case ']':
			dst = append(append(dst, '\n'), prefix...)
			for i := 0; i < indentNum-1; i++ {
				dst = append(dst, indentBytes...)
			}
			dst = append(dst, ']')
			cursor++
			return dst, cursor, nil
		case ',':
			dst = append(dst, ',')
		default:
			return nil, 0, errors.ErrSyntax(
				fmt.Sprintf("invalid character '%c' after array value", src[cursor]),
				cursor+1,
			)
		}
		cursor++
	}
}
