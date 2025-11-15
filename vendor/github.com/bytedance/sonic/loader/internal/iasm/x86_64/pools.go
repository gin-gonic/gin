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

package x86_64

// CreateLabel creates a new Label, it may allocate a new one or grab one from a pool.
func CreateLabel(name string) *Label {
	p := new(Label)

	/* initialize the label */
	p.refs = 1
	p.Name = name
	return p
}

func newProgram(arch *Arch) *Program {
	p := new(Program)

	/* initialize the program */
	p.arch = arch
	return p
}

func newInstruction(name string, argc int, argv Operands) *Instruction {
	p := new(Instruction)

	/* initialize the instruction */
	p.name = name
	p.argc = argc
	p.argv = argv
	return p
}

// CreateMemoryOperand creates a new MemoryOperand, it may allocate a new one or grab one from a pool.
func CreateMemoryOperand() *MemoryOperand {
	p := new(MemoryOperand)

	/* initialize the memory operand */
	p.refs = 1
	return p
}
