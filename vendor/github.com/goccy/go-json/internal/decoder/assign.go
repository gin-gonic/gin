package decoder

import (
	"fmt"
	"reflect"
	"strconv"
)

var (
	nilValue = reflect.ValueOf(nil)
)

func AssignValue(src, dst reflect.Value) error {
	if dst.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("invalid dst type. required pointer type: %T", dst.Type())
	}
	casted, err := castValue(dst.Elem().Type(), src)
	if err != nil {
		return err
	}
	dst.Elem().Set(casted)
	return nil
}

func castValue(t reflect.Type, v reflect.Value) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Int:
		vv, err := castInt(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(int(vv.Int())), nil
	case reflect.Int8:
		vv, err := castInt(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(int8(vv.Int())), nil
	case reflect.Int16:
		vv, err := castInt(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(int16(vv.Int())), nil
	case reflect.Int32:
		vv, err := castInt(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(int32(vv.Int())), nil
	case reflect.Int64:
		return castInt(v)
	case reflect.Uint:
		vv, err := castUint(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(uint(vv.Uint())), nil
	case reflect.Uint8:
		vv, err := castUint(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(uint8(vv.Uint())), nil
	case reflect.Uint16:
		vv, err := castUint(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(uint16(vv.Uint())), nil
	case reflect.Uint32:
		vv, err := castUint(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(uint32(vv.Uint())), nil
	case reflect.Uint64:
		return castUint(v)
	case reflect.Uintptr:
		vv, err := castUint(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(uintptr(vv.Uint())), nil
	case reflect.String:
		return castString(v)
	case reflect.Bool:
		return castBool(v)
	case reflect.Float32:
		vv, err := castFloat(v)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(float32(vv.Float())), nil
	case reflect.Float64:
		return castFloat(v)
	case reflect.Array:
		return castArray(t, v)
	case reflect.Slice:
		return castSlice(t, v)
	case reflect.Map:
		return castMap(t, v)
	case reflect.Struct:
		return castStruct(t, v)
	}
	return v, nil
}

func castInt(v reflect.Value) (reflect.Value, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflect.ValueOf(int64(v.Uint())), nil
	case reflect.String:
		i64, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(i64), nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf(int64(1)), nil
		}
		return reflect.ValueOf(int64(0)), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(int64(v.Float())), nil
	case reflect.Array:
		if v.Len() > 0 {
			return castInt(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to int64 from empty array")
	case reflect.Slice:
		if v.Len() > 0 {
			return castInt(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to int64 from empty slice")
	case reflect.Interface:
		return castInt(reflect.ValueOf(v.Interface()))
	case reflect.Map:
		return nilValue, fmt.Errorf("failed to cast to int64 from map")
	case reflect.Struct:
		return nilValue, fmt.Errorf("failed to cast to int64 from struct")
	case reflect.Ptr:
		return castInt(v.Elem())
	}
	return nilValue, fmt.Errorf("failed to cast to int64 from %s", v.Type().Kind())
}

func castUint(v reflect.Value) (reflect.Value, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(uint64(v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v, nil
	case reflect.String:
		u64, err := strconv.ParseUint(v.String(), 10, 64)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(u64), nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf(uint64(1)), nil
		}
		return reflect.ValueOf(uint64(0)), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(uint64(v.Float())), nil
	case reflect.Array:
		if v.Len() > 0 {
			return castUint(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to uint64 from empty array")
	case reflect.Slice:
		if v.Len() > 0 {
			return castUint(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to uint64 from empty slice")
	case reflect.Interface:
		return castUint(reflect.ValueOf(v.Interface()))
	case reflect.Map:
		return nilValue, fmt.Errorf("failed to cast to uint64 from map")
	case reflect.Struct:
		return nilValue, fmt.Errorf("failed to cast to uint64 from struct")
	case reflect.Ptr:
		return castUint(v.Elem())
	}
	return nilValue, fmt.Errorf("failed to cast to uint64 from %s", v.Type().Kind())
}

func castString(v reflect.Value) (reflect.Value, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(fmt.Sprint(v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflect.ValueOf(fmt.Sprint(v.Uint())), nil
	case reflect.String:
		return v, nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf("true"), nil
		}
		return reflect.ValueOf("false"), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(fmt.Sprint(v.Float())), nil
	case reflect.Array:
		if v.Len() > 0 {
			return castString(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to string from empty array")
	case reflect.Slice:
		if v.Len() > 0 {
			return castString(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to string from empty slice")
	case reflect.Interface:
		return castString(reflect.ValueOf(v.Interface()))
	case reflect.Map:
		return nilValue, fmt.Errorf("failed to cast to string from map")
	case reflect.Struct:
		return nilValue, fmt.Errorf("failed to cast to string from struct")
	case reflect.Ptr:
		return castString(v.Elem())
	}
	return nilValue, fmt.Errorf("failed to cast to string from %s", v.Type().Kind())
}

func castBool(v reflect.Value) (reflect.Value, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Int() {
		case 0:
			return reflect.ValueOf(false), nil
		case 1:
			return reflect.ValueOf(true), nil
		}
		return nilValue, fmt.Errorf("failed to cast to bool from %d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Uint() {
		case 0:
			return reflect.ValueOf(false), nil
		case 1:
			return reflect.ValueOf(true), nil
		}
		return nilValue, fmt.Errorf("failed to cast to bool from %d", v.Uint())
	case reflect.String:
		b, err := strconv.ParseBool(v.String())
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(b), nil
	case reflect.Bool:
		return v, nil
	case reflect.Float32, reflect.Float64:
		switch v.Float() {
		case 0:
			return reflect.ValueOf(false), nil
		case 1:
			return reflect.ValueOf(true), nil
		}
		return nilValue, fmt.Errorf("failed to cast to bool from %f", v.Float())
	case reflect.Array:
		if v.Len() > 0 {
			return castBool(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to string from empty array")
	case reflect.Slice:
		if v.Len() > 0 {
			return castBool(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to string from empty slice")
	case reflect.Interface:
		return castBool(reflect.ValueOf(v.Interface()))
	case reflect.Map:
		return nilValue, fmt.Errorf("failed to cast to string from map")
	case reflect.Struct:
		return nilValue, fmt.Errorf("failed to cast to string from struct")
	case reflect.Ptr:
		return castBool(v.Elem())
	}
	return nilValue, fmt.Errorf("failed to cast to bool from %s", v.Type().Kind())
}

func castFloat(v reflect.Value) (reflect.Value, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(float64(v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflect.ValueOf(float64(v.Uint())), nil
	case reflect.String:
		f64, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(f64), nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf(float64(1)), nil
		}
		return reflect.ValueOf(float64(0)), nil
	case reflect.Float32, reflect.Float64:
		return v, nil
	case reflect.Array:
		if v.Len() > 0 {
			return castFloat(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to float64 from empty array")
	case reflect.Slice:
		if v.Len() > 0 {
			return castFloat(v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to float64 from empty slice")
	case reflect.Interface:
		return castFloat(reflect.ValueOf(v.Interface()))
	case reflect.Map:
		return nilValue, fmt.Errorf("failed to cast to float64 from map")
	case reflect.Struct:
		return nilValue, fmt.Errorf("failed to cast to float64 from struct")
	case reflect.Ptr:
		return castFloat(v.Elem())
	}
	return nilValue, fmt.Errorf("failed to cast to float64 from %s", v.Type().Kind())
}

func castArray(t reflect.Type, v reflect.Value) (reflect.Value, error) {
	kind := v.Type().Kind()
	if kind == reflect.Interface {
		return castArray(t, reflect.ValueOf(v.Interface()))
	}
	if kind != reflect.Slice && kind != reflect.Array {
		return nilValue, fmt.Errorf("failed to cast to array from %s", kind)
	}
	if t.Elem() == v.Type().Elem() {
		return v, nil
	}
	if t.Len() != v.Len() {
		return nilValue, fmt.Errorf("failed to cast [%d]array from slice of %d length", t.Len(), v.Len())
	}
	ret := reflect.New(t).Elem()
	for i := 0; i < v.Len(); i++ {
		vv, err := castValue(t.Elem(), v.Index(i))
		if err != nil {
			return nilValue, err
		}
		ret.Index(i).Set(vv)
	}
	return ret, nil
}

func castSlice(t reflect.Type, v reflect.Value) (reflect.Value, error) {
	kind := v.Type().Kind()
	if kind == reflect.Interface {
		return castSlice(t, reflect.ValueOf(v.Interface()))
	}
	if kind != reflect.Slice && kind != reflect.Array {
		return nilValue, fmt.Errorf("failed to cast to slice from %s", kind)
	}
	if t.Elem() == v.Type().Elem() {
		return v, nil
	}
	ret := reflect.MakeSlice(t, v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		vv, err := castValue(t.Elem(), v.Index(i))
		if err != nil {
			return nilValue, err
		}
		ret.Index(i).Set(vv)
	}
	return ret, nil
}

func castMap(t reflect.Type, v reflect.Value) (reflect.Value, error) {
	ret := reflect.MakeMap(t)
	switch v.Type().Kind() {
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			key, err := castValue(t.Key(), iter.Key())
			if err != nil {
				return nilValue, err
			}
			value, err := castValue(t.Elem(), iter.Value())
			if err != nil {
				return nilValue, err
			}
			ret.SetMapIndex(key, value)
		}
		return ret, nil
	case reflect.Interface:
		return castMap(t, reflect.ValueOf(v.Interface()))
	case reflect.Slice:
		if v.Len() > 0 {
			return castMap(t, v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to map from empty slice")
	}
	return nilValue, fmt.Errorf("failed to cast to map from %s", v.Type().Kind())
}

func castStruct(t reflect.Type, v reflect.Value) (reflect.Value, error) {
	ret := reflect.New(t).Elem()
	switch v.Type().Kind() {
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			k, err := castString(key)
			if err != nil {
				return nilValue, err
			}
			fieldName := k.String()
			field, ok := t.FieldByName(fieldName)
			if ok {
				value, err := castValue(field.Type, iter.Value())
				if err != nil {
					return nilValue, err
				}
				ret.FieldByName(fieldName).Set(value)
			}
		}
		return ret, nil
	case reflect.Struct:
		for i := 0; i < v.Type().NumField(); i++ {
			name := v.Type().Field(i).Name
			ret.FieldByName(name).Set(v.FieldByName(name))
		}
		return ret, nil
	case reflect.Interface:
		return castStruct(t, reflect.ValueOf(v.Interface()))
	case reflect.Slice:
		if v.Len() > 0 {
			return castStruct(t, v.Index(0))
		}
		return nilValue, fmt.Errorf("failed to cast to struct from empty slice")
	default:
		return nilValue, fmt.Errorf("failed to cast to struct from %s", v.Type().Kind())
	}
}
