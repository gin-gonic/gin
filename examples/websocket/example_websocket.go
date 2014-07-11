package main

import (
	"errors"
	"fmt"
	"strconv"

	"code.google.com/p/go.net/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		handler := websocket.Handler(func(conn *websocket.Conn) {
			// This is the simplest possible echoserver
			//
			// For a more advanced handler you can use the
			// conn.Read and conn.Write methods as the
			// websocket.Conn type conforms to io.Reader+io.Writer

			io.Copy(conn, conn)
		})
		handler.ServeHTTP(&c.Writer, c.Req)
	})

	go r.Run(":8080")

	lock := make(chan bool)
	go testServer(100, lock)
	<-lock
}

func testServer(count int, done chan bool) {
	client, err := websocket.Dial("ws://localhost:8080", "", "http://localhost/")
	if err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		out := []byte(strconv.Itoa(i))
		_, err = client.Write(out)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sent: %s\n", out)

		var in = make([]byte, 512)
		_, err = client.Read(in)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n\n", in)
	}

	done <- true
}
