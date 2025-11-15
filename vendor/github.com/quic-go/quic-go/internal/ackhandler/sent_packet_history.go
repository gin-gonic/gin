package ackhandler

import (
	"fmt"
	"iter"
	"slices"

	"github.com/quic-go/quic-go/internal/protocol"
)

const maxSkippedPackets = 4

type sentPacketHistory struct {
	packets          []*packet
	pathProbePackets []packetWithPacketNumber
	skippedPackets   []protocol.PacketNumber

	numOutstanding int

	firstPacketNumber   protocol.PacketNumber
	highestPacketNumber protocol.PacketNumber
}

func newSentPacketHistory(isAppData bool) *sentPacketHistory {
	h := &sentPacketHistory{
		highestPacketNumber: protocol.InvalidPacketNumber,
		firstPacketNumber:   protocol.InvalidPacketNumber,
	}
	if isAppData {
		h.packets = make([]*packet, 0, 32)
		h.skippedPackets = make([]protocol.PacketNumber, 0, maxSkippedPackets)
	} else {
		h.packets = make([]*packet, 0, 6)
	}
	return h
}

func (h *sentPacketHistory) checkSequentialPacketNumberUse(pn protocol.PacketNumber) {
	if h.highestPacketNumber != protocol.InvalidPacketNumber {
		if pn != h.highestPacketNumber+1 {
			panic("non-sequential packet number use")
		}
	}
	h.highestPacketNumber = pn
	if len(h.packets) == 0 {
		h.firstPacketNumber = pn
	}
}

func (h *sentPacketHistory) SkippedPacket(pn protocol.PacketNumber) {
	h.checkSequentialPacketNumberUse(pn)
	if len(h.packets) > 0 {
		h.packets = append(h.packets, nil)
	}
	if len(h.skippedPackets) == maxSkippedPackets {
		h.skippedPackets = slices.Delete(h.skippedPackets, 0, 1)
	}
	h.skippedPackets = append(h.skippedPackets, pn)
}

func (h *sentPacketHistory) SentNonAckElicitingPacket(pn protocol.PacketNumber) {
	h.checkSequentialPacketNumberUse(pn)
	if len(h.packets) > 0 {
		h.packets = append(h.packets, nil)
	}
}

func (h *sentPacketHistory) SentAckElicitingPacket(pn protocol.PacketNumber, p *packet) {
	h.checkSequentialPacketNumberUse(pn)
	h.packets = append(h.packets, p)
	if p.outstanding() {
		h.numOutstanding++
	}
}

func (h *sentPacketHistory) SentPathProbePacket(pn protocol.PacketNumber, p *packet) {
	h.checkSequentialPacketNumberUse(pn)
	h.packets = append(h.packets, &packet{isPathProbePacket: true})
	h.pathProbePackets = append(h.pathProbePackets, packetWithPacketNumber{PacketNumber: pn, packet: p})
}

func (h *sentPacketHistory) Packets() iter.Seq2[protocol.PacketNumber, *packet] {
	return func(yield func(protocol.PacketNumber, *packet) bool) {
		// h.firstPacketNumber might be updated in the yield function,
		// so we need to save it here.
		firstPacketNumber := h.firstPacketNumber
		for i, p := range h.packets {
			if p == nil {
				continue
			}
			if !yield(firstPacketNumber+protocol.PacketNumber(i), p) {
				return
			}
		}
	}
}

func (h *sentPacketHistory) PathProbes() iter.Seq2[protocol.PacketNumber, *packet] {
	return func(yield func(protocol.PacketNumber, *packet) bool) {
		for _, p := range h.pathProbePackets {
			if !yield(p.PacketNumber, p.packet) {
				return
			}
		}
	}
}

// FirstOutstanding returns the first outstanding packet.
func (h *sentPacketHistory) FirstOutstanding() (protocol.PacketNumber, *packet) {
	if !h.HasOutstandingPackets() {
		return protocol.InvalidPacketNumber, nil
	}
	for i, p := range h.packets {
		if p != nil && p.outstanding() {
			return h.firstPacketNumber + protocol.PacketNumber(i), p
		}
	}
	return protocol.InvalidPacketNumber, nil
}

// FirstOutstandingPathProbe returns the first outstanding path probe packet
func (h *sentPacketHistory) FirstOutstandingPathProbe() (protocol.PacketNumber, *packet) {
	if len(h.pathProbePackets) == 0 {
		return protocol.InvalidPacketNumber, nil
	}
	return h.pathProbePackets[0].PacketNumber, h.pathProbePackets[0].packet
}

