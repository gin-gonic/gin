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

import (
	"strconv"
	"unicode"
	"unsafe"
)

type _TokenKind uint8

const (
	_T_end _TokenKind = iota + 1
	_T_int
	_T_punc
	_T_name
)

const (
	_OP2 = 0x80
	_POW = _OP2 | '*'
	_SHL = _OP2 | '<'
	_SHR = _OP2 | '>'
)

type _Slice struct {
	p unsafe.Pointer
	n int
	c int
}

type _Token struct {
	pos int
	ptr *rune
	u64 uint64
	tag _TokenKind
}

func (self _Token) str() (v string) {
	return string(self.rbuf())
}

func (self _Token) rbuf() (v []rune) {
	(*_Slice)(unsafe.Pointer(&v)).c = int(self.u64)
	(*_Slice)(unsafe.Pointer(&v)).n = int(self.u64)
	(*_Slice)(unsafe.Pointer(&v)).p = unsafe.Pointer(self.ptr)
	return
}

func tokenEnd(p int) _Token {
	return _Token{
		pos: p,
		tag: _T_end,
	}
}

func tokenInt(p int, v uint64) _Token {
	return _Token{
		pos: p,
		u64: v,
		tag: _T_int,
	}
}

func tokenPunc(p int, v rune) _Token {
	return _Token{
		pos: p,
		tag: _T_punc,
		u64: uint64(v),
	}
}

func tokenName(p int, v []rune) _Token {
	return _Token{
		pos: p,
		ptr: &v[0],
		tag: _T_name,
		u64: uint64(len(v)),
	}
}

// Repository represents a repository of Term's.
type Repository interface {
	Get(name string) (Term, error)
}

// Parser parses an expression string to it's AST representation.
type Parser struct {
	pos int
	src []rune
}

var binaryOps = [...]func(*Expr, *Expr) *Expr{
	'+':  (*Expr).Add,
	'-':  (*Expr).Sub,
	'*':  (*Expr).Mul,
	'/':  (*Expr).Div,
	'%':  (*Expr).Mod,
	'&':  (*Expr).And,
	'^':  (*Expr).Xor,
	'|':  (*Expr).Or,
	_SHL: (*Expr).Shl,
	_SHR: (*Expr).Shr,
	_POW: (*Expr).Pow,
}

var precedence = [...]map[int]bool{
	{_SHL: true, _SHR: true},
	{'|': true},
	{'^': true},
	{'&': true},
	{'+': true, '-': true},
	{'*': true, '/': true, '%': true},
	{_POW: true},
}

func (self *Parser) ch() rune {
	return self.src[self.pos]
}

func (self *Parser) eof() bool {
	return self.pos >= len(self.src)
}

func (self *Parser) rch() (v rune) {
	v, self.pos = self.src[self.pos], self.pos+1
	return
}

func (self *Parser) hex(ss []rune) bool {
	if len(ss) == 1 && ss[0] == '0' {
		return unicode.ToLower(self.ch()) == 'x'
	} else if len(ss) <= 1 || unicode.ToLower(ss[1]) != 'x' {
		return unicode.IsDigit(self.ch())
	} else {
		return ishexdigit(self.ch())
	}
}

func (self *Parser) int(p int, ss []rune) (_Token, error) {
	var err error
	var val uint64

	/* find all the digits */
	for !self.eof() && self.hex(ss) {
		ss = append(ss, self.rch())
	}

	/* parse the value */
	if val, err = strconv.ParseUint(string(ss), 0, 64); err != nil {
		return _Token{}, err
	} else {
		return tokenInt(p, val), nil
	}
}

func (self *Parser) name(p int, ss []rune) _Token {
	for !self.eof() && isident(self.ch()) {
		ss = append(ss, self.rch())
	}
	return tokenName(p, ss)
}

