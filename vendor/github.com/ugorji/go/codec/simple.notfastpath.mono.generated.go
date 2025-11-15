//go:build !notmono && !codec.notmono  && (notfastpath || codec.notfastpath)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

type fastpathESimpleBytes struct {
	rt    reflect.Type
	encfn func(*encoderSimpleBytes, *encFnInfo, reflect.Value)
}
type fastpathDSimpleBytes struct {
	rt    reflect.Type
	decfn func(*decoderSimpleBytes, *decFnInfo, reflect.Value)
}
type fastpathEsSimpleBytes [0]fastpathESimpleBytes
type fastpathDsSimpleBytes [0]fastpathDSimpleBytes

func (helperEncDriverSimpleBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderSimpleBytes) bool {
	return false
}
func (helperDecDriverSimpleBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderSimpleBytes) bool {
	return false
}

func (helperEncDriverSimpleBytes) fastpathEList() (v *fastpathEsSimpleBytes) { return }
func (helperDecDriverSimpleBytes) fastpathDList() (v *fastpathDsSimpleBytes) { return }

type fastpathESimpleIO struct {
	rt    reflect.Type
	encfn func(*encoderSimpleIO, *encFnInfo, reflect.Value)
}
type fastpathDSimpleIO struct {
	rt    reflect.Type
	decfn func(*decoderSimpleIO, *decFnInfo, reflect.Value)
}
type fastpathEsSimpleIO [0]fastpathESimpleIO
type fastpathDsSimpleIO [0]fastpathDSimpleIO

func (helperEncDriverSimpleIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderSimpleIO) bool {
	return false
}
func (helperDecDriverSimpleIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderSimpleIO) bool {
	return false
}

func (helperEncDriverSimpleIO) fastpathEList() (v *fastpathEsSimpleIO) { return }
func (helperDecDriverSimpleIO) fastpathDList() (v *fastpathDsSimpleIO) { return }
