package magic

import (
	"bytes"
)

var (
	// Woff matches a Web Open Font Format file.
	Woff = prefix([]byte("wOFF"))
	// Woff2 matches a Web Open Font Format version 2 file.
	Woff2 = prefix([]byte("wOF2"))
	// Otf matches an OpenType font file.
	Otf = prefix([]byte{0x4F, 0x54, 0x54, 0x4F, 0x00})
)

// Ttf matches a TrueType font file.
func Ttf(raw []byte, limit uint32) bool {
	if !bytes.HasPrefix(raw, []byte{0x00, 0x01, 0x00, 0x00}) {
		return false
	}
	return !MsAccessAce(raw, limit) && !MsAccessMdb(raw, limit)
}

// Eot matches an Embedded OpenType font file.
func Eot(raw []byte, limit uint32) bool {
	return len(raw) > 35 &&
		bytes.Equal(raw[34:36], []byte{0x4C, 0x50}) &&
		(bytes.Equal(raw[8:11], []byte{0x02, 0x00, 0x01}) ||
			bytes.Equal(raw[8:11], []byte{0x01, 0x00, 0x00}) ||
			bytes.Equal(raw[8:11], []byte{0x02, 0x00, 0x02}))
}

// Ttc matches a TrueType Collection font file.
func Ttc(raw []byte, limit uint32) bool {
	return len(raw) > 7 &&
		bytes.HasPrefix(raw, []byte("ttcf")) &&
		(bytes.Equal(raw[4:8], []byte{0x00, 0x01, 0x00, 0x00}) ||
			bytes.Equal(raw[4:8], []byte{0x00, 0x02, 0x00, 0x00}))
}
