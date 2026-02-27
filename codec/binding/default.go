//go:build !gin_bind_encoding

package bindingcodec

import (
	"reflect"
)

// Description summarizes what this binding codec does. This variable helps
// ensure that build tags are configured to be mutually exclusive
const Description = "Gin default binding api"

// TrySetByInterface uses bindUnmarshaler if implemented, otherwise returns false to revert to gin's default binding logic
func (d bindingApi) TrySetByInterface(inputVal string, valueToSet reflect.Value) (isSet bool, err error) {
	switch v := valueToSet.Addr().Interface().(type) {
	case bindUnmarshaler:
		return true, v.UnmarshalParam(inputVal)
	}
	return false, nil
}
