package decoder

import (
	"context"
	"encoding"
	"encoding/json"
	"reflect"
	"unsafe"
)

type Decoder interface {
	Decode(*RuntimeContext, int64, int64, unsafe.Pointer) (int64, error)
	DecodePath(*RuntimeContext, int64, int64) ([][]byte, int64, error)
	DecodeStream(*Stream, int64, unsafe.Pointer) error
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

type unmarshalerContext interface {
	UnmarshalJSON(context.Context, []byte) error
}

var (
	unmarshalJSONType        = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	unmarshalJSONContextType = reflect.TypeOf((*unmarshalerContext)(nil)).Elem()
	unmarshalTextType        = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)
