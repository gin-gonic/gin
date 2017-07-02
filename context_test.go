// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var _ context.Context = &Context{}

// Unit tests TODO
// func (c *Context) File(filepath string) {
// func (c *Context) Negotiate(code int, config Negotiate) {
// BAD case: func (c *Context) Render(code int, render render.Render, obj ...interface{}) {
// test that information is not leaked when reusing Contexts (using the Pool)

func createMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	must(mw.SetBoundary(boundary))
	must(mw.WriteField("foo", "bar"))
	must(mw.WriteField("bar", "10"))
	must(mw.WriteField("bar", "foo2"))
	must(mw.WriteField("array", "first"))
	must(mw.WriteField("array", "second"))
	must(mw.WriteField("id", ""))
	must(mw.WriteField("time_local", "31/12/2016 14:55"))
	must(mw.WriteField("time_utc", "31/12/2016 14:55"))
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

func TestContextFormFile(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		w.Write([]byte("test"))
	}
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", f.Filename)
	}
}

func TestContextMultipartForm(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.WriteField("foo", "bar")
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.MultipartForm()
	if assert.NoError(t, err) {
		assert.NotNil(t, f)
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
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	c, _ := CreateTestContext(httptest.NewRecorder())
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

func TestContextGetString(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("string", "this is a string")
	assert.Equal(t, "this is a string", c.GetString("string"))
}

func TestContextSetGetBool(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("bool", true)
	assert.Equal(t, true, c.GetBool("bool"))
}

func TestContextGetInt(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("int", 1)
	assert.Equal(t, 1, c.GetInt("int"))
}

func TestContextGetInt64(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("int64", int64(42424242424242))
	assert.Equal(t, int64(42424242424242), c.GetInt64("int64"))
}

func TestContextGetFloat64(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("float64", 4.2)
	assert.Equal(t, 4.2, c.GetFloat64("float64"))
}

func TestContextGetTime(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	t1, _ := time.Parse("1/2/2006 15:04:05", "01/01/2017 12:00:00")
	c.Set("time", t1)
	assert.Equal(t, t1, c.GetTime("time"))
}

func TestContextGetDuration(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("duration", time.Second)
	assert.Equal(t, time.Second, c.GetDuration("duration"))
}

func TestContextGetStringSlice(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("slice", []string{"foo"})
	assert.Equal(t, []string{"foo"}, c.GetStringSlice("slice"))
}

func TestContextGetStringMap(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	var m = make(map[string]interface{})
	m["foo"] = 1
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMap("map"))
	assert.Equal(t, 1, c.GetStringMap("map")["foo"])
}

func TestContextGetStringMapString(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	var m = make(map[string]string)
	m["foo"] = "bar"
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapString("map"))
	assert.Equal(t, "bar", c.GetStringMapString("map")["foo"])
}

func TestContextGetStringMapStringSlice(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	var m = make(map[string][]string)
	m["foo"] = []string{"foo"}
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapStringSlice("map"))
	assert.Equal(t, []string{"foo"}, c.GetStringMapStringSlice("map")["foo"])
}

func TestContextCopy(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.handlers = HandlersChain{func(c *Context) {}, handlerNameTest}

	assert.Regexp(t, "^(.*/vendor/)?github.com/gin-gonic/gin.handlerNameTest$", c.HandlerName())
}

func handlerNameTest(c *Context) {

}

var handlerTest HandlerFunc = func(c *Context) {

}

func TestContextHandler(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.handlers = HandlersChain{func(c *Context) {}, handlerTest}

	assert.Equal(t, reflect.ValueOf(handlerTest).Pointer(), reflect.ValueOf(c.Handler()).Pointer())
}

func TestContextQuery(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10&id=", nil)

	value, ok := c.GetQuery("foo")
	assert.True(t, ok)
	assert.Equal(t, value, "bar")
	assert.Equal(t, c.DefaultQuery("foo", "none"), "bar")
	assert.Equal(t, c.Query("foo"), "bar")

	value, ok = c.GetQuery("page")
	assert.True(t, ok)
	assert.Equal(t, value, "10")
	assert.Equal(t, c.DefaultQuery("page", "0"), "10")
	assert.Equal(t, c.Query("page"), "10")

	value, ok = c.GetQuery("id")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, c.DefaultQuery("id", "nada"), "")
	assert.Empty(t, c.Query("id"))

	value, ok = c.GetQuery("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, c.DefaultQuery("NoKey", "nada"), "nada")
	assert.Empty(t, c.Query("NoKey"))

	// postform should not mess
	value, ok = c.GetPostForm("page")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("foo"))
}

