package http3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3/qlog"
	"github.com/quic-go/quic-go/quicvarint"

	"github.com/quic-go/qpack"
)

const (
	// MethodGet0RTT allows a GET request to be sent using 0-RTT.
	// Note that 0-RTT doesn't provide replay protection and should only be used for idempotent requests.
	MethodGet0RTT = "GET_0RTT"
	// MethodHead0RTT allows a HEAD request to be sent using 0-RTT.
	// Note that 0-RTT doesn't provide replay protection and should only be used for idempotent requests.
	MethodHead0RTT = "HEAD_0RTT"
)

const (
	defaultUserAgent              = "quic-go HTTP/3"
	defaultMaxResponseHeaderBytes = 10 * 1 << 20 // 10 MB
)

type errConnUnusable struct{ e error }

func (e *errConnUnusable) Unwrap() error { return e.e }
func (e *errConnUnusable) Error() string { return fmt.Sprintf("http3: conn unusable: %s", e.e.Error()) }

const max1xxResponses = 5 // arbitrary bound on number of informational responses

var defaultQuicConfig = &quic.Config{
	MaxIncomingStreams: -1, // don't allow the server to create bidirectional streams
	KeepAlivePeriod:    10 * time.Second,
}

// ClientConn is an HTTP/3 client doing requests to a single remote server.
type ClientConn struct {
	conn *Conn

	// Enable support for HTTP/3 datagrams (RFC 9297).
	// If a QUICConfig is set, datagram support also needs to be enabled on the QUIC layer by setting enableDatagrams.
	enableDatagrams bool

	// Additional HTTP/3 settings.
	// It is invalid to specify any settings defined by RFC 9114 (HTTP/3) and RFC 9297 (HTTP Datagrams).
	additionalSettings map[uint64]uint64

	// maxResponseHeaderBytes specifies a limit on how many response bytes are
	// allowed in the server's response header.
	maxResponseHeaderBytes uint64

	// disableCompression, if true, prevents the Transport from requesting compression with an
	// "Accept-Encoding: gzip" request header when the Request contains no existing Accept-Encoding value.
	// If the Transport requests gzip on its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body.
	// However, if the user explicitly requested gzip it is not automatically uncompressed.
	disableCompression bool

	logger *slog.Logger

	requestWriter *requestWriter
	decoder       *qpack.Decoder
}

var _ http.RoundTripper = &ClientConn{}

func newClientConn(
	conn *quic.Conn,
	enableDatagrams bool,
	additionalSettings map[uint64]uint64,
	streamHijacker func(FrameType, quic.ConnectionTracingID, *quic.Stream, error) (hijacked bool, err error),
	uniStreamHijacker func(StreamType, quic.ConnectionTracingID, *quic.ReceiveStream, error) (hijacked bool),
	maxResponseHeaderBytes int64,
	disableCompression bool,
	logger *slog.Logger,
) *ClientConn {
	c := &ClientConn{
		enableDatagrams:    enableDatagrams,
		additionalSettings: additionalSettings,
		disableCompression: disableCompression,
		logger:             logger,
	}
	if maxResponseHeaderBytes <= 0 {
		c.maxResponseHeaderBytes = defaultMaxResponseHeaderBytes
	} else {
		c.maxResponseHeaderBytes = uint64(maxResponseHeaderBytes)
	}
	c.decoder = qpack.NewDecoder(func(hf qpack.HeaderField) {})
	c.requestWriter = newRequestWriter()
	c.conn = newConnection(
		conn.Context(),
		conn,
		c.enableDatagrams,
		false, // client
		c.logger,
		0,
	)
	// send the SETTINGs frame, using 0-RTT data, if possible
	go func() {
		if err := c.setupConn(); err != nil {
			if c.logger != nil {
				c.logger.Debug("Setting up connection failed", "error", err)
			}
			c.conn.CloseWithError(quic.ApplicationErrorCode(ErrCodeInternalError), "")
		}
	}()
	if streamHijacker != nil {
		go c.handleBidirectionalStreams(streamHijacker)
	}
	go c.conn.handleUnidirectionalStreams(uniStreamHijacker)
	return c
}

// OpenRequestStream opens a new request stream on the HTTP/3 connection.
func (c *ClientConn) OpenRequestStream(ctx context.Context) (*RequestStream, error) {
	return c.conn.openRequestStream(ctx, c.requestWriter, nil, c.disableCompression, c.maxResponseHeaderBytes)
}

