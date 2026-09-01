// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package json_test is an external test package on purpose. Each codec
// installs API from an init function, and package-level variables of an
// in-package test would be initialised before that init runs, leaving API
// nil. An importing package is guaranteed to observe a fully initialised
// codec/json, so tests here see API exactly as the rest of gin does.
package json_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin/codec/json"
)

// TestAPIInstalled guards the build constraints across the codec files. Every
// supported combination of tags, toolchain and GOEXPERIMENT must leave exactly
// one codec selected; if the constraints ever fail to overlap, no init runs and
// API stays nil, which builds cleanly and only panics at the first render.
func TestAPIInstalled(t *testing.T) {
	require.NotNil(t, json.API, "no codec init ran: build constraints do not cover this configuration")
	assert.NotEmpty(t, json.Package)
}

type payload struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

func TestMarshal(t *testing.T) {
	got, err := json.API.Marshal(payload{Name: "gin", N: 1})
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"gin","n":1}`, string(got))
}

// TestMarshalEscapesHTML pins the HTML escaping that Context.JSON relies on to
// stay safe when a response is embedded in a document. PureJSON is the only
// render that opts out, via Encoder.SetEscapeHTML.
func TestMarshalEscapesHTML(t *testing.T) {
	got, err := json.API.Marshal(payload{Name: "<script>&"})
	require.NoError(t, err)
	assert.Contains(t, string(got), `\u003cscript\u003e\u0026`)
	assert.NotContains(t, string(got), "<script>")
}

func TestUnmarshal(t *testing.T) {
	var got payload
	require.NoError(t, json.API.Unmarshal([]byte(`{"name":"gin","n":1}`), &got))
	assert.Equal(t, payload{Name: "gin", N: 1}, got)
}

// TestUnmarshalMatchesNamesCaseInsensitively pins v1's case-insensitive member
// matching. encoding/json/v2 matches case-sensitively by default, which would
// silently leave fields zeroed with a nil error rather than reporting a problem.
func TestUnmarshalMatchesNamesCaseInsensitively(t *testing.T) {
	var got payload
	require.NoError(t, json.API.Unmarshal([]byte(`{"NAME":"gin","N":1}`), &got))
	assert.Equal(t, payload{Name: "gin", N: 1}, got)
}

// TestMarshalNilSliceAndMapAsNull pins v1's rendering of nil slices and maps.
// encoding/json/v2 renders them as [] and {} by default, a visible change to
// every response body carrying an unset slice or map.
func TestMarshalNilSliceAndMapAsNull(t *testing.T) {
	got, err := json.API.Marshal(struct {
		S []string          `json:"s"`
		M map[string]string `json:"m"`
	}{})
	require.NoError(t, err)
	assert.JSONEq(t, `{"s":null,"m":null}`, string(got))
}

// TestMarshalIndent asserts the exact indentation bytes, split into lines so
// the comparison is about layout rather than JSON equality.
func TestMarshalIndent(t *testing.T) {
	got, err := json.API.MarshalIndent(payload{Name: "gin", N: 1}, "", "    ")
	require.NoError(t, err)
	assert.Equal(t, []string{
		"{",
		`    "name": "gin",`,
		`    "n": 1`,
		"}",
	}, strings.Split(string(got), "\n"))
}

// TestEncoder covers what render.PureJSON depends on: a trailing newline, and
// SetEscapeHTML(false) actually disabling escaping.
func TestEncoder(t *testing.T) {
	for _, tt := range []struct {
		name       string
		escapeHTML bool
		want       string
	}{
		{"escaped", true, "{\"name\":\"\\u003ca\\u003e\",\"n\":0}\n"},
		{"unescaped", false, "{\"name\":\"<a>\",\"n\":0}\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := json.API.NewEncoder(&buf)
			enc.SetEscapeHTML(tt.escapeHTML)
			require.NoError(t, enc.Encode(payload{Name: "<a>"}))
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestDecoder(t *testing.T) {
	var got payload
	require.NoError(t, json.API.NewDecoder(strings.NewReader(`{"name":"gin","n":1}`)).Decode(&got))
	assert.Equal(t, payload{Name: "gin", N: 1}, got)
}

// TestDecoderDecodesStream pins that consecutive values can be read from one
// decoder, which binding relies on for request bodies that are not exhausted
// by a single value.
func TestDecoderDecodesStream(t *testing.T) {
	dec := json.API.NewDecoder(strings.NewReader(`{"n":1}{"n":2}`))
	for _, want := range []int{1, 2} {
		var got payload
		require.NoError(t, dec.Decode(&got))
		assert.Equal(t, want, got.N)
	}
}

// TestDecoderUseNumber backs binding.EnableDecoderUseNumber. Asserted through
// the textual form rather than a concrete type so it holds for every codec: a
// float64 would render as 1.5 and lose the trailing zero.
func TestDecoderUseNumber(t *testing.T) {
	dec := json.API.NewDecoder(strings.NewReader(`{"n":1.50}`))
	dec.UseNumber()

	var got map[string]any
	require.NoError(t, dec.Decode(&got))
	assert.Equal(t, "1.50", fmt.Sprint(got["n"]))
}

// TestDecoderDisallowUnknownFields backs binding.EnableDecoderDisallowUnknownFields.
func TestDecoderDisallowUnknownFields(t *testing.T) {
	dec := json.API.NewDecoder(strings.NewReader(`{"name":"gin","nope":1}`))
	dec.DisallowUnknownFields()

	var got payload
	assert.Error(t, dec.Decode(&got))
}
