// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"errors"
	"io"
	"os"
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReaderI interface {
	// readx will return a view of the []byte in one of 2 ways:
	//   - direct view into []byte which decoding is happening from (if bytes)
	//   - view into a mutable []byte which the ioReader is using (if IO)
	//
	// Users should directly consume the contents read, and not store for future use.
	readx(n uint) []byte

	// skip n bytes
	skip(n uint)

	readb([]byte)

	// readxb will read n bytes, returning as out, and a flag stating whether
	// an internal buffer (not the view) was used.
	readxb(n uint) (out []byte, usingBuf bool)

	readn1() byte
	readn2() [2]byte
	readn3() [3]byte
	readn4() [4]byte
	readn8() [8]byte
	// readn1eof() (v uint8, eof bool)

	// // read up to 8 bytes at a time
	// readn(num uint8) (v [8]byte)

	numread() uint // number of bytes read

	// skip any whitespace characters, and return the first non-matching byte
	skipWhitespace() (token byte)

	// jsonReadNum will read a sequence of numeric characters, checking from the last
	// read byte. It will return a sequence of numeric characters (v),
	// and the next token character (tok - returned separately),
	//
	// if an EOF is found before the next token is seen, it returns a token value of 0.
	jsonReadNum() (v []byte, token byte)

	// jsonReadAsisChars recognizes 2 terminal characters (" or \).
	// jsonReadAsisChars will read json plain characters until it reaches a terminal char,
	// and returns a slice up to the terminal char (excluded),
	// and also returns the terminal char separately (" or \).
	jsonReadAsisChars() (v []byte, terminal byte)

	// readUntil will read characters until it reaches a ",
	// return a slice up to " (excluded)
	jsonReadUntilDblQuote() (v []byte)

	// skip will skip any byte that matches, and return the first non-matching byte
	// skip(accept *bitset256) (token byte)

	// readTo will read any byte that matches, stopping once no-longer matching.
	// readTo(accept *bitset256) (out []byte)

	// // readUntil will read characters until it reaches a stop char,
	// // return a slice up to the terminal byte (excluded)
	// readUntil(stop byte) (out []byte)

	// only supported when reading from bytes
	// bytesReadFrom(startpos uint) []byte

	// isBytes() bool
	resetIO(r io.Reader, bufsize int, maxInitLen int, blist *bytesFreeList)

	resetBytes(in []byte)

	// nextValueBytes() captures bytes read between a call to startRecording and stopRecording.
	// startRecording will always includes the last byte read.
	startRecording()
	// stopRecording will include all bytes read between the point of startRecording and now.
	stopRecording() []byte
}

// // ------------------------------------------------

const maxConsecutiveEmptyReads = 16 // 2 is sufficient, 16 is enough, 64 is optimal

// const defBufReaderSize = 4096

// --------------------

// ioReaderByteScanner contains the io.Reader and io.ByteScanner interfaces
type ioReaderByteScanner interface {
	io.Reader
	io.ByteScanner
}

// MARKER: why not separate bufioDecReader from ioDecReader?
//
// We tried, but only readn1 of bufioDecReader came close to being
// inlined (at inline cost 82). All other methods were at inline cost >= 90.
//
// Consequently, there's no performance impact from having both together
// (except a single if z.bufio branch, which is likely well predicted and happens
// only once per call (right at the top).

// ioDecReader is a decReader that reads off an io.Reader.
type ioDecReader struct {
	r io.Reader

	blist *bytesFreeList

	maxInitLen uint

	n uint // num read

	bufsize uint

	bufio     bool // are we buffering (rc and wc are valid)
	rbr       bool // r is a byte reader
	recording bool // are we recording (src and erc are valid)
	done      bool // did we reach EOF and are we done?

	// valid when: bufio=false
	b  [1]byte       // tiny buffer for reading single byte (if z.br == nil)
	l  byte          // last byte read
	br io.ByteReader // main reader used for ReadByte

	// valid when: bufio=true
	wc  uint // read cursor
	rc  uint // write cursor
	err error

	// valid when: recording=true
	recc uint // start-recording cursor (valid: recording=true)

	buf []byte // buffer for bufio OR recording (if !bufio)
}

func (z *ioDecReader) resetBytes(in []byte) {
	halt.errorStr("resetBytes unsupported by ioDecReader")
}

