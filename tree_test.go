// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package gin

import (
	"reflect"
	"strings"
	"testing"
)

// Used as a workaround since we can't compare functions or their addresses
var fakeHandlerValue string

func fakeHandler(val string) HandlersChain {
	return HandlersChain{func(c *Context) {
		fakeHandlerValue = val
	}}
}

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
	ps         Params
}

func getParams() *Params {
	ps := make(Params, 0, 20)
	return &ps
}

func getSkippedNodes() *[]skippedNode {
	ps := make([]skippedNode, 0, 20)
	return &ps
}

func checkRequests(t *testing.T, tree *node, requests testRequests, unescapes ...bool) {
	unescape := false
	if len(unescapes) >= 1 {
		unescape = unescapes[0]
	}

	for _, request := range requests {
		value := tree.getValue(request.path, getParams(), getSkippedNodes(), unescape)

		if value.handlers == nil {
			if !request.nilHandler {
				t.Errorf("handle mismatch for route '%s': Expected non-nil handle", request.path)
			}
		} else if request.nilHandler {
			t.Errorf("handle mismatch for route '%s': Expected nil handle", request.path)
		} else {
			value.handlers[0](nil)
			if fakeHandlerValue != request.route {
				t.Errorf("handle mismatch for route '%s': Wrong handle (%s != %s)", request.path, fakeHandlerValue, request.route)
			}
		}

		if value.params != nil {
			if !reflect.DeepEqual(*value.params, request.ps) {
				t.Errorf("Params mismatch for route '%s'", request.path)
			}
		}

	}
}

func checkPriorities(t *testing.T, n *node) uint32 {
	var prio uint32
	for i := range n.children {
		prio += checkPriorities(t, n.children[i])
	}

	if n.handlers != nil {
		prio++
	}

	if n.priority != prio {
		t.Errorf(
			"priority mismatch for node '%s': is %d, should be %d",
			n.path, n.priority, prio,
		)
	}

	return prio
}

func TestCountParams(t *testing.T) {
	if countParams("/path/:param1/static/*catch-all") != 2 {
		t.Fail()
	}
	if countParams(strings.Repeat("/:param", 256)) != 256 {
		t.Fail()
	}
}

func TestTreeAddAndGet(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
	}
	for _, route := range routes {
		tree.addRoute(route, fakeHandler(route))
	}

	checkRequests(t, tree, testRequests{
		{"/a", false, "/a", nil},
		{"/", true, "", nil},
		{"/hi", false, "/hi", nil},
		{"/contact", false, "/contact", nil},
		{"/co", false, "/co", nil},
		{"/con", true, "", nil},  // key mismatch
		{"/cona", true, "", nil}, // key mismatch
		{"/no", true, "", nil},   // no matching child
		{"/ab", false, "/ab", nil},
		{"/α", false, "/α", nil},
		{"/β", false, "/β", nil},
	})

	checkPriorities(t, tree)
}

