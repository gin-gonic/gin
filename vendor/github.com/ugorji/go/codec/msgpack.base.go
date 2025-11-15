// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"fmt"
	"io"
	"net/rpc"
	"reflect"
)

const (
	mpPosFixNumMin byte = 0x00
	mpPosFixNumMax byte = 0x7f
	mpFixMapMin    byte = 0x80
	mpFixMapMax    byte = 0x8f
	mpFixArrayMin  byte = 0x90
	mpFixArrayMax  byte = 0x9f
	mpFixStrMin    byte = 0xa0
	mpFixStrMax    byte = 0xbf
	mpNil          byte = 0xc0
	_              byte = 0xc1
	mpFalse        byte = 0xc2
	mpTrue         byte = 0xc3
	mpFloat        byte = 0xca
	mpDouble       byte = 0xcb
	mpUint8        byte = 0xcc
	mpUint16       byte = 0xcd
	mpUint32       byte = 0xce
	mpUint64       byte = 0xcf
	mpInt8         byte = 0xd0
	mpInt16        byte = 0xd1
	mpInt32        byte = 0xd2
	mpInt64        byte = 0xd3

	// extensions below
	mpBin8     byte = 0xc4
	mpBin16    byte = 0xc5
	mpBin32    byte = 0xc6
	mpExt8     byte = 0xc7
	mpExt16    byte = 0xc8
	mpExt32    byte = 0xc9
	mpFixExt1  byte = 0xd4
	mpFixExt2  byte = 0xd5
	mpFixExt4  byte = 0xd6
	mpFixExt8  byte = 0xd7
	mpFixExt16 byte = 0xd8

	mpStr8  byte = 0xd9 // new
	mpStr16 byte = 0xda
	mpStr32 byte = 0xdb

	mpArray16 byte = 0xdc
	mpArray32 byte = 0xdd

	mpMap16 byte = 0xde
	mpMap32 byte = 0xdf

	mpNegFixNumMin byte = 0xe0
	mpNegFixNumMax byte = 0xff
)

var mpTimeExtTag int8 = -1
var mpTimeExtTagU = uint8(mpTimeExtTag)

var mpdescNames = map[byte]string{
	mpNil:    "nil",
	mpFalse:  "false",
	mpTrue:   "true",
	mpFloat:  "float",
	mpDouble: "float",
	mpUint8:  "uuint",
	mpUint16: "uint",
	mpUint32: "uint",
	mpUint64: "uint",
	mpInt8:   "int",
	mpInt16:  "int",
	mpInt32:  "int",
	mpInt64:  "int",

	mpStr8:  "string|bytes",
	mpStr16: "string|bytes",
	mpStr32: "string|bytes",

	mpBin8:  "bytes",
	mpBin16: "bytes",
	mpBin32: "bytes",

	mpArray16: "array",
	mpArray32: "array",

	mpMap16: "map",
	mpMap32: "map",
}

func mpdesc(bd byte) (s string) {
	s = mpdescNames[bd]
	if s == "" {
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax,
			bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
			s = "int"
		case bd >= mpFixStrMin && bd <= mpFixStrMax:
			s = "string|bytes"
		case bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			s = "array"
		case bd >= mpFixMapMin && bd <= mpFixMapMax:
			s = "map"
		case bd >= mpFixExt1 && bd <= mpFixExt16,
			bd >= mpExt8 && bd <= mpExt32:
			s = "ext"
		default:
			s = "unknown"
		}
	}
	return
}

// MsgpackSpecRpcMultiArgs is a special type which signifies to the MsgpackSpecRpcCodec
// that the backend RPC service takes multiple arguments, which have been arranged
// in sequence in the slice.
//
// The Codec then passes it AS-IS to the rpc service (without wrapping it in an
// array of 1 element).
type MsgpackSpecRpcMultiArgs []interface{}

// A MsgpackContainer type specifies the different types of msgpackContainers.
type msgpackContainerType struct {
	fixCutoff, bFixMin, b8, b16, b32 byte
	// hasFixMin, has8, has8Always bool
}

var (
	msgpackContainerRawLegacy = msgpackContainerType{
		32, mpFixStrMin, 0, mpStr16, mpStr32,
	}
	msgpackContainerStr = msgpackContainerType{
		32, mpFixStrMin, mpStr8, mpStr16, mpStr32, // true, true, false,
	}
	msgpackContainerBin = msgpackContainerType{
		0, 0, mpBin8, mpBin16, mpBin32, // false, true, true,
	}
	msgpackContainerList = msgpackContainerType{
		16, mpFixArrayMin, 0, mpArray16, mpArray32, // true, false, false,
	}
	msgpackContainerMap = msgpackContainerType{
		16, mpFixMapMin, 0, mpMap16, mpMap32, // true, false, false,
	}
)