func (z *ioDecReader) resetIO(r io.Reader, bufsize int, maxInitLen int, blist *bytesFreeList) {
	buf := z.buf
	*z = ioDecReader{}
	z.maxInitLen = max(1024, uint(maxInitLen))
	z.blist = blist
	z.buf = blist.check(buf, max(256, bufsize))
	z.bufsize = uint(max(0, bufsize))
	z.bufio = z.bufsize > 0
	if z.bufio {
		z.buf = z.buf[:cap(z.buf)]
	} else {
		z.buf = z.buf[:0]
	}
	if r == nil {
		z.r = &eofReader
	} else {
		z.r = r
	}
	z.br, z.rbr = z.r.(io.ByteReader)
}

func (z *ioDecReader) numread() uint {
	return z.n
}

func (z *ioDecReader) readn2() [2]byte {
	return ([2]byte)(z.readx(2))
	// using readb forced return bs onto heap, unnecessarily
	// z.readb(bs[:])
	// return
}

func (z *ioDecReader) readn3() [3]byte {
	return ([3]byte)(z.readx(3))
}

func (z *ioDecReader) readn4() [4]byte {
	return ([4]byte)(z.readx(4))
}

func (z *ioDecReader) readn8() [8]byte {
	return ([8]byte)(z.readx(8))
}

func (z *ioDecReader) readx(n uint) (bs []byte) {
	return bytesOK(z.readxb(n))
}

func (z *ioDecReader) readErr() (err error) {
	err, z.err = z.err, nil
	return
}

func (z *ioDecReader) checkErr() {
	halt.onerror(z.readErr())
}

func (z *ioDecReader) unexpectedEOF() {
	z.checkErr()
	// if no error, still halt with unexpected EOF
	halt.error(io.ErrUnexpectedEOF)
}

func (z *ioDecReader) readOne() (b byte, err error) {
	n, err := z.r.Read(z.b[:])
	if n == 1 {
		err = nil
		b = z.b[0]
	}
	return
}

// fillbuf reads a new chunk into the buffer.
func (z *ioDecReader) fillbuf(bufsize uint) (numShift, numRead uint) {
	z.checkErr()
	bufsize = max(bufsize, z.bufsize)

	// Slide existing data to beginning.
	if z.recording {
		numShift = z.recc // recc is always <= rc
	} else {
		numShift = z.rc
	}
	if numShift > 0 {
		numShift-- // never shift last byte read out
	}
	copy(z.buf, z.buf[numShift:z.wc])
	z.wc -= numShift
	z.rc -= numShift
	if z.recording {
		z.recc -= numShift
	}
	// add enough to allow u to read up to bufsize again iff
	// - buf is fully written
	// - NOTE: don't pre-allocate more until needed
	if uint(len(z.buf)) == z.wc {
		if bufsize+z.wc < uint(cap(z.buf)) {
			z.buf = z.buf[:uint(cap(z.buf))]
		} else {
			bufsize = max(uint(cap(z.buf)*3/2), bufsize+z.wc)
			buf := z.blist.get(int(bufsize))
			buf = buf[:cap(buf)]
			copy(buf, z.buf[:z.wc])
			z.blist.put(z.buf)
			z.buf = buf
		}
	}
	// Read new data: try a limited number of times.
	// if n == 0: try up to maxConsecutiveEmptyReads
	// if n > 0 and err == nil: try one more time (to see if we get n == 0 and EOF)
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := z.r.Read(z.buf[z.wc:])
		numRead += uint(n)
		z.wc += uint(n)
		if err != nil {
			z.err = err
			if err == io.EOF {
				z.done = true // leading to UnexpectedEOF if another Read is called
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				// os read deadline, but some bytes read: return (don't store err)
				z.err = nil // allow for a retry next time fillbuf is called
			}
			return
		}

		// if z.wc == uint(len(z.buf)) {
		// 	return
		// }
		// only read one time if results returned
		// if n > 0 && i > 2 {
		// 	i = 2 // try max one more time (to see about getting EOF)
		// }

		// Once you have some data from this read call, move on.
		// Consequently, a blocked Read has less chance of happening.
		if n > 0 {
			return
		}
	}
	z.err = io.ErrNoProgress // either no data read OR not enough data read, without an EOF
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	var err error
	var n int
	if z.bufio {
	BUFIO:
		for z.rc == z.wc {
			z.fillbuf(0)
		}
		n = copy(bs, z.buf[z.rc:z.wc])
		z.rc += uint(n)
		z.n += uint(n)
		if n == len(bs) {
			return
		}
		bs = bs[n:]
		goto BUFIO
	}

	// -------- NOT BUFIO ------

	var nn uint
	bs0 := bs
