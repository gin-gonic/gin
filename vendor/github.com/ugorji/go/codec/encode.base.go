// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"cmp"
	"errors"
	"io"
	"reflect"
	"slices"
	"sync"
	"time"
)

var errEncoderNotInitialized = errors.New("encoder not initialized")

var encBuiltinRtids []uintptr

func init() {
	for _, v := range []interface{}{
		(string)(""),
		(bool)(false),
		(int)(0),
		(int8)(0),
		(int16)(0),
		(int32)(0),
		(int64)(0),
		(uint)(0),
		(uint8)(0),
		(uint16)(0),
		(uint32)(0),
		(uint64)(0),
		(uintptr)(0),
		(float32)(0),
		(float64)(0),
		(complex64)(0),
		(complex128)(0),
		(time.Time{}),
		([]byte)(nil),
		(Raw{}),
		// (interface{})(nil),
	} {
		t := reflect.TypeOf(v)
		encBuiltinRtids = append(encBuiltinRtids, rt2id(t), rt2id(reflect.PointerTo(t)))
	}
	slices.Sort(encBuiltinRtids)
}

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriverI interface {
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	// re is never nil
	EncodeRawExt(re *RawExt)
	// ext is never nil
	EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext)
	// EncodeString using cUTF8, honor'ing StringToRaw flag
	EncodeString(v string)
	EncodeStringNoEscape4Json(v string)
	// encode a non-nil []byte
	EncodeStringBytesRaw(v []byte)
	// encode a []byte as nil, empty or encoded sequence of bytes depending on context
	EncodeBytes(v []byte)
	EncodeTime(time.Time)
	WriteArrayStart(length int)
	WriteArrayEnd()
	WriteMapStart(length int)
	WriteMapEnd()

	// these write a zero-len map or array into the stream
	WriteMapEmpty()
	WriteArrayEmpty()

	writeNilMap()
	writeNilArray()
	writeNilBytes()

	// these are no-op except for json
	encDriverContainerTracker

	// reset will reset current encoding runtime state, and cached information from the handle
	reset()

	atEndOfEncode()
	writerEnd()

	writeBytesAsis(b []byte)
	// writeStringAsisDblQuoted(v string)

	resetOutBytes(out *[]byte)
	resetOutIO(out io.Writer)

	init(h Handle, shared *encoderBase, enc encoderI) (fp interface{})

	// driverStateManager
}

type encInit2er struct{}

func (encInit2er) init2(enc encoderI) {}

type encDriverContainerTracker interface {
	WriteArrayElem(firstTime bool)
	WriteMapElemKey(firstTime bool)
	WriteMapElemValue()
}

type encDriverNoState struct{}

// func (encDriverNoState) captureState() interface{}  { return nil }
// func (encDriverNoState) resetState()                {}
// func (encDriverNoState) restoreState(v interface{}) {}
func (encDriverNoState) reset() {}

type encDriverNoopContainerWriter struct{}

func (encDriverNoopContainerWriter) WriteArrayStart(length int) {}
func (encDriverNoopContainerWriter) WriteArrayEnd()             {}
func (encDriverNoopContainerWriter) WriteMapStart(length int)   {}
func (encDriverNoopContainerWriter) WriteMapEnd()               {}
func (encDriverNoopContainerWriter) atEndOfEncode()             {}

// encStructFieldObj[Slice] is used for sorting when there are missing fields and canonical flag is set
type encStructFieldObj struct {
	key        string
	rv         reflect.Value
	intf       interface{}
	isRv       bool
	noEsc4json bool
	builtin    bool
}

type encStructFieldObjSlice []encStructFieldObj

func (p encStructFieldObjSlice) Len() int      { return len(p) }
func (p encStructFieldObjSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p encStructFieldObjSlice) Less(i, j int) bool {
	return p[uint(i)].key < p[uint(j)].key
}

// ----

type orderedRv[T cmp.Ordered] struct {
	v T
	r reflect.Value
}

func cmpOrderedRv[T cmp.Ordered](v1, v2 orderedRv[T]) int {
	return cmp.Compare(v1.v, v2.v)
}

// ----

type encFnInfo struct {
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	addrE bool
	// addrEf bool // force: if addrE, then encode function MUST take a ptr
}

// ----

