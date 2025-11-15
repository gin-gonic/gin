package quic

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/handshake"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/wire"
)

var errNothingToPack = errors.New("nothing to pack")

type packer interface {
	PackCoalescedPacket(onlyAck bool, maxPacketSize protocol.ByteCount, now monotime.Time, v protocol.Version) (*coalescedPacket, error)
	PackAckOnlyPacket(maxPacketSize protocol.ByteCount, now monotime.Time, v protocol.Version) (shortHeaderPacket, *packetBuffer, error)
	AppendPacket(_ *packetBuffer, maxPacketSize protocol.ByteCount, now monotime.Time, v protocol.Version) (shortHeaderPacket, error)
	PackPTOProbePacket(_ protocol.EncryptionLevel, _ protocol.ByteCount, addPingIfEmpty bool, now monotime.Time, v protocol.Version) (*coalescedPacket, error)
	PackConnectionClose(*qerr.TransportError, protocol.ByteCount, protocol.Version) (*coalescedPacket, error)
	PackApplicationClose(*qerr.ApplicationError, protocol.ByteCount, protocol.Version) (*coalescedPacket, error)
	PackPathProbePacket(protocol.ConnectionID, []ackhandler.Frame, protocol.Version) (shortHeaderPacket, *packetBuffer, error)
	PackMTUProbePacket(ping ackhandler.Frame, size protocol.ByteCount, v protocol.Version) (shortHeaderPacket, *packetBuffer, error)

	SetToken([]byte)
}

type sealer interface {
	handshake.LongHeaderSealer
}

type payload struct {
	streamFrames []ackhandler.StreamFrame
	frames       []ackhandler.Frame
	ack          *wire.AckFrame
	length       protocol.ByteCount
}

type longHeaderPacket struct {
	header       *wire.ExtendedHeader
	ack          *wire.AckFrame
	frames       []ackhandler.Frame
	streamFrames []ackhandler.StreamFrame // only used for 0-RTT packets

	length protocol.ByteCount
}

type shortHeaderPacket struct {
	PacketNumber         protocol.PacketNumber
	Frames               []ackhandler.Frame
	StreamFrames         []ackhandler.StreamFrame
	Ack                  *wire.AckFrame
	Length               protocol.ByteCount
	IsPathMTUProbePacket bool
	IsPathProbePacket    bool

	// used for logging
	DestConnID      protocol.ConnectionID
	PacketNumberLen protocol.PacketNumberLen
	KeyPhase        protocol.KeyPhaseBit
}

func (p *shortHeaderPacket) IsAckEliciting() bool { return ackhandler.HasAckElicitingFrames(p.Frames) }

type coalescedPacket struct {
	buffer         *packetBuffer
	longHdrPackets []*longHeaderPacket
	shortHdrPacket *shortHeaderPacket
}

// IsOnlyShortHeaderPacket says if this packet only contains a short header packet (and no long header packets).
func (p *coalescedPacket) IsOnlyShortHeaderPacket() bool {
	return len(p.longHdrPackets) == 0 && p.shortHdrPacket != nil
}

func (p *longHeaderPacket) EncryptionLevel() protocol.EncryptionLevel {
	//nolint:exhaustive // Will never be called for Retry packets (and they don't have encrypted data).
	switch p.header.Type {
	case protocol.PacketTypeInitial:
		return protocol.EncryptionInitial
	case protocol.PacketTypeHandshake:
		return protocol.EncryptionHandshake
	case protocol.PacketType0RTT:
		return protocol.Encryption0RTT
	default:
		panic("can't determine encryption level")
	}
}

func (p *longHeaderPacket) IsAckEliciting() bool { return ackhandler.HasAckElicitingFrames(p.frames) }

type packetNumberManager interface {
	PeekPacketNumber(protocol.EncryptionLevel) (protocol.PacketNumber, protocol.PacketNumberLen)
	PopPacketNumber(protocol.EncryptionLevel) protocol.PacketNumber
}

type sealingManager interface {
	GetInitialSealer() (handshake.LongHeaderSealer, error)
	GetHandshakeSealer() (handshake.LongHeaderSealer, error)
	Get0RTTSealer() (handshake.LongHeaderSealer, error)
	Get1RTTSealer() (handshake.ShortHeaderSealer, error)
}

type frameSource interface {
	HasData() bool
	Append([]ackhandler.Frame, []ackhandler.StreamFrame, protocol.ByteCount, monotime.Time, protocol.Version) ([]ackhandler.Frame, []ackhandler.StreamFrame, protocol.ByteCount)
}

type ackFrameSource interface {
	GetAckFrame(_ protocol.EncryptionLevel, now monotime.Time, onlyIfQueued bool) *wire.AckFrame
}

type packetPacker struct {
	srcConnID     protocol.ConnectionID
	getDestConnID func() protocol.ConnectionID

	perspective protocol.Perspective
	cryptoSetup sealingManager

	initialStream   *initialCryptoStream
	handshakeStream *cryptoStream

	token []byte

	pnManager           packetNumberManager
	framer              frameSource
	acks                ackFrameSource
	datagramQueue       *datagramQueue
	retransmissionQueue *retransmissionQueue
	rand                rand.Rand

	numNonAckElicitingAcks int
}

var _ packer = &packetPacker{}

