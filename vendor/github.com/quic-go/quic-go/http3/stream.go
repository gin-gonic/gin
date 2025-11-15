package http3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3/qlog"
	"github.com/quic-go/quic-go/qlogwriter"

	"github.com/quic-go/qpack"
)

type datagramStream interface {
	io.ReadWriteCloser
	CancelRead(quic.StreamErrorCode)
	CancelWrite(quic.StreamErrorCode)
	StreamID() quic.StreamID
	Context() context.Context
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	SendDatagram(b []byte) error
	ReceiveDatagram(ctx context.Context) ([]byte, error)

	QUICStream() *quic.Stream
}

// A Stream is an HTTP/3 stream.
//
// When writing to and reading from the stream, data is framed in HTTP/3 DATA frames.
type Stream struct {
	datagramStream
	conn        *Conn
	frameParser *frameParser

	buf []byte // used as a temporary buffer when writing the HTTP/3 frame headers

	bytesRemainingInFrame uint64

	qlogger qlogwriter.Recorder

	parseTrailer  func(io.Reader, *headersFrame) error
	parsedTrailer bool
}

func newStream(
	str datagramStream,
	conn *Conn,
	trace *httptrace.ClientTrace,
	parseTrailer func(io.Reader, *headersFrame) error,
	qlogger qlogwriter.Recorder,
) *Stream {
	return &Stream{
		datagramStream: str,
		conn:           conn,
		buf:            make([]byte, 16),
		qlogger:        qlogger,
		parseTrailer:   parseTrailer,
		frameParser: &frameParser{
			r:         &tracingReader{Reader: str, trace: trace},
			streamID:  str.StreamID(),
			closeConn: conn.CloseWithError,
		},
	}
}

func (s *Stream) Read(b []byte) (int, error) {
	if s.bytesRemainingInFrame == 0 {
	parseLoop:
		for {
			frame, err := s.frameParser.ParseNext(s.qlogger)
			if err != nil {
				return 0, err
			}
			switch f := frame.(type) {
			case *dataFrame:
				if s.parsedTrailer {
					return 0, errors.New("DATA frame received after trailers")
				}
				s.bytesRemainingInFrame = f.Length
				break parseLoop
			case *headersFrame:
				if s.conn.isServer {
					continue
				}
				if s.parsedTrailer {
					maybeQlogInvalidHeadersFrame(s.qlogger, s.StreamID(), f.Length)
					return 0, errors.New("additional HEADERS frame received after trailers")
				}
				s.parsedTrailer = true
				return 0, s.parseTrailer(s.datagramStream, f)
			default:
				s.conn.CloseWithError(quic.ApplicationErrorCode(ErrCodeFrameUnexpected), "")
				// parseNextFrame skips over unknown frame types
				// Therefore, this condition is only entered when we parsed another known frame type.
				return 0, fmt.Errorf("peer sent an unexpected frame: %T", f)
			}
		}
	}

	var n int
	var err error
	if s.bytesRemainingInFrame < uint64(len(b)) {
		n, err = s.datagramStream.Read(b[:s.bytesRemainingInFrame])
	} else {
		n, err = s.datagramStream.Read(b)
	}
	s.bytesRemainingInFrame -= uint64(n)
	return n, err
}

func (s *Stream) hasMoreData() bool {
	return s.bytesRemainingInFrame > 0
}

func (s *Stream) Write(b []byte) (int, error) {
	s.buf = s.buf[:0]
	s.buf = (&dataFrame{Length: uint64(len(b))}).Append(s.buf)
	if s.qlogger != nil {
		s.qlogger.RecordEvent(qlog.FrameCreated{
			StreamID: s.StreamID(),
			Raw: qlog.RawInfo{
				Length:        len(s.buf) + len(b),
				PayloadLength: len(b),
			},
			Frame: qlog.Frame{Frame: qlog.DataFrame{}},
		})
	}
	if _, err := s.datagramStream.Write(s.buf); err != nil {
		return 0, err
	}
	return s.datagramStream.Write(b)
}

func (s *Stream) writeUnframed(b []byte) (int, error) {
	return s.datagramStream.Write(b)
}

func (s *Stream) StreamID() quic.StreamID {
	return s.datagramStream.StreamID()
}

func (s *Stream) SendDatagram(b []byte) error {
	// TODO: reject if datagrams are not negotiated (yet)
	return s.datagramStream.SendDatagram(b)
}

func (s *Stream) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	// TODO: reject if datagrams are not negotiated (yet)
	return s.datagramStream.ReceiveDatagram(ctx)
}

