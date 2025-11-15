package encoder

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

const uintptrSize = 4 << (^uintptr(0) >> 63)

type OpFlags uint16

const (
	AnonymousHeadFlags     OpFlags = 1 << 0
	AnonymousKeyFlags      OpFlags = 1 << 1
	IndirectFlags          OpFlags = 1 << 2
	IsTaggedKeyFlags       OpFlags = 1 << 3
	NilCheckFlags          OpFlags = 1 << 4
	AddrForMarshalerFlags  OpFlags = 1 << 5
	IsNextOpPtrTypeFlags   OpFlags = 1 << 6
	IsNilableTypeFlags     OpFlags = 1 << 7
	MarshalerContextFlags  OpFlags = 1 << 8
	NonEmptyInterfaceFlags OpFlags = 1 << 9
)

type Opcode struct {
	Op         OpType  // operation type
	Idx        uint32  // offset to access ptr
	Next       *Opcode // next opcode
	End        *Opcode // array/slice/struct/map end
	NextField  *Opcode // next struct field
	Key        string  // struct field key
	Offset     uint32  // offset size from struct header
	PtrNum     uint8   // pointer number: e.g. double pointer is 2.
	NumBitSize uint8
	Flags      OpFlags

	Type       *runtime.Type // go type
	Jmp        *CompiledCode // for recursive call
	FieldQuery *FieldQuery   // field query for Interface / MarshalJSON / MarshalText
	ElemIdx    uint32        // offset to access array/slice elem
	Length     uint32        // offset to access slice length or array length
	Indent     uint32        // indent number
	Size       uint32        // array/slice elem size
	DisplayIdx uint32        // opcode index
	DisplayKey string        // key text to display
}

func (c *Opcode) Validate() error {
	var prevIdx uint32
	for code := c; !code.IsEnd(); {
		if prevIdx != 0 {
			if code.DisplayIdx != prevIdx+1 {
				return fmt.Errorf(
					"invalid index. previous display index is %d but next is %d. dump = %s",
					prevIdx, code.DisplayIdx, c.Dump(),
				)
			}
		}
		prevIdx = code.DisplayIdx
		code = code.IterNext()
	}
	return nil
}

func (c *Opcode) IterNext() *Opcode {
	if c == nil {
		return nil
	}
	switch c.Op.CodeType() {
	case CodeArrayElem, CodeSliceElem, CodeMapKey:
		return c.End
	default:
		return c.Next
	}
}

func (c *Opcode) IsEnd() bool {
	if c == nil {
		return true
	}
	return c.Op == OpEnd || c.Op == OpInterfaceEnd || c.Op == OpRecursiveEnd
}

func (c *Opcode) MaxIdx() uint32 {
	max := uint32(0)
	for _, value := range []uint32{
		c.Idx,
		c.ElemIdx,
		c.Length,
		c.Size,
	} {
		if max < value {
			max = value
		}
	}
	return max
}

