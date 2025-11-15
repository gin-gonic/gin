package urn

import (
	"encoding/json"
	"fmt"
	"strings"
)

const errInvalidURN = "invalid URN: %s"

// URN represents an Uniform Resource Name.
//
// The general form represented is:
//
//	urn:<id>:<ss>
//
// Details at https://tools.ietf.org/html/rfc2141.
type URN struct {
	prefix     string // Static prefix. Equal to "urn" when empty.
	ID         string // Namespace identifier (NID)
	SS         string // Namespace specific string (NSS)
	norm       string // Normalized namespace specific string
	kind       Kind
	scim       *SCIM
	rComponent string // RFC8141
	qComponent string // RFC8141
	fComponent string // RFC8141
	rStart     bool   // RFC8141
	qStart     bool   // RFC8141
	tolower    []int
}

// Normalize turns the receiving URN into its norm version.
//
// Which means: lowercase prefix, lowercase namespace identifier, and immutate namespace specific string chars (except <hex> tokens which are lowercased).
func (u *URN) Normalize() *URN {
	return &URN{
		prefix: "urn",
		ID:     strings.ToLower(u.ID),
		SS:     u.norm,
		// rComponent: u.rComponent,
		// qComponent: u.qComponent,
		// fComponent: u.fComponent,
	}
}

// Equal checks the lexical equivalence of the current URN with another one.
func (u *URN) Equal(x *URN) bool {
	if x == nil {
		return false
	}
	nu := u.Normalize()
	nx := x.Normalize()

	return nu.prefix == nx.prefix && nu.ID == nx.ID && nu.SS == nx.SS
}

// String reassembles the URN into a valid URN string.
//
// This requires both ID and SS fields to be non-empty.
// Otherwise it returns an empty string.
//
// Default URN prefix is "urn".
func (u *URN) String() string {
	var res string
	if u.ID != "" && u.SS != "" {
		if u.prefix == "" {
			res += "urn"
		}
		res += u.prefix + ":" + u.ID + ":" + u.SS
		if u.rComponent != "" {
			res += "?+" + u.rComponent
		}
		if u.qComponent != "" {
			res += "?=" + u.qComponent
		}
		if u.fComponent != "" {
			res += "#" + u.fComponent
		}
	}

	return res
}

// Parse is responsible to create an URN instance from a byte array matching the correct URN syntax (RFC 2141).
func Parse(u []byte, options ...Option) (*URN, bool) {
	urn, err := NewMachine(options...).Parse(u)
	if err != nil {
		return nil, false
	}

	return urn, true
}

// MarshalJSON marshals the URN to JSON string form (e.g. `"urn:oid:1.2.3.4"`).
func (u URN) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON unmarshals a URN from JSON string form (e.g. `"urn:oid:1.2.3.4"`).
func (u *URN) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if value, ok := Parse([]byte(str)); !ok {
		return fmt.Errorf(errInvalidURN, str)
	} else {
		*u = *value
	}

	return nil
}

func (u *URN) IsSCIM() bool {
	return u.kind == RFC7643
}

func (u *URN) SCIM() *SCIM {
	if u.kind != RFC7643 {
		return nil
	}

	return u.scim
}

func (u *URN) RFC() Kind {
	return u.kind
}

func (u *URN) FComponent() string {
	return u.fComponent
}

func (u *URN) QComponent() string {
	return u.qComponent
}

func (u *URN) RComponent() string {
	return u.rComponent
}
