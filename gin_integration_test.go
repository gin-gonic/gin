// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// params[0]=url example:http://127.0.0.1:8080/index (cannot be empty)
// params[1]=response status (custom compare status) default:"200 OK"
// params[2]=response body (custom compare content)  default:"it worked"
func testRequest(t *testing.T, params ...string) {

	if len(params) == 0 {
		t.Fatal("url cannot be empty")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(params[0])
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)

	var responseStatus = "200 OK"
	if len(params) > 1 && params[1] != "" {
		responseStatus = params[1]
	}

	var responseBody = "it worked"
	if len(params) > 2 && params[2] != "" {
		responseBody = params[2]
	}

	assert.Equal(t, responseStatus, resp.Status, "should get a "+responseStatus)
	if responseStatus == "200 OK" {
		assert.Equal(t, responseBody, string(body), "resp body should match")
	}
}

func TestRunEmpty(t *testing.T) {
	os.Setenv("PORT", "")
	router := New()
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.Run())
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.Run(":8080"))
	testRequest(t, "http://localhost:8080/example")
}

func TestTrustedCIDRsForRun(t *testing.T) {
	os.Setenv("PORT", "")
	router := New()
	router.TrustedProxies = []string{"hello/world"}
	assert.Error(t, router.Run(":8080"))
}

func TestRunTLS(t *testing.T) {
	router := New()
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

		assert.NoError(t, router.RunTLS(":8443", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	}()

	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.RunTLS(":8443", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	testRequest(t, "https://localhost:8443/example")
}

func TestPusher(t *testing.T) {
	var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

	router := New()
	router.Static("./assets", "./assets")
	router.SetHTMLTemplate(html)

	go func() {
		router.GET("/pusher", func(c *Context) {
			if pusher := c.Writer.Pusher(); pusher != nil {
				err := pusher.Push("/assets/app.js", nil)
				assert.NoError(t, err)
			}
			c.String(http.StatusOK, "it worked")
		})

		assert.NoError(t, router.RunTLS(":8449", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	}()

	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.RunTLS(":8449", "./testdata/certificate/cert.pem", "./testdata/certificate/key.pem"))
	testRequest(t, "https://localhost:8449/pusher")
}

func TestRunEmptyWithEnv(t *testing.T) {
	os.Setenv("PORT", "3123")
	router := New()
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.Run())
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.Run(":3123"))
	testRequest(t, "http://localhost:3123/example")
}

func TestRunTooMuchParams(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		assert.NoError(t, router.Run("2", "2"))
	})
}

func TestRunWithPort(t *testing.T) {
	router := New()
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.Run(":5150"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.Run(":5150"))
	testRequest(t, "http://localhost:5150/example")
}

func TestUnixSocket(t *testing.T) {
	router := New()

	unixTestSocket := filepath.Join(os.TempDir(), "unix_unit_test")

	defer os.Remove(unixTestSocket)

	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.RunUnix(unixTestSocket))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("unix", unixTestSocket)
	assert.NoError(t, err)

	fmt.Fprint(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}

func TestBadUnixSocket(t *testing.T) {
	router := New()
	assert.Error(t, router.RunUnix("#/tmp/unix_unit_test"))
}

func TestFileDescriptor(t *testing.T) {
	router := New()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	socketFile, err := listener.File()
	assert.NoError(t, err)

	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.RunFd(int(socketFile.Fd())))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	fmt.Fprintf(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}

func TestBadFileDescriptor(t *testing.T) {
	router := New()
	assert.Error(t, router.RunFd(0))
}

func TestListener(t *testing.T) {
	router := New()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.RunListener(listener))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	fmt.Fprintf(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}

func TestBadListener(t *testing.T) {
	router := New()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:10086")
	assert.NoError(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	listener.Close()
	assert.Error(t, router.RunListener(listener))
}

func TestWithHttptestWithAutoSelectedPort(t *testing.T) {
	router := New()
	router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

	ts := httptest.NewServer(router)
	defer ts.Close()

	testRequest(t, ts.URL+"/example")
}

func TestConcurrentHandleContext(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.Request.URL.Path = "/example"
		router.HandleContext(c)
	})
	router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

	var wg sync.WaitGroup
	iterations := 200
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			testGetRequestHandler(t, router, "/")
			wg.Done()
		}()
	}
	wg.Wait()
}

// func TestWithHttptestWithSpecifiedPort(t *testing.T) {
// 	router := New()
// 	router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

// 	l, _ := net.Listen("tcp", ":8033")
// 	ts := httptest.Server{
// 		Listener: l,
// 		Config:   &http.Server{Handler: router},
// 	}
// 	ts.Start()
// 	defer ts.Close()

// 	testRequest(t, "http://localhost:8033/example")
// }

func testGetRequestHandler(t *testing.T, h http.Handler, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	assert.Equal(t, "it worked", w.Body.String(), "resp body should match")
	assert.Equal(t, 200, w.Code, "should get a 200")
}

