package encoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

type Code interface {
	Kind() CodeKind
	ToOpcode(*compileContext) Opcodes
	Filter(*FieldQuery) Code
}

type AnonymousCode interface {
	ToAnonymousOpcode(*compileContext) Opcodes
}

type Opcodes []*Opcode

func (o Opcodes) First() *Opcode {
	if len(o) == 0 {
		return nil
	}
	return o[0]
}

func (o Opcodes) Last() *Opcode {
	if len(o) == 0 {
		return nil
	}
	return o[len(o)-1]
}

func (o Opcodes) Add(codes ...*Opcode) Opcodes {
	return append(o, codes...)
}

type CodeKind int

const (
	CodeKindInterface CodeKind = iota
	CodeKindPtr
	CodeKindInt
	CodeKindUint
	CodeKindFloat
	CodeKindString
	CodeKindBool
	CodeKindStruct
	CodeKindMap
	CodeKindSlice
	CodeKindArray
	CodeKindBytes
	CodeKindMarshalJSON
	CodeKindMarshalText
	CodeKindRecursive
)

type IntCode struct {
	typ      *runtime.Type
	bitSize  uint8
	isString bool
	isPtr    bool
}

func (c *IntCode) Kind() CodeKind {
	return CodeKindInt
}

func (c *IntCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		code = newOpCode(ctx, c.typ, OpIntPtr)
	case c.isString:
		code = newOpCode(ctx, c.typ, OpIntString)
	default:
		code = newOpCode(ctx, c.typ, OpInt)
	}
	code.NumBitSize = c.bitSize
	ctx.incIndex()
	return Opcodes{code}
}

func (c *IntCode) Filter(_ *FieldQuery) Code {
	return c
}

type UintCode struct {
	typ      *runtime.Type
	bitSize  uint8
	isString bool
	isPtr    bool
}

func (c *UintCode) Kind() CodeKind {
	return CodeKindUint
}

func (c *UintCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		code = newOpCode(ctx, c.typ, OpUintPtr)
	case c.isString:
		code = newOpCode(ctx, c.typ, OpUintString)
	default:
		code = newOpCode(ctx, c.typ, OpUint)
	}
	code.NumBitSize = c.bitSize
	ctx.incIndex()
	return Opcodes{code}
}

func (c *UintCode) Filter(_ *FieldQuery) Code {
	return c
}

type FloatCode struct {
	typ     *runtime.Type
	bitSize uint8
	isPtr   bool
}

func (c *FloatCode) Kind() CodeKind {
	return CodeKindFloat
}

func (c *FloatCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		switch c.bitSize {
		case 32:
			code = newOpCode(ctx, c.typ, OpFloat32Ptr)
		default:
			code = newOpCode(ctx, c.typ, OpFloat64Ptr)
		}
	default:
		switch c.bitSize {
		case 32:
			code = newOpCode(ctx, c.typ, OpFloat32)
		default:
			code = newOpCode(ctx, c.typ, OpFloat64)
		}
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *FloatCode) Filter(_ *FieldQuery) Code {
	return c
}

type StringCode struct {
	typ   *runtime.Type
	isPtr bool
}

func (c *StringCode) Kind() CodeKind {
	return CodeKindString
}