// EncodeOptions captures configuration options during encode.
type EncodeOptions struct {
	// WriterBufferSize is the size of the buffer used when writing.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	WriterBufferSize int

	// ChanRecvTimeout is the timeout used when selecting from a chan.
	//
	// Configuring this controls how we receive from a chan during the encoding process.
	//   - If ==0, we only consume the elements currently available in the chan.
	//   - if  <0, we consume until the chan is closed.
	//   - If  >0, we consume until this timeout.
	ChanRecvTimeout time.Duration

	// StructToArray specifies to encode a struct as an array, and not as a map
	StructToArray bool

	// Canonical representation means that encoding a value will always result in the same
	// sequence of bytes.
	//
	// This only affects maps, as the iteration order for maps is random.
	//
	// The implementation MAY use the natural sort order for the map keys if possible:
	//
	//     - If there is a natural sort order (ie for number, bool, string or []byte keys),
	//       then the map keys are first sorted in natural order and then written
	//       with corresponding map values to the strema.
	//     - If there is no natural sort order, then the map keys will first be
	//       encoded into []byte, and then sorted,
	//       before writing the sorted keys and the corresponding map values to the stream.
	//
	Canonical bool

	// CheckCircularRef controls whether we check for circular references
	// and error fast during an encode.
	//
	// If enabled, an error is received if a pointer to a struct
	// references itself either directly or through one of its fields (iteratively).
	//
	// This is opt-in, as there may be a performance hit to checking circular references.
	CheckCircularRef bool

	// RecursiveEmptyCheck controls how we determine whether a value is empty.
	//
	// If true, we descend into interfaces and pointers to reursively check if value is empty.
	//
	// We *might* check struct fields one by one to see if empty
	// (if we cannot directly check if a struct value is equal to its zero value).
	// If so, we honor IsZero, Comparable, IsCodecEmpty(), etc.
	// Note: This *may* make OmitEmpty more expensive due to the large number of reflect calls.
	//
	// If false, we check if the value is equal to its zero value (newly allocated state).
	RecursiveEmptyCheck bool

	// Raw controls whether we encode Raw values.
	// This is a "dangerous" option and must be explicitly set.
	// If set, we blindly encode Raw values as-is, without checking
	// if they are a correct representation of a value in that format.
	// If unset, we error out.
	Raw bool

	// StringToRaw controls how strings are encoded.
	//
	// As a go string is just an (immutable) sequence of bytes,
	// it can be encoded either as raw bytes or as a UTF string.
	//
	// By default, strings are encoded as UTF-8.
	// but can be treated as []byte during an encode.
	//
	// Note that things which we know (by definition) to be UTF-8
	// are ALWAYS encoded as UTF-8 strings.
	// These include encoding.TextMarshaler, time.Format calls, struct field names, etc.
	StringToRaw bool

	// OptimumSize controls whether we optimize for the smallest size.
	//
	// Some formats will use this flag to determine whether to encode
	// in the smallest size possible, even if it takes slightly longer.
	//
	// For example, some formats that support half-floats might check if it is possible
	// to store a float64 as a half float. Doing this check has a small performance cost,
	// but the benefit is that the encoded message will be smaller.
	OptimumSize bool

	// NoAddressableReadonly controls whether we try to force a non-addressable value
	// to be addressable so we can call a pointer method on it e.g. for types
	// that support Selfer, json.Marshaler, etc.
	//
	// Use it in the very rare occurrence that your types modify a pointer value when calling
	// an encode callback function e.g. JsonMarshal, TextMarshal, BinaryMarshal or CodecEncodeSelf.
	NoAddressableReadonly bool

	// NilCollectionToZeroLength controls whether we encode nil collections (map, slice, chan)
	// as nil (e.g. null if using JSON) or as zero length collections (e.g. [] or {} if using JSON).
	//
	// This is useful in many scenarios e.g.
	//    - encoding in go, but decoding the encoded stream in python
	//      where context of the type is missing but needed
	//
	// Note: this flag ignores the MapBySlice tag, and will encode nil slices, maps and chan
	// in their natural zero-length formats e.g. a slice in json encoded as []
	// (and not nil or {} if MapBySlice tag).
	NilCollectionToZeroLength bool
}

// ---------------------------------------------

// encoderBase is shared as a field between Encoder and its encDrivers.
// This way, encDrivers need not hold a referece to the Encoder itself.
type encoderBase struct {
	perType encPerType

	h *BasicHandle

	// MARKER: these fields below should belong directly in Encoder.
	// There should not be any pointers here - just values.
	// we pack them here for space efficiency and cache-line optimization.

	rtidFn, rtidFnNoExt *atomicRtidFnSlice

	// se  encoderI
	err error

	blist bytesFreeList

	// js bool // is json encoder?
	// be bool // is binary encoder?

	bytes bool

	c containerState

	calls uint16
	seq   uint16 // sequencer (e.g. used by binc for symbols, etc)

	// ---- cpu cache line boundary
	hh Handle

	// ---- cpu cache line boundary

	// ---- writable fields during execution --- *try* to keep in sep cache line

	ci circularRefChecker

	slist sfiRvFreeList
}

