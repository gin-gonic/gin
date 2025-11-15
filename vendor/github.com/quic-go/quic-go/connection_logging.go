package quic

import (
	"net"
	"net/netip"
	"slices"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlog"
)

// ConvertFrame converts a wire.Frame into a logging.Frame.
// This makes it possible for external packages to access the frames.
// Furthermore, it removes the data slices from CRYPTO and STREAM frames.
func toQlogFrame(frame wire.Frame) qlog.Frame {
	switch f := frame.(type) {
	case *wire.AckFrame:
		// We use a pool for ACK frames.
		// Implementations of the tracer interface may hold on to frames, so we need to make a copy here.
		return qlog.Frame{Frame: toQlogAckFrame(f)}
	case *wire.CryptoFrame:
		return qlog.Frame{
			Frame: &qlog.CryptoFrame{
				Offset: int64(f.Offset),
				Length: int64(len(f.Data)),
			},
		}
	case *wire.StreamFrame:
		return qlog.Frame{
			Frame: &qlog.StreamFrame{
				StreamID: f.StreamID,
				Offset:   int64(f.Offset),
				Length:   int64(f.DataLen()),
				Fin:      f.Fin,
			},
		}
	case *wire.DatagramFrame:
		return qlog.Frame{
			Frame: &qlog.DatagramFrame{
				Length: int64(len(f.Data)),
			},
		}
	default:
		return qlog.Frame{Frame: frame}
	}
}

func toQlogAckFrame(f *wire.AckFrame) *qlog.AckFrame {
	ack := &qlog.AckFrame{
		AckRanges: slices.Clone(f.AckRanges),
		DelayTime: f.DelayTime,
		ECNCE:     f.ECNCE,
		ECT0:      f.ECT0,
		ECT1:      f.ECT1,
	}
	return ack
}

func (c *Conn) logLongHeaderPacket(p *longHeaderPacket, ecn protocol.ECN) {
	// quic-go logging
	if c.logger.Debug() {
		p.header.Log(c.logger)
		if p.ack != nil {
			wire.LogFrame(c.logger, p.ack, true)
		}
		for _, frame := range p.frames {
			wire.LogFrame(c.logger, frame.Frame, true)
		}
		for _, frame := range p.streamFrames {
			wire.LogFrame(c.logger, frame.Frame, true)
		}
	}

	// tracing
	if c.qlogger != nil {
		numFrames := len(p.frames) + len(p.streamFrames)
		if p.ack != nil {
			numFrames++
		}
		frames := make([]qlog.Frame, 0, numFrames)
		if p.ack != nil {
			frames = append(frames, toQlogFrame(p.ack))
		}
		for _, f := range p.frames {
			frames = append(frames, toQlogFrame(f.Frame))
		}
		for _, f := range p.streamFrames {
			frames = append(frames, toQlogFrame(f.Frame))
		}
		c.qlogger.RecordEvent(qlog.PacketSent{
			Header: qlog.PacketHeader{
				PacketType:       toQlogPacketType(p.header.Type),
				KeyPhaseBit:      p.header.KeyPhase,
				PacketNumber:     p.header.PacketNumber,
				Version:          p.header.Version,
				SrcConnectionID:  p.header.SrcConnectionID,
				DestConnectionID: p.header.DestConnectionID,
			},
			Raw: qlog.RawInfo{
				Length:        int(p.length),
				PayloadLength: int(p.header.Length),
			},
			Frames: frames,
			ECN:    toQlogECN(ecn),
		})
	}
}

