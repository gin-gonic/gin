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

// SyntaxError represents a syntax error in the expression.
type SyntaxError struct {
	Pos    int
	Reason string
}

func newSyntaxError(pos int, reason string) *SyntaxError {
	return &SyntaxError{
		Pos:    pos,
		Reason: reason,
	}
}

func (self *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error at position %d: %s", self.Pos, self.Reason)
}

// RuntimeError is an error which would occur at run time.
type RuntimeError struct {
	Reason string
}

func newRuntimeError(reason string) *RuntimeError {
	return &RuntimeError{
		Reason: reason,
	}
}

func (self *RuntimeError) Error() string {
	return "Runtime error: " + self.Reason
}