// A RequestStream is a low-level abstraction representing an HTTP/3 request stream.
// It decouples sending of the HTTP request from reading the HTTP response, allowing
// the application to optimistically use the stream (and, for example, send datagrams)
// before receiving the response.
//
// This is only needed for advanced use case, e.g. WebTransport and the various
// MASQUE proxying protocols.
type RequestStream struct {
	str *Stream

	responseBody io.ReadCloser // set by ReadResponse

	decoder            *qpack.Decoder
	requestWriter      *requestWriter
	maxHeaderBytes     uint64
	reqDone            chan<- struct{}
	disableCompression bool
	response           *http.Response

	sentRequest   bool
	requestedGzip bool
	isConnect     bool
}

func newRequestStream(
	str *Stream,
	requestWriter *requestWriter,
	reqDone chan<- struct{},
	decoder *qpack.Decoder,
	disableCompression bool,
	maxHeaderBytes uint64,
	rsp *http.Response,
) *RequestStream {
	return &RequestStream{
		str:                str,
		requestWriter:      requestWriter,
		reqDone:            reqDone,
		decoder:            decoder,
		disableCompression: disableCompression,
		maxHeaderBytes:     maxHeaderBytes,
		response:           rsp,
	}
}

// Read reads data from the underlying stream.
//
// It can only be used after the request has been sent (using SendRequestHeader)
// and the response has been consumed (using ReadResponse).
func (s *RequestStream) Read(b []byte) (int, error) {
	if s.responseBody == nil {
		return 0, errors.New("http3: invalid use of RequestStream.Read before ReadResponse")
	}
	return s.responseBody.Read(b)
}

// StreamID returns the QUIC stream ID of the underlying QUIC stream.
func (s *RequestStream) StreamID() quic.StreamID {
	return s.str.StreamID()
}

// Write writes data to the stream.
//
// It can only be used after the request has been sent (using SendRequestHeader).
func (s *RequestStream) Write(b []byte) (int, error) {
	if !s.sentRequest {
		return 0, errors.New("http3: invalid use of RequestStream.Write before SendRequestHeader")
	}
	return s.str.Write(b)
}

// Close closes the send-direction of the stream.
// It does not close the receive-direction of the stream.
func (s *RequestStream) Close() error {
	return s.str.Close()
}

// CancelRead aborts receiving on this stream.
// See [quic.Stream.CancelRead] for more details.
func (s *RequestStream) CancelRead(errorCode quic.StreamErrorCode) {
	s.str.CancelRead(errorCode)
}

// CancelWrite aborts sending on this stream.
// See [quic.Stream.CancelWrite] for more details.
func (s *RequestStream) CancelWrite(errorCode quic.StreamErrorCode) {
	s.str.CancelWrite(errorCode)
}

// Context returns a context derived from the underlying QUIC stream's context.
// See [quic.Stream.Context] for more details.
func (s *RequestStream) Context() context.Context {
	return s.str.Context()
}

// SetReadDeadline sets the deadline for Read calls.
func (s *RequestStream) SetReadDeadline(t time.Time) error {
	return s.str.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for Write calls.
func (s *RequestStream) SetWriteDeadline(t time.Time) error {
	return s.str.SetWriteDeadline(t)
}

// SetDeadline sets the read and write deadlines associated with the stream.
// It is equivalent to calling both SetReadDeadline and SetWriteDeadline.
func (s *RequestStream) SetDeadline(t time.Time) error {
	return s.str.SetDeadline(t)
}

// SendDatagrams send a new HTTP Datagram (RFC 9297).
//
// It is only possible to send datagrams if the server enabled support for this extension.
// It is recommended (though not required) to send the request before calling this method,
// as the server might drop datagrams which it can't associate with an existing request.
func (s *RequestStream) SendDatagram(b []byte) error {
	return s.str.SendDatagram(b)
}

// ReceiveDatagram receives HTTP Datagrams (RFC 9297).
//
// It is only possible if support for HTTP Datagrams was enabled, using the EnableDatagram
// option on the [Transport].
func (s *RequestStream) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return s.str.ReceiveDatagram(ctx)
}

// SendRequestHeader sends the HTTP request.
//
// It can only used for requests that don't have a request body.
// It is invalid to call it more than once.
// It is invalid to call it after Write has been called.
func (s *RequestStream) SendRequestHeader(req *http.Request) error {
	if req.Body != nil && req.Body != http.NoBody {
		return errors.New("http3: invalid use of RequestStream.SendRequestHeader with a request that has a request body")
	}
	return s.sendRequestHeader(req)
}

func (s *RequestStream) sendRequestHeader(req *http.Request) error {
	if s.sentRequest {
		return errors.New("http3: invalid duplicate use of RequestStream.SendRequestHeader")
	}
	if !s.disableCompression && req.Method != http.MethodHead &&
		req.Header.Get("Accept-Encoding") == "" && req.Header.Get("Range") == "" {
		s.requestedGzip = true
	}
	s.isConnect = req.Method == http.MethodConnect
	s.sentRequest = true
	return s.requestWriter.WriteRequestHeader(s.str.datagramStream, req, s.requestedGzip, s.str.StreamID(), s.str.qlogger)
}

