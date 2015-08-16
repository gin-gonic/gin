// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	baseError := errors.New("test error")
	err := &Error{
		Err:  baseError,
		Type: ErrorTypePrivate,
	}
	assert.Equal(t, err.Error(), baseError.Error())
	assert.Equal(t, err.JSON(), H{"error": baseError.Error()})

	assert.Equal(t, err.SetType(ErrorTypePublic), err)
	assert.Equal(t, err.Type, ErrorTypePublic)

	assert.Equal(t, err.SetMeta("some data"), err)
	assert.Equal(t, err.Meta, "some data")
	assert.Equal(t, err.JSON(), H{
		"error": baseError.Error(),
		"meta":  "some data",
	})

	jsonBytes, _ := json.Marshal(err)
	assert.Equal(t, string(jsonBytes), "{\"error\":\"test error\",\"meta\":\"some data\"}")

	err.SetMeta(H{
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, err.JSON(), H{
		"error":  baseError.Error(),
		"status": "200",
		"data":   "some data",
	})

	err.SetMeta(H{
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	})
	assert.Equal(t, err.JSON(), H{
		"error":  "custom error",
		"status": "200",
		"data":   "some data",
	})
}

func TestErrorSlice(t *testing.T) {
	errs := errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
		{Err: errors.New("second"), Type: ErrorTypePrivate, Meta: "some data"},
		{Err: errors.New("third"), Type: ErrorTypePublic, Meta: H{"status": "400"}},
	}

	assert.Equal(t, errs, errs.ByType(ErrorTypeAny))
	assert.Equal(t, errs.Last().Error(), "third")
	assert.Equal(t, errs.Errors(), []string{"first", "second", "third"})
	assert.Equal(t, errs.ByType(ErrorTypePublic).Errors(), []string{"third"})
	assert.Equal(t, errs.ByType(ErrorTypePrivate).Errors(), []string{"first", "second"})
	assert.Equal(t, errs.ByType(ErrorTypePublic|ErrorTypePrivate).Errors(), []string{"first", "second", "third"})
	assert.Empty(t, errs.ByType(ErrorTypeBind))
	assert.Empty(t, errs.ByType(ErrorTypeBind).String())

	assert.Equal(t, errs.String(), `Error #01: first
Error #02: second
     Meta: some data
Error #03: third
     Meta: map[status:400]
`)
	assert.Equal(t, errs.JSON(), []interface{}{
		H{"error": "first"},
		H{"error": "second", "meta": "some data"},
		H{"error": "third", "status": "400"},
	})
	jsonBytes, _ := json.Marshal(errs)
	assert.Equal(t, string(jsonBytes), "[{\"error\":\"first\"},{\"error\":\"second\",\"meta\":\"some data\"},{\"error\":\"third\",\"status\":\"400\"}]")
	errs = errorMsgs{
		{Err: errors.New("first"), Type: ErrorTypePrivate},
	}
	assert.Equal(t, errs.JSON(), H{"error": "first"})
	jsonBytes, _ = json.Marshal(errs)
	assert.Equal(t, string(jsonBytes), "{\"error\":\"first\"}")

	errs = errorMsgs{}
	assert.Nil(t, errs.Last())
	assert.Nil(t, errs.JSON())
	assert.Empty(t, errs.String())
}
