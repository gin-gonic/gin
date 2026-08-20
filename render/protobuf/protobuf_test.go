// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protobuf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/gin-gonic/gin/testdata/protoexample"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func sample() *protoexample.Test {
	return &protoexample.Test{
		Label: proto.String("yes"),
		Reps:  []int64{1, 2, 3},
	}
}

func TestRender(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Render(c, http.StatusOK, sample())

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))

	var got protoexample.Test
	require.NoError(t, proto.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "yes", got.GetLabel())
	assert.Equal(t, []int64{1, 2, 3}, got.GetReps())
}

func TestBinding(t *testing.T) {
	body, err := proto.Marshal(sample())
	require.NoError(t, err)

	var out protoexample.Test
	require.NoError(t, Binding.BindBody(body, &out))
	assert.Equal(t, "yes", out.GetLabel())
	assert.Equal(t, "protobuf", Binding.Name())
}

func TestBindingNotProtoMessage(t *testing.T) {
	err := Binding.BindBody([]byte("x"), &struct{ Foo string }{})
	require.Error(t, err)
}

func TestRegistration(t *testing.T) {
	assert.Equal(t, Binding, binding.Default(http.MethodPost, MIMEPROTOBUF))
	r, ok := render.Negotiate(MIMEPROTOBUF, sample())
	assert.True(t, ok)
	assert.NotNil(t, r)
}
