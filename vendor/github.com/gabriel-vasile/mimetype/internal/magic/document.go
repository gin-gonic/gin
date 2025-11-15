package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	// Fdf matches a Forms Data Format file.
	Fdf = prefix([]byte("%FDF"))
	// Mobi matches a Mobi file.
	Mobi = offset([]byte("BOOKMOBI"), 60)
	// Lit matches a Microsoft Lit file.
	Lit = prefix([]byte("ITOLITLS"))
)

// PDF matches a Portable Document Format file.
// The %PDF- header should be the first thing inside the file but many
// implementations don't follow the rule. The PDF spec at Appendix H says the
// signature can be prepended by anything.
// https://bugs.astron.com/view.php?id=446
func PDF(raw []byte, _ uint32) bool {
	raw = raw[:min(len(raw), 1024)]
	return bytes.Contains(raw, []byte("%PDF-"))
}

// DjVu matches a DjVu file.
func DjVu(raw []byte, _ uint32) bool {
	if len(raw) < 12 {
		return false
	}
	if !bytes.HasPrefix(raw, []byte{0x41, 0x54, 0x26, 0x54, 0x46, 0x4F, 0x52, 0x4D}) {
		return false
	}
	return bytes.HasPrefix(raw[12:], []byte("DJVM")) ||
		bytes.HasPrefix(raw[12:], []byte("DJVU")) ||
		bytes.HasPrefix(raw[12:], []byte("DJVI")) ||
		bytes.HasPrefix(raw[12:], []byte("THUM"))
}

// P7s matches an .p7s signature File (PEM, Base64).
func P7s(raw []byte, _ uint32) bool {
	// Check for PEM Encoding.
	if bytes.HasPrefix(raw, []byte("-----BEGIN PKCS7")) {
		return true
	}
	// Check if DER Encoding is long enough.
	if len(raw) < 20 {
		return false
	}
	// Magic Bytes for the signedData ASN.1 encoding.
	startHeader := [][]byte{{0x30, 0x80}, {0x30, 0x81}, {0x30, 0x82}, {0x30, 0x83}, {0x30, 0x84}}
	signedDataMatch := []byte{0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x07}
	// Check if Header is correct. There are multiple valid headers.
	for i, match := range startHeader {
		// If first bytes match, then check for ASN.1 Object Type.
		if bytes.HasPrefix(raw, match) {
			if bytes.HasPrefix(raw[i+2:], signedDataMatch) {
				return true
			}
		}
	}

	return false
}

// Lotus123 matches a Lotus 1-2-3 spreadsheet document.
func Lotus123(raw []byte, _ uint32) bool {
	if len(raw) <= 20 {
		return false
	}
	version := binary.BigEndian.Uint32(raw)
	if version == 0x00000200 {
		return raw[6] != 0 && raw[7] == 0
	}

	return version == 0x00001a00 && raw[20] > 0 && raw[20] < 32
}

// CHM matches a Microsoft Compiled HTML Help file.
func CHM(raw []byte, _ uint32) bool {
	return bytes.HasPrefix(raw, []byte("ITSF\003\000\000\000\x60\000\000\000"))
}
