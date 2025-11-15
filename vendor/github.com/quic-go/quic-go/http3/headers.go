package http3

import (
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"

	"github.com/quic-go/qpack"
)

type header struct {
	// Pseudo header fields defined in RFC 9114
	Path      string
	Method    string
	Authority string
	Scheme    string
	Status    string
	// for Extended connect
	Protocol string
	// parsed and deduplicated. -1 if no Content-Length header is sent
	ContentLength int64
	// all non-pseudo headers
	Headers http.Header
}

// connection-specific header fields must not be sent on HTTP/3
var invalidHeaderFields = [...]string{
	"connection",
	"keep-alive",
	"proxy-connection",
	"transfer-encoding",
	"upgrade",
}

func parseHeaders(headers []qpack.HeaderField, isRequest bool) (header, error) {
	hdr := header{Headers: make(http.Header, len(headers))}
	var readFirstRegularHeader, readContentLength bool
	var contentLengthStr string
	for _, h := range headers {
		// field names need to be lowercase, see section 4.2 of RFC 9114
		if strings.ToLower(h.Name) != h.Name {
			return header{}, fmt.Errorf("header field is not lower-case: %s", h.Name)
		}
		if !httpguts.ValidHeaderFieldValue(h.Value) {
			return header{}, fmt.Errorf("invalid header field value for %s: %q", h.Name, h.Value)
		}
		if h.IsPseudo() {
			if readFirstRegularHeader {
				// all pseudo headers must appear before regular header fields, see section 4.3 of RFC 9114
				return header{}, fmt.Errorf("received pseudo header %s after a regular header field", h.Name)
			}
			var isResponsePseudoHeader bool  // pseudo headers are either valid for requests or for responses
			var isDuplicatePseudoHeader bool // pseudo headers are allowed to appear exactly once
			switch h.Name {
			case ":path":
				isDuplicatePseudoHeader = hdr.Path != ""
				hdr.Path = h.Value
			case ":method":
				isDuplicatePseudoHeader = hdr.Method != ""
				hdr.Method = h.Value
			case ":authority":
				isDuplicatePseudoHeader = hdr.Authority != ""
				hdr.Authority = h.Value
			case ":protocol":
				isDuplicatePseudoHeader = hdr.Protocol != ""
				hdr.Protocol = h.Value
			case ":scheme":
				isDuplicatePseudoHeader = hdr.Scheme != ""
				hdr.Scheme = h.Value
			case ":status":
				isDuplicatePseudoHeader = hdr.Status != ""
				hdr.Status = h.Value
				isResponsePseudoHeader = true
			default:
				return header{}, fmt.Errorf("unknown pseudo header: %s", h.Name)
			}
			if isDuplicatePseudoHeader {
				return header{}, fmt.Errorf("duplicate pseudo header: %s", h.Name)
			}
			if isRequest && isResponsePseudoHeader {
				return header{}, fmt.Errorf("invalid request pseudo header: %s", h.Name)
			}
			if !isRequest && !isResponsePseudoHeader {
				return header{}, fmt.Errorf("invalid response pseudo header: %s", h.Name)
			}
		} else {
			if !httpguts.ValidHeaderFieldName(h.Name) {
				return header{}, fmt.Errorf("invalid header field name: %q", h.Name)
			}
			for _, invalidField := range invalidHeaderFields {
				if h.Name == invalidField {
					return header{}, fmt.Errorf("invalid header field name: %q", h.Name)
				}
			}
			if h.Name == "te" && h.Value != "trailers" {
				return header{}, fmt.Errorf("invalid TE header field value: %q", h.Value)
			}
			readFirstRegularHeader = true
			switch h.Name {
			case "content-length":
				// Ignore duplicate Content-Length headers.
				// Fail if the duplicates differ.
				if !readContentLength {
					readContentLength = true
					contentLengthStr = h.Value
				} else if contentLengthStr != h.Value {
					return header{}, fmt.Errorf("contradicting content lengths (%s and %s)", contentLengthStr, h.Value)
				}
			default:
				hdr.Headers.Add(h.Name, h.Value)
			}
		}
	}
	hdr.ContentLength = -1
	if len(contentLengthStr) > 0 {
		// use ParseUint instead of ParseInt, so that parsing fails on negative values
		cl, err := strconv.ParseUint(contentLengthStr, 10, 63)
		if err != nil {
			return header{}, fmt.Errorf("invalid content length: %w", err)
		}
		hdr.Headers.Set("Content-Length", contentLengthStr)
		hdr.ContentLength = int64(cl)
	}
	return hdr, nil
}