func (c *Opcode) ToHeaderType(isString bool) OpType {
	switch c.Op {
	case OpInt:
		if isString {
			return OpStructHeadIntString
		}
		return OpStructHeadInt
	case OpIntPtr:
		if isString {
			return OpStructHeadIntPtrString
		}
		return OpStructHeadIntPtr
	case OpUint:
		if isString {
			return OpStructHeadUintString
		}
		return OpStructHeadUint
	case OpUintPtr:
		if isString {
			return OpStructHeadUintPtrString
		}
		return OpStructHeadUintPtr
	case OpFloat32:
		if isString {
			return OpStructHeadFloat32String
		}
		return OpStructHeadFloat32
	case OpFloat32Ptr:
		if isString {
			return OpStructHeadFloat32PtrString
		}
		return OpStructHeadFloat32Ptr
	case OpFloat64:
		if isString {
			return OpStructHeadFloat64String
		}
		return OpStructHeadFloat64
	case OpFloat64Ptr:
		if isString {
			return OpStructHeadFloat64PtrString
		}
		return OpStructHeadFloat64Ptr
	case OpString:
		if isString {
			return OpStructHeadStringString
		}
		return OpStructHeadString
	case OpStringPtr:
		if isString {
			return OpStructHeadStringPtrString
		}
		return OpStructHeadStringPtr
	case OpNumber:
		if isString {
			return OpStructHeadNumberString
		}
		return OpStructHeadNumber
	case OpNumberPtr:
		if isString {
			return OpStructHeadNumberPtrString
		}
		return OpStructHeadNumberPtr
	case OpBool:
		if isString {
			return OpStructHeadBoolString
		}
		return OpStructHeadBool
	case OpBoolPtr:
		if isString {
			return OpStructHeadBoolPtrString
		}
		return OpStructHeadBoolPtr
	case OpBytes:
		return OpStructHeadBytes
	case OpBytesPtr:
		return OpStructHeadBytesPtr
	case OpMap:
		return OpStructHeadMap
	case OpMapPtr:
		c.Op = OpMap
		return OpStructHeadMapPtr
	case OpArray:
		return OpStructHeadArray
	case OpArrayPtr:
		c.Op = OpArray
		return OpStructHeadArrayPtr
	case OpSlice:
		return OpStructHeadSlice
	case OpSlicePtr:
		c.Op = OpSlice
		return OpStructHeadSlicePtr
	case OpMarshalJSON:
		return OpStructHeadMarshalJSON
	case OpMarshalJSONPtr:
		return OpStructHeadMarshalJSONPtr
	case OpMarshalText:
		return OpStructHeadMarshalText
	case OpMarshalTextPtr:
		return OpStructHeadMarshalTextPtr
	}
	return OpStructHead
}

func (c *Opcode) ToFieldType(isString bool) OpType {
	switch c.Op {
	case OpInt:
		if isString {
			return OpStructFieldIntString
		}
		return OpStructFieldInt
	case OpIntPtr:
		if isString {
			return OpStructFieldIntPtrString
		}
		return OpStructFieldIntPtr
	case OpUint:
		if isString {
			return OpStructFieldUintString
		}
		return OpStructFieldUint
	case OpUintPtr:
		if isString {
			return OpStructFieldUintPtrString
		}
		return OpStructFieldUintPtr
	case OpFloat32:
		if isString {
			return OpStructFieldFloat32String
		}
		return OpStructFieldFloat32
	case OpFloat32Ptr:
		if isString {
			return OpStructFieldFloat32PtrString
		}
		return OpStructFieldFloat32Ptr
	case OpFloat64:
		if isString {
			return OpStructFieldFloat64String
		}
		return OpStructFieldFloat64
	case OpFloat64Ptr:
		if isString {
			return OpStructFieldFloat64PtrString
		}
		return OpStructFieldFloat64Ptr
	case OpString:
		if isString {
			return OpStructFieldStringString
		}
		return OpStructFieldString
	case OpStringPtr:
		if isString {
			return OpStructFieldStringPtrString
		}
		return OpStructFieldStringPtr
	case OpNumber:
		if isString {
			return OpStructFieldNumberString
		}
		return OpStructFieldNumber
	case OpNumberPtr:
		if isString {
			return OpStructFieldNumberPtrString
		}
		return OpStructFieldNumberPtr
	case OpBool:
		if isString {
			return OpStructFieldBoolString
		}
		return OpStructFieldBool
	case OpBoolPtr:
		if isString {
			return OpStructFieldBoolPtrString
		}
		return OpStructFieldBoolPtr
	case OpBytes:
		return OpStructFieldBytes
	case OpBytesPtr:
		return OpStructFieldBytesPtr
	case OpMap:
		return OpStructFieldMap
	case OpMapPtr:
		c.Op = OpMap
		return OpStructFieldMapPtr
	case OpArray:
		return OpStructFieldArray
	case OpArrayPtr:
		c.Op = OpArray
		return OpStructFieldArrayPtr
	case OpSlice:
		return OpStructFieldSlice
	case OpSlicePtr:
		c.Op = OpSlice
		return OpStructFieldSlicePtr
	case OpMarshalJSON:
		return OpStructFieldMarshalJSON
	case OpMarshalJSONPtr:
		return OpStructFieldMarshalJSONPtr
	case OpMarshalText:
		return OpStructFieldMarshalText
	case OpMarshalTextPtr:
		return OpStructFieldMarshalTextPtr
	}
	return OpStructField
}