func newPacketPacker(
	srcConnID protocol.ConnectionID,
	getDestConnID func() protocol.ConnectionID,
	initialStream *initialCryptoStream,
	handshakeStream *cryptoStream,
	packetNumberManager packetNumberManager,
	retransmissionQueue *retransmissionQueue,
	cryptoSetup sealingManager,
	framer frameSource,
	acks ackFrameSource,
	datagramQueue *datagramQueue,
	perspective protocol.Perspective,
) *packetPacker {
	var b [16]byte
	_, _ = crand.Read(b[:])

	return &packetPacker{
		cryptoSetup:         cryptoSetup,
		getDestConnID:       getDestConnID,
		srcConnID:           srcConnID,
		initialStream:       initialStream,
		handshakeStream:     handshakeStream,
		retransmissionQueue: retransmissionQueue,
		datagramQueue:       datagramQueue,
		perspective:         perspective,
		framer:              framer,
		acks:                acks,
		rand:                *rand.New(rand.NewPCG(binary.BigEndian.Uint64(b[:8]), binary.BigEndian.Uint64(b[8:]))),
		pnManager:           packetNumberManager,
	}
}

// PackConnectionClose packs a packet that closes the connection with a transport error.
func (p *packetPacker) PackConnectionClose(e *qerr.TransportError, maxPacketSize protocol.ByteCount, v protocol.Version) (*coalescedPacket, error) {
	var reason string
	// don't send details of crypto errors
	if !e.ErrorCode.IsCryptoError() {
		reason = e.ErrorMessage
	}
	return p.packConnectionClose(false, uint64(e.ErrorCode), e.FrameType, reason, maxPacketSize, v)
}

// PackApplicationClose packs a packet that closes the connection with an application error.
func (p *packetPacker) PackApplicationClose(e *qerr.ApplicationError, maxPacketSize protocol.ByteCount, v protocol.Version) (*coalescedPacket, error) {
	return p.packConnectionClose(true, uint64(e.ErrorCode), 0, e.ErrorMessage, maxPacketSize, v)
}

func (p *packetPacker) packConnectionClose(
	isApplicationError bool,
	errorCode uint64,
	frameType uint64,
	reason string,
	maxPacketSize protocol.ByteCount,
	v protocol.Version,
) (*coalescedPacket, error) {
	var sealers [4]sealer
	var hdrs [3]*wire.ExtendedHeader
	var payloads [4]payload
	var size protocol.ByteCount
	var connID protocol.ConnectionID
	var oneRTTPacketNumber protocol.PacketNumber
	var oneRTTPacketNumberLen protocol.PacketNumberLen
	var keyPhase protocol.KeyPhaseBit // only set for 1-RTT
	var numLongHdrPackets uint8
	encLevels := [4]protocol.EncryptionLevel{protocol.EncryptionInitial, protocol.EncryptionHandshake, protocol.Encryption0RTT, protocol.Encryption1RTT}
	for i, encLevel := range encLevels {
		if p.perspective == protocol.PerspectiveServer && encLevel == protocol.Encryption0RTT {
			continue
		}
		ccf := &wire.ConnectionCloseFrame{
			IsApplicationError: isApplicationError,
			ErrorCode:          errorCode,
			FrameType:          frameType,
			ReasonPhrase:       reason,
		}
		// don't send application errors in Initial or Handshake packets
		if isApplicationError && (encLevel == protocol.EncryptionInitial || encLevel == protocol.EncryptionHandshake) {
			ccf.IsApplicationError = false
			ccf.ErrorCode = uint64(qerr.ApplicationErrorErrorCode)
			ccf.ReasonPhrase = ""
		}
		pl := payload{
			frames: []ackhandler.Frame{{Frame: ccf}},
			length: ccf.Length(v),
		}

		var sealer sealer
		var err error
		switch encLevel {
		case protocol.EncryptionInitial:
			sealer, err = p.cryptoSetup.GetInitialSealer()
		case protocol.EncryptionHandshake:
			sealer, err = p.cryptoSetup.GetHandshakeSealer()
		case protocol.Encryption0RTT:
			sealer, err = p.cryptoSetup.Get0RTTSealer()
		case protocol.Encryption1RTT:
			var s handshake.ShortHeaderSealer
			s, err = p.cryptoSetup.Get1RTTSealer()
			if err == nil {
				keyPhase = s.KeyPhase()
			}
			sealer = s
		}
		if err == handshake.ErrKeysNotYetAvailable || err == handshake.ErrKeysDropped {
			continue
		}
		if err != nil {
			return nil, err
		}
		sealers[i] = sealer
		var hdr *wire.ExtendedHeader
		if encLevel == protocol.Encryption1RTT {
			connID = p.getDestConnID()
			oneRTTPacketNumber, oneRTTPacketNumberLen = p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
			size += p.shortHeaderPacketLength(connID, oneRTTPacketNumberLen, pl)
		} else {
			hdr = p.getLongHeader(encLevel, v)
			hdrs[i] = hdr
			size += p.longHeaderPacketLength(hdr, pl, v) + protocol.ByteCount(sealer.Overhead())
			numLongHdrPackets++
		}
		payloads[i] = pl
	}
	buffer := getPacketBuffer()
	packet := &coalescedPacket{
		buffer:         buffer,
		longHdrPackets: make([]*longHeaderPacket, 0, numLongHdrPackets),
	}
	for i, encLevel := range encLevels {
		if sealers[i] == nil {
			continue
		}
		if encLevel == protocol.Encryption1RTT {
			shp, err := p.appendShortHeaderPacket(buffer, connID, oneRTTPacketNumber, oneRTTPacketNumberLen, keyPhase, payloads[i], 0, maxPacketSize, sealers[i], false, v)
			if err != nil {
				return nil, err
			}
			packet.shortHdrPacket = &shp
		} else {
			var paddingLen protocol.ByteCount
			if encLevel == protocol.EncryptionInitial {
				paddingLen = p.initialPaddingLen(payloads[i].frames, size, maxPacketSize)
			}
			longHdrPacket, err := p.appendLongHeaderPacket(buffer, hdrs[i], payloads[i], paddingLen, encLevel, sealers[i], v)
			if err != nil {
				return nil, err
			}
			packet.longHdrPackets = append(packet.longHdrPackets, longHdrPacket)
		}
	}
	return packet, nil
}

