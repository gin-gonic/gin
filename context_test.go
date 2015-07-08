// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/manucorporat/sse"
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

func createMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	must(mw.SetBoundary(boundary))
	must(mw.WriteField("foo", "bar"))
	must(mw.WriteField("bar", "foo"))
	must(mw.WriteField("bar", "foo2"))
	must(mw.WriteField("array", "first"))
	must(mw.WriteField("array", "second"))
	req, err := http.NewRequest("POST", "/", body)
	must(err)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func TestContextReset(t *testing.T) {
	router := New()
	c := router.allocateContext()
	assert.Equal(t, c.engine, router)

	c.index = 2
	c.Writer = &responseWriter{ResponseWriter: httptest.NewRecorder()}
	c.Params = Params{Param{}}
	c.Error(errors.New("test"))
	c.Set("foo", "bar")
	c.reset()

	assert.False(t, c.IsAborted())
	assert.Nil(t, c.Keys)
	assert.Nil(t, c.Accepted)
	assert.Len(t, c.Errors, 0)
	assert.Empty(t, c.Errors.Errors())
	assert.Empty(t, c.Errors.ByType(ErrorTypeAny))
	assert.Len(t, c.Params, 0)
	assert.EqualValues(t, c.index, -1)
	assert.Equal(t, c.Writer.(*responseWriter), &c.writermem)
}

func TestContextHandlers(t *testing.T) {
	c, _, _ := createTestContext()
	assert.Nil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	c.handlers = HandlersChain{}
	assert.NotNil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	f := func(c *Context) {}
	g := func(c *Context) {}

	c.handlers = HandlersChain{f}
	compareFunc(t, f, c.handlers.Last())

	c.handlers = HandlersChain{f, g}
	compareFunc(t, g, c.handlers.Last())
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
	assert.Equal(t, cp.index, abortIndex)
	assert.Equal(t, cp.Keys, c.Keys)
	assert.Equal(t, cp.engine, c.engine)
	assert.Equal(t, cp.Params, c.Params)
}

func TestContextHandlerName(t *testing.T) {
	c, _, _ := createTestContext()
	c.handlers = HandlersChain{func(c *Context) {}, handlerNameTest}

	assert.Equal(t, c.HandlerName(), "github.com/gin-gonic/gin.handlerNameTest")
}

func handlerNameTest(c *Context) {

}

func TestContextQuery(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10", nil)

	assert.Equal(t, c.DefaultQuery("foo", "none"), "bar")
	assert.Equal(t, c.Query("foo"), "bar")
	assert.Empty(t, c.PostForm("foo"))

	assert.Equal(t, c.DefaultQuery("page", "0"), "10")
	assert.Equal(t, c.Query("page"), "10")
	assert.Empty(t, c.PostForm("page"))

	assert.Equal(t, c.DefaultQuery("NoKey", "nada"), "nada")
	assert.Empty(t, c.Query("NoKey"))
	assert.Empty(t, c.PostForm("NoKey"))
}

func TestContextQueryAndPostForm(t *testing.T) {
	c, _, _ := createTestContext()
	body := bytes.NewBufferString("foo=bar&page=11&both=POST&foo=second")
	c.Request, _ = http.NewRequest("POST", "/?both=GET&id=main&id=omit&array[]=first&array[]=second", body)
	c.Request.Header.Add("Content-Type", MIMEPOSTForm)

	assert.Equal(t, c.DefaultPostForm("foo", "none"), "bar")
	assert.Equal(t, c.PostForm("foo"), "bar")
	assert.Empty(t, c.Query("foo"))

	assert.Equal(t, c.DefaultPostForm("page", "0"), "11")
	assert.Equal(t, c.PostForm("page"), "11")
	assert.Equal(t, c.Query("page"), "")

	assert.Equal(t, c.PostForm("both"), "POST")
	assert.Equal(t, c.Query("both"), "GET")

	assert.Equal(t, c.DefaultPostForm("id", "000"), "000")
	assert.Equal(t, c.Query("id"), "main")
	assert.Empty(t, c.PostForm("id"))

	assert.Equal(t, c.DefaultPostForm("NoKey", "nada"), "nada")
	assert.Empty(t, c.PostForm("NoKey"))
	assert.Empty(t, c.Query("NoKey"))

	var obj struct {
		Foo   string   `form:"foo"`
		ID    string   `form:"id"`
		Page  string   `form:"page"`
		Both  string   `form:"both"`
		Array []string `form:"array[]"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.ID, "main")
	assert.Equal(t, obj.Page, "11")
	assert.Equal(t, obj.Both, "POST")
	assert.Equal(t, obj.Array, []string{"first", "second"})
}

func TestContextPostFormMultipart(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request = createMultipartRequest()

	var obj struct {
		Foo   string   `form:"foo"`
		Bar   string   `form:"bar"`
		Array []string `form:"array"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Array, []string{"first", "second"})

	assert.Empty(t, c.Query("foo"))
	assert.Empty(t, c.Query("bar"))
	assert.Equal(t, c.PostForm("foo"), "bar")
	assert.Equal(t, c.PostForm("array"), "first")
	assert.Equal(t, c.PostForm("bar"), "foo")
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
// we change the content-type before
func TestContextRenderAPIJSON(t *testing.T) {
	c, w, _ := createTestContext()
	c.Header("Content-Type", "application/vnd.api+json")
	c.JSON(201, H{"foo": "bar"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "{\"foo\":\"bar\"}\n")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/vnd.api+json")
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
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(201, "<html>%s %d</html>", "string", 3)

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

func TestContextRenderSSE(t *testing.T) {
	c, w, _ := createTestContext()
	c.SSEvent("float", 1.5)
	c.Render(-1, sse.Event{
		Id:   "123",
		Data: "text",
	})
	c.SSEvent("chat", H{
		"foo": "bar",
		"bar": "foo",
	})

	assert.Equal(t, w.Body.String(), "event: float\ndata: 1.5\n\nid: 123\ndata: text\n\nevent: chat\ndata: {\"bar\":\"foo\",\"foo\":\"bar\"}\n\n")
}

func TestContextRenderFile(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.File("./gin.go")

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "func New() *Engine {")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

func TestContextHeaders(t *testing.T) {
	c, _, _ := createTestContext()
	c.Header("Content-Type", "text/plain")
	c.Header("X-Custom", "value")

	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "text/plain")
	assert.Equal(t, c.Writer.Header().Get("X-Custom"), "value")

	c.Header("Content-Type", "text/html")
	c.Header("X-Custom", "")

	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "text/html")
	_, exist := c.Writer.Header()["X-Custom"]
	assert.False(t, exist)
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
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEHTML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), "")
}

