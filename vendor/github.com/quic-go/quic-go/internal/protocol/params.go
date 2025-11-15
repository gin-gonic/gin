package protocol

import "time"

// DesiredReceiveBufferSize is the kernel UDP receive buffer size that we'd like to use.
const DesiredReceiveBufferSize = (1 << 20) * 7 // 7 MB

// DesiredSendBufferSize is the kernel UDP send buffer size that we'd like to use.
const DesiredSendBufferSize = (1 << 20) * 7 // 7 MB

// InitialPacketSize is the initial (before Path MTU discovery) maximum packet size used.
const InitialPacketSize = 1280

// MaxCongestionWindowPackets is the maximum congestion window in packet.
const MaxCongestionWindowPackets = 10000

// MaxUndecryptablePackets limits the number of undecryptable packets that are queued in the connection.
const MaxUndecryptablePackets = 32

// ConnectionFlowControlMultiplier determines how much larger the connection flow control windows needs to be relative to any stream's flow control window
// This is the value that Chromium is using
const ConnectionFlowControlMultiplier = 1.5

// DefaultInitialMaxStreamData is the default initial stream-level flow control window for receiving data
const DefaultInitialMaxStreamData = (1 << 10) * 512 // 512 kb

// DefaultInitialMaxData is the connection-level flow control window for receiving data
const DefaultInitialMaxData = ConnectionFlowControlMultiplier * DefaultInitialMaxStreamData

// DefaultMaxReceiveStreamFlowControlWindow is the default maximum stream-level flow control window for receiving data
const DefaultMaxReceiveStreamFlowControlWindow = 6 * (1 << 20) // 6 MB

// DefaultMaxReceiveConnectionFlowControlWindow is the default connection-level flow control window for receiving data
const DefaultMaxReceiveConnectionFlowControlWindow = 15 * (1 << 20) // 15 MB

// WindowUpdateThreshold is the fraction of the receive window that has to be consumed before an higher offset is advertised to the client
const WindowUpdateThreshold = 0.25

// DefaultMaxIncomingStreams is the maximum number of streams that a peer may open
const DefaultMaxIncomingStreams = 100

// DefaultMaxIncomingUniStreams is the maximum number of unidirectional streams that a peer may open
const DefaultMaxIncomingUniStreams = 100

// MaxServerUnprocessedPackets is the max number of packets stored in the server that are not yet processed.
const MaxServerUnprocessedPackets = 1024

// MaxConnUnprocessedPackets is the max number of packets stored in each connection that are not yet processed.
const MaxConnUnprocessedPackets = 256

// SkipPacketInitialPeriod is the initial period length used for packet number skipping to prevent an Optimistic ACK attack.
// Every time a packet number is skipped, the period is doubled, up to SkipPacketMaxPeriod.
const SkipPacketInitialPeriod PacketNumber = 256

// SkipPacketMaxPeriod is the maximum period length used for packet number skipping.
const SkipPacketMaxPeriod PacketNumber = 128 * 1024

// MaxAcceptQueueSize is the maximum number of connections that the server queues for accepting.
// If the queue is full, new connection attempts will be rejected.
const MaxAcceptQueueSize = 32

// TokenValidity is the duration that a (non-retry) token is considered valid
const TokenValidity = 24 * time.Hour

// MaxOutstandingSentPackets is maximum number of packets saved for retransmission.
// When reached, it imposes a soft limit on sending new packets:
// Sending ACKs and retransmission is still allowed, but now new regular packets can be sent.
const MaxOutstandingSentPackets = 2 * MaxCongestionWindowPackets

// MaxTrackedSentPackets is maximum number of sent packets saved for retransmission.
// When reached, no more packets will be sent.
// This value *must* be larger than MaxOutstandingSentPackets.
const MaxTrackedSentPackets = MaxOutstandingSentPackets * 5 / 4

// MaxNonAckElicitingAcks is the maximum number of packets containing an ACK,
// but no ack-eliciting frames, that we send in a row
const MaxNonAckElicitingAcks = 19

// MaxStreamFrameSorterGaps is the maximum number of gaps between received StreamFrames
// prevents DoS attacks against the streamFrameSorter
const MaxStreamFrameSorterGaps = 1000

// MinStreamFrameBufferSize is the minimum data length of a received STREAM frame
// that we use the buffer for. This protects against a DoS where an attacker would send us
// very small STREAM frames to consume a lot of memory.
const MinStreamFrameBufferSize = 128

