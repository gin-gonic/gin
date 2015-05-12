// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

// Unit tests TODO
// func (c *Context) File(filepath string) {
// func (c *Context) Negotiate(code int, config Negotiate) {
// BAD case: func (c *Context) Render(code int, render render.Render, obj ...interface{}) {
// test that information is not leaked when reusing Contexts (using the Pool)

func createTestContext() (c *Context, w *httptest.ResponseRecorder, r *Engine) {
	w = httptest.NewRecorder()
	r = New()
	c = r.allocateContext()
	c.reset()
	c.writermem.reset(w)
	return
}

func TestContextReset(t *testing.T) {
	router := New()
	c := router.allocateContext()
	assert.Equal(t, c.Engine, router)

	c.index = 2
	c.Writer = &responseWriter{ResponseWriter: httptest.NewRecorder()}
	c.Params = Params{Param{}}
	c.Error(errors.New("test"), nil)
	c.Set("foo", "bar")
	c.reset()

	assert.False(t, c.IsAborted())
	assert.Nil(t, c.Keys)
	assert.Nil(t, c.Accepted)
	assert.Len(t, c.Errors, 0)
	assert.Len(t, c.Params, 0)
	assert.Equal(t, c.index, -1)
	assert.Equal(t, c.Writer.(*responseWriter), &c.writermem)
}

// TestContextSetGet tests that a parameter is set correctly on the
// current context and can be retrieved using Get.
func TestContextSetGet(t *testing.T) {
	c, _, _ := createTestContext()
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, value, "bar")
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)

	assert.Equal(t, c.MustGet("foo"), "bar")
	assert.Panics(t, func() { c.MustGet("no_exist") })
}

func TestContextSetGetValues(t *testing.T) {
	c, _, _ := createTestContext()
	c.Set("string", "this is a string")
	c.Set("int32", int32(-42))
	c.Set("int64", int64(42424242424242))
	c.Set("uint64", uint64(42))
	c.Set("float32", float32(4.2))
	c.Set("float64", 4.2)
	var a interface{} = 1
	c.Set("intInterface", a)

	assert.Exactly(t, c.MustGet("string").(string), "this is a string")
	assert.Exactly(t, c.MustGet("int32").(int32), int32(-42))
	assert.Exactly(t, c.MustGet("int64").(int64), int64(42424242424242))
	assert.Exactly(t, c.MustGet("uint64").(uint64), uint64(42))
	assert.Exactly(t, c.MustGet("float32").(float32), float32(4.2))
	assert.Exactly(t, c.MustGet("float64").(float64), 4.2)
	assert.Exactly(t, c.MustGet("intInterface").(int), 1)

}

func TestContextCopy(t *testing.T) {
	c, _, _ := createTestContext()
	c.index = 2
	c.Request, _ = http.NewRequest("POST", "/hola", nil)
	c.handlers = HandlersChain{func(c *Context) {}}
	c.Params = Params{Param{Key: "foo", Value: "bar"}}
	c.Set("foo", "bar")

	cp := c.Copy()
	assert.Nil(t, cp.handlers)
	assert.Nil(t, cp.writermem.ResponseWriter)
	assert.Equal(t, &cp.writermem, cp.Writer.(*responseWriter))
	assert.Equal(t, cp.Request, c.Request)
	assert.Equal(t, cp.index, AbortIndex)
	assert.Equal(t, cp.Keys, c.Keys)
	assert.Equal(t, cp.Engine, c.Engine)
	assert.Equal(t, cp.Params, c.Params)
}

func TestContextFormParse(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10", nil)

	assert.Equal(t, c.DefaultFormValue("foo", "none"), "bar")
	assert.Equal(t, c.FormValue("foo"), "bar")
	assert.Empty(t, c.PostFormValue("foo"))

	assert.Equal(t, c.DefaultFormValue("page", "0"), "10")
	assert.Equal(t, c.FormValue("page"), "10")
	assert.Empty(t, c.PostFormValue("page"))

	assert.Equal(t, c.DefaultFormValue("NoKey", "nada"), "nada")
	assert.Empty(t, c.FormValue("NoKey"))
	assert.Empty(t, c.PostFormValue("NoKey"))

}

