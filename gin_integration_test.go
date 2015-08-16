package gin

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testRequest(t *testing.T, url string) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	assert.NoError(t, err)

	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "it worked", string(body), "resp body should match")
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
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
		router.Run("2", "2")
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

	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, router.RunUnix("/tmp/unix_unit_test"))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("unix", "/tmp/unix_unit_test")
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

func TestBadUnixSocket(t *testing.T) {
	router := New()
	assert.Error(t, router.RunUnix("#/tmp/unix_unit_test"))
}
