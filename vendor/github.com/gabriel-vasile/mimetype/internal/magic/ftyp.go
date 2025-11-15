package magic

import (
	"bytes"
)

var (
	// AVIF matches an AV1 Image File Format still or animated.
	// Wikipedia page seems outdated listing image/avif-sequence for animations.
	// https://github.com/AOMediaCodec/av1-avif/issues/59
	AVIF = ftyp([]byte("avif"), []byte("avis"))
	// ThreeGP matches a 3GPP file.
	ThreeGP = ftyp(
		[]byte("3gp1"), []byte("3gp2"), []byte("3gp3"), []byte("3gp4"),
		[]byte("3gp5"), []byte("3gp6"), []byte("3gp7"), []byte("3gs7"),
		[]byte("3ge6"), []byte("3ge7"), []byte("3gg6"),
	)
	// ThreeG2 matches a 3GPP2 file.
	ThreeG2 = ftyp(
		[]byte("3g24"), []byte("3g25"), []byte("3g26"), []byte("3g2a"),
		[]byte("3g2b"), []byte("3g2c"), []byte("KDDI"),
	)
	// AMp4 matches an audio MP4 file.
	AMp4 = ftyp(
		// audio for Adobe Flash Player 9+
		[]byte("F4A "), []byte("F4B "),
		// Apple iTunes AAC-LC (.M4A) Audio
		[]byte("M4B "), []byte("M4P "),
		// MPEG-4 (.MP4) for SonyPSP
		[]byte("MSNV"),
		// Nero Digital AAC Audio
		[]byte("NDAS"),
	)
	// Mqv matches a Sony / Mobile QuickTime  file.
	Mqv = ftyp([]byte("mqt "))
	// M4a matches an audio M4A file.
	M4a = ftyp([]byte("M4A "))
	// M4v matches an Appl4 M4V video file.
	M4v = ftyp([]byte("M4V "), []byte("M4VH"), []byte("M4VP"))
	// Heic matches a High Efficiency Image Coding (HEIC) file.
	Heic = ftyp([]byte("heic"), []byte("heix"))
	// HeicSequence matches a High Efficiency Image Coding (HEIC) file sequence.
	HeicSequence = ftyp([]byte("hevc"), []byte("hevx"))
	// Heif matches a High Efficiency Image File Format (HEIF) file.
	Heif = ftyp([]byte("mif1"), []byte("heim"), []byte("heis"), []byte("avic"))
	// HeifSequence matches a High Efficiency Image File Format (HEIF) file sequence.
	HeifSequence = ftyp([]byte("msf1"), []byte("hevm"), []byte("hevs"), []byte("avcs"))
	// Mj2 matches a Motion JPEG 2000 file: https://en.wikipedia.org/wiki/Motion_JPEG_2000.
	Mj2 = ftyp([]byte("mj2s"), []byte("mjp2"), []byte("MFSM"), []byte("MGSV"))
	// Dvb matches a Digital Video Broadcasting file: https://dvb.org.
	// https://cconcolato.github.io/mp4ra/filetype.html
	// https://github.com/file/file/blob/512840337ead1076519332d24fefcaa8fac36e06/magic/Magdir/animation#L135-L154
	Dvb = ftyp(
		[]byte("dby1"), []byte("dsms"), []byte("dts1"), []byte("dts2"),
		[]byte("dts3"), []byte("dxo "), []byte("dmb1"), []byte("dmpf"),
		[]byte("drc1"), []byte("dv1a"), []byte("dv1b"), []byte("dv2a"),
		[]byte("dv2b"), []byte("dv3a"), []byte("dv3b"), []byte("dvr1"),
		[]byte("dvt1"), []byte("emsg"))
	// TODO: add support for remaining video formats at ftyps.com.
)

// QuickTime matches a QuickTime File Format file.
// https://www.loc.gov/preservation/digital/formats/fdd/fdd000052.shtml
// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap1/qtff1.html#//apple_ref/doc/uid/TP40000939-CH203-38190
// https://github.com/apache/tika/blob/0f5570691133c75ac4472c3340354a6c4080b104/tika-core/src/main/resources/org/apache/tika/mime/tika-mimetypes.xml#L7758-L7777
func QuickTime(raw []byte, _ uint32) bool {
	if len(raw) < 12 {
		return false
	}
	// First 4 bytes represent the size of the atom as unsigned int.
	// Next 4 bytes are the type of the atom.
	// For `ftyp` atoms check if first byte in size is 0, otherwise, a text file
	// which happens to contain 'ftypqt  ' at index 4 will trigger a false positive.
	if bytes.Equal(raw[4:12], []byte("ftypqt  ")) ||
		bytes.Equal(raw[4:12], []byte("ftypmoov")) {
		return raw[0] == 0x00
	}
	basicAtomTypes := [][]byte{
		[]byte("moov\x00"),
		[]byte("mdat\x00"),
		[]byte("free\x00"),
		[]byte("skip\x00"),
		[]byte("pnot\x00"),
	}
	for _, a := range basicAtomTypes {
		if bytes.Equal(raw[4:9], a) {
			return true
		}
	}
	return bytes.Equal(raw[:8], []byte("\x00\x00\x00\x08wide"))
}

// Mp4 detects an .mp4 file. Mp4 detections only does a basic ftyp check.
// Mp4 has many registered and unregistered code points so it's hard to keep track
// of all. Detection will default on video/mp4 for all ftyp files.
// ISO_IEC_14496-12 is the specification for the iso container.
func Mp4(raw []byte, _ uint32) bool {
	if len(raw) < 12 {
		return false
	}
	// ftyps are made out of boxes. The first 4 bytes of the box represent
	// its size in big-endian uint32. First box is the ftyp box and it is small
	// in size. Check most significant byte is 0 to filter out false positive
	// text files that happen to contain the string "ftyp" at index 4.
	if raw[0] != 0 {
		return false
	}
	return bytes.Equal(raw[4:8], []byte("ftyp"))
}
