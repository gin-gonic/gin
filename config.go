// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"os"

	"github.com/goccy/go-yaml"
)

// Config represents the YAML configuration structure for Gin.
type Config struct {
	// MaxConns limits the maximum number of concurrent connections.
	// 0 means no limit (default behavior).
	MaxConns int64 `yaml:"max_conns"`
}

// LoadConfig reads configuration from a YAML file and returns an OptionFunc
// that applies the configuration to an Engine.
func LoadConfig(path string) (OptionFunc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return func(e *Engine) {
		if cfg.MaxConns > 0 {
			e.MaxConns = cfg.MaxConns
		}
	}, nil
}

// WithMaxConns creates an OptionFunc that sets the maximum concurrent connections.
func WithMaxConns(n int64) OptionFunc {
	return func(e *Engine) {
		e.MaxConns = n
	}
}