func (self *Parser) read(p int, ch rune) (_Token, error) {
	if isdigit(ch) {
		return self.int(p, []rune{ch})
	} else if isident0(ch) {
		return self.name(p, []rune{ch}), nil
	} else if isop2ch(ch) && !self.eof() && self.ch() == ch {
		return tokenPunc(p, _OP2|self.rch()), nil
	} else if isop1ch(ch) {
		return tokenPunc(p, ch), nil
	} else {
		return _Token{}, newSyntaxError(self.pos, "invalid character "+strconv.QuoteRuneToASCII(ch))
	}
}

func (self *Parser) next() (_Token, error) {
	for {
		var p int
		var c rune

		/* check for EOF */
		if self.eof() {
			return tokenEnd(self.pos), nil
		}

		/* read the next char */
		p = self.pos
		c = self.rch()

		/* parse the token if not a space */
		if !unicode.IsSpace(c) {
			return self.read(p, c)
		}
	}
}

func (self *Parser) grab(tk _Token, repo Repository) (*Expr, error) {
	if repo == nil {
		return nil, newSyntaxError(tk.pos, "unresolved symbol: "+tk.str())
	} else if term, err := repo.Get(tk.str()); err != nil {
		return nil, err
	} else {
		return Ref(term), nil
	}
}

func (self *Parser) nest(nest int, repo Repository) (*Expr, error) {
	var err error
	var ret *Expr
	var ntk _Token

	/* evaluate the nested expression */
	if ret, err = self.expr(0, nest+1, repo); err != nil {
		return nil, err
	}

	/* must follows with a ')' */
	if ntk, err = self.next(); err != nil {
		return nil, err
	} else if ntk.tag != _T_punc || ntk.u64 != ')' {
		return nil, newSyntaxError(ntk.pos, "')' expected")
	} else {
		return ret, nil
	}
}

func (self *Parser) unit(nest int, repo Repository) (*Expr, error) {
	if tk, err := self.next(); err != nil {
		return nil, err
	} else if tk.tag == _T_int {
		return Int(int64(tk.u64)), nil
	} else if tk.tag == _T_name {
		return self.grab(tk, repo)
	} else if tk.tag == _T_punc && tk.u64 == '(' {
		return self.nest(nest, repo)
	} else if tk.tag == _T_punc && tk.u64 == '+' {
		return self.unit(nest, repo)
	} else if tk.tag == _T_punc && tk.u64 == '-' {
		return neg2(self.unit(nest, repo))
	} else if tk.tag == _T_punc && tk.u64 == '~' {
		return not2(self.unit(nest, repo))
	} else {
		return nil, newSyntaxError(tk.pos, "integer, unary operator or nested expression expected")
	}
}

func (self *Parser) term(prec int, nest int, repo Repository) (*Expr, error) {
	var err error
	var val *Expr

	/* parse the LHS operand */
	if val, err = self.expr(prec+1, nest, repo); err != nil {
		return nil, err
	}

	/* parse all the operators of the same precedence */
	for {
		var op int
		var rv *Expr
		var tk _Token

		/* peek the next token */
		pp := self.pos
		tk, err = self.next()

		/* check for errors */
		if err != nil {
			return nil, err
		}

		/* encountered EOF */
		if tk.tag == _T_end {
			return val, nil
		}

		/* must be an operator */
		if tk.tag != _T_punc {
			return nil, newSyntaxError(tk.pos, "operators expected")
		}

		/* check for the operator precedence */
		if op = int(tk.u64); !precedence[prec][op] {
			self.pos = pp
			return val, nil
		}

		/* evaluate the RHS operand, and combine the value */
		if rv, err = self.expr(prec+1, nest, repo); err != nil {
			return nil, err
		} else {
			val = binaryOps[op](val, rv)
		}
	}
}

func (self *Parser) expr(prec int, nest int, repo Repository) (*Expr, error) {
	if prec >= len(precedence) {
		return self.unit(nest, repo)
	} else {
		return self.term(prec, nest, repo)
	}
}

// Parse parses the expression, and returns it's AST tree.
func (self *Parser) Parse(repo Repository) (*Expr, error) {
	return self.expr(0, 0, repo)
}

// SetSource resets the expression parser and sets the expression source.
func (self *Parser) SetSource(src string) *Parser {
	self.pos = 0
	self.src = []rune(src)
	return self
}
