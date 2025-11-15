package encoder

import (
	"context"
	"fmt"
	"reflect"
)

var (
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
)

type FieldQuery struct {
	Name   string
	Fields []*FieldQuery
	hash   string
}

func (q *FieldQuery) Hash() string {
	if q.hash != "" {
		return q.hash
	}
	b, _ := Marshal(q)
	q.hash = string(b)
	return q.hash
}

func (q *FieldQuery) MarshalJSON() ([]byte, error) {
	if q.Name != "" {
		if len(q.Fields) > 0 {
			return Marshal(map[string][]*FieldQuery{q.Name: q.Fields})
		}
		return Marshal(q.Name)
	}
	return Marshal(q.Fields)
}

func (q *FieldQuery) QueryString() (FieldQueryString, error) {
	b, err := Marshal(q)
	if err != nil {
		return "", err
	}
	return FieldQueryString(b), nil
}

type FieldQueryString string

func (s FieldQueryString) Build() (*FieldQuery, error) {
	var query interface{}
	if err := Unmarshal([]byte(s), &query); err != nil {
		return nil, err
	}
	return s.build(reflect.ValueOf(query))
}

func (s FieldQueryString) build(v reflect.Value) (*FieldQuery, error) {
	switch v.Type().Kind() {
	case reflect.String:
		return s.buildString(v)
	case reflect.Map:
		return s.buildMap(v)
	case reflect.Slice:
		return s.buildSlice(v)
	case reflect.Interface:
		return s.build(reflect.ValueOf(v.Interface()))
	}
	return nil, fmt.Errorf("failed to build field query")
}

func (s FieldQueryString) buildString(v reflect.Value) (*FieldQuery, error) {
	b := []byte(v.String())
	switch b[0] {
	case '[', '{':
		var query interface{}
		if err := Unmarshal(b, &query); err != nil {
			return nil, err
		}
		if str, ok := query.(string); ok {
			return &FieldQuery{Name: str}, nil
		}
		return s.build(reflect.ValueOf(query))
	}
	return &FieldQuery{Name: string(b)}, nil
}

func (s FieldQueryString) buildSlice(v reflect.Value) (*FieldQuery, error) {
	fields := make([]*FieldQuery, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		def, err := s.build(v.Index(i))
		if err != nil {
			return nil, err
		}
		fields = append(fields, def)
	}
	return &FieldQuery{Fields: fields}, nil
}

func (s FieldQueryString) buildMap(v reflect.Value) (*FieldQuery, error) {
	keys := v.MapKeys()
	if len(keys) != 1 {
		return nil, fmt.Errorf("failed to build field query object")
	}
	key := keys[0]
	if key.Type().Kind() != reflect.String {
		return nil, fmt.Errorf("failed to build field query. invalid object key type")
	}
	name := key.String()
	def, err := s.build(v.MapIndex(key))
	if err != nil {
		return nil, err
	}
	return &FieldQuery{
		Name:   name,
		Fields: def.Fields,
	}, nil
}

type queryKey struct{}

func FieldQueryFromContext(ctx context.Context) *FieldQuery {
	query := ctx.Value(queryKey{})
	if query == nil {
		return nil
	}
	q, ok := query.(*FieldQuery)
	if !ok {
		return nil
	}
	return q
}

func SetFieldQueryToContext(ctx context.Context, query *FieldQuery) context.Context {
	return context.WithValue(ctx, queryKey{}, query)
}
