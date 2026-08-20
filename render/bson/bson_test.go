// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bson

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type payload struct {
	Foo string `bson:"foo"`
	Bar string `bson:"bar"`
}

func TestRender(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Render(c, http.StatusOK, payload{Foo: "foo", Bar: "bar"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/bson", w.Header().Get("Content-Type"))

	var got payload
	require.NoError(t, bson.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "foo", got.Foo)
	assert.Equal(t, "bar", got.Bar)
}

func TestBinding(t *testing.T) {
	body, err := bson.Marshal(payload{Foo: "foo", Bar: "bar"})
	require.NoError(t, err)

	var out payload
	require.NoError(t, Binding.BindBody(body, &out))
	assert.Equal(t, "foo", out.Foo)
	assert.Equal(t, "bar", out.Bar)
	assert.Equal(t, "bson", Binding.Name())
}

func TestRegistration(t *testing.T) {
	assert.Equal(t, Binding, binding.Default(http.MethodPost, MIMEBSON))
	r, ok := render.Negotiate(MIMEBSON, payload{Foo: "x"})
	assert.True(t, ok)
	assert.NotNil(t, r)
}