// longHeaderPacketLength calculates the length of a serialized long header packet.
// It takes into account that packets that have a tiny payload need to be padded,
// such that len(payload) + packet number len >= 4 + AEAD overhead
func (p *packetPacker) longHeaderPacketLength(hdr *wire.ExtendedHeader, pl payload, v protocol.Version) protocol.ByteCount {
	var paddingLen protocol.ByteCount
	pnLen := protocol.ByteCount(hdr.PacketNumberLen)
	if pl.length < 4-pnLen {
		paddingLen = 4 - pnLen - pl.length
	}
	return hdr.GetLength(v) + pl.length + paddingLen
}

// shortHeaderPacketLength calculates the length of a serialized short header packet.
// It takes into account that packets that have a tiny payload need to be padded,
// such that len(payload) + packet number len >= 4 + AEAD overhead
func (p *packetPacker) shortHeaderPacketLength(connID protocol.ConnectionID, pnLen protocol.PacketNumberLen, pl payload) protocol.ByteCount {
	var paddingLen protocol.ByteCount
	if pl.length < 4-protocol.ByteCount(pnLen) {
		paddingLen = 4 - protocol.ByteCount(pnLen) - pl.length
	}
	return wire.ShortHeaderLen(connID, pnLen) + pl.length + paddingLen
}

// size is the expected size of the packet, if no padding was applied.
func (p *packetPacker) initialPaddingLen(frames []ackhandler.Frame, currentSize, maxPacketSize protocol.ByteCount) protocol.ByteCount {
	// For the server, only ack-eliciting Initial packets need to be padded.
	if p.perspective == protocol.PerspectiveServer && !ackhandler.HasAckElicitingFrames(frames) {
		return 0
	}
	if currentSize >= maxPacketSize {
		return 0
	}
	return maxPacketSize - currentSize
}

