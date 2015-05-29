package gin

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	if err != nil {
		println(err)
		t.FailNow()
	}
	fmt.Fprintf(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}
