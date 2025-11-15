package qlogwriter

import (
	"runtime/debug"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/qlogwriter/jsontext"
)

type ConnectionID = protocol.ConnectionID

// Setting of this only works when quic-go is used as a library.
// When building a binary from this repository, the version can be set using the following go build flag:
// -ldflags="-X github.com/quic-go/quic-go/qlogwriter.quicGoVersion=foobar"
var quicGoVersion = "(devel)"

func init() {
	if quicGoVersion != "(devel)" { // variable set by ldflags
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok { // no build info available. This happens when quic-go is not used as a library.
		return
	}
	for _, d := range info.Deps {
		if d.Path == "github.com/quic-go/quic-go" {
			quicGoVersion = d.Version
			if d.Replace != nil {
				if len(d.Replace.Version) > 0 {
					quicGoVersion = d.Version
				} else {
					quicGoVersion += " (replaced)"
				}
			}
			break
		}
	}
}

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

type traceHeader struct {
	VantagePointType string
	GroupID          *ConnectionID
	ReferenceTime    time.Time
	EventSchemas     []string
}

func (l traceHeader) Encode(enc *jsontext.Encoder) error {
	h := encoderHelper{enc: enc}
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("file_schema"))
	h.WriteToken(jsontext.String("urn:ietf:params:qlog:file:sequential"))
	h.WriteToken(jsontext.String("serialization_format"))
	h.WriteToken(jsontext.String("application/qlog+json-seq"))
	h.WriteToken(jsontext.String("title"))
	h.WriteToken(jsontext.String("quic-go qlog"))
	h.WriteToken(jsontext.String("code_version"))
	h.WriteToken(jsontext.String(quicGoVersion))

	h.WriteToken(jsontext.String("trace"))
	// trace
	h.WriteToken(jsontext.BeginObject)
	if len(l.EventSchemas) > 0 {
		h.WriteToken(jsontext.String("event_schemas"))
		h.WriteToken(jsontext.BeginArray)
		for _, schema := range l.EventSchemas {
			h.WriteToken(jsontext.String(schema))
		}
		h.WriteToken(jsontext.EndArray)
	}

	h.WriteToken(jsontext.String("vantage_point"))
	// -- vantage_point
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("type"))
	h.WriteToken(jsontext.String(l.VantagePointType))
	// -- end vantage_point
	h.WriteToken(jsontext.EndObject)

	h.WriteToken(jsontext.String("common_fields"))
	// -- common_fields
	h.WriteToken(jsontext.BeginObject)
	if l.GroupID != nil {
		h.WriteToken(jsontext.String("group_id"))
		h.WriteToken(jsontext.String(l.GroupID.String()))
	}
	h.WriteToken(jsontext.String("reference_time"))
	// ---- reference_time
	h.WriteToken(jsontext.BeginObject)
	h.WriteToken(jsontext.String("clock_type"))
	h.WriteToken(jsontext.String("monotonic"))
	h.WriteToken(jsontext.String("epoch"))
	h.WriteToken(jsontext.String("unknown"))
	h.WriteToken(jsontext.String("wall_clock_time"))
	h.WriteToken(jsontext.String(l.ReferenceTime.Format(time.RFC3339Nano)))
	// ---- end reference_time
	h.WriteToken(jsontext.EndObject)
	// -- end common_fields
	h.WriteToken(jsontext.EndObject)
	// end trace
	h.WriteToken(jsontext.EndObject)

	// The following fields are not required by the qlog draft anymore,
	// but qvis still requires them to be present.
	h.WriteToken(jsontext.String("qlog_format"))
	h.WriteToken(jsontext.String("JSON-SEQ"))
	h.WriteToken(jsontext.String("qlog_version"))
	h.WriteToken(jsontext.String("0.3"))

	h.WriteToken(jsontext.EndObject)
	return h.err
}