func (c *StringCode) ToOpcode(ctx *compileContext) Opcodes {
	isJSONNumberType := c.typ == runtime.Type2RType(jsonNumberType)
	var code *Opcode
	if c.isPtr {
		if isJSONNumberType {
			code = newOpCode(ctx, c.typ, OpNumberPtr)
		} else {
			code = newOpCode(ctx, c.typ, OpStringPtr)
		}
	} else {
		if isJSONNumberType {
			code = newOpCode(ctx, c.typ, OpNumber)
		} else {
			code = newOpCode(ctx, c.typ, OpString)
		}
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *StringCode) Filter(_ *FieldQuery) Code {
	return c
}

type BoolCode struct {
	typ   *runtime.Type
	isPtr bool
}

func (c *BoolCode) Kind() CodeKind {
	return CodeKindBool
}

func (c *BoolCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		code = newOpCode(ctx, c.typ, OpBoolPtr)
	default:
		code = newOpCode(ctx, c.typ, OpBool)
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *BoolCode) Filter(_ *FieldQuery) Code {
	return c
}

type BytesCode struct {
	typ   *runtime.Type
	isPtr bool
}

func (c *BytesCode) Kind() CodeKind {
	return CodeKindBytes
}

func (c *BytesCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		code = newOpCode(ctx, c.typ, OpBytesPtr)
	default:
		code = newOpCode(ctx, c.typ, OpBytes)
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *BytesCode) Filter(_ *FieldQuery) Code {
	return c
}

type SliceCode struct {
	typ   *runtime.Type
	value Code
}

func (c *SliceCode) Kind() CodeKind {
	return CodeKindSlice
}

func (c *SliceCode) ToOpcode(ctx *compileContext) Opcodes {
	// header => opcode => elem => end
	//             ^        |
	//             |________|
	size := c.typ.Elem().Size()
	header := newSliceHeaderCode(ctx, c.typ)
	ctx.incIndex()

	ctx.incIndent()
	codes := c.value.ToOpcode(ctx)
	ctx.decIndent()

	codes.First().Flags |= IndirectFlags
	elemCode := newSliceElemCode(ctx, c.typ.Elem(), header, size)
	ctx.incIndex()
	end := newOpCode(ctx, c.typ, OpSliceEnd)
	ctx.incIndex()
	header.End = end
	header.Next = codes.First()
	codes.Last().Next = elemCode
	elemCode.Next = codes.First()
	elemCode.End = end
	return Opcodes{header}.Add(codes...).Add(elemCode).Add(end)
}

func (c *SliceCode) Filter(_ *FieldQuery) Code {
	return c
}

type ArrayCode struct {
	typ   *runtime.Type
	value Code
}

func (c *ArrayCode) Kind() CodeKind {
	return CodeKindArray
}

func (c *ArrayCode) ToOpcode(ctx *compileContext) Opcodes {
	// header => opcode => elem => end
	//             ^        |
	//             |________|
	elem := c.typ.Elem()
	alen := c.typ.Len()
	size := elem.Size()

	header := newArrayHeaderCode(ctx, c.typ, alen)
	ctx.incIndex()

	ctx.incIndent()
	codes := c.value.ToOpcode(ctx)
	ctx.decIndent()

	codes.First().Flags |= IndirectFlags

	elemCode := newArrayElemCode(ctx, elem, header, alen, size)
	ctx.incIndex()

	end := newOpCode(ctx, c.typ, OpArrayEnd)
	ctx.incIndex()

	header.End = end
	header.Next = codes.First()
	codes.Last().Next = elemCode
	elemCode.Next = codes.First()
	elemCode.End = end

	return Opcodes{header}.Add(codes...).Add(elemCode).Add(end)
}

func (c *ArrayCode) Filter(_ *FieldQuery) Code {
	return c
}

type MapCode struct {
	typ   *runtime.Type
	key   Code
	value Code
}

func (c *MapCode) Kind() CodeKind {
	return CodeKindMap
}

func (c *MapCode) ToOpcode(ctx *compileContext) Opcodes {
	// header => code => value => code => key => code => value => code => end
	//                                     ^                       |
	//                                     |_______________________|
	header := newMapHeaderCode(ctx, c.typ)
	ctx.incIndex()

	keyCodes := c.key.ToOpcode(ctx)

	value := newMapValueCode(ctx, c.typ.Elem(), header)
	ctx.incIndex()

	ctx.incIndent()
	valueCodes := c.value.ToOpcode(ctx)
	ctx.decIndent()

	valueCodes.First().Flags |= IndirectFlags

	key := newMapKeyCode(ctx, c.typ.Key(), header)
	ctx.incIndex()

	end := newMapEndCode(ctx, c.typ, header)
	ctx.incIndex()

	header.Next = keyCodes.First()
	keyCodes.Last().Next = value
	value.Next = valueCodes.First()
	valueCodes.Last().Next = key
	key.Next = keyCodes.First()

	header.End = end
	key.End = end
	value.End = end
	return Opcodes{header}.Add(keyCodes...).Add(value).Add(valueCodes...).Add(key).Add(end)
}

func (c *MapCode) Filter(_ *FieldQuery) Code {
	return c
}

type StructCode struct {
	typ                       *runtime.Type
	fields                    []*StructFieldCode
	isPtr                     bool
	disableIndirectConversion bool
	isIndirect                bool
	isRecursive               bool
}

func (c *StructCode) Kind() CodeKind {
	return CodeKindStruct
}

func (c *StructCode) lastFieldCode(field *StructFieldCode, firstField *Opcode) *Opcode {
	if isEmbeddedStruct(field) {
		return c.lastAnonymousFieldCode(firstField)
	}
	lastField := firstField
	for lastField.NextField != nil {
		lastField = lastField.NextField
	}
	return lastField
}

func (c *StructCode) lastAnonymousFieldCode(firstField *Opcode) *Opcode {
	// firstField is special StructHead operation for anonymous structure.
	// So, StructHead's next operation is truly struct head operation.
	for firstField.Op == OpStructHead || firstField.Op == OpStructField {
		firstField = firstField.Next
	}
	lastField := firstField
	for lastField.NextField != nil {
		lastField = lastField.NextField
	}
	return lastField
}

func (c *StructCode) ToOpcode(ctx *compileContext) Opcodes {
	// header => code => structField => code => end
	//                        ^          |
	//                        |__________|
	if c.isRecursive {
		recursive := newRecursiveCode(ctx, c.typ, &CompiledCode{})
		recursive.Type = c.typ
		ctx.incIndex()
		*ctx.recursiveCodes = append(*ctx.recursiveCodes, recursive)
		return Opcodes{recursive}
	}
	codes := Opcodes{}
	var prevField *Opcode
	ctx.incIndent()
	for idx, field := range c.fields {
		isFirstField := idx == 0
		isEndField := idx == len(c.fields)-1
		fieldCodes := field.ToOpcode(ctx, isFirstField, isEndField)
		for _, code := range fieldCodes {
			if c.isIndirect {
				code.Flags |= IndirectFlags
			}
		}
		firstField := fieldCodes.First()
		if len(codes) > 0 {
			codes.Last().Next = firstField
			firstField.Idx = codes.First().Idx
		}
		if prevField != nil {
			prevField.NextField = firstField
		}
		if isEndField {
			endField := fieldCodes.Last()
			if len(codes) > 0 {
				codes.First().End = endField
			} else {
				firstField.End = endField
			}
			codes = codes.Add(fieldCodes...)
			break
		}
		prevField = c.lastFieldCode(field, firstField)
		codes = codes.Add(fieldCodes...)
	}
	if len(codes) == 0 {
		head := &Opcode{
			Op:         OpStructHead,
			Idx:        opcodeOffset(ctx.ptrIndex),
			Type:       c.typ,
			DisplayIdx: ctx.opcodeIndex,
			Indent:     ctx.indent,
		}
		ctx.incOpcodeIndex()
		end := &Opcode{
			Op:         OpStructEnd,
			Idx:        opcodeOffset(ctx.ptrIndex),
			DisplayIdx: ctx.opcodeIndex,
			Indent:     ctx.indent,
		}
		head.NextField = end
		head.Next = end
		head.End = end
		codes = codes.Add(head, end)
		ctx.incIndex()
	}
	ctx.decIndent()
	ctx.structTypeToCodes[uintptr(unsafe.Pointer(c.typ))] = codes
	return codes
}

func (c *StructCode) ToAnonymousOpcode(ctx *compileContext) Opcodes {
	// header => code => structField => code => end
	//                        ^          |
	//                        |__________|
	if c.isRecursive {
		recursive := newRecursiveCode(ctx, c.typ, &CompiledCode{})
		recursive.Type = c.typ
		ctx.incIndex()
		*ctx.recursiveCodes = append(*ctx.recursiveCodes, recursive)
		return Opcodes{recursive}
	}
	codes := Opcodes{}
	var prevField *Opcode
	for idx, field := range c.fields {
		isFirstField := idx == 0
		isEndField := idx == len(c.fields)-1
		fieldCodes := field.ToAnonymousOpcode(ctx, isFirstField, isEndField)
		for _, code := range fieldCodes {
			if c.isIndirect {
				code.Flags |= IndirectFlags
			}
		}
		firstField := fieldCodes.First()
		if len(codes) > 0 {
			codes.Last().Next = firstField
			firstField.Idx = codes.First().Idx
		}
		if prevField != nil {
			prevField.NextField = firstField
		}
		if isEndField {
			lastField := fieldCodes.Last()
			if len(codes) > 0 {
				codes.First().End = lastField
			} else {
				firstField.End = lastField
			}
		}
		prevField = firstField
		codes = codes.Add(fieldCodes...)
	}
	return codes
}

func (c *StructCode) removeFieldsByTags(tags runtime.StructTags) {
	fields := make([]*StructFieldCode, 0, len(c.fields))
	for _, field := range c.fields {
		if field.isAnonymous {
			structCode := field.getAnonymousStruct()
			if structCode != nil && !structCode.isRecursive {
				structCode.removeFieldsByTags(tags)
				if len(structCode.fields) > 0 {
					fields = append(fields, field)
				}
				continue
			}
		}
		if tags.ExistsKey(field.key) {
			continue
		}
		fields = append(fields, field)
	}
	c.fields = fields
}

func (c *StructCode) enableIndirect() {
	if c.isIndirect {
		return
	}
	c.isIndirect = true
	if len(c.fields) == 0 {
		return
	}
	structCode := c.fields[0].getStruct()
	if structCode == nil {
		return
	}
	structCode.enableIndirect()
}

func (c *StructCode) Filter(query *FieldQuery) Code {
	fieldMap := map[string]*FieldQuery{}
	for _, field := range query.Fields {
		fieldMap[field.Name] = field
	}
	fields := make([]*StructFieldCode, 0, len(c.fields))
	for _, field := range c.fields {
		query, exists := fieldMap[field.key]
		if !exists {
			continue
		}
		fieldCode := &StructFieldCode{
			typ:                field.typ,
			key:                field.key,
			tag:                field.tag,
			value:              field.value,
			offset:             field.offset,
			isAnonymous:        field.isAnonymous,
			isTaggedKey:        field.isTaggedKey,
			isNilableType:      field.isNilableType,
			isNilCheck:         field.isNilCheck,
			isAddrForMarshaler: field.isAddrForMarshaler,
			isNextOpPtrType:    field.isNextOpPtrType,
		}
		if len(query.Fields) > 0 {
			fieldCode.value = fieldCode.value.Filter(query)
		}
		fields = append(fields, fieldCode)
	}
	return &StructCode{
		typ:                       c.typ,
		fields:                    fields,
		isPtr:                     c.isPtr,
		disableIndirectConversion: c.disableIndirectConversion,
		isIndirect:                c.isIndirect,
		isRecursive:               c.isRecursive,
	}
}

type StructFieldCode struct {
	typ                *runtime.Type
	key                string
	tag                *runtime.StructTag
	value              Code
	offset             uintptr
	isAnonymous        bool
	isTaggedKey        bool
	isNilableType      bool
	isNilCheck         bool
	isAddrForMarshaler bool
	isNextOpPtrType    bool
	isMarshalerContext bool
}

func (c *StructFieldCode) getStruct() *StructCode {
	value := c.value
	ptr, ok := value.(*PtrCode)
	if ok {
		value = ptr.value
	}
	structCode, ok := value.(*StructCode)
	if ok {
		return structCode
	}
	return nil
}

func (c *StructFieldCode) getAnonymousStruct() *StructCode {
	if !c.isAnonymous {
		return nil
	}
	return c.getStruct()
}

func optimizeStructHeader(code *Opcode, tag *runtime.StructTag) OpType {
	headType := code.ToHeaderType(tag.IsString)
	if tag.IsOmitEmpty {
		headType = headType.HeadToOmitEmptyHead()
	}
	return headType
}

func optimizeStructField(code *Opcode, tag *runtime.StructTag) OpType {
	fieldType := code.ToFieldType(tag.IsString)
	if tag.IsOmitEmpty {
		fieldType = fieldType.FieldToOmitEmptyField()
	}
	return fieldType
}

func (c *StructFieldCode) headerOpcodes(ctx *compileContext, field *Opcode, valueCodes Opcodes) Opcodes {
	value := valueCodes.First()
	op := optimizeStructHeader(value, c.tag)
	field.Op = op
	if value.Flags&MarshalerContextFlags != 0 {
		field.Flags |= MarshalerContextFlags
	}
	field.NumBitSize = value.NumBitSize
	field.PtrNum = value.PtrNum
	field.FieldQuery = value.FieldQuery
	fieldCodes := Opcodes{field}
	if op.IsMultipleOpHead() {
		field.Next = value
		fieldCodes = fieldCodes.Add(valueCodes...)
	} else {
		ctx.decIndex()
	}
	return fieldCodes
}

func (c *StructFieldCode) fieldOpcodes(ctx *compileContext, field *Opcode, valueCodes Opcodes) Opcodes {
	value := valueCodes.First()
	op := optimizeStructField(value, c.tag)
	field.Op = op
	if value.Flags&MarshalerContextFlags != 0 {
		field.Flags |= MarshalerContextFlags
	}
	field.NumBitSize = value.NumBitSize
	field.PtrNum = value.PtrNum
	field.FieldQuery = value.FieldQuery

	fieldCodes := Opcodes{field}
	if op.IsMultipleOpField() {
		field.Next = value
		fieldCodes = fieldCodes.Add(valueCodes...)
	} else {
		ctx.decIndex()
	}
	return fieldCodes
}

func (c *StructFieldCode) addStructEndCode(ctx *compileContext, codes Opcodes) Opcodes {
	end := &Opcode{
		Op:         OpStructEnd,
		Idx:        opcodeOffset(ctx.ptrIndex),
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
	}
	codes.Last().Next = end
	code := codes.First()
	for code.Op == OpStructField || code.Op == OpStructHead {
		code = code.Next
	}
	for code.NextField != nil {
		code = code.NextField
	}
	code.NextField = end

	codes = codes.Add(end)
	ctx.incOpcodeIndex()
	return codes
}

func (c *StructFieldCode) structKey(ctx *compileContext) string {
	if ctx.escapeKey {
		rctx := &RuntimeContext{Option: &Option{Flag: HTMLEscapeOption}}
		return fmt.Sprintf(`%s:`, string(AppendString(rctx, []byte{}, c.key)))
	}
	return fmt.Sprintf(`"%s":`, c.key)
}

func (c *StructFieldCode) flags() OpFlags {
	var flags OpFlags
	if c.isTaggedKey {
		flags |= IsTaggedKeyFlags
	}
	if c.isNilableType {
		flags |= IsNilableTypeFlags
	}
	if c.isNilCheck {
		flags |= NilCheckFlags
	}
	if c.isAddrForMarshaler {
		flags |= AddrForMarshalerFlags
	}
	if c.isNextOpPtrType {
		flags |= IsNextOpPtrTypeFlags
	}
	if c.isAnonymous {
		flags |= AnonymousKeyFlags
	}
	if c.isMarshalerContext {
		flags |= MarshalerContextFlags
	}
	return flags
}

func (c *StructFieldCode) toValueOpcodes(ctx *compileContext) Opcodes {
	if c.isAnonymous {
		anonymCode, ok := c.value.(AnonymousCode)
		if ok {
			return anonymCode.ToAnonymousOpcode(ctx)
		}
	}
	return c.value.ToOpcode(ctx)
}

func (c *StructFieldCode) ToOpcode(ctx *compileContext, isFirstField, isEndField bool) Opcodes {
	field := &Opcode{
		Idx:        opcodeOffset(ctx.ptrIndex),
		Flags:      c.flags(),
		Key:        c.structKey(ctx),
		Offset:     uint32(c.offset),
		Type:       c.typ,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
		DisplayKey: c.key,
	}
	ctx.incIndex()
	valueCodes := c.toValueOpcodes(ctx)
	if isFirstField {
		codes := c.headerOpcodes(ctx, field, valueCodes)
		if isEndField {
			codes = c.addStructEndCode(ctx, codes)
		}
		return codes
	}
	codes := c.fieldOpcodes(ctx, field, valueCodes)
	if isEndField {
		if isEnableStructEndOptimization(c.value) {
			field.Op = field.Op.FieldToEnd()
		} else {
			codes = c.addStructEndCode(ctx, codes)
		}
	}
	return codes
}

func (c *StructFieldCode) ToAnonymousOpcode(ctx *compileContext, isFirstField, isEndField bool) Opcodes {
	field := &Opcode{
		Idx:        opcodeOffset(ctx.ptrIndex),
		Flags:      c.flags() | AnonymousHeadFlags,
		Key:        c.structKey(ctx),
		Offset:     uint32(c.offset),
		Type:       c.typ,
		DisplayIdx: ctx.opcodeIndex,
		Indent:     ctx.indent,
		DisplayKey: c.key,
	}
	ctx.incIndex()
	valueCodes := c.toValueOpcodes(ctx)
	if isFirstField {
		return c.headerOpcodes(ctx, field, valueCodes)
	}
	return c.fieldOpcodes(ctx, field, valueCodes)
}

func isEnableStructEndOptimization(value Code) bool {
	switch value.Kind() {
	case CodeKindInt,
		CodeKindUint,
		CodeKindFloat,
		CodeKindString,
		CodeKindBool,
		CodeKindBytes:
		return true
	case CodeKindPtr:
		return isEnableStructEndOptimization(value.(*PtrCode).value)
	default:
		return false
	}
}

type InterfaceCode struct {
	typ        *runtime.Type
	fieldQuery *FieldQuery
	isPtr      bool
}

func (c *InterfaceCode) Kind() CodeKind {
	return CodeKindInterface
}

func (c *InterfaceCode) ToOpcode(ctx *compileContext) Opcodes {
	var code *Opcode
	switch {
	case c.isPtr:
		code = newOpCode(ctx, c.typ, OpInterfacePtr)
	default:
		code = newOpCode(ctx, c.typ, OpInterface)
	}
	code.FieldQuery = c.fieldQuery
	if c.typ.NumMethod() > 0 {
		code.Flags |= NonEmptyInterfaceFlags
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *InterfaceCode) Filter(query *FieldQuery) Code {
	return &InterfaceCode{
		typ:        c.typ,
		fieldQuery: query,
		isPtr:      c.isPtr,
	}
}

type MarshalJSONCode struct {
	typ                *runtime.Type
	fieldQuery         *FieldQuery
	isAddrForMarshaler bool
	isNilableType      bool
	isMarshalerContext bool
}

func (c *MarshalJSONCode) Kind() CodeKind {
	return CodeKindMarshalJSON
}

func (c *MarshalJSONCode) ToOpcode(ctx *compileContext) Opcodes {
	code := newOpCode(ctx, c.typ, OpMarshalJSON)
	code.FieldQuery = c.fieldQuery
	if c.isAddrForMarshaler {
		code.Flags |= AddrForMarshalerFlags
	}
	if c.isMarshalerContext {
		code.Flags |= MarshalerContextFlags
	}
	if c.isNilableType {
		code.Flags |= IsNilableTypeFlags
	} else {
		code.Flags &= ^IsNilableTypeFlags
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *MarshalJSONCode) Filter(query *FieldQuery) Code {
	return &MarshalJSONCode{
		typ:                c.typ,
		fieldQuery:         query,
		isAddrForMarshaler: c.isAddrForMarshaler,
		isNilableType:      c.isNilableType,
		isMarshalerContext: c.isMarshalerContext,
	}
}

type MarshalTextCode struct {
	typ                *runtime.Type
	fieldQuery         *FieldQuery
	isAddrForMarshaler bool
	isNilableType      bool
}

func (c *MarshalTextCode) Kind() CodeKind {
	return CodeKindMarshalText
}

func (c *MarshalTextCode) ToOpcode(ctx *compileContext) Opcodes {
	code := newOpCode(ctx, c.typ, OpMarshalText)
	code.FieldQuery = c.fieldQuery
	if c.isAddrForMarshaler {
		code.Flags |= AddrForMarshalerFlags
	}
	if c.isNilableType {
		code.Flags |= IsNilableTypeFlags
	} else {
		code.Flags &= ^IsNilableTypeFlags
	}
	ctx.incIndex()
	return Opcodes{code}
}

func (c *MarshalTextCode) Filter(query *FieldQuery) Code {
	return &MarshalTextCode{
		typ:                c.typ,
		fieldQuery:         query,
		isAddrForMarshaler: c.isAddrForMarshaler,
		isNilableType:      c.isNilableType,
	}
}

type PtrCode struct {
	typ    *runtime.Type
	value  Code
	ptrNum uint8
}

func (c *PtrCode) Kind() CodeKind {
	return CodeKindPtr
}

func (c *PtrCode) ToOpcode(ctx *compileContext) Opcodes {
	codes := c.value.ToOpcode(ctx)
	codes.First().Op = convertPtrOp(codes.First())
	codes.First().PtrNum = c.ptrNum
	return codes
}

func (c *PtrCode) ToAnonymousOpcode(ctx *compileContext) Opcodes {
	var codes Opcodes
	anonymCode, ok := c.value.(AnonymousCode)
	if ok {
		codes = anonymCode.ToAnonymousOpcode(ctx)
	} else {
		codes = c.value.ToOpcode(ctx)
	}
	codes.First().Op = convertPtrOp(codes.First())
	codes.First().PtrNum = c.ptrNum
	return codes
}

func (c *PtrCode) Filter(query *FieldQuery) Code {
	return &PtrCode{
		typ:    c.typ,
		value:  c.value.Filter(query),
		ptrNum: c.ptrNum,
	}
}

func convertPtrOp(code *Opcode) OpType {
	ptrHeadOp := code.Op.HeadToPtrHead()
	if code.Op != ptrHeadOp {
		if code.PtrNum > 0 {
			// ptr field and ptr head
			code.PtrNum--
		}
		return ptrHeadOp
	}
	switch code.Op {
	case OpInt:
		return OpIntPtr
	case OpUint:
		return OpUintPtr
	case OpFloat32:
		return OpFloat32Ptr
	case OpFloat64:
		return OpFloat64Ptr
	case OpString:
		return OpStringPtr
	case OpBool:
		return OpBoolPtr
	case OpBytes:
		return OpBytesPtr
	case OpNumber:
		return OpNumberPtr
	case OpArray:
		return OpArrayPtr
	case OpSlice:
		return OpSlicePtr
	case OpMap:
		return OpMapPtr
	case OpMarshalJSON:
		return OpMarshalJSONPtr
	case OpMarshalText:
		return OpMarshalTextPtr
	case OpInterface:
		return OpInterfacePtr
	case OpRecursive:
		return OpRecursivePtr
	}
	return code.Op
}

func isEmbeddedStruct(field *StructFieldCode) bool {
	if !field.isAnonymous {
		return false
	}
	t := field.typ
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}
