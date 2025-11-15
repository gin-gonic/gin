package ackhandler

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
)

type packetNumberGenerator interface {
	Peek() protocol.PacketNumber
	// Pop pops the packet number.
	// It reports if the packet number (before the one just popped) was skipped.
	// It never skips more than one packet number in a row.
	Pop() (skipped bool, _ protocol.PacketNumber)
}

type sequentialPacketNumberGenerator struct {
	next protocol.PacketNumber
}

var _ packetNumberGenerator = &sequentialPacketNumberGenerator{}

func newSequentialPacketNumberGenerator(initial protocol.PacketNumber) packetNumberGenerator {
	return &sequentialPacketNumberGenerator{next: initial}
}

func (p *sequentialPacketNumberGenerator) Peek() protocol.PacketNumber {
	return p.next
}

func (p *sequentialPacketNumberGenerator) Pop() (bool, protocol.PacketNumber) {
	next := p.next
	p.next++
	return false, next
}

// The skippingPacketNumberGenerator generates the packet number for the next packet
// it randomly skips a packet number every averagePeriod packets (on average).
// It is guaranteed to never skip two consecutive packet numbers.
type skippingPacketNumberGenerator struct {
	period    protocol.PacketNumber
	maxPeriod protocol.PacketNumber

	next       protocol.PacketNumber
	nextToSkip protocol.PacketNumber

	rng utils.Rand
}

var _ packetNumberGenerator = &skippingPacketNumberGenerator{}

func newSkippingPacketNumberGenerator(initial, initialPeriod, maxPeriod protocol.PacketNumber) packetNumberGenerator {
	g := &skippingPacketNumberGenerator{
		next:      initial,
		period:    initialPeriod,
		maxPeriod: maxPeriod,
	}
	g.generateNewSkip()
	return g
}

func (p *skippingPacketNumberGenerator) Peek() protocol.PacketNumber {
	if p.next == p.nextToSkip {
		return p.next + 1
	}
	return p.next
}

func (p *skippingPacketNumberGenerator) Pop() (bool, protocol.PacketNumber) {
	next := p.next
	if p.next == p.nextToSkip {
		next++
		p.next += 2
		p.generateNewSkip()
		return true, next
	}
	p.next++ // generate a new packet number for the next packet
	return false, next
}

func (p *skippingPacketNumberGenerator) generateNewSkip() {
	// make sure that there are never two consecutive packet numbers that are skipped
	p.nextToSkip = p.next + 3 + protocol.PacketNumber(p.rng.Int31n(int32(2*p.period)))
	p.period = min(2*p.period, p.maxPeriod)
}
