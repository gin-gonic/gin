/**
 * Copyright 2023 ByteDance Inc.
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

package loader

import (
    `unsafe`
)

// Function is a function pointer
type Function unsafe.Pointer

// Options used to load a module
type Options struct {
    // NoPreempt is used to disable async preemption for this module
    NoPreempt bool
}

// Loader is a helper used to load a module simply
type Loader struct {
    Name string // module name
    File string // file name
    Options 
}
