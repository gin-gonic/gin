package errorparser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestParseBindError(t *testing.T) {

	_, ok := ParseBindError(fmt.Errorf("not match"))
	assert.False(t, ok)

	_, ok = ParseBindError(validator.ValidationErrors([]validator.FieldError{}))
	assert.True(t, ok)

	_, ok = ParseBindError(&json.SyntaxError{})
	assert.True(t, ok)

}
