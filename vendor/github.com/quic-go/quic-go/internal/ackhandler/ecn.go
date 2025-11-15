package ackhandler

import (
	"fmt"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
)

type ecnState uint8

const (
	ecnStateInitial ecnState = iota
	ecnStateTesting
	ecnStateUnknown
	ecnStateCapable
	ecnStateFailed
)

const (
	// ecnFailedNoECNCounts is emitted when an ACK acknowledges ECN-marked packets,
	// but doesn't contain any ECN counts
	ecnFailedNoECNCounts = "ACK doesn't contain ECN marks"
	// ecnFailedDecreasedECNCounts is emitted when an ACK frame decreases ECN counts
	ecnFailedDecreasedECNCounts = "ACK decreases ECN counts"
	// ecnFailedLostAllTestingPackets is emitted when all ECN testing packets are declared lost
	ecnFailedLostAllTestingPackets = "all ECN testing packets declared lost"
	// ecnFailedMoreECNCountsThanSent is emitted when an ACK contains more ECN counts than ECN-marked packets were sent
	ecnFailedMoreECNCountsThanSent = "ACK contains more ECN counts than ECN-marked packets sent"
	// ecnFailedTooFewECNCounts is emitted when an ACK contains fewer ECN counts than it acknowledges packets
	ecnFailedTooFewECNCounts = "ACK contains fewer new ECN counts than acknowledged ECN-marked packets"
	// ecnFailedManglingDetected is emitted when the path marks all ECN-marked packets as CE
	ecnFailedManglingDetected = "ECN mangling detected"
)

// must fit into an uint8, otherwise numSentTesting and numLostTesting must have a larger type
const numECNTestingPackets = 10

type ecnHandler interface {
	SentPacket(protocol.PacketNumber, protocol.ECN)
	Mode() protocol.ECN
	HandleNewlyAcked(packets []packetWithPacketNumber, ect0, ect1, ecnce int64) (congested bool)
	LostPacket(protocol.PacketNumber)
}

// The ecnTracker performs ECN validation of a path.
// Once failed, it doesn't do any re-validation of the path.
// It is designed only work for 1-RTT packets, it doesn't handle multiple packet number spaces.
// In order to avoid revealing any internal state to on-path observers,
// callers should make sure to start using ECN (i.e. calling Mode) for the very first 1-RTT packet sent.
// The validation logic implemented here strictly follows the algorithm described in RFC 9000 section 13.4.2 and A.4.
type ecnTracker struct {
	state                          ecnState
	numSentTesting, numLostTesting uint8

	firstTestingPacket protocol.PacketNumber
	lastTestingPacket  protocol.PacketNumber
	firstCapablePacket protocol.PacketNumber

	numSentECT0, numSentECT1                  int64
	numAckedECT0, numAckedECT1, numAckedECNCE int64

	qlogger qlogwriter.Recorder
	logger  utils.Logger
}

var _ ecnHandler = &ecnTracker{}

func newECNTracker(logger utils.Logger, qlogger qlogwriter.Recorder) *ecnTracker {
	return &ecnTracker{
		firstTestingPacket: protocol.InvalidPacketNumber,
		lastTestingPacket:  protocol.InvalidPacketNumber,
		firstCapablePacket: protocol.InvalidPacketNumber,
		state:              ecnStateInitial,
		logger:             logger,
		qlogger:            qlogger,
	}
}

func (e *ecnTracker) SentPacket(pn protocol.PacketNumber, ecn protocol.ECN) {
	//nolint:exhaustive // These are the only ones we need to take care of.
	switch ecn {
	case protocol.ECNNon:
		return
	case protocol.ECT0:
		e.numSentECT0++
	case protocol.ECT1:
		e.numSentECT1++
	case protocol.ECNUnsupported:
		if e.state != ecnStateFailed {
			panic("didn't expect ECN to be unsupported")
		}
	default:
		panic(fmt.Sprintf("sent packet with unexpected ECN marking: %s", ecn))
	}

	if e.state == ecnStateCapable && e.firstCapablePacket == protocol.InvalidPacketNumber {
		e.firstCapablePacket = pn
	}

	if e.state != ecnStateTesting {
		return
	}

	e.numSentTesting++
	if e.firstTestingPacket == protocol.InvalidPacketNumber {
		e.firstTestingPacket = pn
	}
	if e.numSentECT0+e.numSentECT1 >= numECNTestingPackets {
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State: qlog.ECNStateUnknown,
			})
		}
		e.state = ecnStateUnknown
		e.lastTestingPacket = pn
	}
}

