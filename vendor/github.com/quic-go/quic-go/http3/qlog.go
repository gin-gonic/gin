package http3

import (
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3/qlog"
	"github.com/quic-go/quic-go/qlogwriter"

	"github.com/quic-go/qpack"
)

func maybeQlogInvalidHeadersFrame(qlogger qlogwriter.Recorder, streamID quic.StreamID, l uint64) {
	if qlogger != nil {
		qlogger.RecordEvent(qlog.FrameParsed{
			StreamID: streamID,
			Raw:      qlog.RawInfo{PayloadLength: int(l)},
			Frame:    qlog.Frame{Frame: qlog.HeadersFrame{}},
		})
	}
}

func qlogParsedHeadersFrame(qlogger qlogwriter.Recorder, streamID quic.StreamID, hf *headersFrame, hfs []qpack.HeaderField) {
	headerFields := make([]qlog.HeaderField, len(hfs))
	for i, hf := range hfs {
		headerFields[i] = qlog.HeaderField{
			Name:  hf.Name,
			Value: hf.Value,
		}
	}
	qlogger.RecordEvent(qlog.FrameParsed{
		StreamID: streamID,
		Raw: qlog.RawInfo{
			Length:        int(hf.Length) + hf.headerLen,
			PayloadLength: int(hf.Length),
		},
		Frame: qlog.Frame{Frame: qlog.HeadersFrame{
			HeaderFields: headerFields,
		}},
	})
}

func qlogCreatedHeadersFrame(qlogger qlogwriter.Recorder, streamID quic.StreamID, length, payloadLength int, hfs []qlog.HeaderField) {
	headerFields := make([]qlog.HeaderField, len(hfs))
	for i, hf := range hfs {
		headerFields[i] = qlog.HeaderField{
			Name:  hf.Name,
			Value: hf.Value,
		}
	}
	qlogger.RecordEvent(qlog.FrameCreated{
		StreamID: streamID,
		Raw:      qlog.RawInfo{Length: length, PayloadLength: payloadLength},
		Frame: qlog.Frame{Frame: qlog.HeadersFrame{
			HeaderFields: headerFields,
		}},
	})
}
