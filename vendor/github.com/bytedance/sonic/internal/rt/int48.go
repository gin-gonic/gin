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

const (
    MinInt48 int64 = -(1 << 47)
    MaxInt48 int64 = +(1 << 47) - 1
)

func PackInt(v int) uint64 {
    if u := uint64(v); int64(v) < MinInt48 || int64(v) > MaxInt48 {
        panic("int48 out of range")
    } else {
        return ((u >> 63) << 47) | (u & 0x00007fffffffffff)
    }
}

func UnpackInt(v uint64) int {
    v &= 0x0000ffffffffffff
    v |= (v >> 47) * (0xffff << 48)
    return int(v)
}
