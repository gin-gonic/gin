package magic

import "bytes"

var (
	// Png matches a Portable Network Graphics file.
	// https://www.w3.org/TR/PNG/
	Png = prefix([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// Apng matches an Animated Portable Network Graphics file.
	// https://wiki.mozilla.org/APNG_Specification
	Apng = offset([]byte("acTL"), 37)
	// Jpg matches a Joint Photographic Experts Group file.
	Jpg = prefix([]byte{0xFF, 0xD8, 0xFF})
	// Jp2 matches a JPEG 2000 Image file (ISO 15444-1).
	Jp2 = jpeg2k([]byte{0x6a, 0x70, 0x32, 0x20})
	// Jpx matches a JPEG 2000 Image file (ISO 15444-2).
	Jpx = jpeg2k([]byte{0x6a, 0x70, 0x78, 0x20})
	// Jpm matches a JPEG 2000 Image file (ISO 15444-6).
	Jpm = jpeg2k([]byte{0x6a, 0x70, 0x6D, 0x20})
	// Gif matches a Graphics Interchange Format file.
	Gif = prefix([]byte("GIF87a"), []byte("GIF89a"))
	// Bmp matches a bitmap image file.
	Bmp = prefix([]byte{0x42, 0x4D})
	// Ps matches a PostScript file.
	Ps = prefix([]byte("%!PS-Adobe-"))
	// Psd matches a Photoshop Document file.
	Psd = prefix([]byte("8BPS"))
	// Ico matches an ICO file.
	Ico = prefix([]byte{0x00, 0x00, 0x01, 0x00}, []byte{0x00, 0x00, 0x02, 0x00})
	// Icns matches an ICNS (Apple Icon Image format) file.
	Icns = prefix([]byte("icns"))
	// Tiff matches a Tagged Image File Format file.
	Tiff = prefix([]byte{0x49, 0x49, 0x2A, 0x00}, []byte{0x4D, 0x4D, 0x00, 0x2A})
	// Bpg matches a Better Portable Graphics file.
	Bpg = prefix([]byte{0x42, 0x50, 0x47, 0xFB})
	// Xcf matches GIMP image data.
	Xcf = prefix([]byte("gimp xcf"))
	// Pat matches GIMP pattern data.
	Pat = offset([]byte("GPAT"), 20)
	// Gbr matches GIMP brush data.
	Gbr = offset([]byte("GIMP"), 20)
	// Hdr matches Radiance HDR image.
	// https://web.archive.org/web/20060913152809/http://local.wasp.uwa.edu.au/~pbourke/dataformats/pic/
	Hdr = prefix([]byte("#?RADIANCE\n"))
	// Xpm matches X PixMap image data.
	Xpm = prefix([]byte{0x2F, 0x2A, 0x20, 0x58, 0x50, 0x4D, 0x20, 0x2A, 0x2F})
	// Jxs matches a JPEG XS coded image file (ISO/IEC 21122-3).
	Jxs = prefix([]byte{0x00, 0x00, 0x00, 0x0C, 0x4A, 0x58, 0x53, 0x20, 0x0D, 0x0A, 0x87, 0x0A})
	// Jxr matches Microsoft HD JXR photo file.
	Jxr = prefix([]byte{0x49, 0x49, 0xBC, 0x01})
)

func jpeg2k(sig []byte) Detector {
	return func(raw []byte, _ uint32) bool {
		if len(raw) < 24 {
			return false
		}

		if !bytes.Equal(raw[4:8], []byte{0x6A, 0x50, 0x20, 0x20}) &&
			!bytes.Equal(raw[4:8], []byte{0x6A, 0x50, 0x32, 0x20}) {
			return false
		}
		return bytes.Equal(raw[20:24], sig)
	}
}

// Webp matches a WebP file.
func Webp(raw []byte, _ uint32) bool {
	return len(raw) > 12 &&
		bytes.Equal(raw[0:4], []byte("RIFF")) &&
		bytes.Equal(raw[8:12], []byte{0x57, 0x45, 0x42, 0x50})
}

// Dwg matches a CAD drawing file.
func Dwg(raw []byte, _ uint32) bool {
	if len(raw) < 6 || raw[0] != 0x41 || raw[1] != 0x43 {
		return false
	}
	dwgVersions := [][]byte{
		{0x31, 0x2E, 0x34, 0x30},
		{0x31, 0x2E, 0x35, 0x30},
		{0x32, 0x2E, 0x31, 0x30},
		{0x31, 0x30, 0x30, 0x32},
		{0x31, 0x30, 0x30, 0x33},
		{0x31, 0x30, 0x30, 0x34},
		{0x31, 0x30, 0x30, 0x36},
		{0x31, 0x30, 0x30, 0x39},
		{0x31, 0x30, 0x31, 0x32},
		{0x31, 0x30, 0x31, 0x34},
		{0x31, 0x30, 0x31, 0x35},
		{0x31, 0x30, 0x31, 0x38},
		{0x31, 0x30, 0x32, 0x31},
		{0x31, 0x30, 0x32, 0x34},
		{0x31, 0x30, 0x33, 0x32},
	}

	for _, d := range dwgVersions {
		if bytes.Equal(raw[2:6], d) {
			return true
		}
	}

	return false
}

// Jxl matches JPEG XL image file.
func Jxl(raw []byte, _ uint32) bool {
	return bytes.HasPrefix(raw, []byte{0xFF, 0x0A}) ||
		bytes.HasPrefix(raw, []byte("\x00\x00\x00\x0cJXL\x20\x0d\x0a\x87\x0a"))
}
