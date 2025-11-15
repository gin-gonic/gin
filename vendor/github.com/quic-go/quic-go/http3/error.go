package http3

import (
	"errors"
	"fmt"

	"github.com/quic-go/quic-go"
)

// Error is returned from the round tripper (for HTTP clients)
// and inside the HTTP handler (for HTTP servers) if an HTTP/3 error occurs.
// See section 8 of RFC 9114.
type Error struct {
	Remote       bool
	ErrorCode    ErrCode
	ErrorMessage string
}

var _ error = &Error{}

func (e *Error) Error() string {
	s := e.ErrorCode.string()
	if s == "" {
		s = fmt.Sprintf("H3 error (%#x)", uint64(e.ErrorCode))
	}
	// Usually errors are remote. Only make it explicit for local errors.
	if !e.Remote {
		s += " (local)"
	}
	if e.ErrorMessage != "" {
		s += ": " + e.ErrorMessage
	}
	return s
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && e.ErrorCode == t.ErrorCode && e.Remote == t.Remote
}

func maybeReplaceError(err error) error {
	if err == nil {
		return nil
	}

	var (
		e      Error
		strErr *quic.StreamError
		appErr *quic.ApplicationError
	)
	switch {
	default:
		return err
	case errors.As(err, &strErr):
		e.Remote = strErr.Remote
		e.ErrorCode = ErrCode(strErr.ErrorCode)
	case errors.As(err, &appErr):
		e.Remote = appErr.Remote
		e.ErrorCode = ErrCode(appErr.ErrorCode)
		e.ErrorMessage = appErr.ErrorMessage
	}
	return &e
}