func TestContextQueryAndPostForm(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString("foo=bar&page=11&both=&foo=second")
	c.Request, _ = http.NewRequest("POST", "/?both=GET&id=main&id=omit&array[]=first&array[]=second", body)
	c.Request.Header.Add("Content-Type", MIMEPOSTForm)

	assert.Equal(t, c.DefaultPostForm("foo", "none"), "bar")
	assert.Equal(t, c.PostForm("foo"), "bar")
	assert.Empty(t, c.Query("foo"))

	value, ok := c.GetPostForm("page")
	assert.True(t, ok)
	assert.Equal(t, value, "11")
	assert.Equal(t, c.DefaultPostForm("page", "0"), "11")
	assert.Equal(t, c.PostForm("page"), "11")
	assert.Equal(t, c.Query("page"), "")

	value, ok = c.GetPostForm("both")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("both"))
	assert.Equal(t, c.DefaultPostForm("both", "nothing"), "")
	assert.Equal(t, c.Query("both"), "GET")

	value, ok = c.GetQuery("id")
	assert.True(t, ok)
	assert.Equal(t, value, "main")
	assert.Equal(t, c.DefaultPostForm("id", "000"), "000")
	assert.Equal(t, c.Query("id"), "main")
	assert.Empty(t, c.PostForm("id"))

	value, ok = c.GetQuery("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	value, ok = c.GetPostForm("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, c.DefaultPostForm("NoKey", "nada"), "nada")
	assert.Equal(t, c.DefaultQuery("NoKey", "nothing"), "nothing")
	assert.Empty(t, c.PostForm("NoKey"))
	assert.Empty(t, c.Query("NoKey"))

	var obj struct {
		Foo   string   `form:"foo"`
		ID    string   `form:"id"`
		Page  int      `form:"page"`
		Both  string   `form:"both"`
		Array []string `form:"array[]"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.ID, "main")
	assert.Equal(t, obj.Page, 11)
	assert.Equal(t, obj.Both, "")
	assert.Equal(t, obj.Array, []string{"first", "second"})

	values, ok := c.GetQueryArray("array[]")
	assert.True(t, ok)
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.QueryArray("array[]")
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.QueryArray("nokey")
	assert.Equal(t, 0, len(values))

	values = c.QueryArray("both")
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "GET", values[0])
}

func TestContextPostFormMultipart(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request = createMultipartRequest()

	var obj struct {
		Foo       string    `form:"foo"`
		Bar       string    `form:"bar"`
		BarAsInt  int       `form:"bar"`
		Array     []string  `form:"array"`
		ID        string    `form:"id"`
		TimeLocal time.Time `form:"time_local" time_format:"02/01/2006 15:04"`
		TimeUTC   time.Time `form:"time_utc" time_format:"02/01/2006 15:04" time_utc:"1"`
		BlankTime time.Time `form:"blank_time" time_format:"02/01/2006 15:04"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "10")
	assert.Equal(t, obj.BarAsInt, 10)
	assert.Equal(t, obj.Array, []string{"first", "second"})
	assert.Equal(t, obj.ID, "")
	assert.Equal(t, obj.TimeLocal.Format("02/01/2006 15:04"), "31/12/2016 14:55")
	assert.Equal(t, obj.TimeLocal.Location(), time.Local)
	assert.Equal(t, obj.TimeUTC.Format("02/01/2006 15:04"), "31/12/2016 14:55")
	assert.Equal(t, obj.TimeUTC.Location(), time.UTC)
	assert.True(t, obj.BlankTime.IsZero())

	value, ok := c.GetQuery("foo")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.Query("bar"))
	assert.Equal(t, c.DefaultQuery("id", "nothing"), "nothing")

	value, ok = c.GetPostForm("foo")
	assert.True(t, ok)
	assert.Equal(t, value, "bar")
	assert.Equal(t, c.PostForm("foo"), "bar")

	value, ok = c.GetPostForm("array")
	assert.True(t, ok)
	assert.Equal(t, value, "first")
	assert.Equal(t, c.PostForm("array"), "first")

	assert.Equal(t, c.DefaultPostForm("bar", "nothing"), "10")

	value, ok = c.GetPostForm("id")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("id"))
	assert.Empty(t, c.DefaultPostForm("id", "nothing"))

	value, ok = c.GetPostForm("nokey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, c.DefaultPostForm("nokey", "nothing"), "nothing")

	values, ok := c.GetPostFormArray("array")
	assert.True(t, ok)
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.PostFormArray("array")
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.PostFormArray("nokey")
	assert.Equal(t, 0, len(values))

	values = c.PostFormArray("foo")
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "bar", values[0])
}

