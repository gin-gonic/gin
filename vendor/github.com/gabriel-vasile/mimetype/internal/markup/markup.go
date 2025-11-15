// Package markup implements functions for extracting info from
// HTML and XML documents.
package markup

import (
	"bytes"

	"github.com/gabriel-vasile/mimetype/internal/scan"
)

func GetAnAttribute(s *scan.Bytes) (name, val string, hasMore bool) {
	for scan.ByteIsWS(s.Peek()) || s.Peek() == '/' {
		s.Advance(1)
	}
	if s.Peek() == '>' {
		return "", "", false
	}
	// Allocate 10 to avoid resizes.
	// Attribute names and values are continuous slices of bytes in input,
	// so we could do without allocating and returning slices of input.
	nameB := make([]byte, 0, 10)
	// step 4 and 5
	for {
		// bap means byte at position in the specification.
		bap := s.Pop()
		if bap == 0 {
			return "", "", false
		}
		if bap == '=' && len(nameB) > 0 {
			val, hasMore := getAValue(s)
			return string(nameB), string(val), hasMore
		} else if scan.ByteIsWS(bap) {
			for scan.ByteIsWS(s.Peek()) {
				s.Advance(1)
			}
			if s.Peek() != '=' {
				return string(nameB), "", true
			}
			s.Advance(1)
			for scan.ByteIsWS(s.Peek()) {
				s.Advance(1)
			}
			val, hasMore := getAValue(s)
			return string(nameB), string(val), hasMore
		} else if bap == '/' || bap == '>' {
			return string(nameB), "", false
		} else if bap >= 'A' && bap <= 'Z' {
			nameB = append(nameB, bap+0x20)
		} else {
			nameB = append(nameB, bap)
		}
	}
}

func getAValue(s *scan.Bytes) (_ []byte, hasMore bool) {
	for scan.ByteIsWS(s.Peek()) {
		s.Advance(1)
	}
	origS, end := *s, 0
	bap := s.Pop()
	if bap == 0 {
		return nil, false
	}
	end++
	// Step 10
	switch bap {
	case '"', '\'':
		val := s.PopUntil(bap)
		if s.Pop() != bap {
			return nil, false
		}
		return val, s.Peek() != 0 && s.Peek() != '>'
	case '>':
		return nil, false
	}

	// Step 11
	for {
		bap = s.Pop()
		if bap == 0 {
			return nil, false
		}
		switch {
		case scan.ByteIsWS(bap):
			return origS[:end], true
		case bap == '>':
			return origS[:end], false
		default:
			end++
		}
	}
}

func SkipAComment(s *scan.Bytes) (skipped bool) {
	if bytes.HasPrefix(*s, []byte("<!--")) {
		// Offset by 2 len(<!) because the starting and ending -- can be the same.
		if i := bytes.Index((*s)[2:], []byte("-->")); i != -1 {
			s.Advance(i + 2 + 3) // 2 comes from len(<!) and 3 comes from len(-->).
			return true
		}
	}
	return false
}
