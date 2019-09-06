package gin

import "testing"

func oldCleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

var result string

func BenchmarkCleanPath(b *testing.B) {
	functions := []struct {
		test string
		fun  func(string) string
	}{
		{"OldCleanPath", oldCleanPath},
		{"NewCleanPath", cleanPath},
	}

	paths := []string{
		"/", "/abc", "/a/b/c", "/abc/", "/a/b/c/",
		"", "a/", "abc", "abc/def", "a/b/c",
		"//", "/abc//", "/abc/def//", "/a/b/c//",
		"/abc//def//ghi", "//abc", "///abc",
		"//abc//", ".", "./", "/abc/./def",
		"/./abc/def", "/abc/.", "..", "../",
		"../../", "../..", "../../abc",
		"/abc/def/ghi/../jkl", "/abc/def/../ghi/../jkl",
		"/abc/def/..", "/abc/def/../..",
		"/abc/def/../../..", "/abc/def/../../..",
		"/abc/def/../../../ghi/jkl/../../../mno",
		"abc/..def", "abc/./...", "abc/a../...", "abc/a..z",
		"abc/./../def", "abc//./../def", "abc/../../././../def",
		"abc/./../..def", "abc/../.../..def", "abc/.//../..def",
	}

	for _, function := range functions {
		b.Run(function.test, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				var r string
				for _, path := range paths {
					r = function.fun(path)
				}
				result = r
			}
		})
	}
}
