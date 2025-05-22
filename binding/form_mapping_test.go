// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/hex"
	"errors"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMappingBaseTypes(t *testing.T) {
	intPtr := func(i int) *int {
		return &i
	}
	for _, tt := range []struct {
		name   string
		value  any
		form   string
		expect any
	}{
		{"base type", struct{ F int }{}, "9", int(9)},
		{"base type", struct{ F int8 }{}, "9", int8(9)},
		{"base type", struct{ F int16 }{}, "9", int16(9)},
		{"base type", struct{ F int32 }{}, "9", int32(9)},
		{"base type", struct{ F int64 }{}, "9", int64(9)},
		{"base type", struct{ F uint }{}, "9", uint(9)},
		{"base type", struct{ F uint8 }{}, "9", uint8(9)},
		{"base type", struct{ F uint16 }{}, "9", uint16(9)},
		{"base type", struct{ F uint32 }{}, "9", uint32(9)},
		{"base type", struct{ F uint64 }{}, "9", uint64(9)},
		{"base type", struct{ F bool }{}, "True", true},
		{"base type", struct{ F float32 }{}, "9.1", float32(9.1)},
		{"base type", struct{ F float64 }{}, "9.1", float64(9.1)},
		{"base type", struct{ F string }{}, "test", string("test")},
		{"base type", struct{ F *int }{}, "9", intPtr(9)},

		// zero values
		{"zero value", struct{ F int }{}, "", int(0)},
		{"zero value", struct{ F uint }{}, "", uint(0)},
		{"zero value", struct{ F bool }{}, "", false},
		{"zero value", struct{ F float32 }{}, "", float32(0)},
		{"file value", struct{ F *multipart.FileHeader }{}, "", &multipart.FileHeader{}},
	} {
		tp := reflect.TypeOf(tt.value)
		testName := tt.name + ":" + tp.Field(0).Type.String()

		val := reflect.New(reflect.TypeOf(tt.value))
		val.Elem().Set(reflect.ValueOf(tt.value))

		field := val.Elem().Type().Field(0)

		_, err := mapping(val, emptyField, formSource{field.Name: {tt.form}}, "form")
		require.NoError(t, err, testName)

		actual := val.Elem().Field(0).Interface()
		assert.Equal(t, tt.expect, actual, testName)
	}
}

func TestMappingDefault(t *testing.T) {
	var s struct {
		Str   string `form:",default=defaultVal"`
		Int   int    `form:",default=9"`
		Slice []int  `form:",default=9"`
		Array [1]int `form:",default=9"`
	}
	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)

	assert.Equal(t, "defaultVal", s.Str)
	assert.Equal(t, 9, s.Int)
	assert.Equal(t, []int{9}, s.Slice)
	assert.Equal(t, [1]int{9}, s.Array)
}

func TestMappingSkipField(t *testing.T) {
	var s struct {
		A int
	}
	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)

	assert.Equal(t, 0, s.A)
}

func TestMappingIgnoreField(t *testing.T) {
	var s struct {
		A int `form:"A"`
		B int `form:"-"`
	}
	err := mappingByPtr(&s, formSource{"A": {"9"}, "B": {"9"}}, "form")
	require.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.B)
}

func TestMappingUnexportedField(t *testing.T) {
	var s struct {
		A int `form:"a"`
		b int `form:"b"`
	}
	err := mappingByPtr(&s, formSource{"a": {"9"}, "b": {"9"}}, "form")
	require.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.b)
}

func TestMappingPrivateField(t *testing.T) {
	var s struct {
		f int `form:"field"`
	}
	err := mappingByPtr(&s, formSource{"field": {"6"}}, "form")
	require.NoError(t, err)
	assert.Equal(t, 0, s.f)
}

func TestMappingUnknownFieldType(t *testing.T) {
	var s struct {
		U uintptr
	}

	err := mappingByPtr(&s, formSource{"U": {"unknown"}}, "form")
	require.Error(t, err)
	assert.Equal(t, errUnknownType, err)
}

func TestMappingURI(t *testing.T) {
	var s struct {
		F int `uri:"field"`
	}
	err := mapURI(&s, map[string][]string{"field": {"6"}})
	require.NoError(t, err)
	assert.Equal(t, 6, s.F)
}

func TestMappingForm(t *testing.T) {
	var s struct {
		F int `form:"field"`
	}
	err := mapForm(&s, map[string][]string{"field": {"6"}})
	require.NoError(t, err)
	assert.Equal(t, 6, s.F)
}

