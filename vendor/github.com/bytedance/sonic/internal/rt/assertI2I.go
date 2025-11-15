// +build go1.17

/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rt

import (
    _ `unsafe`
)

func AssertI2I(t *GoType, i GoIface) (r GoIface) {
    inter := IfaceType(t)
	tab := i.Itab
	if tab == nil {
		return
	}
	if (*GoInterfaceType)(tab.it) != inter {
		tab = GetItab(inter, tab.Vt, true)
		if tab == nil {
			return
		}
	}
	r.Itab = tab
	r.Value = i.Value
	return
}