func newOpCode(ctx *compileContext, typ *runtime.Type, op OpType) *Opcode {
	return newOpCodeWithNext(ctx, typ, op, newEndOp(ctx, typ))
}

func opcodeOffset(idx int) uint32 {
	return uint32(idx) * uintptrSize
}

func getCodeAddrByIdx(head *Opcode, idx uint32) *Opcode {
	addr := uintptr(unsafe.Pointer(head)) + uintptr(idx)*unsafe.Sizeof(Opcode{})
	return *(**Opcode)(unsafe.Pointer(&addr))
}

func copyOpcode(code *Opcode) *Opcode {
	codeNum := ToEndCode(code).DisplayIdx + 1
	codeSlice := make([]Opcode, codeNum)
	head := (*Opcode)((*runtime.SliceHeader)(unsafe.Pointer(&codeSlice)).Data)
	ptr := head
	c := code
	for {
		*ptr = Opcode{
			Op:         c.Op,
			Key:        c.Key,
			PtrNum:     c.PtrNum,
			NumBitSize: c.NumBitSize,
			Flags:      c.Flags,
			Idx:        c.Idx,
			Offset:     c.Offset,
			Type:       c.Type,
			FieldQuery: c.FieldQuery,
			DisplayIdx: c.DisplayIdx,
			DisplayKey: c.DisplayKey,
			ElemIdx:    c.ElemIdx,
			Length:     c.Length,
			Size:       c.Size,
			Indent:     c.Indent,
			Jmp:        c.Jmp,
		}
		if c.End != nil {
			ptr.End = getCodeAddrByIdx(head, c.End.DisplayIdx)
		}
		if c.NextField != nil {
			ptr.NextField = getCodeAddrByIdx(head, c.NextField.DisplayIdx)
		}
		if c.Next != nil {
			ptr.Next = getCodeAddrByIdx(head, c.Next.DisplayIdx)
		}
		if c.IsEnd() {
			break
		}
		ptr = getCodeAddrByIdx(head, c.DisplayIdx+1)
		c = c.IterNext()
	}
	return head
}

func setTotalLengthToInterfaceOp(code *Opcode) {
	for c := code; !c.IsEnd(); {
		if c.Op == OpInterface || c.Op == OpInterfacePtr {
			c.Length = uint32(code.TotalLength())
		}
		c = c.IterNext()
	}
}

func ToEndCode(code *Opcode) *Opcode {
	c := code
	for !c.IsEnd() {
		c = c.IterNext()
	}
	return c
}

func copyToInterfaceOpcode(code *Opcode) *Opcode {
	copied := copyOpcode(code)
	c := copied
	c = ToEndCode(c)
	c.Idx += uintptrSize
	c.ElemIdx = c.Idx + uintptrSize
	c.Length = c.Idx + 2*uintptrSize
	c.Op = OpInterfaceEnd
	return copied
}

func newOpCodeWithNext(ctx *compileContext, typ *runtime.Type, op OpType, next *Opcode) *Opcode {
	return &Opcode{
		Op:         op,
		Idx:        opcodeOffset(ctx.ptrIndex),
		Next:       next,
		Type:       typ,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
	}
}

func newEndOp(ctx *compileContext, typ *runtime.Type) *Opcode {
	return newOpCodeWithNext(ctx, typ, OpEnd, nil)
}

func (c *Opcode) TotalLength() int {
	var idx int
	code := c
	for !code.IsEnd() {
		maxIdx := int(code.MaxIdx() / uintptrSize)
		if idx < maxIdx {
			idx = maxIdx
		}
		if code.Op == OpRecursiveEnd {
			break
		}
		code = code.IterNext()
	}
	maxIdx := int(code.MaxIdx() / uintptrSize)
	if idx < maxIdx {
		idx = maxIdx
	}
	return idx + 1
}

