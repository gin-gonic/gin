package binding

import (
	"net/http"
	"strings"
)

type allBinding struct{}

var _ BindingMany = allBinding{}

func (allBinding) Name() string {
	return "all"
}

func (allBinding) BindMany(req *http.Request, uriParams map[string][]string, obj any) error {
	// from binding.Header
	if err := mapHeader(obj, req.Header); err != nil {
		return err
	}

	// from binding.Uri
	if err := mapURI(obj, uriParams); err != nil {
		return err
	}

	// from binding.Query
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}

	// from context.Bind (for body/post-form/anything else)
	contentType := req.Header.Get("Content-Type")
	// trim contentType parameters, e.g. "application/json; charset=utf-8" -> "application/json"
	contentTypeLastIdx := strings.IndexAny(contentType, " ;")
	if contentTypeLastIdx != -1 {
		contentType = contentType[:contentTypeLastIdx]
	}
	b := Default(req.Method, contentType)
	// final validation done by whatever binding is selected here
	return b.Bind(req, obj)
}