// ReadResponse reads the HTTP response from the stream.
//
// It must be called after sending the request (using SendRequestHeader).
// It is invalid to call it more than once.
// It doesn't set Response.Request and Response.TLS.
// It is invalid to call it after Read has been called.
func (s *RequestStream) ReadResponse() (*http.Response, error) {
	if !s.sentRequest {
		return nil, errors.New("http3: invalid duplicate use of RequestStream.ReadResponse before SendRequestHeader")
	}
	frame, err := s.str.frameParser.ParseNext(s.str.qlogger)
	if err != nil {
		s.str.CancelRead(quic.StreamErrorCode(ErrCodeFrameError))
		s.str.CancelWrite(quic.StreamErrorCode(ErrCodeFrameError))
		return nil, fmt.Errorf("http3: parsing frame failed: %w", err)
	}
	hf, ok := frame.(*headersFrame)
	if !ok {
		s.str.conn.CloseWithError(quic.ApplicationErrorCode(ErrCodeFrameUnexpected), "expected first frame to be a HEADERS frame")
		return nil, errors.New("http3: expected first frame to be a HEADERS frame")
	}
	if hf.Length > s.maxHeaderBytes {
		maybeQlogInvalidHeadersFrame(s.str.qlogger, s.str.StreamID(), hf.Length)
		s.str.CancelRead(quic.StreamErrorCode(ErrCodeFrameError))
		s.str.CancelWrite(quic.StreamErrorCode(ErrCodeFrameError))
		return nil, fmt.Errorf("http3: HEADERS frame too large: %d bytes (max: %d)", hf.Length, s.maxHeaderBytes)
	}
	headerBlock := make([]byte, hf.Length)
	if _, err := io.ReadFull(s.str.datagramStream, headerBlock); err != nil {
		maybeQlogInvalidHeadersFrame(s.str.qlogger, s.str.StreamID(), hf.Length)
		s.str.CancelRead(quic.StreamErrorCode(ErrCodeRequestIncomplete))
		s.str.CancelWrite(quic.StreamErrorCode(ErrCodeRequestIncomplete))
		return nil, fmt.Errorf("http3: failed to read response headers: %w", err)
	}
	hfs, err := s.decoder.DecodeFull(headerBlock)
	if err != nil {
		maybeQlogInvalidHeadersFrame(s.str.qlogger, s.str.StreamID(), hf.Length)
		// TODO: use the right error code
		s.str.conn.CloseWithError(quic.ApplicationErrorCode(ErrCodeGeneralProtocolError), "")
		return nil, fmt.Errorf("http3: failed to decode response headers: %w", err)
	}
	if s.str.qlogger != nil {
		qlogParsedHeadersFrame(s.str.qlogger, s.str.StreamID(), hf, hfs)
	}
	res := s.response
	if err := updateResponseFromHeaders(res, hfs); err != nil {
		s.str.CancelRead(quic.StreamErrorCode(ErrCodeMessageError))
		s.str.CancelWrite(quic.StreamErrorCode(ErrCodeMessageError))
		return nil, fmt.Errorf("http3: invalid response: %w", err)
	}

	// Check that the server doesn't send more data in DATA frames than indicated by the Content-Length header (if set).
	// See section 4.1.2 of RFC 9114.
	respBody := newResponseBody(s.str, res.ContentLength, s.reqDone)

	// Rules for when to set Content-Length are defined in https://tools.ietf.org/html/rfc7230#section-3.3.2.
	isInformational := res.StatusCode >= 100 && res.StatusCode < 200
	isNoContent := res.StatusCode == http.StatusNoContent
	isSuccessfulConnect := s.isConnect && res.StatusCode >= 200 && res.StatusCode < 300
	if (isInformational || isNoContent || isSuccessfulConnect) && res.ContentLength == -1 {
		res.ContentLength = 0
	}
	if s.requestedGzip && res.Header.Get("Content-Encoding") == "gzip" {
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		s.responseBody = newGzipReader(respBody)
		res.Uncompressed = true
	} else {
		s.responseBody = respBody
	}
	res.Body = s.responseBody
	return res, nil
}

type tracingReader struct {
	io.Reader
	readFirst bool
	trace     *httptrace.ClientTrace
}

func (r *tracingReader) Read(b []byte) (int, error) {
	n, err := r.Reader.Read(b)
	if n > 0 && !r.readFirst {
		traceGotFirstResponseByte(r.trace)
		r.readFirst = true
	}
	return n, err
}