// PackCoalescedPacket packs a new packet.
// It packs an Initial / Handshake if there is data to send in these packet number spaces.
// It should only be called before the handshake is confirmed.
func (p *packetPacker) PackCoalescedPacket(onlyAck bool, maxSize protocol.ByteCount, now monotime.Time, v protocol.Version) (*coalescedPacket, error) {
	var (
		initialHdr, handshakeHdr, zeroRTTHdr                            *wire.ExtendedHeader
		initialPayload, handshakePayload, zeroRTTPayload, oneRTTPayload payload
		oneRTTPacketNumber                                              protocol.PacketNumber
		oneRTTPacketNumberLen                                           protocol.PacketNumberLen
	)
	// Try packing an Initial packet.
	initialSealer, err := p.cryptoSetup.GetInitialSealer()
	if err != nil && err != handshake.ErrKeysDropped {
		return nil, err
	}
	var size protocol.ByteCount
	if initialSealer != nil {
		initialHdr, initialPayload = p.maybeGetCryptoPacket(
			maxSize-protocol.ByteCount(initialSealer.Overhead()),
			protocol.EncryptionInitial,
			now,
			false,
			onlyAck,
			true,
			v,
		)
		if initialPayload.length > 0 {
			size += p.longHeaderPacketLength(initialHdr, initialPayload, v) + protocol.ByteCount(initialSealer.Overhead())
		}
	}

	// Add a Handshake packet.
	var handshakeSealer sealer
	if (onlyAck && size == 0) || (!onlyAck && size < maxSize-protocol.MinCoalescedPacketSize) {
		var err error
		handshakeSealer, err = p.cryptoSetup.GetHandshakeSealer()
		if err != nil && err != handshake.ErrKeysDropped && err != handshake.ErrKeysNotYetAvailable {
			return nil, err
		}
		if handshakeSealer != nil {
			handshakeHdr, handshakePayload = p.maybeGetCryptoPacket(
				maxSize-size-protocol.ByteCount(handshakeSealer.Overhead()),
				protocol.EncryptionHandshake,
				now,
				false,
				onlyAck,
				size == 0,
				v,
			)
			if handshakePayload.length > 0 {
				s := p.longHeaderPacketLength(handshakeHdr, handshakePayload, v) + protocol.ByteCount(handshakeSealer.Overhead())
				size += s
			}
		}
	}

	// Add a 0-RTT / 1-RTT packet.
	var zeroRTTSealer sealer
	var oneRTTSealer handshake.ShortHeaderSealer
	var connID protocol.ConnectionID
	var kp protocol.KeyPhaseBit
	if (onlyAck && size == 0) || (!onlyAck && size < maxSize-protocol.MinCoalescedPacketSize) {
		var err error
		oneRTTSealer, err = p.cryptoSetup.Get1RTTSealer()
		if err != nil && err != handshake.ErrKeysDropped && err != handshake.ErrKeysNotYetAvailable {
			return nil, err
		}
		if err == nil { // 1-RTT
			kp = oneRTTSealer.KeyPhase()
			connID = p.getDestConnID()
			oneRTTPacketNumber, oneRTTPacketNumberLen = p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
			hdrLen := wire.ShortHeaderLen(connID, oneRTTPacketNumberLen)
			oneRTTPayload = p.maybeGetShortHeaderPacket(oneRTTSealer, hdrLen, maxSize-size, onlyAck, size == 0, now, v)
			if oneRTTPayload.length > 0 {
				size += p.shortHeaderPacketLength(connID, oneRTTPacketNumberLen, oneRTTPayload) + protocol.ByteCount(oneRTTSealer.Overhead())
			}
		} else if p.perspective == protocol.PerspectiveClient && !onlyAck { // 0-RTT packets can't contain ACK frames
			var err error
			zeroRTTSealer, err = p.cryptoSetup.Get0RTTSealer()
			if err != nil && err != handshake.ErrKeysDropped && err != handshake.ErrKeysNotYetAvailable {
				return nil, err
			}
			if zeroRTTSealer != nil {
				zeroRTTHdr, zeroRTTPayload = p.maybeGetAppDataPacketFor0RTT(zeroRTTSealer, maxSize-size, now, v)
				if zeroRTTPayload.length > 0 {
					size += p.longHeaderPacketLength(zeroRTTHdr, zeroRTTPayload, v) + protocol.ByteCount(zeroRTTSealer.Overhead())
				}
			}
		}
	}

	if initialPayload.length == 0 && handshakePayload.length == 0 && zeroRTTPayload.length == 0 && oneRTTPayload.length == 0 {
		return nil, nil
	}

	buffer := getPacketBuffer()
	packet := &coalescedPacket{
		buffer:         buffer,
		longHdrPackets: make([]*longHeaderPacket, 0, 3),
	}
	if initialPayload.length > 0 {
		padding := p.initialPaddingLen(initialPayload.frames, size, maxSize)
		cont, err := p.appendLongHeaderPacket(buffer, initialHdr, initialPayload, padding, protocol.EncryptionInitial, initialSealer, v)
		if err != nil {
			return nil, err
		}
		packet.longHdrPackets = append(packet.longHdrPackets, cont)
	}
	if handshakePayload.length > 0 {
		cont, err := p.appendLongHeaderPacket(buffer, handshakeHdr, handshakePayload, 0, protocol.EncryptionHandshake, handshakeSealer, v)
		if err != nil {
			return nil, err
		}
		packet.longHdrPackets = append(packet.longHdrPackets, cont)
	}
	if zeroRTTPayload.length > 0 {
		longHdrPacket, err := p.appendLongHeaderPacket(buffer, zeroRTTHdr, zeroRTTPayload, 0, protocol.Encryption0RTT, zeroRTTSealer, v)
		if err != nil {
			return nil, err
		}
		packet.longHdrPackets = append(packet.longHdrPackets, longHdrPacket)
	} else if oneRTTPayload.length > 0 {
		shp, err := p.appendShortHeaderPacket(buffer, connID, oneRTTPacketNumber, oneRTTPacketNumberLen, kp, oneRTTPayload, 0, maxSize, oneRTTSealer, false, v)
		if err != nil {
			return nil, err
		}
		packet.shortHdrPacket = &shp
	}
	return packet, nil
}

// PackAckOnlyPacket packs a packet containing only an ACK in the application data packet number space.
// It should be called after the handshake is confirmed.
func (p *packetPacker) PackAckOnlyPacket(maxSize protocol.ByteCount, now monotime.Time, v protocol.Version) (shortHeaderPacket, *packetBuffer, error) {
	buf := getPacketBuffer()
	packet, err := p.appendPacket(buf, true, maxSize, now, v)
	return packet, buf, err
}

// AppendPacket packs a packet in the application data packet number space.
// It should be called after the handshake is confirmed.
func (p *packetPacker) AppendPacket(buf *packetBuffer, maxSize protocol.ByteCount, now monotime.Time, v protocol.Version) (shortHeaderPacket, error) {
	return p.appendPacket(buf, false, maxSize, now, v)
}

func (p *packetPacker) appendPacket(
	buf *packetBuffer,
	onlyAck bool,
	maxPacketSize protocol.ByteCount,
	now monotime.Time,
	v protocol.Version,
) (shortHeaderPacket, error) {
	sealer, err := p.cryptoSetup.Get1RTTSealer()
	if err != nil {
		return shortHeaderPacket{}, err
	}
	pn, pnLen := p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
	connID := p.getDestConnID()
	hdrLen := wire.ShortHeaderLen(connID, pnLen)
	pl := p.maybeGetShortHeaderPacket(sealer, hdrLen, maxPacketSize, onlyAck, true, now, v)
	if pl.length == 0 {
		return shortHeaderPacket{}, errNothingToPack
	}
	kp := sealer.KeyPhase()

	return p.appendShortHeaderPacket(buf, connID, pn, pnLen, kp, pl, 0, maxPacketSize, sealer, false, v)
}