func (c *Conn) logShortHeaderPacket(
	destConnID protocol.ConnectionID,
	ackFrame *wire.AckFrame,
	frames []ackhandler.Frame,
	streamFrames []ackhandler.StreamFrame,
	pn protocol.PacketNumber,
	pnLen protocol.PacketNumberLen,
	kp protocol.KeyPhaseBit,
	ecn protocol.ECN,
	size protocol.ByteCount,
	isCoalesced bool,
) {
	if c.logger.Debug() && !isCoalesced {
		c.logger.Debugf("-> Sending packet %d (%d bytes) for connection %s, 1-RTT (ECN: %s)", pn, size, c.logID, ecn)
	}
	// quic-go logging
	if c.logger.Debug() {
		wire.LogShortHeader(c.logger, destConnID, pn, pnLen, kp)
		if ackFrame != nil {
			wire.LogFrame(c.logger, ackFrame, true)
		}
		for _, f := range frames {
			wire.LogFrame(c.logger, f.Frame, true)
		}
		for _, f := range streamFrames {
			wire.LogFrame(c.logger, f.Frame, true)
		}
	}

	// tracing
	if c.qlogger != nil {
		numFrames := len(frames) + len(streamFrames)
		if ackFrame != nil {
			numFrames++
		}
		fs := make([]qlog.Frame, 0, numFrames)
		if ackFrame != nil {
			fs = append(fs, toQlogFrame(ackFrame))
		}
		for _, f := range frames {
			fs = append(fs, toQlogFrame(f.Frame))
		}
		for _, f := range streamFrames {
			fs = append(fs, toQlogFrame(f.Frame))
		}
		c.qlogger.RecordEvent(qlog.PacketSent{
			Header: qlog.PacketHeader{
				PacketType:       qlog.PacketType1RTT,
				KeyPhaseBit:      kp,
				PacketNumber:     pn,
				Version:          c.version,
				DestConnectionID: destConnID,
			},
			Raw: qlog.RawInfo{
				Length:        int(size),
				PayloadLength: int(size - wire.ShortHeaderLen(destConnID, pnLen)),
			},
			Frames: fs,
			ECN:    toQlogECN(ecn),
		})
	}
}

func (c *Conn) logCoalescedPacket(packet *coalescedPacket, ecn protocol.ECN) {
	if c.logger.Debug() {
		// There's a short period between dropping both Initial and Handshake keys and completion of the handshake,
		// during which we might call PackCoalescedPacket but just pack a short header packet.
		if len(packet.longHdrPackets) == 0 && packet.shortHdrPacket != nil {
			c.logShortHeaderPacket(
				packet.shortHdrPacket.DestConnID,
				packet.shortHdrPacket.Ack,
				packet.shortHdrPacket.Frames,
				packet.shortHdrPacket.StreamFrames,
				packet.shortHdrPacket.PacketNumber,
				packet.shortHdrPacket.PacketNumberLen,
				packet.shortHdrPacket.KeyPhase,
				ecn,
				packet.shortHdrPacket.Length,
				false,
			)
			return
		}
		if len(packet.longHdrPackets) > 1 {
			c.logger.Debugf("-> Sending coalesced packet (%d parts, %d bytes) for connection %s", len(packet.longHdrPackets), packet.buffer.Len(), c.logID)
		} else {
			c.logger.Debugf("-> Sending packet %d (%d bytes) for connection %s, %s", packet.longHdrPackets[0].header.PacketNumber, packet.buffer.Len(), c.logID, packet.longHdrPackets[0].EncryptionLevel())
		}
	}
	for _, p := range packet.longHdrPackets {
		c.logLongHeaderPacket(p, ecn)
	}
	if p := packet.shortHdrPacket; p != nil {
		c.logShortHeaderPacket(p.DestConnID, p.Ack, p.Frames, p.StreamFrames, p.PacketNumber, p.PacketNumberLen, p.KeyPhase, ecn, p.Length, true)
	}
}

