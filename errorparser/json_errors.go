package errorparser

import (
	"encoding/json"
)

func parseJsonDecodeError(err error) (errs []ParseError, match bool) {

	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		return parseJsonUnmarshalTypeError(typeErr), true
	}

	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		return parseJsonSyntaxError(syntaxErr), true
	}

	return nil, false
}

func parseJsonUnmarshalTypeError(err *json.UnmarshalTypeError) (errs []ParseError) {

	errs = []ParseError{}

	item := NewParseError(
		err.Field,
		ParseErrorTypeMismatch,
		err,
	)

	errs = append(errs, item)

	return errs
}

func parseJsonSyntaxError(err *json.SyntaxError) (errs []ParseError) {

	errs = []ParseError{}

	item := NewParseError(
		"",
		ParseErrorTypeBadInput,
		err,
	)

	errs = append(errs, item)

	return errs
}
