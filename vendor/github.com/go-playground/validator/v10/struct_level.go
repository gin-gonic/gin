package validator

import (
	"context"
	"reflect"
)

// StructLevelFunc accepts all values needed for struct level validation
type StructLevelFunc func(sl StructLevel)

// StructLevelFuncCtx accepts all values needed for struct level validation
// but also allows passing of contextual validation information via context.Context.
type StructLevelFuncCtx func(ctx context.Context, sl StructLevel)

// wrapStructLevelFunc wraps normal StructLevelFunc makes it compatible with StructLevelFuncCtx
func wrapStructLevelFunc(fn StructLevelFunc) StructLevelFuncCtx {
	return func(ctx context.Context, sl StructLevel) {
		fn(sl)
	}
}

// StructLevel contains all the information and helper functions
// to validate a struct
type StructLevel interface {

	// Validator returns the main validation object, in case one wants to call validations internally.
	// this is so you don't have to use anonymous functions to get access to the validate
	// instance.
	Validator() *Validate

	// Top returns the top level struct, if any
	Top() reflect.Value

	// Parent returns the current fields parent struct, if any
	Parent() reflect.Value

	// Current returns the current struct.
	Current() reflect.Value

	// ExtractType gets the actual underlying type of field value.
	// It will dive into pointers, customTypes and return you the
	// underlying value and its kind.
	ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool)

	// ReportError reports an error just by passing the field and tag information
	//
	// NOTES:
	//
	// fieldName and structFieldName get appended to the existing
	// namespace that validator is on. e.g. pass 'FirstName' or
	// 'Names[0]' depending on the nesting
	//
	// tag can be an existing validation tag or just something you make up
	// and process on the flip side it's up to you.
	ReportError(field interface{}, fieldName, structFieldName string, tag, param string)

	// ReportValidationErrors reports an error just by passing ValidationErrors
	//
	// NOTES:
	//
	// relativeNamespace and relativeActualNamespace get appended to the
	// existing namespace that validator is on.
	// e.g. pass 'User.FirstName' or 'Users[0].FirstName' depending
	// on the nesting. most of the time they will be blank, unless you validate
	// at a level lower the current field depth
	ReportValidationErrors(relativeNamespace, relativeActualNamespace string, errs ValidationErrors)
}

var _ StructLevel = new(validate)

// Top returns the top level struct
//
// NOTE: this can be the same as the current struct being validated
// if not is a nested struct.
//
// this is only called when within Struct and Field Level validation and
// should not be relied upon for an accurate value otherwise.
func (v *validate) Top() reflect.Value {
	return v.top
}

// Parent returns the current structs parent
//
// NOTE: this can be the same as the current struct being validated
// if not is a nested struct.
//
// this is only called when within Struct and Field Level validation and
// should not be relied upon for an accurate value otherwise.
func (v *validate) Parent() reflect.Value {
	return v.slflParent
}

// Current returns the current struct.
func (v *validate) Current() reflect.Value {
	return v.slCurrent
}

// Validator returns the main validation object, in case one want to call validations internally.
func (v *validate) Validator() *Validate {
	return v.v
}

// ExtractType gets the actual underlying type of field value.
func (v *validate) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return v.extractTypeInternal(field, false)
}

// ReportError reports an error just by passing the field and tag information
func (v *validate) ReportError(field interface{}, fieldName, structFieldName, tag, param string) {
	fv, kind, _ := v.extractTypeInternal(reflect.ValueOf(field), false)

	if len(structFieldName) == 0 {
		structFieldName = fieldName
	}

	v.str1 = string(append(v.ns, fieldName...))

	if v.v.hasTagNameFunc || fieldName != structFieldName {
		v.str2 = string(append(v.actualNs, structFieldName...))
	} else {
		v.str2 = v.str1
	}

	if kind == reflect.Invalid {
		v.errs = append(v.errs,
			&fieldError{
				v:              v.v,
				tag:            tag,
				actualTag:      tag,
				ns:             v.str1,
				structNs:       v.str2,
				fieldLen:       uint8(len(fieldName)),
				structfieldLen: uint8(len(structFieldName)),
				param:          param,
				kind:           kind,
			},
		)
		return
	}

	v.errs = append(v.errs,
		&fieldError{
			v:              v.v,
			tag:            tag,
			actualTag:      tag,
			ns:             v.str1,
			structNs:       v.str2,
			fieldLen:       uint8(len(fieldName)),
			structfieldLen: uint8(len(structFieldName)),
			value:          getValue(fv),
			param:          param,
			kind:           kind,
			typ:            fv.Type(),
		},
	)
}

// ReportValidationErrors reports ValidationErrors obtained from running validations within the Struct Level validation.
//
// NOTE: this function prepends the current namespace to the relative ones.
func (v *validate) ReportValidationErrors(relativeNamespace, relativeStructNamespace string, errs ValidationErrors) {
	var err *fieldError

	for i := 0; i < len(errs); i++ {
		err = errs[i].(*fieldError)
		err.ns = string(append(append(v.ns, relativeNamespace...), err.ns...))
		err.structNs = string(append(append(v.actualNs, relativeStructNamespace...), err.structNs...))

		v.errs = append(v.errs, err)
	}
}
