// Copyright 2013 Julien Schmidt. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package gin

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cleanPathTest struct {
	path, result string
}

var cleanTests = []cleanPathTest{
	// Already clean
	{"/", "/"},
	{"/abc", "/abc"},
	{"/a/b/c", "/a/b/c"},
	{"/abc/", "/abc/"},
	{"/a/b/c/", "/a/b/c/"},

	// missing root
	{"", "/"},
	{"a/", "/a/"},
	{"abc", "/abc"},
	{"abc/def", "/abc/def"},
	{"a/b/c", "/a/b/c"},

	// Remove doubled slash
	{"//", "/"},
	{"/abc//", "/abc/"},
	{"/abc/def//", "/abc/def/"},
	{"/a/b/c//", "/a/b/c/"},
	{"/abc//def//ghi", "/abc/def/ghi"},
	{"//abc", "/abc"},
	{"///abc", "/abc"},
	{"//abc//", "/abc/"},

	// Remove . elements
	{".", "/"},
	{"./", "/"},
	{"/abc/./def", "/abc/def"},
	{"/./abc/def", "/abc/def"},
	{"/abc/.", "/abc/"},

	// Remove .. elements
	{"..", "/"},
	{"../", "/"},
	{"../../", "/"},
	{"../..", "/"},
	{"../../abc", "/abc"},
	{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
	{"/abc/def/../ghi/../jkl", "/abc/jkl"},
	{"/abc/def/..", "/abc"},
	{"/abc/def/../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

	// Combinations
	{"abc/./../def", "/def"},
	{"abc//./../def", "/def"},
	{"abc/../../././../def", "/def"},
}

func TestPathClean(t *testing.T) {
	for _, test := range cleanTests {
		assert.Equal(t, test.result, cleanPath(test.path))
		assert.Equal(t, test.result, cleanPath(test.result))
	}
}

func TestPathCleanMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}

	if runtime.GOMAXPROCS(0) > 1 {
		t.Skip("skipping malloc count; GOMAXPROCS>1")
	}

	for _, test := range cleanTests {
		allocs := testing.AllocsPerRun(100, func() { cleanPath(test.result) })
		assert.InDelta(t, 0, allocs, 0.01)
	}
}

func BenchmarkPathClean(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		for _, test := range cleanTests {
			cleanPath(test.path)
		}
	}
}

func genLongPaths() (testPaths []cleanPathTest) {
	for i := 1; i <= 1234; i++ {
		ss := strings.Repeat("a", i)

		correctPath := "/" + ss
		testPaths = append(testPaths, cleanPathTest{
			path:   correctPath,
			result: correctPath,
		}, cleanPathTest{
			path:   ss,
			result: correctPath,
		}, cleanPathTest{
			path:   "//" + ss,
			result: correctPath,
		}, cleanPathTest{
			path:   "/" + ss + "/b/..",
			result: correctPath,
		})
	}
	return
}

func TestPathCleanLong(t *testing.T) {
	cleanTests := genLongPaths()

	for _, test := range cleanTests {
		assert.Equal(t, test.result, cleanPath(test.path))
		assert.Equal(t, test.result, cleanPath(test.result))
	}
}

func TestPathCleanFastPathHelpers(t *testing.T) {
	cleanSegmentTests := []struct {
		name string
		path string
		want bool
	}{
		{"empty", "", false},
		{"dot", ".", false},
		{"dotdot", "..", false},
		{"nested", "abc/def", false},
		{"clean", "abc", true},
	}
	for _, test := range cleanSegmentTests {
		t.Run("single segment "+test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isSingleCleanPathSegment(test.path))
		})
	}

	repeatedSlashTests := []struct {
		name    string
		path    string
		want    string
		wantHit bool
	}{
		{"only slashes", "///", "/", true},
		{"dot segment", "//.", "/", true},
		{"dotdot segment", "//..", "/", true},
		{"single segment", "//abc", "/abc", true},
		{"nested segment", "//abc/def", "", false},
	}
	for _, test := range repeatedSlashTests {
		t.Run("leading slash "+test.name, func(t *testing.T) {
			got, ok := cleanRepeatedLeadingSlash(test.path)
			assert.Equal(t, test.wantHit, ok)
			assert.Equal(t, test.want, got)
		})
	}

	parentPathTests := []struct {
		name    string
		path    string
		want    string
		wantHit bool
	}{
		{"root parent", "/abc/..", "/", true},
		{"clean parent", "/abc/b/..", "/abc", true},
		{"not trailing parent", "/abc/b", "", false},
		{"unclean prefix", "/abc/./b/..", "", false},
	}
	for _, test := range parentPathTests {
		t.Run("trailing parent "+test.name, func(t *testing.T) {
			got, ok := cleanTrailingParentPath(test.path)
			assert.Equal(t, test.wantHit, ok)
			assert.Equal(t, test.want, got)
		})
	}

	absolutePathTests := []struct {
		name string
		path string
		want bool
	}{
		{"empty", "", false},
		{"relative", "abc", false},
		{"root", "/", true},
		{"single segment", "/abc", true},
		{"nested", "/abc/def", true},
		{"trailing slash", "/abc/", false},
		{"dot segment", "/abc/./def", false},
		{"dotdot segment", "/abc/../def", false},
		{"empty segment", "/abc//def", false},
	}
	for _, test := range absolutePathTests {
		t.Run("absolute path "+test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isCleanAbsolutePath(test.path))
		})
	}
}

func BenchmarkPathCleanLong(b *testing.B) {
	cleanTests := genLongPaths()

	b.ReportAllocs()

	for b.Loop() {
		for _, test := range cleanTests {
			cleanPath(test.path)
		}
	}
}

func TestRemoveRepeatedChar(t *testing.T) {
	testCases := []struct {
		name string
		str  string
		char byte
		want string
	}{
		{
			name: "empty",
			str:  "",
			char: 'a',
			want: "",
		},
		{
			name: "noSlash",
			str:  "abc",
			char: ',',
			want: "abc",
		},
		{
			name: "withSlash",
			str:  "/a/b/c/",
			char: '/',
			want: "/a/b/c/",
		},
		{
			name: "withRepeatedSlashes",
			str:  "/a//b///c////",
			char: '/',
			want: "/a/b/c/",
		},
		{
			name: "threeSlashes",
			str:  "///",
			char: '/',
			want: "/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := removeRepeatedChar(tc.str, tc.char)
			assert.Equal(t, tc.want, res)
		})
	}
}
