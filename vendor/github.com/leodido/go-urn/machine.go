package urn

import (
	"fmt"

	scimschema "github.com/leodido/go-urn/scim/schema"
)

var (
	errPrefix              = "expecting the prefix to be the \"urn\" string (whatever case) [col %d]"
	errIdentifier          = "expecting the identifier to be string (1..31 alnum chars, also containing dashes but not at its beginning) [col %d]"
	errSpecificString      = "expecting the specific string to be a string containing alnum, hex, or others ([()+,-.:=@;$_!*']) chars [col %d]"
	errNoUrnWithinID       = "expecting the identifier to not contain the \"urn\" reserved string [col %d]"
	errHex                 = "expecting the percent encoded chars to be well-formed (%%alnum{2}) [col %d]"
	errSCIMNamespace       = "expecing the SCIM namespace identifier (ietf:params:scim) [col %d]"
	errSCIMType            = "expecting a correct SCIM type (schemas, api, param) [col %d]"
	errSCIMName            = "expecting one or more alnum char in the SCIM name part [col %d]"
	errSCIMOther           = "expecting a well-formed other SCIM part [col %d]"
	errSCIMOtherIncomplete = "expecting a not empty SCIM other part after colon [col %d]"
	err8141InformalID      = "informal URN namespace must be in the form urn-[1-9][0-9] [col %d]"
	err8141SpecificString  = "expecting the specific string to contain alnum, hex, or others ([~&()+,-.:=@;$_!*'] or [/?] not in first position) chars [col %d]"
	err8141Identifier      = "expecting the indentifier to be a string with (length 2 to 32 chars) containing alnum (or dashes) not starting or ending with a dash [col %d]"
	err8141RComponentStart = "expecting only one r-component (starting with the ?+ sequence) [col %d]"
	err8141QComponentStart = "expecting only one q-component (starting with the ?= sequence) [col %d]"
	err8141MalformedRComp  = "expecting a non-empty r-component containing alnum, hex, or others ([~&()+,-.:=@;$_!*'] or [/?] but not at its beginning) [col %d]"
	err8141MalformedQComp  = "expecting a non-empty q-component containing alnum, hex, or others ([~&()+,-.:=@;$_!*'] or [/?] but not at its beginning) [col %d]"
)
var _toStateActions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 33, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0,
}

var _eofActions []byte = []byte{
	0, 1, 1, 1, 1, 4, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 8, 9,
	9, 4, 4, 11, 1, 1, 1, 1,
	12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12,
	12, 14, 14, 14, 14, 16, 18, 20,
	20, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 1, 1, 1, 1, 21,
	22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22,
	23, 24, 24, 25, 25, 0, 26, 28,
	28, 29, 29, 30, 30, 26, 26, 31,
	31, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 21,
	21, 22, 22, 22, 34, 34, 35, 37,
	37, 38, 40, 41, 41, 38, 42, 42,
	42, 44, 42, 48, 48, 48, 50, 44,
	50, 0,
}

const start int = 1
const firstFinal int = 172

const enScimOnly int = 44
const enRfc8141Only int = 83
const enFail int = 193
const enMain int = 1

// Machine is the interface representing the FSM
type Machine interface {
	Error() error
	Parse(input []byte) (*URN, error)
	WithParsingMode(ParsingMode)
}

type machine struct {
	data           []byte
	cs             int
	p, pe, eof, pb int
	err            error
	startParsingAt int
	parsingMode    ParsingMode
	parsingModeSet bool
}

// NewMachine creates a new FSM able to parse RFC 2141 strings.
func NewMachine(options ...Option) Machine {
	m := &machine{
		parsingModeSet: false,
	}

	for _, o := range options {
		o(m)
	}
	// Set default parsing mode
	if !m.parsingModeSet {
		m.WithParsingMode(DefaultParsingMode)
	}

	return m
}

