package errorparser

func ParseBindError(err error) (errs []ParseError, match bool) {

	if errs, ok := parseValidatorError(err); ok {
		return errs, true
	}

	if errs, ok := parseJsonDecodeError(err); ok {
		return errs, true
	}

	// todo: protobuf
	// todo: xml
	// todo: yaml

	return nil, false
}