func TestContextSetCookie(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.SetCookie("user", "gin", 1, "/", "localhost", true, true)
	assert.Equal(t, c.Writer.Header().Get("Set-Cookie"), "user=gin; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure")
}

func TestContextSetCookiePathEmpty(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.SetCookie("user", "gin", 1, "", "localhost", true, true)
	assert.Equal(t, c.Writer.Header().Get("Set-Cookie"), "user=gin; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure")
}

func TestContextGetCookie(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/get", nil)
	c.Request.Header.Set("Cookie", "user=gin")
	cookie, _ := c.Cookie("user")
	assert.Equal(t, cookie, "gin")

	_, err := c.Cookie("nokey")
	assert.Error(t, err)
}

func TestContextBodyAllowedForStatus(t *testing.T) {
	assert.Equal(t, false, bodyAllowedForStatus(102))
	assert.Equal(t, false, bodyAllowedForStatus(204))
	assert.Equal(t, false, bodyAllowedForStatus(304))
	assert.Equal(t, true, bodyAllowedForStatus(500))
}

type TestPanicRender struct {
}

func (*TestPanicRender) Render(http.ResponseWriter) error {
	return errors.New("TestPanicRender")
}

func (*TestPanicRender) WriteContentType(http.ResponseWriter) {}

func TestContextRenderPanicIfErr(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, "TestPanicRender", fmt.Sprint(r))
	}()
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Render(http.StatusOK, &TestPanicRender{})

	assert.Fail(t, "Panic not detected")
}

// Tests that the response is serialized as JSON
// and Content-Type is set to application/json
func TestContextRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.JSON(201, H{"foo": "bar"})

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

// Tests that no JSON is rendered if code is 204
func TestContextRenderNoContentJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.JSON(204, H{"foo": "bar"})

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

// Tests that the response is serialized as JSON
// we change the content-type before
func TestContextRenderAPIJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Header("Content-Type", "application/vnd.api+json")
	c.JSON(201, H{"foo": "bar"})

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/vnd.api+json", w.HeaderMap.Get("Content-Type"))
}

// Tests that no Custom JSON is rendered if code is 204
func TestContextRenderNoContentAPIJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Header("Content-Type", "application/vnd.api+json")
	c.JSON(204, H{"foo": "bar"})

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/vnd.api+json")
}

// Tests that the response is serialized as JSON
// and Content-Type is set to application/json
func TestContextRenderIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.IndentedJSON(201, H{"foo": "bar", "bar": "foo", "nested": H{"foo": "bar"}})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "{\n    \"bar\": \"foo\",\n    \"foo\": \"bar\",\n    \"nested\": {\n        \"foo\": \"bar\"\n    }\n}")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

// Tests that no Custom JSON is rendered if code is 204
func TestContextRenderNoContentIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.IndentedJSON(204, H{"foo": "bar", "bar": "foo", "nested": H{"foo": "bar"}})

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

// Tests that the response executes the templates
// and responds with Content-Type set to text/html
func TestContextRenderHTML(t *testing.T) {
	w := httptest.NewRecorder()
	c, router := CreateTestContext(w)
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
	router.SetHTMLTemplate(templ)

	c.HTML(201, "t", H{"name": "alexandernyquist"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "Hello alexandernyquist")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// Tests that no HTML is rendered if code is 204
func TestContextRenderNoContentHTML(t *testing.T) {
	w := httptest.NewRecorder()
	c, router := CreateTestContext(w)
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
	router.SetHTMLTemplate(templ)

	c.HTML(204, "t", H{"name": "alexandernyquist"})

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// TestContextXML tests that the response is serialized as XML
// and Content-Type is set to application/xml
func TestContextRenderXML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.XML(201, H{"foo": "bar"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "<map><foo>bar</foo></map>")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/xml; charset=utf-8")
}

// Tests that no XML is rendered if code is 204
func TestContextRenderNoContentXML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.XML(204, H{"foo": "bar"})

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/xml; charset=utf-8")
}

// TestContextString tests that the response is returned
// with Content-Type set to text/plain
func TestContextRenderString(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.String(201, "test %s %d", "string", 2)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "test string 2")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

// Tests that no String is rendered if code is 204
func TestContextRenderNoContentString(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.String(204, "test %s %d", "string", 2)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

// TestContextString tests that the response is returned
// with Content-Type set to text/html
func TestContextRenderHTMLString(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(201, "<html>%s %d</html>", "string", 3)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "<html>string 3</html>")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// Tests that no HTML String is rendered if code is 204
func TestContextRenderNoContentHTMLString(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(204, "<html>%s %d</html>", "string", 3)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html; charset=utf-8")
}

// TestContextData tests that the response can be written from `bytesting`
// with specified MIME type
func TestContextRenderData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Data(201, "text/csv", []byte(`foo,bar`))

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "foo,bar")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/csv")
}

// Tests that no Custom Data is rendered if code is 204
func TestContextRenderNoContentData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Data(204, "text/csv", []byte(`foo,bar`))

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/csv")
}

