// Package scan has functions for scanning byte slices.
package scan

import (
	"bytes"
	"encoding/binary"
)

// Bytes is a byte slice with helper methods for easier scanning.
type Bytes []byte

func (b *Bytes) Advance(n int) bool {
	if n < 0 || len(*b) < n {
		return false
	}
	*b = (*b)[n:]
	return true
}

// TrimLWS trims whitespace from beginning of the bytes.
func (b *Bytes) TrimLWS() {
	firstNonWS := 0
	for ; firstNonWS < len(*b) && ByteIsWS((*b)[firstNonWS]); firstNonWS++ {
	}

	*b = (*b)[firstNonWS:]
}

// TrimRWS trims whitespace from the end of the bytes.
func (b *Bytes) TrimRWS() {
	lb := len(*b)
	for lb > 0 && ByteIsWS((*b)[lb-1]) {
		*b = (*b)[:lb-1]
		lb--
	}
}

// Peek one byte from b or 0x00 if b is empty.
func (b *Bytes) Peek() byte {
	if len(*b) > 0 {
		return (*b)[0]
	}
	return 0
}

// Pop one byte from b or 0x00 if b is empty.
func (b *Bytes) Pop() byte {
	if len(*b) > 0 {
		ret := (*b)[0]
		*b = (*b)[1:]
		return ret
	}
	return 0
}

// PopN pops n bytes from b or nil if b is empty.
func (b *Bytes) PopN(n int) []byte {
	if len(*b) >= n {
		ret := (*b)[:n]
		*b = (*b)[n:]
		return ret
	}
	return nil
}

// PopUntil will advance b until, but not including, the first occurence of stopAt
// character. If no occurence is found, then it will advance until the end of b.
// The returned Bytes is a slice of all the bytes that we're advanced over.
func (b *Bytes) PopUntil(stopAt ...byte) Bytes {
	if len(*b) == 0 {
		return Bytes{}
	}
	i := bytes.IndexAny(*b, string(stopAt))
	if i == -1 {
		i = len(*b)
	}

	prefix := (*b)[:i]
	*b = (*b)[i:]
	return Bytes(prefix)
}

// ReadSlice is the same as PopUntil, but the returned value includes stopAt as well.
func (b *Bytes) ReadSlice(stopAt byte) Bytes {
	if len(*b) == 0 {
		return Bytes{}
	}
	i := bytes.IndexByte(*b, stopAt)
	if i == -1 {
		i = len(*b)
	} else {
		i++
	}

	prefix := (*b)[:i]
	*b = (*b)[i:]
	return Bytes(prefix)
}

// Line returns the first line from b and advances b with the length of the
// line. One new line character is trimmed after the line if it exists.
func (b *Bytes) Line() Bytes {
	line := b.PopUntil('\n')
	lline := len(line)
	if lline > 0 && line[lline-1] == '\r' {
		line = line[:lline-1]
	}
	b.Advance(1)
	return line
}

// DropLastLine drops the last incomplete line from b.
//
// mimetype limits itself to ReadLimit bytes when performing a detection.
// This means, for file formats like CSV for NDJSON, the last line of the input
// can be an incomplete line.
// If b length is less than readLimit, it means we received an incomplete file
// and proceed with dropping the last line.
func (b *Bytes) DropLastLine(readLimit uint32) {
	if readLimit == 0 || uint32(len(*b)) < readLimit {
		return
	}

	for i := len(*b) - 1; i > 0; i-- {
		if (*b)[i] == '\n' {
			*b = (*b)[:i]
			return
		}
	}
}

func (b *Bytes) Uint16() (uint16, bool) {
	if len(*b) < 2 {
		return 0, false
	}
	v := binary.LittleEndian.Uint16(*b)
	*b = (*b)[2:]
	return v, true
}

const (
	CompactWS = 1 << iota
	IgnoreCase
)

// Search for occurences of pattern p inside b at any index.
func (b Bytes) Search(p []byte, flags int) int {
	if flags == 0 {
		return bytes.Index(b, p)
	}

	lb, lp := len(b), len(p)
	for i := range b {
		if lb-i < lp {
			return -1
		}
		if b[i:].Match(p, flags) {
			return i
		}
	}

	return 0
}

// Match pattern p at index 0 of b.
func (b Bytes) Match(p []byte, flags int) bool {
	for len(b) > 0 {
		// If we finished all we we're looking for from p.
		if len(p) == 0 {
			return true
		}
		if flags&IgnoreCase > 0 && isUpper(p[0]) {
			if upper(b[0]) != p[0] {
				return false
			}
			b, p = b[1:], p[1:]
		} else if flags&CompactWS > 0 && ByteIsWS(p[0]) {
			p = p[1:]
			if !ByteIsWS(b[0]) {
				return false
			}
			b = b[1:]
			if !ByteIsWS(p[0]) {
				b.TrimLWS()
			}
		} else {
			if b[0] != p[0] {
				return false
			}
			b, p = b[1:], p[1:]
		}
	}
	return true
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}
func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}

func ByteIsWS(b byte) bool {
	return b == '\t' || b == '\n' || b == '\x0c' || b == '\r' || b == ' '
}

var (
	ASCIISpaces = []byte{' ', '\r', '\n', '\x0c', '\t'}
	ASCIIDigits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)
