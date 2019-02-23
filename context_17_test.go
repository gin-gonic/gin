// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that the response is serialized as JSON
// and Content-Type is set to application/json
// and special HTML characters are preserved
func TestContextRenderPureJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.PureJSON(http.StatusCreated, H{"foo": "bar", "html": "<b>"})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"<b>\"}\n", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}