READER:
	n, err = z.r.Read(bs)
	if n > 0 {
		z.l = bs[n-1]
		nn += uint(n)
		bs = bs[n:]
	}
	if len(bs) != 0 && err == nil {
		goto READER
	}
	if z.recording {
		z.buf = append(z.buf, bs0[:nn]...)
	}
	z.n += nn
	if len(bs) != 0 {
		halt.onerror(err)
		halt.errorf("ioDecReader.readb read %d out of %d bytes requested", nn, len(bs0))
	}
	return
}

func (z *ioDecReader) readn1() (b uint8) {
	if z.bufio {
		for z.rc == z.wc {
			z.fillbuf(0)
		}
		b = z.buf[z.rc]
		z.rc++
		z.n++
		return
	}

	// -------- NOT BUFIO ------

	var err error
	if z.rbr {
		b, err = z.br.ReadByte()
	} else {
		b, err = z.readOne()
	}
	halt.onerror(err)
	z.l = b
	z.n++
	if z.recording {
		z.buf = append(z.buf, b)
	}
	return
}

func (z *ioDecReader) readxb(n uint) (out []byte, useBuf bool) {
	if n == 0 {
		return zeroByteSlice, false
	}

	if z.bufio {
	BUFIO:
		nn := int(n+z.rc) - int(z.wc)
		if nn > 0 {
			z.fillbuf(decInferLen(nn, z.maxInitLen, 1))
			goto BUFIO
		}
		pos := z.rc
		z.rc += uint(n)
		z.n += uint(n)
		out = z.buf[pos:z.rc]
		useBuf = true
		return
	}

	// -------- NOT BUFIO ------

	var n3 int
	var err error
	useBuf = true
	out = z.buf
	r0 := uint(len(out))
	r := r0
	nn := int(n)
	for nn > 0 {
		halt.onerror(err) // check error whenever there's more to read
		n2 := r + decInferLen(int(nn), z.maxInitLen, 1)
		if cap(out) < int(n2) {
			out2 := z.blist.putGet(out, int(n2))[:n2] // make([]byte, len2+len3)
			copy(out2, out)
			out = out2
		} else {
			out = out[:n2]
		}
		n3, err = z.r.Read(out[r:n2])
		if n3 > 0 {
			z.l = out[r+uint(n3)-1]
			nn -= n3
			r += uint(n3)
		}
	}
	z.buf = out[:r0+n]
	out = out[r0 : r0+n]
	z.n += n
	return
}

func (z *ioDecReader) skip(n uint) {
	if n == 0 {
		return
	}

	if z.bufio {
	BUFIO:
		n2 := min(n, z.wc-z.rc)
		// handle in-line, so z.buf doesn't grow much (since we're skipping)
		// ie by setting z.rc, fillbuf should keep shifting left (unless recording)
		z.rc += n2
		z.n += n2
		n -= n2
		if n > 0 {
			z.fillbuf(decInferLen(int(n+z.rc)-int(z.wc), z.maxInitLen, 1))
			goto BUFIO
		}
		return
	}

	// -------- NOT BUFIO ------

	var out []byte
	var fromBlist bool
	if z.recording {
		out = z.buf
	} else {
		nn := int(decInferLen(int(n), z.maxInitLen, 1))
		if cap(z.buf) >= nn/2 {
			out = z.buf[:cap(z.buf)]
		} else {
			fromBlist = true
			out = z.blist.get(nn)
		}
	}

	var r uint
	var n3 int
	var err error
	nn := int(n)
	for nn > 0 {
		halt.onerror(err)
		n2 := uint(nn)
		if z.recording {
			r = uint(len(out))
			n2 = r + decInferLen(int(nn), z.maxInitLen, 1)
			if cap(out) < int(n2) {
				out2 := z.blist.putGet(out, int(n2))[:n2] // make([]byte, len2+len3)
				copy(out2, out)
				out = out2
			} else {
				out = out[:n2]
			}
		}
		n3, err = z.r.Read(out[r:n2])
		if n3 > 0 {
			z.l = out[r+uint(n3)-1]
			z.n += uint(n3)
			nn -= n3
		}
	}
	if z.recording {
		z.buf = out
	} else if fromBlist {
		z.blist.put(out)
	}
	return
}

// ---- JSON SPECIFIC HELPERS HERE ----

