// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"io"
	"math"
	"reflect"
	"slices"
	"sync"
	"time"
)

func init() {
	for _, v := range []interface{}{
		(*string)(nil),
		(*bool)(nil),
		(*int)(nil),
		(*int8)(nil),
		(*int16)(nil),
		(*int32)(nil),
		(*int64)(nil),
		(*uint)(nil),
		(*uint8)(nil),
		(*uint16)(nil),
		(*uint32)(nil),
		(*uint64)(nil),
		(*uintptr)(nil),
		(*float32)(nil),
		(*float64)(nil),
		(*complex64)(nil),
		(*complex128)(nil),
		(*[]byte)(nil),
		([]byte)(nil),
		(*time.Time)(nil),
		(*Raw)(nil),
		(*interface{})(nil),
	} {
		decBuiltinRtids = append(decBuiltinRtids, i2rtid(v))
	}
	slices.Sort(decBuiltinRtids)
}

const msgBadDesc = "unrecognized descriptor byte"

var decBuiltinRtids []uintptr

// decDriver calls (DecodeBytes and DecodeStringAsBytes) return a state
// of the view they return, allowing consumers to handle appropriately.
//
// sequencing of this is intentional:
//   - mutable if <= dBytesAttachBuffer  (buf | view | invalid)
//   - noCopy if >= dBytesAttachViewZerocopy
type dBytesAttachState uint8

const (
	dBytesAttachInvalid      dBytesAttachState = iota
	dBytesAttachView                           // (bytes && !zerocopy && !buf)
	dBytesAttachBuffer                         // (buf)
	dBytesAttachViewZerocopy                   // (bytes && zerocopy && !buf)
	dBytesDetach                               // (!bytes && !buf)
)

type dBytesIntoState uint8

const (
	dBytesIntoNoChange dBytesIntoState = iota
	dBytesIntoParamOut
	dBytesIntoParamOutSlice
	dBytesIntoNew
)

func (x dBytesAttachState) String() string {
	switch x {
	case dBytesAttachInvalid:
		return "invalid"
	case dBytesAttachView:
		return "view"
	case dBytesAttachBuffer:
		return "buffer"
	case dBytesAttachViewZerocopy:
		return "view-zerocopy"
	case dBytesDetach:
		return "detach"
	}
	return "unknown"
}

const (
	decDefMaxDepth         = 1024        // maximum depth
	decDefChanCap          = 64          // should be large, as cap cannot be expanded
	decScratchByteArrayLen = (4 + 3) * 8 // around cacheLineSize ie ~64, depending on Decoder size

	// MARKER: massage decScratchByteArrayLen to ensure xxxDecDriver structs fit within cacheLine*N

	// decFailNonEmptyIntf configures whether we error
	// when decoding naked into a non-empty interface.
	//
	// Typically, we cannot decode non-nil stream value into
	// nil interface with methods (e.g. io.Reader).
	// However, in some scenarios, this should be allowed:
	//   - MapType
	//   - SliceType
	//   - Extensions
	//
	// Consequently, we should relax this. Put it behind a const flag for now.
	decFailNonEmptyIntf = false

	// decUseTransient says whether we should use the transient optimization.
	//
	// There's potential for GC corruption or memory overwrites if transient isn't
	// used carefully, so this flag helps turn it off quickly if needed.
	//
	// Use it everywhere needed so we can completely remove unused code blocks.
	decUseTransient = true
)

var (
	errNeedMapOrArrayDecodeToStruct = errors.New("only encoded map or array can decode into struct")
	errCannotDecodeIntoNil          = errors.New("cannot decode into nil")

	errExpandSliceCannotChange = errors.New("expand slice: cannot change")

	errDecoderNotInitialized = errors.New("Decoder not initialized")

	errDecUnreadByteNothingToRead   = errors.New("cannot unread - nothing has been read")
	errDecUnreadByteLastByteNotRead = errors.New("cannot unread - last byte has not been read")
	errDecUnreadByteUnknown         = errors.New("cannot unread - reason unknown")
	errMaxDepthExceeded             = errors.New("maximum decoding depth exceeded")
)

type decNotDecodeableReason uint8

