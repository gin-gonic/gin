package http3

import (
	"crypto/tls"
	"net"
	"net/http/httptrace"
	"net/textproto"
	"time"

	"github.com/quic-go/quic-go"
)

func traceGetConn(trace *httptrace.ClientTrace, hostPort string) {
	if trace != nil && trace.GetConn != nil {
		trace.GetConn(hostPort)
	}
}

// fakeConn is a wrapper for quic.EarlyConnection
// because the quic connection does not implement net.Conn.
type fakeConn struct {
	conn *quic.Conn
}

func (c *fakeConn) Close() error                       { panic("connection operation prohibited") }
func (c *fakeConn) Read(p []byte) (int, error)         { panic("connection operation prohibited") }
func (c *fakeConn) Write(p []byte) (int, error)        { panic("connection operation prohibited") }
func (c *fakeConn) SetDeadline(t time.Time) error      { panic("connection operation prohibited") }
func (c *fakeConn) SetReadDeadline(t time.Time) error  { panic("connection operation prohibited") }
func (c *fakeConn) SetWriteDeadline(t time.Time) error { panic("connection operation prohibited") }
func (c *fakeConn) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *fakeConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }

func traceGotConn(trace *httptrace.ClientTrace, conn *quic.Conn, reused bool) {
	if trace != nil && trace.GotConn != nil {
		trace.GotConn(httptrace.GotConnInfo{
			Conn:   &fakeConn{conn: conn},
			Reused: reused,
		})
	}
}

func traceGotFirstResponseByte(trace *httptrace.ClientTrace) {
	if trace != nil && trace.GotFirstResponseByte != nil {
		trace.GotFirstResponseByte()
	}
}

func traceGot1xxResponse(trace *httptrace.ClientTrace, code int, header textproto.MIMEHeader) {
	if trace != nil && trace.Got1xxResponse != nil {
		trace.Got1xxResponse(code, header)
	}
}

func traceGot100Continue(trace *httptrace.ClientTrace) {
	if trace != nil && trace.Got100Continue != nil {
		trace.Got100Continue()
	}
}

func traceHasWroteHeaderField(trace *httptrace.ClientTrace) bool {
	return trace != nil && trace.WroteHeaderField != nil
}

func traceWroteHeaderField(trace *httptrace.ClientTrace, k, v string) {
	if trace != nil && trace.WroteHeaderField != nil {
		trace.WroteHeaderField(k, []string{v})
	}
}

func traceWroteHeaders(trace *httptrace.ClientTrace) {
	if trace != nil && trace.WroteHeaders != nil {
		trace.WroteHeaders()
	}
}

func traceWroteRequest(trace *httptrace.ClientTrace, err error) {
	if trace != nil && trace.WroteRequest != nil {
		trace.WroteRequest(httptrace.WroteRequestInfo{Err: err})
	}
}

func traceConnectStart(trace *httptrace.ClientTrace, network, addr string) {
	if trace != nil && trace.ConnectStart != nil {
		trace.ConnectStart(network, addr)
	}
}

func traceConnectDone(trace *httptrace.ClientTrace, network, addr string, err error) {
	if trace != nil && trace.ConnectDone != nil {
		trace.ConnectDone(network, addr, err)
	}
}

func traceTLSHandshakeStart(trace *httptrace.ClientTrace) {
	if trace != nil && trace.TLSHandshakeStart != nil {
		trace.TLSHandshakeStart()
	}
}

func traceTLSHandshakeDone(trace *httptrace.ClientTrace, state tls.ConnectionState, err error) {
	if trace != nil && trace.TLSHandshakeDone != nil {
		trace.TLSHandshakeDone(state, err)
	}
}