func TestContextRenderSSE(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.SSEvent("float", 1.5)
	c.Render(-1, sse.Event{
		Id:   "123",
		Data: "text",
	})
	c.SSEvent("chat", H{
		"foo": "bar",
		"bar": "foo",
	})

	assert.Equal(t, strings.Replace(w.Body.String(), " ", "", -1), strings.Replace("event:float\ndata:1.5\n\nid:123\ndata:text\n\nevent:chat\ndata:{\"bar\":\"foo\",\"foo\":\"bar\"}\n\n", " ", "", -1))
}

func TestContextRenderFile(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.File("./gin.go")

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "func New() *Engine {")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
}

// TestContextRenderYAML tests that the response is serialized as YAML
// and Content-Type is set to application/x-yaml
func TestContextRenderYAML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.YAML(201, H{"foo": "bar"})

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), "foo: bar\n")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/x-yaml; charset=utf-8")
}

func TestContextHeaders(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	assert.Panics(t, func() { c.Redirect(299, "/new_path") })
	assert.Panics(t, func() { c.Redirect(309, "/new_path") })

	c.Redirect(301, "/path")
	c.Writer.WriteHeaderNow()
	assert.Equal(t, w.Code, 301)
	assert.Equal(t, w.Header().Get("Location"), "/path")
}

func TestContextRenderRedirectWithAbsolutePath(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	c.Redirect(302, "http://google.com")
	c.Writer.WriteHeaderNow()

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.Header().Get("Location"), "http://google.com")
}

func TestContextRenderRedirectWith201(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	c.Redirect(201, "/resource")
	c.Writer.WriteHeaderNow()

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Header().Get("Location"), "/resource")
}

func TestContextRenderRedirectAll(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	assert.Panics(t, func() { c.Redirect(200, "/resource") })
	assert.Panics(t, func() { c.Redirect(202, "/resource") })
	assert.Panics(t, func() { c.Redirect(299, "/resource") })
	assert.Panics(t, func() { c.Redirect(309, "/resource") })
	assert.NotPanics(t, func() { c.Redirect(300, "/resource") })
	assert.NotPanics(t, func() { c.Redirect(308, "/resource") })
}

func TestContextNegotiationWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "", nil)

	c.Negotiate(200, Negotiate{
		Offered: []string{MIMEJSON, MIMEXML},
		Data:    H{"foo": "bar"},
	})

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"foo\":\"bar\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

func TestContextNegotiationWithXML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "", nil)

	c.Negotiate(200, Negotiate{
		Offered: []string{MIMEXML, MIMEJSON},
		Data:    H{"foo": "bar"},
	})

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<map><foo>bar</foo></map>", w.Body.String())
	assert.Equal(t, "application/xml; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

func TestContextNegotiationWithHTML(t *testing.T) {
	w := httptest.NewRecorder()
	c, router := CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "", nil)
	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
	router.SetHTMLTemplate(templ)

	c.Negotiate(200, Negotiate{
		Offered:  []string{MIMEHTML},
		Data:     H{"name": "gin"},
		HTMLName: "t",
	})

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Hello gin", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.HeaderMap.Get("Content-Type"))
}

func TestContextNegotiationNotSupport(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "", nil)

	c.Negotiate(200, Negotiate{
		Offered: []string{MIMEPOSTForm},
	})

	assert.Equal(t, 406, w.Code)
	assert.Equal(t, c.index, abortIndex)
	assert.True(t, c.IsAborted())
}