const (
	decNotDecodeableReasonUnknown decNotDecodeableReason = iota
	decNotDecodeableReasonBadKind
	decNotDecodeableReasonNonAddrValue
	decNotDecodeableReasonNilReference
)

type decDriverI interface {

	// this will check if the next token is a break.
	CheckBreak() bool

	// TryNil tries to decode as nil.
	// If a nil is in the stream, it consumes it and returns true.
	//
	// Note: if TryNil returns true, that must be handled.
	TryNil() bool

	// ContainerType returns one of: Bytes, String, Nil, Slice or Map.
	//
	// Return unSet if not known.
	//
	// Note: Implementations MUST fully consume sentinel container types, specifically Nil.
	ContainerType() (vt valueType)

	// DecodeNaked will decode primitives (number, bool, string, []byte) and RawExt.
	// For maps and arrays, it will not do the decoding in-band, but will signal
	// the decoder, so that is done later, by setting the fauxUnion.valueType field.
	//
	// Note: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	// for extensions, DecodeNaked must read the tag and the []byte if it exists.
	// if the []byte is not read, then kInterfaceNaked will treat it as a Handle
	// that stores the subsequent value in-band, and complete reading the RawExt.
	//
	// extensions should also use readx to decode them, for efficiency.
	// kInterface will extract the detached byte slice if it has to pass it outside its realm.
	DecodeNaked()

	DecodeInt64() (i int64)
	DecodeUint64() (ui uint64)

	DecodeFloat32() (f float32)
	DecodeFloat64() (f float64)

	DecodeBool() (b bool)

	// DecodeStringAsBytes returns the bytes representing a string.
	// It will return a view into scratch buffer or input []byte (if applicable).
	//
	// Note: This can also decode symbols, if supported.
	//
	// Users should consume it right away and not store it for later use.
	DecodeStringAsBytes() (v []byte, state dBytesAttachState)

	// DecodeBytes returns the bytes representing a binary value.
	// It will return a view into scratch buffer or input []byte (if applicable).
	DecodeBytes() (out []byte, state dBytesAttachState)
	// DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte)

	// DecodeExt will decode into an extension.
	// ext is never nil.
	DecodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)

	// DecodeRawExt will decode into a *RawExt
	DecodeRawExt(re *RawExt)

	DecodeTime() (t time.Time)

	// ReadArrayStart will return the length of the array.
	// If the format doesn't prefix the length, it returns containerLenUnknown.
	// If the expected array was a nil in the stream, it returns containerLenNil.
	ReadArrayStart() int

	// ReadMapStart will return the length of the array.
	// If the format doesn't prefix the length, it returns containerLenUnknown.
	// If the expected array was a nil in the stream, it returns containerLenNil.
	ReadMapStart() int

	decDriverContainerTracker

	reset()

	// atEndOfDecode()

	// nextValueBytes will return the bytes representing the next value in the stream.
	// It generally will include the last byte read, as that is a part of the next value
	// in the stream.
	nextValueBytes() []byte

	// descBd will describe the token descriptor that signifies what type was decoded
	descBd() string

	// isBytes() bool

	resetInBytes(in []byte)
	resetInIO(r io.Reader)

	NumBytesRead() int

	init(h Handle, shared *decoderBase, dec decoderI) (fp interface{})

	// driverStateManager
	decNegintPosintFloatNumber
}

type decInit2er struct{}

func (decInit2er) init2(dec decoderI) {}

type decDriverContainerTracker interface {
	ReadArrayElem(firstTime bool)
	ReadMapElemKey(firstTime bool)
	ReadMapElemValue()
	ReadArrayEnd()
	ReadMapEnd()
}

type decNegintPosintFloatNumber interface {
	decInteger() (ui uint64, neg, ok bool)
	decFloat() (f float64, ok bool)
}

type decDriverNoopNumberHelper struct{}

func (x decDriverNoopNumberHelper) decInteger() (ui uint64, neg, ok bool) {
	panic("decInteger unsupported")
}
func (x decDriverNoopNumberHelper) decFloat() (f float64, ok bool) { panic("decFloat unsupported") }

type decDriverNoopContainerReader struct{}