func parseTrailers(headers []qpack.HeaderField) (http.Header, error) {
	h := make(http.Header, len(headers))
	for _, field := range headers {
		if field.IsPseudo() {
			return nil, fmt.Errorf("http3: received pseudo header in trailer: %s", field.Name)
		}
		h.Add(field.Name, field.Value)
	}
	return h, nil
}

func requestFromHeaders(headerFields []qpack.HeaderField) (*http.Request, error) {
	hdr, err := parseHeaders(headerFields, true)
	if err != nil {
		return nil, err
	}
	// concatenate cookie headers, see https://tools.ietf.org/html/rfc6265#section-5.4
	if len(hdr.Headers["Cookie"]) > 0 {
		hdr.Headers.Set("Cookie", strings.Join(hdr.Headers["Cookie"], "; "))
	}

	isConnect := hdr.Method == http.MethodConnect
	// Extended CONNECT, see https://datatracker.ietf.org/doc/html/rfc8441#section-4
	isExtendedConnected := isConnect && hdr.Protocol != ""
	if isExtendedConnected {
		if hdr.Scheme == "" || hdr.Path == "" || hdr.Authority == "" {
			return nil, errors.New("extended CONNECT: :scheme, :path and :authority must not be empty")
		}
	} else if isConnect {
		if hdr.Path != "" || hdr.Authority == "" { // normal CONNECT
			return nil, errors.New(":path must be empty and :authority must not be empty")
		}
	} else if len(hdr.Path) == 0 || len(hdr.Authority) == 0 || len(hdr.Method) == 0 {
		return nil, errors.New(":path, :authority and :method must not be empty")
	}

	if !isExtendedConnected && len(hdr.Protocol) > 0 {
		return nil, errors.New(":protocol must be empty")
	}

	var u *url.URL
	var requestURI string

	protocol := "HTTP/3.0"

	if isConnect {
		u = &url.URL{}
		if isExtendedConnected {
			u, err = url.ParseRequestURI(hdr.Path)
			if err != nil {
				return nil, err
			}
			protocol = hdr.Protocol
		} else {
			u.Path = hdr.Path
		}
		u.Scheme = hdr.Scheme
		u.Host = hdr.Authority
		requestURI = hdr.Authority
	} else {
		u, err = url.ParseRequestURI(hdr.Path)
		if err != nil {
			return nil, fmt.Errorf("invalid content length: %w", err)
		}
		requestURI = hdr.Path
	}

	return &http.Request{
		Method:        hdr.Method,
		URL:           u,
		Proto:         protocol,
		ProtoMajor:    3,
		ProtoMinor:    0,
		Header:        hdr.Headers,
		Body:          nil,
		ContentLength: hdr.ContentLength,
		Host:          hdr.Authority,
		RequestURI:    requestURI,
	}, nil
}

// updateResponseFromHeaders sets up http.Response as an HTTP/3 response,
// using the decoded qpack header filed.
// It is only called for the HTTP header (and not the HTTP trailer).
// It takes an http.Response as an argument to allow the caller to set the trailer later on.
func updateResponseFromHeaders(rsp *http.Response, headerFields []qpack.HeaderField) error {
	hdr, err := parseHeaders(headerFields, false)
	if err != nil {
		return err
	}
	if hdr.Status == "" {
		return errors.New("missing :status field")
	}
	rsp.Proto = "HTTP/3.0"
	rsp.ProtoMajor = 3
	rsp.Header = hdr.Headers
	processTrailers(rsp)
	rsp.ContentLength = hdr.ContentLength

	status, err := strconv.Atoi(hdr.Status)
	if err != nil {
		return fmt.Errorf("invalid status code: %w", err)
	}
	rsp.StatusCode = status
	rsp.Status = hdr.Status + " " + http.StatusText(status)
	return nil
}

// processTrailers initializes the rsp.Trailer map, and adds keys for every announced header value.
// The Trailer header is removed from the http.Response.Header map.
// It handles both duplicate as well as comma-separated values for the Trailer header.
// For example:
//
//	Trailer: Trailer1, Trailer2
//	Trailer: Trailer3
//
// Will result in a http.Response.Trailer map containing the keys "Trailer1", "Trailer2", "Trailer3".
func processTrailers(rsp *http.Response) {
	rawTrailers, ok := rsp.Header["Trailer"]
	if !ok {
		return
	}

	rsp.Trailer = make(http.Header)
	for _, rawVal := range rawTrailers {
		for _, val := range strings.Split(rawVal, ",") {
			rsp.Trailer[http.CanonicalHeaderKey(textproto.TrimString(val))] = nil
		}
	}
	delete(rsp.Header, "Trailer")
}
