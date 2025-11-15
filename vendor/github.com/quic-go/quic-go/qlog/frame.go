package qlog

import (
	"encoding/hex"

	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlogwriter/jsontext"
)

type Frame struct {
	Frame any
}

type frames []Frame

type (
	// An AckFrame is an ACK frame.
	AckFrame = wire.AckFrame
	// A ConnectionCloseFrame is a CONNECTION_CLOSE frame.
	ConnectionCloseFrame = wire.ConnectionCloseFrame
	// A DataBlockedFrame is a DATA_BLOCKED frame.
	DataBlockedFrame = wire.DataBlockedFrame
	// A HandshakeDoneFrame is a HANDSHAKE_DONE frame.
	HandshakeDoneFrame = wire.HandshakeDoneFrame
	// A MaxDataFrame is a MAX_DATA frame.
	MaxDataFrame = wire.MaxDataFrame
	// A MaxStreamDataFrame is a MAX_STREAM_DATA frame.
	MaxStreamDataFrame = wire.MaxStreamDataFrame
	// A MaxStreamsFrame is a MAX_STREAMS_FRAME.
	MaxStreamsFrame = wire.MaxStreamsFrame
	// A NewConnectionIDFrame is a NEW_CONNECTION_ID frame.
	NewConnectionIDFrame = wire.NewConnectionIDFrame
	// A NewTokenFrame is a NEW_TOKEN frame.
	NewTokenFrame = wire.NewTokenFrame
	// A PathChallengeFrame is a PATH_CHALLENGE frame.
	PathChallengeFrame = wire.PathChallengeFrame
	// A PathResponseFrame is a PATH_RESPONSE frame.
	PathResponseFrame = wire.PathResponseFrame
	// A PingFrame is a PING frame.
	PingFrame = wire.PingFrame
	// A ResetStreamFrame is a RESET_STREAM frame.
	ResetStreamFrame = wire.ResetStreamFrame
	// A RetireConnectionIDFrame is a RETIRE_CONNECTION_ID frame.
	RetireConnectionIDFrame = wire.RetireConnectionIDFrame
	// A StopSendingFrame is a STOP_SENDING frame.
	StopSendingFrame = wire.StopSendingFrame
	// A StreamsBlockedFrame is a STREAMS_BLOCKED frame.
	StreamsBlockedFrame = wire.StreamsBlockedFrame
	// A StreamDataBlockedFrame is a STREAM_DATA_BLOCKED frame.
	StreamDataBlockedFrame = wire.StreamDataBlockedFrame
	// An AckFrequencyFrame is an ACK_FREQUENCY frame.
	AckFrequencyFrame = wire.AckFrequencyFrame
	// An ImmediateAckFrame is an IMMEDIATE_ACK frame.
	ImmediateAckFrame = wire.ImmediateAckFrame
)

type AckRange = wire.AckRange

// A CryptoFrame is a CRYPTO frame.
type CryptoFrame struct {
	Offset int64
	Length int64
}

// A StreamFrame is a STREAM frame.
type StreamFrame struct {
	StreamID StreamID
	Offset   int64
	Length   int64
	Fin      bool
}

// A DatagramFrame is a DATAGRAM frame.
type DatagramFrame struct {
	Length int64
}

func (fs frames) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginArray)
	for _, f := range fs {
		if err := f.Encode(enc); err != nil {
			return err
		}
	}
	h.WriteToken(jsontext.EndArray)
	return h.err
}

func (f Frame) Encode(enc *jsontext.Encoder) error {
	switch frame := f.Frame.(type) {
	case *PingFrame:
		return encodePingFrame(enc, frame)
	case *AckFrame:
		return encodeAckFrame(enc, frame)
	case *ResetStreamFrame:
		return encodeResetStreamFrame(enc, frame)
	case *StopSendingFrame:
		return encodeStopSendingFrame(enc, frame)
	case *CryptoFrame:
		return encodeCryptoFrame(enc, frame)
	case *NewTokenFrame:
		return encodeNewTokenFrame(enc, frame)
	case *StreamFrame:
		return encodeStreamFrame(enc, frame)
	case *MaxDataFrame:
		return encodeMaxDataFrame(enc, frame)
	case *MaxStreamDataFrame:
		return encodeMaxStreamDataFrame(enc, frame)
	case *MaxStreamsFrame:
		return encodeMaxStreamsFrame(enc, frame)
	case *DataBlockedFrame:
		return encodeDataBlockedFrame(enc, frame)
	case *StreamDataBlockedFrame:
		return encodeStreamDataBlockedFrame(enc, frame)
	case *StreamsBlockedFrame:
		return encodeStreamsBlockedFrame(enc, frame)
	case *NewConnectionIDFrame:
		return encodeNewConnectionIDFrame(enc, frame)
	case *RetireConnectionIDFrame:
		return encodeRetireConnectionIDFrame(enc, frame)
	case *PathChallengeFrame:
		return encodePathChallengeFrame(enc, frame)
	case *PathResponseFrame:
		return encodePathResponseFrame(enc, frame)
	case *ConnectionCloseFrame:
		return encodeConnectionCloseFrame(enc, frame)
	case *HandshakeDoneFrame:
		return encodeHandshakeDoneFrame(enc, frame)
	case *DatagramFrame:
		return encodeDatagramFrame(enc, frame)
	case *AckFrequencyFrame:
		return encodeAckFrequencyFrame(enc, frame)
	case *ImmediateAckFrame:
		return encodeImmediateAckFrame(enc, frame)
	default:
		panic("unknown frame type")
	}
}

