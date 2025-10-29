// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build jsoniter

package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// Package indicates what library is being used for JSON encoding.
const Package = "github.com/json-iterator/go"

func containsAF(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'f' {
			return true
		}
	}
	return false
}

// HexStringEncoder 自定义编码器将 int64 类型编码为十六进制字符串或者把16进制转为int64
type HexStringEncoder struct {
	IsInt64Pointer bool //源数据为 *int64
}

// Encode 实现 jsoniter.ValEncoder 接口
func (e *HexStringEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if ptr == nil {
		stream.WriteNil()
		return
	}

	// Convert *int64 value to a 16-byte hexadecimal string
	if e.IsInt64Pointer {
		pp := (**int64)(ptr)
		if *pp == nil {
			stream.WriteNil()
			return
		}
		stream.WriteString(fmt.Sprintf("%x", **pp))
		return
	}

	// Convert int64 value to a 16-byte hexadecimal string
	value := *(*int64)(ptr)
	if value == 0 {
		stream.WriteString("0") //0值特殊处理
	} else {
		stream.WriteString(fmt.Sprintf("%x", value))
	}
}

func (e *HexStringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return ptr == nil || *(*int64)(ptr) == 0
}

func (codec *HexStringEncoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	if valueType == jsoniter.StringValue {
		str := iter.ReadString()
		var i int64
		var err error
		if len(str) < 17 || containsAF(str) {
			i, err = strconv.ParseInt(str, 16, 64)
			if err != nil {
				i = 0
			}
		} else {
			i, err = strconv.ParseInt(str, 10, 64)
			if err != nil {
				i = 0
			}
		}

		*((*int64)(ptr)) = i
	} else if valueType == jsoniter.NumberValue {
		*((*int64)(ptr)) = iter.ReadInt64()
	} else {
		*((*int64)(ptr)) = 0
	}
}

// EmptyObjectEncoder 实现一个编码器，当字段值为nil时，写入空对象{}
type EmptyObjectEncoder struct {
	encoder jsoniter.ValEncoder
}

func (encoder *EmptyObjectEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	// If the pointer points to nil, write an empty object.
	if *(*uintptr)(ptr) == 0 {
		stream.WriteRaw("{}")
		return
	}
	// Fallback to default encoding.
	encoder.encoder.Encode(ptr, stream)
}

func (encoder *EmptyObjectEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(ptr)
}

type EmptyArrayEncoder struct {
	encoder jsoniter.ValEncoder
}

func (encoder *EmptyArrayEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	// If the pointer points to nil, write an empty object.
	if *(*uintptr)(ptr) == 0 {
		stream.WriteRaw("[]")
		return
	}
	// Fallback to default encoding.
	encoder.encoder.Encode(ptr, stream)
}

func (encoder *EmptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(ptr)
}

// 空数组64位数组
type EmptyArrayInt64Encoder struct {
	encoder jsoniter.ValEncoder
	decoder jsoniter.ValDecoder
}

func (encoder *EmptyArrayInt64Encoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	// If the pointer points to nil, write an empty object.
	if *(*uintptr)(ptr) == 0 {
		stream.WriteRaw("[]")
		return
	}

	//循环数组
	slice := (*[]int64)(ptr)
	strSlice := make([]string, len(*slice))
	for i, v := range *slice {
		if v == 0 {
			strSlice[i] = "0"
		} else {
			strSlice[i] = fmt.Sprintf("%x", v)
		}
	}

	jsonData, err := jsonInstance.Marshal(strSlice)
	if err != nil {
		encoder.encoder.Encode(ptr, stream)
		return
	}
	stream.WriteRaw(string(jsonData))
}

func (codec *EmptyArrayInt64Encoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	//str := iter.ReadString()
	valueList := []int64{}
	for iter.ReadArray() {
		val := iter.Read()
		if iter.Error != nil {
			break
		}

		if str, ok := val.(string); ok {
			//含有a-f或者正好16位使用16进制解析
			if len(str) < 17 || containsAF(str) {
				intVal, err := strconv.ParseInt(str, 16, 64)
				if err != nil {
					continue
				}
				valueList = append(valueList, intVal)
			} else {
				intVal, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					continue
				}
				valueList = append(valueList, intVal)
			}
		}
	}

	// 将ptr解析为*[]int64类型的指针
	ptrToSlice := (*[]int64)(ptr)
	// 使用reflect包将ptr的内容替换为slice
	reflect.ValueOf(ptrToSlice).Elem().Set(reflect.ValueOf(valueList))
}

