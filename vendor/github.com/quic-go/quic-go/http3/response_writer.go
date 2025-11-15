package http3

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/quic-go/qpack"
	"github.com/quic-go/quic-go/http3/qlog"

	"golang.org/x/net/http/httpguts"
)

// The HTTPStreamer allows taking over a HTTP/3 stream. The interface is implemented by the http.ResponseWriter.
// When a stream is taken over, it's the caller's responsibility to close the stream.
type HTTPStreamer interface {
	HTTPStream() *Stream
}

// The maximum length of an encoded HTTP/3 frame header is 16:
// The frame has a type and length field, both QUIC varints (maximum 8 bytes in length)
const frameHeaderLen = 16

const maxSmallResponseSize = 4096

type responseWriter struct {
	str *Stream

	conn     *Conn
	header   http.Header
	trailers map[string]struct{}
	buf      []byte
	status   int // status code passed to WriteHeader

	// for responses smaller than maxSmallResponseSize, we buffer calls to Write,
	// and automatically add the Content-Length header
	smallResponseBuf []byte

	contentLen     int64 // if handler set valid Content-Length header
	numWritten     int64 // bytes written
	headerComplete bool  // set once WriteHeader is called with a status code >= 200
	headerWritten  bool  // set once the response header has been serialized to the stream
	isHead         bool
	trailerWritten bool // set once the response trailers has been serialized to the stream

	hijacked bool // set on HTTPStream is called

	logger *slog.Logger
}

var (
	_ http.ResponseWriter = &responseWriter{}
	_ http.Flusher        = &responseWriter{}
	_ Hijacker            = &responseWriter{}
	_ HTTPStreamer        = &responseWriter{}
	// make sure that we implement (some of the) methods used by the http.ResponseController
	_ interface {
		SetReadDeadline(time.Time) error
		SetWriteDeadline(time.Time) error
		Flush()
		FlushError() error
	} = &responseWriter{}
)

func newResponseWriter(str *Stream, conn *Conn, isHead bool, logger *slog.Logger) *responseWriter {
	return &responseWriter{
		str:    str,
		conn:   conn,
		header: http.Header{},
		buf:    make([]byte, frameHeaderLen),
		isHead: isHead,
		logger: logger,
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) WriteHeader(status int) {
	if w.headerComplete {
		return
	}

	// http status must be 3 digits
	if status < 100 || status > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", status))
	}
	w.status = status

	// immediately write 1xx headers
	if status < 200 {
		w.writeHeader(status)
		return
	}

	// We're done with headers once we write a status >= 200.
	w.headerComplete = true
	// Add Date header.
	// This is what the standard library does.
	// Can be disabled by setting the Date header to nil.
	if _, ok := w.header["Date"]; !ok {
		w.header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}
	// Content-Length checking
	// use ParseUint instead of ParseInt, as negative values are invalid
	if clen := w.header.Get("Content-Length"); clen != "" {
		if cl, err := strconv.ParseUint(clen, 10, 63); err == nil {
			w.contentLen = int64(cl)
		} else {
			// emit a warning for malformed Content-Length and remove it
			logger := w.logger
			if logger == nil {
				logger = slog.Default()
			}
			logger.Error("Malformed Content-Length", "value", clen)
			w.header.Del("Content-Length")
		}
	}
}

func (w *responseWriter) sniffContentType(p []byte) {
	// If no content type, apply sniffing algorithm to body.
	// We can't use `w.header.Get` here since if the Content-Type was set to nil, we shouldn't do sniffing.
	_, haveType := w.header["Content-Type"]

	// If the Content-Encoding was set and is non-blank, we shouldn't sniff the body.
	hasCE := w.header.Get("Content-Encoding") != ""
	if !hasCE && !haveType && len(p) > 0 {
		w.header.Set("Content-Type", http.DetectContentType(p))
	}
}