func encodePingFrame(enc *jsontext.Encoder, _ *PingFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("ping"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type ackRanges []wire.AckRange

func (ars ackRanges) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginArray)
	for _, r := range ars {
		if err := ackRange(r).encode(enc); err != nil {
			return err
		}
	}
	h.WriteToken(jsontext.EndArray)
	return h.err
}

type ackRange wire.AckRange

func (ar ackRange) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginArray)
	h.WriteToken(jsontext.Int(int64(ar.Smallest)))
	if ar.Smallest != ar.Largest {
		h.WriteToken(jsontext.Int(int64(ar.Largest)))
	}
	h.WriteToken(jsontext.EndArray)
	return h.err
}

func encodeAckFrame(enc *jsontext.Encoder, f *AckFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("ack"))
	if f.DelayTime > 0 {
		h.WriteToken(jsontext.String("ack_delay"))
		h.WriteToken(jsontext.Float(milliseconds(f.DelayTime)))
	}
	h.WriteToken(jsontext.String("acked_ranges"))
	if err := ackRanges(f.AckRanges).encode(enc); err != nil {
		return err
	}
	hasECN := f.ECT0 > 0 || f.ECT1 > 0 || f.ECNCE > 0
	if hasECN {
		h.WriteToken(jsontext.String("ect0"))
		h.WriteToken(jsontext.Uint(f.ECT0))
		h.WriteToken(jsontext.String("ect1"))
		h.WriteToken(jsontext.Uint(f.ECT1))
		h.WriteToken(jsontext.String("ce"))
		h.WriteToken(jsontext.Uint(f.ECNCE))
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeResetStreamFrame(enc *jsontext.Encoder, f *ResetStreamFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	if f.ReliableSize > 0 {
		h.WriteToken(jsontext.String("reset_stream_at"))
	} else {
		h.WriteToken(jsontext.String("reset_stream"))
	}
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Int(int64(f.StreamID)))
	h.WriteToken(jsontext.String("error_code"))
	h.WriteToken(jsontext.Int(int64(f.ErrorCode)))
	h.WriteToken(jsontext.String("final_size"))
	h.WriteToken(jsontext.Int(int64(f.FinalSize)))
	if f.ReliableSize > 0 {
		h.WriteToken(jsontext.String("reliable_size"))
		h.WriteToken(jsontext.Int(int64(f.ReliableSize)))
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeStopSendingFrame(enc *jsontext.Encoder, f *StopSendingFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("stop_sending"))
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Int(int64(f.StreamID)))
	h.WriteToken(jsontext.String("error_code"))
	h.WriteToken(jsontext.Int(int64(f.ErrorCode)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeCryptoFrame(enc *jsontext.Encoder, f *CryptoFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("crypto"))
	h.WriteToken(jsontext.String("offset"))
	h.WriteToken(jsontext.Int(f.Offset))
	h.WriteToken(jsontext.String("length"))
	h.WriteToken(jsontext.Int(f.Length))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeNewTokenFrame(enc *jsontext.Encoder, f *NewTokenFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("new_token"))
	h.WriteToken(jsontext.String("token"))
	if err := (Token{Raw: f.Token}).encode(enc); err != nil {
		return err
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeStreamFrame(enc *jsontext.Encoder, f *StreamFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("stream"))
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Int(int64(f.StreamID)))
	h.WriteToken(jsontext.String("offset"))
	h.WriteToken(jsontext.Int(f.Offset))
	h.WriteToken(jsontext.String("length"))
	h.WriteToken(jsontext.Int(f.Length))
	if f.Fin {
		h.WriteToken(jsontext.String("fin"))
		h.WriteToken(jsontext.True)
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeMaxDataFrame(enc *jsontext.Encoder, f *MaxDataFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("max_data"))
	h.WriteToken(jsontext.String("maximum"))
	h.WriteToken(jsontext.Int(int64(f.MaximumData)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeMaxStreamDataFrame(enc *jsontext.Encoder, f *MaxStreamDataFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("max_stream_data"))
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Int(int64(f.StreamID)))
	h.WriteToken(jsontext.String("maximum"))
	h.WriteToken(jsontext.Int(int64(f.MaximumStreamData)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeMaxStreamsFrame(enc *jsontext.Encoder, f *MaxStreamsFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("max_streams"))
	h.WriteToken(jsontext.String("stream_type"))
	h.WriteToken(jsontext.String(streamType(f.Type).String()))
	h.WriteToken(jsontext.String("maximum"))
	h.WriteToken(jsontext.Int(int64(f.MaxStreamNum)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeDataBlockedFrame(enc *jsontext.Encoder, f *DataBlockedFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("data_blocked"))
	h.WriteToken(jsontext.String("limit"))
	h.WriteToken(jsontext.Int(int64(f.MaximumData)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeStreamDataBlockedFrame(enc *jsontext.Encoder, f *StreamDataBlockedFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("stream_data_blocked"))
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Int(int64(f.StreamID)))
	h.WriteToken(jsontext.String("limit"))
	h.WriteToken(jsontext.Int(int64(f.MaximumStreamData)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeStreamsBlockedFrame(enc *jsontext.Encoder, f *StreamsBlockedFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("streams_blocked"))
	h.WriteToken(jsontext.String("stream_type"))
	h.WriteToken(jsontext.String(streamType(f.Type).String()))
	h.WriteToken(jsontext.String("limit"))
	h.WriteToken(jsontext.Int(int64(f.StreamLimit)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeNewConnectionIDFrame(enc *jsontext.Encoder, f *NewConnectionIDFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("new_connection_id"))
	h.WriteToken(jsontext.String("sequence_number"))
	h.WriteToken(jsontext.Uint(f.SequenceNumber))
	h.WriteToken(jsontext.String("retire_prior_to"))
	h.WriteToken(jsontext.Uint(f.RetirePriorTo))
	h.WriteToken(jsontext.String("length"))
	h.WriteToken(jsontext.Int(int64(f.ConnectionID.Len())))
	h.WriteToken(jsontext.String("connection_id"))
	h.WriteToken(jsontext.String(f.ConnectionID.String()))
	h.WriteToken(jsontext.String("stateless_reset_token"))
	h.WriteToken(jsontext.String(hex.EncodeToString(f.StatelessResetToken[:])))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeRetireConnectionIDFrame(enc *jsontext.Encoder, f *RetireConnectionIDFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("retire_connection_id"))
	h.WriteToken(jsontext.String("sequence_number"))
	h.WriteToken(jsontext.Uint(f.SequenceNumber))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodePathChallengeFrame(enc *jsontext.Encoder, f *PathChallengeFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("path_challenge"))
	h.WriteToken(jsontext.String("data"))
	h.WriteToken(jsontext.String(hex.EncodeToString(f.Data[:])))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodePathResponseFrame(enc *jsontext.Encoder, f *PathResponseFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("path_response"))
	h.WriteToken(jsontext.String("data"))
	h.WriteToken(jsontext.String(hex.EncodeToString(f.Data[:])))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeConnectionCloseFrame(enc *jsontext.Encoder, f *ConnectionCloseFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("connection_close"))
	h.WriteToken(jsontext.String("error_space"))
	errorSpace := "transport"
	if f.IsApplicationError {
		errorSpace = "application"
	}
	h.WriteToken(jsontext.String(errorSpace))
	errName := transportError(f.ErrorCode).String()
	if len(errName) > 0 {
		h.WriteToken(jsontext.String("error_code"))
		h.WriteToken(jsontext.String(errName))
	} else {
		h.WriteToken(jsontext.String("error_code"))
		h.WriteToken(jsontext.Uint(f.ErrorCode))
	}
	h.WriteToken(jsontext.String("raw_error_code"))
	h.WriteToken(jsontext.Uint(f.ErrorCode))
	h.WriteToken(jsontext.String("reason"))
	h.WriteToken(jsontext.String(f.ReasonPhrase))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeHandshakeDoneFrame(enc *jsontext.Encoder, _ *HandshakeDoneFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("handshake_done"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeDatagramFrame(enc *jsontext.Encoder, f *DatagramFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("datagram"))
	h.WriteToken(jsontext.String("length"))
	h.WriteToken(jsontext.Int(f.Length))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeAckFrequencyFrame(enc *jsontext.Encoder, f *AckFrequencyFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("ack_frequency"))
	h.WriteToken(jsontext.String("sequence_number"))
	h.WriteToken(jsontext.Uint(f.SequenceNumber))
	h.WriteToken(jsontext.String("ack_eliciting_threshold"))
	h.WriteToken(jsontext.Uint(f.AckElicitingThreshold))
	h.WriteToken(jsontext.String("request_max_ack_delay"))
	h.WriteToken(jsontext.Float(milliseconds(f.RequestMaxAckDelay)))
	h.WriteToken(jsontext.String("reordering_threshold"))
	h.WriteToken(jsontext.Int(int64(f.ReorderingThreshold)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

func encodeImmediateAckFrame(enc *jsontext.Encoder, _ *ImmediateAckFrame) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("immediate_ack"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}