func (z *ioDecReader) jsonReadNum() (bs []byte, token byte) {
	var start, pos, end uint
	if z.bufio {
		// read and fill into buf, then take substring
		start = z.rc - 1 // include last byte read
		pos = start
	BUFIO:
		if pos == z.wc {
			if z.done {
				end = pos
				goto END
			}
			numshift, numread := z.fillbuf(0)
			start -= numshift
			pos -= numshift
			if numread == 0 {
				end = pos
				goto END
			}
		}
		token = z.buf[pos]
		pos++
		if isNumberChar(token) {
			goto BUFIO
		}
		end = pos - 1
	END:
		z.n += (pos - z.rc)
		z.rc = pos
		return z.buf[start:end], token
	}

	// if not recording, add the last read byte into buf
	if !z.recording {
		z.buf = append(z.buf[:0], z.l)
	}
	start = uint(len(z.buf) - 1) // incl last byte in z.buf
	var b byte
	var err error

READER:
	if z.rbr {
		b, err = z.br.ReadByte()
	} else {
		b, err = z.readOne()
	}
	if err == io.EOF {
		return z.buf[start:], 0
	}
	halt.onerror(err)
	z.l = b
	z.n++
	z.buf = append(z.buf, b)
	if isNumberChar(b) {
		goto READER
	}
	return z.buf[start : len(z.buf)-1], b
}

func (z *ioDecReader) skipWhitespace() (tok byte) {
	var pos uint
	if z.bufio {
		pos = z.rc
	BUFIO:
		if pos == z.wc {
			if z.done {
				z.unexpectedEOF()
			}
			numshift, numread := z.fillbuf(0)
			pos -= numshift
			if numread == 0 {
				z.unexpectedEOF()
			}
		}
		tok = z.buf[pos]
		pos++
		if isWhitespaceChar(tok) {
			goto BUFIO
		}
		z.n += (pos - z.rc)
		z.rc = pos
		return tok
	}

	var err error
READER:
	if z.rbr {
		tok, err = z.br.ReadByte()
	} else {
		tok, err = z.readOne()
	}
	halt.onerror(err)
	z.n++
	z.l = tok
	if z.recording {
		z.buf = append(z.buf, tok)
	}
	if isWhitespaceChar(tok) {
		goto READER
	}
	return tok
}

func (z *ioDecReader) readUntil(stop1, stop2 byte) (bs []byte, tok byte) {
	var start, pos uint
	if z.bufio {
		start = z.rc
		pos = start
	BUFIO:
		if pos == z.wc {
			if z.done {
				z.unexpectedEOF()
			}
			numshift, numread := z.fillbuf(0)
			start -= numshift
			pos -= numshift
			if numread == 0 {
				z.unexpectedEOF()
			}
		}
		tok = z.buf[pos]
		pos++
		if tok == stop1 || tok == stop2 {
			z.n += (pos - z.rc)
			z.rc = pos
			return z.buf[start : pos-1], tok
		}
		goto BUFIO
	}

	var err error
	if !z.recording {
		z.buf = z.buf[:0]
	}
	start = uint(len(z.buf))
READER:
	if z.rbr {
		tok, err = z.br.ReadByte()
	} else {
		tok, err = z.readOne()
	}
	halt.onerror(err)
	z.n++
	z.l = tok
	z.buf = append(z.buf, tok)
	if tok == stop1 || tok == stop2 {
		return z.buf[start : len(z.buf)-1], tok
	}
	goto READER
}

func (z *ioDecReader) jsonReadAsisChars() (bs []byte, tok byte) {
	return z.readUntil('"', '\\')
}

func (z *ioDecReader) jsonReadUntilDblQuote() (bs []byte) {
	bs, _ = z.readUntil('"', 0)
	return
}

// ---- start/stop recording ----

func (z *ioDecReader) startRecording() {
	z.recording = true
	// always include last byte read
	if z.bufio {
		z.recc = z.rc - 1
	} else {
		z.buf = append(z.buf[:0], z.l)
	}
}

func (z *ioDecReader) stopRecording() (v []byte) {
	z.recording = false
	if z.bufio {
		v = z.buf[z.recc:z.rc]
		z.recc = 0
	} else {
		v = z.buf
		z.buf = z.buf[:0]
	}
	return
}

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
//
// Note: we do not try to convert index'ing out of bounds to an io error.
// instead, we let it bubble up to the exported Encode/Decode method
// and recover it as an io error.
//
// Every function here MUST defensively check bounds either explicitly
// or via a bounds check.
//
// see panicValToErr(...) function in helper.go.
type bytesDecReader struct {
	b  []byte // data
	c  uint   // cursor
	r  uint   // recording cursor
	xb []byte // buffer for readxb
}

