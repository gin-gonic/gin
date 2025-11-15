//
// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package expr

var op1ch = [...]bool{
	'+': true,
	'-': true,
	'*': true,
	'/': true,
	'%': true,
	'&': true,
	'|': true,
	'^': true,
	'~': true,
	'(': true,
	')': true,
}

var op2ch = [...]bool{
	'*': true,
	'<': true,
	'>': true,
}

func neg2(v *Expr, err error) (*Expr, error) {
	if err != nil {
		return nil, err
	} else {
		return v.Neg(), nil
	}
}

func not2(v *Expr, err error) (*Expr, error) {
	if err != nil {
		return nil, err
	} else {
		return v.Not(), nil
	}
}

func isop1ch(ch rune) bool {
	return ch >= 0 && int(ch) < len(op1ch) && op1ch[ch]
}

func isop2ch(ch rune) bool {
	return ch >= 0 && int(ch) < len(op2ch) && op2ch[ch]
}

func isdigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isident(ch rune) bool {
	return isdigit(ch) || isident0(ch)
}

func isident0(ch rune) bool {
	return (ch == '_') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func ishexdigit(ch rune) bool {
	return isdigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
