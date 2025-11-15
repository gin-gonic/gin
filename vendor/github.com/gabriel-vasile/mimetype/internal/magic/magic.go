// Package magic holds the matching functions used to find MIME types.
package magic

import (
	"bytes"
	"fmt"

	"github.com/gabriel-vasile/mimetype/internal/scan"
)

type (
	// Detector receiveѕ the raw data of a file and returns whether the data
	// meets any conditions. The limit parameter is an upper limit to the number
	// of bytes received and is used to tell if the byte slice represents the
	// whole file or is just the header of a file: len(raw) < limit or len(raw)>limit.
	Detector func(raw []byte, limit uint32) bool
	xmlSig   struct {
		// the local name of the root tag
		localName []byte
		// the namespace of the XML document
		xmlns []byte
	}
)

// prefix creates a Detector which returns true if any of the provided signatures
// is the prefix of the raw input.
func prefix(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		for _, s := range sigs {
			if bytes.HasPrefix(raw, s) {
				return true
			}
		}
		return false
	}
}

// offset creates a Detector which returns true if the provided signature can be
// found at offset in the raw input.
func offset(sig []byte, offset int) Detector {
	return func(raw []byte, limit uint32) bool {
		return len(raw) > offset && bytes.HasPrefix(raw[offset:], sig)
	}
}

// ciPrefix is like prefix but the check is case insensitive.
func ciPrefix(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		for _, s := range sigs {
			if ciCheck(s, raw) {
				return true
			}
		}
		return false
	}
}
func ciCheck(sig, raw []byte) bool {
	if len(raw) < len(sig)+1 {
		return false
	}
	// perform case insensitive check
	for i, b := range sig {
		db := raw[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return false
		}
	}

	return true
}

// xml creates a Detector which returns true if any of the provided XML signatures
// matches the raw input.
func xml(sigs ...xmlSig) Detector {
	return func(raw []byte, limit uint32) bool {
		b := scan.Bytes(raw)
		b.TrimLWS()
		if len(b) == 0 {
			return false
		}
		for _, s := range sigs {
			if xmlCheck(s, b) {
				return true
			}
		}
		return false
	}
}
func xmlCheck(sig xmlSig, raw []byte) bool {
	raw = raw[:min(len(raw), 512)]

	if len(sig.localName) == 0 {
		return bytes.Index(raw, sig.xmlns) > 0
	}
	if len(sig.xmlns) == 0 {
		return bytes.Index(raw, sig.localName) > 0
	}

	localNameIndex := bytes.Index(raw, sig.localName)
	return localNameIndex != -1 && localNameIndex < bytes.Index(raw, sig.xmlns)
}

// markup creates a Detector which returns true is any of the HTML signatures
// matches the raw input.
func markup(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		b := scan.Bytes(raw)
		if bytes.HasPrefix(b, []byte{0xEF, 0xBB, 0xBF}) {
			// We skip the UTF-8 BOM if present to ensure we correctly
			// process any leading whitespace. The presence of the BOM
			// is taken into account during charset detection in charset.go.
			b.Advance(3)
		}
		b.TrimLWS()
		if len(b) == 0 {
			return false
		}
		for _, s := range sigs {
			if markupCheck(s, b) {
				return true
			}
		}
		return false
	}
}
func markupCheck(sig, raw []byte) bool {
	if len(raw) < len(sig)+1 {
		return false
	}

	// perform case insensitive check
	for i, b := range sig {
		db := raw[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return false
		}
	}
	// Next byte must be space or right angle bracket.
	if db := raw[len(sig)]; !scan.ByteIsWS(db) && db != '>' {
		return false
	}

	return true
}

// ftyp creates a Detector which returns true if any of the FTYP signatures
// matches the raw input.
func ftyp(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		if len(raw) < 12 {
			return false
		}
		for _, s := range sigs {
			if bytes.Equal(raw[8:12], s) {
				return true
			}
		}
		return false
	}
}

func newXMLSig(localName, xmlns string) xmlSig {
	ret := xmlSig{xmlns: []byte(xmlns)}
	if localName != "" {
		ret.localName = []byte(fmt.Sprintf("<%s", localName))
	}

	return ret
}

// A valid shebang starts with the "#!" characters,
// followed by any number of spaces,
// followed by the path to the interpreter,
// and, optionally, followed by the arguments for the interpreter.
//
// Ex:
//
//	#! /usr/bin/env php
//
// /usr/bin/env is the interpreter, php is the first and only argument.
func shebang(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		b := scan.Bytes(raw)
		line := b.Line()
		for _, s := range sigs {
			if shebangCheck(s, line) {
				return true
			}
		}
		return false
	}
}

func shebangCheck(sig []byte, raw scan.Bytes) bool {
	if len(raw) < len(sig)+2 {
		return false
	}
	if raw[0] != '#' || raw[1] != '!' {
		return false
	}

	raw.Advance(2) // skip #! we checked above
	raw.TrimLWS()
	raw.TrimRWS()
	return bytes.Equal(raw, sig)
}
