package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	// SevenZ matches a 7z archive.
	SevenZ = prefix([]byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C})
	// Gzip matches gzip files based on http://www.zlib.org/rfc-gzip.html#header-trailer.
	Gzip = prefix([]byte{0x1f, 0x8b})
	// Fits matches an Flexible Image Transport System file.
	Fits = prefix([]byte{
		0x53, 0x49, 0x4D, 0x50, 0x4C, 0x45, 0x20, 0x20, 0x3D, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54,
	})
	// Xar matches an eXtensible ARchive format file.
	Xar = prefix([]byte{0x78, 0x61, 0x72, 0x21})
	// Bz2 matches a bzip2 file.
	Bz2 = prefix([]byte{0x42, 0x5A, 0x68})
	// Ar matches an ar (Unix) archive file.
	Ar = prefix([]byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E})
	// Deb matches a Debian package file.
	Deb = offset([]byte{
		0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D,
		0x62, 0x69, 0x6E, 0x61, 0x72, 0x79,
	}, 8)
	// Warc matches a Web ARChive file.
	Warc = prefix([]byte("WARC/1.0"), []byte("WARC/1.1"))
	// Cab matches a Microsoft Cabinet archive file.
	Cab = prefix([]byte("MSCF\x00\x00\x00\x00"))
	// Xz matches an xz compressed stream based on https://tukaani.org/xz/xz-file-format.txt.
	Xz = prefix([]byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00})
	// Lzip matches an Lzip compressed file.
	Lzip = prefix([]byte{0x4c, 0x5a, 0x49, 0x50})
	// RPM matches an RPM or Delta RPM package file.
	RPM = prefix([]byte{0xed, 0xab, 0xee, 0xdb}, []byte("drpm"))
	// Cpio matches a cpio archive file.
	Cpio = prefix([]byte("070707"), []byte("070701"), []byte("070702"))
	// RAR matches a RAR archive file.
	RAR = prefix([]byte("Rar!\x1A\x07\x00"), []byte("Rar!\x1A\x07\x01\x00"))
)

// InstallShieldCab matches an InstallShield Cabinet archive file.
func InstallShieldCab(raw []byte, _ uint32) bool {
	return len(raw) > 7 &&
		bytes.Equal(raw[0:4], []byte("ISc(")) &&
		raw[6] == 0 &&
		(raw[7] == 1 || raw[7] == 2 || raw[7] == 4)
}

// Zstd matches a Zstandard archive file.
// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md
func Zstd(raw []byte, limit uint32) bool {
	if len(raw) < 4 {
		return false
	}
	sig := binary.LittleEndian.Uint32(raw)
	// Check for Zstandard frames and skippable frames.
	return (sig >= 0xFD2FB522 && sig <= 0xFD2FB528) ||
		(sig >= 0x184D2A50 && sig <= 0x184D2A5F)
}

// CRX matches a Chrome extension file: a zip archive prepended by a package header.
func CRX(raw []byte, limit uint32) bool {
	const minHeaderLen = 16
	if len(raw) < minHeaderLen || !bytes.HasPrefix(raw, []byte("Cr24")) {
		return false
	}
	pubkeyLen := binary.LittleEndian.Uint32(raw[8:12])
	sigLen := binary.LittleEndian.Uint32(raw[12:16])
	zipOffset := minHeaderLen + pubkeyLen + sigLen
	if uint32(len(raw)) < zipOffset {
		return false
	}
	return Zip(raw[zipOffset:], limit)
}

// Tar matches a (t)ape (ar)chive file.
// Tar files are divided into 512 bytes records. First record contains a 257
// bytes header padded with NUL.
func Tar(raw []byte, _ uint32) bool {
	const sizeRecord = 512

	// The structure of a tar header:
	// type TarHeader struct {
	// 	Name     [100]byte
	// 	Mode     [8]byte
	// 	Uid      [8]byte
	// 	Gid      [8]byte
	// 	Size     [12]byte
	// 	Mtime    [12]byte
	// 	Chksum   [8]byte
	// 	Linkflag byte
	// 	Linkname [100]byte
	// 	Magic    [8]byte
	// 	Uname    [32]byte
	// 	Gname    [32]byte
	// 	Devmajor [8]byte
	// 	Devminor [8]byte
	// }

	if len(raw) < sizeRecord {
		return false
	}
	raw = raw[:sizeRecord]

	// First 100 bytes of the header represent the file name.
	// Check if file looks like Gentoo GLEP binary package.
	if bytes.Contains(raw[:100], []byte("/gpkg-1\x00")) {
		return false
	}

	// Get the checksum recorded into the file.
	recsum := tarParseOctal(raw[148:156])
	if recsum == -1 {
		return false
	}
	sum1, sum2 := tarChksum(raw)
	return recsum == sum1 || recsum == sum2
}

// tarParseOctal converts octal string to decimal int.
func tarParseOctal(b []byte) int64 {
	// Because unused fields are filled with NULs, we need to skip leading NULs.
	// Fields may also be padded with spaces or NULs.
	// So we remove leading and trailing NULs and spaces to be sure.
	b = bytes.Trim(b, " \x00")

	if len(b) == 0 {
		return -1
	}
	ret := int64(0)
	for _, b := range b {
		if b == 0 {
			break
		}
		if b < '0' || b > '7' {
			return -1
		}
		ret = (ret << 3) | int64(b-'0')
	}
	return ret
}

// tarChksum computes the checksum for the header block b.
// The actual checksum is written to same b block after it has been calculated.
// Before calculation the bytes from b reserved for checksum have placeholder
// value of ASCII space 0x20.
// POSIX specifies a sum of the unsigned byte values, but the Sun tar used
// signed byte values. We compute and return both.
func tarChksum(b []byte) (unsigned, signed int64) {
	for i, c := range b {
		if 148 <= i && i < 156 {
			c = ' ' // Treat the checksum field itself as all spaces.
		}
		unsigned += int64(c)
		signed += int64(int8(c))
	}
	return unsigned, signed
}
