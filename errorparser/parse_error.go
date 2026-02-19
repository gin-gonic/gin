package errorparser

type ParseError struct {
	ParamName    string
	ErrorType    ParseErrorType
	InitialError error
}

func NewParseError(
	paramName string,
	errorType ParseErrorType,
	initialError error,
) ParseError {
	return ParseError{
		ParamName:    paramName,
		ErrorType:    errorType,
		InitialError: initialError,
	}
}

type ParseErrorType string

const (
	ParseErrorTypeNone       ParseErrorType = ""
	ParseErrorTypeBadInput   ParseErrorType = "bad_input"
	ParseErrorTypeMismatch   ParseErrorType = "type_mismatch"
	ParseErrorTypeValidation ParseErrorType = "validation"
)
