package charset

import (
	"bytes"
	"unicode/utf8"

	"github.com/gabriel-vasile/mimetype/internal/markup"
	"github.com/gabriel-vasile/mimetype/internal/scan"
)

const (
	F = 0 /* character never appears in text */
	T = 1 /* character appears in plain ASCII text */
	I = 2 /* character appears in ISO-8859 text */
	X = 3 /* character appears in non-ISO extended ASCII (Mac, IBM PC) */
)

var (
	boms = []struct {
		bom []byte
		enc string
	}{
		{[]byte{0xEF, 0xBB, 0xBF}, "utf-8"},
		{[]byte{0x00, 0x00, 0xFE, 0xFF}, "utf-32be"},
		{[]byte{0xFF, 0xFE, 0x00, 0x00}, "utf-32le"},
		{[]byte{0xFE, 0xFF}, "utf-16be"},
		{[]byte{0xFF, 0xFE}, "utf-16le"},
	}

	// https://github.com/file/file/blob/fa93fb9f7d21935f1c7644c47d2975d31f12b812/src/encoding.c#L241
	textChars = [256]byte{
		/*                  BEL BS HT LF VT FF CR    */
		F, F, F, F, F, F, F, T, T, T, T, T, T, T, F, F, /* 0x0X */
		/*                              ESC          */
		F, F, F, F, F, F, F, F, F, F, F, T, F, F, F, F, /* 0x1X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, /* 0x2X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, /* 0x3X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, /* 0x4X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, /* 0x5X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, /* 0x6X */
		T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, F, /* 0x7X */
		/*            NEL                            */
		X, X, X, X, X, T, X, X, X, X, X, X, X, X, X, X, /* 0x8X */
		X, X, X, X, X, X, X, X, X, X, X, X, X, X, X, X, /* 0x9X */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xaX */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xbX */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xcX */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xdX */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xeX */
		I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, I, /* 0xfX */
	}
)

// FromBOM returns the charset declared in the BOM of content.
func FromBOM(content []byte) string {
	for _, b := range boms {
		if bytes.HasPrefix(content, b.bom) {
			return b.enc
		}
	}
	return ""
}

// FromPlain returns the charset of a plain text. It relies on BOM presence
// and it falls back on checking each byte in content.
func FromPlain(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	if cset := FromBOM(content); cset != "" {
		return cset
	}
	origContent := content
	// Try to detect UTF-8.
	// First eliminate any partial rune at the end.
	for i := len(content) - 1; i >= 0 && i > len(content)-4; i-- {
		b := content[i]
		if b < 0x80 {
			break
		}
		if utf8.RuneStart(b) {
			content = content[:i]
			break
		}
	}
	hasHighBit := false
	for _, c := range content {
		if c >= 0x80 {
			hasHighBit = true
			break
		}
	}
	if hasHighBit && utf8.Valid(content) {
		return "utf-8"
	}

	// ASCII is a subset of UTF8. Follow W3C recommendation and replace with UTF8.
	if ascii(origContent) {
		return "utf-8"
	}

	return latin(origContent)
}

func latin(content []byte) string {
	hasControlBytes := false
	for _, b := range content {
		t := textChars[b]
		if t != T && t != I {
			return ""
		}
		if b >= 0x80 && b <= 0x9F {
			hasControlBytes = true
		}
	}
	// Code range 0x80 to 0x9F is reserved for control characters in ISO-8859-1
	// (so-called C1 Controls). Windows 1252, however, has printable punctuation
	// characters in this range.
	if hasControlBytes {
		return "windows-1252"
	}
	return "iso-8859-1"
}

func ascii(content []byte) bool {
	for _, b := range content {
		if textChars[b] != T {
			return false
		}
	}
	return true
}

