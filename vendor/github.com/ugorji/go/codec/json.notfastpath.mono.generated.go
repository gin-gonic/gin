//go:build !notmono && !codec.notmono  && (notfastpath || codec.notfastpath)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

type fastpathEJsonBytes struct {
	rt    reflect.Type
	encfn func(*encoderJsonBytes, *encFnInfo, reflect.Value)
}
type fastpathDJsonBytes struct {
	rt    reflect.Type
	decfn func(*decoderJsonBytes, *decFnInfo, reflect.Value)
}
type fastpathEsJsonBytes [0]fastpathEJsonBytes
type fastpathDsJsonBytes [0]fastpathDJsonBytes

func (helperEncDriverJsonBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderJsonBytes) bool {
	return false
}
func (helperDecDriverJsonBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderJsonBytes) bool {
	return false
}

func (helperEncDriverJsonBytes) fastpathEList() (v *fastpathEsJsonBytes) { return }
func (helperDecDriverJsonBytes) fastpathDList() (v *fastpathDsJsonBytes) { return }

type fastpathEJsonIO struct {
	rt    reflect.Type
	encfn func(*encoderJsonIO, *encFnInfo, reflect.Value)
}
type fastpathDJsonIO struct {
	rt    reflect.Type
	decfn func(*decoderJsonIO, *decFnInfo, reflect.Value)
}
type fastpathEsJsonIO [0]fastpathEJsonIO
type fastpathDsJsonIO [0]fastpathDJsonIO

func (helperEncDriverJsonIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderJsonIO) bool {
	return false
}
func (helperDecDriverJsonIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderJsonIO) bool {
	return false
}

func (helperEncDriverJsonIO) fastpathEList() (v *fastpathEsJsonIO) { return }
func (helperDecDriverJsonIO) fastpathDList() (v *fastpathDsJsonIO) { return }
