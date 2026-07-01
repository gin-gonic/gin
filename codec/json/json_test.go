// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	if API == nil {
		t.Fatal("API is nil")
	}
	if Package == "" {
		t.Fatal("Package is empty")
	}
}

func TestAPIMarshalAndUnmarshal(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	data, err := API.Marshal(payload{Name: "gin"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"name":"gin"}` {
		t.Fatalf("unexpected marshal output: %s", data)
	}

	var decoded payload
	if err := API.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Name != "gin" {
		t.Fatalf("unexpected decoded payload: %#v", decoded)
	}
}

func TestAPIMarshalIndent(t *testing.T) {
	data, err := API.MarshalIndent(map[string]string{"name": "gin"}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("\n  ")) {
		t.Fatalf("expected indented JSON, got %q", data)
	}
}

func TestAPIEncoder(t *testing.T) {
	var buf bytes.Buffer
	encoder := API.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode("<gin>"); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "\"<gin>\"\n" {
		t.Fatalf("unexpected encoded JSON: %q", got)
	}
}

func TestAPIDecoder(t *testing.T) {
	decoder := API.NewDecoder(strings.NewReader(`{"known": 1, "extra": 2}`))
	decoder.UseNumber()
	decoder.DisallowUnknownFields()

	var dst struct {
		Known json.Number `json:"known"`
	}
	if err := decoder.Decode(&dst); err == nil {
		t.Fatal("expected unknown field error")
	}
}