func (c *Conn) qlogTransportParameters(tp *wire.TransportParameters, sentBy protocol.Perspective, restore bool) {
	ev := qlog.ParametersSet{
		Restore:                         restore,
		OriginalDestinationConnectionID: tp.OriginalDestinationConnectionID,
		InitialSourceConnectionID:       tp.InitialSourceConnectionID,
		RetrySourceConnectionID:         tp.RetrySourceConnectionID,
		StatelessResetToken:             tp.StatelessResetToken,
		DisableActiveMigration:          tp.DisableActiveMigration,
		MaxIdleTimeout:                  tp.MaxIdleTimeout,
		MaxUDPPayloadSize:               tp.MaxUDPPayloadSize,
		AckDelayExponent:                tp.AckDelayExponent,
		MaxAckDelay:                     tp.MaxAckDelay,
		ActiveConnectionIDLimit:         tp.ActiveConnectionIDLimit,
		InitialMaxData:                  tp.InitialMaxData,
		InitialMaxStreamDataBidiLocal:   tp.InitialMaxStreamDataBidiLocal,
		InitialMaxStreamDataBidiRemote:  tp.InitialMaxStreamDataBidiRemote,
		InitialMaxStreamDataUni:         tp.InitialMaxStreamDataUni,
		InitialMaxStreamsBidi:           int64(tp.MaxBidiStreamNum),
		InitialMaxStreamsUni:            int64(tp.MaxUniStreamNum),
		MaxDatagramFrameSize:            tp.MaxDatagramFrameSize,
		EnableResetStreamAt:             tp.EnableResetStreamAt,
	}
	if sentBy == c.perspective {
		ev.Initiator = qlog.InitiatorLocal
	} else {
		ev.Initiator = qlog.InitiatorRemote
	}
	if tp.PreferredAddress != nil {
		ev.PreferredAddress = &qlog.PreferredAddress{
			IPv4:                tp.PreferredAddress.IPv4,
			IPv6:                tp.PreferredAddress.IPv6,
			ConnectionID:        tp.PreferredAddress.ConnectionID,
			StatelessResetToken: tp.PreferredAddress.StatelessResetToken,
		}
	}
	c.qlogger.RecordEvent(ev)
}

func toQlogECN(ecn protocol.ECN) qlog.ECN {
	//nolint:exhaustive // only need to handle the 3 valid values
	switch ecn {
	case protocol.ECT0:
		return qlog.ECT0
	case protocol.ECT1:
		return qlog.ECT1
	case protocol.ECNCE:
		return qlog.ECNCE
	default:
		return qlog.ECNUnsupported
	}
}

func toQlogPacketType(pt protocol.PacketType) qlog.PacketType {
	var qpt qlog.PacketType
	switch pt {
	case protocol.PacketTypeInitial:
		qpt = qlog.PacketTypeInitial
	case protocol.PacketTypeHandshake:
		qpt = qlog.PacketTypeHandshake
	case protocol.PacketType0RTT:
		qpt = qlog.PacketType0RTT
	case protocol.PacketTypeRetry:
		qpt = qlog.PacketTypeRetry
	}
	return qpt
}

func toPathEndpointInfo(addr *net.UDPAddr) qlog.PathEndpointInfo {
	if addr == nil {
		return qlog.PathEndpointInfo{}
	}

	var info qlog.PathEndpointInfo
	if addr.IP == nil || addr.IP.To4() != nil {
		addrPort := netip.AddrPortFrom(netip.AddrFrom4([4]byte(addr.IP.To4())), uint16(addr.Port))
		if addrPort.IsValid() {
			info.IPv4 = addrPort
		}
	} else {
		addrPort := netip.AddrPortFrom(netip.AddrFrom16([16]byte(addr.IP.To16())), uint16(addr.Port))
		if addrPort.IsValid() {
			info.IPv6 = addrPort
		}
	}
	return info
}

// startedConnectionEvent builds a StartedConnection event using consistent logic
// for both endpoints. If the local address is unspecified (e.g., dual-stack
// listener), it selects the family based on the remote address and uses the
// unspecified address of that family with the local port.
func startedConnectionEvent(local, remote *net.UDPAddr) qlog.StartedConnection {
	var localInfo, remoteInfo qlog.PathEndpointInfo
	if remote != nil {
		remoteInfo = toPathEndpointInfo(remote)
	}
	if local != nil {
		if local.IP == nil || local.IP.IsUnspecified() {
			// Choose local family based on the remote address family.
			if remote != nil && remote.IP.To4() != nil {
				ap := netip.AddrPortFrom(netip.AddrFrom4([4]byte{}), uint16(local.Port))
				if ap.IsValid() {
					localInfo.IPv4 = ap
				}
			} else if remote != nil && remote.IP.To16() != nil && remote.IP.To4() == nil {
				ap := netip.AddrPortFrom(netip.AddrFrom16([16]byte{}), uint16(local.Port))
				if ap.IsValid() {
					localInfo.IPv6 = ap
				}
			}
		} else {
			localInfo = toPathEndpointInfo(local)
		}
	}
	return qlog.StartedConnection{Local: localInfo, Remote: remoteInfo}
}