func TestContextNegotiationFormat(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "", nil)

	assert.Panics(t, func() { c.NegotiateFormat() })
	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEJSON)
	assert.Equal(t, c.NegotiateFormat(MIMEHTML, MIMEJSON), MIMEHTML)
}

func TestContextNegotiationFormatWithAccept(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEHTML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), "")
}

func TestContextNegotiationFormatCustum(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	c.Accepted = nil
	c.SetAccepted(MIMEJSON, MIMEXML)

	assert.Equal(t, c.NegotiateFormat(MIMEJSON, MIMEXML), MIMEJSON)
	assert.Equal(t, c.NegotiateFormat(MIMEXML, MIMEHTML), MIMEXML)
	assert.Equal(t, c.NegotiateFormat(MIMEJSON), MIMEJSON)
}

func TestContextIsAborted(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.index = 4
	c.AbortWithStatus(401)

	assert.Equal(t, c.index, abortIndex)
	assert.Equal(t, c.Writer.Status(), 401)
	assert.Equal(t, w.Code, 401)
	assert.True(t, c.IsAborted())
}

type testJSONAbortMsg struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func TestContextAbortWithStatusJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.index = 4

	in := new(testJSONAbortMsg)
	in.Bar = "barValue"
	in.Foo = "fooValue"

	c.AbortWithStatusJSON(415, in)

	assert.Equal(t, c.index, abortIndex)
	assert.Equal(t, c.Writer.Status(), 415)
	assert.Equal(t, w.Code, 415)
	assert.True(t, c.IsAborted())

	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, contentType, "application/json; charset=utf-8")

	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Body)
	jsonStringBody := buf.String()
	assert.Equal(t, fmt.Sprint(`{"foo":"fooValue","bar":"barValue"}`), jsonStringBody)
}

func TestContextError(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.AbortWithError(401, errors.New("bad input")).SetMeta("some input")

	assert.Equal(t, w.Code, 401)
	assert.Equal(t, c.index, abortIndex)
	assert.True(t, c.IsAborted())
}

func TestContextClientIP(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)

	c.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	c.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	c.Request.Header.Set("X-Appengine-Remote-Addr", "50.50.50.50")
	c.Request.RemoteAddr = "  40.40.40.40:42123 "

	assert.Equal(t, "20.20.20.20", c.ClientIP())

	c.Request.Header.Del("X-Forwarded-For")
	assert.Equal(t, "10.10.10.10", c.ClientIP())

	c.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assert.Equal(t, "30.30.30.30", c.ClientIP())

	c.Request.Header.Del("X-Forwarded-For")
	c.Request.Header.Del("X-Real-IP")
	c.engine.AppEngine = true
	assert.Equal(t, "50.50.50.50", c.ClientIP())

	c.Request.Header.Del("X-Appengine-Remote-Addr")
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// no port
	c.Request.RemoteAddr = "50.50.50.50"
	assert.Equal(t, "", c.ClientIP())
}

func TestContextContentType(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, c.ContentType(), "application/json")
}

func TestContextAutoBindJSON(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
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
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

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
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

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
	c, _ := CreateTestContext(httptest.NewRecorder())
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

func TestWebsocketsRequired(t *testing.T) {
	// Example request from spec: https://tools.ietf.org/html/rfc6455#section-1.2
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/chat", nil)
	c.Request.Header.Set("Host", "server.example.com")
	c.Request.Header.Set("Upgrade", "websocket")
	c.Request.Header.Set("Connection", "Upgrade")
	c.Request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	c.Request.Header.Set("Origin", "http://example.com")
	c.Request.Header.Set("Sec-WebSocket-Protocol", "chat, superchat")
	c.Request.Header.Set("Sec-WebSocket-Version", "13")

	assert.True(t, c.IsWebsocket())

	// Normal request, no websocket required.
	c, _ = CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/chat", nil)
	c.Request.Header.Set("Host", "server.example.com")

	assert.False(t, c.IsWebsocket())
}

func TestGetRequestHeaderValue(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/chat", nil)
	c.Request.Header.Set("Gin-Version", "1.0.0")

	assert.Equal(t, "1.0.0", c.GetHeader("Gin-Version"))
	assert.Equal(t, "", c.GetHeader("Connection"))
}

func TestContextGetRawData(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString("Fetch binary post data")
	c.Request, _ = http.NewRequest("POST", "/", body)
	c.Request.Header.Add("Content-Type", MIMEPOSTForm)

	data, err := c.GetRawData()
	assert.Nil(t, err)
	assert.Equal(t, "Fetch binary post data", string(data))
}
