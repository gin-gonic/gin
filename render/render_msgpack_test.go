// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

package render

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugorji/go/codec"
)

func TestRenderMsgPack(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
	}

	(MsgPack{data}).WriteContentType(w)
	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))

	err := (MsgPack{data}).Render(w)

	require.NoError(t, err)

	h := new(codec.MsgpackHandle)
	assert.NotNil(t, h)
	buf := bytes.NewBuffer([]byte{})
	assert.NotNil(t, buf)
	err = codec.NewEncoder(buf, h).Encode(data)

	require.NoError(t, err)
	assert.Equal(t, w.Body.String(), buf.String())
	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestWriteMsgPack(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
		"num": 42,
	}

	err := WriteMsgPack(w, data)
	require.NoError(t, err)

	assert.Equal(t, "application/msgpack; charset=utf-8", w.Header().Get("Content-Type"))

	// Verify the encoded data is correct
	h := new(codec.MsgpackHandle)
	buf := bytes.NewBuffer([]byte{})
	err = codec.NewEncoder(buf, h).Encode(data)
	require.NoError(t, err)

	assert.Equal(t, buf.String(), w.Body.String())
}

type failWriter struct {
	*httptest.ResponseRecorder
}

func (w *failWriter) Write(data []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestRenderMsgPackError(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]any{
		"foo": "bar",
	}

	err := (MsgPack{data}).Render(&failWriter{w})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "write error")
}
