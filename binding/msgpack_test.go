// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugorji/go/codec"
)

func TestMsgpackBindingBindBody(t *testing.T) {
	type teststruct struct {
		Foo string `msgpack:"foo"`
	}
	var s teststruct
	err := msgpackBinding{}.BindBody(msgpackBody(t, teststruct{"FOO"}), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}

func TestMsgpackBindingNoValidate(t *testing.T) {
	type teststruct struct {
		Foo      string `msgpack:"foo"`
		Required string `msgpack:"required" binding:"required"`
	}
	body := msgpackBody(t, teststruct{Foo: "FOO"})

	var fromRequest teststruct
	req := requestWithBody(http.MethodPost, "/", string(body))
	require.NoError(t, msgpackBinding{}.BindNoValidate(req, &fromRequest))
	assert.Equal(t, "FOO", fromRequest.Foo)
	require.Error(t, Validator.ValidateStruct(&fromRequest))

	var fromBody teststruct
	require.NoError(t, msgpackBinding{}.BindBodyNoValidate(body, &fromBody))
	assert.Equal(t, fromRequest, fromBody)
}

func msgpackBody(t *testing.T, obj any) []byte {
	var bs bytes.Buffer
	h := &codec.MsgpackHandle{}
	err := codec.NewEncoder(&bs, h).Encode(obj)
	require.NoError(t, err)
	return bs.Bytes()
}
