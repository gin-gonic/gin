// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestBodyBindingsWithoutValidation(t *testing.T) {
	type formatRequest struct {
		Name     string `json:"name" xml:"name" yaml:"name" toml:"name"`
		Required string `json:"required" xml:"required" yaml:"required" toml:"required" binding:"required"`
	}

	tests := []struct {
		name   string
		body   string
		binder BindingBodyNoValidate
	}{
		{"JSON", `{"name":"Ada"}`, JSON},
		{"XML", `<request><name>Ada</name></request>`, XML},
		{"YAML", "name: Ada\n", YAML},
		{"TOML", `name = "Ada"`, TOML},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fromRequest formatRequest
			req := requestWithBody(http.MethodPost, "/", tt.body)
			require.NoError(t, tt.binder.BindNoValidate(req, &fromRequest))
			require.Equal(t, "Ada", fromRequest.Name)
			require.Error(t, Validator.ValidateStruct(&fromRequest))

			var fromBody formatRequest
			require.NoError(t, tt.binder.BindBodyNoValidate([]byte(tt.body), &fromBody))
			require.Equal(t, fromRequest, fromBody)
		})
	}

	var validated formatRequest
	require.Error(t, JSON.BindBody([]byte(`{"name":"Ada"}`), &validated))
}

func TestPlainBindingWithoutValidation(t *testing.T) {
	var fromRequest string
	req := requestWithBody(http.MethodPost, "/", "hello")
	require.NoError(t, Plain.BindNoValidate(req, &fromRequest))
	require.Equal(t, "hello", fromRequest)

	var fromBody string
	require.NoError(t, Plain.BindBodyNoValidate([]byte("hello"), &fromBody))
	require.Equal(t, fromRequest, fromBody)
}

func TestProtobufBindingWithoutValidation(t *testing.T) {
	body, err := proto.Marshal(&protoexample.Test{Label: proto.String("yes")})
	require.NoError(t, err)

	var fromRequest protoexample.Test
	req := requestWithBody(http.MethodPost, "/", string(body))
	require.NoError(t, ProtoBuf.BindNoValidate(req, &fromRequest))
	require.Equal(t, "yes", fromRequest.GetLabel())

	var fromBody protoexample.Test
	require.NoError(t, ProtoBuf.BindBodyNoValidate(body, &fromBody))
	require.Equal(t, fromRequest.GetLabel(), fromBody.GetLabel())

	var validated protoexample.Test
	require.NoError(t, ProtoBuf.BindBody(body, &validated))
	require.Equal(t, "yes", validated.GetLabel())
}
