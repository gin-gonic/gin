/**
 * Copyright 2024 ByteDance Inc.
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

package vars

import (
	"os"
	"unsafe"
)

const (
	MaxStack = 4096 // 4k states
	StackSize = unsafe.Sizeof(Stack{})
	StateSize  = int64(unsafe.Sizeof(State{}))
	StackLimit = MaxStack * StateSize
)

const (
	MAX_ILBUF  = 100000 // cutoff at 100k of IL instructions
	MAX_FIELDS = 50     // cutoff at 50 fields struct
)

var (
	DebugSyncGC   = os.Getenv("SONIC_SYNC_GC") != ""
	DebugAsyncGC  = os.Getenv("SONIC_NO_ASYNC_GC") == ""
	DebugCheckPtr = os.Getenv("SONIC_CHECK_POINTER") != ""
)

var UseVM = os.Getenv("SONIC_ENCODER_USE_VM") != ""
