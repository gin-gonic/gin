//go:build !notmono && !codec.notmono  && (notfastpath || codec.notfastpath)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

type fastpathEBincBytes struct {
	rt    reflect.Type
	encfn func(*encoderBincBytes, *encFnInfo, reflect.Value)
}
type fastpathDBincBytes struct {
	rt    reflect.Type
	decfn func(*decoderBincBytes, *decFnInfo, reflect.Value)
}
type fastpathEsBincBytes [0]fastpathEBincBytes
type fastpathDsBincBytes [0]fastpathDBincBytes

func (helperEncDriverBincBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincBytes) bool {
	return false
}
func (helperDecDriverBincBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincBytes) bool {
	return false
}

func (helperEncDriverBincBytes) fastpathEList() (v *fastpathEsBincBytes) { return }
func (helperDecDriverBincBytes) fastpathDList() (v *fastpathDsBincBytes) { return }

type fastpathEBincIO struct {
	rt    reflect.Type
	encfn func(*encoderBincIO, *encFnInfo, reflect.Value)
}
type fastpathDBincIO struct {
	rt    reflect.Type
	decfn func(*decoderBincIO, *decFnInfo, reflect.Value)
}
type fastpathEsBincIO [0]fastpathEBincIO
type fastpathDsBincIO [0]fastpathDBincIO

func (helperEncDriverBincIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincIO) bool {
	return false
}
func (helperDecDriverBincIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincIO) bool {
	return false
}

func (helperEncDriverBincIO) fastpathEList() (v *fastpathEsBincIO) { return }
func (helperDecDriverBincIO) fastpathDList() (v *fastpathDsBincIO) { return }
