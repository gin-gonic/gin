package validator

import "reflect"

// FieldLevel contains all the information and helper functions
// to validate a field
type FieldLevel interface {

	// Top returns the top level struct, if any
	Top() reflect.Value

	// Parent returns the current fields parent struct, if any or
	// the comparison value if called 'VarWithValue'
	Parent() reflect.Value

	// Field returns current field for validation
	Field() reflect.Value

	// FieldName returns the field's name with the tag
	// name taking precedence over the fields actual name.
	FieldName() string

	// StructFieldName returns the struct field's name
	StructFieldName() string

	// Param returns param for validation against current field
	Param() string

	// GetTag returns the current validations tag name
	GetTag() string

	// ExtractType gets the actual underlying type of field value.
	// It will dive into pointers, customTypes and return you the
	// underlying value and it's kind.
	ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool)

	// GetStructFieldOK traverses the parent struct to retrieve a specific field denoted by the provided namespace
	// in the param and returns the field, field kind and whether is was successful in retrieving
	// the field at all.
	//
	// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
	// could not be retrieved because it didn't exist.
	//
	// Deprecated: Use GetStructFieldOK2() instead which also return if the value is nullable.
	GetStructFieldOK() (reflect.Value, reflect.Kind, bool)

	// GetStructFieldOKAdvanced is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
	// the field and namespace allowing more extensibility for validators.
	//
	// Deprecated: Use GetStructFieldOKAdvanced2() instead which also return if the value is nullable.
	GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool)

	// GetStructFieldOK2 traverses the parent struct to retrieve a specific field denoted by the provided namespace
	// in the param and returns the field, field kind, if it's a nullable type and whether is was successful in retrieving
	// the field at all.
	//
	// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
	// could not be retrieved because it didn't exist.
	GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool)

	// GetStructFieldOKAdvanced2 is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
	// the field and namespace allowing more extensibility for validators.
	GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool)
}

var _ FieldLevel = new(validate)

// Field returns current field for validation
func (v *validate) Field() reflect.Value {
	return v.flField
}

// FieldName returns the field's name with the tag
// name taking precedence over the fields actual name.
func (v *validate) FieldName() string {
	return v.cf.altName
}

// GetTag returns the current validations tag name
func (v *validate) GetTag() string {
	return v.ct.tag
}

// StructFieldName returns the struct field's name
func (v *validate) StructFieldName() string {
	return v.cf.name
}

// Param returns param for validation against current field
func (v *validate) Param() string {
	return v.ct.param
}

// GetStructFieldOK returns Param returns param for validation against current field
//
// Deprecated: Use GetStructFieldOK2() instead which also return if the value is nullable.
func (v *validate) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	current, kind, _, found := v.getStructFieldOKInternal(v.slflParent, v.ct.param)
	return current, kind, found
}

// GetStructFieldOKAdvanced is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
// the field and namespace allowing more extensibility for validators.
//
// Deprecated: Use GetStructFieldOKAdvanced2() instead which also return if the value is nullable.
func (v *validate) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	current, kind, _, found := v.GetStructFieldOKAdvanced2(val, namespace)
	return current, kind, found
}

// GetStructFieldOK2 returns Param returns param for validation against current field
func (v *validate) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return v.getStructFieldOKInternal(v.slflParent, v.ct.param)
}

// GetStructFieldOKAdvanced2 is the same as GetStructFieldOK except that it accepts the parent struct to start looking for
// the field and namespace allowing more extensibility for validators.
func (v *validate) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return v.getStructFieldOKInternal(val, namespace)
}
