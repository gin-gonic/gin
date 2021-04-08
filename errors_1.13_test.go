//go:build go1.13
// +build go1.13

package gin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestErr string

func (e TestErr) Error() string { return string(e) }

// TestErrorUnwrap tests the behavior of gin.Error with "errors.Is()" and "errors.As()".
// "errors.Is()" and "errors.As()" have been added to the standard library in go 1.13,
// hence the "// +build go1.13" directive at the beginning of this file.
func TestErrorUnwrap(t *testing.T) {
	innerErr := TestErr("somme error")

	// 2 layers of wrapping : use 'fmt.Errorf("%w")' to wrap a gin.Error{}, which itself wraps innerErr
	err := fmt.Errorf("wrapped: %w", &Error{
		Err:  innerErr,
		Type: ErrorTypeAny,
	})

	// check that 'errors.Is()' and 'errors.As()' behave as expected :
	assert.True(t, errors.Is(err, innerErr))
	var testErr TestErr
	assert.True(t, errors.As(err, &testErr))
}