func (x decDriverNoopContainerReader) ReadArrayStart() (v int)       { panic("ReadArrayStart unsupported") }
func (x decDriverNoopContainerReader) ReadMapStart() (v int)         { panic("ReadMapStart unsupported") }
func (x decDriverNoopContainerReader) ReadArrayEnd()                 {}
func (x decDriverNoopContainerReader) ReadMapEnd()                   {}
func (x decDriverNoopContainerReader) ReadArrayElem(firstTime bool)  {}
func (x decDriverNoopContainerReader) ReadMapElemKey(firstTime bool) {}
func (x decDriverNoopContainerReader) ReadMapElemValue()             {}
func (x decDriverNoopContainerReader) CheckBreak() (v bool)          { return }

// ----

type decFnInfo struct {
	ti     *typeInfo
	xfFn   Ext
	xfTag  uint64
	addrD  bool // decoding into a pointer is preferred
	addrDf bool // force: if addrD, then decode function MUST take a ptr
}

// DecodeOptions captures configuration options during decode.
type DecodeOptions struct {
	// MapType specifies type to use during schema-less decoding of a map in the stream.
	// If nil (unset), we default to map[string]interface{} iff json handle and MapKeyAsString=true,
	// else map[interface{}]interface{}.
	MapType reflect.Type

	// SliceType specifies type to use during schema-less decoding of an array in the stream.
	// If nil (unset), we default to []interface{} for all formats.
	SliceType reflect.Type

	// MaxInitLen defines the maxinum initial length that we "make" a collection
	// (string, slice, map, chan). If 0 or negative, we default to a sensible value
	// based on the size of an element in the collection.
	//
	// For example, when decoding, a stream may say that it has 2^64 elements.
	// We should not auto-matically provision a slice of that size, to prevent Out-Of-Memory crash.
	// Instead, we provision up to MaxInitLen, fill that up, and start appending after that.
	MaxInitLen int

	// ReaderBufferSize is the size of the buffer used when reading.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	ReaderBufferSize int

	// MaxDepth defines the maximum depth when decoding nested
	// maps and slices. If 0 or negative, we default to a suitably large number (currently 1024).
	MaxDepth int16

	// If ErrorIfNoField, return an error when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool

	// If ErrorIfNoArrayExpand, return an error when decoding a slice/array that cannot be expanded.
	// For example, the stream contains an array of 8 items, but you are decoding into a [4]T array,
	// or you are decoding into a slice of length 4 which is non-addressable (and so cannot be set).
	ErrorIfNoArrayExpand bool

	// If SignedInteger, use the int64 during schema-less decoding of unsigned values (not uint64).
	SignedInteger bool

	// MapValueReset controls how we decode into a map value.
	//
	// By default, we MAY retrieve the mapping for a key, and then decode into that.
	// However, especially with big maps, that retrieval may be expensive and unnecessary
	// if the stream already contains all that is necessary to recreate the value.
	//
	// If true, we will never retrieve the previous mapping,
	// but rather decode into a new value and set that in the map.
	//
	// If false, we will retrieve the previous mapping if necessary e.g.
	// the previous mapping is a pointer, or is a struct or array with pre-set state,
	// or is an interface.
	MapValueReset bool

	// SliceElementReset: on decoding a slice, reset the element to a zero value first.
	//
	// concern: if the slice already contained some garbage, we will decode into that garbage.
	SliceElementReset bool

	// InterfaceReset controls how we decode into an interface.
	//
	// By default, when we see a field that is an interface{...},
	// or a map with interface{...} value, we will attempt decoding into the
	// "contained" value.
	//
	// However, this prevents us from reading a string into an interface{}
	// that formerly contained a number.
	//
	// If true, we will decode into a new "blank" value, and set that in the interface.
	// If false, we will decode into whatever is contained in the interface.
	InterfaceReset bool

	// InternString controls interning of strings during decoding.
	//
	// Some handles, e.g. json, typically will read map keys as strings.
	// If the set of keys are finite, it may help reduce allocation to
	// look them up from a map (than to allocate them afresh).
	//
	// Note: Handles will be smart when using the intern functionality.
	// Every string should not be interned.
	// An excellent use-case for interning is struct field names,
	// or map keys where key type is string.
	InternString bool

	// PreferArrayOverSlice controls whether to decode to an array or a slice.
	//
	// This only impacts decoding into a nil interface{}.
	//
	// Consequently, it has no effect on codecgen.
	//
	// *Note*: This only applies if using go1.5 and above,
	// as it requires reflect.ArrayOf support which was absent before go1.5.
	PreferArrayOverSlice bool

	// DeleteOnNilMapValue controls how to decode a nil value in the stream.
	//
	// If true, we will delete the mapping of the key.
	// Else, just set the mapping to the zero value of the type.
	//
	// Deprecated: This does NOTHING and is left behind for compiling compatibility.
	// This change is necessitated because 'nil' in a stream now consistently
	// means the zero value (ie reset the value to its zero state).
	DeleteOnNilMapValue bool

	// RawToString controls how raw bytes in a stream are decoded into a nil interface{}.
	// By default, they are decoded as []byte, but can be decoded as string (if configured).
	RawToString bool

	// ZeroCopy controls whether decoded values of []byte or string type
	// point into the input []byte parameter passed to a NewDecoderBytes/ResetBytes(...) call.
	//
	// To illustrate, if ZeroCopy and decoding from a []byte (not io.Writer),
	// then a []byte or string in the output result may just be a slice of (point into)
	// the input bytes.
	//
	// This optimization prevents unnecessary copying.
	//
	// However, it is made optional, as the caller MUST ensure that the input parameter []byte is
	// not modified after the Decode() happens, as any changes are mirrored in the decoded result.
	ZeroCopy bool

	// PreferPointerForStructOrArray controls whether a struct or array
	// is stored in a nil interface{}, or a pointer to it.
	//
	// This mostly impacts when we decode registered extensions.
	PreferPointerForStructOrArray bool

	// ValidateUnicode controls will cause decoding to fail if an expected unicode
	// string is well-formed but include invalid codepoints.
	//
	// This could have a performance impact.
	ValidateUnicode bool
}