func (e *encoderBase) HandleName() string {
	return e.hh.Name()
}

// Release is a no-op.
//
// Deprecated: Pooled resources are not used with an Encoder.
// This method is kept for compatibility reasons only.
func (e *encoderBase) Release() {
}

func (e *encoderBase) setContainerState(cs containerState) {
	if cs != 0 {
		e.c = cs
	}
}

func (e *encoderBase) haltOnMbsOddLen(length int) {
	if length&1 != 0 { // similar to &1==1 or %2 == 1
		halt.errorInt("mapBySlice requires even slice length, but got ", int64(length))
	}
}

// addrRV returns a addressable value given that rv is not addressable
func (e *encoderBase) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
	// if rv.CanAddr() {
	// 	return rvAddr(rv, ptrType)
	// }
	if e.h.NoAddressableReadonly {
		rva = reflect.New(typ)
		rvSetDirect(rva.Elem(), rv)
		return
	}
	return rvAddr(e.perType.AddressableRO(rv), ptrType)
}

func (e *encoderBase) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}

func (e *encoderBase) kErr(_ *encFnInfo, rv reflect.Value) {
	halt.errorf("unsupported encoding kind: %s, for %#v", rv.Kind(), any(rv))
}

func chanToSlice(rv reflect.Value, rtslice reflect.Type, timeout time.Duration) (rvcs reflect.Value) {
	rvcs = rvZeroK(rtslice, reflect.Slice)
	if timeout < 0 { // consume until close
		for {
			recv, recvOk := rv.Recv()
			if !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	} else {
		cases := make([]reflect.SelectCase, 2)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: rv}
		if timeout == 0 {
			cases[1] = reflect.SelectCase{Dir: reflect.SelectDefault}
		} else {
			tt := time.NewTimer(timeout)
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(tt.C)}
		}
		for {
			chosen, recv, recvOk := reflect.Select(cases)
			if chosen == 1 || !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	}
	return
}

type encoderI interface {
	Encode(v interface{}) error
	MustEncode(v interface{})
	Release()
	Reset(w io.Writer)
	ResetBytes(out *[]byte)

	wrapErr(v error, err *error)
	atEndOfEncode()
	writerEnd()

	encodeI(v interface{})
	encodeR(v reflect.Value)
	encodeAs(v interface{}, t reflect.Type, ext bool)

	setContainerState(cs containerState) // needed for canonical encoding via side encoder
}

var errEncNoResetBytesWithWriter = errors.New("cannot reset an Encoder which outputs to []byte with a io.Writer")
var errEncNoResetWriterWithBytes = errors.New("cannot reset an Encoder which outputs to io.Writer with a []byte")

type encDriverContainerNoTrackerT struct{}

func (encDriverContainerNoTrackerT) WriteArrayElem(firstTime bool)  {}
func (encDriverContainerNoTrackerT) WriteMapElemKey(firstTime bool) {}
func (encDriverContainerNoTrackerT) WriteMapElemValue()             {}

type Encoder struct {
	encoderI
}

// NewEncoder returns an Encoder for encoding into an io.Writer.
//
// For efficiency, Users are encouraged to configure WriterBufferSize on the handle
// OR pass in a memory buffered writer (eg bufio.Writer, bytes.Buffer).
func NewEncoder(w io.Writer, h Handle) *Encoder {
	return &Encoder{h.newEncoder(w)}
}

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder {
	return &Encoder{h.newEncoderBytes(out)}
}

// ----

func sideEncode(h Handle, p *sync.Pool, fn func(encoderI)) {
	var s encoderI
	if usePoolForSideEncode {
		s = p.Get().(encoderI)
		defer p.Put(s)
	} else {
		// initialization cycle error
		// s = NewEncoderBytes(nil, h).encoderI
		s = p.New().(encoderI)
	}
	fn(s)
}

func oneOffEncode(se encoderI, v interface{}, out *[]byte, basetype reflect.Type, ext bool) {
	se.ResetBytes(out)
	se.encodeAs(v, basetype, ext)
	se.atEndOfEncode()
	se.writerEnd()
	// e.sideEncoder(&bs)
	// e.sideEncode(v, basetype, 0)
}
