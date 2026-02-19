package errorparser

import (
	"github.com/go-playground/validator/v10"
)

func parseValidatorError(err error) (errs []ParseError, match bool) {
	if vErr, ok := err.(validator.ValidationErrors); ok {
		return parseValidatorValidationErrors(vErr), true
	}
	return nil, false
}

func parseValidatorValidationErrors(vErr validator.ValidationErrors) (errs []ParseError) {
	fErrs := []validator.FieldError(vErr)
	errs = make([]ParseError, 0, len(fErrs))
	for _, fErr := range fErrs {
		item := NewParseError(
			fErr.Field(),
			ParseErrorTypeValidation,
			fErr,
		)

		errs = append(errs, item)
	}
	return errs
}
