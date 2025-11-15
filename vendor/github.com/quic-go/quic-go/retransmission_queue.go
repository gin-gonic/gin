package quic

import (
	"fmt"

	"github.com/quic-go/quic-go/internal/ackhandler"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
)

type framesToRetransmit struct {
	crypto []*wire.CryptoFrame
	other  []wire.Frame
}

type retransmissionQueue struct {
	initial   *framesToRetransmit
	handshake *framesToRetransmit
	appData   framesToRetransmit
}

func newRetransmissionQueue() *retransmissionQueue {
	return &retransmissionQueue{
		initial:   &framesToRetransmit{},
		handshake: &framesToRetransmit{},
	}
}

func (q *retransmissionQueue) addInitial(f wire.Frame) {
	if q.initial == nil {
		return
	}
	if cf, ok := f.(*wire.CryptoFrame); ok {
		q.initial.crypto = append(q.initial.crypto, cf)
		return
	}
	q.initial.other = append(q.initial.other, f)
}

func (q *retransmissionQueue) addHandshake(f wire.Frame) {
	if q.handshake == nil {
		return
	}
	if cf, ok := f.(*wire.CryptoFrame); ok {
		q.handshake.crypto = append(q.handshake.crypto, cf)
		return
	}
	q.handshake.other = append(q.handshake.other, f)
}

func (q *retransmissionQueue) addAppData(f wire.Frame) {
	switch f := f.(type) {
	case *wire.StreamFrame:
		panic("STREAM frames are handled with their respective streams.")
	case *wire.CryptoFrame:
		q.appData.crypto = append(q.appData.crypto, f)
	default:
		q.appData.other = append(q.appData.other, f)
	}
}

func (q *retransmissionQueue) HasData(encLevel protocol.EncryptionLevel) bool {
	//nolint:exhaustive // 0-RTT data is retransmitted in 1-RTT packets.
	switch encLevel {
	case protocol.EncryptionInitial:
		return q.initial != nil &&
			(len(q.initial.crypto) > 0 || len(q.initial.other) > 0)
	case protocol.EncryptionHandshake:
		return q.handshake != nil &&
			(len(q.handshake.crypto) > 0 || len(q.handshake.other) > 0)
	case protocol.Encryption1RTT:
		return len(q.appData.crypto) > 0 || len(q.appData.other) > 0
	}
	return false
}

func (q *retransmissionQueue) GetFrame(encLevel protocol.EncryptionLevel, maxLen protocol.ByteCount, v protocol.Version) wire.Frame {
	var r *framesToRetransmit
	//nolint:exhaustive // 0-RTT data is retransmitted in 1-RTT packets.
	switch encLevel {
	case protocol.EncryptionInitial:
		r = q.initial
	case protocol.EncryptionHandshake:
		r = q.handshake
	case protocol.Encryption1RTT:
		r = &q.appData
	}
	if r == nil {
		return nil
	}

	if len(r.crypto) > 0 {
		f := r.crypto[0]
		newFrame, needsSplit := f.MaybeSplitOffFrame(maxLen, v)
		if newFrame == nil && !needsSplit { // the whole frame fits
			r.crypto = r.crypto[1:]
			return f
		}
		if newFrame != nil { // frame was split. Leave the original frame in the queue.
			return newFrame
		}
	}
	if len(r.other) == 0 {
		return nil
	}
	f := r.other[0]
	if f.Length(v) > maxLen {
		return nil
	}
	r.other = r.other[1:]
	return f
}

func (q *retransmissionQueue) DropPackets(encLevel protocol.EncryptionLevel) {
	//nolint:exhaustive // Can only drop Initial and Handshake packet number space.
	switch encLevel {
	case protocol.EncryptionInitial:
		q.initial = nil
	case protocol.EncryptionHandshake:
		q.handshake = nil
	default:
		panic(fmt.Sprintf("unexpected encryption level: %s", encLevel))
	}
}

func (q *retransmissionQueue) AckHandler(encLevel protocol.EncryptionLevel) ackhandler.FrameHandler {
	switch encLevel {
	case protocol.EncryptionInitial:
		return (*retransmissionQueueInitialAckHandler)(q)
	case protocol.EncryptionHandshake:
		return (*retransmissionQueueHandshakeAckHandler)(q)
	case protocol.Encryption0RTT, protocol.Encryption1RTT:
		return (*retransmissionQueueAppDataAckHandler)(q)
	}
	return nil
}

type retransmissionQueueInitialAckHandler retransmissionQueue

func (q *retransmissionQueueInitialAckHandler) OnAcked(wire.Frame) {}
func (q *retransmissionQueueInitialAckHandler) OnLost(f wire.Frame) {
	(*retransmissionQueue)(q).addInitial(f)
}

type retransmissionQueueHandshakeAckHandler retransmissionQueue

func (q *retransmissionQueueHandshakeAckHandler) OnAcked(wire.Frame) {}
func (q *retransmissionQueueHandshakeAckHandler) OnLost(f wire.Frame) {
	(*retransmissionQueue)(q).addHandshake(f)
}

type retransmissionQueueAppDataAckHandler retransmissionQueue

func (q *retransmissionQueueAppDataAckHandler) OnAcked(wire.Frame) {}
func (q *retransmissionQueueAppDataAckHandler) OnLost(f wire.Frame) {
	(*retransmissionQueue)(q).addAppData(f)
}