func TestContextPostFormParse(t *testing.T) {
	c, _, _ := createTestContext()
	body := bytes.NewBufferString("foo=bar&page=11&both=POST")
	c.Request, _ = http.NewRequest("POST", "http://example.com/?both=GET&id=main", body)
	c.Request.Header.Add("Content-Type", MIMEPOSTForm)

	assert.Equal(t, c.DefaultPostFormValue("foo", "none"), "bar")
	assert.Equal(t, c.PostFormValue("foo"), "bar")
	assert.Equal(t, c.FormValue("foo"), "bar")

	assert.Equal(t, c.DefaultPostFormValue("page", "0"), "11")
	assert.Equal(t, c.PostFormValue("page"), "11")
	assert.Equal(t, c.FormValue("page"), "11")

	assert.Equal(t, c.PostFormValue("both"), "POST")
	assert.Equal(t, c.FormValue("both"), "POST")

	assert.Equal(t, c.FormValue("id"), "main")
	assert.Empty(t, c.PostFormValue("id"))

	assert.Equal(t, c.DefaultPostFormValue("NoKey", "nada"), "nada")
	assert.Empty(t, c.PostFormValue("NoKey"))
	assert.Empty(t, c.FormValue("NoKey"))
}

// Tests that the response is serialized as JSON
// and Content-Type is set to application/json
func TestContextRenderJSON(t *testing.T) {
	c, w, _ := createTestContext()
	c.JSON(201, H{"foo": "bar"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "{\"foo\":\"bar\"}\n")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

// Tests that the response is serialized as JSON
// and Content-Type is set to application/json
func TestContextRenderIndentedJSON(t *testing.T) {
	c, w, _ := createTestContext()
	c.IndentedJSON(201, H{"foo": "bar", "bar": "foo", "nested": H{"foo": "bar"}})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\",\n    \"nested\": {\n        \"foo\": \"bar\"\n    }\n}")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

// Tests that the response executes the templates
// and responds with Content-Type set to text/html
func TestContextRenderHTML(t *testing.T) {
	c, w, router := createTestContext()
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
	router.SetHTMLTemplate(templ)

	c.HTML(201, "t", H{"name": "alexandernyquist"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "Hello alexandernyquist")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// TestContextXML tests that the response is serialized as XML
// and Content-Type is set to application/xml
func TestContextRenderXML(t *testing.T) {
	c, w, _ := createTestContext()
	c.XML(201, H{"foo": "bar"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "<map><foo>bar</foo></map>")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/xml; charset=utf-8")
}

// TestContextString tests that the response is returned
// with Content-Type set to text/plain
func TestContextRenderString(t *testing.T) {
	c, w, _ := createTestContext()
	c.String(201, "test %s %d", "string", 2)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "test string 2")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

// TestContextString tests that the response is returned
// with Content-Type set to text/html
func TestContextRenderHTMLString(t *testing.T) {
	c, w, _ := createTestContext()
	c.HTMLString(201, "<html>%s %d</html>", "string", 3)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "<html>string 3</html>")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// TestContextData tests that the response can be written from `bytesting`
// with specified MIME type
func TestContextRenderData(t *testing.T) {
	c, w, _ := createTestContext()
	c.Data(201, "text/csv", []byte(`foo,bar`))

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "foo,bar")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/csv")
}

// TODO
func TestContextRenderRedirectWithRelativePath(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	assert.Panics(t, func() { c.Redirect(299, "/new_path") })
	assert.Panics(t, func() { c.Redirect(309, "/new_path") })

	c.Redirect(302, "/path")
	c.Writer.WriteHeaderNow()
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.Header().Get("Location"), "/path")
}

func TestContextRenderRedirectWithAbsolutePath(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	c.Redirect(302, "http://google.com")
	c.Writer.WriteHeaderNow()

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.Header().Get("Location"), "http://google.com")
}

func TestContextNegotiationFormat(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "", nil)

	assert.Panics(t, func() { c.NegotiateFormat() })
	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEJSON)
	assert.Equal(t, c.NegotiateFormat(MIMEHTML, MIMEJSON), MIMEHTML)
}

func TestContextNegotiationFormatWithAccept(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEHTML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), "")
}

func TestContextNegotiationFormatCustum(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	c.Accepted = nil
	c.SetAccepted(MIMEJSON, MIMEXML)

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEJSON)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), MIMEJSON)
}

