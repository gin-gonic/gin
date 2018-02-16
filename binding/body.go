package binding

import "net/url"

// Body type map[string]interface{}
type Body map[string]interface{}

// ToBody transforms object into Body type.
func ToBody(in interface{}) *Body {
	switch in.(type) {
	case Body:
		res := in.(Body)

		return &res
	case url.Values:
		result := make(Body)

		for key, value := range in.(url.Values) {
			if len(value) == 1 {
				result[key] = value[0]
			} else {
				result[key] = value
			}
		}

		return &result
	}

	return new(Body)
}