// Err returns the error that occurred on the last call to Parse.
//
// If the result is nil, then the line was parsed successfully.
func (m *machine) Error() error {
	return m.err
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

// Parse parses the input byte array as a RFC 2141 or RFC7643 string.
func (m *machine) Parse(input []byte) (*URN, error) {
	m.data = input
	m.p = 0
	m.pb = 0
	m.pe = len(input)
	m.eof = len(input)
	m.err = nil
	m.cs = m.startParsingAt
	output := &URN{
		tolower: []int{},
	}
	{
		if (m.p) == (m.pe) {
			goto _testEof
		}
		if m.cs == 0 {
			goto _out
		}
	_resume:
		switch m.cs {
		case 1:
			switch (m.data)[(m.p)] {
			case 85:
				goto tr1
			case 117:
				goto tr1
			}
			goto tr0
		case 0:
			goto _out
		case 2:
			switch (m.data)[(m.p)] {
			case 82:
				goto tr2
			case 114:
				goto tr2
			}
			goto tr0
		case 3:
			switch (m.data)[(m.p)] {
			case 78:
				goto tr3
			case 110:
				goto tr3
			}
			goto tr0
		case 4:
			if (m.data)[(m.p)] == 58 {
				goto tr4
			}
			goto tr0
		case 5:
			switch (m.data)[(m.p)] {
			case 85:
				goto tr7
			case 117:
				goto tr7
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr6
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr6
				}
			default:
				goto tr6
			}
			goto tr5
		case 6:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr9
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr9
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr9
				}
			default:
				goto tr9
			}
			goto tr8
		case 7:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr11
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr11
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr11
				}
			default:
				goto tr11
			}
			goto tr8
		case 8:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr12
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr12
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr12
				}
			default:
				goto tr12
			}
			goto tr8
		case 9:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr13
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr13
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr13
				}
			default:
				goto tr13
			}
			goto tr8
		case 10:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr14
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr14
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr14
				}
			default:
				goto tr14
			}
			goto tr8
		case 11:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr15
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr15
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr15
				}
			default:
				goto tr15
			}
			goto tr8
		case 12:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr16
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr16
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr16
				}
			default:
				goto tr16
			}
			goto tr8
		case 13:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr17
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr17
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr17
				}
			default:
				goto tr17
			}
			goto tr8
		case 14:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr18
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr18
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr18
				}
			default:
				goto tr18
			}
			goto tr8
		case 15:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr19
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr19
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr19
				}
			default:
				goto tr19
			}
			goto tr8
		case 16:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr20
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr20
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr20
				}
			default:
				goto tr20
			}
			goto tr8
		case 17:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr21
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr21
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr21
				}
			default:
				goto tr21
			}
			goto tr8
		case 18:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr22
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr22
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr22
				}
			default:
				goto tr22
			}
			goto tr8
		case 19:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr23
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr23
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr23
				}
			default:
				goto tr23
			}
			goto tr8
		case 20:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr24
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr24
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr24
				}
			default:
				goto tr24
			}
			goto tr8
		case 21:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr25
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr25
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr25
				}
			default:
				goto tr25
			}
			goto tr8
		case 22:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr26
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr26
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr26
				}
			default:
				goto tr26
			}
			goto tr8
		case 23:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr27
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr27
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr27
				}
			default:
				goto tr27
			}
			goto tr8
		case 24:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr28
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr28
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr28
				}
			default:
				goto tr28
			}
			goto tr8
		case 25:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr29
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr29
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr29
				}
			default:
				goto tr29
			}
			goto tr8
		case 26:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr30
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr30
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr30
				}
			default:
				goto tr30
			}
			goto tr8
		case 27:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr31
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr31
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr31
				}
			default:
				goto tr31
			}
			goto tr8
		case 28:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr32
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr32
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr32
				}
			default:
				goto tr32
			}
			goto tr8
		case 29:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr33
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr33
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr33
				}
			default:
				goto tr33
			}
			goto tr8
		case 30:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr34
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr34
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr34
				}
			default:
				goto tr34
			}
			goto tr8
		case 31:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr35
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr35
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr35
				}
			default:
				goto tr35
			}
			goto tr8
		case 32:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr36
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr36
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr36
				}
			default:
				goto tr36
			}
			goto tr8
		case 33:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr37
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr37
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr37
				}
			default:
				goto tr37
			}
			goto tr8
		case 34:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr38
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr38
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr38
				}
			default:
				goto tr38
			}
			goto tr8
		case 35:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr39
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr39
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr39
				}
			default:
				goto tr39
			}
			goto tr8
		case 36:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr40
			case 58:
				goto tr10
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr40
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr40
				}
			default:
				goto tr40
			}
			goto tr8
		case 37:
			if (m.data)[(m.p)] == 58 {
				goto tr10
			}
			goto tr8
		case 38:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr42
			case 36:
				goto tr42
			case 37:
				goto tr43
			case 61:
				goto tr42
			case 95:
				goto tr42
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr42
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr42
					}
				case (m.data)[(m.p)] >= 64:
					goto tr42
				}
			default:
				goto tr42
			}
			goto tr41
		case 172:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr212
			case 36:
				goto tr212
			case 37:
				goto tr213
			case 61:
				goto tr212
			case 95:
				goto tr212
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr212
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr212
					}
				case (m.data)[(m.p)] >= 64:
					goto tr212
				}
			default:
				goto tr212
			}
			goto tr41
		case 39:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr45
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr45
				}
			default:
				goto tr46
			}
			goto tr44
		case 40:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr47
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr47
				}
			default:
				goto tr48
			}
			goto tr44
		case 173:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr212
			case 36:
				goto tr212
			case 37:
				goto tr213
			case 61:
				goto tr212
			case 95:
				goto tr212
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr212
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr212
					}
				case (m.data)[(m.p)] >= 64:
					goto tr212
				}
			default:
				goto tr212
			}
			goto tr44
		case 41:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr9
			case 58:
				goto tr10
			case 82:
				goto tr49
			case 114:
				goto tr49
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr9
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr9
				}
			default:
				goto tr9
			}
			goto tr5
		case 42:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr11
			case 58:
				goto tr10
			case 78:
				goto tr50
			case 110:
				goto tr50
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr11
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr11
				}
			default:
				goto tr11
			}
			goto tr5
		case 43:
			if (m.data)[(m.p)] == 45 {
				goto tr12
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr12
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr12
				}
			default:
				goto tr12
			}
			goto tr51
		case 44:
			switch (m.data)[(m.p)] {
			case 85:
				goto tr52
			case 117:
				goto tr52
			}
			goto tr0
		case 45:
			switch (m.data)[(m.p)] {
			case 82:
				goto tr53
			case 114:
				goto tr53
			}
			goto tr0
		case 46:
			switch (m.data)[(m.p)] {
			case 78:
				goto tr54
			case 110:
				goto tr54
			}
			goto tr0
		case 47:
			if (m.data)[(m.p)] == 58 {
				goto tr55
			}
			goto tr0
		case 48:
			if (m.data)[(m.p)] == 105 {
				goto tr57
			}
			goto tr56
		case 49:
			if (m.data)[(m.p)] == 101 {
				goto tr58
			}
			goto tr56
		case 50:
			if (m.data)[(m.p)] == 116 {
				goto tr59
			}
			goto tr56
		case 51:
			if (m.data)[(m.p)] == 102 {
				goto tr60
			}
			goto tr56
		case 52:
			if (m.data)[(m.p)] == 58 {
				goto tr61
			}
			goto tr56
		case 53:
			if (m.data)[(m.p)] == 112 {
				goto tr62
			}
			goto tr56
		case 54:
			if (m.data)[(m.p)] == 97 {
				goto tr63
			}
			goto tr56
		case 55:
			if (m.data)[(m.p)] == 114 {
				goto tr64
			}
			goto tr56
		case 56:
			if (m.data)[(m.p)] == 97 {
				goto tr65
			}
			goto tr56
		case 57:
			if (m.data)[(m.p)] == 109 {
				goto tr66
			}
			goto tr56
		case 58:
			if (m.data)[(m.p)] == 115 {
				goto tr67
			}
			goto tr56
		case 59:
			if (m.data)[(m.p)] == 58 {
				goto tr68
			}
			goto tr56
		case 60:
			if (m.data)[(m.p)] == 115 {
				goto tr69
			}
			goto tr56
		case 61:
			if (m.data)[(m.p)] == 99 {
				goto tr70
			}
			goto tr56
		case 62:
			if (m.data)[(m.p)] == 105 {
				goto tr71
			}
			goto tr56
		case 63:
			if (m.data)[(m.p)] == 109 {
				goto tr72
			}
			goto tr56
		case 64:
			if (m.data)[(m.p)] == 58 {
				goto tr73
			}
			goto tr56
		case 65:
			switch (m.data)[(m.p)] {
			case 97:
				goto tr75
			case 112:
				goto tr76
			case 115:
				goto tr77
			}
			goto tr74
		case 66:
			if (m.data)[(m.p)] == 112 {
				goto tr78
			}
			goto tr74
		case 67:
			if (m.data)[(m.p)] == 105 {
				goto tr79
			}
			goto tr74
		case 68:
			if (m.data)[(m.p)] == 58 {
				goto tr80
			}
			goto tr74
		case 69:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr82
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr82
				}
			default:
				goto tr82
			}
			goto tr81
		case 174:
			if (m.data)[(m.p)] == 58 {
				goto tr215
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr214
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr214
				}
			default:
				goto tr214
			}
			goto tr81
		case 70:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr84
			case 36:
				goto tr84
			case 37:
				goto tr85
			case 61:
				goto tr84
			case 95:
				goto tr84
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr84
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr84
					}
				case (m.data)[(m.p)] >= 64:
					goto tr84
				}
			default:
				goto tr84
			}
			goto tr83
		case 175:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr216
			case 36:
				goto tr216
			case 37:
				goto tr217
			case 61:
				goto tr216
			case 95:
				goto tr216
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr216
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr216
					}
				case (m.data)[(m.p)] >= 64:
					goto tr216
				}
			default:
				goto tr216
			}
			goto tr83
		case 71:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr87
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr87
				}
			default:
				goto tr88
			}
			goto tr86
		case 72:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr89
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr89
				}
			default:
				goto tr90
			}
			goto tr86
		case 176:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr216
			case 36:
				goto tr216
			case 37:
				goto tr217
			case 61:
				goto tr216
			case 95:
				goto tr216
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 39 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr216
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr216
					}
				case (m.data)[(m.p)] >= 64:
					goto tr216
				}
			default:
				goto tr216
			}
			goto tr86
		case 73:
			if (m.data)[(m.p)] == 97 {
				goto tr91
			}
			goto tr74
		case 74:
			if (m.data)[(m.p)] == 114 {
				goto tr92
			}
			goto tr74
		case 75:
			if (m.data)[(m.p)] == 97 {
				goto tr93
			}
			goto tr74
		case 76:
			if (m.data)[(m.p)] == 109 {
				goto tr79
			}
			goto tr74
		case 77:
			if (m.data)[(m.p)] == 99 {
				goto tr94
			}
			goto tr74
		case 78:
			if (m.data)[(m.p)] == 104 {
				goto tr95
			}
			goto tr74
		case 79:
			if (m.data)[(m.p)] == 101 {
				goto tr96
			}
			goto tr74
		case 80:
			if (m.data)[(m.p)] == 109 {
				goto tr97
			}
			goto tr74
		case 81:
			if (m.data)[(m.p)] == 97 {
				goto tr98
			}
			goto tr74
		case 82:
			if (m.data)[(m.p)] == 115 {
				goto tr79
			}
			goto tr74
		case 83:
			switch (m.data)[(m.p)] {
			case 85:
				goto tr99
			case 117:
				goto tr99
			}
			goto tr0
		case 84:
			switch (m.data)[(m.p)] {
			case 82:
				goto tr100
			case 114:
				goto tr100
			}
			goto tr0
		case 85:
			switch (m.data)[(m.p)] {
			case 78:
				goto tr101
			case 110:
				goto tr101
			}
			goto tr0
		case 86:
			if (m.data)[(m.p)] == 58 {
				goto tr102
			}
			goto tr0
		case 87:
			switch (m.data)[(m.p)] {
			case 85:
				goto tr105
			case 117:
				goto tr105
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr104
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr104
				}
			default:
				goto tr104
			}
			goto tr103
		case 88:
			if (m.data)[(m.p)] == 45 {
				goto tr107
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr108
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr108
				}
			default:
				goto tr108
			}
			goto tr106
		case 89:
			if (m.data)[(m.p)] == 45 {
				goto tr109
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr110
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr110
				}
			default:
				goto tr110
			}
			goto tr106
		case 90:
			if (m.data)[(m.p)] == 45 {
				goto tr111
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr112
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr112
				}
			default:
				goto tr112
			}
			goto tr106
		case 91:
			if (m.data)[(m.p)] == 45 {
				goto tr113
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr114
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr114
				}
			default:
				goto tr114
			}
			goto tr106
		case 92:
			if (m.data)[(m.p)] == 45 {
				goto tr115
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr116
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr116
				}
			default:
				goto tr116
			}
			goto tr106
		case 93:
			if (m.data)[(m.p)] == 45 {
				goto tr117
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr118
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr118
				}
			default:
				goto tr118
			}
			goto tr106
		case 94:
			if (m.data)[(m.p)] == 45 {
				goto tr119
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr120
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr120
				}
			default:
				goto tr120
			}
			goto tr106
		case 95:
			if (m.data)[(m.p)] == 45 {
				goto tr121
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr122
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr122
				}
			default:
				goto tr122
			}
			goto tr106
		case 96:
			if (m.data)[(m.p)] == 45 {
				goto tr123
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr124
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr124
				}
			default:
				goto tr124
			}
			goto tr106
		case 97:
			if (m.data)[(m.p)] == 45 {
				goto tr125
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr126
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr126
				}
			default:
				goto tr126
			}
			goto tr106
		case 98:
			if (m.data)[(m.p)] == 45 {
				goto tr127
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr128
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr128
				}
			default:
				goto tr128
			}
			goto tr106
		case 99:
			if (m.data)[(m.p)] == 45 {
				goto tr129
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr130
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr130
				}
			default:
				goto tr130
			}
			goto tr106
		case 100:
			if (m.data)[(m.p)] == 45 {
				goto tr131
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr132
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr132
				}
			default:
				goto tr132
			}
			goto tr106
		case 101:
			if (m.data)[(m.p)] == 45 {
				goto tr133
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr134
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr134
				}
			default:
				goto tr134
			}
			goto tr106
		case 102:
			if (m.data)[(m.p)] == 45 {
				goto tr135
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr136
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr136
				}
			default:
				goto tr136
			}
			goto tr106
		case 103:
			if (m.data)[(m.p)] == 45 {
				goto tr137
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr138
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr138
				}
			default:
				goto tr138
			}
			goto tr106
		case 104:
			if (m.data)[(m.p)] == 45 {
				goto tr139
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr140
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr140
				}
			default:
				goto tr140
			}
			goto tr106
		case 105:
			if (m.data)[(m.p)] == 45 {
				goto tr141
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr142
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr142
				}
			default:
				goto tr142
			}
			goto tr106
		case 106:
			if (m.data)[(m.p)] == 45 {
				goto tr143
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr144
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr144
				}
			default:
				goto tr144
			}
			goto tr106
		case 107:
			if (m.data)[(m.p)] == 45 {
				goto tr145
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr146
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr146
				}
			default:
				goto tr146
			}
			goto tr106
		case 108:
			if (m.data)[(m.p)] == 45 {
				goto tr147
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr148
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr148
				}
			default:
				goto tr148
			}
			goto tr106
		case 109:
			if (m.data)[(m.p)] == 45 {
				goto tr149
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr150
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr150
				}
			default:
				goto tr150
			}
			goto tr106
		case 110:
			if (m.data)[(m.p)] == 45 {
				goto tr151
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr152
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr152
				}
			default:
				goto tr152
			}
			goto tr106
		case 111:
			if (m.data)[(m.p)] == 45 {
				goto tr153
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr154
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr154
				}
			default:
				goto tr154
			}
			goto tr106
		case 112:
			if (m.data)[(m.p)] == 45 {
				goto tr155
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr156
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr156
				}
			default:
				goto tr156
			}
			goto tr106
		case 113:
			if (m.data)[(m.p)] == 45 {
				goto tr157
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr158
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr158
				}
			default:
				goto tr158
			}
			goto tr106
		case 114:
			if (m.data)[(m.p)] == 45 {
				goto tr159
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr160
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr160
				}
			default:
				goto tr160
			}
			goto tr106
		case 115:
			if (m.data)[(m.p)] == 45 {
				goto tr161
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr162
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr162
				}
			default:
				goto tr162
			}
			goto tr106
		case 116:
			if (m.data)[(m.p)] == 45 {
				goto tr163
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr164
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr164
				}
			default:
				goto tr164
			}
			goto tr106
		case 117:
			if (m.data)[(m.p)] == 45 {
				goto tr165
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr166
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr166
				}
			default:
				goto tr166
			}
			goto tr106
		case 118:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr167
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr167
				}
			default:
				goto tr167
			}
			goto tr106
		case 119:
			if (m.data)[(m.p)] == 58 {
				goto tr168
			}
			goto tr106
		case 120:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr170
			case 37:
				goto tr171
			case 61:
				goto tr170
			case 95:
				goto tr170
			case 126:
				goto tr170
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr170
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr170
					}
				case (m.data)[(m.p)] >= 64:
					goto tr170
				}
			default:
				goto tr170
			}
			goto tr169
		case 177:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr218
			case 35:
				goto tr219
			case 37:
				goto tr220
			case 61:
				goto tr218
			case 63:
				goto tr221
			case 95:
				goto tr218
			case 126:
				goto tr218
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr218
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr218
				}
			default:
				goto tr218
			}
			goto tr169
		case 178:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr222
			case 37:
				goto tr223
			case 61:
				goto tr222
			case 95:
				goto tr222
			case 126:
				goto tr222
			}
			switch {
			case (m.data)[(m.p)] < 63:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr222
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr222
				}
			default:
				goto tr222
			}
			goto tr183
		case 179:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr224
			case 37:
				goto tr225
			case 61:
				goto tr224
			case 95:
				goto tr224
			case 126:
				goto tr224
			}
			switch {
			case (m.data)[(m.p)] < 63:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr224
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr224
				}
			default:
				goto tr224
			}
			goto tr183
		case 121:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr173
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr173
				}
			default:
				goto tr174
			}
			goto tr172
		case 122:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr175
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr175
				}
			default:
				goto tr176
			}
			goto tr172
		case 180:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr224
			case 37:
				goto tr225
			case 61:
				goto tr224
			case 95:
				goto tr224
			case 126:
				goto tr224
			}
			switch {
			case (m.data)[(m.p)] < 63:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr224
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr224
				}
			default:
				goto tr224
			}
			goto tr172
		case 123:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr178
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr178
				}
			default:
				goto tr179
			}
			goto tr177
		case 124:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr180
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr180
				}
			default:
				goto tr181
			}
			goto tr177
		case 181:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr218
			case 35:
				goto tr219
			case 37:
				goto tr220
			case 61:
				goto tr218
			case 63:
				goto tr221
			case 95:
				goto tr218
			case 126:
				goto tr218
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr218
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr218
				}
			default:
				goto tr218
			}
			goto tr177
		case 125:
			switch (m.data)[(m.p)] {
			case 43:
				goto tr182
			case 61:
				goto tr184
			}
			goto tr183
		case 126:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr186
			case 37:
				goto tr187
			case 61:
				goto tr186
			case 63:
				goto tr188
			case 95:
				goto tr186
			case 126:
				goto tr186
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr186
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr186
					}
				case (m.data)[(m.p)] >= 64:
					goto tr186
				}
			default:
				goto tr186
			}
			goto tr185
		case 182:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr226
			case 35:
				goto tr227
			case 37:
				goto tr228
			case 61:
				goto tr226
			case 63:
				goto tr229
			case 95:
				goto tr226
			case 126:
				goto tr226
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr226
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr226
				}
			default:
				goto tr226
			}
			goto tr185
		case 127:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr190
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr190
				}
			default:
				goto tr191
			}
			goto tr189
		case 128:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr192
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr192
				}
			default:
				goto tr193
			}
			goto tr189
		case 183:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr226
			case 35:
				goto tr227
			case 37:
				goto tr228
			case 61:
				goto tr226
			case 63:
				goto tr229
			case 95:
				goto tr226
			case 126:
				goto tr226
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr226
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr226
				}
			default:
				goto tr226
			}
			goto tr189
		case 184:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr226
			case 35:
				goto tr227
			case 37:
				goto tr228
			case 43:
				goto tr230
			case 61:
				goto tr231
			case 63:
				goto tr229
			case 95:
				goto tr226
			case 126:
				goto tr226
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr226
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr226
				}
			default:
				goto tr226
			}
			goto tr185
		case 185:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr232
			case 35:
				goto tr233
			case 37:
				goto tr234
			case 47:
				goto tr226
			case 61:
				goto tr232
			case 63:
				goto tr235
			case 95:
				goto tr232
			case 126:
				goto tr232
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr232
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr232
				}
			default:
				goto tr232
			}
			goto tr185
		case 186:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr204
			case 35:
				goto tr227
			case 37:
				goto tr237
			case 47:
				goto tr226
			case 61:
				goto tr204
			case 63:
				goto tr229
			case 95:
				goto tr204
			case 126:
				goto tr204
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr204
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr204
				}
			default:
				goto tr204
			}
			goto tr236
		case 187:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr238
			case 35:
				goto tr239
			case 37:
				goto tr240
			case 61:
				goto tr238
			case 63:
				goto tr241
			case 95:
				goto tr238
			case 126:
				goto tr238
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr238
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr238
				}
			default:
				goto tr238
			}
			goto tr203
		case 129:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr195
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr195
				}
			default:
				goto tr196
			}
			goto tr194
		case 130:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr197
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr197
				}
			default:
				goto tr198
			}
			goto tr194
		case 188:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr238
			case 35:
				goto tr239
			case 37:
				goto tr240
			case 61:
				goto tr238
			case 63:
				goto tr241
			case 95:
				goto tr238
			case 126:
				goto tr238
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr238
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr238
				}
			default:
				goto tr238
			}
			goto tr194
		case 189:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr238
			case 35:
				goto tr239
			case 37:
				goto tr240
			case 61:
				goto tr242
			case 63:
				goto tr241
			case 95:
				goto tr238
			case 126:
				goto tr238
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr238
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr238
				}
			default:
				goto tr238
			}
			goto tr203
		case 190:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr243
			case 35:
				goto tr244
			case 37:
				goto tr245
			case 47:
				goto tr238
			case 61:
				goto tr243
			case 63:
				goto tr246
			case 95:
				goto tr243
			case 126:
				goto tr243
			}
			switch {
			case (m.data)[(m.p)] < 64:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 59 {
					goto tr243
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr243
				}
			default:
				goto tr243
			}
			goto tr203
		case 131:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr200
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr200
				}
			default:
				goto tr201
			}
			goto tr199
		case 132:
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr197
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr197
				}
			default:
				goto tr198
			}
			goto tr199
		case 133:
			if (m.data)[(m.p)] == 43 {
				goto tr202
			}
			goto tr185
		case 191:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr232
			case 35:
				goto tr233
			case 37:
				goto tr234
			case 61:
				goto tr232
			case 63:
				goto tr247
			case 95:
				goto tr232
			case 126:
				goto tr232
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr232
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr232
					}
				case (m.data)[(m.p)] >= 64:
					goto tr232
				}
			default:
				goto tr232
			}
			goto tr185
		case 134:
			switch (m.data)[(m.p)] {
			case 43:
				goto tr202
			case 61:
				goto tr184
			}
			goto tr185
		case 135:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr204
			case 37:
				goto tr205
			case 61:
				goto tr204
			case 63:
				goto tr206
			case 95:
				goto tr204
			case 126:
				goto tr204
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr204
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr204
					}
				case (m.data)[(m.p)] >= 64:
					goto tr204
				}
			default:
				goto tr204
			}
			goto tr203
		case 136:
			if (m.data)[(m.p)] == 61 {
				goto tr207
			}
			goto tr203
		case 192:
			switch (m.data)[(m.p)] {
			case 33:
				goto tr243
			case 35:
				goto tr244
			case 37:
				goto tr245
			case 61:
				goto tr243
			case 63:
				goto tr248
			case 95:
				goto tr243
			case 126:
				goto tr243
			}
			switch {
			case (m.data)[(m.p)] < 48:
				if 36 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 46 {
					goto tr243
				}
			case (m.data)[(m.p)] > 59:
				switch {
				case (m.data)[(m.p)] > 90:
					if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
						goto tr243
					}
				case (m.data)[(m.p)] >= 64:
					goto tr243
				}
			default:
				goto tr243
			}
			goto tr203
		case 137:
			if (m.data)[(m.p)] == 58 {
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr167
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr167
				}
			default:
				goto tr167
			}
			goto tr106
		case 138:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr165
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr166
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr166
				}
			default:
				goto tr166
			}
			goto tr106
		case 139:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr163
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr164
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr164
				}
			default:
				goto tr164
			}
			goto tr106
		case 140:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr161
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr162
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr162
				}
			default:
				goto tr162
			}
			goto tr106
		case 141:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr159
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr160
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr160
				}
			default:
				goto tr160
			}
			goto tr106
		case 142:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr157
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr158
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr158
				}
			default:
				goto tr158
			}
			goto tr106
		case 143:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr155
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr156
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr156
				}
			default:
				goto tr156
			}
			goto tr106
		case 144:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr153
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr154
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr154
				}
			default:
				goto tr154
			}
			goto tr106
		case 145:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr151
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr152
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr152
				}
			default:
				goto tr152
			}
			goto tr106
		case 146:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr149
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr150
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr150
				}
			default:
				goto tr150
			}
			goto tr106
		case 147:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr147
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr148
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr148
				}
			default:
				goto tr148
			}
			goto tr106
		case 148:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr145
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr146
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr146
				}
			default:
				goto tr146
			}
			goto tr106
		case 149:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr143
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr144
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr144
				}
			default:
				goto tr144
			}
			goto tr106
		case 150:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr141
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr142
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr142
				}
			default:
				goto tr142
			}
			goto tr106
		case 151:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr139
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr140
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr140
				}
			default:
				goto tr140
			}
			goto tr106
		case 152:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr137
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr138
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr138
				}
			default:
				goto tr138
			}
			goto tr106
		case 153:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr135
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr136
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr136
				}
			default:
				goto tr136
			}
			goto tr106
		case 154:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr133
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr134
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr134
				}
			default:
				goto tr134
			}
			goto tr106
		case 155:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr131
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr132
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr132
				}
			default:
				goto tr132
			}
			goto tr106
		case 156:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr129
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr130
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr130
				}
			default:
				goto tr130
			}
			goto tr106
		case 157:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr127
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr128
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr128
				}
			default:
				goto tr128
			}
			goto tr106
		case 158:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr125
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr126
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr126
				}
			default:
				goto tr126
			}
			goto tr106
		case 159:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr123
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr124
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr124
				}
			default:
				goto tr124
			}
			goto tr106
		case 160:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr121
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr122
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr122
				}
			default:
				goto tr122
			}
			goto tr106
		case 161:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr119
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr120
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr120
				}
			default:
				goto tr120
			}
			goto tr106
		case 162:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr117
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr118
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr118
				}
			default:
				goto tr118
			}
			goto tr106
		case 163:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr115
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr116
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr116
				}
			default:
				goto tr116
			}
			goto tr106
		case 164:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr113
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr114
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr114
				}
			default:
				goto tr114
			}
			goto tr106
		case 165:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr111
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr112
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr112
				}
			default:
				goto tr112
			}
			goto tr106
		case 166:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr109
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr110
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr110
				}
			default:
				goto tr110
			}
			goto tr106
		case 167:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr107
			case 82:
				goto tr208
			case 114:
				goto tr208
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr108
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr108
				}
			default:
				goto tr108
			}
			goto tr103
		case 168:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr109
			case 58:
				goto tr168
			case 78:
				goto tr209
			case 110:
				goto tr209
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr110
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr110
				}
			default:
				goto tr110
			}
			goto tr103
		case 169:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr210
			case 58:
				goto tr168
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr112
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr112
				}
			default:
				goto tr112
			}
			goto tr106
		case 170:
			switch (m.data)[(m.p)] {
			case 45:
				goto tr113
			case 48:
				goto tr211
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 49 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr114
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr211
				}
			default:
				goto tr211
			}
			goto tr106
		case 171:
			if (m.data)[(m.p)] == 45 {
				goto tr115
			}
			switch {
			case (m.data)[(m.p)] < 65:
				if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
					goto tr116
				}
			case (m.data)[(m.p)] > 90:
				if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
					goto tr116
				}
			default:
				goto tr116
			}
			goto tr106
		case 193:
			switch (m.data)[(m.p)] {
			case 10:
				goto tr183
			case 13:
				goto tr183
			}
			goto tr249
		}

	tr183:
		m.cs = 0
		goto _again
	tr0:
		m.cs = 0
		goto f0
	tr5:
		m.cs = 0
		goto f3
	tr8:
		m.cs = 0
		goto f5
	tr41:
		m.cs = 0
		goto f7
	tr44:
		m.cs = 0
		goto f8
	tr51:
		m.cs = 0
		goto f10
	tr56:
		m.cs = 0
		goto f11
	tr74:
		m.cs = 0
		goto f13
	tr81:
		m.cs = 0
		goto f15
	tr83:
		m.cs = 0
		goto f17
	tr86:
		m.cs = 0
		goto f19
	tr103:
		m.cs = 0
		goto f20
	tr106:
		m.cs = 0
		goto f21
	tr169:
		m.cs = 0
		goto f22
	tr172:
		m.cs = 0
		goto f23
	tr177:
		m.cs = 0
		goto f24
	tr185:
		m.cs = 0
		goto f25
	tr189:
		m.cs = 0
		goto f27
	tr194:
		m.cs = 0
		goto f28
	tr199:
		m.cs = 0
		goto f29
	tr203:
		m.cs = 0
		goto f30
	tr236:
		m.cs = 0
		goto f46
	tr1:
		m.cs = 2
		goto f1
	tr2:
		m.cs = 3
		goto _again
	tr3:
		m.cs = 4
		goto _again
	tr4:
		m.cs = 5
		goto f2
	tr6:
		m.cs = 6
		goto f4
	tr9:
		m.cs = 7
		goto _again
	tr11:
		m.cs = 8
		goto _again
	tr12:
		m.cs = 9
		goto _again
	tr13:
		m.cs = 10
		goto _again
	tr14:
		m.cs = 11
		goto _again
	tr15:
		m.cs = 12
		goto _again
	tr16:
		m.cs = 13
		goto _again
	tr17:
		m.cs = 14
		goto _again
	tr18:
		m.cs = 15
		goto _again
	tr19:
		m.cs = 16
		goto _again
	tr20:
		m.cs = 17
		goto _again
	tr21:
		m.cs = 18
		goto _again
	tr22:
		m.cs = 19
		goto _again
	tr23:
		m.cs = 20
		goto _again
	tr24:
		m.cs = 21
		goto _again
	tr25:
		m.cs = 22
		goto _again
	tr26:
		m.cs = 23
		goto _again
	tr27:
		m.cs = 24
		goto _again
	tr28:
		m.cs = 25
		goto _again
	tr29:
		m.cs = 26
		goto _again
	tr30:
		m.cs = 27
		goto _again
	tr31:
		m.cs = 28
		goto _again
	tr32:
		m.cs = 29
		goto _again
	tr33:
		m.cs = 30
		goto _again
	tr34:
		m.cs = 31
		goto _again
	tr35:
		m.cs = 32
		goto _again
	tr36:
		m.cs = 33
		goto _again
	tr37:
		m.cs = 34
		goto _again
	tr38:
		m.cs = 35
		goto _again
	tr39:
		m.cs = 36
		goto _again
	tr40:
		m.cs = 37
		goto _again
	tr10:
		m.cs = 38
		goto f6
	tr213:
		m.cs = 39
		goto _again
	tr43:
		m.cs = 39
		goto f4
	tr45:
		m.cs = 40
		goto _again
	tr46:
		m.cs = 40
		goto f9
	tr7:
		m.cs = 41
		goto f1
	tr49:
		m.cs = 42
		goto _again
	tr50:
		m.cs = 43
		goto _again
	tr52:
		m.cs = 45
		goto f1
	tr53:
		m.cs = 46
		goto _again
	tr54:
		m.cs = 47
		goto _again
	tr55:
		m.cs = 48
		goto f2
	tr57:
		m.cs = 49
		goto f4
	tr58:
		m.cs = 50
		goto _again
	tr59:
		m.cs = 51
		goto _again
	tr60:
		m.cs = 52
		goto _again
	tr61:
		m.cs = 53
		goto _again
	tr62:
		m.cs = 54
		goto _again
	tr63:
		m.cs = 55
		goto _again
	tr64:
		m.cs = 56
		goto _again
	tr65:
		m.cs = 57
		goto _again
	tr66:
		m.cs = 58
		goto _again
	tr67:
		m.cs = 59
		goto _again
	tr68:
		m.cs = 60
		goto _again
	tr69:
		m.cs = 61
		goto _again
	tr70:
		m.cs = 62
		goto _again
	tr71:
		m.cs = 63
		goto _again
	tr72:
		m.cs = 64
		goto _again
	tr73:
		m.cs = 65
		goto f12
	tr75:
		m.cs = 66
		goto f4
	tr78:
		m.cs = 67
		goto _again
	tr79:
		m.cs = 68
		goto _again
	tr80:
		m.cs = 69
		goto f14
	tr215:
		m.cs = 70
		goto f35
	tr217:
		m.cs = 71
		goto _again
	tr85:
		m.cs = 71
		goto f18
	tr87:
		m.cs = 72
		goto _again
	tr88:
		m.cs = 72
		goto f9
	tr76:
		m.cs = 73
		goto f4
	tr91:
		m.cs = 74
		goto _again
	tr92:
		m.cs = 75
		goto _again
	tr93:
		m.cs = 76
		goto _again
	tr77:
		m.cs = 77
		goto f4
	tr94:
		m.cs = 78
		goto _again
	tr95:
		m.cs = 79
		goto _again
	tr96:
		m.cs = 80
		goto _again
	tr97:
		m.cs = 81
		goto _again
	tr98:
		m.cs = 82
		goto _again
	tr99:
		m.cs = 84
		goto f1
	tr100:
		m.cs = 85
		goto _again
	tr101:
		m.cs = 86
		goto _again
	tr102:
		m.cs = 87
		goto f2
	tr104:
		m.cs = 88
		goto f4
	tr107:
		m.cs = 89
		goto _again
	tr109:
		m.cs = 90
		goto _again
	tr111:
		m.cs = 91
		goto _again
	tr113:
		m.cs = 92
		goto _again
	tr115:
		m.cs = 93
		goto _again
	tr117:
		m.cs = 94
		goto _again
	tr119:
		m.cs = 95
		goto _again
	tr121:
		m.cs = 96
		goto _again
	tr123:
		m.cs = 97
		goto _again
	tr125:
		m.cs = 98
		goto _again
	tr127:
		m.cs = 99
		goto _again
	tr129:
		m.cs = 100
		goto _again
	tr131:
		m.cs = 101
		goto _again
	tr133:
		m.cs = 102
		goto _again
	tr135:
		m.cs = 103
		goto _again
	tr137:
		m.cs = 104
		goto _again
	tr139:
		m.cs = 105
		goto _again
	tr141:
		m.cs = 106
		goto _again
	tr143:
		m.cs = 107
		goto _again
	tr145:
		m.cs = 108
		goto _again
	tr147:
		m.cs = 109
		goto _again
	tr149:
		m.cs = 110
		goto _again
	tr151:
		m.cs = 111
		goto _again
	tr153:
		m.cs = 112
		goto _again
	tr155:
		m.cs = 113
		goto _again
	tr157:
		m.cs = 114
		goto _again
	tr159:
		m.cs = 115
		goto _again
	tr161:
		m.cs = 116
		goto _again
	tr163:
		m.cs = 117
		goto _again
	tr165:
		m.cs = 118
		goto _again
	tr167:
		m.cs = 119
		goto _again
	tr168:
		m.cs = 120
		goto f6
	tr225:
		m.cs = 121
		goto _again
	tr223:
		m.cs = 121
		goto f4
	tr173:
		m.cs = 122
		goto _again
	tr174:
		m.cs = 122
		goto f9
	tr220:
		m.cs = 123
		goto _again
	tr171:
		m.cs = 123
		goto f4
	tr178:
		m.cs = 124
		goto _again
	tr179:
		m.cs = 124
		goto f9
	tr221:
		m.cs = 125
		goto f38
	tr182:
		m.cs = 126
		goto _again
	tr228:
		m.cs = 127
		goto _again
	tr187:
		m.cs = 127
		goto f26
	tr234:
		m.cs = 127
		goto f44
	tr190:
		m.cs = 128
		goto _again
	tr191:
		m.cs = 128
		goto f9
	tr240:
		m.cs = 129
		goto _again
	tr205:
		m.cs = 129
		goto f31
	tr245:
		m.cs = 129
		goto f50
	tr195:
		m.cs = 130
		goto _again
	tr196:
		m.cs = 130
		goto f9
	tr237:
		m.cs = 131
		goto f31
	tr200:
		m.cs = 132
		goto _again
	tr201:
		m.cs = 132
		goto f9
	tr188:
		m.cs = 133
		goto f26
	tr247:
		m.cs = 134
		goto f45
	tr184:
		m.cs = 135
		goto _again
	tr206:
		m.cs = 136
		goto f31
	tr248:
		m.cs = 136
		goto f50
	tr166:
		m.cs = 137
		goto _again
	tr164:
		m.cs = 138
		goto _again
	tr162:
		m.cs = 139
		goto _again
	tr160:
		m.cs = 140
		goto _again
	tr158:
		m.cs = 141
		goto _again
	tr156:
		m.cs = 142
		goto _again
	tr154:
		m.cs = 143
		goto _again
	tr152:
		m.cs = 144
		goto _again
	tr150:
		m.cs = 145
		goto _again
	tr148:
		m.cs = 146
		goto _again
	tr146:
		m.cs = 147
		goto _again
	tr144:
		m.cs = 148
		goto _again
	tr142:
		m.cs = 149
		goto _again
	tr140:
		m.cs = 150
		goto _again
	tr138:
		m.cs = 151
		goto _again
	tr136:
		m.cs = 152
		goto _again
	tr134:
		m.cs = 153
		goto _again
	tr132:
		m.cs = 154
		goto _again
	tr130:
		m.cs = 155
		goto _again
	tr128:
		m.cs = 156
		goto _again
	tr126:
		m.cs = 157
		goto _again
	tr124:
		m.cs = 158
		goto _again
	tr122:
		m.cs = 159
		goto _again
	tr120:
		m.cs = 160
		goto _again
	tr118:
		m.cs = 161
		goto _again
	tr116:
		m.cs = 162
		goto _again
	tr114:
		m.cs = 163
		goto _again
	tr112:
		m.cs = 164
		goto _again
	tr110:
		m.cs = 165
		goto _again
	tr108:
		m.cs = 166
		goto _again
	tr105:
		m.cs = 167
		goto f1
	tr208:
		m.cs = 168
		goto _again
	tr209:
		m.cs = 169
		goto _again
	tr210:
		m.cs = 170
		goto f2
	tr211:
		m.cs = 171
		goto _again
	tr212:
		m.cs = 172
		goto _again
	tr42:
		m.cs = 172
		goto f4
	tr47:
		m.cs = 173
		goto _again
	tr48:
		m.cs = 173
		goto f9
	tr214:
		m.cs = 174
		goto _again
	tr82:
		m.cs = 174
		goto f16
	tr216:
		m.cs = 175
		goto _again
	tr84:
		m.cs = 175
		goto f18
	tr89:
		m.cs = 176
		goto _again
	tr90:
		m.cs = 176
		goto f9
	tr218:
		m.cs = 177
		goto _again
	tr170:
		m.cs = 177
		goto f4
	tr219:
		m.cs = 178
		goto f38
	tr227:
		m.cs = 178
		goto f42
	tr233:
		m.cs = 178
		goto f45
	tr239:
		m.cs = 178
		goto f48
	tr244:
		m.cs = 178
		goto f51
	tr224:
		m.cs = 179
		goto _again
	tr222:
		m.cs = 179
		goto f4
	tr175:
		m.cs = 180
		goto _again
	tr176:
		m.cs = 180
		goto f9
	tr180:
		m.cs = 181
		goto _again
	tr181:
		m.cs = 181
		goto f9
	tr226:
		m.cs = 182
		goto _again
	tr186:
		m.cs = 182
		goto f26
	tr232:
		m.cs = 182
		goto f44
	tr192:
		m.cs = 183
		goto _again
	tr193:
		m.cs = 183
		goto f9
	tr229:
		m.cs = 184
		goto f42
	tr235:
		m.cs = 184
		goto f45
	tr230:
		m.cs = 185
		goto _again
	tr231:
		m.cs = 186
		goto _again
	tr238:
		m.cs = 187
		goto _again
	tr204:
		m.cs = 187
		goto f31
	tr243:
		m.cs = 187
		goto f50
	tr197:
		m.cs = 188
		goto _again
	tr198:
		m.cs = 188
		goto f9
	tr241:
		m.cs = 189
		goto _again
	tr246:
		m.cs = 189
		goto f50
	tr242:
		m.cs = 190
		goto _again
	tr202:
		m.cs = 191
		goto _again
	tr207:
		m.cs = 192
		goto _again
	tr249:
		m.cs = 193
		goto _again

	f4:

		m.pb = m.p

		goto _again
	f9:

		// List of positions in the buffer to later lowercase
		output.tolower = append(output.tolower, m.p-m.pb)

		goto _again
	f2:

		output.prefix = string(m.text())

		goto _again
	f6:

		output.ID = string(m.text())

		goto _again
	f38:

		output.SS = string(m.text())
		// Iterate upper letters lowering them
		for _, i := range output.tolower {
			m.data[m.pb+i] = m.data[m.pb+i] + 32
		}
		output.norm = string(m.text())
		// Revert the buffer to the original
		for _, i := range output.tolower {
			m.data[m.pb+i] = m.data[m.pb+i] - 32
		}

		goto _again
	f0:

		m.err = fmt.Errorf(errPrefix, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f5:

		m.err = fmt.Errorf(errIdentifier, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f7:

		m.err = fmt.Errorf(errSpecificString, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f23:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		goto _again
	f11:

		m.err = fmt.Errorf(errSCIMNamespace, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f13:

		m.err = fmt.Errorf(errSCIMType, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f15:

		m.err = fmt.Errorf(errSCIMName, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f17:

		if m.p == m.pe {
			m.err = fmt.Errorf(errSCIMOtherIncomplete, m.p-1)
		} else {
			m.err = fmt.Errorf(errSCIMOther, m.p)
		}
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f14:

		output.scim.Type = scimschema.TypeFromString(string(m.text()))

		goto _again
	f16:

		output.scim.pos = m.p

		goto _again
	f35:

		output.scim.Name = string(m.data[output.scim.pos:m.p])

		goto _again
	f18:

		output.scim.pos = m.p

		goto _again
	f22:

		m.err = fmt.Errorf(err8141SpecificString, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f21:

		m.err = fmt.Errorf(err8141Identifier, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f42:

		output.rComponent = string(m.text())

		goto _again
	f48:

		output.qComponent = string(m.text())

		goto _again
	f44:

		if output.rStart {
			m.err = fmt.Errorf(err8141RComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.rStart = true

		goto _again
	f50:

		if output.qStart {
			m.err = fmt.Errorf(err8141QComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.qStart = true

		goto _again
	f25:

		m.err = fmt.Errorf(err8141MalformedRComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f30:

		m.err = fmt.Errorf(err8141MalformedQComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f1:

		m.pb = m.p

		if m.parsingMode != RFC8141Only {
			// Throw an error when:
			// - we are entering here matching the the prefix in the namespace identifier part
			// - looking ahead (3 chars) we find a colon
			if pos := m.p + 3; pos < m.pe && m.data[pos] == 58 && output.prefix != "" {
				m.err = fmt.Errorf(errNoUrnWithinID, pos)
				(m.p)--

				m.cs = 193
				goto _again

			}
		}

		goto _again
	f12:

		output.ID = string(m.text())

		output.scim = &SCIM{}

		goto _again
	f3:

		m.err = fmt.Errorf(errIdentifier, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		m.err = fmt.Errorf(errPrefix, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f10:

		m.err = fmt.Errorf(errIdentifier, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		m.err = fmt.Errorf(errNoUrnWithinID, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f8:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		m.err = fmt.Errorf(errSpecificString, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f19:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		if m.p == m.pe {
			m.err = fmt.Errorf(errSCIMOtherIncomplete, m.p-1)
		} else {
			m.err = fmt.Errorf(errSCIMOther, m.p)
		}
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f24:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		m.err = fmt.Errorf(err8141SpecificString, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f27:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		m.err = fmt.Errorf(err8141MalformedRComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f28:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		m.err = fmt.Errorf(err8141MalformedQComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f20:

		m.err = fmt.Errorf(err8141Identifier, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		m.err = fmt.Errorf(errPrefix, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f26:

		if output.rStart {
			m.err = fmt.Errorf(err8141RComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.rStart = true

		m.pb = m.p

		goto _again
	f45:

		if output.rStart {
			m.err = fmt.Errorf(err8141RComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.rStart = true

		output.rComponent = string(m.text())

		goto _again
	f31:

		if output.qStart {
			m.err = fmt.Errorf(err8141QComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.qStart = true

		m.pb = m.p

		goto _again
	f51:

		if output.qStart {
			m.err = fmt.Errorf(err8141QComponentStart, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}
		output.qStart = true

		output.qComponent = string(m.text())

		goto _again
	f46:

		m.err = fmt.Errorf(err8141MalformedRComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		m.err = fmt.Errorf(err8141MalformedQComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again
	f29:

		if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
			m.err = fmt.Errorf(errHex, m.p)
			(m.p)--

			m.cs = 193
			goto _again

		}

		m.err = fmt.Errorf(err8141MalformedRComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		m.err = fmt.Errorf(err8141MalformedQComp, m.p)
		(m.p)--

		m.cs = 193
		goto _again

		goto _again

	_again:
		switch _toStateActions[m.cs] {
		case 33:

			(m.p)--

			m.err = fmt.Errorf(err8141InformalID, m.p)
			m.cs = 193
			goto _again
		}

		if m.cs == 0 {
			goto _out
		}
		if (m.p)++; (m.p) != (m.pe) {
			goto _resume
		}
	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch _eofActions[m.cs] {
			case 1:

				m.err = fmt.Errorf(errPrefix, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 6:

				m.err = fmt.Errorf(errIdentifier, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 8:

				m.err = fmt.Errorf(errSpecificString, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 24:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

			case 12:

				m.err = fmt.Errorf(errSCIMNamespace, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 14:

				m.err = fmt.Errorf(errSCIMType, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 16:

				m.err = fmt.Errorf(errSCIMName, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 18:

				if m.p == m.pe {
					m.err = fmt.Errorf(errSCIMOtherIncomplete, m.p-1)
				} else {
					m.err = fmt.Errorf(errSCIMOther, m.p)
				}
				(m.p)--

				m.cs = 193
				goto _again

			case 23:

				m.err = fmt.Errorf(err8141SpecificString, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 22:

				m.err = fmt.Errorf(err8141Identifier, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 26:

				m.err = fmt.Errorf(err8141MalformedRComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 31:

				m.err = fmt.Errorf(err8141MalformedQComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 34:

				output.SS = string(m.text())
				// Iterate upper letters lowering them
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] + 32
				}
				output.norm = string(m.text())
				// Revert the buffer to the original
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] - 32
				}

				output.kind = RFC2141

			case 38:

				output.SS = string(m.text())
				// Iterate upper letters lowering them
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] + 32
				}
				output.norm = string(m.text())
				// Revert the buffer to the original
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] - 32
				}

				output.kind = RFC8141

			case 4:

				m.err = fmt.Errorf(errIdentifier, m.p)
				(m.p)--

				m.cs = 193
				goto _again

				m.err = fmt.Errorf(errPrefix, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 11:

				m.err = fmt.Errorf(errIdentifier, m.p)
				(m.p)--

				m.cs = 193
				goto _again

				m.err = fmt.Errorf(errNoUrnWithinID, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 9:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				m.err = fmt.Errorf(errSpecificString, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 20:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				if m.p == m.pe {
					m.err = fmt.Errorf(errSCIMOtherIncomplete, m.p-1)
				} else {
					m.err = fmt.Errorf(errSCIMOther, m.p)
				}
				(m.p)--

				m.cs = 193
				goto _again

			case 25:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				m.err = fmt.Errorf(err8141SpecificString, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 28:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				m.err = fmt.Errorf(err8141MalformedRComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 29:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				m.err = fmt.Errorf(err8141MalformedQComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 21:

				m.err = fmt.Errorf(err8141Identifier, m.p)
				(m.p)--

				m.cs = 193
				goto _again

				m.err = fmt.Errorf(errPrefix, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 42:

				output.rComponent = string(m.text())

				output.kind = RFC8141

			case 48:

				output.qComponent = string(m.text())

				output.kind = RFC8141

			case 41:

				output.fComponent = string(m.text())

				output.kind = RFC8141

			case 40:

				m.pb = m.p

				output.fComponent = string(m.text())

				output.kind = RFC8141

			case 30:

				if m.parsingMode == RFC2141Only || m.parsingMode == RFC8141Only {
					m.err = fmt.Errorf(errHex, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}

				m.err = fmt.Errorf(err8141MalformedRComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

				m.err = fmt.Errorf(err8141MalformedQComp, m.p)
				(m.p)--

				m.cs = 193
				goto _again

			case 35:

				output.scim.Name = string(m.data[output.scim.pos:m.p])

				output.SS = string(m.text())
				// Iterate upper letters lowering them
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] + 32
				}
				output.norm = string(m.text())
				// Revert the buffer to the original
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] - 32
				}

				output.kind = RFC7643

			case 37:

				output.scim.Other = string(m.data[output.scim.pos:m.p])

				output.SS = string(m.text())
				// Iterate upper letters lowering them
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] + 32
				}
				output.norm = string(m.text())
				// Revert the buffer to the original
				for _, i := range output.tolower {
					m.data[m.pb+i] = m.data[m.pb+i] - 32
				}

				output.kind = RFC7643

			case 44:

				if output.rStart {
					m.err = fmt.Errorf(err8141RComponentStart, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}
				output.rStart = true

				output.rComponent = string(m.text())

				output.kind = RFC8141

			case 50:

				if output.qStart {
					m.err = fmt.Errorf(err8141QComponentStart, m.p)
					(m.p)--

					m.cs = 193
					goto _again

				}
				output.qStart = true

				output.qComponent = string(m.text())

				output.kind = RFC8141
			}
		}

	_out:
		{
		}
	}

	if m.cs < firstFinal || m.cs == enFail {
		return nil, m.err
	}

	return output, nil
}

func (m *machine) WithParsingMode(x ParsingMode) {
	m.parsingMode = x
	switch m.parsingMode {
	case RFC2141Only:
		m.startParsingAt = enMain
	case RFC8141Only:
		m.startParsingAt = enRfc8141Only
	case RFC7643Only:
		m.startParsingAt = enScimOnly
	}
	m.parsingModeSet = true
}