func (h *sentPacketHistory) SkippedPackets() iter.Seq[protocol.PacketNumber] {
	return func(yield func(protocol.PacketNumber) bool) {
		for _, p := range h.skippedPackets {
			if !yield(p) {
				return
			}
		}
	}
}

func (h *sentPacketHistory) Len() int {
	return len(h.packets)
}

// Remove removes a packet from the sent packet history.
// It must not be used for skipped packet numbers.
func (h *sentPacketHistory) Remove(pn protocol.PacketNumber) error {
	idx, ok := h.getIndex(pn)
	if !ok {
		return fmt.Errorf("packet %d not found in sent packet history", pn)
	}
	p := h.packets[idx]
	if p.outstanding() {
		h.numOutstanding--
		if h.numOutstanding < 0 {
			panic("negative number of outstanding packets")
		}
	}
	h.packets[idx] = nil
	// clean up all skipped packets directly before this packet number
	var hasPacketBefore bool
	for idx > 0 {
		idx--
		if h.packets[idx] != nil {
			hasPacketBefore = true
			break
		}
	}
	if !hasPacketBefore {
		h.cleanupStart()
	}
	if len(h.packets) > 0 && h.packets[0] == nil {
		panic("cleanup failed")
	}
	return nil
}

// RemovePathProbe removes a path probe packet.
// It scales O(N), but that's ok, since we don't expect to send many path probe packets.
// It is not valid to call this function in IteratePathProbes.
func (h *sentPacketHistory) RemovePathProbe(pn protocol.PacketNumber) *packet {
	var packetToDelete *packet
	idx := -1
	for i, p := range h.pathProbePackets {
		if p.PacketNumber == pn {
			packetToDelete = p.packet
			idx = i
			break
		}
	}
	if idx != -1 {
		// don't use slices.Delete, because it zeros the deleted element
		copy(h.pathProbePackets[idx:], h.pathProbePackets[idx+1:])
		h.pathProbePackets = h.pathProbePackets[:len(h.pathProbePackets)-1]
	}
	return packetToDelete
}

// getIndex gets the index of packet p in the packets slice.
func (h *sentPacketHistory) getIndex(p protocol.PacketNumber) (int, bool) {
	if len(h.packets) == 0 {
		return 0, false
	}
	if p < h.firstPacketNumber {
		return 0, false
	}
	index := int(p - h.firstPacketNumber)
	if index > len(h.packets)-1 {
		return 0, false
	}
	return index, true
}

func (h *sentPacketHistory) HasOutstandingPackets() bool {
	return h.numOutstanding > 0
}

func (h *sentPacketHistory) HasOutstandingPathProbes() bool {
	return len(h.pathProbePackets) > 0
}

// delete all nil entries at the beginning of the packets slice
func (h *sentPacketHistory) cleanupStart() {
	for i, p := range h.packets {
		if p != nil {
			h.packets = h.packets[i:]
			h.firstPacketNumber += protocol.PacketNumber(i)
			return
		}
	}
	h.packets = h.packets[:0]
	h.firstPacketNumber = protocol.InvalidPacketNumber
}

func (h *sentPacketHistory) LowestPacketNumber() protocol.PacketNumber {
	if len(h.packets) == 0 {
		return protocol.InvalidPacketNumber
	}
	return h.firstPacketNumber
}

func (h *sentPacketHistory) DeclareLost(pn protocol.PacketNumber) {
	idx, ok := h.getIndex(pn)
	if !ok {
		return
	}
	p := h.packets[idx]
	if p.outstanding() {
		h.numOutstanding--
		if h.numOutstanding < 0 {
			panic("negative number of outstanding packets")
		}
	}
	h.packets[idx] = nil
	if idx == 0 {
		h.cleanupStart()
	}
}

// Difference returns the difference between two packet numbers a and b (a - b),
// taking into account any skipped packet numbers between them.
//
// Note that old skipped packets are garbage collected at some point,
// so this function is not guaranteed to return the correct result after a while.
func (h *sentPacketHistory) Difference(a, b protocol.PacketNumber) protocol.PacketNumber {
	diff := a - b
	if len(h.skippedPackets) == 0 {
		return diff
	}
	if a < h.skippedPackets[0] || b > h.skippedPackets[len(h.skippedPackets)-1] {
		return diff
	}
	for _, p := range h.skippedPackets {
		if p > b && p < a {
			diff--
		}
	}
	return diff
}