func (p *packetPacker) maybeGetCryptoPacket(
	maxPacketSize protocol.ByteCount,
	encLevel protocol.EncryptionLevel,
	now monotime.Time,
	addPingIfEmpty bool,
	onlyAck, ackAllowed bool,
	v protocol.Version,
) (*wire.ExtendedHeader, payload) {
	if onlyAck {
		if ack := p.acks.GetAckFrame(encLevel, now, true); ack != nil {
			return p.getLongHeader(encLevel, v), payload{
				ack:    ack,
				length: ack.Length(v),
			}
		}
		return nil, payload{}
	}

	var hasCryptoData func() bool
	var popCryptoFrame func(maxLen protocol.ByteCount) *wire.CryptoFrame
	//nolint:exhaustive // Initial and Handshake are the only two encryption levels here.
	switch encLevel {
	case protocol.EncryptionInitial:
		hasCryptoData = p.initialStream.HasData
		popCryptoFrame = p.initialStream.PopCryptoFrame
	case protocol.EncryptionHandshake:
		hasCryptoData = p.handshakeStream.HasData
		popCryptoFrame = p.handshakeStream.PopCryptoFrame
	}
	handler := p.retransmissionQueue.AckHandler(encLevel)
	hasRetransmission := p.retransmissionQueue.HasData(encLevel)

	var ack *wire.AckFrame
	if ackAllowed {
		ack = p.acks.GetAckFrame(encLevel, now, !hasRetransmission && !hasCryptoData())
	}
	var pl payload
	if !hasCryptoData() && !hasRetransmission && ack == nil {
		if !addPingIfEmpty {
			// nothing to send
			return nil, payload{}
		}
		ping := &wire.PingFrame{}
		pl.frames = append(pl.frames, ackhandler.Frame{Frame: ping, Handler: emptyHandler{}})
		pl.length += ping.Length(v)
	}

	if ack != nil {
		pl.ack = ack
		pl.length = ack.Length(v)
		maxPacketSize -= pl.length
	}
	hdr := p.getLongHeader(encLevel, v)
	maxPacketSize -= hdr.GetLength(v)
	if hasRetransmission {
		for {
			frame := p.retransmissionQueue.GetFrame(encLevel, maxPacketSize, v)
			if frame == nil {
				break
			}
			pl.frames = append(pl.frames, ackhandler.Frame{
				Frame:   frame,
				Handler: p.retransmissionQueue.AckHandler(encLevel),
			})
			frameLen := frame.Length(v)
			pl.length += frameLen
			maxPacketSize -= frameLen
		}
		return hdr, pl
	} else {
		for hasCryptoData() {
			cf := popCryptoFrame(maxPacketSize)
			if cf == nil {
				break
			}
			pl.frames = append(pl.frames, ackhandler.Frame{Frame: cf, Handler: handler})
			pl.length += cf.Length(v)
			maxPacketSize -= cf.Length(v)
		}
	}
	return hdr, pl
}

func (p *packetPacker) maybeGetAppDataPacketFor0RTT(sealer sealer, maxSize protocol.ByteCount, now monotime.Time, v protocol.Version) (*wire.ExtendedHeader, payload) {
	if p.perspective != protocol.PerspectiveClient {
		return nil, payload{}
	}

	hdr := p.getLongHeader(protocol.Encryption0RTT, v)
	maxPayloadSize := maxSize - hdr.GetLength(v) - protocol.ByteCount(sealer.Overhead())
	return hdr, p.maybeGetAppDataPacket(maxPayloadSize, false, false, now, v)
}

func (p *packetPacker) maybeGetShortHeaderPacket(
	sealer handshake.ShortHeaderSealer,
	hdrLen, maxPacketSize protocol.ByteCount,
	onlyAck, ackAllowed bool,
	now monotime.Time,
	v protocol.Version,
) payload {
	maxPayloadSize := maxPacketSize - hdrLen - protocol.ByteCount(sealer.Overhead())
	return p.maybeGetAppDataPacket(maxPayloadSize, onlyAck, ackAllowed, now, v)
}

func (p *packetPacker) maybeGetAppDataPacket(
	maxPayloadSize protocol.ByteCount,
	onlyAck, ackAllowed bool,
	now monotime.Time,
	v protocol.Version,
) payload {
	pl := p.composeNextPacket(maxPayloadSize, onlyAck, ackAllowed, now, v)

	// check if we have anything to send
	if len(pl.frames) == 0 && len(pl.streamFrames) == 0 {
		if pl.ack == nil {
			return payload{}
		}
		// the packet only contains an ACK
		if p.numNonAckElicitingAcks >= protocol.MaxNonAckElicitingAcks {
			ping := &wire.PingFrame{}
			pl.frames = append(pl.frames, ackhandler.Frame{Frame: ping})
			pl.length += ping.Length(v)
			p.numNonAckElicitingAcks = 0
		} else {
			p.numNonAckElicitingAcks++
		}
	} else {
		p.numNonAckElicitingAcks = 0
	}
	return pl
}

