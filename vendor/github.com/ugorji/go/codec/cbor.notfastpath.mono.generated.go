//go:build !notmono && !codec.notmono  && (notfastpath || codec.notfastpath)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

type fastpathECborBytes struct {
	rt    reflect.Type
	encfn func(*encoderCborBytes, *encFnInfo, reflect.Value)
}
type fastpathDCborBytes struct {
	rt    reflect.Type
	decfn func(*decoderCborBytes, *decFnInfo, reflect.Value)
}
type fastpathEsCborBytes [0]fastpathECborBytes
type fastpathDsCborBytes [0]fastpathDCborBytes

func (helperEncDriverCborBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborBytes) bool {
	return false
}
func (helperDecDriverCborBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborBytes) bool {
	return false
}

func (helperEncDriverCborBytes) fastpathEList() (v *fastpathEsCborBytes) { return }
func (helperDecDriverCborBytes) fastpathDList() (v *fastpathDsCborBytes) { return }

type fastpathECborIO struct {
	rt    reflect.Type
	encfn func(*encoderCborIO, *encFnInfo, reflect.Value)
}
type fastpathDCborIO struct {
	rt    reflect.Type
	decfn func(*decoderCborIO, *decFnInfo, reflect.Value)
}
type fastpathEsCborIO [0]fastpathECborIO
type fastpathDsCborIO [0]fastpathDCborIO

func (helperEncDriverCborIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborIO) bool {
	return false
}
func (helperDecDriverCborIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborIO) bool {
	return false
}

func (helperEncDriverCborIO) fastpathEList() (v *fastpathEsCborIO) { return }
func (helperDecDriverCborIO) fastpathDList() (v *fastpathDsCborIO) { return }