func TestTreeRunDynamicRouting(t *testing.T) {
	router := New()
	router.GET("/aa/*xx", func(c *Context) { c.String(http.StatusOK, "/aa/*xx") })
	router.GET("/ab/*xx", func(c *Context) { c.String(http.StatusOK, "/ab/*xx") })
	router.GET("/", func(c *Context) { c.String(http.StatusOK, "home") })
	router.GET("/:cc", func(c *Context) { c.String(http.StatusOK, "/:cc") })
	router.GET("/:cc/cc", func(c *Context) { c.String(http.StatusOK, "/:cc/cc") })
	router.GET("/:cc/:dd/ee", func(c *Context) { c.String(http.StatusOK, "/:cc/:dd/ee") })
	router.GET("/:cc/:dd/:ee/ff", func(c *Context) { c.String(http.StatusOK, "/:cc/:dd/:ee/ff") })
	router.GET("/:cc/:dd/:ee/:ff/gg", func(c *Context) { c.String(http.StatusOK, "/:cc/:dd/:ee/:ff/gg") })
	router.GET("/:cc/:dd/:ee/:ff/:gg/hh", func(c *Context) { c.String(http.StatusOK, "/:cc/:dd/:ee/:ff/:gg/hh") })
	router.GET("/get/test/abc/", func(c *Context) { c.String(http.StatusOK, "/get/test/abc/") })
	router.GET("/get/:param/abc/", func(c *Context) { c.String(http.StatusOK, "/get/:param/abc/") })
	router.GET("/something/:paramname/thirdthing", func(c *Context) { c.String(http.StatusOK, "/something/:paramname/thirdthing") })
	router.GET("/something/secondthing/test", func(c *Context) { c.String(http.StatusOK, "/something/secondthing/test") })
	router.GET("/get/abc", func(c *Context) { c.String(http.StatusOK, "/get/abc") })
	router.GET("/get/:param", func(c *Context) { c.String(http.StatusOK, "/get/:param") })
	router.GET("/get/abc/123abc", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc") })
	router.GET("/get/abc/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/:param") })
	router.GET("/get/abc/123abc/xxx8", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8") })
	router.GET("/get/abc/123abc/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/:param") })
	router.GET("/get/abc/123abc/xxx8/1234", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/1234") })
	router.GET("/get/abc/123abc/xxx8/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/:param") })
	router.GET("/get/abc/123abc/xxx8/1234/ffas", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/1234/ffas") })
	router.GET("/get/abc/123abc/xxx8/1234/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/1234/:param") })
	router.GET("/get/abc/123abc/xxx8/1234/kkdd/12c", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/1234/kkdd/12c") })
	router.GET("/get/abc/123abc/xxx8/1234/kkdd/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abc/xxx8/1234/kkdd/:param") })
	router.GET("/get/abc/:param/test", func(c *Context) { c.String(http.StatusOK, "/get/abc/:param/test") })
	router.GET("/get/abc/123abd/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abd/:param") })
	router.GET("/get/abc/123abddd/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abddd/:param") })
	router.GET("/get/abc/123/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123/:param") })
	router.GET("/get/abc/123abg/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abg/:param") })
	router.GET("/get/abc/123abf/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abf/:param") })
	router.GET("/get/abc/123abfff/:param", func(c *Context) { c.String(http.StatusOK, "/get/abc/123abfff/:param") })

	ts := httptest.NewServer(router)
	defer ts.Close()

	testRequest(t, ts.URL+"/", "", "home")
	testRequest(t, ts.URL+"/aa/aa", "", "/aa/*xx")
	testRequest(t, ts.URL+"/ab/ab", "", "/ab/*xx")
	testRequest(t, ts.URL+"/all", "", "/:cc")
	testRequest(t, ts.URL+"/all/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/a/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/c/d/ee", "", "/:cc/:dd/ee")
	testRequest(t, ts.URL+"/c/d/e/ff", "", "/:cc/:dd/:ee/ff")
	testRequest(t, ts.URL+"/c/d/e/f/gg", "", "/:cc/:dd/:ee/:ff/gg")
	testRequest(t, ts.URL+"/c/d/e/f/g/hh", "", "/:cc/:dd/:ee/:ff/:gg/hh")
	testRequest(t, ts.URL+"/cc/dd/ee/ff/gg/hh", "", "/:cc/:dd/:ee/:ff/:gg/hh")
	testRequest(t, ts.URL+"/a", "", "/:cc")
	testRequest(t, ts.URL+"/d", "", "/:cc")
	testRequest(t, ts.URL+"/ad", "", "/:cc")
	testRequest(t, ts.URL+"/dd", "", "/:cc")
	testRequest(t, ts.URL+"/aa", "", "/:cc")
	testRequest(t, ts.URL+"/aaa", "", "/:cc")
	testRequest(t, ts.URL+"/aaa/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/ab", "", "/:cc")
	testRequest(t, ts.URL+"/abb", "", "/:cc")
	testRequest(t, ts.URL+"/abb/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/dddaa", "", "/:cc")
	testRequest(t, ts.URL+"/allxxxx", "", "/:cc")
	testRequest(t, ts.URL+"/alldd", "", "/:cc")
	testRequest(t, ts.URL+"/cc/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/ccc/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/deedwjfs/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/acllcc/cc", "", "/:cc/cc")
	testRequest(t, ts.URL+"/get/test/abc/", "", "/get/test/abc/")
	testRequest(t, ts.URL+"/get/testaa/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/te/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/xx/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/tt/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/a/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/t/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/aa/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/get/abas/abc/", "", "/get/:param/abc/")
	testRequest(t, ts.URL+"/something/secondthing/test", "", "/something/secondthing/test")
	testRequest(t, ts.URL+"/something/secondthingaaaa/thirdthing", "", "/something/:paramname/thirdthing")
	testRequest(t, ts.URL+"/something/abcdad/thirdthing", "", "/something/:paramname/thirdthing")
	testRequest(t, ts.URL+"/something/se/thirdthing", "", "/something/:paramname/thirdthing")
	testRequest(t, ts.URL+"/something/s/thirdthing", "", "/something/:paramname/thirdthing")
	testRequest(t, ts.URL+"/something/secondthing/thirdthing", "", "/something/:paramname/thirdthing")
	testRequest(t, ts.URL+"/get/abc", "", "/get/abc")
	testRequest(t, ts.URL+"/get/a", "", "/get/:param")
	testRequest(t, ts.URL+"/get/abz", "", "/get/:param")
	testRequest(t, ts.URL+"/get/12a", "", "/get/:param")
	testRequest(t, ts.URL+"/get/abcd", "", "/get/:param")
	testRequest(t, ts.URL+"/get/abc/123abc", "", "/get/abc/123abc")
	testRequest(t, ts.URL+"/get/abc/12", "", "/get/abc/:param")
	testRequest(t, ts.URL+"/get/abc/123ab", "", "/get/abc/:param")
	testRequest(t, ts.URL+"/get/abc/xyz", "", "/get/abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abcddxx", "", "/get/abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8", "", "/get/abc/123abc/xxx8")
	testRequest(t, ts.URL+"/get/abc/123abc/x", "", "/get/abc/123abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx", "", "/get/abc/123abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/abc", "", "/get/abc/123abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8xxas", "", "/get/abc/123abc/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234", "", "/get/abc/123abc/xxx8/1234")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1", "", "/get/abc/123abc/xxx8/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/123", "", "/get/abc/123abc/xxx8/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/78k", "", "/get/abc/123abc/xxx8/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234xxxd", "", "/get/abc/123abc/xxx8/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/ffas", "", "/get/abc/123abc/xxx8/1234/ffas")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/f", "", "/get/abc/123abc/xxx8/1234/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/ffa", "", "/get/abc/123abc/xxx8/1234/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kka", "", "/get/abc/123abc/xxx8/1234/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/ffas321", "", "/get/abc/123abc/xxx8/1234/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/12c", "", "/get/abc/123abc/xxx8/1234/kkdd/12c")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/1", "", "/get/abc/123abc/xxx8/1234/kkdd/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/12", "", "/get/abc/123abc/xxx8/1234/kkdd/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/12b", "", "/get/abc/123abc/xxx8/1234/kkdd/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/34", "", "/get/abc/123abc/xxx8/1234/kkdd/:param")
	testRequest(t, ts.URL+"/get/abc/123abc/xxx8/1234/kkdd/12c2e3", "", "/get/abc/123abc/xxx8/1234/kkdd/:param")
	testRequest(t, ts.URL+"/get/abc/12/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abdd/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abdddf/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123ab/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abgg/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abff/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abffff/test", "", "/get/abc/:param/test")
	testRequest(t, ts.URL+"/get/abc/123abd/test", "", "/get/abc/123abd/:param")
	testRequest(t, ts.URL+"/get/abc/123abddd/test", "", "/get/abc/123abddd/:param")
	testRequest(t, ts.URL+"/get/abc/123/test22", "", "/get/abc/123/:param")
	testRequest(t, ts.URL+"/get/abc/123abg/test", "", "/get/abc/123abg/:param")
	testRequest(t, ts.URL+"/get/abc/123abf/testss", "", "/get/abc/123abf/:param")
	testRequest(t, ts.URL+"/get/abc/123abfff/te", "", "/get/abc/123abfff/:param")
	// 404 not found
	testRequest(t, ts.URL+"/a/dd", "404 Not Found")
	testRequest(t, ts.URL+"/addr/dd/aa", "404 Not Found")
	testRequest(t, ts.URL+"/something/secondthing/121", "404 Not Found")
}
