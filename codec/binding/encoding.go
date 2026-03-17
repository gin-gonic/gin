//go:build gin_bind_encoding

package bindingcodec

import (
	"encoding"
	"log"
	"reflect"
	"time"
)

// Description summarizes what this binding codec does. This variable helps
// ensure that build tags are configured to be mutually exclusive
const Description = "Gin binding api using encoding.TextUnmarshaler"

func init() {
	log.Println("[GIN-debug] gin_bind_encoding flag active. Gin will bind using encoding.TextUnmarshaler by default")
}

// TrySetByInterface checks for bindUnmarshaler first, then falls back to encoding.TextUnmarshaler.
// This allows types which implement TextUnmarshaler (like uuid.UUID) to be bound
// automatically without requiring an explicit parser tag.
//
// Note: time.Time is excluded from automatic TextUnmarshaler handling because gin
// provides dedicated time parsing via time_format, time_utc, and time_location tags.
func (d bindingApi) TrySetByInterface(inputVal string, valueToSet reflect.Value) (isSet bool, err error) {
	switch v := valueToSet.Addr().Interface().(type) {
	case bindUnmarshaler:
		return true, v.UnmarshalParam(inputVal)
	case encoding.TextUnmarshaler:
		// Skip time.Time — it has dedicated handling in setWithProperType via setTimeField
		if _, isTime := valueToSet.Interface().(time.Time); !isTime {
			return true, v.UnmarshalText([]byte(inputVal))
		}
	}
	return false, nil
}
