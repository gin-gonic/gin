package gin

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	defaultAutoBinderErrorHandler = func(ctx *Context, err error) {
		ctx.Error(err)
		ctx.Abort()
	}
)

type binderType func(obj any) error

func isFunc(obj any) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Func
}

func isGinContext(rt reflect.Type) bool {
	return rt == reflect.TypeOf((*Context)(nil))
}

func isPtr(rt reflect.Type) bool {
	return rt.Kind() == reflect.Pointer
}

func isStruct(rt reflect.Type) bool {
	return rt.Kind() == reflect.Struct
}

func constructStruct(prt reflect.Type, binder binderType) (reflect.Value, error) {
	var pInstancePtr any

	if isPtr(prt) {
		pInstancePtr = reflect.New(prt.Elem()).Interface()
	} else {
		pInstancePtr = reflect.New(prt).Interface()
	}

	if err := binder(pInstancePtr); err != nil {
		return reflect.Value{}, err
	}

	if prt.Kind() == reflect.Pointer {
		return reflect.ValueOf(pInstancePtr), nil
	}

	return reflect.ValueOf(pInstancePtr).Elem(), nil
}

func callHandler(rt reflect.Type, rv reflect.Value, ctx *Context, binder binderType) error {
	numberOfParams := rt.NumIn()

	var args []reflect.Value

	for i := 0; i < numberOfParams; i++ {
		prt := rt.In(i)

		if isGinContext(prt) {
			args = append(args, reflect.ValueOf(ctx))
			continue
		}

		if isStruct(prt) || isStruct(prt.Elem()) {
			if prv, err := constructStruct(prt, binder); err != nil {
				return err
			} else {
				args = append(args, prv)
			}
		}
	}

	rv.Call(args)

	return nil
}

// AutoBinder is a handler wrapper that binds the actual handler's request.
//
// Example: func MyGetHandler(ctx *gin.Context, request *MyRequest) {}
//
// engine.GET("/endpoint", gin.AutoBinder(MyGetHandler)) and you can handel the errors by passing a handler
//
// engine.GET("/endpoint", gin.AutoBinder(MyGetHandler, func(ctx *gin.Context, err error) {}))
func AutoBinder(handler any, errorHandler ...func(*Context, error)) HandlerFunc {
	rt := reflect.TypeOf(handler)

	if rt.Kind() != reflect.Func {
		panic(errors.New("invalid handler type"))
	}

	if rt.NumIn() == 0 {
		panic(fmt.Errorf("handler should have at least one parameter, handler: %v", rt.Name()))
	}

	return func(ctx *Context) {
		selectedErrorHandler := defaultAutoBinderErrorHandler
		if len(errorHandler) > 0 && errorHandler[0] != nil {
			selectedErrorHandler = errorHandler[0]
		}

		rt := reflect.TypeOf(handler)
		rv := reflect.ValueOf(handler)

		if err := callHandler(rt, rv, ctx, func(obj any) error {
			return ctx.ShouldBind(obj)
		}); err != nil {
			selectedErrorHandler(ctx, err)
		}
	}
}
