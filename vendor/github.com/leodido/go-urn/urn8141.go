package urn

import (
	"encoding/json"
	"fmt"
)

const errInvalidURN8141 = "invalid URN per RFC 8141: %s"

type URN8141 struct {
	*URN
}

func (u URN8141) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URN8141) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if value, ok := Parse([]byte(str), WithParsingMode(RFC8141Only)); !ok {
		return fmt.Errorf(errInvalidURN8141, str)
	} else {
		*u = URN8141{value}
	}

	return nil
}
