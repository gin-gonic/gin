// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !notmono && !codec.notmono

package codec

import "io"

func callMake(v interface{}) {}

type encWriter interface{ encWriterI }
type decReader interface{ decReaderI }
type encDriver interface{ encDriverI }
type decDriver interface{ decDriverI }

func (h *SimpleHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverSimpleBytes{}.newEncoderBytes(out, h)
}

func (h *SimpleHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverSimpleIO{}.newEncoderIO(w, h)
}

func (h *SimpleHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverSimpleBytes{}.newDecoderBytes(in, h)
}

func (h *SimpleHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverSimpleIO{}.newDecoderIO(r, h)
}

func (h *JsonHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverJsonBytes{}.newEncoderBytes(out, h)
}

func (h *JsonHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverJsonIO{}.newEncoderIO(w, h)
}

func (h *JsonHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverJsonBytes{}.newDecoderBytes(in, h)
}

func (h *JsonHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverJsonIO{}.newDecoderIO(r, h)
}

func (h *MsgpackHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverMsgpackBytes{}.newEncoderBytes(out, h)
}

func (h *MsgpackHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverMsgpackIO{}.newEncoderIO(w, h)
}

func (h *MsgpackHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverMsgpackBytes{}.newDecoderBytes(in, h)
}

func (h *MsgpackHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverMsgpackIO{}.newDecoderIO(r, h)
}

func (h *BincHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverBincBytes{}.newEncoderBytes(out, h)
}

func (h *BincHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverBincIO{}.newEncoderIO(w, h)
}

func (h *BincHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverBincBytes{}.newDecoderBytes(in, h)
}

func (h *BincHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverBincIO{}.newDecoderIO(r, h)
}

func (h *CborHandle) newEncoderBytes(out *[]byte) encoderI {
	return helperEncDriverCborBytes{}.newEncoderBytes(out, h)
}

func (h *CborHandle) newEncoder(w io.Writer) encoderI {
	return helperEncDriverCborIO{}.newEncoderIO(w, h)
}

func (h *CborHandle) newDecoderBytes(in []byte) decoderI {
	return helperDecDriverCborBytes{}.newDecoderBytes(in, h)
}

func (h *CborHandle) newDecoder(r io.Reader) decoderI {
	return helperDecDriverCborIO{}.newDecoderIO(r, h)
}

var (
	bincFpEncIO    = helperEncDriverBincIO{}.fastpathEList()
	bincFpEncBytes = helperEncDriverBincBytes{}.fastpathEList()
	bincFpDecIO    = helperDecDriverBincIO{}.fastpathDList()
	bincFpDecBytes = helperDecDriverBincBytes{}.fastpathDList()
)

var (
	cborFpEncIO    = helperEncDriverCborIO{}.fastpathEList()
	cborFpEncBytes = helperEncDriverCborBytes{}.fastpathEList()
	cborFpDecIO    = helperDecDriverCborIO{}.fastpathDList()
	cborFpDecBytes = helperDecDriverCborBytes{}.fastpathDList()
)

var (
	jsonFpEncIO    = helperEncDriverJsonIO{}.fastpathEList()
	jsonFpEncBytes = helperEncDriverJsonBytes{}.fastpathEList()
	jsonFpDecIO    = helperDecDriverJsonIO{}.fastpathDList()
	jsonFpDecBytes = helperDecDriverJsonBytes{}.fastpathDList()
)

var (
	msgpackFpEncIO    = helperEncDriverMsgpackIO{}.fastpathEList()
	msgpackFpEncBytes = helperEncDriverMsgpackBytes{}.fastpathEList()
	msgpackFpDecIO    = helperDecDriverMsgpackIO{}.fastpathDList()
	msgpackFpDecBytes = helperDecDriverMsgpackBytes{}.fastpathDList()
)

var (
	simpleFpEncIO    = helperEncDriverSimpleIO{}.fastpathEList()
	simpleFpEncBytes = helperEncDriverSimpleBytes{}.fastpathEList()
	simpleFpDecIO    = helperDecDriverSimpleIO{}.fastpathDList()
	simpleFpDecBytes = helperDecDriverSimpleBytes{}.fastpathDList()
)