func (c *ClientConn) setupConn() error {
	// open the control stream
	str, err := c.conn.OpenUniStream()
	if err != nil {
		return err
	}
	b := make([]byte, 0, 64)
	b = quicvarint.Append(b, streamTypeControlStream)
	// send the SETTINGS frame
	b = (&settingsFrame{Datagram: c.enableDatagrams, Other: c.additionalSettings}).Append(b)
	if c.conn.qlogger != nil {
		sf := qlog.SettingsFrame{
			Other: maps.Clone(c.additionalSettings),
		}
		if c.enableDatagrams {
			sf.Datagram = pointer(true)
		}
		c.conn.qlogger.RecordEvent(qlog.FrameCreated{
			StreamID: str.StreamID(),
			Raw:      qlog.RawInfo{Length: len(b)},
			Frame:    qlog.Frame{Frame: sf},
		})
	}
	_, err = str.Write(b)
	return err
}

func (c *ClientConn) handleBidirectionalStreams(streamHijacker func(FrameType, quic.ConnectionTracingID, *quic.Stream, error) (hijacked bool, err error)) {
	for {
		str, err := c.conn.conn.AcceptStream(context.Background())
		if err != nil {
			if c.logger != nil {
				c.logger.Debug("accepting bidirectional stream failed", "error", err)
			}
			return
		}
		fp := &frameParser{
			r:         str,
			closeConn: c.conn.CloseWithError,
			unknownFrameHandler: func(ft FrameType, e error) (processed bool, err error) {
				id := c.conn.Context().Value(quic.ConnectionTracingKey).(quic.ConnectionTracingID)
				return streamHijacker(ft, id, str, e)
			},
		}
		go func() {
			if _, err := fp.ParseNext(c.conn.qlogger); err == errHijacked {
				return
			}
			if err != nil {
				if c.logger != nil {
					c.logger.Debug("error handling stream", "error", err)
				}
			}
			c.conn.CloseWithError(quic.ApplicationErrorCode(ErrCodeFrameUnexpected), "received HTTP/3 frame on bidirectional stream")
		}()
	}
}

// RoundTrip executes a request and returns a response
func (c *ClientConn) RoundTrip(req *http.Request) (*http.Response, error) {
	rsp, err := c.roundTrip(req)
	if err != nil && req.Context().Err() != nil {
		// if the context was canceled, return the context cancellation error
		err = req.Context().Err()
	}
	return rsp, err
}

func (c *ClientConn) roundTrip(req *http.Request) (*http.Response, error) {
	// Immediately send out this request, if this is a 0-RTT request.
	switch req.Method {
	case MethodGet0RTT:
		// don't modify the original request
		reqCopy := *req
		req = &reqCopy
		req.Method = http.MethodGet
	case MethodHead0RTT:
		// don't modify the original request
		reqCopy := *req
		req = &reqCopy
		req.Method = http.MethodHead
	default:
		// wait for the handshake to complete
		select {
		case <-c.conn.HandshakeComplete():
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	// It is only possible to send an Extended CONNECT request once the SETTINGS were received.
	// See section 3 of RFC 8441.
	if isExtendedConnectRequest(req) {
		connCtx := c.conn.Context()
		// wait for the server's SETTINGS frame to arrive
		select {
		case <-c.conn.ReceivedSettings():
		case <-connCtx.Done():
			return nil, context.Cause(connCtx)
		}
		if !c.conn.Settings().EnableExtendedConnect {
			return nil, errors.New("http3: server didn't enable Extended CONNECT")
		}
	}

	reqDone := make(chan struct{})
	str, err := c.conn.openRequestStream(
		req.Context(),
		c.requestWriter,
		reqDone,
		c.disableCompression,
		c.maxResponseHeaderBytes,
	)
	if err != nil {
		return nil, &errConnUnusable{e: err}
	}

	// Request Cancellation:
	// This go routine keeps running even after RoundTripOpt() returns.
	// It is shut down when the application is done processing the body.
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-req.Context().Done():
			str.CancelWrite(quic.StreamErrorCode(ErrCodeRequestCanceled))
			str.CancelRead(quic.StreamErrorCode(ErrCodeRequestCanceled))
		case <-reqDone:
		}
	}()

	rsp, err := c.doRequest(req, str)
	if err != nil { // if any error occurred
		close(reqDone)
		<-done
		return nil, maybeReplaceError(err)
	}
	return rsp, maybeReplaceError(err)
}