// MinCoalescedPacketSize is the minimum size of a coalesced packet that we pack.
// If a packet has less than this number of bytes, we won't coalesce any more packets onto it.
const MinCoalescedPacketSize = 128

// MaxCryptoStreamOffset is the maximum offset allowed on any of the crypto streams.
// This limits the size of the ClientHello and Certificates that can be received.
const MaxCryptoStreamOffset = 16 * (1 << 10)

// MinRemoteIdleTimeout is the minimum value that we accept for the remote idle timeout
const MinRemoteIdleTimeout = 5 * time.Second

// DefaultIdleTimeout is the default idle timeout
const DefaultIdleTimeout = 30 * time.Second

// DefaultHandshakeIdleTimeout is the default idle timeout used before handshake completion.
const DefaultHandshakeIdleTimeout = 5 * time.Second

// MinStreamFrameSize is the minimum size that has to be left in a packet, so that we add another STREAM frame.
// This avoids splitting up STREAM frames into small pieces, which has 2 advantages:
// 1. it reduces the framing overhead
// 2. it reduces the head-of-line blocking, when a packet is lost
const MinStreamFrameSize ByteCount = 128

// MaxPostHandshakeCryptoFrameSize is the maximum size of CRYPTO frames
// we send after the handshake completes.
const MaxPostHandshakeCryptoFrameSize = 1000

// MaxAckFrameSize is the maximum size for an ACK frame that we write
// Due to the varint encoding, ACK frames can grow (almost) indefinitely large.
// The MaxAckFrameSize should be large enough to encode many ACK range,
// but must ensure that a maximum size ACK frame fits into one packet.
const MaxAckFrameSize ByteCount = 1000

// MaxNumAckRanges is the maximum number of ACK ranges that we send in an ACK frame.
// It also serves as a limit for the packet history.
// If at any point we keep track of more ranges, old ranges are discarded.
const MaxNumAckRanges = 32

// MinPacingDelay is the minimum duration that is used for packet pacing
// If the packet packing frequency is higher, multiple packets might be sent at once.
// Example: For a packet pacing delay of 200μs, we would send 5 packets at once, wait for 1ms, and so forth.
const MinPacingDelay = time.Millisecond

// DefaultConnectionIDLength is the connection ID length that is used for multiplexed connections
// if no other value is configured.
const DefaultConnectionIDLength = 4

// MaxActiveConnectionIDs is the number of connection IDs that we're storing.
const MaxActiveConnectionIDs = 4

// MaxIssuedConnectionIDs is the maximum number of connection IDs that we're issuing at the same time.
const MaxIssuedConnectionIDs = 6

// PacketsPerConnectionID is the number of packets we send using one connection ID.
// If the peer provices us with enough new connection IDs, we switch to a new connection ID.
const PacketsPerConnectionID = 10000

// AckDelayExponent is the ack delay exponent used when sending ACKs.
const AckDelayExponent = 3

// Estimated timer granularity.
// The loss detection timer will not be set to a value smaller than granularity.
const TimerGranularity = time.Millisecond

// MaxAckDelay is the maximum time by which we delay sending ACKs.
const MaxAckDelay = 25 * time.Millisecond

// MaxAckDelayInclGranularity is the max_ack_delay including the timer granularity.
// This is the value that should be advertised to the peer.
const MaxAckDelayInclGranularity = MaxAckDelay + TimerGranularity

// KeyUpdateInterval is the maximum number of packets we send or receive before initiating a key update.
const KeyUpdateInterval = 100 * 1000

// Max0RTTQueueingDuration is the maximum time that we store 0-RTT packets in order to wait for the corresponding Initial to be received.
const Max0RTTQueueingDuration = 100 * time.Millisecond

// Max0RTTQueues is the maximum number of connections that we buffer 0-RTT packets for.
const Max0RTTQueues = 32

// Max0RTTQueueLen is the maximum number of 0-RTT packets that we buffer for each connection.
// When a new connection is created, all buffered packets are passed to the connection immediately.
// To avoid blocking, this value has to be smaller than MaxConnUnprocessedPackets.
// To avoid packets being dropped as undecryptable by the connection, this value has to be smaller than MaxUndecryptablePackets.
const Max0RTTQueueLen = 31
