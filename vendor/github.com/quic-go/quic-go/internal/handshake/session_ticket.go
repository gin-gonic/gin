package handshake

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/quicvarint"
)

const sessionTicketRevision = 5

type sessionTicket struct {
	Parameters *wire.TransportParameters
}

func (t *sessionTicket) Marshal() []byte {
	b := make([]byte, 0, 256)
	b = quicvarint.Append(b, sessionTicketRevision)
	return t.Parameters.MarshalForSessionTicket(b)
}

func (t *sessionTicket) Unmarshal(b []byte) error {
	rev, l, err := quicvarint.Parse(b)
	if err != nil {
		return errors.New("failed to read session ticket revision")
	}
	b = b[l:]
	if rev != sessionTicketRevision {
		return fmt.Errorf("unknown session ticket revision: %d", rev)
	}
	var tp wire.TransportParameters
	if err := tp.UnmarshalFromSessionTicket(b); err != nil {
		return fmt.Errorf("unmarshaling transport parameters from session ticket failed: %s", err.Error())
	}
	t.Parameters = &tp
	return nil
}

const extraPrefix = "quic-go1"

func addSessionStateExtraPrefix(b []byte) []byte {
	return append([]byte(extraPrefix), b...)
}

func findSessionStateExtraData(extras [][]byte) []byte {
	prefix := []byte(extraPrefix)
	for _, extra := range extras {
		if len(extra) < len(prefix) || !bytes.Equal(prefix, extra[:len(prefix)]) {
			continue
		}
		return extra[len(prefix):]
	}
	return nil
}
