// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build notmono || codec.notmono

package codec

import (
	"io"
)

// This contains all the iniatializations of generics.
// Putting it into one file, ensures that we can go generics or not.

type maker interface{ Make() }

func callMake(v interface{}) {
	v.(maker).Make()
}

// ---- (writer.go)

type encWriter interface {
	bufioEncWriterM | bytesEncAppenderM
	encWriterI
}

type bytesEncAppenderM struct {
	*bytesEncAppender
}

func (z *bytesEncAppenderM) Make() {
	z.bytesEncAppender = new(bytesEncAppender)
	z.out = &bytesEncAppenderDefOut
}

type bufioEncWriterM struct {
	*bufioEncWriter
}

func (z *bufioEncWriterM) Make() {
	z.bufioEncWriter = new(bufioEncWriter)
	z.w = io.Discard
}

// ---- reader.go

type decReader interface {
	bytesDecReaderM | ioDecReaderM

	decReaderI
}

type bytesDecReaderM struct {
	*bytesDecReader
}

func (z *bytesDecReaderM) Make() {
	z.bytesDecReader = new(bytesDecReader)
}

type ioDecReaderM struct {
	*ioDecReader
}

func (z *ioDecReaderM) Make() {
	z.ioDecReader = new(ioDecReader)
}

// type helperEncWriter[T encWriter] struct{}
// type helperDecReader[T decReader] struct{}
// func (helperDecReader[T]) decByteSlice(r T, clen, maxInitLen int, bs []byte) (bsOut []byte) {

// ---- (encode.go)

type encDriver interface {
	simpleEncDriverM[bufioEncWriterM] |
		simpleEncDriverM[bytesEncAppenderM] |
		jsonEncDriverM[bufioEncWriterM] |
		jsonEncDriverM[bytesEncAppenderM] |
		cborEncDriverM[bufioEncWriterM] |
		cborEncDriverM[bytesEncAppenderM] |
		msgpackEncDriverM[bufioEncWriterM] |
		msgpackEncDriverM[bytesEncAppenderM] |
		bincEncDriverM[bufioEncWriterM] |
		bincEncDriverM[bytesEncAppenderM]

	encDriverI
}

// ---- (decode.go)

type decDriver interface {
	simpleDecDriverM[bytesDecReaderM] |
		simpleDecDriverM[ioDecReaderM] |
		jsonDecDriverM[bytesDecReaderM] |
		jsonDecDriverM[ioDecReaderM] |
		cborDecDriverM[bytesDecReaderM] |
		cborDecDriverM[ioDecReaderM] |
		msgpackDecDriverM[bytesDecReaderM] |
		msgpackDecDriverM[ioDecReaderM] |
		bincDecDriverM[bytesDecReaderM] |
		bincDecDriverM[ioDecReaderM]

	decDriverI
}

// Below: <format>.go files

// ---- (binc.go)

type bincEncDriverM[T encWriter] struct {
	*bincEncDriver[T]
}

func (d *bincEncDriverM[T]) Make() {
	d.bincEncDriver = new(bincEncDriver[T])
}

type bincDecDriverM[T decReader] struct {
	*bincDecDriver[T]
}

func (d *bincDecDriverM[T]) Make() {
	d.bincDecDriver = new(bincDecDriver[T])
}