//--------------------------------------------------

// MsgpackHandle is a Handle for the Msgpack Schema-Free Encoding Format.
type MsgpackHandle struct {
	binaryEncodingType
	notJsonType
	BasicHandle

	// NoFixedNum says to output all signed integers as 2-bytes, never as 1-byte fixednum.
	NoFixedNum bool

	// WriteExt controls whether the new spec is honored.
	//
	// With WriteExt=true, we can encode configured extensions with extension tags
	// and encode string/[]byte/extensions in a way compatible with the new spec
	// but incompatible with the old spec.
	//
	// For compatibility with the old spec, set WriteExt=false.
	//
	// With WriteExt=false:
	//    configured extensions are serialized as raw bytes (not msgpack extensions).
	//    reserved byte descriptors like Str8 and those enabling the new msgpack Binary type
	//    are not encoded.
	WriteExt bool

	// PositiveIntUnsigned says to encode positive integers as unsigned.
	PositiveIntUnsigned bool
}

// Name returns the name of the handle: msgpack
func (h *MsgpackHandle) Name() string { return "msgpack" }

func (h *MsgpackHandle) desc(bd byte) string { return mpdesc(bd) }

// SetBytesExt sets an extension
func (h *MsgpackHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

//--------------------------------------------------

type msgpackSpecRpcCodec struct {
	*rpcCodec
}

// /////////////// Spec RPC Codec ///////////////////
func (c *msgpackSpecRpcCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	// WriteRequest can write to both a Go service, and other services that do
	// not abide by the 1 argument rule of a Go service.
	// We discriminate based on if the body is a MsgpackSpecRpcMultiArgs
	var bodyArr []interface{}
	if m, ok := body.(MsgpackSpecRpcMultiArgs); ok {
		bodyArr = ([]interface{})(m)
	} else {
		bodyArr = []interface{}{body}
	}
	r2 := []interface{}{0, uint32(r.Seq), r.ServiceMethod, bodyArr}
	return c.write(r2)
}

func (c *msgpackSpecRpcCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	var moe interface{}
	if r.Error != "" {
		moe = r.Error
	}
	if moe != nil && body != nil {
		body = nil
	}
	r2 := []interface{}{1, uint32(r.Seq), moe, body}
	return c.write(r2)
}

func (c *msgpackSpecRpcCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.parseCustomHeader(1, &r.Seq, &r.Error)
}

func (c *msgpackSpecRpcCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.parseCustomHeader(0, &r.Seq, &r.ServiceMethod)
}

func (c *msgpackSpecRpcCodec) ReadRequestBody(body interface{}) error {
	if body == nil { // read and discard
		return c.read(nil)
	}
	bodyArr := []interface{}{body}
	return c.read(&bodyArr)
}

func (c *msgpackSpecRpcCodec) parseCustomHeader(expectTypeByte byte, msgid *uint64, methodOrError *string) (err error) {
	if c.cls.Load().closed {
		return io.ErrUnexpectedEOF
	}

	// We read the response header by hand
	// so that the body can be decoded on its own from the stream at a later time.

	const fia byte = 0x94 //four item array descriptor value

	var ba [1]byte
	var n int
	for {
		n, err = c.r.Read(ba[:])
		if err != nil {
			return
		}
		if n == 1 {
			break
		}
	}

	var b = ba[0]
	if b != fia {
		err = fmt.Errorf("not array - %s %x/%s", msgBadDesc, b, mpdesc(b))
	} else {
		err = c.read(&b)
		if err == nil {
			if b != expectTypeByte {
				err = fmt.Errorf("%s - expecting %v but got %x/%s", msgBadDesc, expectTypeByte, b, mpdesc(b))
			} else {
				err = c.read(msgid)
				if err == nil {
					err = c.read(methodOrError)
				}
			}
		}
	}
	return
}

//--------------------------------------------------

// msgpackSpecRpc is the implementation of Rpc that uses custom communication protocol
// as defined in the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md
type msgpackSpecRpc struct{}

// MsgpackSpecRpc implements Rpc using the communication protocol defined in
// the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md .
//
// See GoRpc documentation, for information on buffering for better performance.
var MsgpackSpecRpc msgpackSpecRpc

func (x msgpackSpecRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec {
	return &msgpackSpecRpcCodec{newRPCCodec(conn, h)}
}

func (x msgpackSpecRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec {
	return &msgpackSpecRpcCodec{newRPCCodec(conn, h)}
}