func TestMappingFormFieldNotSent(t *testing.T) {
	var s struct {
		F string `form:"field,default=defVal"`
	}
	err := mapForm(&s, map[string][]string{})
	require.NoError(t, err)
	assert.Equal(t, "defVal", s.F)
}

func TestMappingFormWithEmptyToDefault(t *testing.T) {
	var s struct {
		F string `form:"field,default=DefVal"`
	}
	err := mapForm(&s, map[string][]string{"field": {""}})
	require.NoError(t, err)
	assert.Equal(t, "DefVal", s.F)
}

func TestMapFormWithTag(t *testing.T) {
	var s struct {
		F int `externalTag:"field"`
	}
	err := MapFormWithTag(&s, map[string][]string{"field": {"6"}}, "externalTag")
	require.NoError(t, err)
	assert.Equal(t, 6, s.F)
}

func TestMappingTime(t *testing.T) {
	var s struct {
		Time      time.Time
		LocalTime time.Time `time_format:"2006-01-02"`
		ZeroValue time.Time
		CSTTime   time.Time `time_format:"2006-01-02" time_location:"Asia/Shanghai"`
		UTCTime   time.Time `time_format:"2006-01-02" time_utc:"1"`
	}

	var err error
	time.Local, err = time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)

	err = mapForm(&s, map[string][]string{
		"Time":      {"2019-01-20T16:02:58Z"},
		"LocalTime": {"2019-01-20"},
		"ZeroValue": {},
		"CSTTime":   {"2019-01-20"},
		"UTCTime":   {"2019-01-20"},
	})
	require.NoError(t, err)

	assert.Equal(t, "2019-01-20 16:02:58 +0000 UTC", s.Time.String())
	assert.Equal(t, "2019-01-20 00:00:00 +0100 CET", s.LocalTime.String())
	assert.Equal(t, "2019-01-19 23:00:00 +0000 UTC", s.LocalTime.UTC().String())
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", s.ZeroValue.String())
	assert.Equal(t, "2019-01-20 00:00:00 +0800 CST", s.CSTTime.String())
	assert.Equal(t, "2019-01-19 16:00:00 +0000 UTC", s.CSTTime.UTC().String())
	assert.Equal(t, "2019-01-20 00:00:00 +0000 UTC", s.UTCTime.String())

	// wrong location
	var wrongLoc struct {
		Time time.Time `time_location:"wrong"`
	}
	err = mapForm(&wrongLoc, map[string][]string{"Time": {"2019-01-20T16:02:58Z"}})
	require.Error(t, err)

	// wrong time value
	var wrongTime struct {
		Time time.Time
	}
	err = mapForm(&wrongTime, map[string][]string{"Time": {"wrong"}})
	require.Error(t, err)
}

func TestMappingTimeDuration(t *testing.T) {
	var s struct {
		D time.Duration
	}

	// ok
	err := mappingByPtr(&s, formSource{"D": {"5s"}}, "form")
	require.NoError(t, err)
	assert.Equal(t, 5*time.Second, s.D)

	// error
	err = mappingByPtr(&s, formSource{"D": {"wrong"}}, "form")
	require.Error(t, err)
}

func TestMappingSlice(t *testing.T) {
	var s struct {
		Slice []int `form:"slice,default=9"`
	}

	// default value
	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)
	assert.Equal(t, []int{9}, s.Slice)

	// ok
	err = mappingByPtr(&s, formSource{"slice": {"3", "4"}}, "form")
	require.NoError(t, err)
	assert.Equal(t, []int{3, 4}, s.Slice)

	// error
	err = mappingByPtr(&s, formSource{"slice": {"wrong"}}, "form")
	require.Error(t, err)
}

func TestMappingArray(t *testing.T) {
	var s struct {
		Array [2]int `form:"array,default=9"`
	}

	// wrong default
	err := mappingByPtr(&s, formSource{}, "form")
	require.Error(t, err)

	// ok
	err = mappingByPtr(&s, formSource{"array": {"3", "4"}}, "form")
	require.NoError(t, err)
	assert.Equal(t, [2]int{3, 4}, s.Array)

	// error - not enough vals
	err = mappingByPtr(&s, formSource{"array": {"3"}}, "form")
	require.Error(t, err)

	// error - wrong value
	err = mappingByPtr(&s, formSource{"array": {"wrong"}}, "form")
	require.Error(t, err)
}

