package rt

import (
	"reflect"
	"unsafe"
	"encoding/json"
)

func AsGoType(t uintptr) *GoType {
	return (*GoType)(unsafe.Pointer(t))
}

var (
	BoolType				= UnpackType(reflect.TypeOf(false))
	ByteType                = UnpackType(reflect.TypeOf(byte(0)))
	IncntType               = UnpackType(reflect.TypeOf(int(0)))
	Int8Type                = UnpackType(reflect.TypeOf(int8(0)))
	Int16Type               = UnpackType(reflect.TypeOf(int16(0)))
	Int32Type               = UnpackType(reflect.TypeOf(int32(0)))
	Int64Type               = UnpackType(reflect.TypeOf(int64(0)))
	UintType                = UnpackType(reflect.TypeOf(uint(0)))
	Uint8Type               = UnpackType(reflect.TypeOf(uint8(0)))
	Uint16Type              = UnpackType(reflect.TypeOf(uint16(0)))
	Uint32Type              = UnpackType(reflect.TypeOf(uint32(0)))
	Uint64Type              = UnpackType(reflect.TypeOf(uint64(0)))
	Float32Type             = UnpackType(reflect.TypeOf(float32(0)))
	Float64Type             = UnpackType(reflect.TypeOf(float64(0)))

	StringType              = UnpackType(reflect.TypeOf(""))
	BytesType               = UnpackType(reflect.TypeOf([]byte(nil)))
	JsonNumberType          = UnpackType(reflect.TypeOf(json.Number("")))

	SliceEfaceType          = UnpackType(reflect.TypeOf([]interface{}(nil)))
	SliceStringType         = UnpackType(reflect.TypeOf([]string(nil)))
	SliceI32Type          	= UnpackType(reflect.TypeOf([]int32(nil)))
	SliceI64Type          	= UnpackType(reflect.TypeOf([]int64(nil)))
	SliceU32Type          	= UnpackType(reflect.TypeOf([]uint32(nil)))
	SliceU64Type          	= UnpackType(reflect.TypeOf([]uint64(nil)))

	AnyType    				= UnpackType(reflect.TypeOf((*interface{})(nil)).Elem())
	MapEfaceType    		= UnpackType(reflect.TypeOf(map[string]interface{}(nil)))
	MapStringType    		= UnpackType(reflect.TypeOf(map[string]string(nil)))

	MapEfaceMapType    		= MapType(UnpackType(reflect.TypeOf(map[string]interface{}(nil))))
)
