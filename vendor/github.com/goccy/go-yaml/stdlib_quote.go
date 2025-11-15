// Copied and trimmed down from https://github.com/golang/go/blob/e3769299cd3484e018e0e2a6e1b95c2b18ce4f41/src/strconv/quote.go
// We want to use the standard library's private "quoteWith" function rather than write our own so that we get robust unicode support.
// Every private function called by quoteWith was copied.
// There are 2 modifications to simplify the code:
// 1. The unicode.IsPrint function was substituted for the custom implementation of IsPrint
// 2. All code paths reachable only when ASCIIonly or grphicOnly are set to true were removed.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaml

import (
	"unicode"
	"unicode/utf8"
)

const (
	lowerhex = "0123456789abcdef"
)

func quoteWith(s string, quote byte) string {
	return string(appendQuotedWith(make([]byte, 0, 3*len(s)/2), s, quote))
}

func appendQuotedWith(buf []byte, s string, quote byte) []byte {
	// Often called with big strings, so preallocate. If there's quoting,
	// this is conservative but still helps a lot.
	if cap(buf)-len(buf) < len(s) {
		nBuf := make([]byte, len(buf), len(buf)+1+len(s)+1)
		copy(nBuf, buf)
		buf = nBuf
	}
	buf = append(buf, quote)
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[s[0]>>4])
			buf = append(buf, lowerhex[s[0]&0xF])
			continue
		}
		buf = appendEscapedRune(buf, r, quote)
	}
	buf = append(buf, quote)
	return buf
}

func appendEscapedRune(buf []byte, r rune, quote byte) []byte {
	var runeTmp [utf8.UTFMax]byte
	// goccy/go-yaml patch on top of the standard library's appendEscapedRune function.
	//
	// We use this to implement the YAML single-quoted string, where the only escape sequence is '', which represents a single quote.
	// The below snippet from the standard library is for escaping e.g. \ with \\, which is not what we want for the single-quoted string.
	//
	// if r == rune(quote) || r == '\\' { // always backslashed
	// 	buf = append(buf, '\\')
	// 	buf = append(buf, byte(r))
	// 	return buf
	// }
	if r == rune(quote) {
		buf = append(buf, byte(r))
		buf = append(buf, byte(r))
		return buf
	}
	if unicode.IsPrint(r) {
		n := utf8.EncodeRune(runeTmp[:], r)
		buf = append(buf, runeTmp[:n]...)
		return buf
	}
	switch r {
	case '\a':
		buf = append(buf, `\a`...)
	case '\b':
		buf = append(buf, `\b`...)
	case '\f':
		buf = append(buf, `\f`...)
	case '\n':
		buf = append(buf, `\n`...)
	case '\r':
		buf = append(buf, `\r`...)
	case '\t':
		buf = append(buf, `\t`...)
	case '\v':
		buf = append(buf, `\v`...)
	default:
		switch {
		case r < ' ':
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(r)>>4])
			buf = append(buf, lowerhex[byte(r)&0xF])
		case r > utf8.MaxRune:
			r = 0xFFFD
			fallthrough
		case r < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		}
	}
	return buf
}