func TestMappingCollectionFormat(t *testing.T) {
	var s struct {
		SliceMulti []int  `form:"slice_multi" collection_format:"multi"`
		SliceCsv   []int  `form:"slice_csv" collection_format:"csv"`
		SliceSsv   []int  `form:"slice_ssv" collection_format:"ssv"`
		SliceTsv   []int  `form:"slice_tsv" collection_format:"tsv"`
		SlicePipes []int  `form:"slice_pipes" collection_format:"pipes"`
		ArrayMulti [2]int `form:"array_multi" collection_format:"multi"`
		ArrayCsv   [2]int `form:"array_csv" collection_format:"csv"`
		ArraySsv   [2]int `form:"array_ssv" collection_format:"ssv"`
		ArrayTsv   [2]int `form:"array_tsv" collection_format:"tsv"`
		ArrayPipes [2]int `form:"array_pipes" collection_format:"pipes"`
	}
	err := mappingByPtr(&s, formSource{
		"slice_multi": {"1", "2"},
		"slice_csv":   {"1,2"},
		"slice_ssv":   {"1 2"},
		"slice_tsv":   {"1	2"},
		"slice_pipes": {"1|2"},
		"array_multi": {"1", "2"},
		"array_csv":   {"1,2"},
		"array_ssv":   {"1 2"},
		"array_tsv":   {"1	2"},
		"array_pipes": {"1|2"},
	}, "form")
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2}, s.SliceMulti)
	assert.Equal(t, []int{1, 2}, s.SliceCsv)
	assert.Equal(t, []int{1, 2}, s.SliceSsv)
	assert.Equal(t, []int{1, 2}, s.SliceTsv)
	assert.Equal(t, []int{1, 2}, s.SlicePipes)
	assert.Equal(t, [2]int{1, 2}, s.ArrayMulti)
	assert.Equal(t, [2]int{1, 2}, s.ArrayCsv)
	assert.Equal(t, [2]int{1, 2}, s.ArraySsv)
	assert.Equal(t, [2]int{1, 2}, s.ArrayTsv)
	assert.Equal(t, [2]int{1, 2}, s.ArrayPipes)
}

func TestMappingCollectionFormatInvalid(t *testing.T) {
	var s struct {
		SliceCsv []int `form:"slice_csv" collection_format:"xxx"`
	}
	err := mappingByPtr(&s, formSource{
		"slice_csv": {"1,2"},
	}, "form")
	require.Error(t, err)

	var s2 struct {
		ArrayCsv [2]int `form:"array_csv" collection_format:"xxx"`
	}
	err = mappingByPtr(&s2, formSource{
		"array_csv": {"1,2"},
	}, "form")
	require.Error(t, err)
}

func TestMappingMultipleDefaultWithCollectionFormat(t *testing.T) {
	var s struct {
		SliceMulti       []int     `form:",default=1;2;3" collection_format:"multi"`
		SliceCsv         []int     `form:",default=1;2;3" collection_format:"csv"`
		SliceSsv         []int     `form:",default=1 2 3" collection_format:"ssv"`
		SliceTsv         []int     `form:",default=1\t2\t3" collection_format:"tsv"`
		SlicePipes       []int     `form:",default=1|2|3" collection_format:"pipes"`
		ArrayMulti       [2]int    `form:",default=1;2" collection_format:"multi"`
		ArrayCsv         [2]int    `form:",default=1;2" collection_format:"csv"`
		ArraySsv         [2]int    `form:",default=1 2" collection_format:"ssv"`
		ArrayTsv         [2]int    `form:",default=1\t2" collection_format:"tsv"`
		ArrayPipes       [2]int    `form:",default=1|2" collection_format:"pipes"`
		SliceStringMulti []string  `form:",default=1;2;3" collection_format:"multi"`
		SliceStringCsv   []string  `form:",default=1;2;3" collection_format:"csv"`
		SliceStringSsv   []string  `form:",default=1 2 3" collection_format:"ssv"`
		SliceStringTsv   []string  `form:",default=1\t2\t3" collection_format:"tsv"`
		SliceStringPipes []string  `form:",default=1|2|3" collection_format:"pipes"`
		ArrayStringMulti [2]string `form:",default=1;2" collection_format:"multi"`
		ArrayStringCsv   [2]string `form:",default=1;2" collection_format:"csv"`
		ArrayStringSsv   [2]string `form:",default=1 2" collection_format:"ssv"`
		ArrayStringTsv   [2]string `form:",default=1\t2" collection_format:"tsv"`
		ArrayStringPipes [2]string `form:",default=1|2" collection_format:"pipes"`
	}
	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3}, s.SliceMulti)
	assert.Equal(t, []int{1, 2, 3}, s.SliceCsv)
	assert.Equal(t, []int{1, 2, 3}, s.SliceSsv)
	assert.Equal(t, []int{1, 2, 3}, s.SliceTsv)
	assert.Equal(t, []int{1, 2, 3}, s.SlicePipes)
	assert.Equal(t, [2]int{1, 2}, s.ArrayMulti)
	assert.Equal(t, [2]int{1, 2}, s.ArrayCsv)
	assert.Equal(t, [2]int{1, 2}, s.ArraySsv)
	assert.Equal(t, [2]int{1, 2}, s.ArrayTsv)
	assert.Equal(t, [2]int{1, 2}, s.ArrayPipes)
	assert.Equal(t, []string{"1", "2", "3"}, s.SliceStringMulti)
	assert.Equal(t, []string{"1", "2", "3"}, s.SliceStringCsv)
	assert.Equal(t, []string{"1", "2", "3"}, s.SliceStringSsv)
	assert.Equal(t, []string{"1", "2", "3"}, s.SliceStringTsv)
	assert.Equal(t, []string{"1", "2", "3"}, s.SliceStringPipes)
	assert.Equal(t, [2]string{"1", "2"}, s.ArrayStringMulti)
	assert.Equal(t, [2]string{"1", "2"}, s.ArrayStringCsv)
	assert.Equal(t, [2]string{"1", "2"}, s.ArrayStringSsv)
	assert.Equal(t, [2]string{"1", "2"}, s.ArrayStringTsv)
	assert.Equal(t, [2]string{"1", "2"}, s.ArrayStringPipes)
}

