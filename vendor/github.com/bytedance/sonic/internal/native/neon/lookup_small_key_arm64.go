/**
 * Copyright 2024 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package neon

import (
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

//go:nosplit
func lookup_small_key(key *string, table *[]byte, lowerOff int) (ret int) {
    return __lookup_small_key(rt.NoEscape(unsafe.Pointer(key)), rt.NoEscape(unsafe.Pointer(table)), lowerOff)
}

//go:nosplit
func __lookup_small_key(key unsafe.Pointer, table unsafe.Pointer, lowerOff int) (ret int)
