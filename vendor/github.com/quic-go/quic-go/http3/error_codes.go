package http3

import (
	"fmt"

	"github.com/quic-go/quic-go"
)

type ErrCode quic.ApplicationErrorCode

const (
	ErrCodeNoError              ErrCode = 0x100
	ErrCodeGeneralProtocolError ErrCode = 0x101
	ErrCodeInternalError        ErrCode = 0x102
	ErrCodeStreamCreationError  ErrCode = 0x103
	ErrCodeClosedCriticalStream ErrCode = 0x104
	ErrCodeFrameUnexpected      ErrCode = 0x105
	ErrCodeFrameError           ErrCode = 0x106
	ErrCodeExcessiveLoad        ErrCode = 0x107
	ErrCodeIDError              ErrCode = 0x108
	ErrCodeSettingsError        ErrCode = 0x109
	ErrCodeMissingSettings      ErrCode = 0x10a
	ErrCodeRequestRejected      ErrCode = 0x10b
	ErrCodeRequestCanceled      ErrCode = 0x10c
	ErrCodeRequestIncomplete    ErrCode = 0x10d
	ErrCodeMessageError         ErrCode = 0x10e
	ErrCodeConnectError         ErrCode = 0x10f
	ErrCodeVersionFallback      ErrCode = 0x110
	ErrCodeDatagramError        ErrCode = 0x33
)

func (e ErrCode) String() string {
	s := e.string()
	if s != "" {
		return s
	}
	return fmt.Sprintf("unknown error code: %#x", uint16(e))
}

func (e ErrCode) string() string {
	switch e {
	case ErrCodeNoError:
		return "H3_NO_ERROR"
	case ErrCodeGeneralProtocolError:
		return "H3_GENERAL_PROTOCOL_ERROR"
	case ErrCodeInternalError:
		return "H3_INTERNAL_ERROR"
	case ErrCodeStreamCreationError:
		return "H3_STREAM_CREATION_ERROR"
	case ErrCodeClosedCriticalStream:
		return "H3_CLOSED_CRITICAL_STREAM"
	case ErrCodeFrameUnexpected:
		return "H3_FRAME_UNEXPECTED"
	case ErrCodeFrameError:
		return "H3_FRAME_ERROR"
	case ErrCodeExcessiveLoad:
		return "H3_EXCESSIVE_LOAD"
	case ErrCodeIDError:
		return "H3_ID_ERROR"
	case ErrCodeSettingsError:
		return "H3_SETTINGS_ERROR"
	case ErrCodeMissingSettings:
		return "H3_MISSING_SETTINGS"
	case ErrCodeRequestRejected:
		return "H3_REQUEST_REJECTED"
	case ErrCodeRequestCanceled:
		return "H3_REQUEST_CANCELLED"
	case ErrCodeRequestIncomplete:
		return "H3_INCOMPLETE_REQUEST"
	case ErrCodeMessageError:
		return "H3_MESSAGE_ERROR"
	case ErrCodeConnectError:
		return "H3_CONNECT_ERROR"
	case ErrCodeVersionFallback:
		return "H3_VERSION_FALLBACK"
	case ErrCodeDatagramError:
		return "H3_DATAGRAM_ERROR"
	default:
		return ""
	}
}
