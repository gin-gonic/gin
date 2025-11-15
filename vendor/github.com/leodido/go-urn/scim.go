package urn

import (
	"encoding/json"
	"fmt"

	scimschema "github.com/leodido/go-urn/scim/schema"
)

const errInvalidSCIMURN = "invalid SCIM URN: %s"

type SCIM struct {
	Type  scimschema.Type
	Name  string
	Other string
	pos   int
}

func (s SCIM) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SCIM) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	// Parse as SCIM
	value, ok := Parse([]byte(str), WithParsingMode(RFC7643Only))
	if !ok {
		return fmt.Errorf(errInvalidSCIMURN, str)
	}
	if value.RFC() != RFC7643 {
		return fmt.Errorf(errInvalidSCIMURN, str)
	}
	*s = *value.SCIM()

	return nil
}

func (s *SCIM) String() string {
	ret := fmt.Sprintf("urn:ietf:params:scim:%s:%s", s.Type.String(), s.Name)
	if s.Other != "" {
		ret += fmt.Sprintf(":%s", s.Other)
	}

	return ret
}
