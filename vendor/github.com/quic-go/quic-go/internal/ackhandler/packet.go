package ackhandler

import (
	"sync"

	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
)

type packetWithPacketNumber struct {
	PacketNumber protocol.PacketNumber
	*packet
}

// A Packet is a packet
type packet struct {
	SendTime        monotime.Time
	StreamFrames    []StreamFrame
	Frames          []Frame
	LargestAcked    protocol.PacketNumber // InvalidPacketNumber if the packet doesn't contain an ACK
	Length          protocol.ByteCount
	EncryptionLevel protocol.EncryptionLevel

	IsPathMTUProbePacket bool // We don't report the loss of Path MTU probe packets to the congestion controller.

	includedInBytesInFlight bool
	declaredLost            bool
	isPathProbePacket       bool
}

func (p *packet) outstanding() bool {
	return !p.declaredLost && !p.IsPathMTUProbePacket && !p.isPathProbePacket
}

var packetPool = sync.Pool{New: func() any { return &packet{} }}

func getPacket() *packet {
	p := packetPool.Get().(*packet)
	p.StreamFrames = nil
	p.Frames = nil
	p.LargestAcked = 0
	p.Length = 0
	p.EncryptionLevel = protocol.EncryptionLevel(0)
	p.SendTime = 0
	p.IsPathMTUProbePacket = false
	p.includedInBytesInFlight = false
	p.declaredLost = false
	p.isPathProbePacket = false
	return p
}

// We currently only return Packets back into the pool when they're acknowledged (not when they're lost).
// This simplifies the code, and gives the vast majority of the performance benefit we can gain from using the pool.
func putPacket(p *packet) {
	p.Frames = nil
	p.StreamFrames = nil
	packetPool.Put(p)
}
