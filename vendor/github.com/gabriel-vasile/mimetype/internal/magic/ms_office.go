package magic

import (
	"bytes"
	"encoding/binary"
)

// Xlsx matches a Microsoft Excel 2007 file.
func Xlsx(raw []byte, limit uint32) bool {
	return msoxml(raw, zipEntries{{
		name: []byte("xl/"),
		dir:  true,
	}}, 100)
}

// Docx matches a Microsoft Word 2007 file.
func Docx(raw []byte, limit uint32) bool {
	return msoxml(raw, zipEntries{{
		name: []byte("word/"),
		dir:  true,
	}}, 100)
}

// Pptx matches a Microsoft PowerPoint 2007 file.
func Pptx(raw []byte, limit uint32) bool {
	return msoxml(raw, zipEntries{{
		name: []byte("ppt/"),
		dir:  true,
	}}, 100)
}

// Visio matches a Microsoft Visio 2013+ file.
func Visio(raw []byte, limit uint32) bool {
	return msoxml(raw, zipEntries{{
		name: []byte("visio/"),
		dir:  true,
	}}, 100)
}

// Ole matches an Open Linking and Embedding file.
//
// https://en.wikipedia.org/wiki/Object_Linking_and_Embedding
func Ole(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
}

// Aaf matches an Advanced Authoring Format file.
// See: https://pyaaf.readthedocs.io/en/latest/about.html
// See: https://en.wikipedia.org/wiki/Advanced_Authoring_Format
func Aaf(raw []byte, limit uint32) bool {
	if len(raw) < 31 {
		return false
	}
	return bytes.HasPrefix(raw[8:], []byte{0x41, 0x41, 0x46, 0x42, 0x0D, 0x00, 0x4F, 0x4D}) &&
		(raw[30] == 0x09 || raw[30] == 0x0C)
}

// Doc matches a Microsoft Word 97-2003 file.
// See: https://github.com/decalage2/oletools/blob/412ee36ae45e70f42123e835871bac956d958461/oletools/common/clsid.py
func Doc(raw []byte, _ uint32) bool {
	clsids := [][]byte{
		// Microsoft Word 97-2003 Document (Word.Document.8)
		{0x06, 0x09, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
		// Microsoft Word 6.0-7.0 Document (Word.Document.6)
		{0x00, 0x09, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
		// Microsoft Word Picture (Word.Picture.8)
		{0x07, 0x09, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
	}

	for _, clsid := range clsids {
		if matchOleClsid(raw, clsid) {
			return true
		}
	}

	return false
}

// Ppt matches a Microsoft PowerPoint 97-2003 file or a PowerPoint 95 presentation.
func Ppt(raw []byte, limit uint32) bool {
	// Root CLSID test is the safest way to detect identify OLE, however, the format
	// often places the root CLSID at the end of the file.
	if matchOleClsid(raw, []byte{
		0x10, 0x8d, 0x81, 0x64, 0x9b, 0x4f, 0xcf, 0x11,
		0x86, 0xea, 0x00, 0xaa, 0x00, 0xb9, 0x29, 0xe8,
	}) || matchOleClsid(raw, []byte{
		0x70, 0xae, 0x7b, 0xea, 0x3b, 0xfb, 0xcd, 0x11,
		0xa9, 0x03, 0x00, 0xaa, 0x00, 0x51, 0x0e, 0xa3,
	}) {
		return true
	}

	lin := len(raw)
	if lin < 520 {
		return false
	}
	pptSubHeaders := [][]byte{
		{0xA0, 0x46, 0x1D, 0xF0},
		{0x00, 0x6E, 0x1E, 0xF0},
		{0x0F, 0x00, 0xE8, 0x03},
	}
	for _, h := range pptSubHeaders {
		if bytes.HasPrefix(raw[512:], h) {
			return true
		}
	}

	if bytes.HasPrefix(raw[512:], []byte{0xFD, 0xFF, 0xFF, 0xFF}) &&
		raw[518] == 0x00 && raw[519] == 0x00 {
		return true
	}

	return lin > 1152 && bytes.Contains(raw[1152:min(4096, lin)],
		[]byte("P\x00o\x00w\x00e\x00r\x00P\x00o\x00i\x00n\x00t\x00 D\x00o\x00c\x00u\x00m\x00e\x00n\x00t"))
}

// Xls matches a Microsoft Excel 97-2003 file.
func Xls(raw []byte, limit uint32) bool {
	// Root CLSID test is the safest way to detect identify OLE, however, the format
	// often places the root CLSID at the end of the file.
	if matchOleClsid(raw, []byte{
		0x10, 0x08, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) || matchOleClsid(raw, []byte{
		0x20, 0x08, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) {
		return true
	}

	lin := len(raw)
	if lin < 520 {
		return false
	}
	xlsSubHeaders := [][]byte{
		{0x09, 0x08, 0x10, 0x00, 0x00, 0x06, 0x05, 0x00},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x10},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x1F},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x22},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x23},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x28},
		{0xFD, 0xFF, 0xFF, 0xFF, 0x29},
	}
	for _, h := range xlsSubHeaders {
		if bytes.HasPrefix(raw[512:], h) {
			return true
		}
	}

	return lin > 1152 && bytes.Contains(raw[1152:min(4096, lin)],
		[]byte("W\x00k\x00s\x00S\x00S\x00W\x00o\x00r\x00k\x00B\x00o\x00o\x00k"))
}

// Pub matches a Microsoft Publisher file.
func Pub(raw []byte, limit uint32) bool {
	return matchOleClsid(raw, []byte{
		0x01, 0x12, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46,
	})
}

// Msg matches a Microsoft Outlook email file.
func Msg(raw []byte, limit uint32) bool {
	return matchOleClsid(raw, []byte{
		0x0B, 0x0D, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46,
	})
}

// Msi matches a Microsoft Windows Installer file.
// http://fileformats.archiveteam.org/wiki/Microsoft_Compound_File
func Msi(raw []byte, limit uint32) bool {
	return matchOleClsid(raw, []byte{
		0x84, 0x10, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46,
	})
}

// One matches a Microsoft OneNote file.
func One(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{
		0xe4, 0x52, 0x5c, 0x7b, 0x8c, 0xd8, 0xa7, 0x4d,
		0xae, 0xb1, 0x53, 0x78, 0xd0, 0x29, 0x96, 0xd3,
	})
}

// Helper to match by a specific CLSID of a compound file.
//
// http://fileformats.archiveteam.org/wiki/Microsoft_Compound_File
func matchOleClsid(in []byte, clsid []byte) bool {
	// Microsoft Compound files v3 have a sector length of 512, while v4 has 4096.
	// Change sector offset depending on file version.
	// https://www.loc.gov/preservation/digital/formats/fdd/fdd000392.shtml
	sectorLength := 512
	if len(in) < sectorLength {
		return false
	}
	if in[26] == 0x04 && in[27] == 0x00 {
		sectorLength = 4096
	}

	// SecID of first sector of the directory stream.
	firstSecID := int(binary.LittleEndian.Uint32(in[48:52]))

	// Expected offset of CLSID for root storage object.
	clsidOffset := sectorLength*(1+firstSecID) + 80

	if len(in) <= clsidOffset+16 {
		return false
	}

	return bytes.HasPrefix(in[clsidOffset:], clsid)
}
