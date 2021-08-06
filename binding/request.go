package binding

import (
	"errors"
	"net/http"
	"reflect"
)

var ErrInvalidTagInRequestBody = errors.New("body struct should not contain tag `query`, `header`, `cookie`, `uri` in binding request api")

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

	if err := b.bindingQuery(req, obj); err != nil {
		return err
	}

	binders := []Binding{Header, Cookie}
	for _, binder := range binders {
		if err := binder.BindOnly(req, obj); err != nil {
			return err
		}
	}

	// body decode
	bodyObj := extractBody(obj)
	if bodyObj == nil {
		return nil
	}

	// default json
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = MIMEJSON
	}
	bb := Default(req.Method, contentType)
	return bb.BindOnly(req, bodyObj)

}

func (b requestBinding) bindingQuery(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	return mapFormByTag(obj, values, "query")
}

// extractBody return body object
func extractBody(obj interface{}) interface{} {

	// pre-check obj
	rv := reflect.ValueOf(obj)
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return nil
	}

	return extract(rv)
}

func extract(rv reflect.Value) interface{} {

	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		tf := typ.Field(i)
		vf := rv.Field(i)

		_, ok := tf.Tag.Lookup("body")
		if !ok {
			continue
		}

		// find body struct
		if reflect.Indirect(vf).Kind() == reflect.Struct {
			// body must not has tag "query"
			if hasTag(vf, "query") || hasTag(vf, "header") ||
				hasTag(vf, "cookie") || hasTag(vf, "uri") {
				panic(ErrInvalidTagInRequestBody)
			}

			return vf.Addr().Interface()
		}
	}

	return nil
}

func hasTag(rv reflect.Value, tag string) bool {
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return false
	}

	typ := rv.Type()
	for i := 0; i < typ.NumField(); i++ {
		_, ok := typ.Field(i).Tag.Lookup(tag)
		if ok {
			return true
		}
	}

	return false
}