func (e *ecnTracker) Mode() protocol.ECN {
	switch e.state {
	case ecnStateInitial:
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State: qlog.ECNStateTesting,
			})
		}
		e.state = ecnStateTesting
		return e.Mode()
	case ecnStateTesting, ecnStateCapable:
		return protocol.ECT0
	case ecnStateUnknown, ecnStateFailed:
		return protocol.ECNNon
	default:
		panic(fmt.Sprintf("unknown ECN state: %d", e.state))
	}
}

func (e *ecnTracker) LostPacket(pn protocol.PacketNumber) {
	if e.state != ecnStateTesting && e.state != ecnStateUnknown {
		return
	}
	if !e.isTestingPacket(pn) {
		return
	}
	e.numLostTesting++
	// Only proceed if we have sent all 10 testing packets.
	if e.state != ecnStateUnknown {
		return
	}
	if e.numLostTesting >= e.numSentTesting {
		e.logger.Debugf("Disabling ECN. All testing packets were lost.")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedLostAllTestingPackets,
			})
		}
		e.state = ecnStateFailed
		return
	}
	// Path validation also fails if some testing packets are lost, and all other testing packets where CE-marked
	e.failIfMangled()
}

// HandleNewlyAcked handles the ECN counts on an ACK frame.
// It must only be called for ACK frames that increase the largest acknowledged packet number,
// see section 13.4.2.1 of RFC 9000.
func (e *ecnTracker) HandleNewlyAcked(packets []packetWithPacketNumber, ect0, ect1, ecnce int64) (congested bool) {
	if e.state == ecnStateFailed {
		return false
	}

	// ECN validation can fail if the received total count for either ECT(0) or ECT(1) exceeds
	// the total number of packets sent with each corresponding ECT codepoint.
	if ect0 > e.numSentECT0 || ect1 > e.numSentECT1 {
		e.logger.Debugf("Disabling ECN. Received more ECT(0) / ECT(1) acknowledgements than packets sent.")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedMoreECNCountsThanSent,
			})
		}
		e.state = ecnStateFailed
		return false
	}

	// Count ECT0 and ECT1 marks that we used when sending the packets that are now being acknowledged.
	var ackedECT0, ackedECT1 int64
	for _, p := range packets {
		//nolint:exhaustive // We only ever send ECT(0) and ECT(1).
		switch e.ecnMarking(p.PacketNumber) {
		case protocol.ECT0:
			ackedECT0++
		case protocol.ECT1:
			ackedECT1++
		}
	}

	// If an ACK frame newly acknowledges a packet that the endpoint sent with either the ECT(0) or ECT(1)
	// codepoint set, ECN validation fails if the corresponding ECN counts are not present in the ACK frame.
	// This check detects:
	// * paths that bleach all ECN marks, and
	// * peers that don't report any ECN counts
	if (ackedECT0 > 0 || ackedECT1 > 0) && ect0 == 0 && ect1 == 0 && ecnce == 0 {
		e.logger.Debugf("Disabling ECN. ECN-marked packet acknowledged, but no ECN counts on ACK frame.")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedNoECNCounts,
			})
		}
		e.state = ecnStateFailed
		return false
	}

	// Determine the increase in ECT0, ECT1 and ECNCE marks
	newECT0 := ect0 - e.numAckedECT0
	newECT1 := ect1 - e.numAckedECT1
	newECNCE := ecnce - e.numAckedECNCE

	// We're only processing ACKs that increase the Largest Acked.
	// Therefore, the ECN counters should only ever increase.
	// Any decrease means that the peer's counting logic is broken.
	if newECT0 < 0 || newECT1 < 0 || newECNCE < 0 {
		e.logger.Debugf("Disabling ECN. ECN counts decreased unexpectedly.")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedDecreasedECNCounts,
			})
		}
		e.state = ecnStateFailed
		return false
	}

	// ECN validation also fails if the sum of the increase in ECT(0) and ECN-CE counts is less than the number
	// of newly acknowledged packets that were originally sent with an ECT(0) marking.
	// This could be the result of (partial) bleaching.
	if newECT0+newECNCE < ackedECT0 {
		e.logger.Debugf("Disabling ECN. Received less ECT(0) + ECN-CE than packets sent with ECT(0).")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedTooFewECNCounts,
			})
		}
		e.state = ecnStateFailed
		return false
	}
	// Similarly, ECN validation fails if the sum of the increases to ECT(1) and ECN-CE counts is less than
	// the number of newly acknowledged packets sent with an ECT(1) marking.
	if newECT1+newECNCE < ackedECT1 {
		e.logger.Debugf("Disabling ECN. Received less ECT(1) + ECN-CE than packets sent with ECT(1).")
		if e.qlogger != nil {
			e.qlogger.RecordEvent(qlog.ECNStateUpdated{
				State:   qlog.ECNStateFailed,
				Trigger: ecnFailedTooFewECNCounts,
			})
		}
		e.state = ecnStateFailed
		return false
	}

	// update our counters
	e.numAckedECT0 = ect0
	e.numAckedECT1 = ect1
	e.numAckedECNCE = ecnce

	// Detect mangling (a path remarking all ECN-marked testing packets as CE),
	// once all 10 testing packets have been sent out.
	if e.state == ecnStateUnknown {
		e.failIfMangled()
		if e.state == ecnStateFailed {
			return false
		}
	}
	if e.state == ecnStateTesting || e.state == ecnStateUnknown {
		var ackedTestingPacket bool
		for _, p := range packets {
			if e.isTestingPacket(p.PacketNumber) {
				ackedTestingPacket = true
				break
			}
		}
		// This check won't succeed if the path is mangling ECN-marks (i.e. rewrites all ECN-marked packets to CE).
		if ackedTestingPacket && (newECT0 > 0 || newECT1 > 0) {
			e.logger.Debugf("ECN capability confirmed.")
			if e.qlogger != nil {
				e.qlogger.RecordEvent(qlog.ECNStateUpdated{
					State: qlog.ECNStateCapable,
				})
			}
			e.state = ecnStateCapable
		}
	}

	// Don't trust CE marks before having confirmed ECN capability of the path.
	// Otherwise, mangling would be misinterpreted as actual congestion.
	return e.state == ecnStateCapable && newECNCE > 0
}