func (p *packetPacker) composeNextPacket(
	maxPayloadSize protocol.ByteCount,
	onlyAck, ackAllowed bool,
	now monotime.Time,
	v protocol.Version,
) payload {
	if onlyAck {
		if ack := p.acks.GetAckFrame(protocol.Encryption1RTT, now, true); ack != nil {
			return payload{ack: ack, length: ack.Length(v)}
		}
		return payload{}
	}

	hasData := p.framer.HasData()
	hasRetransmission := p.retransmissionQueue.HasData(protocol.Encryption1RTT)

	var hasAck bool
	var pl payload
	if ackAllowed {
		if ack := p.acks.GetAckFrame(protocol.Encryption1RTT, now, !hasRetransmission && !hasData); ack != nil {
			pl.ack = ack
			pl.length += ack.Length(v)
			hasAck = true
		}
	}

	if p.datagramQueue != nil {
		if f := p.datagramQueue.Peek(); f != nil {
			size := f.Length(v)
			if size <= maxPayloadSize-pl.length { // DATAGRAM frame fits
				pl.frames = append(pl.frames, ackhandler.Frame{Frame: f})
				pl.length += size
				p.datagramQueue.Pop()
			} else if !hasAck {
				// The DATAGRAM frame doesn't fit, and the packet doesn't contain an ACK.
				// Discard this frame. There's no point in retrying this in the next packet,
				// as it's unlikely that the available packet size will increase.
				p.datagramQueue.Pop()
			}
			// If the DATAGRAM frame was too large and the packet contained an ACK, we'll try to send it out later.
		}
	}

	if hasAck && !hasData && !hasRetransmission {
		return pl
	}

	if hasRetransmission {
		for {
			remainingLen := maxPayloadSize - pl.length
			if remainingLen < protocol.MinStreamFrameSize {
				break
			}
			f := p.retransmissionQueue.GetFrame(protocol.Encryption1RTT, remainingLen, v)
			if f == nil {
				break
			}
			pl.frames = append(pl.frames, ackhandler.Frame{Frame: f, Handler: p.retransmissionQueue.AckHandler(protocol.Encryption1RTT)})
			pl.length += f.Length(v)
		}
	}

	if hasData {
		var lengthAdded protocol.ByteCount
		startLen := len(pl.frames)
		pl.frames, pl.streamFrames, lengthAdded = p.framer.Append(pl.frames, pl.streamFrames, maxPayloadSize-pl.length, now, v)
		pl.length += lengthAdded
		// add handlers for the control frames that were added
		for i := startLen; i < len(pl.frames); i++ {
			if pl.frames[i].Handler != nil {
				continue
			}
			switch pl.frames[i].Frame.(type) {
			case *wire.PathChallengeFrame, *wire.PathResponseFrame:
				// Path probing is currently not supported, therefore we don't need to set the OnAcked callback yet.
				// PATH_CHALLENGE and PATH_RESPONSE are never retransmitted.
			default:
				// we might be packing a 0-RTT packet, but we need to use the 1-RTT ack handler anyway
				pl.frames[i].Handler = p.retransmissionQueue.AckHandler(protocol.Encryption1RTT)
			}
		}
	}
	return pl
}

func (p *packetPacker) PackPTOProbePacket(
	encLevel protocol.EncryptionLevel,
	maxPacketSize protocol.ByteCount,
	addPingIfEmpty bool,
	now monotime.Time,
	v protocol.Version,
) (*coalescedPacket, error) {
	if encLevel == protocol.Encryption1RTT {
		return p.packPTOProbePacket1RTT(maxPacketSize, addPingIfEmpty, now, v)
	}

	var sealer handshake.LongHeaderSealer
	//nolint:exhaustive // Probe packets are never sent for 0-RTT.
	switch encLevel {
	case protocol.EncryptionInitial:
		var err error
		sealer, err = p.cryptoSetup.GetInitialSealer()
		if err != nil {
			return nil, err
		}
	case protocol.EncryptionHandshake:
		var err error
		sealer, err = p.cryptoSetup.GetHandshakeSealer()
		if err != nil {
			return nil, err
		}
	default:
		panic("unknown encryption level")
	}
	hdr, pl := p.maybeGetCryptoPacket(
		maxPacketSize-protocol.ByteCount(sealer.Overhead()),
		encLevel,
		now,
		addPingIfEmpty,
		false,
		true,
		v,
	)
	if pl.length == 0 {
		return nil, nil
	}
	buffer := getPacketBuffer()
	packet := &coalescedPacket{buffer: buffer}
	size := p.longHeaderPacketLength(hdr, pl, v) + protocol.ByteCount(sealer.Overhead())
	var padding protocol.ByteCount
	if encLevel == protocol.EncryptionInitial {
		padding = p.initialPaddingLen(pl.frames, size, maxPacketSize)
	}

	longHdrPacket, err := p.appendLongHeaderPacket(buffer, hdr, pl, padding, encLevel, sealer, v)
	if err != nil {
		return nil, err
	}
	packet.longHdrPackets = []*longHeaderPacket{longHdrPacket}
	return packet, nil
}

