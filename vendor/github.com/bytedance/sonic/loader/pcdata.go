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
    `encoding/binary`
)

const (
    _N_PCDATA   = 4

    _PCDATA_UnsafePoint   = 0
    _PCDATA_StackMapIndex = 1
    _PCDATA_InlTreeIndex  = 2
    _PCDATA_ArgLiveIndex  = 3

    _PCDATA_INVALID_OFFSET = 0
)

const (
    // PCDATA_UnsafePoint values.
    PCDATA_UnsafePointSafe   = -1 // Safe for async preemption
    PCDATA_UnsafePointUnsafe = -2 // Unsafe for async preemption

    // PCDATA_Restart1(2) apply on a sequence of instructions, within
    // which if an async preemption happens, we should back off the PC
    // to the start of the sequence when resume.
    // We need two so we can distinguish the start/end of the sequence
    // in case that two sequences are next to each other.
    PCDATA_Restart1 = -3
    PCDATA_Restart2 = -4

    // Like PCDATA_RestartAtEntry, but back to function entry if async
    // preempted.
    PCDATA_RestartAtEntry = -5

    _PCDATA_START_VAL = -1
)

var emptyByte byte

// Pcvalue is the program count corresponding to the value Val
//   WARN: we use relative value here (to function entry)
type Pcvalue struct {
    PC  uint32 // program count relative to function entry
    Val int32  // value relative to the value in function entry
}

// Pcdata represents pc->value mapping table.
//   WARN: we use ** [Pcdata[i].PC, Pcdata[i+1].PC) **
//   as the range where the Pcdata[i].Val is effective.
type Pcdata []Pcvalue

// see https://docs.google.com/document/d/1lyPIbmsYbXnpNj57a261hgOYVpNRcgydurVQIyZOz_o/pub
func (self Pcdata) MarshalBinary() (data []byte, err error) {
    // delta value always starts from -1
    sv := int32(_PCDATA_START_VAL)
    sp := uint32(0)
    buf := make([]byte, binary.MaxVarintLen32)
    for _, v := range self {
        if v.PC < sp {
            panic("PC must be in ascending order!")
        }
        dp := uint64(v.PC - sp)
        dv := int64(v.Val - sv)
        if dv == 0 || dp == 0 {
            continue
        }
        n := binary.PutVarint(buf, dv)
        data = append(data, buf[:n]...)
        n2 := binary.PutUvarint(buf, dp)
        data = append(data, buf[:n2]...)
        sp = v.PC
        sv = v.Val
    }
    // put 0 to indicate ends
    data = append(data, 0)
    return
}
