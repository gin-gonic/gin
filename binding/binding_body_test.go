package binding

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindingBody(t *testing.T) {
	for _, tt := range []struct {
		name    string
		binding BindingBody
		body    string
		want    string
	}{
		{
			name:    "JSON binding",
			binding: JSON,
			body:    `{"foo":"FOO"}`,
		},
		{
			name:    "XML binding",
			binding: XML,
			body: `<?xml version="1.0" encoding="UTF-8"?>
<root>
   <foo>FOO</foo>
</root>`,
		},
		{
			name:    "YAML binding",
			binding: YAML,
			body:    `foo: FOO`,
		},
	} {
		t.Logf("testing: %s", tt.name)
		req := requestWithBody("POST", "/", tt.body)
		form := FooStruct{}
		body, _ := ioutil.ReadAll(req.Body)
		assert.NoError(t, tt.binding.BindBody(body, &form))
		assert.Equal(t, FooStruct{"FOO"}, form)
	}
}