// ----------------------------------------

type decoderBase struct {
	perType decPerType

	h *BasicHandle

	rtidFn, rtidFnNoExt *atomicRtidFnSlice

	buf []byte

	// used for interning strings
	is internerMap

	err error

	// sd decoderI

	blist bytesFreeList

	mtr  bool // is maptype a known type?
	str  bool // is slicetype a known type?
	jsms bool // is json handle, and MapKeyAsString

	bytes bool // uses a bytes reader
	bufio bool // uses a ioDecReader with buffer size > 0

	// ---- cpu cache line boundary?
	// ---- writable fields during execution --- *try* to keep in sep cache line
	maxdepth int16
	depth    int16

	// Extensions can call Decode() within a current Decode() call.
	// We need to know when the top level Decode() call returns,
	// so we can decide whether to Release() or not.
	calls uint16 // what depth in mustDecode are we in now.

	c containerState

	// decByteState

	n fauxUnion

	// b is an always-available scratch buffer used by Decoder and decDrivers.
	// By being always-available, it can be used for one-off things without
	// having to get from freelist, use, and return back to freelist.
	//
	// Use it for a narrow set of things e.g.
	//   - binc uses it for parsing numbers, represented at 8 or less bytes
	//   - uses as potential buffer for struct field names
	b [decScratchByteArrayLen]byte

	hh Handle
	// cache the mapTypeId and sliceTypeId for faster comparisons
	mtid uintptr
	stid uintptr
}

func (d *decoderBase) maxInitLen() uint {
	return uint(max(1024, d.h.MaxInitLen))
}

func (d *decoderBase) naked() *fauxUnion {
	return &d.n
}

func (d *decoderBase) fauxUnionReadRawBytes(dr decDriverI, asString, rawToString bool) { //, handleZeroCopy bool) {
	// fauxUnion is only used within DecodeNaked calls; consequently, we should try to intern.
	d.n.l, d.n.a = dr.DecodeBytes()
	if asString || rawToString {
		d.n.v = valueTypeString
		d.n.s = d.detach2Str(d.n.l, d.n.a)
	} else {
		d.n.v = valueTypeBytes
		d.n.l = d.detach2Bytes(d.n.l, d.n.a)
	}
}