func (encoder *EmptyArrayInt64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(ptr)
}

type ToStringEncoder struct{}

func (codec *ToStringEncoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	if valueType == jsoniter.StringValue {
		str := iter.ReadString()
		*((*string)(ptr)) = str
	} else {
		valueAsString := iter.ReadAny().ToString()
		*((*string)(ptr)) = valueAsString
	}
}

type ToBoolEncoder struct {
	decoder    jsoniter.ValDecoder
	DefaultVal bool
}

func (codec *ToBoolEncoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	if valueType == jsoniter.BoolValue {
		codec.decoder.Decode(ptr, iter)
	} else if valueType == jsoniter.NumberValue {
		i := iter.ReadInt()
		var myBool bool
		if i > 0 {
			myBool = true
		} else {
			myBool = false
		}
		*((*bool)(ptr)) = myBool
	} else if valueType == jsoniter.StringValue {
		str := strings.ToLower(iter.ReadString())
		var myBool bool
		if str == "true" || str == "1" {
			myBool = true
		} else if str == "false" || str == "0" || str == "-1" {
			myBool = false
		} else {
			myBool = codec.DefaultVal
		}
		*((*bool)(ptr)) = myBool
	} else {
		str := strings.ToLower(iter.ReadAny().ToString())
		var myBool bool
		if str == "true" || str == "1" {
			myBool = true
		} else if str == "false" || str == "0" || str == "-1" {
			myBool = false
		} else {
			myBool = codec.DefaultVal
		}
		*((*bool)(ptr)) = myBool
	}
}

// HexStringExtension 检查 struct 字段tags，为相应的 int64 字段应用 HexStringEncoder
type ApipostExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor 修改 struct 字段的编码/解码器
func (extension *ApipostExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		// 检查字段类型和 tag
		if binding.Field.Type().Kind() == reflect.Int64 {
			//处理64位转换
			if strings.Contains(binding.Field.Tag().Get("json"), "hexstring") {
				binding.Encoder = &HexStringEncoder{}
				binding.Decoder = &HexStringEncoder{}
			}
		} else if binding.Field.Type().Kind() == reflect.Ptr || binding.Field.Type().Kind() == reflect.Interface {
			if strings.Contains(binding.Field.Tag().Get("json"), "hexstring") && binding.Field.Type().Type1().Elem().Kind() == reflect.Int64 { //处理*int64位转换
				binding.Encoder = &HexStringEncoder{IsInt64Pointer: true}
			} else if strings.Contains(binding.Field.Tag().Get("json"), "emptyobject") { //处理空对象
				binding.Encoder = &EmptyObjectEncoder{binding.Encoder}
			}
		} else if binding.Field.Type().Kind() == reflect.Slice || binding.Field.Type().Kind() == reflect.Array {
			//处理空数组
			if binding.Field.Type().Type1().Elem().String() == "int64" {
				//强制转64数组
				int64SliceEncode := &EmptyArrayInt64Encoder{binding.Encoder, binding.Decoder}
				binding.Encoder = int64SliceEncode
				binding.Decoder = int64SliceEncode
			} else if strings.Contains(binding.Field.Tag().Get("json"), "emptyarray") {
				binding.Encoder = &EmptyArrayEncoder{binding.Encoder}
			}
		} else if binding.Field.Type().Kind() == reflect.String {
			if strings.Contains(binding.Field.Tag().Get("json"), "tostring") {
				binding.Decoder = &ToStringEncoder{}
			}
		} else if binding.Field.Type().Kind() == reflect.Bool {
			tagStr := binding.Field.Tag().Get("json")
			if strings.Contains(tagStr, "tofalse") {
				binding.Decoder = &ToBoolEncoder{binding.Decoder, false}
			} else if strings.Contains(tagStr, "totrue") {
				binding.Decoder = &ToBoolEncoder{binding.Decoder, true}
			}
		}
	}
}

func init() {
	json.RegisterExtension(&ApipostExtension{})
	API = jsoniterApi{}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsoniterApi struct{}

func (j jsoniterApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsoniterApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j jsoniterApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j jsoniterApi) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j jsoniterApi) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}