func (c *Opcode) dumpHead(code *Opcode) string {
	var length uint32
	if code.Op.CodeType() == CodeArrayHead {
		length = code.Length
	} else {
		length = code.Length / uintptrSize
	}
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d][elemIdx:%d][length:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
		code.ElemIdx/uintptrSize,
		length,
	)
}

func (c *Opcode) dumpMapHead(code *Opcode) string {
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
	)
}

func (c *Opcode) dumpMapEnd(code *Opcode) string {
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
	)
}

func (c *Opcode) dumpElem(code *Opcode) string {
	var length uint32
	if code.Op.CodeType() == CodeArrayElem {
		length = code.Length
	} else {
		length = code.Length / uintptrSize
	}
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d][elemIdx:%d][length:%d][size:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
		code.ElemIdx/uintptrSize,
		length,
		code.Size,
	)
}

func (c *Opcode) dumpField(code *Opcode) string {
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d][key:%s][offset:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
		code.DisplayKey,
		code.Offset,
	)
}

func (c *Opcode) dumpKey(code *Opcode) string {
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
	)
}

func (c *Opcode) dumpValue(code *Opcode) string {
	return fmt.Sprintf(
		`[%03d]%s%s ([idx:%d])`,
		code.DisplayIdx,
		strings.Repeat("-", int(code.Indent)),
		code.Op,
		code.Idx/uintptrSize,
	)
}

func (c *Opcode) Dump() string {
	codes := []string{}
	for code := c; !code.IsEnd(); {
		switch code.Op.CodeType() {
		case CodeSliceHead:
			codes = append(codes, c.dumpHead(code))
			code = code.Next
		case CodeMapHead:
			codes = append(codes, c.dumpMapHead(code))
			code = code.Next
		case CodeArrayElem, CodeSliceElem:
			codes = append(codes, c.dumpElem(code))
			code = code.End
		case CodeMapKey:
			codes = append(codes, c.dumpKey(code))
			code = code.End
		case CodeMapValue:
			codes = append(codes, c.dumpValue(code))
			code = code.Next
		case CodeMapEnd:
			codes = append(codes, c.dumpMapEnd(code))
			code = code.Next
		case CodeStructField:
			codes = append(codes, c.dumpField(code))
			code = code.Next
		case CodeStructEnd:
			codes = append(codes, c.dumpField(code))
			code = code.Next
		default:
			codes = append(codes, fmt.Sprintf(
				"[%03d]%s%s ([idx:%d])",
				code.DisplayIdx,
				strings.Repeat("-", int(code.Indent)),
				code.Op,
				code.Idx/uintptrSize,
			))
			code = code.Next
		}
	}
	return strings.Join(codes, "\n")
}

func (c *Opcode) DumpDOT() string {
	type edge struct {
		from, to *Opcode
		label    string
		weight   int
	}
	var edges []edge

	b := &bytes.Buffer{}
	fmt.Fprintf(b, "digraph \"%p\" {\n", c.Type)
	fmt.Fprintln(b, "mclimit=1.5;\nrankdir=TD;\nordering=out;\nnode[shape=box];")
	for code := c; !code.IsEnd(); {
		label := code.Op.String()
		fmt.Fprintf(b, "\"%p\" [label=%q];\n", code, label)
		if p := code.Next; p != nil {
			edges = append(edges, edge{
				from:   code,
				to:     p,
				label:  "Next",
				weight: 10,
			})
		}
		if p := code.NextField; p != nil {
			edges = append(edges, edge{
				from:   code,
				to:     p,
				label:  "NextField",
				weight: 2,
			})
		}
		if p := code.End; p != nil {
			edges = append(edges, edge{
				from:   code,
				to:     p,
				label:  "End",
				weight: 1,
			})
		}
		if p := code.Jmp; p != nil {
			edges = append(edges, edge{
				from:   code,
				to:     p.Code,
				label:  "Jmp",
				weight: 1,
			})
		}

		switch code.Op.CodeType() {
		case CodeSliceHead:
			code = code.Next
		case CodeMapHead:
			code = code.Next
		case CodeArrayElem, CodeSliceElem:
			code = code.End
		case CodeMapKey:
			code = code.End
		case CodeMapValue:
			code = code.Next
		case CodeMapEnd:
			code = code.Next
		case CodeStructField:
			code = code.Next
		case CodeStructEnd:
			code = code.Next
		default:
			code = code.Next
		}
		if code.IsEnd() {
			fmt.Fprintf(b, "\"%p\" [label=%q];\n", code, code.Op.String())
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].to.DisplayIdx < edges[j].to.DisplayIdx
	})
	for _, e := range edges {
		fmt.Fprintf(b, "\"%p\" -> \"%p\" [label=%q][weight=%d];\n", e.from, e.to, e.label, e.weight)
	}
	fmt.Fprint(b, "}")
	return b.String()
}

