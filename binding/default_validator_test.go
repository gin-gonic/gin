// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestSliceValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  SliceValidationError
		want string
	}{
		{"has nil elements", SliceValidationError{errors.New("test error"), nil}, "[0]: test error"},
		{"has zero elements", SliceValidationError{}, ""},
		{"has one element", SliceValidationError{errors.New("test one error")}, "[0]: test one error"},
		{
			"has two elements",
			SliceValidationError{
				errors.New("first error"),
				errors.New("second error"),
			},
			"[0]: first error\n[1]: second error",
		},
		{
			"has many elements",
			SliceValidationError{
				errors.New("first error"),
				errors.New("second error"),
				nil,
				nil,
				nil,
				errors.New("last error"),
			},
			"[0]: first error\n[1]: second error\n[5]: last error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("SliceValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValidator(t *testing.T) {
	type exampleStruct struct {
		A string `binding:"max=8"`
		B int    `binding:"gt=0"`
	}
	tests := []struct {
		name    string
		v       *defaultValidator
		obj     any
		wantErr bool
	}{
		{"validate nil obj", &defaultValidator{}, nil, false},
		{"validate int obj", &defaultValidator{}, 3, false},
		{"validate struct failed-1", &defaultValidator{}, exampleStruct{A: "123456789", B: 1}, true},
		{"validate struct failed-2", &defaultValidator{}, exampleStruct{A: "12345678", B: 0}, true},
		{"validate struct passed", &defaultValidator{}, exampleStruct{A: "12345678", B: 1}, false},
		{"validate *struct failed-1", &defaultValidator{}, &exampleStruct{A: "123456789", B: 1}, true},
		{"validate *struct failed-2", &defaultValidator{}, &exampleStruct{A: "12345678", B: 0}, true},
		{"validate *struct passed", &defaultValidator{}, &exampleStruct{A: "12345678", B: 1}, false},
		{"validate []struct failed-1", &defaultValidator{}, []exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []struct failed-2", &defaultValidator{}, []exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []struct passed", &defaultValidator{}, []exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate []*struct failed-1", &defaultValidator{}, []*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate []*struct failed-2", &defaultValidator{}, []*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate []*struct passed", &defaultValidator{}, []*exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]struct failed-1", &defaultValidator{}, &[]exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]struct failed-2", &defaultValidator{}, &[]exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]struct passed", &defaultValidator{}, &[]exampleStruct{{A: "12345678", B: 1}}, false},
		{"validate *[]*struct failed-1", &defaultValidator{}, &[]*exampleStruct{{A: "123456789", B: 1}}, true},
		{"validate *[]*struct failed-2", &defaultValidator{}, &[]*exampleStruct{{A: "12345678", B: 0}}, true},
		{"validate *[]*struct passed", &defaultValidator{}, &[]*exampleStruct{{A: "12345678", B: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.ValidateStruct(tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("defaultValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStructValidatorRequiredStruct(t *testing.T) {
	type myUUID struct {
		S string
	}
	type createReq struct {
		ID myUUID `binding:"required"`
	}

	// without the option, required is a no-op on non-pointer struct fields
	if err := (&defaultValidator{}).ValidateStruct(createReq{}); err != nil {
		t.Fatalf("default validator: unexpected error: %v", err)
	}

	v := NewStructValidator(validator.WithRequiredStructEnabled())
	if err := v.ValidateStruct(createReq{}); err == nil {
		t.Fatal("WithRequiredStructEnabled: expected required error on zero struct field")
	}
	if err := v.ValidateStruct(createReq{ID: myUUID{S: "x"}}); err != nil {
		t.Fatalf("WithRequiredStructEnabled: unexpected error for non-zero field: %v", err)
	}

	engine, ok := v.Engine().(*validator.Validate)
	if !ok || engine == nil {
		t.Fatalf("Engine() = %T, want *validator.Validate", v.Engine())
	}
}
