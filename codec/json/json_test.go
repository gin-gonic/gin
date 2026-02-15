// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONMarshal(t *testing.T) {
	data := map[string]string{"key": "value"}
	result, err := API.Marshal(data)
	require.NoError(t, err)
	assert.JSONEq(t, `{"key":"value"}`, string(result))
}

func TestJSONUnmarshal(t *testing.T) {
	var data map[string]string
	err := API.Unmarshal([]byte(`{"key":"value"}`), &data)
	require.NoError(t, err)
	assert.Equal(t, "value", data["key"])
}

func TestJSONMarshalIndent(t *testing.T) {
	data := map[string]string{"key": "value"}
	result, err := API.MarshalIndent(data, "", "  ")
	require.NoError(t, err)
	assert.Contains(t, string(result), `"key": "value"`)
}

func TestJSONNewEncoder(t *testing.T) {
	var buf bytes.Buffer
	encoder := API.NewEncoder(&buf)
	require.NotNil(t, encoder)
	err := encoder.Encode(map[string]string{"key": "value"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"key":"value"}`, buf.String())
}

func TestJSONNewDecoder(t *testing.T) {
	buf := bytes.NewBufferString(`{"key":"value"}`)
	decoder := API.NewDecoder(buf)
	require.NotNil(t, decoder)
	var data map[string]string
	err := decoder.Decode(&data)
	require.NoError(t, err)
	assert.Equal(t, "value", data["key"])
}
