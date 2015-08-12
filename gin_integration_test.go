package gin

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	go func() {
		router.Use(LoggerWithWriter(buffer))
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		router.Run(":5150")
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	assert.Error(t, router.Run(":5150"))

	resp, err := http.Get("http://localhost:5150/example")
	defer resp.Body.Close()
	assert.NoError(t, err)

	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "it worked", string(body[:]), "resp body should match")
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
}

func TestUnixSocket(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()

	go func() {
		router.Use(LoggerWithWriter(buffer))
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		router.RunUnix("/tmp/unix_unit_test")
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