// Return a fixed (detached) string representation of a []byte.
//
// Possibly get an interned version of a string,
// iff InternString=true and decoding a map key.
//
// This should mostly be used for map keys, struct field names, etc
// where the key type is string. This is because keys of a map/struct are
// typically reused across many objects.
func (d *decoderBase) detach2Str(v []byte, state dBytesAttachState) (s string) {
	// note: string([]byte) checks - and optimizes - for len 0 and len 1
	if len(v) <= 1 {
		s = string(v)
	} else if state >= dBytesAttachViewZerocopy { // !scratchBuf && d.bytes && d.h.ZeroCopy
		s = stringView(v)
	} else if d.is == nil || d.c != containerMapKey || len(v) > internMaxStrLen {
		s = string(v)
	} else {
		s = d.is.string(v)
	}
	return
}

func (d *decoderBase) usableStructFieldNameBytes(buf, v []byte, state dBytesAttachState) (out []byte) {
	// In JSON, mapElemValue reads a colon and spaces.
	// In bufio mode of ioDecReader, fillbuf could overwrite the read buffer
	// which readXXX() calls return sub-slices from.
	//
	// Consequently, we detach the bytes in this special case.
	//
	// Note: ioDecReader (non-bufio) and bytesDecReader do not have
	// this issue (as no fillbuf exists where bytes might be returned from).
	if d.bufio && d.h.jsonHandle && state < dBytesAttachViewZerocopy {
		if cap(buf) > len(v) {
			out = buf[:len(v)]
		} else if len(d.b) > len(v) {
			out = d.b[:len(v)]
		} else {
			out = make([]byte, len(v), max(64, len(v)))
		}
		copy(out, v)
		return
	}
	return v
}

func (d *decoderBase) detach2Bytes(in []byte, state dBytesAttachState) (out []byte) {
	if cap(in) == 0 || state >= dBytesAttachViewZerocopy {
		return in
	}
	if len(in) == 0 {
		return zeroByteSlice
	}
	out = make([]byte, len(in))
	copy(out, in)
	return out
}

func (d *decoderBase) attachState(usingBufFromReader bool) (r dBytesAttachState) {
	if usingBufFromReader {
		r = dBytesAttachBuffer
	} else if !d.bytes {
		r = dBytesDetach
	} else if d.h.ZeroCopy {
		r = dBytesAttachViewZerocopy
	} else {
		r = dBytesAttachView
	}
	return
}

func (d *decoderBase) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *decoderBase) HandleName() string {
	return d.hh.Name()
}

func (d *decoderBase) isBytes() bool {
	return d.bytes
}

type decoderI interface {
	Decode(v interface{}) (err error)
	HandleName() string
	MustDecode(v interface{})
	NumBytesRead() int
	Release() // deprecated
	Reset(r io.Reader)
	ResetBytes(in []byte)
	ResetString(s string)

	isBytes() bool
	wrapErr(v error, err *error)
	swallow()

	nextValueBytes() []byte // wrapper method, for use in tests
	// getDecDriver() decDriverI

	decode(v interface{})
	decodeAs(v interface{}, t reflect.Type, ext bool)

	interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt)
}

var errDecNoResetBytesWithReader = errors.New("cannot reset an Decoder reading from []byte with a io.Reader")
var errDecNoResetReaderWithBytes = errors.New("cannot reset an Decoder reading from io.Reader with a []byte")

func setZero(iv interface{}) {
	rv, isnil := isNil(iv, false)
	if isnil {
		return
	}
	if !rv.IsValid() {
		rv = reflect.ValueOf(iv)
	}
	if isnilBitset.isset(byte(rv.Kind())) && rvIsNil(rv) {
		return
	}
	// var canDecode bool
	switch v := iv.(type) {
	case *string:
		*v = ""
	case *bool:
		*v = false
	case *int:
		*v = 0
	case *int8:
		*v = 0
	case *int16:
		*v = 0
	case *int32:
		*v = 0
	case *int64:
		*v = 0
	case *uint:
		*v = 0
	case *uint8:
		*v = 0
	case *uint16:
		*v = 0
	case *uint32:
		*v = 0
	case *uint64:
		*v = 0
	case *float32:
		*v = 0
	case *float64:
		*v = 0
	case *complex64:
		*v = 0
	case *complex128:
		*v = 0
	case *[]byte:
		*v = nil
	case *Raw:
		*v = nil
	case *time.Time:
		*v = time.Time{}
	case reflect.Value:
		decSetNonNilRV2Zero(v)
	default:
		if !fastpathDecodeSetZeroTypeSwitch(iv) {
			decSetNonNilRV2Zero(rv)
		}
	}
}