var (
	bincFpEncIO    = helperEncDriver[bincEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	bincFpEncBytes = helperEncDriver[bincEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	bincFpDecIO    = helperDecDriver[bincDecDriverM[ioDecReaderM]]{}.fastpathDList()
	bincFpDecBytes = helperDecDriver[bincDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (cbor.go)

type cborEncDriverM[T encWriter] struct {
	*cborEncDriver[T]
}

func (d *cborEncDriverM[T]) Make() {
	d.cborEncDriver = new(cborEncDriver[T])
}

type cborDecDriverM[T decReader] struct {
	*cborDecDriver[T]
}

func (d *cborDecDriverM[T]) Make() {
	d.cborDecDriver = new(cborDecDriver[T])
}

var (
	cborFpEncIO    = helperEncDriver[cborEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	cborFpEncBytes = helperEncDriver[cborEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	cborFpDecIO    = helperDecDriver[cborDecDriverM[ioDecReaderM]]{}.fastpathDList()
	cborFpDecBytes = helperDecDriver[cborDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (json.go)

type jsonEncDriverM[T encWriter] struct {
	*jsonEncDriver[T]
}

func (d *jsonEncDriverM[T]) Make() {
	d.jsonEncDriver = new(jsonEncDriver[T])
}

type jsonDecDriverM[T decReader] struct {
	*jsonDecDriver[T]
}

func (d *jsonDecDriverM[T]) Make() {
	d.jsonDecDriver = new(jsonDecDriver[T])
}

var (
	jsonFpEncIO    = helperEncDriver[jsonEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	jsonFpEncBytes = helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	jsonFpDecIO    = helperDecDriver[jsonDecDriverM[ioDecReaderM]]{}.fastpathDList()
	jsonFpDecBytes = helperDecDriver[jsonDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (msgpack.go)

type msgpackEncDriverM[T encWriter] struct {
	*msgpackEncDriver[T]
}

func (d *msgpackEncDriverM[T]) Make() {
	d.msgpackEncDriver = new(msgpackEncDriver[T])
}

type msgpackDecDriverM[T decReader] struct {
	*msgpackDecDriver[T]
}

func (d *msgpackDecDriverM[T]) Make() {
	d.msgpackDecDriver = new(msgpackDecDriver[T])
}

var (
	msgpackFpEncIO    = helperEncDriver[msgpackEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriver[msgpackDecDriverM[ioDecReaderM]]{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

// ---- (simple.go)

type simpleEncDriverM[T encWriter] struct {
	*simpleEncDriver[T]
}

func (d *simpleEncDriverM[T]) Make() {
	d.simpleEncDriver = new(simpleEncDriver[T])
}

type simpleDecDriverM[T decReader] struct {
	*simpleDecDriver[T]
}

func (d *simpleDecDriverM[T]) Make() {
	d.simpleDecDriver = new(simpleDecDriver[T])
}

var (
	simpleFpEncIO    = helperEncDriver[simpleEncDriverM[bufioEncWriterM]]{}.fastpathEList()
	simpleFpEncBytes = helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]{}.fastpathEList()
	simpleFpDecIO    = helperDecDriver[simpleDecDriverM[ioDecReaderM]]{}.fastpathDList()
	simpleFpDecBytes = helperDecDriver[simpleDecDriverM[bytesDecReaderM]]{}.fastpathDList()
)

func (h *SimpleHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[simpleEncDriverM[bytesEncAppenderM]]{}.newEncoderBytes(out, h)
}

func (h *SimpleHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[simpleEncDriverM[bufioEncWriterM]]{}.newEncoderIO(w, h)
}

func (h *SimpleHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[simpleDecDriverM[bytesDecReaderM]]{}.newDecoderBytes(in, h)
}

func (h *SimpleHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[simpleDecDriverM[ioDecReaderM]]{}.newDecoderIO(r, h)
}

func (h *JsonHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[jsonEncDriverM[bytesEncAppenderM]]{}.newEncoderBytes(out, h)
}

func (h *JsonHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[jsonEncDriverM[bufioEncWriterM]]{}.newEncoderIO(w, h)
}

func (h *JsonHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[jsonDecDriverM[bytesDecReaderM]]{}.newDecoderBytes(in, h)
}

func (h *JsonHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[jsonDecDriverM[ioDecReaderM]]{}.newDecoderIO(r, h)
}

func (h *MsgpackHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[msgpackEncDriverM[bytesEncAppenderM]]{}.newEncoderBytes(out, h)
}

func (h *MsgpackHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[msgpackEncDriverM[bufioEncWriterM]]{}.newEncoderIO(w, h)
}

func (h *MsgpackHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[msgpackDecDriverM[bytesDecReaderM]]{}.newDecoderBytes(in, h)
}

func (h *MsgpackHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[msgpackDecDriverM[ioDecReaderM]]{}.newDecoderIO(r, h)
}

func (h *CborHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[cborEncDriverM[bytesEncAppenderM]]{}.newEncoderBytes(out, h)
}

func (h *CborHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[cborEncDriverM[bufioEncWriterM]]{}.newEncoderIO(w, h)
}

func (h *CborHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[cborDecDriverM[bytesDecReaderM]]{}.newDecoderBytes(in, h)
}

func (h *CborHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[cborDecDriverM[ioDecReaderM]]{}.newDecoderIO(r, h)
}

func (h *BincHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriver[bincEncDriverM[bytesEncAppenderM]]{}.newEncoderBytes(out, h)
}

func (h *BincHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriver[bincEncDriverM[bufioEncWriterM]]{}.newEncoderIO(w, h)
}

func (h *BincHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriver[bincDecDriverM[bytesDecReaderM]]{}.newDecoderBytes(in, h)
}

func (h *BincHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriver[bincDecDriverM[ioDecReaderM]]{}.newDecoderIO(r, h)
}
