package json

import (
	"bytes"
	"sync"
)

const (
	QueryNone    = "json"
	QueryGeo     = "geo"
	QueryHAR     = "har"
	QueryGLTF    = "gltf"
	maxRecursion = 4096
)

var queries = map[string][]query{
	QueryNone: nil,
	QueryGeo: {{
		SearchPath: [][]byte{[]byte("type")},
		SearchVals: [][]byte{
			[]byte(`"Feature"`),
			[]byte(`"FeatureCollection"`),
			[]byte(`"Point"`),
			[]byte(`"LineString"`),
			[]byte(`"Polygon"`),
			[]byte(`"MultiPoint"`),
			[]byte(`"MultiLineString"`),
			[]byte(`"MultiPolygon"`),
			[]byte(`"GeometryCollection"`),
		},
	}},
	QueryHAR: {{
		SearchPath: [][]byte{[]byte("log"), []byte("version")},
	}, {
		SearchPath: [][]byte{[]byte("log"), []byte("creator")},
	}, {
		SearchPath: [][]byte{[]byte("log"), []byte("entries")},
	}},
	QueryGLTF: {{
		SearchPath: [][]byte{[]byte("asset"), []byte("version")},
		SearchVals: [][]byte{[]byte(`"1.0"`), []byte(`"2.0"`)},
	}},
}

var parserPool = sync.Pool{
	New: func() any {
		return &parserState{maxRecursion: maxRecursion}
	},
}

// parserState holds the state of JSON parsing. The number of inspected bytes,
// the current path inside the JSON object, etc.
type parserState struct {
	// ib represents the number of inspected bytes.
	// Because mimetype limits itself to only reading the header of the file,
	// it means sometimes the input JSON can be truncated. In that case, we want
	// to still detect it as JSON, even if it's invalid/truncated.
	// When ib == len(input) it means the JSON was valid (at least the header).
	ib           int
	maxRecursion int
	// currPath keeps a track of the JSON keys parsed up.
	// It works only for JSON objects. JSON arrays are ignored
	// mainly because the functionality is not needed.
	currPath [][]byte
	// firstToken stores the first JSON token encountered in input.
	// TODO: performance would be better if we would stop parsing as soon
	// as we see that first token is not what we are interested in.
	firstToken int
	// querySatisfied is true if both path and value of any queries passed to
	// consumeAny are satisfied.
	querySatisfied bool
}

// query holds information about a combination of {"key": "val"} that we're trying
// to search for inside the JSON.
type query struct {
	// SearchPath represents the whole path to look for inside the JSON.
	// ex: [][]byte{[]byte("foo"), []byte("bar")} matches {"foo": {"bar": "baz"}}
	SearchPath [][]byte
	// SearchVals represents values to look for when the SearchPath is found.
	// Each SearchVal element is tried until one of them matches (logical OR.)
	SearchVals [][]byte
}

func eq(path1, path2 [][]byte) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i := range path1 {
		if !bytes.Equal(path1[i], path2[i]) {
			return false
		}
	}
	return true
}

// LooksLikeObjectOrArray reports if first non white space character from raw
// is either { or [. Parsing raw as JSON is a heavy operation. When receiving some
// text input we can skip parsing if the input does not even look like JSON.
func LooksLikeObjectOrArray(raw []byte) bool {
	for i := range raw {
		if isSpace(raw[i]) {
			continue
		}
		return raw[i] == '{' || raw[i] == '['
	}

	return false
}

// Parse will take out a parser from the pool depending on queryType and tries
// to parse raw bytes as JSON.
func Parse(queryType string, raw []byte) (parsed, inspected, firstToken int, querySatisfied bool) {
	p := parserPool.Get().(*parserState)
	defer func() {
		// Avoid hanging on to too much memory in extreme input cases.
		if len(p.currPath) > 128 {
			p.currPath = nil
		}
		parserPool.Put(p)
	}()
	p.reset()

	qs := queries[queryType]
	got := p.consumeAny(raw, qs, 0)
	return got, p.ib, p.firstToken, p.querySatisfied
}

func (p *parserState) reset() {
	p.ib = 0
	p.currPath = p.currPath[0:0]
	p.firstToken = TokInvalid
	p.querySatisfied = false
}

func (p *parserState) consumeSpace(b []byte) (n int) {
	for len(b) > 0 && isSpace(b[0]) {
		b = b[1:]
		n++
		p.ib++
	}
	return n
}

func (p *parserState) consumeConst(b, cnst []byte) int {
	lb := len(b)
	for i, c := range cnst {
		if lb > i && b[i] == c {
			p.ib++
		} else {
			return 0
		}
	}
	return len(cnst)
}

func (p *parserState) consumeString(b []byte) (n int) {
	var c byte
	for len(b[n:]) > 0 {
		c, n = b[n], n+1
		p.ib++
		switch c {
		case '\\':
			if len(b[n:]) == 0 {
				return 0
			}
			switch b[n] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				n++
				p.ib++
				continue
			case 'u':
				n++
				p.ib++
				for j := 0; j < 4 && len(b[n:]) > 0; j++ {
					if !isXDigit(b[n]) {
						return 0
					}
					n++
					p.ib++
				}
				continue
			default:
				return 0
			}
		case '"':
			return n
		default:
			continue
		}
	}
	return 0
}

