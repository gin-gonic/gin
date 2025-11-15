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

 package optdec

 import (
	 "encoding/json"
	 "errors"
	 "reflect"
	 "strconv"
 
	 "github.com/bytedance/sonic/internal/rt"
 )

 /** JIT Error Helpers **/
 
 var stackOverflow = &json.UnsupportedValueError{
	 Str:   "Value nesting too deep",
	 Value: reflect.ValueOf("..."),
 }
 
 func error_type(vt *rt.GoType) error {
	 return &json.UnmarshalTypeError{Type: vt.Pack()}
 }
 
 func error_mismatch(node Node, ctx *context, typ reflect.Type) error {
	 return MismatchTypeError{
		 Pos:  node.Position(),
		 Src:  ctx.Parser.Json,
		 Type: typ,
	 }
 }
 
 func newUnmatched(pos int, vt *rt.GoType) error {
	 return MismatchTypeError{
		Pos:  pos,
		Src:  "",
		Type: vt.Pack(),
	 }
 }

 func error_field(name string) error {
	 return errors.New("json: unknown field " + strconv.Quote(name))
 }
 
 func error_value(value string, vtype reflect.Type) error {
	 return &json.UnmarshalTypeError{
		 Type:  vtype,
		 Value: value,
	 }
 }
 
 func error_syntax(pos int, src string, msg string) error {
	 return SyntaxError{
		 Pos: pos,
		 Src: src,
		 Msg: msg,
	 }
 }

 func error_unsuppoted(typ *rt.GoType) error {
	return &json.UnsupportedTypeError{
		Type: typ.Pack(),
	}
}
 