func (z *bytesDecReader) resetIO(r io.Reader, bufsize int, maxInitLen int, blist *bytesFreeList) {
	halt.errorStr("resetIO unsupported by bytesDecReader")
}

func (z *bytesDecReader) resetBytes(in []byte) {
	// it's ok to resize a nil slice, so long as it's not past 0
	z.b = in[:len(in):len(in)] // reslicing must not go past capacity
	z.c = 0
}

func (z *bytesDecReader) numread() uint {
	return z.c
}

// Note: slicing from a non-constant start position is more expensive,
// as more computation is required to decipher the pointer start position.
// However, we do it only once, and it's better than reslicing both z.b and return value.

func (z *bytesDecReader) readx(n uint) (bs []byte) {
	bs = z.b[z.c : z.c+n]
	z.c += n
	return
}

func (z *bytesDecReader) skip(n uint) {
	if z.c+n > uint(cap(z.b)) {
		halt.error(&outOfBoundsError{uint(cap(z.b)), z.c + n})
	}
	z.c += n
}

func (z *bytesDecReader) readxb(n uint) (out []byte, usingBuf bool) {
	return z.readx(n), false
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(uint(len(bs))))
}

func (z *bytesDecReader) readn1() (v uint8) {
	v = z.b[z.c]
	z.c++
	return
}

func (z *bytesDecReader) readn2() (bs [2]byte) {
	bs = [2]byte(z.b[z.c:])
	z.c += 2
	return
}

func (z *bytesDecReader) readn3() (bs [3]byte) {
	bs = [3]byte(z.b[z.c:])
	z.c += 3
	return
}

func (z *bytesDecReader) readn4() (bs [4]byte) {
	bs = [4]byte(z.b[z.c:])
	z.c += 4
	return
}

func (z *bytesDecReader) readn8() (bs [8]byte) {
	bs = [8]byte(z.b[z.c:])
	z.c += 8
	return
}

func (z *bytesDecReader) jsonReadNum() (bs []byte, token byte) {
	start := z.c - 1 // include last byte
	i := start
LOOP:
	// gracefully handle end of slice (~= EOF)
	if i < uint(len(z.b)) {
		if isNumberChar(z.b[i]) {
			i++
			goto LOOP
		}
		token = z.b[i]
	}
	z.c = i + 1
	bs = z.b[start:i] // byteSliceOf(z.b, start, i)
	return
}

func (z *bytesDecReader) jsonReadAsisChars() (bs []byte, token byte) {
	i := z.c
LOOP:
	token = z.b[i]
	i++
	if token == '"' || token == '\\' {
		// z.c, i = i, z.c
		// return byteSliceOf(z.b, i, z.c-1), token
		bs = z.b[z.c : i-1]
		z.c = i
		return
		// return z.b[i : z.c-1], token
	}
	goto LOOP
}

func (z *bytesDecReader) skipWhitespace() (token byte) {
	i := z.c
LOOP:
	// setting token before check reduces inlining cost,
	// making containerNext inlineable
	token = z.b[i]
	if !isWhitespaceChar(token) {
		z.c = i + 1
		return
	}
	i++
	goto LOOP
}

func (z *bytesDecReader) jsonReadUntilDblQuote() (out []byte) {
	i := z.c
LOOP:
	if z.b[i] == '"' {
		out = z.b[z.c:i] // byteSliceOf(z.b, z.c, i)
		z.c = i + 1
		return
	}
	i++
	goto LOOP
}

func (z *bytesDecReader) startRecording() {
	z.r = z.c - 1
}

func (z *bytesDecReader) stopRecording() (v []byte) {
	v = z.b[z.r:z.c]
	z.r = 0
	return
}

type devNullReader struct{}

func (devNullReader) Read(p []byte) (int, error) { return 0, io.EOF }
func (devNullReader) Close() error               { return nil }
func (devNullReader) ReadByte() (byte, error)    { return 0, io.EOF }
func (devNullReader) UnreadByte() error          { return io.EOF }

// MARKER: readn{1,2,3,4,8} should throw an out of bounds error if past length.
// MARKER: readn1: explicitly ensure bounds check is done
// MARKER: readn{2,3,4,8}: ensure you slice z.b completely so we get bounds error if past end.