func (p *parserState) consumeNumber(b []byte) (n int) {
	got := false
	var i int

	if len(b) == 0 {
		goto out
	}
	if b[0] == '-' {
		b, i = b[1:], i+1
		p.ib++
	}

	for len(b) > 0 {
		if !isDigit(b[0]) {
			break
		}
		got = true
		b, i = b[1:], i+1
		p.ib++
	}
	if len(b) == 0 {
		goto out
	}
	if b[0] == '.' {
		b, i = b[1:], i+1
		p.ib++
	}
	for len(b) > 0 {
		if !isDigit(b[0]) {
			break
		}
		got = true
		b, i = b[1:], i+1
		p.ib++
	}
	if len(b) == 0 {
		goto out
	}
	if got && (b[0] == 'e' || b[0] == 'E') {
		b, i = b[1:], i+1
		p.ib++
		got = false
		if len(b) == 0 {
			goto out
		}
		if b[0] == '+' || b[0] == '-' {
			b, i = b[1:], i+1
			p.ib++
		}
		for len(b) > 0 {
			if !isDigit(b[0]) {
				break
			}
			got = true
			b, i = b[1:], i+1
			p.ib++
		}
	}
out:
	if got {
		return i
	}
	return 0
}

func (p *parserState) consumeArray(b []byte, qs []query, lvl int) (n int) {
	p.appendPath([]byte{'['}, qs)
	if len(b) == 0 {
		return 0
	}

	for n < len(b) {
		n += p.consumeSpace(b[n:])
		if len(b[n:]) == 0 {
			return 0
		}
		if b[n] == ']' {
			p.ib++
			p.popLastPath(qs)
			return n + 1
		}
		innerParsed := p.consumeAny(b[n:], qs, lvl)
		if innerParsed == 0 {
			return 0
		}
		n += innerParsed
		if len(b[n:]) == 0 {
			return 0
		}
		switch b[n] {
		case ',':
			n += 1
			p.ib++
			continue
		case ']':
			p.ib++
			return n + 1
		default:
			return 0
		}
	}
	return 0
}

func queryPathMatch(qs []query, path [][]byte) int {
	for i := range qs {
		if eq(qs[i].SearchPath, path) {
			return i
		}
	}
	return -1
}

// appendPath will append a path fragment if queries is not empty.
// If we don't need query functionality (just checking if a JSON is valid),
// then we can skip keeping track of the path we're currently in.
func (p *parserState) appendPath(path []byte, qs []query) {
	if len(qs) != 0 {
		p.currPath = append(p.currPath, path)
	}
}
func (p *parserState) popLastPath(qs []query) {
	if len(qs) != 0 {
		p.currPath = p.currPath[:len(p.currPath)-1]
	}
}

func (p *parserState) consumeObject(b []byte, qs []query, lvl int) (n int) {
	for n < len(b) {
		n += p.consumeSpace(b[n:])
		if len(b[n:]) == 0 {
			return 0
		}
		if b[n] == '}' {
			p.ib++
			return n + 1
		}
		if b[n] != '"' {
			return 0
		} else {
			n += 1
			p.ib++
		}
		// queryMatched stores the index of the query satisfying the current path.
		queryMatched := -1
		if keyLen := p.consumeString(b[n:]); keyLen == 0 {
			return 0
		} else {
			p.appendPath(b[n:n+keyLen-1], qs)
			if !p.querySatisfied {
				queryMatched = queryPathMatch(qs, p.currPath)
			}
			n += keyLen
		}
		n += p.consumeSpace(b[n:])
		if len(b[n:]) == 0 {
			return 0
		}
		if b[n] != ':' {
			return 0
		} else {
			n += 1
			p.ib++
		}
		n += p.consumeSpace(b[n:])
		if len(b[n:]) == 0 {
			return 0
		}

		if valLen := p.consumeAny(b[n:], qs, lvl); valLen == 0 {
			return 0
		} else {
			if queryMatched != -1 {
				q := qs[queryMatched]
				if len(q.SearchVals) == 0 {
					p.querySatisfied = true
				}
				for _, val := range q.SearchVals {
					if bytes.Equal(val, bytes.TrimSpace(b[n:n+valLen])) {
						p.querySatisfied = true
					}
				}
			}
			n += valLen
		}
		if len(b[n:]) == 0 {
			return 0
		}
		switch b[n] {
		case ',':
			p.popLastPath(qs)
			n++
			p.ib++
			continue
		case '}':
			p.popLastPath(qs)
			p.ib++
			return n + 1
		default:
			return 0
		}
	}
	return 0
}

func (p *parserState) consumeAny(b []byte, qs []query, lvl int) (n int) {
	// Avoid too much recursion.
	if p.maxRecursion != 0 && lvl > p.maxRecursion {
		return 0
	}
	if len(qs) == 0 {
		p.querySatisfied = true
	}
	n += p.consumeSpace(b)
	if len(b[n:]) == 0 {
		return 0
	}

	var t, rv int
	switch b[n] {
	case '"':
		n++
		p.ib++
		rv = p.consumeString(b[n:])
		t = TokString
	case '[':
		n++
		p.ib++
		rv = p.consumeArray(b[n:], qs, lvl+1)
		t = TokArray
	case '{':
		n++
		p.ib++
		rv = p.consumeObject(b[n:], qs, lvl+1)
		t = TokObject
	case 't':
		rv = p.consumeConst(b[n:], []byte("true"))
		t = TokTrue
	case 'f':
		rv = p.consumeConst(b[n:], []byte("false"))
		t = TokFalse
	case 'n':
		rv = p.consumeConst(b[n:], []byte("null"))
		t = TokNull
	default:
		rv = p.consumeNumber(b[n:])
		t = TokNumber
	}
	if lvl == 0 {
		p.firstToken = t
	}
	if rv <= 0 {
		return n
	}
	n += rv
	n += p.consumeSpace(b[n:])
	return n
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isXDigit(c byte) bool {
	if isDigit(c) {
		return true
	}
	return ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

const (
	TokInvalid = 0
	TokNull    = 1 << iota
	TokTrue
	TokFalse
	TokNumber
	TokString
	TokArray
	TokObject
	TokComma
)
