//go:build notfastpath || (codec.notfastpath && (notmono || codec.notmono))

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "reflect"

// type fastpathT struct{}
type fastpathE[T encDriver] struct {
	rt    reflect.Type
	encfn func(*encoder[T], *encFnInfo, reflect.Value)
}
type fastpathD[T decDriver] struct {
	rt    reflect.Type
	decfn func(*decoder[T], *decFnInfo, reflect.Value)
}
type fastpathEs[T encDriver] [0]fastpathE[T]
type fastpathDs[T decDriver] [0]fastpathD[T]

func (helperEncDriver[T]) fastpathEncodeTypeSwitch(iv interface{}, e *encoder[T]) bool { return false }
func (helperDecDriver[T]) fastpathDecodeTypeSwitch(iv interface{}, d *decoder[T]) bool { return false }

func (helperEncDriver[T]) fastpathEList() (v *fastpathEs[T]) { return }
func (helperDecDriver[T]) fastpathDList() (v *fastpathDs[T]) { return }
