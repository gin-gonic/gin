package mimetype

import (
	"mime"

	"github.com/gabriel-vasile/mimetype/internal/charset"
	"github.com/gabriel-vasile/mimetype/internal/magic"
)

// MIME struct holds information about a file format: the string representation
// of the MIME type, the extension and the parent file format.
type MIME struct {
	mime      string
	aliases   []string
	extension string
	// detector receives the raw input and a limit for the number of bytes it is
	// allowed to check. It returns whether the input matches a signature or not.
	detector magic.Detector
	children []*MIME
	parent   *MIME
}

// String returns the string representation of the MIME type, e.g., "application/zip".
func (m *MIME) String() string {
	return m.mime
}

// Extension returns the file extension associated with the MIME type.
// It includes the leading dot, as in ".html". When the file format does not
// have an extension, the empty string is returned.
func (m *MIME) Extension() string {
	return m.extension
}

// Parent returns the parent MIME type from the hierarchy.
// Each MIME type has a non-nil parent, except for the root MIME type.
//
// For example, the application/json and text/html MIME types have text/plain as
// their parent because they are text files who happen to contain JSON or HTML.
// Another example is the ZIP format, which is used as container
// for Microsoft Office files, EPUB files, JAR files, and others.
func (m *MIME) Parent() *MIME {
	return m.parent
}

// Is checks whether this MIME type, or any of its aliases, is equal to the
// expected MIME type. MIME type equality test is done on the "type/subtype"
// section, ignores any optional MIME parameters, ignores any leading and
// trailing whitespace, and is case insensitive.
func (m *MIME) Is(expectedMIME string) bool {
	// Parsing is needed because some detected MIME types contain parameters
	// that need to be stripped for the comparison.
	expectedMIME, _, _ = mime.ParseMediaType(expectedMIME)
	found, _, _ := mime.ParseMediaType(m.mime)

	if expectedMIME == found {
		return true
	}

	for _, alias := range m.aliases {
		if alias == expectedMIME {
			return true
		}
	}

	return false
}

func newMIME(
	mime, extension string,
	detector magic.Detector,
	children ...*MIME) *MIME {
	m := &MIME{
		mime:      mime,
		extension: extension,
		detector:  detector,
		children:  children,
	}

	for _, c := range children {
		c.parent = m
	}

	return m
}

func (m *MIME) alias(aliases ...string) *MIME {
	m.aliases = aliases
	return m
}

// match does a depth-first search on the signature tree. It returns the deepest
// successful node for which all the children detection functions fail.
func (m *MIME) match(in []byte, readLimit uint32) *MIME {
	for _, c := range m.children {
		if c.detector(in, readLimit) {
			return c.match(in, readLimit)
		}
	}

	needsCharset := map[string]func([]byte) string{
		"text/plain": charset.FromPlain,
		"text/html":  charset.FromHTML,
		"text/xml":   charset.FromXML,
	}
	charset := ""
	if f, ok := needsCharset[m.mime]; ok {
		// The charset comes from BOM, from HTML headers, from XML headers.
		// Limit the number of bytes searched for to 1024.
		charset = f(in[:min(len(in), 1024)])
	}
	if m == root {
		return m
	}

	return m.cloneHierarchy(charset)
}

// flatten transforms an hierarchy of MIMEs into a slice of MIMEs.
func (m *MIME) flatten() []*MIME {
	out := []*MIME{m}
	for _, c := range m.children {
		out = append(out, c.flatten()...)
	}

	return out
}

// clone creates a new MIME with the provided optional MIME parameters.
func (m *MIME) clone(charset string) *MIME {
	clonedMIME := m.mime
	if charset != "" {
		clonedMIME = m.mime + "; charset=" + charset
	}

	return &MIME{
		mime:      clonedMIME,
		aliases:   m.aliases,
		extension: m.extension,
	}
}

// cloneHierarchy creates a clone of m and all its ancestors. The optional MIME
// parameters are set on the last child of the hierarchy.
func (m *MIME) cloneHierarchy(charset string) *MIME {
	ret := m.clone(charset)
	lastChild := ret
	for p := m.Parent(); p != nil; p = p.Parent() {
		pClone := p.clone("")
		lastChild.parent = pClone
		lastChild = pClone
	}

	return ret
}

func (m *MIME) lookup(mime string) *MIME {
	for _, n := range append(m.aliases, m.mime) {
		if n == mime {
			return m
		}
	}

	for _, c := range m.children {
		if m := c.lookup(mime); m != nil {
			return m
		}
	}
	return nil
}

// Extend adds detection for a sub-format. The detector is a function
// returning true when the raw input file satisfies a signature.
// The sub-format will be detected if all the detectors in the parent chain return true.
// The extension should include the leading dot, as in ".html".
func (m *MIME) Extend(detector func(raw []byte, limit uint32) bool, mime, extension string, aliases ...string) {
	c := &MIME{
		mime:      mime,
		extension: extension,
		detector:  detector,
		parent:    m,
		aliases:   aliases,
	}

	mu.Lock()
	m.children = append([]*MIME{c}, m.children...)
	mu.Unlock()
}
