//go:build !notmono && !codec.notmono  && (notfastpath || codec.notfastpath)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

type fastpathEMsgpackBytes struct {
	rt    reflect.Type
	encfn func(*encoderMsgpackBytes, *encFnInfo, reflect.Value)
}
type fastpathDMsgpackBytes struct {
	rt    reflect.Type
	decfn func(*decoderMsgpackBytes, *decFnInfo, reflect.Value)
}
type fastpathEsMsgpackBytes [0]fastpathEMsgpackBytes
type fastpathDsMsgpackBytes [0]fastpathDMsgpackBytes

func (helperEncDriverMsgpackBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderMsgpackBytes) bool {
	return false
}
func (helperDecDriverMsgpackBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderMsgpackBytes) bool {
	return false
}

func (helperEncDriverMsgpackBytes) fastpathEList() (v *fastpathEsMsgpackBytes) { return }
func (helperDecDriverMsgpackBytes) fastpathDList() (v *fastpathDsMsgpackBytes) { return }

type fastpathEMsgpackIO struct {
	rt    reflect.Type
	encfn func(*encoderMsgpackIO, *encFnInfo, reflect.Value)
}
type fastpathDMsgpackIO struct {
	rt    reflect.Type
	decfn func(*decoderMsgpackIO, *decFnInfo, reflect.Value)
}
type fastpathEsMsgpackIO [0]fastpathEMsgpackIO
type fastpathDsMsgpackIO [0]fastpathDMsgpackIO

func (helperEncDriverMsgpackIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderMsgpackIO) bool {
	return false
}
func (helperDecDriverMsgpackIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderMsgpackIO) bool {
	return false
}

func (helperEncDriverMsgpackIO) fastpathEList() (v *fastpathEsMsgpackIO) { return }
func (helperDecDriverMsgpackIO) fastpathDList() (v *fastpathDsMsgpackIO) { return }
