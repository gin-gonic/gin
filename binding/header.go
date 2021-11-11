package binding

import (
	"context"
	"net/http"
	"net/textproto"
	"reflect"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (b headerBinding) Bind(req *http.Request, obj interface{}) error {
	return b.BindContext(context.Background(), req, obj)
}

func (headerBinding) BindContext(ctx context.Context, req *http.Request, obj interface{}) error {
	if err := mapHeader(obj, req.Header); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}

func mapHeader(ptr interface{}, h map[string][]string) error {
	return mappingByPtr(ptr, headerSource(h), "header")
}

type headerSource map[string][]string

var _ setter = headerSource(nil)

func (hs headerSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (bool, error) {
	return setByForm(value, field, hs, textproto.CanonicalMIMEHeaderKey(tagValue), opt)
}