// ReceivedSettings returns a channel that is closed once the server's HTTP/3 settings were received.
// Settings can be obtained from the Settings method after the channel was closed.
func (c *ClientConn) ReceivedSettings() <-chan struct{} {
	return c.conn.ReceivedSettings()
}

// Settings returns the HTTP/3 settings for this connection.
// It is only valid to call this function after the channel returned by ReceivedSettings was closed.
func (c *ClientConn) Settings() *Settings {
	return c.conn.Settings()
}

// CloseWithError closes the connection with the given error code and message.
// It is invalid to call this function after the connection was closed.
func (c *ClientConn) CloseWithError(code ErrCode, msg string) error {
	return c.conn.CloseWithError(quic.ApplicationErrorCode(code), msg)
}

// Context returns a context that is cancelled when the connection is closed.
func (c *ClientConn) Context() context.Context {
	return c.conn.Context()
}

// cancelingReader reads from the io.Reader.
// It cancels writing on the stream if any error other than io.EOF occurs.
type cancelingReader struct {
	r   io.Reader
	str *RequestStream
}

func (r *cancelingReader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	if err != nil && err != io.EOF {
		r.str.CancelWrite(quic.StreamErrorCode(ErrCodeRequestCanceled))
	}
	return n, err
}

func (c *ClientConn) sendRequestBody(str *RequestStream, body io.ReadCloser, contentLength int64) error {
	defer body.Close()
	buf := make([]byte, bodyCopyBufferSize)
	sr := &cancelingReader{str: str, r: body}
	if contentLength == -1 {
		_, err := io.CopyBuffer(str, sr, buf)
		return err
	}

	// make sure we don't send more bytes than the content length
	n, err := io.CopyBuffer(str, io.LimitReader(sr, contentLength), buf)
	if err != nil {
		return err
	}
	var extra int64
	extra, err = io.CopyBuffer(io.Discard, sr, buf)
	n += extra
	if n > contentLength {
		str.CancelWrite(quic.StreamErrorCode(ErrCodeRequestCanceled))
		return fmt.Errorf("http: ContentLength=%d with Body length %d", contentLength, n)
	}
	return err
}

func (c *ClientConn) doRequest(req *http.Request, str *RequestStream) (*http.Response, error) {
	trace := httptrace.ContextClientTrace(req.Context())
	if err := str.sendRequestHeader(req); err != nil {
		traceWroteRequest(trace, err)
		return nil, err
	}
	if req.Body == nil {
		traceWroteRequest(trace, nil)
		str.Close()
	} else {
		// send the request body asynchronously
		go func() {
			contentLength := int64(-1)
			// According to the documentation for http.Request.ContentLength,
			// a value of 0 with a non-nil Body is also treated as unknown content length.
			if req.ContentLength > 0 {
				contentLength = req.ContentLength
			}
			err := c.sendRequestBody(str, req.Body, contentLength)
			traceWroteRequest(trace, err)
			if err != nil {
				if c.logger != nil {
					c.logger.Debug("error writing request", "error", err)
				}
			}
			str.Close()
		}()
	}

	// copy from net/http: support 1xx responses
	var num1xx int // number of informational 1xx headers received
	var res *http.Response
	for {
		var err error
		res, err = str.ReadResponse()
		if err != nil {
			return nil, err
		}
		resCode := res.StatusCode
		is1xx := 100 <= resCode && resCode <= 199
		// treat 101 as a terminal status, see https://github.com/golang/go/issues/26161
		is1xxNonTerminal := is1xx && resCode != http.StatusSwitchingProtocols
		if is1xxNonTerminal {
			num1xx++
			if num1xx > max1xxResponses {
				str.CancelRead(quic.StreamErrorCode(ErrCodeExcessiveLoad))
				str.CancelWrite(quic.StreamErrorCode(ErrCodeExcessiveLoad))
				return nil, errors.New("http3: too many 1xx informational responses")
			}
			traceGot1xxResponse(trace, resCode, textproto.MIMEHeader(res.Header))
			if resCode == http.StatusContinue {
				traceGot100Continue(trace)
			}
			continue
		}
		break
	}
	connState := c.conn.ConnectionState().TLS
	res.TLS = &connState
	res.Request = req
	return res, nil
}

// Conn returns the underlying HTTP/3 connection.
// This method is only useful for advanced use cases, such as when the application needs to
// open streams on the HTTP/3 connection (e.g. WebTransport).
func (c *ClientConn) Conn() *Conn {
	return c.conn
}