func (p *packetPacker) packPTOProbePacket1RTT(maxPacketSize protocol.ByteCount, addPingIfEmpty bool, now monotime.Time, v protocol.Version) (*coalescedPacket, error) {
	s, err := p.cryptoSetup.Get1RTTSealer()
	if err != nil {
		return nil, err
	}
	kp := s.KeyPhase()
	connID := p.getDestConnID()
	pn, pnLen := p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
	hdrLen := wire.ShortHeaderLen(connID, pnLen)
	pl := p.maybeGetAppDataPacket(maxPacketSize-protocol.ByteCount(s.Overhead())-hdrLen, false, true, now, v)
	if pl.length == 0 {
		if !addPingIfEmpty {
			return nil, nil
		}
		ping := &wire.PingFrame{}
		pl.frames = append(pl.frames, ackhandler.Frame{Frame: ping, Handler: emptyHandler{}})
		pl.length += ping.Length(v)
	}
	buffer := getPacketBuffer()
	packet := &coalescedPacket{buffer: buffer}
	shp, err := p.appendShortHeaderPacket(buffer, connID, pn, pnLen, kp, pl, 0, maxPacketSize, s, false, v)
	if err != nil {
		return nil, err
	}
	packet.shortHdrPacket = &shp
	return packet, nil
}

func (p *packetPacker) PackMTUProbePacket(ping ackhandler.Frame, size protocol.ByteCount, v protocol.Version) (shortHeaderPacket, *packetBuffer, error) {
	pl := payload{
		frames: []ackhandler.Frame{ping},
		length: ping.Frame.Length(v),
	}
	buffer := getPacketBuffer()
	s, err := p.cryptoSetup.Get1RTTSealer()
	if err != nil {
		return shortHeaderPacket{}, nil, err
	}
	connID := p.getDestConnID()
	pn, pnLen := p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
	padding := size - p.shortHeaderPacketLength(connID, pnLen, pl) - protocol.ByteCount(s.Overhead())
	kp := s.KeyPhase()
	packet, err := p.appendShortHeaderPacket(buffer, connID, pn, pnLen, kp, pl, padding, size, s, true, v)
	return packet, buffer, err
}

func (p *packetPacker) PackPathProbePacket(connID protocol.ConnectionID, frames []ackhandler.Frame, v protocol.Version) (shortHeaderPacket, *packetBuffer, error) {
	pn, pnLen := p.pnManager.PeekPacketNumber(protocol.Encryption1RTT)
	buf := getPacketBuffer()
	s, err := p.cryptoSetup.Get1RTTSealer()
	if err != nil {
		return shortHeaderPacket{}, nil, err
	}
	var l protocol.ByteCount
	for _, f := range frames {
		l += f.Frame.Length(v)
	}
	payload := payload{
		frames: frames,
		length: l,
	}
	padding := protocol.MinInitialPacketSize - p.shortHeaderPacketLength(connID, pnLen, payload) - protocol.ByteCount(s.Overhead())
	packet, err := p.appendShortHeaderPacket(buf, connID, pn, pnLen, s.KeyPhase(), payload, padding, protocol.MinInitialPacketSize, s, false, v)
	if err != nil {
		return shortHeaderPacket{}, nil, err
	}
	packet.IsPathProbePacket = true
	return packet, buf, err
}

func (p *packetPacker) getLongHeader(encLevel protocol.EncryptionLevel, v protocol.Version) *wire.ExtendedHeader {
	pn, pnLen := p.pnManager.PeekPacketNumber(encLevel)
	hdr := &wire.ExtendedHeader{
		PacketNumber:    pn,
		PacketNumberLen: pnLen,
	}
	hdr.Version = v
	hdr.SrcConnectionID = p.srcConnID
	hdr.DestConnectionID = p.getDestConnID()

	//nolint:exhaustive // 1-RTT packets are not long header packets.
	switch encLevel {
	case protocol.EncryptionInitial:
		hdr.Type = protocol.PacketTypeInitial
		hdr.Token = p.token
	case protocol.EncryptionHandshake:
		hdr.Type = protocol.PacketTypeHandshake
	case protocol.Encryption0RTT:
		hdr.Type = protocol.PacketType0RTT
	}
	return hdr
}

func (p *packetPacker) appendLongHeaderPacket(buffer *packetBuffer, header *wire.ExtendedHeader, pl payload, padding protocol.ByteCount, encLevel protocol.EncryptionLevel, sealer sealer, v protocol.Version) (*longHeaderPacket, error) {
	var paddingLen protocol.ByteCount
	pnLen := protocol.ByteCount(header.PacketNumberLen)
	if pl.length < 4-pnLen {
		paddingLen = 4 - pnLen - pl.length
	}
	paddingLen += padding
	header.Length = pnLen + protocol.ByteCount(sealer.Overhead()) + pl.length + paddingLen

	startLen := len(buffer.Data)
	raw := buffer.Data[startLen:]
	raw, err := header.Append(raw, v)
	if err != nil {
		return nil, err
	}
	payloadOffset := protocol.ByteCount(len(raw))

	raw, err = p.appendPacketPayload(raw, pl, paddingLen, v)
	if err != nil {
		return nil, err
	}
	raw = p.encryptPacket(raw, sealer, header.PacketNumber, payloadOffset, pnLen)
	buffer.Data = buffer.Data[:len(buffer.Data)+len(raw)]

	if pn := p.pnManager.PopPacketNumber(encLevel); pn != header.PacketNumber {
		return nil, fmt.Errorf("packetPacker BUG: Peeked and Popped packet numbers do not match: expected %d, got %d", pn, header.PacketNumber)
	}
	return &longHeaderPacket{
		header:       header,
		ack:          pl.ack,
		frames:       pl.frames,
		streamFrames: pl.streamFrames,
		length:       protocol.ByteCount(len(raw)),
	}, nil
}

