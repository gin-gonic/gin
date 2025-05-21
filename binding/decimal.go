package binding

import (
	"github.com/shopspring/decimal"
	"strings"
)

// CustomDecimal represents a decimal number that can be bound from form values.
// It supports values with leading dots (e.g. ".1" is parsed as "0.1").
type CustomDecimal struct {
	decimal.Decimal
}

// UnmarshalParam implements the binding.BindUnmarshaler interface.
// It converts form values to decimal.Decimal, with special handling for
// values that start with a dot (e.g. ".1" becomes "0.1").
func (cd *CustomDecimal) UnmarshalParam(val string) error {
	if strings.HasPrefix(val, ".") {
		val = "0" + val
	}

	dec, err := decimal.NewFromString(val)
	if err != nil {
		return err
	}

	cd.Decimal = dec
	return nil
}
