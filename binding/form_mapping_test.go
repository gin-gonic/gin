// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding"
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

// ====  BindUmarshaler tests START ====

type customHexUnmarshalParam int

func (f *customHexUnmarshalParam) UnmarshalParam(param string) error {
	v, err := strconv.ParseInt(param, 16, 64)
	if err != nil {
		return err
	}
	*f = customHexUnmarshalParam(v)
	return nil
}

func TestMappingCustomHexUnmarshalParam(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[customHexUnmarshalParam](
		t,
		`f5`,
		func(hex customHexUnmarshalParam, t *testing.T) {
			assert.EqualValues(t, 245, hex)
		},
	)

	// verify default binding works with UnmarshalParam
	var sDefaultValue struct {
		Field1 customHexUnmarshalParam `form:"field1,default=f5"`
	}
	err := mappingByPtr(&sDefaultValue, formSource{"field1": {}}, "form")
	require.NoError(t, err)
	assert.EqualValues(t, 0xf5, sDefaultValue.Field1)
}

type customTypeUnmarshalParam struct {
	Protocol string
	Path     string
	Name     string
}

func (f *customTypeUnmarshalParam) UnmarshalParam(param string) error {
	parts := strings.Split(param, ":")
	if len(parts) != 3 {
		return errors.New("invalid format")
	}
	f.Protocol = parts[0]
	f.Path = parts[1]
	f.Name = parts[2]
	return nil
}

func TestMappingCustomStructType(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[customTypeUnmarshalParam](
		t,
		`file:/foo:happiness`,
		func(data customTypeUnmarshalParam, t *testing.T) {
			assert.Equal(t, "file", data.Protocol)
			assert.Equal(t, "/foo", data.Path)
			assert.Equal(t, "happiness", data.Name)
		},
	)
}

func TestMappingCustomPointerStructType(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[*customTypeUnmarshalParam](
		t,
		`file:/foo:happiness`,
		func(data *customTypeUnmarshalParam, t *testing.T) {
			assert.Equal(t, "file", data.Protocol)
			assert.Equal(t, "/foo", data.Path)
			assert.Equal(t, "happiness", data.Name)
		},
	)
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

func TestMappingCustomSlice(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[customPath](
		t,
		`bar/foo`,
		func(path customPath, t *testing.T) {
			assert.Equal(t, "bar", path[0])
			assert.Equal(t, "foo", path[1])
		},
	)
}

func TestMappingCustomSliceStopsWhenError(t *testing.T) {
	var sForm struct {
		Field1 customPath `form:"field1"`
	}
	err := mappingByPtr(&sForm, formSource{"field1": {"invalid"}}, "form")
	require.ErrorContains(t, err, "invalid format")
	require.Empty(t, sForm.Field1)
}

func TestMappingCustomSliceOfSlice(t *testing.T) {
	val := `bar/foo,bar/foo/spam`
	expected := []customPath{{"bar", "foo"}, {"bar", "foo", "spam"}}

	var sUri struct {
		Field1 []customPath `uri:"field1" collection_format:"csv"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {val}}, "uri")
	require.NoError(t, err)
	assert.Equal(t, expected, sUri.Field1)

	var sForm struct {
		Field1 []customPath `form:"field1" collection_format:"csv"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {val}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sForm.Field1)
}

type objectID [12]byte

func (o *objectID) UnmarshalParam(param string) error {
	oid, err := convertTo[objectID](param)
	if err != nil {
		return err
	}

	*o = oid
	return nil
}

func convertTo[T ~[12]byte](s string) (T, error) {
	var nilObjectID T
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

func TestMappingCustomArray(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[objectID](
		t,
		`664a062ac74a8ad104e0e80f`,
		func(oid objectID, t *testing.T) {
			expected, _ := convertTo[objectID](`664a062ac74a8ad104e0e80f`)
			assert.Equal(t, expected, oid)
		},
	)
}

func TestMappingCustomArrayOfArray(t *testing.T) {
	val := `664a062ac74a8ad104e0e80e,664a062ac74a8ad104e0e80f`
	expected1, _ := convertTo[objectID](`664a062ac74a8ad104e0e80e`)
	expected2, _ := convertTo[objectID](`664a062ac74a8ad104e0e80f`)
	expected := []objectID{expected1, expected2}

	var sUri struct {
		Field1 []objectID `uri:"field1" collection_format:"csv"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {val}}, "uri")
	require.NoError(t, err)
	assert.Equal(t, expected, sUri.Field1)

	var sForm struct {
		Field1 []objectID `form:"field1" collection_format:"csv"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {val}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sForm.Field1)

	var sDefaultValue struct {
		Field1 []objectID `form:"field1,default=664a062ac74a8ad104e0e80e;664a062ac74a8ad104e0e80f" collection_format:"csv"`
	}
	err = mappingByPtr(&sDefaultValue, formSource{"field1": {}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sDefaultValue.Field1)
}

// RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam declares a struct with a field of the given generic type T
// and runs a mapping test using the given value for both the uri and form tag. Any asserts that should be done on the
// result are passed as a function in the last parameter.
//
// This method eliminates the need for writing duplicate tests to verify both form+uri tags for BindUnmarshaler tests
func RunMappingUsingUriAndFormTagAndAssertForUnmarshalParam[T any](
	t *testing.T,
	valueToBind string,
	assertsToRunAfterBind func(T, *testing.T),
) {
	var sUri struct {
		Field1 T `uri:"field1"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {valueToBind}}, "uri")
	require.NoError(t, err)
	assertsToRunAfterBind(sUri.Field1, t)

	var sForm struct {
		Field1 T `form:"field1"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {valueToBind}}, "form")
	require.NoError(t, err)
	assertsToRunAfterBind(sForm.Field1, t)
}

// ====  BindUmarshaler tests END ====

// ====  TextUnmarshaler tests START ====

type customHexUnmarshalText int

func (f *customHexUnmarshalText) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 16, 64)
	if err != nil {
		return err
	}
	*f = customHexUnmarshalText(v)
	return nil
}

