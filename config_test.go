// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	path := filepath.Join("testdata", "config.yaml")
	opt, err := LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, opt)

	engine := New()
	opt(engine)
	assert.Equal(t, int64(100), engine.MaxConns)
}

func TestLoadConfigNotFound(t *testing.T) {
	opt, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
	assert.Nil(t, opt)
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	tmpFile, err := os.CreateTemp("", "invalid_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("max_conns: [invalid")
	require.NoError(t, err)
	tmpFile.Close()

	opt, err := LoadConfig(tmpFile.Name())
	assert.Error(t, err)
	assert.Nil(t, opt)
}

func TestLoadConfigEmpty(t *testing.T) {
	// Create a temporary empty config file
	tmpFile, err := os.CreateTemp("", "empty_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	opt, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, opt)

	engine := New()
	opt(engine)
	assert.Equal(t, int64(0), engine.MaxConns)
}

func TestWithMaxConns(t *testing.T) {
	engine := New(WithMaxConns(50))
	assert.Equal(t, int64(50), engine.MaxConns)
}

func TestWithMaxConnsZero(t *testing.T) {
	engine := New(WithMaxConns(0))
	assert.Equal(t, int64(0), engine.MaxConns)
}