func (w *responseWriter) Write(p []byte) (int, error) {
	bodyAllowed := bodyAllowedForStatus(w.status)
	if !w.headerComplete {
		w.sniffContentType(p)
		w.WriteHeader(http.StatusOK)
		bodyAllowed = true
	}
	if !bodyAllowed {
		return 0, http.ErrBodyNotAllowed
	}

	w.numWritten += int64(len(p))
	if w.contentLen != 0 && w.numWritten > w.contentLen {
		return 0, http.ErrContentLength
	}

	if w.isHead {
		return len(p), nil
	}

	if !w.headerWritten {
		// Buffer small responses.
		// This allows us to automatically set the Content-Length field.
		if len(w.smallResponseBuf)+len(p) < maxSmallResponseSize {
			w.smallResponseBuf = append(w.smallResponseBuf, p...)
			return len(p), nil
		}
	}
	return w.doWrite(p)
}

func (w *responseWriter) doWrite(p []byte) (int, error) {
	if !w.headerWritten {
		w.sniffContentType(w.smallResponseBuf)
		if err := w.writeHeader(w.status); err != nil {
			return 0, maybeReplaceError(err)
		}
		w.headerWritten = true
	}

	l := uint64(len(w.smallResponseBuf) + len(p))
	if l == 0 {
		return 0, nil
	}
	df := &dataFrame{Length: l}
	w.buf = w.buf[:0]
	w.buf = df.Append(w.buf)
	if w.str.qlogger != nil {
		w.str.qlogger.RecordEvent(qlog.FrameCreated{
			StreamID: w.str.StreamID(),
			Raw:      qlog.RawInfo{Length: len(w.buf) + int(l), PayloadLength: int(l)},
			Frame:    qlog.Frame{Frame: qlog.DataFrame{}},
		})
	}
	if _, err := w.str.writeUnframed(w.buf); err != nil {
		return 0, maybeReplaceError(err)
	}
	if len(w.smallResponseBuf) > 0 {
		if _, err := w.str.writeUnframed(w.smallResponseBuf); err != nil {
			return 0, maybeReplaceError(err)
		}
		w.smallResponseBuf = nil
	}
	var n int
	if len(p) > 0 {
		var err error
		n, err = w.str.writeUnframed(p)
		if err != nil {
			return n, maybeReplaceError(err)
		}
	}
	return n, nil
}

func (w *responseWriter) writeHeader(status int) error {
	var headerFields []qlog.HeaderField // only used for qlog
	var headers bytes.Buffer
	enc := qpack.NewEncoder(&headers)
	if err := enc.WriteField(qpack.HeaderField{Name: ":status", Value: strconv.Itoa(status)}); err != nil {
		return err
	}
	if w.str.qlogger != nil {
		headerFields = append(headerFields, qlog.HeaderField{Name: ":status", Value: strconv.Itoa(status)})
	}

	// Handle trailer fields
	if vals, ok := w.header["Trailer"]; ok {
		for _, val := range vals {
			for _, trailer := range strings.Split(val, ",") {
				// We need to convert to the canonical header key value here because this will be called when using
				// headers.Add or headers.Set.
				trailer = textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(trailer))
				w.declareTrailer(trailer)
			}
		}
	}

	for k, v := range w.header {
		if _, excluded := w.trailers[k]; excluded {
			continue
		}
		// Ignore "Trailer:" prefixed headers
		if strings.HasPrefix(k, http.TrailerPrefix) {
			continue
		}
		for index := range v {
			name := strings.ToLower(k)
			value := v[index]
			if err := enc.WriteField(qpack.HeaderField{Name: name, Value: value}); err != nil {
				return err
			}
			if w.str.qlogger != nil {
				headerFields = append(headerFields, qlog.HeaderField{Name: name, Value: value})
			}
		}
	}

	buf := make([]byte, 0, frameHeaderLen+headers.Len())
	buf = (&headersFrame{Length: uint64(headers.Len())}).Append(buf)
	buf = append(buf, headers.Bytes()...)

	if w.str.qlogger != nil {
		qlogCreatedHeadersFrame(w.str.qlogger, w.str.StreamID(), len(buf), headers.Len(), headerFields)
	}

	_, err := w.str.writeUnframed(buf)
	return err
}

