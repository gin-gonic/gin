package binding

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"
)

type plainBinding struct{}

func (plainBinding) Name() string {
	return "plain"
}

func (plainBinding) Bind(req *http.Request, obj interface{}) error {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	all, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if v.Kind() == reflect.String {
		v.SetString(*(*string)(unsafe.Pointer(&all)))
		return nil
	}

	if _, ok := v.Interface().([]byte); ok {
		v.SetBytes(all)
		return nil
	}

	return fmt.Errorf("type (%T) unkown type", v)
}
