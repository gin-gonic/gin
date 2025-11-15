package http3

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/http/httpguts"

	"github.com/quic-go/quic-go"
)

// Settings are HTTP/3 settings that apply to the underlying connection.
type Settings struct {
	// Support for HTTP/3 datagrams (RFC 9297)
	EnableDatagrams bool
	// Extended CONNECT, RFC 9220
	EnableExtendedConnect bool
	// Other settings, defined by the application
	Other map[uint64]uint64
}

// RoundTripOpt are options for the Transport.RoundTripOpt method.
type RoundTripOpt struct {
	// OnlyCachedConn controls whether the Transport may create a new QUIC connection.
	// If set true and no cached connection is available, RoundTripOpt will return ErrNoCachedConn.
	OnlyCachedConn bool
}

type clientConn interface {
	OpenRequestStream(context.Context) (*RequestStream, error)
	RoundTrip(*http.Request) (*http.Response, error)
}

type roundTripperWithCount struct {
	cancel     context.CancelFunc
	dialing    chan struct{} // closed as soon as quic.Dial(Early) returned
	dialErr    error
	conn       *quic.Conn
	clientConn clientConn

	useCount atomic.Int64
}

func (r *roundTripperWithCount) Close() error {
	r.cancel()
	<-r.dialing
	if r.conn != nil {
		return r.conn.CloseWithError(0, "")
	}
	return nil
}

// Transport implements the http.RoundTripper interface
type Transport struct {
	// TLSClientConfig specifies the TLS configuration to use with
	// tls.Client. If nil, the default configuration is used.
	TLSClientConfig *tls.Config

	// QUICConfig is the quic.Config used for dialing new connections.
	// If nil, reasonable default values will be used.
	QUICConfig *quic.Config

	// Dial specifies an optional dial function for creating QUIC
	// connections for requests.
	// If Dial is nil, a UDPConn will be created at the first request
	// and will be reused for subsequent connections to other servers.
	Dial func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (*quic.Conn, error)

	// Enable support for HTTP/3 datagrams (RFC 9297).
	// If a QUICConfig is set, datagram support also needs to be enabled on the QUIC layer by setting EnableDatagrams.
	EnableDatagrams bool

	// Additional HTTP/3 settings.
	// It is invalid to specify any settings defined by RFC 9114 (HTTP/3) and RFC 9297 (HTTP Datagrams).
	AdditionalSettings map[uint64]uint64

	// MaxResponseHeaderBytes specifies a limit on how many response bytes are
	// allowed in the server's response header.
	// Zero means to use a default limit.
	MaxResponseHeaderBytes int64

	// DisableCompression, if true, prevents the Transport from requesting compression with an
	// "Accept-Encoding: gzip" request header when the Request contains no existing Accept-Encoding value.
	// If the Transport requests gzip on its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body.
	// However, if the user explicitly requested gzip it is not automatically uncompressed.
	DisableCompression bool

	StreamHijacker    func(FrameType, quic.ConnectionTracingID, *quic.Stream, error) (hijacked bool, err error)
	UniStreamHijacker func(StreamType, quic.ConnectionTracingID, *quic.ReceiveStream, error) (hijacked bool)

	Logger *slog.Logger

	mutex sync.Mutex

	initOnce sync.Once
	initErr  error

	newClientConn func(*quic.Conn) clientConn

	clients   map[string]*roundTripperWithCount
	transport *quic.Transport
	closed    bool
}

var (
	_ http.RoundTripper = &Transport{}
	_ io.Closer         = &Transport{}
)

var (
	// ErrNoCachedConn is returned when Transport.OnlyCachedConn is set
	ErrNoCachedConn = errors.New("http3: no cached connection was available")
	// ErrTransportClosed is returned when attempting to use a closed Transport
	ErrTransportClosed = errors.New("http3: transport is closed")
)

func (t *Transport) init() error {
	if t.newClientConn == nil {
		t.newClientConn = func(conn *quic.Conn) clientConn {
			return newClientConn(
				conn,
				t.EnableDatagrams,
				t.AdditionalSettings,
				t.StreamHijacker,
				t.UniStreamHijacker,
				t.MaxResponseHeaderBytes,
				t.DisableCompression,
				t.Logger,
			)
		}
	}
	if t.QUICConfig == nil {
		t.QUICConfig = defaultQuicConfig.Clone()
		t.QUICConfig.EnableDatagrams = t.EnableDatagrams
	}
	if t.EnableDatagrams && !t.QUICConfig.EnableDatagrams {
		return errors.New("HTTP Datagrams enabled, but QUIC Datagrams disabled")
	}
	if len(t.QUICConfig.Versions) == 0 {
		t.QUICConfig = t.QUICConfig.Clone()
		t.QUICConfig.Versions = []quic.Version{quic.SupportedVersions()[0]}
	}
	if len(t.QUICConfig.Versions) != 1 {
		return errors.New("can only use a single QUIC version for dialing a HTTP/3 connection")
	}
	if t.QUICConfig.MaxIncomingStreams == 0 {
		t.QUICConfig.MaxIncomingStreams = -1 // don't allow any bidirectional streams
	}
	if t.Dial == nil {
		udpConn, err := net.ListenUDP("udp", nil)
		if err != nil {
			return err
		}
		t.transport = &quic.Transport{Conn: udpConn}
	}
	return nil
}

