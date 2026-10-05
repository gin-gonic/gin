// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timeFormatReq struct {
	Timestamp time.Time `json:"Timestamp" time_format:"2006-01-02 15:04:05"`
	UnixTime  time.Time `json:"UnixTime" time_format:"unix"`
	Plain     time.Time `json:"Plain"`
}

func TestJSONBindingTimeFormat(t *testing.T) {
	var s timeFormatReq
	err := jsonBinding{}.BindBody([]byte(`{
		"Timestamp": "2001-11-11 11:11:11",
		"UnixTime": 1575528300,
		"Plain": "2001-11-11T11:11:11Z"
	}`), &s)
	require.NoError(t, err)

	// Custom time_format tag is honored.
	assert.Equal(t, time.Date(2001, 11, 11, 11, 11, 11, 0, time.Local), s.Timestamp)

	// unix time_format tag is honored.
	assert.Equal(t, time.Unix(1575528300, 0), s.UnixTime)

	// Fields without a time_format tag keep default RFC3339 parsing.
	assert.Equal(t, time.Date(2001, 11, 11, 11, 11, 11, 0, time.UTC), s.Plain)
}

type timeFormatPtrReq struct {
	Timestamp *time.Time `json:"Timestamp" time_format:"2006-01-02 15:04:05"`
}

func TestJSONBindingTimeFormatPointer(t *testing.T) {
	var s timeFormatPtrReq
	err := jsonBinding{}.BindBody([]byte(`{"Timestamp": "2020-01-02 03:04:05"}`), &s)
	require.NoError(t, err)
	require.NotNil(t, s.Timestamp)
	assert.Equal(t, time.Date(2020, 1, 2, 3, 4, 5, 0, time.Local), *s.Timestamp)
}

type timeFormatNested struct {
	Inner struct {
		When time.Time `json:"when" time_format:"2006-01-02"`
	} `json:"inner"`
}

func TestJSONBindingTimeFormatNested(t *testing.T) {
	var s timeFormatNested
	err := jsonBinding{}.BindBody([]byte(`{"inner": {"when": "2019-12-07"}}`), &s)
	require.NoError(t, err)
	assert.Equal(t, time.Date(2019, 12, 7, 0, 0, 0, 0, time.Local), s.Inner.When)
}