package ackhandler

import (
	"iter"
	"slices"

	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
)

type lostPacket struct {
	PacketNumber protocol.PacketNumber
	SendTime     monotime.Time
}

type lostPacketTracker struct {
	maxLength   int
	lostPackets []lostPacket
}

func newLostPacketTracker(maxLength int) *lostPacketTracker {
	return &lostPacketTracker{
		maxLength: maxLength,
		// Preallocate a small slice only.
		// Hopefully we won't lose many packets.
		lostPackets: make([]lostPacket, 0, 4),
	}
}

func (t *lostPacketTracker) Add(p protocol.PacketNumber, sendTime monotime.Time) {
	if len(t.lostPackets) == t.maxLength {
		t.lostPackets = t.lostPackets[1:]
	}
	t.lostPackets = append(t.lostPackets, lostPacket{
		PacketNumber: p,
		SendTime:     sendTime,
	})
}

// Delete deletes a packet from the lost packet tracker.
// This function is not optimized for performance if many packets are lost,
// but it is only used when a spurious loss is detected, which is rare.
func (t *lostPacketTracker) Delete(pn protocol.PacketNumber) {
	t.lostPackets = slices.DeleteFunc(t.lostPackets, func(p lostPacket) bool {
		return p.PacketNumber == pn
	})
}

func (t *lostPacketTracker) All() iter.Seq2[protocol.PacketNumber, monotime.Time] {
	return func(yield func(protocol.PacketNumber, monotime.Time) bool) {
		for _, p := range t.lostPackets {
			if !yield(p.PacketNumber, p.SendTime) {
				return
			}
		}
	}
}

func (t *lostPacketTracker) DeleteBefore(ti monotime.Time) {
	if len(t.lostPackets) == 0 {
		return
	}
	if !t.lostPackets[0].SendTime.Before(ti) {
		return
	}
	var idx int
	for ; idx < len(t.lostPackets); idx++ {
		if !t.lostPackets[idx].SendTime.Before(ti) {
			break
		}
	}
	t.lostPackets = slices.Delete(t.lostPackets, 0, idx)
}
