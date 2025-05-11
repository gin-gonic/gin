// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin/internal/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	baseError := errors.New("test error")
	err := &Error{
		Err:  baseError,
		Type: ErrorTypePrivate,
	}
	assert.Equal(t, err.Error(), baseError.Error())
	assert.Equal(t, H{"error": baseError.Error()}, err.JSON())

	assert.Equal(t, err.SetType(ErrorTypePublic), err)
	assert.Equal(t, ErrorTypePublic, err.Type)

	assert.Equal(t, err.SetMeta("some data"), err)
	assert.Equal(t, "some data", err.Meta)
	assert.Equal(t, H{
		"error": baseError.Error(),
		"meta":  "some data",
	}, err.JSON())

	jsonBytes, _ := json.Marshal(err)
	assert.Equal(t, "{\"error\":\"test error\",\"meta\":\"some data\"}", string(jsonBytes))

	err.SetMeta(H{ //nolint: errcheck
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, H{
		"error":  baseError.Error(),
		"status": "200",
		"data":   "some data",
	}, err.JSON())

	err.SetMeta(H{ //nolint: errcheck
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, H{
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	}, err.JSON())

	type customError struct {
		status string
		data   string
	}
	err.SetMeta(customError{status: "200", data: "other data"}) //nolint: errcheck
	assert.Equal(t, customError{status: "200", data: "other data"}, err.JSON())
}

func TestErrorSlice(t *testing.T) {
	errs := errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
		{Err: errors.New("second"), Type: ErrorTypePrivate, Meta: "some data"},
		{Err: errors.New("third"), Type: ErrorTypePublic, Meta: H{"status": "400"}},
	}

	assert.Equal(t, errs, errs.ByType(ErrorTypeAny))
	assert.Equal(t, "third", errs.Last().Error())
	assert.Equal(t, []string{"first", "second", "third"}, errs.Errors())
	assert.Equal(t, []string{"third"}, errs.ByType(ErrorTypePublic).Errors())
	assert.Equal(t, []string{"first", "second"}, errs.ByType(ErrorTypePrivate).Errors())
	assert.Equal(t, []string{"first", "second", "third"}, errs.ByType(ErrorTypePublic|ErrorTypePrivate).Errors())
	assert.Empty(t, errs.ByType(ErrorTypeBind))
	assert.Empty(t, errs.ByType(ErrorTypeBind).String())

	assert.Equal(t, `Error #01: first
Error #02: second
     Meta: some data
Error #03: third
     Meta: map[status:400]
`, errs.String())
	assert.Equal(t, []any{
		H{"error": "first"},
		H{"error": "second", "meta": "some data"},
		H{"error": "third", "status": "400"},
	}, errs.JSON())
	jsonBytes, _ := json.Marshal(errs)
	assert.Equal(t, "[{\"error\":\"first\"},{\"error\":\"second\",\"meta\":\"some data\"},{\"error\":\"third\",\"status\":\"400\"}]", string(jsonBytes))
	errs = errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
	}
	assert.Equal(t, H{"error": "first"}, errs.JSON())
	jsonBytes, _ = json.Marshal(errs)
	assert.Equal(t, "{\"error\":\"first\"}", string(jsonBytes))

	errs = errorMsgs{}
	assert.Nil(t, errs.Last())
	assert.Nil(t, errs.JSON())
	assert.Empty(t, errs.String())
}

type TestErr string

func (e TestErr) Error() string { return string(e) }

// TestErrorUnwrap tests the behavior of gin.Error with "errors.Is()" and "errors.As()".
// "errors.Is()" and "errors.As()" have been added to the standard library in go 1.13.
func TestErrorUnwrap(t *testing.T) {
	innerErr := TestErr("some error")

	// 2 layers of wrapping : use 'fmt.Errorf("%w")' to wrap a gin.Error{}, which itself wraps innerErr
	err := fmt.Errorf("wrapped: %w", &Error{
		Err:  innerErr,
		Type: ErrorTypeAny,
	})

	// check that 'errors.Is()' and 'errors.As()' behave as expected :
	require.ErrorIs(t, err, innerErr)
	var testErr TestErr
	require.ErrorAs(t, err, &testErr)
}
