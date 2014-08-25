package gin

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
)

type (
	ResponseWriter interface {
		http.ResponseWriter
		Status() int
		Written() bool
		WriteHeaderNow()
		Hijack() (net.Conn, *bufio.ReadWriter, error)
	}

	responseWriter struct {
		http.ResponseWriter
		status  int
		written bool
	}
)

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.status = 200
	w.written = false
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 {
		w.status = code
		if w.written {
			log.Println("[GIN] WARNING. Headers were already written!")
		}
	}
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.written {
		w.written = true
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Written() bool {
	return w.written
}

// allow connection hijacking
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
