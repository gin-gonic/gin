package qlog

import (
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlogwriter/jsontext"
)

// Frame represents an HTTP/3 frame.
type Frame struct {
	Frame any
}

func (f Frame) encode(enc *jsontext.Encoder) error {
	switch frame := f.Frame.(type) {
	case DataFrame:
		return frame.encode(enc)
	case HeadersFrame:
		return frame.encode(enc)
	case GoAwayFrame:
		return frame.encode(enc)
	case SettingsFrame:
		return frame.encode(enc)
	case PushPromiseFrame:
		return frame.encode(enc)
	case CancelPushFrame:
		return frame.encode(enc)
	case MaxPushIDFrame:
		return frame.encode(enc)
	case ReservedFrame:
		return frame.encode(enc)
	case UnknownFrame:
		return frame.encode(enc)
	}
	// This shouldn't happen if the code is correctly logging frames.
	// Write a null token to produce valid JSON.
	return enc.WriteToken(jsontext.Null)
}

// A DataFrame is a DATA frame
type DataFrame struct{}

func (f *DataFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("data"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type HeaderField struct {
	Name  string
	Value string
}

// A HeadersFrame is a HEADERS frame
type HeadersFrame struct {
	HeaderFields []HeaderField
}

func (f *HeadersFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("headers"))
	if len(f.HeaderFields) > 0 {
		h.WriteToken(jsontext.String("header_fields"))
		h.WriteToken(jsontext.BeginArray)
		for _, f := range f.HeaderFields {
			h.WriteToken(jsontext.BeginObject)
			h.WriteToken(jsontext.String("name"))
			h.WriteToken(jsontext.String(f.Name))
			h.WriteToken(jsontext.String("value"))
			h.WriteToken(jsontext.String(f.Value))
			h.WriteToken(jsontext.EndObject)
		}
		h.WriteToken(jsontext.EndArray)
	}
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// A GoAwayFrame is a GOAWAY frame
type GoAwayFrame struct {
	StreamID quic.StreamID
}

func (f *GoAwayFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("goaway"))
	h.WriteToken(jsontext.String("id"))
	h.WriteToken(jsontext.Uint(uint64(f.StreamID)))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

type SettingsFrame struct {
	Datagram        *bool
	ExtendedConnect *bool
	Other           map[uint64]uint64
}

func (f *SettingsFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("settings"))
	h.WriteToken(jsontext.String("settings"))
	h.WriteToken(jsontext.BeginArray)
	if f.Datagram != nil {
		h.WriteToken(jsontext.BeginObject)
		h.WriteToken(jsontext.String("name"))
		h.WriteToken(jsontext.String("settings_h3_datagram"))
		h.WriteToken(jsontext.String("value"))
		h.WriteToken(jsontext.Bool(*f.Datagram))
		h.WriteToken(jsontext.EndObject)
	}
	if f.ExtendedConnect != nil {
		h.WriteToken(jsontext.BeginObject)
		h.WriteToken(jsontext.String("name"))
		h.WriteToken(jsontext.String("settings_enable_connect_protocol"))
		h.WriteToken(jsontext.String("value"))
		h.WriteToken(jsontext.Bool(*f.ExtendedConnect))
		h.WriteToken(jsontext.EndObject)
	}
	if len(f.Other) > 0 {
		for k, v := range f.Other {
			h.WriteToken(jsontext.BeginObject)
			h.WriteToken(jsontext.String("name"))
			h.WriteToken(jsontext.String("unknown"))
			h.WriteToken(jsontext.String("name_bytes"))
			h.WriteToken(jsontext.Uint(k))
			h.WriteToken(jsontext.String("value"))
			h.WriteToken(jsontext.Uint(v))
			h.WriteToken(jsontext.EndObject)
		}
	}
	h.WriteToken(jsontext.EndArray)
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// A PushPromiseFrame is a PUSH_PROMISE frame
type PushPromiseFrame struct{}

func (f *PushPromiseFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("push_promise"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// A CancelPushFrame is a CANCEL_PUSH frame
type CancelPushFrame struct{}

func (f *CancelPushFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("cancel_push"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// A MaxPushIDFrame is a MAX_PUSH_ID frame
type MaxPushIDFrame struct{}

func (f *MaxPushIDFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("max_push_id"))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// A ReservedFrame is one of the reserved frame types
type ReservedFrame struct {
	Type uint64
}

func (f *ReservedFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("reserved"))
	h.WriteToken(jsontext.String("frame_type_bytes"))
	h.WriteToken(jsontext.Uint(f.Type))
	h.WriteToken(jsontext.EndObject)
	return h.err
}

// An UnknownFrame is an unknown frame type
type UnknownFrame struct {
	Type uint64
}

func (f *UnknownFrame) encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("frame_type"))
	h.WriteToken(jsontext.String("unknown"))
	h.WriteToken(jsontext.String("frame_type_bytes"))
	h.WriteToken(jsontext.Uint(f.Type))
	h.WriteToken(jsontext.EndObject)
	return h.err
}
