// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvroBindingBindBody(t *testing.T) {
	var s struct {
		A int64  `avro:"a"`
		B string `avro:"b"`
	}
	schema := `{
	"type": "record",
	"name": "test",
	"fields" : [
		{"name": "a", "type": "long"},
	    {"name": "b", "type": "string"}
		]
	}`
	avroBody := `{"a": 27, "b": "foo"}`

	err := avroBinding{}.BindBody([]byte(avroBody), schema, &s)
	require.NoError(t, err)
	assert.Equal(t, "foo", s.B)
}
