package gin

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

const (
	StatusUnset int = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	// net/http.Response.Write only has two options: 200 or 500
	// we will follow that lead and defer to their logic

	// check if the write gave an error and set status accordingly
	size, err := w.ResponseWriter.Write(data)
	if err != nil {
		// error on write, we give a 500
		w.status = http.StatusInternalServerError
	} else if w.WasWritten() == false {
		// everything went okay and we never set a custom
		// status so 200 it is
		w.status = http.StatusOK
	}

	// can easily tap into Content-Length here with 'size'
	return size, err
}

// returns the status of the given response
func (w *ResponseWriter) Status() int {
	return w.status
}

// return a boolean acknowledging if a status code has all ready been set
func (w *ResponseWriter) WasWritten() bool {
	return w.status == StatusUnset
}

// allow connection hijacking
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