// verify type implements TextUnmarshaler
var _ encoding.TextUnmarshaler = (*customHexUnmarshalText)(nil)

func TestMappingCustomHexUnmarshalText(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[customHexUnmarshalText](
		t,
		`f5`,
		func(hex customHexUnmarshalText, t *testing.T) {
			assert.EqualValues(t, 245, hex)
		},
	)

	// verify default binding works with UnmarshalText
	var sDefaultValue struct {
		Field1 customHexUnmarshalText `form:"field1,default=f5,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&sDefaultValue, formSource{"field1": {}}, "form")
	require.NoError(t, err)
	assert.EqualValues(t, 0xf5, sDefaultValue.Field1)
}

type customTypeUnmarshalText struct {
	Protocol string
	Path     string
	Name     string
}

func (f *customTypeUnmarshalText) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), ":")
	if len(parts) != 3 {
		return errors.New("invalid format")
	}
	f.Protocol = parts[0]
	f.Path = parts[1]
	f.Name = parts[2]
	return nil
}

var _ encoding.TextUnmarshaler = (*customTypeUnmarshalText)(nil)

func TestMappingCustomStructTypeUnmarshalText(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[customTypeUnmarshalText](
		t,
		`file:/foo:happiness`,
		func(data customTypeUnmarshalText, t *testing.T) {
			assert.Equal(t, "file", data.Protocol)
			assert.Equal(t, "/foo", data.Path)
			assert.Equal(t, "happiness", data.Name)
		},
	)
}

func TestMappingCustomPointerStructTypeUnmarshalText(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[*customTypeUnmarshalText](
		t,
		`file:/foo:happiness`,
		func(data *customTypeUnmarshalText, t *testing.T) {
			assert.Equal(t, "file", data.Protocol)
			assert.Equal(t, "/foo", data.Path)
			assert.Equal(t, "happiness", data.Name)
		},
	)
}

type customPathUnmarshalText []string

func (p *customPathUnmarshalText) UnmarshalText(text []byte) error {
	elems := strings.Split(string(text), "/")
	n := len(elems)
	if n < 2 {
		return errors.New("invalid format")
	}

	*p = elems
	return nil
}

var _ encoding.TextUnmarshaler = (*customPathUnmarshalText)(nil)

func TestMappingCustomSliceUnmarshalText(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[customPathUnmarshalText](
		t,
		`bar/foo`,
		func(path customPathUnmarshalText, t *testing.T) {
			assert.Equal(t, "bar", path[0])
			assert.Equal(t, "foo", path[1])
		},
	)
}

func TestMappingCustomSliceUnmarshalTextStopsWhenError(t *testing.T) {
	var sForm struct {
		Field1 customPathUnmarshalText `form:"field1,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&sForm, formSource{"field1": {"invalid"}}, "form")
	require.ErrorContains(t, err, "invalid format")
	require.Empty(t, sForm.Field1)
}

func TestMappingCustomSliceOfSliceUnmarshalText(t *testing.T) {
	val := `bar/foo,bar/foo/spam`
	expected := []customPathUnmarshalText{{"bar", "foo"}, {"bar", "foo", "spam"}}

	var sUri struct {
		Field1 []customPathUnmarshalText `uri:"field1,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {val}}, "uri")
	require.NoError(t, err)
	assert.Equal(t, expected, sUri.Field1)

	var sForm struct {
		Field1 []customPathUnmarshalText `form:"field1,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {val}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sForm.Field1)

	var sDefaultValue struct {
		Field1 []customPathUnmarshalText `form:"field1,default=bar/foo;bar/foo/spam,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err = mappingByPtr(&sDefaultValue, formSource{"field1": {}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sDefaultValue.Field1)
}

type objectIDUnmarshalText [12]byte

func (o *objectIDUnmarshalText) UnmarshalText(text []byte) error {
	oid, err := convertTo[objectIDUnmarshalText](string(text))
	if err != nil {
		return err
	}

	*o = oid
	return nil
}

var _ encoding.TextUnmarshaler = (*objectIDUnmarshalText)(nil)

func TestMappingCustomArrayUnmarshalText(t *testing.T) {
	RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[objectIDUnmarshalText](
		t,
		`664a062ac74a8ad104e0e80f`,
		func(oid objectIDUnmarshalText, t *testing.T) {
			expected, _ := convertTo[objectIDUnmarshalText](`664a062ac74a8ad104e0e80f`)
			assert.Equal(t, expected, oid)
		},
	)
}

func TestMappingCustomArrayOfArrayUnmarshalText(t *testing.T) {
	val := `664a062ac74a8ad104e0e80e,664a062ac74a8ad104e0e80f`
	expected1, _ := convertTo[objectIDUnmarshalText](`664a062ac74a8ad104e0e80e`)
	expected2, _ := convertTo[objectIDUnmarshalText](`664a062ac74a8ad104e0e80f`)
	expected := []objectIDUnmarshalText{expected1, expected2}

	var sUri struct {
		Field1 []objectIDUnmarshalText `uri:"field1,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {val}}, "uri")
	require.NoError(t, err)
	assert.Equal(t, expected, sUri.Field1)

	var sForm struct {
		Field1 []objectIDUnmarshalText `form:"field1,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {val}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sForm.Field1)

	var sDefaultValue struct {
		Field1 []objectIDUnmarshalText `form:"field1,default=664a062ac74a8ad104e0e80e;664a062ac74a8ad104e0e80f,parser=encoding.TextUnmarshaler" collection_format:"csv"`
	}
	err = mappingByPtr(&sDefaultValue, formSource{"field1": {}}, "form")
	require.NoError(t, err)
	assert.Equal(t, expected, sDefaultValue.Field1)
}

// RunMappingUsingUriAndFormTagAndAssertForUnmarshalText declares a struct with a field of the given generic type T
// and runs a mapping test using the given value for both the uri and form tag. Any asserts that should be done on the
// result are passed as a function in the last parameter.
//
// This method eliminates the need for writing duplicate tests to verify both form+uri tags for TextUnmarshaler tests
func RunMappingUsingUriAndFormTagAndAssertForUnmarshalText[T any](
	t *testing.T,
	valueToBind string,
	assertsToRunAfterBind func(T, *testing.T),
) {
	var sUri struct {
		Field1 T `uri:"field1,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&sUri, formSource{"field1": {valueToBind}}, "uri")
	require.NoError(t, err)
	assertsToRunAfterBind(sUri.Field1, t)

	var sForm struct {
		Field1 T `form:"field1,parser=encoding.TextUnmarshaler"`
	}
	err = mappingByPtr(&sForm, formSource{"field1": {valueToBind}}, "form")
	require.NoError(t, err)
	assertsToRunAfterBind(sForm.Field1, t)
}

// If someone specifies parser=TextUnmarshaler and it's not defined for the type, gin should revert to using its default
// binding logic.
func TestMappingUsingBindUnmarshalerAndTextUnmarshalerWhenOnlyBindUnmarshalerDefined(t *testing.T) {
	var s struct {
		Hex                customHexUnmarshalParam `form:"hex"`
		HexByUnmarshalText customHexUnmarshalParam `form:"hex2,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&s, formSource{
		"hex":  {`f5`},
		"hex2": {`f5`},
	}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 0xf5, s.Hex)
	assert.EqualValues(t, 0xf5, s.HexByUnmarshalText) // reverts to BindUnmarshaler binding
}

// If someone does not specify parser=TextUnmarshaler even when it's defined for the type, gin should ignore the
// UnmarshalText logic and continue using its default binding logic. (This ensures gin does not break backwards
// compatibility)
func TestMappingUsingBindUnmarshalerAndTextUnmarshalerWhenOnlyTextUnmarshalerDefined(t *testing.T) {
	var s struct {
		Hex                customHexUnmarshalText `form:"hex"`
		HexByUnmarshalText customHexUnmarshalText `form:"hex2,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&s, formSource{
		"hex":  {`11`},
		"hex2": {`11`},
	}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 11, s.Hex)                  // this is using default int binding, not our custom hex binding. 0x11 should be 17 in decimal
	assert.EqualValues(t, 0x11, s.HexByUnmarshalText) // correct expected value for hex binding
}

type customHexUnmarshalParamAndUnmarshalText int

func (f *customHexUnmarshalParamAndUnmarshalText) UnmarshalParam(param string) error {
	return errors.New("should not be called in unit test if parser tag present")
}

func (f *customHexUnmarshalParamAndUnmarshalText) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 16, 64)
	if err != nil {
		return err
	}
	*f = customHexUnmarshalParamAndUnmarshalText(v)
	return nil
}

// If a type has both UnmarshalParam and UnmarshalText methods defined, but the parser tag is set to TextUnmarshaler,
// then only the UnmarshalText method should be invoked.
func TestMappingUsingTextUnmarshalerWhenBindUnmarshalerAlsoDefined(t *testing.T) {
	var s struct {
		Hex customHexUnmarshalParamAndUnmarshalText `form:"hex,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&s, formSource{
		"hex": {`f5`},
	}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 0xf5, s.Hex)
}

// ====  TextUnmarshaler tests END ====
