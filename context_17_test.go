// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package gin

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestContextHTTPContext(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	c.Request = req.WithContext(ctx)

	assert.NoError(t, c.Err())
	assert.NotNil(t, c.Done())
	select {
	case <-c.Done():
		assert.Fail(t, "context should not be canceled")
	default:
	}

	ti, ok := c.Deadline()
	assert.Equal(t, ti, time.Time{})
	assert.False(t, ok)
	assert.Equal(t, c.Value(0), c.Request)

	cancelFunc()
	assert.NotNil(t, c.Done())
	select {
	case <-c.Done():
	default:
		assert.Fail(t, "context should be canceled")
	}
}

func TestContextHTTPContextWithDeadline(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	location, _ := time.LoadLocation("Europe/Paris")
	assert.NotNil(t, location)
	date := time.Date(2031, 12, 27, 16, 00, 00, 00, location)
	ctx, cancelFunc := context.WithDeadline(context.Background(), date)
	defer cancelFunc()
	c.Request = req.WithContext(ctx)

	assert.NoError(t, c.Err())

	ti, ok := c.Deadline()
	assert.Equal(t, ti, date)
	assert.True(t, ok)
}
