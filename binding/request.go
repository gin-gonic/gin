package binding

import (
	"net/http"
	"reflect"
)

type requestBinding struct{}

func (requestBinding) Name() string {
	return "request"
}

func (b requestBinding) Bind(obj interface{}, req *http.Request, form map[string][]string) error {
	if err := b.BindOnly(obj, req, form); err != nil {
		return err
	}

	return validate(obj)
}

func (b requestBinding) BindOnly(obj interface{}, req *http.Request, uriMap map[string][]string) error {

	if err := Uri.BindOnly(uriMap, obj); err != nil {
		return err
	}

	binders := []interface{}{Header, Query, Cookie}
	for _, binder := range binders {
		if b, ok := binder.(Binding); ok {
			if err := b.BindOnly(req, obj); err != nil {
				return err
			}
		}
	}

	// body decode
	bodyPtr := reflectx(obj)

	// default json
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = MIMEJSON
	}
	bb := Default(req.Method, contentType)
	return bb.BindOnly(req, bodyPtr)

}

func reflectx(obj interface{}) interface{} {

	// pre-check obj
	rv := reflect.ValueOf(obj)
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return nil
	}

	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		tf := typ.Field(i)
		vf := rv.Field(i)
		_, ok := tf.Tag.Lookup("body")
		if !ok {
			continue
		}

		// find body struct
		if vf.Kind() == reflect.Ptr && vf.Elem().Kind() == reflect.Struct {
			return vf.Interface()
		}
	}

	return nil
}
