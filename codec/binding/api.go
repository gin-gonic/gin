package bindingcodec

import "reflect"

func init() {
	API = bindingApi{}
}

type bindingApi struct{}

// API the binding codec in use
var API Core

// Core the api for binding codec
type Core interface {
	// TrySetByInterface tries to set valueToSet from inputVal using one of the optional interfaces that gin supports
	//
	// Returns:
	//   - isSet: whether the value was set successfully
	//   - err: any error that occurred during the setting process
	TrySetByInterface(inputVal string, valueToSet reflect.Value) (isSet bool, err error)
}

// bindUnmarshaler duplicates binding.BindUnmarshaler to avoid an import cycle
// This must match binding.BindUnmarshaler exactly to maintain consistent behavior
type bindUnmarshaler interface {
	UnmarshalParam(param string) error
}