func (p *packetPacker) appendShortHeaderPacket(
	buffer *packetBuffer,
	connID protocol.ConnectionID,
	pn protocol.PacketNumber,
	pnLen protocol.PacketNumberLen,
	kp protocol.KeyPhaseBit,
	pl payload,
	padding, maxPacketSize protocol.ByteCount,
	sealer sealer,
	isMTUProbePacket bool,
	v protocol.Version,
) (shortHeaderPacket, error) {
	var paddingLen protocol.ByteCount
	if pl.length < 4-protocol.ByteCount(pnLen) {
		paddingLen = 4 - protocol.ByteCount(pnLen) - pl.length
	}
	paddingLen += padding

	startLen := len(buffer.Data)
	raw := buffer.Data[startLen:]
	raw, err := wire.AppendShortHeader(raw, connID, pn, pnLen, kp)
	if err != nil {
		return shortHeaderPacket{}, err
	}
	payloadOffset := protocol.ByteCount(len(raw))

	raw, err = p.appendPacketPayload(raw, pl, paddingLen, v)
	if err != nil {
		return shortHeaderPacket{}, err
	}
	if !isMTUProbePacket {
		if size := protocol.ByteCount(len(raw) + sealer.Overhead()); size > maxPacketSize {
			return shortHeaderPacket{}, fmt.Errorf("PacketPacker BUG: packet too large (%d bytes, allowed %d bytes)", size, maxPacketSize)
		}
	}
	raw = p.encryptPacket(raw, sealer, pn, payloadOffset, protocol.ByteCount(pnLen))
	buffer.Data = buffer.Data[:len(buffer.Data)+len(raw)]

	if newPN := p.pnManager.PopPacketNumber(protocol.Encryption1RTT); newPN != pn {
		return shortHeaderPacket{}, fmt.Errorf("packetPacker BUG: Peeked and Popped packet numbers do not match: expected %d, got %d", pn, newPN)
	}
	return shortHeaderPacket{
		PacketNumber:         pn,
		PacketNumberLen:      pnLen,
		KeyPhase:             kp,
		StreamFrames:         pl.streamFrames,
		Frames:               pl.frames,
		Ack:                  pl.ack,
		Length:               protocol.ByteCount(len(raw)),
		DestConnID:           connID,
		IsPathMTUProbePacket: isMTUProbePacket,
	}, nil
}

// appendPacketPayload serializes the payload of a packet into the raw byte slice.
// It modifies the order of payload.frames.
func (p *packetPacker) appendPacketPayload(raw []byte, pl payload, paddingLen protocol.ByteCount, v protocol.Version) ([]byte, error) {
	payloadOffset := len(raw)
	if pl.ack != nil {
		var err error
		raw, err = pl.ack.Append(raw, v)
		if err != nil {
			return nil, err
		}
	}
	if paddingLen > 0 {
		raw = append(raw, make([]byte, paddingLen)...)
	}
	// Randomize the order of the control frames.
	// This makes sure that the receiver doesn't rely on the order in which frames are packed.
	if len(pl.frames) > 1 {
		p.rand.Shuffle(len(pl.frames), func(i, j int) { pl.frames[i], pl.frames[j] = pl.frames[j], pl.frames[i] })
	}
	for _, f := range pl.frames {
		var err error
		raw, err = f.Frame.Append(raw, v)
		if err != nil {
			return nil, err
		}
	}
	for _, f := range pl.streamFrames {
		var err error
		raw, err = f.Frame.Append(raw, v)
		if err != nil {
			return nil, err
		}
	}

	if payloadSize := protocol.ByteCount(len(raw)-payloadOffset) - paddingLen; payloadSize != pl.length {
		return nil, fmt.Errorf("PacketPacker BUG: payload size inconsistent (expected %d, got %d bytes)", pl.length, payloadSize)
	}
	return raw, nil
}

func (p *packetPacker) encryptPacket(raw []byte, sealer sealer, pn protocol.PacketNumber, payloadOffset, pnLen protocol.ByteCount) []byte {
	_ = sealer.Seal(raw[payloadOffset:payloadOffset], raw[payloadOffset:], pn, raw[:payloadOffset])
	raw = raw[:len(raw)+sealer.Overhead()]
	// apply header protection
	pnOffset := payloadOffset - pnLen
	sealer.EncryptHeader(raw[pnOffset+4:pnOffset+4+16], &raw[0], raw[pnOffset:payloadOffset])
	return raw
}

func (p *packetPacker) SetToken(token []byte) {
	p.token = token
}

type emptyHandler struct{}

var _ ackhandler.FrameHandler = emptyHandler{}

func (emptyHandler) OnAcked(wire.Frame) {}
func (emptyHandler) OnLost(wire.Frame)  {}
