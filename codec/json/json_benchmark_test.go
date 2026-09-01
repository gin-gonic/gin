// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json_test

import (
	"io"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/codec/json"
)

// These benchmarks compare codecs, and on go1.27 also compare the engine behind
// encoding/json: the jsonv2 GOEXPERIMENT is on by default there, so v1 is itself
// implemented over encoding/json/v2. To see a codec or engine difference, run
// the same benchmark under each configuration and compare with benchstat, e.g.
//
//	go test -bench . -count 10 ./codec/json/ > new.txt
//	GOEXPERIMENT=nojsonv2 go test -bench . -count 10 ./codec/json/ > old.txt
//	benchstat old.txt new.txt
//
// Marshal and Unmarshal do not necessarily move in the same direction, so read
// them separately rather than as one number.

type benchPayload struct {
	ID     int               `json:"id"`
	Name   string            `json:"name"`
	Tags   []string          `json:"tags"`
	Meta   map[string]string `json:"meta"`
	Nested []benchPayload    `json:"nested,omitempty"`
}

var nested = benchPayload{
	ID:   1,
	Name: "gin <framework>",
	Tags: []string{"a", "b", "c"},
	Meta: map[string]string{"k1": "v1", "k2": "v2"},
	Nested: []benchPayload{
		{ID: 2, Name: "x", Tags: []string{"z"}},
		{ID: 3, Name: "y", Meta: map[string]string{"q": "r"}},
	},
}

// flat carries no map, so it isolates struct and slice encoding from the key
// sorting that marshaling a map requires.
type benchFlat struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	N    float64  `json:"n"`
	OK   bool     `json:"ok"`
}

var flat = []benchFlat{
	{1, "alpha <a>", []string{"x", "y"}, 1.5, true},
	{2, "beta", []string{"z"}, 2.25, false},
	{3, "gamma", nil, 3.75, true},
}

// Safe at package scope: this is an external test package, so codec/json is
// fully initialised and API is installed before these run.
var nestedJSON, _ = json.API.Marshal(nested)

func BenchmarkMarshal(b *testing.B) {
	for b.Loop() {
		if _, err := json.API.Marshal(nested); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalFlat(b *testing.B) {
	for b.Loop() {
		if _, err := json.API.Marshal(flat); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalIndent(b *testing.B) {
	for b.Loop() {
		if _, err := json.API.MarshalIndent(nested, "", "    "); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for b.Loop() {
		var out benchPayload
		if err := json.API.Unmarshal(nestedJSON, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncoder(b *testing.B) {
	for b.Loop() {
		if err := json.API.NewEncoder(io.Discard).Encode(nested); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoder(b *testing.B) {
	r := strings.NewReader("")
	for b.Loop() {
		r.Reset(string(nestedJSON))

		var out benchPayload
		if err := json.API.NewDecoder(r).Decode(&out); err != nil {
			b.Fatal(err)
		}
	}
}