func newSliceHeaderCode(ctx *compileContext, typ *runtime.Type) *Opcode {
	idx := opcodeOffset(ctx.ptrIndex)
	ctx.incPtrIndex()
	elemIdx := opcodeOffset(ctx.ptrIndex)
	ctx.incPtrIndex()
	length := opcodeOffset(ctx.ptrIndex)
	return &Opcode{
		Op:         OpSlice,
		Type:       typ,
		Idx:        idx,
		DisplayIdx: ctx.opcodeIndex,
		ElemIdx:    elemIdx,
		Length:     length,
		Indent:     ctx.indent,
	}
}

func newSliceElemCode(ctx *compileContext, typ *runtime.Type, head *Opcode, size uintptr) *Opcode {
	return &Opcode{
		Op:         OpSliceElem,
		Type:       typ,
		Idx:        head.Idx,
		DisplayIdx: ctx.opcodeIndex,
		ElemIdx:    head.ElemIdx,
		Length:     head.Length,
		Indent:     ctx.indent,
		Size:       uint32(size),
	}
}

func newArrayHeaderCode(ctx *compileContext, typ *runtime.Type, alen int) *Opcode {
	idx := opcodeOffset(ctx.ptrIndex)
	ctx.incPtrIndex()
	elemIdx := opcodeOffset(ctx.ptrIndex)
	return &Opcode{
		Op:         OpArray,
		Type:       typ,
		Idx:        idx,
		DisplayIdx: ctx.opcodeIndex,
		ElemIdx:    elemIdx,
		Indent:     ctx.indent,
		Length:     uint32(alen),
	}
}

func newArrayElemCode(ctx *compileContext, typ *runtime.Type, head *Opcode, length int, size uintptr) *Opcode {
	return &Opcode{
		Op:         OpArrayElem,
		Type:       typ,
		Idx:        head.Idx,
		DisplayIdx: ctx.opcodeIndex,
		ElemIdx:    head.ElemIdx,
		Length:     uint32(length),
		Indent:     ctx.indent,
		Size:       uint32(size),
	}
}

func newMapHeaderCode(ctx *compileContext, typ *runtime.Type) *Opcode {
	idx := opcodeOffset(ctx.ptrIndex)
	ctx.incPtrIndex()
	return &Opcode{
		Op:         OpMap,
		Type:       typ,
		Idx:        idx,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
	}
}

func newMapKeyCode(ctx *compileContext, typ *runtime.Type, head *Opcode) *Opcode {
	return &Opcode{
		Op:         OpMapKey,
		Type:       typ,
		Idx:        head.Idx,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
	}
}

func newMapValueCode(ctx *compileContext, typ *runtime.Type, head *Opcode) *Opcode {
	return &Opcode{
		Op:         OpMapValue,
		Type:       typ,
		Idx:        head.Idx,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
	}
}

func newMapEndCode(ctx *compileContext, typ *runtime.Type, head *Opcode) *Opcode {
	return &Opcode{
		Op:         OpMapEnd,
		Type:       typ,
		Idx:        head.Idx,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
		Next:       newEndOp(ctx, typ),
	}
}

func newRecursiveCode(ctx *compileContext, typ *runtime.Type, jmp *CompiledCode) *Opcode {
	return &Opcode{
		Op:         OpRecursive,
		Type:       typ,
		Idx:        opcodeOffset(ctx.ptrIndex),
		Next:       newEndOp(ctx, typ),
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
		Jmp:        jmp,
	}
}