// decSetNonNilRV2Zero will set the non-nil value to its zero value.
func decSetNonNilRV2Zero(v reflect.Value) {
	// If not decodeable (settable), we do not touch it.
	// We considered empty'ing it if not decodeable e.g.
	//    - if chan, drain it
	//    - if map, clear it
	//    - if slice or array, zero all elements up to len
	//
	// However, we decided instead that we either will set the
	// whole value to the zero value, or leave AS IS.

	k := v.Kind()
	if k == reflect.Interface {
		decSetNonNilRV2Zero4Intf(v)
	} else if k == reflect.Ptr {
		decSetNonNilRV2Zero4Ptr(v)
	} else if v.CanSet() {
		rvSetDirectZero(v)
	}
}

func decSetNonNilRV2Zero4Ptr(v reflect.Value) {
	ve := v.Elem()
	if ve.CanSet() {
		rvSetZero(ve) // we can have a pointer to an interface
	} else if v.CanSet() {
		rvSetZero(v)
	}
}

func decSetNonNilRV2Zero4Intf(v reflect.Value) {
	ve := v.Elem()
	if ve.CanSet() {
		rvSetDirectZero(ve) // interfaces always have element as a non-interface
	} else if v.CanSet() {
		rvSetZero(v)
	}
}

func (d *decoderBase) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		halt.errorf("cannot expand array len during decode from %v to %v", any(sliceLen), any(streamLen))
	}
}

//go:noinline
func (d *decoderBase) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		halt.onerror(errCannotDecodeIntoNil)
	}
	// check if an interface can be retrieved, before grabbing an interface
	if !rv.CanInterface() {
		halt.errorf("cannot decode into a value without an interface: %v", rv)
	}
	halt.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *decoderBase) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		halt.onerror(errMaxDepthExceeded)
	}
}

func (d *decoderBase) depthDecr() {
	d.depth--
}

func (d *decoderBase) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *decoderBase) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	// MARKER 2025: is this slow for calling oneShot?
	if decUseTransient && d.h.getTypeInfo4RT(baseRT(rvt)).flagCanTransient {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}

// decNegintPosintFloatNumberHelper is used for formats that are binary
// and have distinct ways of storing positive integers vs negative integers
// vs floats, which are uniquely identified by the byte descriptor.
//
// Currently, these formats are binc, cbor and simple.
type decNegintPosintFloatNumberHelper struct {
	d decDriverI
}

func (x decNegintPosintFloatNumberHelper) uint64(ui uint64, neg, ok bool) uint64 {
	if ok && !neg {
		return ui
	}
	return x.uint64TryFloat(ok)
}

func (x decNegintPosintFloatNumberHelper) uint64TryFloat(neg bool) (ui uint64) {
	if neg { // neg = true
		halt.errorStr("assigning negative signed value to unsigned type")
	}
	f, ok := x.d.decFloat()
	if !(ok && f >= 0 && noFrac64(math.Float64bits(f))) {
		halt.errorStr2("invalid number loading uint64, with descriptor: ", x.d.descBd())
	}
	return uint64(f)
}

func (x decNegintPosintFloatNumberHelper) int64(ui uint64, neg, ok, cbor bool) (i int64) {
	if ok {
		return decNegintPosintFloatNumberHelperInt64v(ui, neg, cbor)
	}
	// 	return x.int64TryFloat()
	// }
	// func (x decNegintPosintFloatNumberHelper) int64TryFloat() (i int64) {
	f, ok := x.d.decFloat()
	if !(ok && noFrac64(math.Float64bits(f))) {
		halt.errorf("invalid number loading uint64 (%v), with descriptor: %s", f, x.d.descBd())
	}
	return int64(f)
}

func (x decNegintPosintFloatNumberHelper) float64(f float64, ok, cbor bool) float64 {
	if ok {
		return f
	}
	return x.float64TryInteger(cbor)
}

