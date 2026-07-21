// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"testing"

	"github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblemMarshalJSON(t *testing.T) {
	p := Problem{
		Type:     "https://example.com/probs/out-of-credit",
		Title:    "You do not have enough credit.",
		Status:   403,
		Detail:   "Your current balance is 30, but that costs 50.",
		Instance: "/account/12345/msgs/abc",
	}

	jsonBytes, err := json.API.Marshal(p)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"type": "https://example.com/probs/out-of-credit",
		"title": "You do not have enough credit.",
		"status": 403,
		"detail": "Your current balance is 30, but that costs 50.",
		"instance": "/account/12345/msgs/abc"
	}`, string(jsonBytes))
}

func TestProblemMarshalJSONOmitsEmptyMembers(t *testing.T) {
	jsonBytes, err := json.API.Marshal(Problem{Status: 404})

	require.NoError(t, err)
	assert.JSONEq(t, `{"status":404}`, string(jsonBytes))
}

func TestProblemMarshalJSONExtensions(t *testing.T) {
	p := Problem{
		Status: 403,
		Detail: "Your current balance is 30, but that costs 50.",
		Extensions: map[string]any{
			"balance": 30,
			"status":  "extension members must not override standard members",
		},
	}

	jsonBytes, err := json.API.Marshal(p)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"status": 403,
		"detail": "Your current balance is 30, but that costs 50.",
		"balance": 30
	}`, string(jsonBytes))
}
