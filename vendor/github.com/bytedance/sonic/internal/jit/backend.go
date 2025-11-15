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

package jit

import (
    `fmt`
    `sync`
    _ `unsafe`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/loader`
    `github.com/twitchyliquid64/golang-asm/asm/arch`
    `github.com/twitchyliquid64/golang-asm/obj`
    `github.com/twitchyliquid64/golang-asm/obj/x86`
    `github.com/twitchyliquid64/golang-asm/objabi`
)

type Backend struct {
    Ctxt *obj.Link
    Arch *arch.Arch
    Head *obj.Prog
    Tail *obj.Prog
    Prog []*obj.Prog
}

var (
    _progPool sync.Pool
)

func newProg() *obj.Prog {
    if val := _progPool.Get(); val == nil {
        return new(obj.Prog)
    } else {
        return remProg(val.(*obj.Prog))
    }
}

func remProg(p *obj.Prog) *obj.Prog {
    *p = obj.Prog{}
    return p
}

func newBackend(name string) (ret *Backend) {
    ret      = new(Backend)
    ret.Arch = arch.Set(name)
    ret.Ctxt = newLinkContext(ret.Arch.LinkArch)
    ret.Arch.Init(ret.Ctxt)
    return
}

func newLinkContext(arch *obj.LinkArch) (ret *obj.Link) {
    ret          = obj.Linknew(arch)
    ret.Headtype = objabi.Hlinux
    ret.DiagFunc = diagLinkContext
    return
}

func diagLinkContext(str string, args ...interface{}) {
    rt.Throw(fmt.Sprintf(str, args...))
}

func (self *Backend) New() (ret *obj.Prog) {
    ret = newProg()
    ret.Ctxt = self.Ctxt
    self.Prog = append(self.Prog, ret)
    return
}

func (self *Backend) Append(p *obj.Prog) {
    if self.Head == nil {
        self.Head = p
        self.Tail = p
    } else {
        self.Tail.Link = p
        self.Tail = p
    }
}

func (self *Backend) Release() {
    self.Arch = nil
    self.Ctxt = nil

    /* return all the progs into pool */
    for _, p := range self.Prog {
        _progPool.Put(p)
    }

    /* clear all the references */
    self.Head = nil
    self.Tail = nil
    self.Prog = nil
}

func (self *Backend) Assemble() ([]byte, loader.Pcdata) {
    var sym obj.LSym
    var fnv obj.FuncInfo

    /* construct the function */
    sym.Func = &fnv
    fnv.Text = self.Head

    /* call the assembler */
    self.Arch.Assemble(self.Ctxt, &sym, self.New)
    pcdata := self.GetPcspTable(self.Ctxt, &sym, self.New)
    return sym.P, pcdata
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func nextPc(p *obj.Prog) uint32 {
    if p.Link != nil && p.Pc + int64(p.Isize) != p.Link.Pc {
        panic("p.PC + p.Isize != p.Link.PC")
    }
    return uint32(p.Pc + int64(p.Isize))
}

// NOTE: copied from https://github.com/twitchyliquid64/golang-asm/blob/8d7f1f783b11f9a00f5bcdfcae17f5ac8f22512e/obj/x86/obj6.go#L811.
// we add two instructions such as subq/addq %rsp, $imm to the table.
func (self *Backend) GetPcspTable(ctxt *obj.Link, cursym *obj.LSym, newprog obj.ProgAlloc) loader.Pcdata {
    pcdata := loader.Pcdata{}
    var deltasp int32
    var maxdepth int32
    foundRet := false
    p := cursym.Func.Text
	for ; p != nil; p = p.Link {
        if foundRet {
            break
        }
		switch p.As {
		default: continue
		case x86.APUSHL, x86.APUSHFL:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
			deltasp += 4
            maxdepth = max(maxdepth, deltasp)
			continue

		case x86.APUSHQ, x86.APUSHFQ:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            deltasp += 8
            maxdepth = max(maxdepth, deltasp)
			continue

		case x86.APUSHW, x86.APUSHFW:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            deltasp += 2
            maxdepth = max(maxdepth, deltasp)
			continue

		case x86.APOPL, x86.APOPFL:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            deltasp -= 4
			continue

		case x86.APOPQ, x86.APOPFQ:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            deltasp -= 8
			continue

		case x86.APOPW, x86.APOPFW:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            deltasp -= 2
			continue

		case x86.AADJSP:
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
			deltasp += int32(p.From.Offset)
            maxdepth = max(maxdepth, deltasp)
			continue

        case x86.ASUBQ:
            // subq %rsp, $imm
            if p.To.Reg == x86.REG_SP && p.To.Type == obj.TYPE_REG {
                pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
                deltasp += int32(p.From.Offset)
                maxdepth = max(maxdepth, deltasp)
            }
            continue
        case x86.AADDQ:
            // addq %rsp, $imm
            if p.To.Reg == x86.REG_SP && p.To.Type == obj.TYPE_REG {
                pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
                deltasp -= int32(p.From.Offset)
            }
            continue
		case obj.ARET:
            if deltasp != 0 {
                panic("unbalanced PUSH/POP")
            }
            pcdata = append(pcdata, loader.Pcvalue{PC: nextPc(p), Val: int32(deltasp)})
            foundRet = true
		}
	}

    // the instructions after the RET instruction
    if p != nil {
        pcdata = append(pcdata, loader.Pcvalue{PC: uint32(cursym.Size), Val: int32(maxdepth)})
    }
    
    return pcdata
}
