package binding

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Params struct {
	Name          string `uri:"name"`
	Age           int    `query:"age,default=18"`
	Money         int32  `query:"money" binding:"required"`
	Authorization string `cookie:"Authorization"`
	UserAgent     string `header:"User-Agent"`
	Data          struct {
		Replicas *int32 `json:"replicas" yaml:"replicas" xml:"replicas" form:"replicas"`
	} `body:"body"`
}

func TestBindRequest(t *testing.T) {
	h := Request
	assert.Equal(t, "request", h.Name())

	path := "/hello/:name?money=1000"
	mocks := []struct {
		body   string
		mime   string
		method string
	}{
		{body: `replicas: 5`, mime: MIMEYAML, method: http.MethodPost},
		{body: `{"replicas": 5}`, mime: MIMEJSON, method: http.MethodPut},
		{
			body:   `<?xml version="1.0" encoding="UTF-8" ?><map><replicas>5</replicas></map>`,
			mime:   MIMEXML2,
			method: http.MethodDelete,
		},
		{
			body:   `<map><replicas>5</replicas></map>`,
			mime:   MIMEXML,
			method: http.MethodDelete,
		},
		{body: `replicas=5`, mime: MIMEPOSTForm, method: http.MethodPatch},
		{body: `replicas=5`, mime: MIMEPOSTForm, method: http.MethodGet},
	}

	for _, mock := range mocks {
		path := path

		t.Run(fmt.Sprintf("%s_%s", mock.method, mock.mime), func(t *testing.T) {

			req := httpRequest(mock.method, path, mock.body)
			mockUri := map[string][]string{
				"name": {"zhangsan"},
			}
			req.Header.Add("Content-Type", mock.mime)
			req.Header.Add("User-Agent", "go-client")
			req.AddCookie(
				&http.Cookie{Name: "Authorization", Value: "token 123123123"},
			)

			params := &Params{}
			err := h.Bind(params, req, mockUri)
			if err != nil {
				panic(err)
			}

			assert.NoError(t, err)
			assert.Equal(t, "zhangsan", params.Name)                 // uri
			assert.Equal(t, 18, params.Age)                          // form,defualt
			assert.Equal(t, int32(1000), params.Money)               // form,required
			assert.Equal(t, "token 123123123", params.Authorization) // cookie
			assert.Equal(t, "go-client", params.UserAgent)           // header
			assert.Equal(t, int32(5), *params.Data.Replicas)         // body,ptr
		})

	}
}

func httpRequest(method string, path string, body string) *http.Request {

	if method == http.MethodGet {
		path = fmt.Sprintf("%s&%s", path, body)
		req, _ := http.NewRequest(method, path, nil)

		return req
	}

	return requestWithBody(method, path, body)
}
