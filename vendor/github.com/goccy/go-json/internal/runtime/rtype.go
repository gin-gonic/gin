package runtime

import (
	"reflect"
	"unsafe"
)

// Type representing reflect.rtype for noescape trick
type Type struct{}

//go:linkname rtype_Align reflect.(*rtype).Align
//go:noescape
func rtype_Align(*Type) int

func (t *Type) Align() int {
	return rtype_Align(t)
}

//go:linkname rtype_FieldAlign reflect.(*rtype).FieldAlign
//go:noescape
func rtype_FieldAlign(*Type) int

func (t *Type) FieldAlign() int {
	return rtype_FieldAlign(t)
}

//go:linkname rtype_Method reflect.(*rtype).Method
//go:noescape
func rtype_Method(*Type, int) reflect.Method

func (t *Type) Method(a0 int) reflect.Method {
	return rtype_Method(t, a0)
}

//go:linkname rtype_MethodByName reflect.(*rtype).MethodByName
//go:noescape
func rtype_MethodByName(*Type, string) (reflect.Method, bool)

func (t *Type) MethodByName(a0 string) (reflect.Method, bool) {
	return rtype_MethodByName(t, a0)
}

//go:linkname rtype_NumMethod reflect.(*rtype).NumMethod
//go:noescape
func rtype_NumMethod(*Type) int

func (t *Type) NumMethod() int {
	return rtype_NumMethod(t)
}

//go:linkname rtype_Name reflect.(*rtype).Name
//go:noescape
func rtype_Name(*Type) string

func (t *Type) Name() string {
	return rtype_Name(t)
}

//go:linkname rtype_PkgPath reflect.(*rtype).PkgPath
//go:noescape
func rtype_PkgPath(*Type) string

func (t *Type) PkgPath() string {
	return rtype_PkgPath(t)
}

//go:linkname rtype_Size reflect.(*rtype).Size
//go:noescape
func rtype_Size(*Type) uintptr

func (t *Type) Size() uintptr {
	return rtype_Size(t)
}

//go:linkname rtype_String reflect.(*rtype).String
//go:noescape
func rtype_String(*Type) string

func (t *Type) String() string {
	return rtype_String(t)
}

//go:linkname rtype_Kind reflect.(*rtype).Kind
//go:noescape
func rtype_Kind(*Type) reflect.Kind

func (t *Type) Kind() reflect.Kind {
	return rtype_Kind(t)
}

//go:linkname rtype_Implements reflect.(*rtype).Implements
//go:noescape
func rtype_Implements(*Type, reflect.Type) bool

func (t *Type) Implements(u reflect.Type) bool {
	return rtype_Implements(t, u)
}

//go:linkname rtype_AssignableTo reflect.(*rtype).AssignableTo
//go:noescape
func rtype_AssignableTo(*Type, reflect.Type) bool

func (t *Type) AssignableTo(u reflect.Type) bool {
	return rtype_AssignableTo(t, u)
}

//go:linkname rtype_ConvertibleTo reflect.(*rtype).ConvertibleTo
//go:noescape
func rtype_ConvertibleTo(*Type, reflect.Type) bool

func (t *Type) ConvertibleTo(u reflect.Type) bool {
	return rtype_ConvertibleTo(t, u)
}

//go:linkname rtype_Comparable reflect.(*rtype).Comparable
//go:noescape
func rtype_Comparable(*Type) bool

func (t *Type) Comparable() bool {
	return rtype_Comparable(t)
}

//go:linkname rtype_Bits reflect.(*rtype).Bits
//go:noescape
func rtype_Bits(*Type) int

func (t *Type) Bits() int {
	return rtype_Bits(t)
}

//go:linkname rtype_ChanDir reflect.(*rtype).ChanDir
//go:noescape
func rtype_ChanDir(*Type) reflect.ChanDir

func (t *Type) ChanDir() reflect.ChanDir {
	return rtype_ChanDir(t)
}

//go:linkname rtype_IsVariadic reflect.(*rtype).IsVariadic
//go:noescape
func rtype_IsVariadic(*Type) bool

func (t *Type) IsVariadic() bool {
	return rtype_IsVariadic(t)
}

//go:linkname rtype_Elem reflect.(*rtype).Elem
//go:noescape
func rtype_Elem(*Type) reflect.Type

func (t *Type) Elem() *Type {
	return Type2RType(rtype_Elem(t))
}

//go:linkname rtype_Field reflect.(*rtype).Field
//go:noescape
func rtype_Field(*Type, int) reflect.StructField

func (t *Type) Field(i int) reflect.StructField {
	return rtype_Field(t, i)
}

//go:linkname rtype_FieldByIndex reflect.(*rtype).FieldByIndex
//go:noescape
func rtype_FieldByIndex(*Type, []int) reflect.StructField

func (t *Type) FieldByIndex(index []int) reflect.StructField {
	return rtype_FieldByIndex(t, index)
}

//go:linkname rtype_FieldByName reflect.(*rtype).FieldByName
//go:noescape
func rtype_FieldByName(*Type, string) (reflect.StructField, bool)

func (t *Type) FieldByName(name string) (reflect.StructField, bool) {
	return rtype_FieldByName(t, name)
}

//go:linkname rtype_FieldByNameFunc reflect.(*rtype).FieldByNameFunc
//go:noescape
func rtype_FieldByNameFunc(*Type, func(string) bool) (reflect.StructField, bool)

func (t *Type) FieldByNameFunc(match func(string) bool) (reflect.StructField, bool) {
	return rtype_FieldByNameFunc(t, match)
}

//go:linkname rtype_In reflect.(*rtype).In
//go:noescape
func rtype_In(*Type, int) reflect.Type

func (t *Type) In(i int) reflect.Type {
	return rtype_In(t, i)
}

//go:linkname rtype_Key reflect.(*rtype).Key
//go:noescape
func rtype_Key(*Type) reflect.Type

func (t *Type) Key() *Type {
	return Type2RType(rtype_Key(t))
}

//go:linkname rtype_Len reflect.(*rtype).Len
//go:noescape
func rtype_Len(*Type) int

func (t *Type) Len() int {
	return rtype_Len(t)
}

//go:linkname rtype_NumField reflect.(*rtype).NumField
//go:noescape
func rtype_NumField(*Type) int

func (t *Type) NumField() int {
	return rtype_NumField(t)
}

//go:linkname rtype_NumIn reflect.(*rtype).NumIn
//go:noescape
func rtype_NumIn(*Type) int

func (t *Type) NumIn() int {
	return rtype_NumIn(t)
}

//go:linkname rtype_NumOut reflect.(*rtype).NumOut
//go:noescape
func rtype_NumOut(*Type) int

func (t *Type) NumOut() int {
	return rtype_NumOut(t)
}

//go:linkname rtype_Out reflect.(*rtype).Out
//go:noescape
func rtype_Out(*Type, int) reflect.Type

//go:linkname PtrTo reflect.(*rtype).ptrTo
//go:noescape
func PtrTo(*Type) *Type

func (t *Type) Out(i int) reflect.Type {
	return rtype_Out(t, i)
}

//go:linkname IfaceIndir reflect.ifaceIndir
//go:noescape
func IfaceIndir(*Type) bool

//go:linkname RType2Type reflect.toType
//go:noescape
func RType2Type(t *Type) reflect.Type

//go:nolint structcheck
type emptyInterface struct {
	_   *Type
	ptr unsafe.Pointer
}

func Type2RType(t reflect.Type) *Type {
	return (*Type)(((*emptyInterface)(unsafe.Pointer(&t))).ptr)
}
