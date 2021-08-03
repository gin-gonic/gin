package binding

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Params struct {
	Name          string     `uri:"name"`
	Age           int        `form:"age,default=18"`
	Money         int32      `form:"money" binding:"required"`
	Authorization string     `cookie:"Authorization"`
	UserAgent     string     `header:"User-Agent"`
	Data          *ParamBody `body:"body"`
}

type ParamBody struct {
	Replicas *int32 `json:"replicas"`
}

func NewParams() *Params {
	return &Params{
		Data: &ParamBody{},
	}
}

func TestBindRequest(t *testing.T) {
	h := Request
	assert.Equal(t, "request", h.Name())

	params := NewParams()

	mocks := []struct {
		body   string
		mime   string
		method string
	}{
		{body: `replicas: 5`, mime: MIMEYAML, method: http.MethodPost},
		{body: `{"replicas": 5}`, mime: MIMEJSON, method: http.MethodPut},
		{body: `replicas=5`, mime: MIMEPOSTForm, method: http.MethodPatch},
		{body: `replicas=5`, mime: MIMEPOSTForm, method: http.MethodGet},
		{
			body:   `<?xml version="1.0" encoding="UTF-8" ?><replicas>5</replicas>`,
			mime:   MIMEXML2,
			method: http.MethodDelete,
		},
		{
			body:   `<replicas>5</replicas>`,
			mime:   MIMEXML,
			method: http.MethodDelete,
		},
	}

	for _, mock := range mocks {

		t.Run(mock.mime, func(t *testing.T) {
			req := requestWithBody(mock.method, "/hello/:name?money=1000", mock.body)
			mockUri := map[string][]string{
				"name": {"zhangsan"},
			}
			req.Header.Add("Content-Type", mock.mime)
			req.Header.Add("User-Agent", "go-client")
			req.AddCookie(
				&http.Cookie{Name: "Authorization", Value: "token 123123123"},
			)

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