// RoundTripOpt is like RoundTrip, but takes options.
func (t *Transport) RoundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	rsp, err := t.roundTripOpt(req, opt)
	if err != nil {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, err
	}
	return rsp, nil
}

func (t *Transport) roundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	t.initOnce.Do(func() { t.initErr = t.init() })
	if t.initErr != nil {
		return nil, t.initErr
	}

	if req.URL == nil {
		return nil, errors.New("http3: nil Request.URL")
	}
	if req.URL.Scheme != "https" {
		return nil, fmt.Errorf("http3: unsupported protocol scheme: %s", req.URL.Scheme)
	}
	if req.URL.Host == "" {
		return nil, errors.New("http3: no Host in request URL")
	}
	if req.Header == nil {
		return nil, errors.New("http3: nil Request.Header")
	}
	if req.Method != "" && !validMethod(req.Method) {
		return nil, fmt.Errorf("http3: invalid method %q", req.Method)
	}
	for k, vv := range req.Header {
		if !httpguts.ValidHeaderFieldName(k) {
			return nil, fmt.Errorf("http3: invalid http header field name %q", k)
		}
		for _, v := range vv {
			if !httpguts.ValidHeaderFieldValue(v) {
				return nil, fmt.Errorf("http3: invalid http header field value %q for key %v", v, k)
			}
		}
	}

	return t.doRoundTripOpt(req, opt, false)
}

func (t *Transport) doRoundTripOpt(req *http.Request, opt RoundTripOpt, isRetried bool) (*http.Response, error) {
	hostname := authorityAddr(hostnameFromURL(req.URL))
	trace := httptrace.ContextClientTrace(req.Context())
	traceGetConn(trace, hostname)
	cl, isReused, err := t.getClient(req.Context(), hostname, opt.OnlyCachedConn)
	if err != nil {
		return nil, err
	}

	select {
	case <-cl.dialing:
	case <-req.Context().Done():
		return nil, context.Cause(req.Context())
	}

	if cl.dialErr != nil {
		t.removeClient(hostname)
		return nil, cl.dialErr
	}
	defer cl.useCount.Add(-1)
	traceGotConn(trace, cl.conn, isReused)
	rsp, err := cl.clientConn.RoundTrip(req)
	if err != nil {
		// request aborted due to context cancellation
		select {
		case <-req.Context().Done():
			return nil, err
		default:
		}
		if isRetried {
			return nil, err
		}

		t.removeClient(hostname)
		req, err = canRetryRequest(err, req)
		if err != nil {
			return nil, err
		}
		return t.doRoundTripOpt(req, opt, true)
	}
	return rsp, nil
}

func canRetryRequest(err error, req *http.Request) (*http.Request, error) {
	// error occurred while opening the stream, we can be sure that the request wasn't sent out
	var connErr *errConnUnusable
	if errors.As(err, &connErr) {
		return req, nil
	}

	// If the request stream is reset, we can only be sure that the request wasn't processed
	// if the error code is H3_REQUEST_REJECTED.
	var e *Error
	if !errors.As(err, &e) || e.ErrorCode != ErrCodeRequestRejected {
		return nil, err
	}
	// if the body is nil (or http.NoBody), it's safe to reuse this request and its body
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil
	}
	// if the request body can be reset back to its original state via req.GetBody, do that
	if req.GetBody != nil {
		newBody, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		reqCopy := *req
		reqCopy.Body = newBody
		req = &reqCopy
		return &reqCopy, nil
	}
	return nil, fmt.Errorf("http3: Transport: cannot retry err [%w] after Request.Body was written; define Request.GetBody to avoid this error", err)
}

// RoundTrip does a round trip.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripOpt(req, RoundTripOpt{})
}

func (t *Transport) getClient(ctx context.Context, hostname string, onlyCached bool) (rtc *roundTripperWithCount, isReused bool, err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.closed {
		return nil, false, ErrTransportClosed
	}

	if t.clients == nil {
		t.clients = make(map[string]*roundTripperWithCount)
	}

	cl, ok := t.clients[hostname]
	if !ok {
		if onlyCached {
			return nil, false, ErrNoCachedConn
		}
		ctx, cancel := context.WithCancel(ctx)
		cl = &roundTripperWithCount{
			dialing: make(chan struct{}),
			cancel:  cancel,
		}
		go func() {
			defer close(cl.dialing)
			defer cancel()
			conn, rt, err := t.dial(ctx, hostname)
			if err != nil {
				cl.dialErr = err
				return
			}
			cl.conn = conn
			cl.clientConn = rt
		}()
		t.clients[hostname] = cl
	}
	select {
	case <-cl.dialing:
		if cl.dialErr != nil {
			delete(t.clients, hostname)
			return nil, false, cl.dialErr
		}
		select {
		case <-cl.conn.HandshakeComplete():
			isReused = true
		default:
		}
	default:
	}
	cl.useCount.Add(1)
	return cl, isReused, nil
}