func TestContextNegotiationFormatCustum(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	c.Accepted = nil
	c.SetAccepted(MIMEJSON, MIMEXML)

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEJSON)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), MIMEJSON)
}

func TestContextIsAborted(t *testing.T) {
	c, _, _ := createTestContext()
	assert.False(t, c.IsAborted())

	c.Abort()
	assert.True(t, c.IsAborted())

	c.Next()
	assert.True(t, c.IsAborted())

	c.index++
	assert.True(t, c.IsAborted())
}

// TestContextData tests that the response can be written from `bytesting`
// with specified MIME type
func TestContextAbortWithStatus(t *testing.T) {
	c, w, _ := createTestContext()
	c.index = 4
	c.AbortWithStatus(401)
	c.Writer.WriteHeaderNow()

	assert.Equal(t, c.index, abortIndex)
	assert.Equal(t, c.Writer.Status(), 401)
	assert.Equal(t, w.Code, 401)
	assert.True(t, c.IsAborted())
}

func TestContextError(t *testing.T) {
	c, _, _ := createTestContext()
	assert.Empty(t, c.Errors)

	c.Error(errors.New("first error"))
	assert.Len(t, c.Errors, 1)
	assert.Equal(t, c.Errors.String(), "Error #01: first error\n")

	c.Error(&Error{
		Err:  errors.New("second error"),
		Meta: "some data 2",
		Type: ErrorTypePublic,
	})
	assert.Len(t, c.Errors, 2)

	assert.Equal(t, c.Errors[0].Err, errors.New("first error"))
	assert.Nil(t, c.Errors[0].Meta)
	assert.Equal(t, c.Errors[0].Type, ErrorTypePrivate)

	assert.Equal(t, c.Errors[1].Err, errors.New("second error"))
	assert.Equal(t, c.Errors[1].Meta, "some data 2")
	assert.Equal(t, c.Errors[1].Type, ErrorTypePublic)

	assert.Equal(t, c.Errors.Last(), c.Errors[1])
}

func TestContextTypedError(t *testing.T) {
	c, _, _ := createTestContext()
	c.Error(errors.New("externo 0")).SetType(ErrorTypePublic)
	c.Error(errors.New("interno 0")).SetType(ErrorTypePrivate)

	for _, err := range c.Errors.ByType(ErrorTypePublic) {
		assert.Equal(t, err.Type, ErrorTypePublic)
	}
	for _, err := range c.Errors.ByType(ErrorTypePrivate) {
		assert.Equal(t, err.Type, ErrorTypePrivate)
	}
	assert.Equal(t, c.Errors.Errors(), []string{"externo 0", "interno 0"})
}

func TestContextAbortWithError(t *testing.T) {
	c, w, _ := createTestContext()
	c.AbortWithError(401, errors.New("bad input")).SetMeta("some input")
	c.Writer.WriteHeaderNow()

	assert.Equal(t, w.Code, 401)
	assert.Equal(t, c.index, abortIndex)
	assert.True(t, c.IsAborted())
}

func TestContextClientIP(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", nil)

	c.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	c.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	c.Request.RemoteAddr = "  40.40.40.40 "

	assert.Equal(t, c.ClientIP(), "10.10.10.10")

	c.Request.Header.Del("X-Real-IP")
	assert.Equal(t, c.ClientIP(), "20.20.20.20")

	c.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assert.Equal(t, c.ClientIP(), "30.30.30.30")

	c.Request.Header.Del("X-Forwarded-For")
	assert.Equal(t, c.ClientIP(), "40.40.40.40")
}

func TestContextContentType(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, c.ContentType(), "application/json")
}

func TestContextAutoBindJSON(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", MIMEJSON)

	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, obj.Bar, "foo")
	assert.Equal(t, obj.Foo, "bar")
	assert.Empty(t, c.Errors)
}

func TestContextBindWithJSON(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", MIMEXML) // set fake content-type

	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	assert.NoError(t, c.BindJSON(&obj))
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
	assert.Error(t, c.Bind(&obj))
	c.Writer.WriteHeaderNow()

	assert.Empty(t, obj.Bar)
	assert.Empty(t, obj.Foo)
	assert.Equal(t, w.Code, 400)
	assert.True(t, c.IsAborted())
}

func TestContextGolangContext(t *testing.T) {
	c, _, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	assert.NoError(t, c.Err())
	assert.Nil(t, c.Done())
	ti, ok := c.Deadline()
	assert.Equal(t, ti, time.Time{})
	assert.False(t, ok)
	assert.Equal(t, c.Value(0), c.Request)
	assert.Nil(t, c.Value("foo"))

	c.Set("foo", "bar")
	assert.Equal(t, c.Value("foo"), "bar")
	assert.Nil(t, c.Value(1))
}
