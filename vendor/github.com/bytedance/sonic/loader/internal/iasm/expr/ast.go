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
	"fmt"
)

// Type is tyep expression type.
type Type int

const (
	// CONST indicates that the expression is a constant.
	CONST Type = iota

	// TERM indicates that the expression is a Term reference.
	TERM

	// EXPR indicates that the expression is a unary or binary expression.
	EXPR
)

var typeNames = map[Type]string{
	EXPR:  "Expr",
	TERM:  "Term",
	CONST: "Const",
}

// String returns the string representation of a Type.
func (self Type) String() string {
	if v, ok := typeNames[self]; ok {
		return v
	} else {
		return fmt.Sprintf("expr.Type(%d)", self)
	}
}

// Operator represents an operation to perform when Type is EXPR.
type Operator uint8

const (
	// ADD performs "Add Expr.Left and Expr.Right".
	ADD Operator = iota

	// SUB performs "Subtract Expr.Left by Expr.Right".
	SUB

	// MUL performs "Multiply Expr.Left by Expr.Right".
	MUL

	// DIV performs "Divide Expr.Left by Expr.Right".
	DIV

	// MOD performs "Modulo Expr.Left by Expr.Right".
	MOD

	// AND performs "Bitwise AND Expr.Left and Expr.Right".
	AND

	// OR performs "Bitwise OR Expr.Left and Expr.Right".
	OR

	// XOR performs "Bitwise XOR Expr.Left and Expr.Right".
	XOR

	// SHL performs "Bitwise Shift Expr.Left to the Left by Expr.Right Bits".
	SHL

	// SHR performs "Bitwise Shift Expr.Left to the Right by Expr.Right Bits".
	SHR

	// POW performs "Raise Expr.Left to the power of Expr.Right"
	POW

	// NOT performs "Bitwise Invert Expr.Left".
	NOT

	// NEG performs "Negate Expr.Left".
	NEG
)

var operatorNames = map[Operator]string{
	ADD: "Add",
	SUB: "Subtract",
	MUL: "Multiply",
	DIV: "Divide",
	MOD: "Modulo",
	AND: "And",
	OR:  "Or",
	XOR: "ExclusiveOr",
	SHL: "ShiftLeft",
	SHR: "ShiftRight",
	POW: "Power",
	NOT: "Invert",
	NEG: "Negate",
}

// String returns the string representation of a Type.
func (self Operator) String() string {
	if v, ok := operatorNames[self]; ok {
		return v
	} else {
		return fmt.Sprintf("expr.Operator(%d)", self)
	}
}

// Expr represents an expression node.
type Expr struct {
	Type  Type
	Term  Term
	Op    Operator
	Left  *Expr
	Right *Expr
	Const int64
}

// Ref creates an expression from a Term.
func Ref(t Term) (p *Expr) {
	p = newExpression()
	p.Term = t
	p.Type = TERM
	return
}

// Int creates an expression from an integer.
func Int(v int64) (p *Expr) {
	p = newExpression()
	p.Type = CONST
	p.Const = v
	return
}

func (self *Expr) clear() {
	if self.Term != nil {
		self.Term.Free()
	}
	if self.Left != nil {
		self.Left.Free()
	}
	if self.Right != nil {
		self.Right.Free()
	}
}

// Free returns the Expr into pool.
// Any operation performed after Free is undefined behavior.
func (self *Expr) Free() {
	self.clear()
	freeExpression(self)
}

// Evaluate evaluates the expression into an integer.
// It also implements the Term interface.
func (self *Expr) Evaluate() (int64, error) {
	switch self.Type {
	case EXPR:
		return self.eval()
	case TERM:
		return self.Term.Evaluate()
	case CONST:
		return self.Const, nil
	default:
		panic("invalid expression type: " + self.Type.String())
	}
}

/** Expression Combinator **/

func combine(a *Expr, op Operator, b *Expr) (r *Expr) {
	r = newExpression()
	r.Op = op
	r.Type = EXPR
	r.Left = a
	r.Right = b
	return
}

func (self *Expr) Add(v *Expr) *Expr { return combine(self, ADD, v) }
func (self *Expr) Sub(v *Expr) *Expr { return combine(self, SUB, v) }
func (self *Expr) Mul(v *Expr) *Expr { return combine(self, MUL, v) }
func (self *Expr) Div(v *Expr) *Expr { return combine(self, DIV, v) }
func (self *Expr) Mod(v *Expr) *Expr { return combine(self, MOD, v) }
func (self *Expr) And(v *Expr) *Expr { return combine(self, AND, v) }
func (self *Expr) Or(v *Expr) *Expr  { return combine(self, OR, v) }
func (self *Expr) Xor(v *Expr) *Expr { return combine(self, XOR, v) }
func (self *Expr) Shl(v *Expr) *Expr { return combine(self, SHL, v) }
func (self *Expr) Shr(v *Expr) *Expr { return combine(self, SHR, v) }
func (self *Expr) Pow(v *Expr) *Expr { return combine(self, POW, v) }
func (self *Expr) Not() *Expr        { return combine(self, NOT, nil) }
func (self *Expr) Neg() *Expr        { return combine(self, NEG, nil) }

/** Expression Evaluator **/

var binaryEvaluators = [256]func(int64, int64) (int64, error){
	ADD: func(a, b int64) (int64, error) { return a + b, nil },
	SUB: func(a, b int64) (int64, error) { return a - b, nil },
	MUL: func(a, b int64) (int64, error) { return a * b, nil },
	DIV: idiv,
	MOD: imod,
	AND: func(a, b int64) (int64, error) { return a & b, nil },
	OR:  func(a, b int64) (int64, error) { return a | b, nil },
	XOR: func(a, b int64) (int64, error) { return a ^ b, nil },
	SHL: func(a, b int64) (int64, error) { return a << b, nil },
	SHR: func(a, b int64) (int64, error) { return a >> b, nil },
	POW: ipow,
}

func (self *Expr) eval() (int64, error) {
	var lhs int64
	var rhs int64
	var err error
	var vfn func(int64, int64) (int64, error)

	/* evaluate LHS */
	if lhs, err = self.Left.Evaluate(); err != nil {
		return 0, err
	}

	/* check for unary operators */
	switch self.Op {
	case NOT:
		return self.unaryNot(lhs)
	case NEG:
		return self.unaryNeg(lhs)
	}

	/* check for operators */
	if vfn = binaryEvaluators[self.Op]; vfn == nil {
		panic("invalid operator: " + self.Op.String())
	}

	/* must be a binary expression */
	if self.Right == nil {
		panic("operator " + self.Op.String() + " is a binary operator")
	}

	/* evaluate RHS, and call the operator */
	if rhs, err = self.Right.Evaluate(); err != nil {
		return 0, err
	} else {
		return vfn(lhs, rhs)
	}
}

func (self *Expr) unaryNot(v int64) (int64, error) {
	if self.Right == nil {
		return ^v, nil
	} else {
		panic("operator Invert is an unary operator")
	}
}

func (self *Expr) unaryNeg(v int64) (int64, error) {
	if self.Right == nil {
		return -v, nil
	} else {
		panic("operator Negate is an unary operator")
	}
}
