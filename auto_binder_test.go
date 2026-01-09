package gin

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myRequest struct {
	Field1 string `json:"field_1"`
}

func TestAutoBinder_isFunc(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{
			"valid function",
			func(string, int) error { return nil },
			true,
		},
		{
			"valid zero-param function",
			func() error { return nil },
			true,
		},
		{
			"invalid function",
			func() string { return "" }(),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isFunc(tt.input)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestAutoBinder_isGinContext(t *testing.T) {
	assert.True(t, isGinContext(reflect.TypeOf(&Context{})))
	assert.False(t, isGinContext(reflect.TypeOf(Context{})))
	assert.False(t, isGinContext(reflect.TypeOf([]string{})))
}

func TestAutoBinder_constructStruct_pointer(t *testing.T) {
	type myType struct {
		Field int `json:"field"`
	}

	rv, err := constructStruct(reflect.TypeOf(&myType{}), func(obj any) error {
		assert.True(t, isPtr(reflect.TypeOf(obj)))

		return json.Unmarshal(
			[]byte(`{"field": 10}`),
			obj,
		)
	})

	assert.NoError(t, err)

	instance, ok := rv.Interface().(*myType)

	assert.True(t, ok)

	assert.Equal(t, 10, instance.Field)
}

func TestAutoBinder_constructStruct_nonPointer(t *testing.T) {
	type myType struct {
		Field int `json:"field"`
	}

	rv, err := constructStruct(reflect.TypeOf(myType{}), func(obj any) error {
		assert.True(t, isPtr(reflect.TypeOf(obj)))

		return json.Unmarshal(
			[]byte(`{"field": 10}`),
			obj,
		)
	})

	assert.NoError(t, err)

	instance, ok := rv.Interface().(myType)

	assert.True(t, ok)

	assert.Equal(t, 10, instance.Field)
}

func TestAutoBinder_constructStruct_nonStruct(t *testing.T) {
	_, err := constructStruct(reflect.TypeOf("string test"), func(obj any) error {
		assert.True(t, isPtr(reflect.TypeOf(obj)))

		return json.Unmarshal(
			[]byte(`{"field": 10}`),
			obj,
		)
	})

	assert.Error(t, err)
}

func TestAutoBinder_callHandler(t *testing.T) {
	called := false

	handler := func(ctx *Context, req *myRequest) {
		if ctx == nil {
			t.Errorf("ctx should not passed as nil")
			return
		}

		if req.Field1 != "value1" {
			t.Errorf("expected %v, actual %v", "value1", req.Field1)
		}

		called = true
	}

	rt := reflect.TypeOf(handler)
	rv := reflect.ValueOf(handler)

	ctx := &Context{}

	err := callHandler(rt, rv, ctx, func(obj any) error {
		return json.Unmarshal([]byte(`{"field_1": "value1"}`), obj)
	})

	if err != nil {
		panic(err)
	}

	if !called {
		t.Error("handler should be called")
	}

}
