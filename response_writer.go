package gin

import (
	"net/http"
)

type (
	ResponseWriter interface {
		http.ResponseWriter
		Status() int
		Written() bool
		// Before allows for a function to be called before the ResponseWriter has been written to.
		Before(BeforeFunc)

		// private
		setStatus(int)
	}

	responseWriter struct {
		http.ResponseWriter
		status      int
		written     bool
		beforeFuncs []BeforeFunc
	}

	BeforeFunc func(ResponseWriter)
)

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.status = 0
	w.written = false
	w.beforeFuncs = w.beforeFuncs[:0]
}

func (w *responseWriter) setStatus(code int) {
	w.status = code
}

func (w *responseWriter) WriteHeader(code int) {
	w.callBefore()
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Written() bool {
	return w.written
}

func (w *responseWriter) Before(before BeforeFunc) {
	w.beforeFuncs = append(w.beforeFuncs, before)
}

func (w *responseWriter) callBefore() {
	for i := len(w.beforeFuncs) - 1; i >= 0; i-- {
		w.beforeFuncs[i](w)
	}
}
