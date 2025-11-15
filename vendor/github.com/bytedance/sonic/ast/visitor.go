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

package ast

import (
    `encoding/json`
    `errors`

    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/unquote`
)

// Visitor handles the callbacks during preorder traversal of a JSON AST.
//
// According to the JSON RFC8259, a JSON AST can be defined by
// the following rules without separator / whitespace tokens.
//
//  JSON-AST  = value
//  value     = false / null / true / object / array / number / string
//  object    = begin-object [ member *( member ) ] end-object
//  member    = string value
//  array     = begin-array [ value *( value ) ] end-array
//
type Visitor interface {

    // OnNull handles a JSON null value.
    OnNull() error

    // OnBool handles a JSON true / false value.
    OnBool(v bool) error

    // OnString handles a JSON string value.
    OnString(v string) error

    // OnInt64 handles a JSON number value with int64 type.
    OnInt64(v int64, n json.Number) error

    // OnFloat64 handles a JSON number value with float64 type.
    OnFloat64(v float64, n json.Number) error

    // OnObjectBegin handles the beginning of a JSON object value with a
    // suggested capacity that can be used to make your custom object container.
    //
    // After this point the visitor will receive a sequence of callbacks like
    // [string, value, string, value, ......, ObjectEnd].
    //
    // Note:
    // 1. This is a recursive definition which means the value can
    // also be a JSON object / array described by a sequence of callbacks.
    // 2. The suggested capacity will be 0 if current object is empty.
    // 3. Currently sonic use a fixed capacity for non-empty object (keep in
    // sync with ast.Node) which might not be very suitable. This may be
    // improved in future version.
    OnObjectBegin(capacity int) error

    // OnObjectKey handles a JSON object key string in member.
    OnObjectKey(key string) error

    // OnObjectEnd handles the ending of a JSON object value.
    OnObjectEnd() error

    // OnArrayBegin handles the beginning of a JSON array value with a
    // suggested capacity that can be used to make your custom array container.
    //
    // After this point the visitor will receive a sequence of callbacks like
    // [value, value, value, ......, ArrayEnd].
    //
    // Note:
    // 1. This is a recursive definition which means the value can
    // also be a JSON object / array described by a sequence of callbacks.
    // 2. The suggested capacity will be 0 if current array is empty.
    // 3. Currently sonic use a fixed capacity for non-empty array (keep in
    // sync with ast.Node) which might not be very suitable. This may be
    // improved in future version.
    OnArrayBegin(capacity int) error

    // OnArrayEnd handles the ending of a JSON array value.
    OnArrayEnd() error
}

// VisitorOptions contains all Visitor's options. The default value is an
// empty VisitorOptions{}.
type VisitorOptions struct {
    // OnlyNumber indicates parser to directly return number value without
    // conversion, then the first argument of OnInt64 / OnFloat64 will always
    // be zero.
    OnlyNumber bool
}

var defaultVisitorOptions = &VisitorOptions{}

// Preorder decodes the whole JSON string and callbacks each AST node to visitor
// during preorder traversal. Any visitor method with an error returned will
// break the traversal and the given error will be directly returned. The opts
// argument can be reused after every call.
func Preorder(str string, visitor Visitor, opts *VisitorOptions) error {
    if opts == nil {
        opts = defaultVisitorOptions
    }
    // process VisitorOptions first to guarantee that all options will be
    // constant during decoding and make options more readable.
    var (
        optDecodeNumber = !opts.OnlyNumber
    )

    tv := &traverser{
        parser: Parser{
            s:         str,
            noLazy:    true,
            skipValue: false,
        },
        visitor: visitor,
    }

    if optDecodeNumber {
        tv.parser.decodeNumber(true)
    }

    err := tv.decodeValue()

    if optDecodeNumber {
        tv.parser.decodeNumber(false)
    }
    return err
}

type traverser struct {
    parser  Parser
    visitor Visitor
}

// NOTE: keep in sync with (*Parser).Parse method.
func (self *traverser) decodeValue() error {
    switch val := self.parser.decodeValue(); val.Vt {
    case types.V_EOF:
        return types.ERR_EOF
    case types.V_NULL:
        return self.visitor.OnNull()
    case types.V_TRUE:
        return self.visitor.OnBool(true)
    case types.V_FALSE:
        return self.visitor.OnBool(false)
    case types.V_STRING:
        return self.decodeString(val.Iv, val.Ep)
    case types.V_DOUBLE:
        return self.visitor.OnFloat64(val.Dv,
            json.Number(self.parser.s[val.Ep:self.parser.p]))
    case types.V_INTEGER:
        return self.visitor.OnInt64(val.Iv,
            json.Number(self.parser.s[val.Ep:self.parser.p]))
    case types.V_ARRAY:
        return self.decodeArray()
    case types.V_OBJECT:
        return self.decodeObject()
    default:
        return types.ParsingError(-val.Vt)
    }
}

// NOTE: keep in sync with (*Parser).decodeArray method.
func (self *traverser) decodeArray() error {
    sp := self.parser.p
    ns := len(self.parser.s)

    /* allocate array space and parse every element */
    if err := self.visitor.OnArrayBegin(_DEFAULT_NODE_CAP); err != nil {
        if err == VisitOPSkip {
            // NOTICE: for user needs to skip entry object
            self.parser.p -= 1
            if _, e := self.parser.skipFast(); e != 0 {
                return e
            }
            return self.visitor.OnArrayEnd()
        }
        return err
    }

    /* check for EOF */
    self.parser.p = self.parser.lspace(sp)
    if self.parser.p >= ns {
        return types.ERR_EOF
    }

    /* check for empty array */
    if self.parser.s[self.parser.p] == ']' {
        self.parser.p++
        return self.visitor.OnArrayEnd()
    }

    for {
        /* decode the value */
        if err := self.decodeValue(); err != nil {
            return err
        }
        self.parser.p = self.parser.lspace(self.parser.p)

        /* check for EOF */
        if self.parser.p >= ns {
            return types.ERR_EOF
        }

        /* check for the next character */
        switch self.parser.s[self.parser.p] {
        case ',':
            self.parser.p++
        case ']':
            self.parser.p++
            return self.visitor.OnArrayEnd()
        default:
            return types.ERR_INVALID_CHAR
        }
    }
}

// NOTE: keep in sync with (*Parser).decodeObject method.
func (self *traverser) decodeObject() error {
    sp := self.parser.p
    ns := len(self.parser.s)

    /* allocate object space and decode each pair */
    if err := self.visitor.OnObjectBegin(_DEFAULT_NODE_CAP); err != nil {
        if err == VisitOPSkip {
            // NOTICE: for user needs to skip entry object
            self.parser.p -= 1
            if _, e := self.parser.skipFast(); e != 0 {
                return e
            }
            return self.visitor.OnObjectEnd()
        }
        return err
    }

    /* check for EOF */
    self.parser.p = self.parser.lspace(sp)
    if self.parser.p >= ns {
        return types.ERR_EOF
    }

    /* check for empty object */
    if self.parser.s[self.parser.p] == '}' {
        self.parser.p++
        return self.visitor.OnObjectEnd()
    }

    for {
        var njs types.JsonState
        var err types.ParsingError

        /* decode the key */
        if njs = self.parser.decodeValue(); njs.Vt != types.V_STRING {
            return types.ERR_INVALID_CHAR
        }

        /* extract the key */
        idx := self.parser.p - 1
        key := self.parser.s[njs.Iv:idx]

        /* check for escape sequence */
        if njs.Ep != -1 {
            if key, err = unquote.String(key); err != 0 {
                return err
            }
        }

        if err := self.visitor.OnObjectKey(key); err != nil {
            return err
        }

        /* expect a ':' delimiter */
        if err = self.parser.delim(); err != 0 {
            return err
        }

        /* decode the value */
        if err := self.decodeValue(); err != nil {
            return err
        }

        self.parser.p = self.parser.lspace(self.parser.p)

        /* check for EOF */
        if self.parser.p >= ns {
            return types.ERR_EOF
        }

        /* check for the next character */
        switch self.parser.s[self.parser.p] {
        case ',':
            self.parser.p++
        case '}':
            self.parser.p++
            return self.visitor.OnObjectEnd()
        default:
            return types.ERR_INVALID_CHAR
        }
    }
}

// NOTE: keep in sync with (*Parser).decodeString method.
func (self *traverser) decodeString(iv int64, ep int) error {
    p := self.parser.p - 1
    s := self.parser.s[iv:p]

    /* fast path: no escape sequence */
    if ep == -1 {
        return self.visitor.OnString(s)
    }

    /* unquote the string */
    out, err := unquote.String(s)
    if err != 0 {
        return err
    }
    return self.visitor.OnString(out)
}

// If visitor return this error on `OnObjectBegin()` or `OnArrayBegin()`,
// the traverser will skip entry object or array
var VisitOPSkip = errors.New("")