func TestTreeWildcard(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/",
		"/cmd/:tool/",
		"/cmd/:tool/:sub",
		"/cmd/whoami",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/src/*filepath",
		"/search/",
		"/search/:query",
		"/search/gin-gonic",
		"/search/google",
		"/user_:name",
		"/user_:name/about",
		"/files/:dir/*filepath",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
		"/info/:user/project/golang",
		"/aa/*xx",
		"/ab/*xx",
		"/:cc",
		"/c1/:dd/e",
		"/c1/:dd/e1",
		"/:cc/cc",
		"/:cc/:dd/ee",
		"/:cc/:dd/:ee/ff",
		"/:cc/:dd/:ee/:ff/gg",
		"/:cc/:dd/:ee/:ff/:gg/hh",
		"/get/test/abc/",
		"/get/:param/abc/",
		"/something/:paramname/thirdthing",
		"/something/secondthing/test",
		"/get/abc",
		"/get/:param",
		"/get/abc/123abc",
		"/get/abc/:param",
		"/get/abc/123abc/xxx8",
		"/get/abc/123abc/:param",
		"/get/abc/123abc/xxx8/1234",
		"/get/abc/123abc/xxx8/:param",
		"/get/abc/123abc/xxx8/1234/ffas",
		"/get/abc/123abc/xxx8/1234/:param",
		"/get/abc/123abc/xxx8/1234/kkdd/12c",
		"/get/abc/123abc/xxx8/1234/kkdd/:param",
		"/get/abc/:param/test",
		"/get/abc/123abd/:param",
		"/get/abc/123abddd/:param",
		"/get/abc/123/:param",
		"/get/abc/123abg/:param",
		"/get/abc/123abf/:param",
		"/get/abc/123abfff/:param",
		"/get/abc/escaped_colon/test\\:param",
	}
	for _, route := range routes {
		tree.addRoute(route, fakeHandler(route))
	}

	checkRequests(t, tree, testRequests{
		{"/", false, "/", nil},
		{"/cmd/test", true, "/cmd/:tool/", Params{Param{"tool", "test"}}},
		{"/cmd/test/", false, "/cmd/:tool/", Params{Param{"tool", "test"}}},
		{"/cmd/test/3", false, "/cmd/:tool/:sub", Params{Param{Key: "tool", Value: "test"}, Param{Key: "sub", Value: "3"}}},
		{"/cmd/who", true, "/cmd/:tool/", Params{Param{"tool", "who"}}},
		{"/cmd/who/", false, "/cmd/:tool/", Params{Param{"tool", "who"}}},
		{"/cmd/whoami", false, "/cmd/whoami", nil},
		{"/cmd/whoami/", true, "/cmd/whoami", nil},
		{"/cmd/whoami/r", false, "/cmd/:tool/:sub", Params{Param{Key: "tool", Value: "whoami"}, Param{Key: "sub", Value: "r"}}},
		{"/cmd/whoami/r/", true, "/cmd/:tool/:sub", Params{Param{Key: "tool", Value: "whoami"}, Param{Key: "sub", Value: "r"}}},
		{"/cmd/whoami/root", false, "/cmd/whoami/root", nil},
		{"/cmd/whoami/root/", false, "/cmd/whoami/root/", nil},
		{"/src/", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/"}}},
		{"/src/some/file.png", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/some/file.png"}}},
		{"/search/", false, "/search/", nil},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query", Params{Param{Key: "query", Value: "someth!ng+in+ünìcodé"}}},
		{"/search/someth!ng+in+ünìcodé/", true, "", Params{Param{Key: "query", Value: "someth!ng+in+ünìcodé"}}},
		{"/search/gin", false, "/search/:query", Params{Param{"query", "gin"}}},
		{"/search/gin-gonic", false, "/search/gin-gonic", nil},
		{"/search/google", false, "/search/google", nil},
		{"/user_gopher", false, "/user_:name", Params{Param{Key: "name", Value: "gopher"}}},
		{"/user_gopher/about", false, "/user_:name/about", Params{Param{Key: "name", Value: "gopher"}}},
		{"/files/js/inc/framework.js", false, "/files/:dir/*filepath", Params{Param{Key: "dir", Value: "js"}, Param{Key: "filepath", Value: "/inc/framework.js"}}},
		{"/info/gordon/public", false, "/info/:user/public", Params{Param{Key: "user", Value: "gordon"}}},
		{"/info/gordon/project/go", false, "/info/:user/project/:project", Params{Param{Key: "user", Value: "gordon"}, Param{Key: "project", Value: "go"}}},
		{"/info/gordon/project/golang", false, "/info/:user/project/golang", Params{Param{Key: "user", Value: "gordon"}}},
		{"/aa/aa", false, "/aa/*xx", Params{Param{Key: "xx", Value: "/aa"}}},
		{"/ab/ab", false, "/ab/*xx", Params{Param{Key: "xx", Value: "/ab"}}},
		{"/a", false, "/:cc", Params{Param{Key: "cc", Value: "a"}}},
		// * Error with argument being intercepted
		// new PR handle (/all /all/cc /a/cc)
		// fix PR: https://github.com/gin-gonic/gin/pull/2796
		{"/all", false, "/:cc", Params{Param{Key: "cc", Value: "all"}}},
		{"/d", false, "/:cc", Params{Param{Key: "cc", Value: "d"}}},
		{"/ad", false, "/:cc", Params{Param{Key: "cc", Value: "ad"}}},
		{"/dd", false, "/:cc", Params{Param{Key: "cc", Value: "dd"}}},
		{"/dddaa", false, "/:cc", Params{Param{Key: "cc", Value: "dddaa"}}},
		{"/aa", false, "/:cc", Params{Param{Key: "cc", Value: "aa"}}},
		{"/aaa", false, "/:cc", Params{Param{Key: "cc", Value: "aaa"}}},
		{"/aaa/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "aaa"}}},
		{"/ab", false, "/:cc", Params{Param{Key: "cc", Value: "ab"}}},
		{"/abb", false, "/:cc", Params{Param{Key: "cc", Value: "abb"}}},
		{"/abb/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "abb"}}},
		{"/allxxxx", false, "/:cc", Params{Param{Key: "cc", Value: "allxxxx"}}},
		{"/alldd", false, "/:cc", Params{Param{Key: "cc", Value: "alldd"}}},
		{"/all/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "all"}}},
		{"/a/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "a"}}},
		{"/c1/d/e", false, "/c1/:dd/e", Params{Param{Key: "dd", Value: "d"}}},
		{"/c1/d/e1", false, "/c1/:dd/e1", Params{Param{Key: "dd", Value: "d"}}},
		{"/c1/d/ee", false, "/:cc/:dd/ee", Params{Param{Key: "cc", Value: "c1"}, Param{Key: "dd", Value: "d"}}},
		{"/cc/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "cc"}}},
		{"/ccc/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "ccc"}}},
		{"/deedwjfs/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "deedwjfs"}}},
		{"/acllcc/cc", false, "/:cc/cc", Params{Param{Key: "cc", Value: "acllcc"}}},
		{"/get/test/abc/", false, "/get/test/abc/", nil},
		{"/get/te/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "te"}}},
		{"/get/testaa/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "testaa"}}},
		{"/get/xx/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "xx"}}},
		{"/get/tt/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "tt"}}},
		{"/get/a/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "a"}}},
		{"/get/t/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "t"}}},
		{"/get/aa/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "aa"}}},
		{"/get/abas/abc/", false, "/get/:param/abc/", Params{Param{Key: "param", Value: "abas"}}},
		{"/something/secondthing/test", false, "/something/secondthing/test", nil},
		{"/something/abcdad/thirdthing", false, "/something/:paramname/thirdthing", Params{Param{Key: "paramname", Value: "abcdad"}}},
		{"/something/secondthingaaaa/thirdthing", false, "/something/:paramname/thirdthing", Params{Param{Key: "paramname", Value: "secondthingaaaa"}}},
		{"/something/se/thirdthing", false, "/something/:paramname/thirdthing", Params{Param{Key: "paramname", Value: "se"}}},
		{"/something/s/thirdthing", false, "/something/:paramname/thirdthing", Params{Param{Key: "paramname", Value: "s"}}},
		{"/c/d/ee", false, "/:cc/:dd/ee", Params{Param{Key: "cc", Value: "c"}, Param{Key: "dd", Value: "d"}}},
		{"/c/d/e/ff", false, "/:cc/:dd/:ee/ff", Params{Param{Key: "cc", Value: "c"}, Param{Key: "dd", Value: "d"}, Param{Key: "ee", Value: "e"}}},
		{"/c/d/e/f/gg", false, "/:cc/:dd/:ee/:ff/gg", Params{Param{Key: "cc", Value: "c"}, Param{Key: "dd", Value: "d"}, Param{Key: "ee", Value: "e"}, Param{Key: "ff", Value: "f"}}},
		{"/c/d/e/f/g/hh", false, "/:cc/:dd/:ee/:ff/:gg/hh", Params{Param{Key: "cc", Value: "c"}, Param{Key: "dd", Value: "d"}, Param{Key: "ee", Value: "e"}, Param{Key: "ff", Value: "f"}, Param{Key: "gg", Value: "g"}}},
		{"/cc/dd/ee/ff/gg/hh", false, "/:cc/:dd/:ee/:ff/:gg/hh", Params{Param{Key: "cc", Value: "cc"}, Param{Key: "dd", Value: "dd"}, Param{Key: "ee", Value: "ee"}, Param{Key: "ff", Value: "ff"}, Param{Key: "gg", Value: "gg"}}},
		{"/get/abc", false, "/get/abc", nil},
		{"/get/a", false, "/get/:param", Params{Param{Key: "param", Value: "a"}}},
		{"/get/abz", false, "/get/:param", Params{Param{Key: "param", Value: "abz"}}},
		{"/get/12a", false, "/get/:param", Params{Param{Key: "param", Value: "12a"}}},
		{"/get/abcd", false, "/get/:param", Params{Param{Key: "param", Value: "abcd"}}},
		{"/get/abc/123abc", false, "/get/abc/123abc", nil},
		{"/get/abc/12", false, "/get/abc/:param", Params{Param{Key: "param", Value: "12"}}},
		{"/get/abc/123ab", false, "/get/abc/:param", Params{Param{Key: "param", Value: "123ab"}}},
		{"/get/abc/xyz", false, "/get/abc/:param", Params{Param{Key: "param", Value: "xyz"}}},
		{"/get/abc/123abcddxx", false, "/get/abc/:param", Params{Param{Key: "param", Value: "123abcddxx"}}},
		{"/get/abc/123abc/xxx8", false, "/get/abc/123abc/xxx8", nil},
		{"/get/abc/123abc/x", false, "/get/abc/123abc/:param", Params{Param{Key: "param", Value: "x"}}},
		{"/get/abc/123abc/xxx", false, "/get/abc/123abc/:param", Params{Param{Key: "param", Value: "xxx"}}},
		{"/get/abc/123abc/abc", false, "/get/abc/123abc/:param", Params{Param{Key: "param", Value: "abc"}}},
		{"/get/abc/123abc/xxx8xxas", false, "/get/abc/123abc/:param", Params{Param{Key: "param", Value: "xxx8xxas"}}},
		{"/get/abc/123abc/xxx8/1234", false, "/get/abc/123abc/xxx8/1234", nil},
		{"/get/abc/123abc/xxx8/1", false, "/get/abc/123abc/xxx8/:param", Params{Param{Key: "param", Value: "1"}}},
		{"/get/abc/123abc/xxx8/123", false, "/get/abc/123abc/xxx8/:param", Params{Param{Key: "param", Value: "123"}}},
		{"/get/abc/123abc/xxx8/78k", false, "/get/abc/123abc/xxx8/:param", Params{Param{Key: "param", Value: "78k"}}},
		{"/get/abc/123abc/xxx8/1234xxxd", false, "/get/abc/123abc/xxx8/:param", Params{Param{Key: "param", Value: "1234xxxd"}}},
		{"/get/abc/123abc/xxx8/1234/ffas", false, "/get/abc/123abc/xxx8/1234/ffas", nil},
		{"/get/abc/123abc/xxx8/1234/f", false, "/get/abc/123abc/xxx8/1234/:param", Params{Param{Key: "param", Value: "f"}}},
		{"/get/abc/123abc/xxx8/1234/ffa", false, "/get/abc/123abc/xxx8/1234/:param", Params{Param{Key: "param", Value: "ffa"}}},
		{"/get/abc/123abc/xxx8/1234/kka", false, "/get/abc/123abc/xxx8/1234/:param", Params{Param{Key: "param", Value: "kka"}}},
		{"/get/abc/123abc/xxx8/1234/ffas321", false, "/get/abc/123abc/xxx8/1234/:param", Params{Param{Key: "param", Value: "ffas321"}}},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c", false, "/get/abc/123abc/xxx8/1234/kkdd/12c", nil},
		{"/get/abc/123abc/xxx8/1234/kkdd/1", false, "/get/abc/123abc/xxx8/1234/kkdd/:param", Params{Param{Key: "param", Value: "1"}}},
		{"/get/abc/123abc/xxx8/1234/kkdd/12", false, "/get/abc/123abc/xxx8/1234/kkdd/:param", Params{Param{Key: "param", Value: "12"}}},
		{"/get/abc/123abc/xxx8/1234/kkdd/12b", false, "/get/abc/123abc/xxx8/1234/kkdd/:param", Params{Param{Key: "param", Value: "12b"}}},
		{"/get/abc/123abc/xxx8/1234/kkdd/34", false, "/get/abc/123abc/xxx8/1234/kkdd/:param", Params{Param{Key: "param", Value: "34"}}},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c2e3", false, "/get/abc/123abc/xxx8/1234/kkdd/:param", Params{Param{Key: "param", Value: "12c2e3"}}},
		{"/get/abc/12/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "12"}}},
		{"/get/abc/123abdd/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123abdd"}}},
		{"/get/abc/123abdddf/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123abdddf"}}},
		{"/get/abc/123ab/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123ab"}}},
		{"/get/abc/123abgg/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123abgg"}}},
		{"/get/abc/123abff/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123abff"}}},
		{"/get/abc/123abffff/test", false, "/get/abc/:param/test", Params{Param{Key: "param", Value: "123abffff"}}},
		{"/get/abc/123abd/test", false, "/get/abc/123abd/:param", Params{Param{Key: "param", Value: "test"}}},
		{"/get/abc/123abddd/test", false, "/get/abc/123abddd/:param", Params{Param{Key: "param", Value: "test"}}},
		{"/get/abc/123/test22", false, "/get/abc/123/:param", Params{Param{Key: "param", Value: "test22"}}},
		{"/get/abc/123abg/test", false, "/get/abc/123abg/:param", Params{Param{Key: "param", Value: "test"}}},
		{"/get/abc/123abf/testss", false, "/get/abc/123abf/:param", Params{Param{Key: "param", Value: "testss"}}},
		{"/get/abc/123abfff/te", false, "/get/abc/123abfff/:param", Params{Param{Key: "param", Value: "te"}}},
		{"/get/abc/escaped_colon/test\\:param", false, "/get/abc/escaped_colon/test\\:param", nil},
	})

	checkPriorities(t, tree)
}

func TestUnescapeParameters(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/",
		"/cmd/:tool/:sub",
		"/cmd/:tool/",
		"/src/*filepath",
		"/search/:query",
		"/files/:dir/*filepath",
		"/info/:user/project/:project",
		"/info/:user",
	}
	for _, route := range routes {
		tree.addRoute(route, fakeHandler(route))
	}

	unescape := true
	checkRequests(t, tree, testRequests{
		{"/", false, "/", nil},
		{"/cmd/test/", false, "/cmd/:tool/", Params{Param{Key: "tool", Value: "test"}}},
		{"/cmd/test", true, "", Params{Param{Key: "tool", Value: "test"}}},
		{"/src/some/file.png", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/some/file.png"}}},
		{"/src/some/file+test.png", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/some/file test.png"}}},
		{"/src/some/file++++%%%%test.png", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/some/file++++%%%%test.png"}}},
		{"/src/some/file%2Ftest.png", false, "/src/*filepath", Params{Param{Key: "filepath", Value: "/some/file/test.png"}}},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query", Params{Param{Key: "query", Value: "someth!ng in ünìcodé"}}},
		{"/info/gordon/project/go", false, "/info/:user/project/:project", Params{Param{Key: "user", Value: "gordon"}, Param{Key: "project", Value: "go"}}},
		{"/info/slash%2Fgordon", false, "/info/:user", Params{Param{Key: "user", Value: "slash/gordon"}}},
		{"/info/slash%2Fgordon/project/Project%20%231", false, "/info/:user/project/:project", Params{Param{Key: "user", Value: "slash/gordon"}, Param{Key: "project", Value: "Project #1"}}},
		{"/info/slash%%%%", false, "/info/:user", Params{Param{Key: "user", Value: "slash%%%%"}}},
		{"/info/slash%%%%2Fgordon/project/Project%%%%20%231", false, "/info/:user/project/:project", Params{Param{Key: "user", Value: "slash%%%%2Fgordon"}, Param{Key: "project", Value: "Project%%%%20%231"}}},
	}, unescape)

	checkPriorities(t, tree)
}

// TestSecureParameterHandling tests the fixes for path traversal vulnerabilities:
// PRISMA-2022-0393: Path traversal due to multiple encodings in path parameters
// PRISMA-2022-0394: Path traversal due to wildcard parameters incorrectly decoding URIs recursively
func TestSecureParameterHandling(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/info/:user",
		"/files/*filepath",
	}
	for _, route := range routes {
		tree.addRoute(route, fakeHandler(route))
	}

	// Test cases for PRISMA-2022-0393 (path parameters)
	// These test that double-encoded parameters are not double-decoded
	unescape := true
	checkRequests(t, tree, testRequests{
		// Normal case - single encoding works as expected
		{"/info/user%2Fprofile", false, "/info/:user", Params{Param{Key: "user", Value: "user/profile"}}},

		// Double encoding - should only decode once
		{"/info/user%252Fprofile", false, "/info/:user", Params{Param{Key: "user", Value: "user%2Fprofile"}}},

		// Triple encoding - should only decode once
		{"/info/user%25252Fprofile", false, "/info/:user", Params{Param{Key: "user", Value: "user%252Fprofile"}}},

		// Mixed encoding - should only decode once
		{"/info/%2Fuser%252Fprofile", false, "/info/:user", Params{Param{Key: "user", Value: "/user%2Fprofile"}}},

		// No encoding - should pass through unchanged
		{"/info/user", false, "/info/:user", Params{Param{Key: "user", Value: "user"}}},
	}, unescape)

	// Test cases for PRISMA-2022-0394 (wildcard parameters)
	// These test that double-encoded wildcards are not double-decoded
	checkRequests(t, tree, testRequests{
		// Normal case - single encoding works as expected
		{"/files/path%2Fto%2Ffile.txt", false, "/files/*filepath", Params{Param{Key: "filepath", Value: "/path/to/file.txt"}}},

		// Double encoding - should only decode once
		{"/files/path%252Fto%252Ffile.txt", false, "/files/*filepath", Params{Param{Key: "filepath", Value: "/path%2Fto%2Ffile.txt"}}},

		// Triple encoding - should only decode once
		{"/files/path%25252Fto%25252Ffile.txt", false, "/files/*filepath", Params{Param{Key: "filepath", Value: "/path%252Fto%252Ffile.txt"}}},

		// Mixed encoding - should only decode once
		{"/files/%2Fpath%252Fto%2Ffile.txt", false, "/files/*filepath", Params{Param{Key: "filepath", Value: "//path%2Fto/file.txt"}}},

		// No encoding - should pass through unchanged
		{"/files/normal/file.txt", false, "/files/*filepath", Params{Param{Key: "filepath", Value: "/normal/file.txt"}}},
	}, unescape)
}

func catchPanic(testFunc func()) (recv any) {
	defer func() {
		recv = recover()
	}()

	testFunc()
	return
}

type testRoute struct {
	path     string
	conflict bool
}

func testRoutes(t *testing.T, routes []testRoute) {
	tree := &node{}

	for _, route := range routes {
		recv := catchPanic(func() {
			tree.addRoute(route.path, nil)
		})

		if route.conflict {
			if recv == nil {
				t.Errorf("no panic for conflicting route '%s'", route.path)
			}
		} else if recv != nil {
			t.Errorf("unexpected panic for route '%s': %v", route.path, recv)
		}
	}
}

func TestTreeWildcardConflict(t *testing.T) {
	routes := []testRoute{
		{"/cmd/:tool/:sub", false},
		{"/cmd/vet", false},
		{"/foo/bar", false},
		{"/foo/:name", false},
		{"/foo/:names", true},
		{"/cmd/*path", true},
		{"/cmd/:badvar", true},
		{"/cmd/:tool/names", false},
		{"/cmd/:tool/:badsub/details", true},
		{"/src/*filepath", false},
		{"/src/:file", true},
		{"/src/static.json", true},
		{"/src/*filepathx", true},
		{"/src/", true},
		{"/src/foo/bar", true},
		{"/src1/", false},
		{"/src1/*filepath", true},
		{"/src2*filepath", true},
		{"/src2/*filepath", false},
		{"/search/:query", false},
		{"/search/valid", false},
		{"/user_:name", false},
		{"/user_x", false},
		{"/user_:name", false},
		{"/id:id", false},
		{"/id/:id", false},
		{"/static/*file", false},
		{"/static/", true},
		{"/escape/test\\:d1", false},
		{"/escape/test\\:d2", false},
		{"/escape/test:param", false},
	}
	testRoutes(t, routes)
}

func TestCatchAllAfterSlash(t *testing.T) {
	routes := []testRoute{
		{"/non-leading-*catchall", true},
	}
	testRoutes(t, routes)
}

func TestTreeChildConflict(t *testing.T) {
	routes := []testRoute{
		{"/cmd/vet", false},
		{"/cmd/:tool", false},
		{"/cmd/:tool/:sub", false},
		{"/cmd/:tool/misc", false},
		{"/cmd/:tool/:othersub", true},
	}
	testRoutes(t, routes)
}