// FromXML returns the charset of an XML document. It relies on the XML
// header <?xml version="1.0" encoding="UTF-8"?> and falls back on the plain
// text content.
func FromXML(content []byte) string {
	if cset := fromXML(content); cset != "" {
		return cset
	}
	return FromPlain(content)
}
func fromXML(s scan.Bytes) string {
	xml := []byte("<?XML")
	lxml := len(xml)
	for {
		if len(s) == 0 {
			return ""
		}
		for scan.ByteIsWS(s.Peek()) {
			s.Advance(1)
		}
		if len(s) <= lxml {
			return ""
		}
		if !s.Match(xml, scan.IgnoreCase) {
			s = s[1:] // safe to slice instead of s.Advance(1) because bounds are checked
			continue
		}
		aName, aVal, hasMore := "", "", true
		for hasMore {
			aName, aVal, hasMore = markup.GetAnAttribute(&s)
			if aName == "encoding" && aVal != "" {
				return aVal
			}
		}
	}
}

// FromHTML returns the charset of an HTML document. It first looks if a BOM is
// present and if so uses it to determine the charset. If no BOM is present,
// it relies on the meta tag <meta charset="UTF-8"> and falls back on the
// plain text content.
func FromHTML(content []byte) string {
	if cset := FromBOM(content); cset != "" {
		return cset
	}
	if cset := fromHTML(content); cset != "" {
		return cset
	}
	return FromPlain(content)
}

func fromHTML(s scan.Bytes) string {
	const (
		dontKnow = iota
		doNeedPragma
		doNotNeedPragma
	)
	meta := []byte("<META")
	body := []byte("<BODY")
	lmeta := len(meta)
	for {
		if markup.SkipAComment(&s) {
			continue
		}
		if len(s) <= lmeta {
			return ""
		}
		// Abort when <body is reached.
		if s.Match(body, scan.IgnoreCase) {
			return ""
		}
		if !s.Match(meta, scan.IgnoreCase) {
			s = s[1:] // safe to slice instead of s.Advance(1) because bounds are checked
			continue
		}
		s = s[lmeta:]
		c := s.Pop()
		if c == 0 || (!scan.ByteIsWS(c) && c != '/') {
			return ""
		}
		attrList := make(map[string]bool)
		gotPragma := false
		needPragma := dontKnow

		charset := ""
		aName, aVal, hasMore := "", "", true
		for hasMore {
			aName, aVal, hasMore = markup.GetAnAttribute(&s)
			if attrList[aName] {
				continue
			}
			// processing step
			if len(aName) == 0 && len(aVal) == 0 {
				if needPragma == dontKnow {
					continue
				}
				if needPragma == doNeedPragma && !gotPragma {
					continue
				}
			}
			attrList[aName] = true
			if aName == "http-equiv" && scan.Bytes(aVal).Match([]byte("CONTENT-TYPE"), scan.IgnoreCase) {
				gotPragma = true
			} else if aName == "content" {
				charset = string(extractCharsetFromMeta(scan.Bytes(aVal)))
				if len(charset) != 0 {
					needPragma = doNeedPragma
				}
			} else if aName == "charset" {
				charset = aVal
				needPragma = doNotNeedPragma
			}
		}

		if needPragma == dontKnow || needPragma == doNeedPragma && !gotPragma {
			continue
		}

		return charset
	}
}

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#algorithm-for-extracting-a-character-encoding-from-a-meta-element
func extractCharsetFromMeta(s scan.Bytes) []byte {
	for {
		i := bytes.Index(s, []byte("charset"))
		if i == -1 {
			return nil
		}
		s.Advance(i + len("charset"))
		for scan.ByteIsWS(s.Peek()) {
			s.Advance(1)
		}
		if s.Pop() != '=' {
			continue
		}
		for scan.ByteIsWS(s.Peek()) {
			s.Advance(1)
		}
		quote := s.Peek()
		if quote == 0 {
			return nil
		}
		if quote == '"' || quote == '\'' {
			s.Advance(1)
			return bytes.TrimSpace(s.PopUntil(quote))
		}

		return bytes.TrimSpace(s.PopUntil(';', '\t', '\n', '\x0c', '\r', ' '))
	}
}