func (t *Transport) dial(ctx context.Context, hostname string) (*quic.Conn, clientConn, error) {
	var tlsConf *tls.Config
	if t.TLSClientConfig == nil {
		tlsConf = &tls.Config{}
	} else {
		tlsConf = t.TLSClientConfig.Clone()
	}
	if tlsConf.ServerName == "" {
		sni, _, err := net.SplitHostPort(hostname)
		if err != nil {
			// It's ok if net.SplitHostPort returns an error - it could be a hostname/IP address without a port.
			sni = hostname
		}
		tlsConf.ServerName = sni
	}
	// Replace existing ALPNs by H3
	tlsConf.NextProtos = []string{NextProtoH3}

	dial := t.Dial
	if dial == nil {
		dial = func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
			network := "udp"
			udpAddr, err := t.resolveUDPAddr(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			trace := httptrace.ContextClientTrace(ctx)
			traceConnectStart(trace, network, udpAddr.String())
			traceTLSHandshakeStart(trace)
			conn, err := t.transport.DialEarly(ctx, udpAddr, tlsCfg, cfg)
			var state tls.ConnectionState
			if conn != nil {
				state = conn.ConnectionState().TLS
			}
			traceTLSHandshakeDone(trace, state, err)
			traceConnectDone(trace, network, udpAddr.String(), err)
			return conn, err
		}
	}
	conn, err := dial(ctx, hostname, tlsConf, t.QUICConfig)
	if err != nil {
		return nil, nil, err
	}
	return conn, t.newClientConn(conn), nil
}

func (t *Transport) resolveUDPAddr(ctx context.Context, network, addr string) (*net.UDPAddr, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := net.LookupPort(network, portStr)
	if err != nil {
		return nil, err
	}
	resolver := net.DefaultResolver
	ipAddrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	addrs := addrList(ipAddrs)
	ip := addrs.forResolve(network, addr)
	return &net.UDPAddr{IP: ip.IP, Port: port, Zone: ip.Zone}, nil
}

func (t *Transport) removeClient(hostname string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.clients == nil {
		return
	}
	delete(t.clients, hostname)
}

// NewClientConn creates a new HTTP/3 client connection on top of a QUIC connection.
// Most users should use RoundTrip instead of creating a connection directly.
// Specifically, it is not needed to perform GET, POST, HEAD and CONNECT requests.
//
// Obtaining a ClientConn is only needed for more advanced use cases, such as
// using Extended CONNECT for WebTransport or the various MASQUE protocols.
func (t *Transport) NewClientConn(conn *quic.Conn) *ClientConn {
	return newClientConn(
		conn,
		t.EnableDatagrams,
		t.AdditionalSettings,
		t.StreamHijacker,
		t.UniStreamHijacker,
		t.MaxResponseHeaderBytes,
		t.DisableCompression,
		t.Logger,
	)
}

// Close closes the QUIC connections that this Transport has used.
// A Transport cannot be used after it has been closed.
func (t *Transport) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, cl := range t.clients {
		if err := cl.Close(); err != nil {
			return err
		}
	}
	t.clients = nil
	if t.transport != nil {
		if err := t.transport.Close(); err != nil {
			return err
		}
		if err := t.transport.Conn.Close(); err != nil {
			return err
		}
		t.transport = nil
	}
	t.closed = true
	return nil
}

func hostnameFromURL(url *url.URL) string {
	if url != nil {
		return url.Host
	}
	return ""
}

func validMethod(method string) bool {
	/*
				     Method         = "OPTIONS"                ; Section 9.2
		   		                    | "GET"                    ; Section 9.3
		   		                    | "HEAD"                   ; Section 9.4
		   		                    | "POST"                   ; Section 9.5
		   		                    | "PUT"                    ; Section 9.6
		   		                    | "DELETE"                 ; Section 9.7
		   		                    | "TRACE"                  ; Section 9.8
		   		                    | "CONNECT"                ; Section 9.9
		   		                    | extension-method
		   		   extension-method = token
		   		     token          = 1*<any CHAR except CTLs or separators>
	*/
	return len(method) > 0 && strings.IndexFunc(method, isNotToken) == -1
}

// copied from net/http/http.go
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

// CloseIdleConnections closes any QUIC connections in the transport's pool that are currently idle.
// An idle connection is one that was previously used for requests but is now sitting unused.
// This method does not interrupt any connections currently in use.
// It also does not affect connections obtained via NewClientConn.
func (t *Transport) CloseIdleConnections() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for hostname, cl := range t.clients {
		if cl.useCount.Load() == 0 {
			cl.Close()
			delete(t.clients, hostname)
		}
	}
}
