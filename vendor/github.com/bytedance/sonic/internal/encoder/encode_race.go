//go:build race
// +build race

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

package encoder

import (
    `encoding/json`

    `github.com/bytedance/sonic/internal/rt`
)


func helpDetectDataRace(val interface{}) {
    var out []byte
    defer func() {
        if v := recover(); v != nil {
            // NOTICE: help user to locate where panic occurs
            println("panic when encoding on: ", truncate(out))
            panic(v)
        }
    }()
    out, _ = json.Marshal(val)
}

func encodeIntoCheckRace(buf *[]byte, val interface{}, opts Options) error {
	err := encodeInto(buf, val, opts)
    /* put last to make the panic from sonic will always be caught at first */
    helpDetectDataRace(val)
    return err
}

func truncate(json []byte) string {
    if len(json) <= 256 {
        return rt.Mem2Str(json)
    } else {
        return rt.Mem2Str(json[len(json)-256:])
    }
}