// TestContextData tests that the response can be written from `bytesting`
// with specified MIME type
func TestContextAbortWithStatus(t *testing.T) {
	c, w, _ := createTestContext()
	c.index = 4
	c.AbortWithStatus(401)
	c.Writer.WriteHeaderNow()

	assert.Equal(t, c.index, AbortIndex)
	assert.Equal(t, c.Writer.Status(), 401)
	assert.Equal(t, w.Code, 401)
	assert.True(t, c.IsAborted())
}

func TestContextError(t *testing.T) {
	c, _, _ := createTestContext()
	assert.Nil(t, c.LastError())
	assert.Empty(t, c.Errors.String())

	c.Error(errors.New("first error"), "some data")
	assert.Equal(t, c.LastError().Error(), "first error")
	assert.Len(t, c.Errors, 1)
	assert.Equal(t, c.Errors.String(), "Error #01: first error\n     Meta: some data\n")

	c.Error(errors.New("second error"), "some data 2")
	assert.Equal(t, c.LastError().Error(), "second error")
	assert.Len(t, c.Errors, 2)
	assert.Equal(t, c.Errors.String(), "Error #01: first error\n     Meta: some data\n"+
		"Error #02: second error\n     Meta: some data 2\n")

	assert.Equal(t, c.Errors[0].Error, errors.New("first error"))
	assert.Equal(t, c.Errors[0].Meta, "some data")
	assert.Equal(t, c.Errors[0].Flags, ErrorTypeExternal)

	assert.Equal(t, c.Errors[1].Error, errors.New("second error"))
	assert.Equal(t, c.Errors[1].Meta, "some data 2")
	assert.Equal(t, c.Errors[1].Flags, ErrorTypeExternal)
}

func TestContextTypedError(t *testing.T) {
	c, _, _ := createTestContext()
	c.ErrorTyped(errors.New("externo 0"), ErrorTypeExternal, nil)
	c.ErrorTyped(errors.New("externo 1"), ErrorTypeExternal, nil)
	c.ErrorTyped(errors.New("interno 0"), ErrorTypeInternal, nil)
	c.ErrorTyped(errors.New("externo 2"), ErrorTypeExternal, nil)
	c.ErrorTyped(errors.New("interno 1"), ErrorTypeInternal, nil)
	c.ErrorTyped(errors.New("interno 2"), ErrorTypeInternal, nil)

	for _, err := range c.Errors.ByType(ErrorTypeExternal) {
		assert.Equal(t, err.Flags, ErrorTypeExternal)
	}

	for _, err := range c.Errors.ByType(ErrorTypeInternal) {
		assert.Equal(t, err.Flags, ErrorTypeInternal)
	}
}

func TestContextFail(t *testing.T) {
	c, w, _ := createTestContext()
	c.Fail(401, errors.New("bad input"))
	c.Writer.WriteHeaderNow()

	assert.Equal(t, w.Code, 401)
	assert.Equal(t, c.LastError().Error(), "bad input")
	assert.Equal(t, c.index, AbortIndex)
	assert.True(t, c.IsAborted())
}

func TestContextClientIP(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "", nil)

	c.Request.Header.Set("X-Real-IP", "10.10.10.10")
	c.Request.Header.Set("X-Forwarded-For", "20.20.20.20 , 30.30.30.30")
	c.Request.RemoteAddr = "40.40.40.40"

	assert.Equal(t, c.ClientIP(), "10.10.10.10")
	c.Request.Header.Del("X-Real-IP")
	assert.Equal(t, c.ClientIP(), "20.20.20.20")
	c.Request.Header.Del("X-Forwarded-For")
	assert.Equal(t, c.ClientIP(), "40.40.40.40")
}

func TestContextContentType(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, c.ContentType(), "application/json")
}

func TestContextAutoBind(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", MIMEJSON)
	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	assert.True(t, c.Bind(&obj))
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, w.Body.Len(), 0)
}

func TestContextBadAutoBind(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", bytes.NewBufferString("\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", MIMEJSON)
	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	assert.False(t, c.IsAborted())
	assert.False(t, c.Bind(&obj))
	c.Writer.WriteHeaderNow()

	assert.Empty(t, obj.Bar)
	assert.Empty(t, obj.Foo)
	assert.Equal(t, w.Code, 400)
	assert.True(t, c.IsAborted())
}

func TestContextBindWith(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", MIMEXML)
	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	assert.True(t, c.BindWith(&obj, binding.JSON))
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, w.Body.Len(), 0)
}
