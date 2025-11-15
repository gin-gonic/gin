// +build !amd64,!arm64 go1.26 !go1.17 arm64,!go1.20

/*
* Copyright 2022 ByteDance Inc.
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

package ast

import (
    `encoding/json`
    `unicode/utf8`

    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/compat`
)

func init() {
    compat.Warn("sonic/ast")
}

func quote(buf *[]byte, val string) {
    quoteString(buf, val)
}

func (self *Parser) decodeValue() (val types.JsonState) {
    e, v := decodeValue(self.s, self.p, self.dbuf == nil)
    if e < 0 {
        return v
    }
    self.p = e
    return v
}

func (self *Parser) skip() (int, types.ParsingError) {
    e, s := skipValue(self.s, self.p)
    if e < 0 {
        return self.p, types.ParsingError(-e)
    }
    self.p = e
    return s, 0
}

func (self *Parser) skipFast() (int, types.ParsingError) {
    e, s := skipValueFast(self.s, self.p)
    if e < 0 {
        return self.p, types.ParsingError(-e)
    }
    self.p = e
    return s, 0
}

func (self *Node) encodeInterface(buf *[]byte) error {
    out, err := json.Marshal(self.packAny())
    if err != nil {
        return err
    }
    *buf = append(*buf, out...)
    return nil
}

func (self *Parser) getByPath(validate bool, path ...interface{}) (int, types.ParsingError) {
    for _, p := range path {
        if idx, ok := p.(int); ok && idx >= 0 {
            if err := self.searchIndex(idx); err != 0 {
                return self.p, err
            }
        } else if key, ok := p.(string); ok {
            if err := self.searchKey(key); err != 0 {
                return self.p, err
            }
        } else {
            panic("path must be either int(>=0) or string")
        }
    }

    var start int
    var e types.ParsingError
    if validate {
        start, e = self.skip()
    } else {
        start, e = self.skipFast()
    }
    if e != 0 {
        return self.p, e
    }
    return start, 0
}

func validate_utf8(str string) bool {
    return utf8.ValidString(str)
}