func TestMappingStructField(t *testing.T) {
	var s struct {
		J struct {
			I int
		}
	}

	err := mappingByPtr(&s, formSource{"J": {`{"I": 9}`}}, "form")
	require.NoError(t, err)
	assert.Equal(t, 9, s.J.I)
}

func TestMappingPtrField(t *testing.T) {
	type ptrStruct struct {
		Key int64 `json:"key"`
	}

	type ptrRequest struct {
		Items []*ptrStruct `json:"items" form:"items"`
	}

	var err error

	// With 0 items.
	var req0 ptrRequest
	err = mappingByPtr(&req0, formSource{}, "form")
	require.NoError(t, err)
	assert.Empty(t, req0.Items)

	// With 1 item.
	var req1 ptrRequest
	err = mappingByPtr(&req1, formSource{"items": {`{"key": 1}`}}, "form")
	require.NoError(t, err)
	assert.Len(t, req1.Items, 1)
	assert.EqualValues(t, 1, req1.Items[0].Key)

	// With 2 items.
	var req2 ptrRequest
	err = mappingByPtr(&req2, formSource{"items": {`{"key": 1}`, `{"key": 2}`}}, "form")
	require.NoError(t, err)
	assert.Len(t, req2.Items, 2)
	assert.EqualValues(t, 1, req2.Items[0].Key)
	assert.EqualValues(t, 2, req2.Items[1].Key)
}

func TestMappingMapField(t *testing.T) {
	var s struct {
		M map[string]int
	}

	err := mappingByPtr(&s, formSource{"M": {`{"one": 1}`}}, "form")
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"one": 1}, s.M)
}

func TestMappingIgnoredCircularRef(t *testing.T) {
	type S struct {
		S *S `form:"-"`
	}
	var s S

	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)
}

type customUnmarshalParamHex int

func (f *customUnmarshalParamHex) UnmarshalParam(param string) error {
	v, err := strconv.ParseInt(param, 16, 64)
	if err != nil {
		return err
	}
	*f = customUnmarshalParamHex(v)
	return nil
}

func TestMappingCustomUnmarshalParamHexWithFormTag(t *testing.T) {
	var s struct {
		Foo customUnmarshalParamHex `form:"foo"`
	}
	err := mappingByPtr(&s, formSource{"foo": {`f5`}}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 245, s.Foo)
}

func TestMappingCustomUnmarshalParamHexWithURITag(t *testing.T) {
	var s struct {
		Foo customUnmarshalParamHex `uri:"foo"`
	}
	err := mappingByPtr(&s, formSource{"foo": {`f5`}}, "uri")
	require.NoError(t, err)

	assert.EqualValues(t, 245, s.Foo)
}

