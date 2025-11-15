package magic

import (
	"bytes"
	"strconv"

	"github.com/gabriel-vasile/mimetype/internal/scan"
)

// NetPBM matches a Netpbm Portable BitMap ASCII/Binary file.
//
// See: https://en.wikipedia.org/wiki/Netpbm
func NetPBM(raw []byte, _ uint32) bool {
	return netp(raw, "P1\n", "P4\n")
}

// NetPGM matches a Netpbm Portable GrayMap ASCII/Binary file.
//
// See: https://en.wikipedia.org/wiki/Netpbm
func NetPGM(raw []byte, _ uint32) bool {
	return netp(raw, "P2\n", "P5\n")
}

// NetPPM matches a Netpbm Portable PixMap ASCII/Binary file.
//
// See: https://en.wikipedia.org/wiki/Netpbm
func NetPPM(raw []byte, _ uint32) bool {
	return netp(raw, "P3\n", "P6\n")
}

// NetPAM matches a Netpbm Portable Arbitrary Map file.
//
// See: https://en.wikipedia.org/wiki/Netpbm
func NetPAM(raw []byte, _ uint32) bool {
	if !bytes.HasPrefix(raw, []byte("P7\n")) {
		return false
	}
	w, h, d, m, e := false, false, false, false, false
	s := scan.Bytes(raw)
	var l scan.Bytes
	// Read line by line.
	for i := 0; i < 128; i++ {
		l = s.Line()
		// If the line is empty or a comment, skip.
		if len(l) == 0 || l.Peek() == '#' {
			if len(s) == 0 {
				return false
			}
			continue
		} else if bytes.HasPrefix(l, []byte("TUPLTYPE")) {
			continue
		} else if bytes.HasPrefix(l, []byte("WIDTH ")) {
			w = true
		} else if bytes.HasPrefix(l, []byte("HEIGHT ")) {
			h = true
		} else if bytes.HasPrefix(l, []byte("DEPTH ")) {
			d = true
		} else if bytes.HasPrefix(l, []byte("MAXVAL ")) {
			m = true
		} else if bytes.HasPrefix(l, []byte("ENDHDR")) {
			e = true
		}
		// When we reached header, return true if we collected all four required headers.
		// WIDTH, HEIGHT, DEPTH and MAXVAL.
		if e {
			return w && h && d && m
		}
	}
	return false
}

func netp(s scan.Bytes, prefixes ...string) bool {
	foundPrefix := ""
	for _, p := range prefixes {
		if bytes.HasPrefix(s, []byte(p)) {
			foundPrefix = p
		}
	}
	if foundPrefix == "" {
		return false
	}
	s.Advance(len(foundPrefix)) // jump over P1, P2, P3, etc.

	var l scan.Bytes
	// Read line by line.
	for i := 0; i < 128; i++ {
		l = s.Line()
		// If the line is a comment, skip.
		if l.Peek() == '#' {
			continue
		}
		// If line has leading whitespace, then skip over whitespace.
		for scan.ByteIsWS(l.Peek()) {
			l.Advance(1)
		}
		if len(s) == 0 || len(l) > 0 {
			break
		}
	}

	// At this point l should be the two integers denoting the size of the matrix.
	width := l.PopUntil(scan.ASCIISpaces...)
	for scan.ByteIsWS(l.Peek()) {
		l.Advance(1)
	}
	height := l.PopUntil(scan.ASCIISpaces...)

	w, errw := strconv.ParseInt(string(width), 10, 64)
	h, errh := strconv.ParseInt(string(height), 10, 64)
	return errw == nil && errh == nil && w > 0 && h > 0
}