func (x decNegintPosintFloatNumberHelper) float64TryInteger(cbor bool) float64 {
	ui, neg, ok := x.d.decInteger()
	if !ok {
		halt.errorStr2("invalid descriptor for float: ", x.d.descBd())
	}
	return float64(decNegintPosintFloatNumberHelperInt64v(ui, neg, cbor))
}

func decNegintPosintFloatNumberHelperInt64v(ui uint64, neg, incrIfNeg bool) (i int64) {
	if neg && incrIfNeg {
		ui++
	}
	i = chkOvf.SignedIntV(ui)
	if neg {
		i = -i
	}
	return
}

// isDecodeable checks if value can be decoded into
//
// decode can take any reflect.Value that is a inherently addressable i.e.
//   - non-nil chan    (we will SEND to it)
//   - non-nil slice   (we will set its elements)
//   - non-nil map     (we will put into it)
//   - non-nil pointer (we can "update" it)
//   - func: no
//   - interface: no
//   - array:                   if canAddr=true
//   - any other value pointer: if canAddr=true
func isDecodeable(rv reflect.Value) (canDecode bool, reason decNotDecodeableReason) {
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Chan, reflect.Map:
		canDecode = !rvIsNil(rv)
		reason = decNotDecodeableReasonNilReference
	case reflect.Func, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
		reason = decNotDecodeableReasonBadKind
	default:
		canDecode = rv.CanAddr()
		reason = decNotDecodeableReasonNonAddrValue
	}
	return
}

// decInferLen will infer a sensible length, given the following:
//   - clen: length wanted.
//   - maxlen: max length to be returned.
//     if <= 0, it is unset, and we infer it based on the unit size
//   - unit: number of bytes for each element of the collection
func decInferLen(clen int, maxlen, unit uint) (n uint) {
	// anecdotal testing showed increase in allocation with map length of 16.
	// We saw same typical alloc from 0-8, then a 20% increase at 16.
	// Thus, we set it to 8.

	const (
		minLenIfUnset = 8
		maxMem        = 1024 * 1024 // 1 MB Memory
	)

	// handle when maxlen is not set i.e. <= 0

	// clen==0:           use 0
	// maxlen<=0, clen<0: use default
	// maxlen> 0, clen<0: use default
	// maxlen<=0, clen>0: infer maxlen, and cap on it
	// maxlen> 0, clen>0: cap at maxlen

	if clen == 0 || clen == containerLenNil {
		return 0
	}
	if clen < 0 {
		// if unspecified, return 64 for bytes, ... 8 for uint64, ... and everything else
		return max(64/unit, minLenIfUnset)
	}
	if unit == 0 {
		return uint(clen)
	}
	if maxlen == 0 {
		maxlen = maxMem / unit
	}
	return min(uint(clen), maxlen)
}

type Decoder struct {
	decoderI
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to configure ReaderBufferSize on the handle
// OR pass in a memory buffered reader (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder {
	return &Decoder{h.newDecoder(r)}
}

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder {
	return &Decoder{h.newDecoderBytes(in)}
}

// NewDecoderString returns a Decoder which efficiently decodes directly
// from a string with zero copying.
//
// It is a convenience function that calls NewDecoderBytes with a
// []byte view into the string.
//
// This can be an efficient zero-copy if using default mode i.e. without codec.safe tag.
func NewDecoderString(s string, h Handle) *Decoder {
	return NewDecoderBytes(bytesView(s), h)
}

// ----

func sideDecode(h Handle, p *sync.Pool, fn func(decoderI)) {
	var s decoderI
	if usePoolForSideDecode {
		s = p.Get().(decoderI)
		defer p.Put(s)
	} else {
		// initialization cycle error
		// s = NewDecoderBytes(nil, h).decoderI
		s = p.New().(decoderI)
	}
	fn(s)
}

func oneOffDecode(sd decoderI, v interface{}, in []byte, basetype reflect.Type, ext bool) {
	sd.ResetBytes(in)
	sd.decodeAs(v, basetype, ext)
	// d.sideDecoder(xbs)
	// d.sideDecode(rv, basetype)
}

func bytesOKdbi(v []byte, _ dBytesIntoState) []byte {
	return v
}

func bytesOKs(bs []byte, _ dBytesAttachState) []byte {
	return bs
}