type customUnmarshalParamType struct {
	Protocol string
	Path     string
	Name     string
}

func (f *customUnmarshalParamType) UnmarshalParam(param string) error {
	parts := strings.Split(param, ":")
	if len(parts) != 3 {
		return errors.New("invalid format")
	}
	f.Protocol = parts[0]
	f.Path = parts[1]
	f.Name = parts[2]
	return nil
}

func TestMappingCustomStructTypeWithFormTag(t *testing.T) {
	var s struct {
		FileData customUnmarshalParamType `form:"data"`
	}
	err := mappingByPtr(&s, formSource{"data": {`file:/foo:happiness`}}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, "file", s.FileData.Protocol)
	assert.EqualValues(t, "/foo", s.FileData.Path)
	assert.EqualValues(t, "happiness", s.FileData.Name)
}

func TestMappingCustomStructTypeWithURITag(t *testing.T) {
	var s struct {
		FileData customUnmarshalParamType `uri:"data"`
	}
	err := mappingByPtr(&s, formSource{"data": {`file:/foo:happiness`}}, "uri")
	require.NoError(t, err)

	assert.EqualValues(t, "file", s.FileData.Protocol)
	assert.EqualValues(t, "/foo", s.FileData.Path)
	assert.EqualValues(t, "happiness", s.FileData.Name)
}

func TestMappingCustomPointerStructTypeWithFormTag(t *testing.T) {
	var s struct {
		FileData *customUnmarshalParamType `form:"data"`
	}
	err := mappingByPtr(&s, formSource{"data": {`file:/foo:happiness`}}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, "file", s.FileData.Protocol)
	assert.EqualValues(t, "/foo", s.FileData.Path)
	assert.EqualValues(t, "happiness", s.FileData.Name)
}

func TestMappingCustomPointerStructTypeWithURITag(t *testing.T) {
	var s struct {
		FileData *customUnmarshalParamType `uri:"data"`
	}
	err := mappingByPtr(&s, formSource{"data": {`file:/foo:happiness`}}, "uri")
	require.NoError(t, err)

	assert.EqualValues(t, "file", s.FileData.Protocol)
	assert.EqualValues(t, "/foo", s.FileData.Path)
	assert.EqualValues(t, "happiness", s.FileData.Name)
}

type customPath []string

func (p *customPath) UnmarshalParam(param string) error {
	elems := strings.Split(param, "/")
	n := len(elems)
	if n < 2 {
		return errors.New("invalid format")
	}

	*p = elems
	return nil
}

func TestMappingCustomSliceUri(t *testing.T) {
	var s struct {
		FileData customPath `uri:"path"`
	}
	err := mappingByPtr(&s, formSource{"path": {`bar/foo`}}, "uri")
	require.NoError(t, err)

	assert.EqualValues(t, "bar", s.FileData[0])
	assert.EqualValues(t, "foo", s.FileData[1])
}

func TestMappingCustomSliceForm(t *testing.T) {
	var s struct {
		FileData customPath `form:"path"`
	}
	err := mappingByPtr(&s, formSource{"path": {`bar/foo`}}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, "bar", s.FileData[0])
	assert.EqualValues(t, "foo", s.FileData[1])
}

type objectID [12]byte

func (o *objectID) UnmarshalParam(param string) error {
	oid, err := convertTo(param)
	if err != nil {
		return err
	}

	*o = oid
	return nil
}

func convertTo(s string) (objectID, error) {
	var nilObjectID objectID
	if len(s) != 24 {
		return nilObjectID, errors.New("invalid format")
	}

	var oid [12]byte
	_, err := hex.Decode(oid[:], []byte(s))
	if err != nil {
		return nilObjectID, err
	}

	return oid, nil
}

func TestMappingCustomArrayUri(t *testing.T) {
	var s struct {
		FileData objectID `uri:"id"`
	}
	val := `664a062ac74a8ad104e0e80f`
	err := mappingByPtr(&s, formSource{"id": {val}}, "uri")
	require.NoError(t, err)

	expected, _ := convertTo(val)
	assert.EqualValues(t, expected, s.FileData)
}

func TestMappingCustomArrayForm(t *testing.T) {
	var s struct {
		FileData objectID `form:"id"`
	}
	val := `664a062ac74a8ad104e0e80f`
	err := mappingByPtr(&s, formSource{"id": {val}}, "form")
	require.NoError(t, err)

	expected, _ := convertTo(val)
	assert.EqualValues(t, expected, s.FileData)
}
