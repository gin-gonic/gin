// Copyright 2022 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTOMLBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `toml:"foo"`
	}
	tomlBody := `foo="FOO"`
	err := tomlBinding{}.BindBody([]byte(tomlBody), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}
