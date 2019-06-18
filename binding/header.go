package binding

import "net/http"

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req *http.Request, obj interface{}) error {

	if err := mapHeader(obj, req.Header); err != nil {
		return err
	}

	return validate(obj)
}