// failIfMangled fails ECN validation if all testing packets are lost or CE-marked.
func (e *ecnTracker) failIfMangled() {
	numAckedECNCE := e.numAckedECNCE + int64(e.numLostTesting)
	if e.numSentECT0+e.numSentECT1 > numAckedECNCE {
		return
	}
	if e.qlogger != nil {
		e.qlogger.RecordEvent(qlog.ECNStateUpdated{
			State:   qlog.ECNStateFailed,
			Trigger: ecnFailedManglingDetected,
		})
	}
	e.state = ecnStateFailed
}

func (e *ecnTracker) ecnMarking(pn protocol.PacketNumber) protocol.ECN {
	if pn < e.firstTestingPacket || e.firstTestingPacket == protocol.InvalidPacketNumber {
		return protocol.ECNNon
	}
	if pn < e.lastTestingPacket || e.lastTestingPacket == protocol.InvalidPacketNumber {
		return protocol.ECT0
	}
	if pn < e.firstCapablePacket || e.firstCapablePacket == protocol.InvalidPacketNumber {
		return protocol.ECNNon
	}
	// We don't need to deal with the case when ECN validation fails,
	// since we're ignoring any ECN counts reported in ACK frames in that case.
	return protocol.ECT0
}

func (e *ecnTracker) isTestingPacket(pn protocol.PacketNumber) bool {
	if e.firstTestingPacket == protocol.InvalidPacketNumber {
		return false
	}
	return pn >= e.firstTestingPacket && (pn <= e.lastTestingPacket || e.lastTestingPacket == protocol.InvalidPacketNumber)
}
