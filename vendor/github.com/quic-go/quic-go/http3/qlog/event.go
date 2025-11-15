package qlog

import (
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlogwriter/jsontext"
)

type encoderHelper struct {
	enc *jsontext.Encoder
	err error
}

func (h *encoderHelper) WriteToken(t jsontext.Token) {
	if h.err != nil {
		return
	}
	h.err = h.enc.WriteToken(t)
}

type RawInfo struct {
	Length        int // full packet length, including header and AEAD authentication tag
	PayloadLength int // length of the packet payload, excluding AEAD tag
}

func (i RawInfo) HasValues() bool {
	return i.Length != 0 || i.PayloadLength != 0
}

func (i RawInfo) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	if i.Length != 0 {
		h.WriteToken(jsontext.String("length"))
		h.WriteToken(jsontext.Uint(uint64(i.Length)))
	}
	if i.PayloadLength != 0 {
		h.WriteToken(jsontext.String("payload_length"))
		h.WriteToken(jsontext.Uint(uint64(i.PayloadLength)))
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type FrameParsed struct {
	StreamID quic.StreamID
	Raw      RawInfo
	Frame    Frame
}

func (e FrameParsed) Name() string { return "http3:frame_parsed" }

func (e FrameParsed) Encode(enc *jsontext.Encoder, _ time.Time) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Uint(uint64(e.StreamID)))
	if e.Raw.HasValues() {
		h.WriteToken(jsontext.String("raw"))
		if err := e.Raw.encode(enc); err != nil {
			return err
		}
	}
	h.WriteToken(jsontext.String("frame"))
	if err := e.Frame.encode(enc); err != nil {
		return err
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type FrameCreated struct {
	StreamID quic.StreamID
	Raw      RawInfo
	Frame    Frame
}

func (e FrameCreated) Name() string { return "http3:frame_created" }

func (e FrameCreated) Encode(enc *jsontext.Encoder, _ time.Time) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("stream_id"))
	h.WriteToken(jsontext.Uint(uint64(e.StreamID)))
	if e.Raw.HasValues() {
		h.WriteToken(jsontext.String("raw"))
		if err := e.Raw.encode(enc); err != nil {
			return err
		}
	}
	h.WriteToken(jsontext.String("frame"))
	if err := e.Frame.encode(enc); err != nil {
		return err
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type DatagramCreated struct {
	QuaterStreamID uint64
	Raw            RawInfo
}

func (e DatagramCreated) Name() string { return "http3:datagram_created" }

func (e DatagramCreated) Encode(enc *jsontext.Encoder, _ time.Time) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("quater_stream_id"))
	h.WriteToken(jsontext.Uint(e.QuaterStreamID))
	h.WriteToken(jsontext.String("raw"))
	if err := e.Raw.encode(enc); err != nil {
		return err
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type DatagramParsed struct {
	QuaterStreamID uint64
	Raw            RawInfo
}

func (e DatagramParsed) Name() string { return "http3:datagram_parsed" }

func (e DatagramParsed) Encode(enc *jsontext.Encoder, _ time.Time) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("quater_stream_id"))
	h.WriteToken(jsontext.Uint(e.QuaterStreamID))
	h.WriteToken(jsontext.String("raw"))
	if err := e.Raw.encode(enc); err != nil {
		return err
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}