func (w *responseWriter) FlushError() error {
	if !w.headerComplete {
		w.WriteHeader(http.StatusOK)
	}
	_, err := w.doWrite(nil)
	return err
}

func (w *responseWriter) flushTrailers() {
	if w.trailerWritten {
		return
	}
	if err := w.writeTrailers(); err != nil {
		w.logger.Debug("could not write trailers", "error", err)
	}
}

func (w *responseWriter) Flush() {
	if err := w.FlushError(); err != nil {
		if w.logger != nil {
			w.logger.Debug("could not flush to stream", "error", err)
		}
	}
}

// declareTrailer adds a trailer to the trailer list, while also validating that the trailer has a
// valid name.
func (w *responseWriter) declareTrailer(k string) {
	if !httpguts.ValidTrailerHeader(k) {
		// Forbidden by RFC 9110, section 6.5.1.
		w.logger.Debug("ignoring invalid trailer", slog.String("header", k))
		return
	}
	if w.trailers == nil {
		w.trailers = make(map[string]struct{})
	}
	w.trailers[k] = struct{}{}
}

// hasNonEmptyTrailers checks to see if there are any trailers with an actual
// value set. This is possible by adding trailers to the "Trailers" header
// but never actually setting those names as trailers in the course of handling
// the request. In that case, this check may save us some allocations.
func (w *responseWriter) hasNonEmptyTrailers() bool {
	for trailer := range w.trailers {
		if _, ok := w.header[trailer]; ok {
			return true
		}
	}
	return false
}

// writeTrailers will write trailers to the stream if there are any.
func (w *responseWriter) writeTrailers() error {
	// promote headers added via "Trailer:" convention as trailers, these can be added after
	// streaming the status/headers have been written.
	for k := range w.header {
		// Handle "Trailer:" prefix
		if strings.HasPrefix(k, http.TrailerPrefix) {
			w.declareTrailer(k)
		}
	}

	if !w.hasNonEmptyTrailers() {
		return nil
	}

	var b bytes.Buffer
	var headerFields []qlog.HeaderField
	enc := qpack.NewEncoder(&b)
	for trailer := range w.trailers {
		trailerName := strings.ToLower(strings.TrimPrefix(trailer, http.TrailerPrefix))
		if vals, ok := w.header[trailer]; ok {
			for _, val := range vals {
				if err := enc.WriteField(qpack.HeaderField{Name: trailerName, Value: val}); err != nil {
					return err
				}
				if w.str.qlogger != nil {
					headerFields = append(headerFields, qlog.HeaderField{Name: trailerName, Value: val})
				}
			}
		}
	}

	buf := make([]byte, 0, frameHeaderLen+b.Len())
	buf = (&headersFrame{Length: uint64(b.Len())}).Append(buf)
	buf = append(buf, b.Bytes()...)
	if w.str.qlogger != nil {
		qlogCreatedHeadersFrame(w.str.qlogger, w.str.StreamID(), len(buf), b.Len(), headerFields)
	}
	_, err := w.str.writeUnframed(buf)
	w.trailerWritten = true
	return err
}

func (w *responseWriter) HTTPStream() *Stream {
	w.hijacked = true
	w.Flush()
	return w.str
}

func (w *responseWriter) wasStreamHijacked() bool { return w.hijacked }

func (w *responseWriter) Connection() *Conn {
	return w.conn
}

func (w *responseWriter) SetReadDeadline(deadline time.Time) error {
	return w.str.SetReadDeadline(deadline)
}

func (w *responseWriter) SetWriteDeadline(deadline time.Time) error {
	return w.str.SetWriteDeadline(deadline)
}

// copied from http2/http2.go
// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 2616, section 4.4.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
