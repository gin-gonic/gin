package magic

import (
	"bytes"

	"github.com/gabriel-vasile/mimetype/internal/scan"
)

var (
	// Odt matches an OpenDocument Text file.
	Odt = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.text"), 30)
	// Ott matches an OpenDocument Text Template file.
	Ott = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.text-template"), 30)
	// Ods matches an OpenDocument Spreadsheet file.
	Ods = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet"), 30)
	// Ots matches an OpenDocument Spreadsheet Template file.
	Ots = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet-template"), 30)
	// Odp matches an OpenDocument Presentation file.
	Odp = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.presentation"), 30)
	// Otp matches an OpenDocument Presentation Template file.
	Otp = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.presentation-template"), 30)
	// Odg matches an OpenDocument Drawing file.
	Odg = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.graphics"), 30)
	// Otg matches an OpenDocument Drawing Template file.
	Otg = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.graphics-template"), 30)
	// Odf matches an OpenDocument Formula file.
	Odf = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.formula"), 30)
	// Odc matches an OpenDocument Chart file.
	Odc = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.chart"), 30)
	// Epub matches an EPUB file.
	Epub = offset([]byte("mimetypeapplication/epub+zip"), 30)
	// Sxc matches an OpenOffice Spreadsheet file.
	Sxc = offset([]byte("mimetypeapplication/vnd.sun.xml.calc"), 30)
)

// Zip matches a zip archive.
func Zip(raw []byte, limit uint32) bool {
	return len(raw) > 3 &&
		raw[0] == 0x50 && raw[1] == 0x4B &&
		(raw[2] == 0x3 || raw[2] == 0x5 || raw[2] == 0x7) &&
		(raw[3] == 0x4 || raw[3] == 0x6 || raw[3] == 0x8)
}

// Jar matches a Java archive file. There are two types of Jar files:
// 1. the ones that can be opened with jexec and have 0xCAFE optional flag
// https://stackoverflow.com/tags/executable-jar/info
// 2. regular jars, same as above, just without the executable flag
// https://bugs.freebsd.org/bugzilla/show_bug.cgi?id=262278#c0
// There is an argument to only check for manifest, since it's the common nominator
// for both executable and non-executable versions. But the traversing zip entries
// is unreliable because it does linear search for signatures
// (instead of relying on offsets told by the file.)
func Jar(raw []byte, limit uint32) bool {
	return executableJar(raw) ||
		zipHas(raw, zipEntries{{
			name: []byte("META-INF/MANIFEST.MF"),
		}, {
			name: []byte("META-INF/"),
		}}, 1)
}

// KMZ matches a zipped KML file, which is "doc.kml" by convention.
func KMZ(raw []byte, _ uint32) bool {
	return zipHas(raw, zipEntries{{
		name: []byte("doc.kml"),
	}}, 100)
}

// An executable Jar has a 0xCAFE flag enabled in the first zip entry.
// The rule from file/file is:
// >(26.s+30)	leshort	0xcafe		Java archive data (JAR)
func executableJar(b scan.Bytes) bool {
	b.Advance(0x1A)
	offset, ok := b.Uint16()
	if !ok {
		return false
	}
	b.Advance(int(offset) + 2)

	cafe, ok := b.Uint16()
	return ok && cafe == 0xCAFE
}

// zipIterator iterates over a zip file returning the name of the zip entries
// in that file.
type zipIterator struct {
	b scan.Bytes
}

type zipEntries []struct {
	name []byte
	dir  bool // dir means checking just the prefix of the entry, not the whole path
}

func (z zipEntries) match(file []byte) bool {
	for i := range z {
		if z[i].dir && bytes.HasPrefix(file, z[i].name) {
			return true
		}
		if bytes.Equal(file, z[i].name) {
			return true
		}
	}
	return false
}

func zipHas(raw scan.Bytes, searchFor zipEntries, stopAfter int) bool {
	iter := zipIterator{raw}
	for i := 0; i < stopAfter; i++ {
		f := iter.next()
		if len(f) == 0 {
			break
		}
		if searchFor.match(f) {
			return true
		}
	}

	return false
}

// msoxml behaves like zipHas, but it puts restrictions on what the first zip
// entry can be.
func msoxml(raw scan.Bytes, searchFor zipEntries, stopAfter int) bool {
	iter := zipIterator{raw}
	for i := 0; i < stopAfter; i++ {
		f := iter.next()
		if len(f) == 0 {
			break
		}
		if searchFor.match(f) {
			return true
		}
		// If the first is not one of the next usually expected entries,
		// then abort this check.
		if i == 0 {
			if !bytes.Equal(f, []byte("[Content_Types].xml")) &&
				!bytes.Equal(f, []byte("_rels/.rels")) &&
				!bytes.Equal(f, []byte("docProps")) &&
				!bytes.Equal(f, []byte("customXml")) &&
				!bytes.Equal(f, []byte("[trash]")) {
				return false
			}
		}
	}

	return false
}

// next extracts the name of the next zip entry.
func (i *zipIterator) next() []byte {
	pk := []byte("PK\003\004")

	n := bytes.Index(i.b, pk)
	if n == -1 {
		return nil
	}
	i.b.Advance(n)
	if !i.b.Advance(0x1A) {
		return nil
	}
	l, ok := i.b.Uint16()
	if !ok {
		return nil
	}
	if !i.b.Advance(0x02) {
		return nil
	}
	if len(i.b) < int(l) {
		return nil
	}
	return i.b[:l]
}

// APK matches an Android Package Archive.
// The source of signatures is https://github.com/file/file/blob/1778642b8ba3d947a779a36fcd81f8e807220a19/magic/Magdir/archive#L1820-L1887
func APK(raw []byte, _ uint32) bool {
	return zipHas(raw, zipEntries{{
		name: []byte("AndroidManifest.xml"),
	}, {
		name: []byte("META-INF/com/android/build/gradle/app-metadata.properties"),
	}, {
		name: []byte("classes.dex"),
	}, {
		name: []byte("resources.arsc"),
	}, {
		name: []byte("res/drawable"),
	}}, 100)
}
