package scimschema

type Type int

const (
	Unsupported Type = iota
	Schemas
	API
	Param
)

func (t Type) String() string {
	switch t {
	case Schemas:
		return "schemas"
	case API:
		return "api"
	case Param:
		return "param"
	}

	return ""
}

func TypeFromString(input string) Type {
	switch input {
	case "schemas":
		return Schemas
	case "api":
		return API
	case "param":
		return Param
	}

	return Unsupported
}
