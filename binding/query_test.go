package binding

import (
	"reflect"
	"testing"
	"time"
)

type typeTest struct {
	needParse []string
	ptr       interface{}
	want      interface{}
}

func TestParseTypeVar(t *testing.T) {
	var (
		b            bool
		i            int
		i8           int8
		i16          int16
		i32          int32
		i64          int64
		u            uint
		u8           uint8
		u16          uint16
		u32          uint32
		u64          uint64
		s            string
		f32          float32
		f64          float64
		duration     time.Duration
		stringSlice  []string
		intSlice     []int
		float32Slice []float32
		stringArray  [3]string
		tm           time.Time
	)

	tv := []typeTest{
		{needParse: []string{"1"}, ptr: &i, want: 1},
		{needParse: []string{"2"}, ptr: &i8, want: int8(2)},
		{needParse: []string{"3"}, ptr: &i16, want: int16(3)},
		{needParse: []string{"4"}, ptr: &i32, want: int32(4)},
		{needParse: []string{"5"}, ptr: &i64, want: int64(5)},
		{needParse: []string{"6"}, ptr: &u, want: uint(6)},
		{needParse: []string{"7"}, ptr: &u8, want: uint8(7)},
		{needParse: []string{"8"}, ptr: &u16, want: uint16(8)},
		{needParse: []string{"9"}, ptr: &u32, want: uint32(9)},
		{needParse: []string{"10"}, ptr: &u64, want: uint64(10)},
		{needParse: []string{"test"}, ptr: &s, want: "test"},
		{needParse: []string{"1.1"}, ptr: &f32, want: float32(1.1)},
		{needParse: []string{"2.2"}, ptr: &f64, want: float64(2.2)},
		{needParse: []string{"1", "2", "3"}, ptr: &stringSlice, want: []string{"1", "2", "3"}},
		{needParse: []string{"1", "2", "3"}, ptr: &intSlice, want: []int{1, 2, 3}},
		{needParse: []string{"4.1", "5.1", "6.1"}, ptr: &float32Slice, want: []float32{4.1, 5.1, 6.1}},
		{needParse: []string{"a1", "a2", "a3"}, ptr: &stringArray, want: [3]string{"a1", "a2", "a3"}},
		{needParse: []string{"true"}, ptr: &b, want: true},
		{needParse: []string{"1s"}, ptr: &duration, want: time.Second},
		{needParse: []string{"2006-01-02T15:04:05Z"}, ptr: &tm, want: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)},
	}

	for k := range tv {
		if err := parseTypeVar(reflect.ValueOf(tv[k].ptr), tv[k].needParse); err != nil {
			t.Errorf("parseBaseTypeVar %T fail:%s\n", tv[k].want, err)
		}

		if !reflect.DeepEqual(reflect.ValueOf(tv[k].ptr).Elem().Interface(), tv[k].want) {
			t.Errorf("parseBaseTypeVar %T fail got:%v, want:%v\n", tv[k].ptr, tv[k].ptr, tv[k].want)
		}
	}

}
