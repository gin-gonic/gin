// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	} {
		tp := reflect.TypeOf(tt.value)
		testName := tt.name + ":" + tp.Field(0).Type.String()

		val := reflect.New(reflect.TypeOf(tt.value))
		val.Elem().Set(reflect.ValueOf(tt.value))

		field := val.Elem().Type().Field(0)

		_, err := mapping(val, emptyField, formSource{field.Name: {tt.form}}, "form")
		assert.NoError(t, err, testName)

		actual := val.Elem().Field(0).Interface()
		assert.Equal(t, tt.expect, actual, testName)
	}
}

func TestMappingDefault(t *testing.T) {
	var s struct {
		Int   int    `form:",default=9"`
		Slice []int  `form:",default=9"`
		Array [1]int `form:",default=9"`
	}
	err := mappingByPtr(&s, formSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.Int)
	assert.Equal(t, []int{9}, s.Slice)
	assert.Equal(t, [1]int{9}, s.Array)
}

func TestMappingSkipField(t *testing.T) {
	var s struct {
		A int
	}
	err := mappingByPtr(&s, formSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 0, s.A)
}

func TestMappingIgnoreField(t *testing.T) {
	var s struct {
		A int `form:"A"`
		B int `form:"-"`
	}
	err := mappingByPtr(&s, formSource{"A": {"9"}, "B": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.B)
}

func TestMappingUnexportedField(t *testing.T) {
	var s struct {
		A int `form:"a"`
		b int `form:"b"`
	}
	err := mappingByPtr(&s, formSource{"a": {"9"}, "b": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.b)
}

func TestMappingPrivateField(t *testing.T) {
	var s struct {
		f int `form:"field"`
	}
	err := mappingByPtr(&s, formSource{"field": {"6"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, int(0), s.f)
}

func TestMappingUnknownFieldType(t *testing.T) {
	var s struct {
		U uintptr
	}

	err := mappingByPtr(&s, formSource{"U": {"unknown"}}, "form")
	assert.Error(t, err)
	assert.Equal(t, errUnknownType, err)
}

func TestMappingURI(t *testing.T) {
	var s struct {
		F int `uri:"field"`
	}
	err := mapURI(&s, map[string][]string{"field": {"6"}})
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
}

func TestMappingForm(t *testing.T) {
	var s struct {
		F int `form:"field"`
	}
	err := mapForm(&s, map[string][]string{"field": {"6"}})
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
}

func TestMapFormWithTag(t *testing.T) {
	var s struct {
		F int `externalTag:"field"`
	}
	err := MapFormWithTag(&s, map[string][]string{"field": {"6"}}, "externalTag")
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
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
	assert.NoError(t, err)

	err = mapForm(&s, map[string][]string{
		"Time":      {"2019-01-20T16:02:58Z"},
		"LocalTime": {"2019-01-20"},
		"ZeroValue": {},
		"CSTTime":   {"2019-01-20"},
		"UTCTime":   {"2019-01-20"},
	})
	assert.NoError(t, err)

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
	assert.Error(t, err)

	// wrong time value
	var wrongTime struct {
		Time time.Time
	}
	err = mapForm(&wrongTime, map[string][]string{"Time": {"wrong"}})
	assert.Error(t, err)
}

func TestMappingTimeDuration(t *testing.T) {
	var s struct {
		D time.Duration
	}

	// ok
	err := mappingByPtr(&s, formSource{"D": {"5s"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 5*time.Second, s.D)

	// error
	err = mappingByPtr(&s, formSource{"D": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingSlice(t *testing.T) {
	var s struct {
		Slice []int `form:"slice,default=9"`
	}

	// default value
	err := mappingByPtr(&s, formSource{}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{9}, s.Slice)

	// ok
	err = mappingByPtr(&s, formSource{"slice": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4}, s.Slice)

	// error
	err = mappingByPtr(&s, formSource{"slice": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingArray(t *testing.T) {
	var s struct {
		Array [2]int `form:"array,default=9"`
	}

	// wrong default
	err := mappingByPtr(&s, formSource{}, "form")
	assert.Error(t, err)

	// ok
	err = mappingByPtr(&s, formSource{"array": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, [2]int{3, 4}, s.Array)

	// error - not enough vals
	err = mappingByPtr(&s, formSource{"array": {"3"}}, "form")
	assert.Error(t, err)

	// error - wrong value
	err = mappingByPtr(&s, formSource{"array": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingStructField(t *testing.T) {
	var s struct {
		J struct {
			I int
		}
	}

	err := mappingByPtr(&s, formSource{"J": {`{"I": 9}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 9, s.J.I)
}

func TestMappingMapField(t *testing.T) {
	var s struct {
		M map[string]int
	}

	err := mappingByPtr(&s, formSource{"M": {`{"one": 1}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"one": 1}, s.M)
}

func TestMappingIgnoredCircularRef(t *testing.T) {
	type S struct {
		S *S `form:"-"`
	}
	var s S

	err := mappingByPtr(&s, formSource{}, "form")
	assert.NoError(t, err)
}
