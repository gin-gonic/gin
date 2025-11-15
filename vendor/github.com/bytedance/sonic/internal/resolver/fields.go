/**
 * Copyright 2025 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resolver

import (
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytedance/sonic/internal/encoder/alg"
)

type StdField struct {
    name        string
    nameBytes   []byte
    nameNonEsc  string
    nameEscHTML string
    tag         bool
    index       []int
    typ         reflect.Type
    omitEmpty   bool
    omitZero    bool
    isZero      func(reflect.Value) bool
    quoted      bool
}

type StdStructFields struct {
    list      []StdField
    nameIndex map[string]*StdField
    byFoldedName map[string]*StdField
}

func typeFields(t reflect.Type) StdStructFields {
	// Anonymous fields to explore at the current level and the next.
	current := []StdField{}
	next := []StdField{{typ: t}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []StdField

	// Buffer to run appendHTMLEscape on field names.
	var nameEscBuf []byte

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Only strings, floats, integers, and booleans can be quoted.
				quoted := false
				if opts.Contains("string") {
					switch ft.Kind() {
					case reflect.Bool,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
						reflect.Float32, reflect.Float64,
						reflect.String:
						quoted = true
					}
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := StdField{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
						omitZero:  opts.Contains("omitzero"),
						quoted:    quoted,
					}
					field.nameBytes = []byte(field.name)

					// Build nameEscHTML and nameNonEsc ahead of time.
					nameEscBuf = alg.HtmlEscape(nameEscBuf[:0], field.nameBytes)
					field.nameEscHTML = `"` + string(nameEscBuf) + `":`
					field.nameNonEsc = `"` + field.name + `":`

					if field.omitZero {
						t := sf.Type
						// Provide a function that uses a type's IsZero method.
						switch {
						case t.Kind() == reflect.Interface && t.Implements(isZeroerType):
							field.isZero = func(v reflect.Value) bool {
								// Avoid panics calling IsZero on a nil interface or
								// non-nil interface with nil pointer.
								return v.IsNil() ||
									(v.Elem().Kind() == reflect.Pointer && v.Elem().IsNil()) ||
									v.Interface().(isZeroer).IsZero()
							}
						case t.Kind() == reflect.Pointer && t.Implements(isZeroerType):
							field.isZero = func(v reflect.Value) bool {
								// Avoid panics calling IsZero on nil pointer.
								return v.IsNil() || v.Interface().(isZeroer).IsZero()
							}
						case t.Implements(isZeroerType):
							field.isZero = func(v reflect.Value) bool {
								return v.Interface().(isZeroer).IsZero()
							}
						case reflect.PointerTo(t).Implements(isZeroerType):
							field.isZero = func(v reflect.Value) bool {
								if !v.CanAddr() {
									// Temporarily box v so we can take the address.
									v2 := reflect.New(v.Type()).Elem()
									v2.Set(v)
									v = v2
								}
								return v.Addr().Interface().(isZeroer).IsZero()
							}
						}
					}

					fields = append(fields, field)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 and 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, StdField{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		a, b := fields[i], fields[j]
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if c := strings.Compare(a.name, b.name); c != 0 {
			return c < 0
		}
		if len(a.index) != len(b.index) {
			return len(a.index) < len(b.index)
		}
		if a.tag != b.tag {
			if a.tag {
				return true
			}
			return false
		}
		return compare(a.index, b.index) < 0
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Slice(fields, func(i, j int) bool {
		a, b := fields[i], fields[j]
		return compare(a.index, b.index) < 0
	})

	exactNameIndex := make(map[string]*StdField, len(fields))
	foldedNameIndex := make(map[string]*StdField, len(fields))
	for i, field := range fields {
		exactNameIndex[field.name] = &fields[i]
		// For historical reasons, first folded match takes precedence.
		if _, ok := foldedNameIndex[string(foldName(field.nameBytes))]; !ok {
			foldedNameIndex[string(foldName(field.nameBytes))] = &fields[i]
		}
	}
	return StdStructFields{fields, exactNameIndex, foldedNameIndex}
}

func compare(s1, s2 []int) int {
	for i, v1 := range s1 {
		if i >= len(s2) {
			return +1
		}
		v2 := s2[i]
		if v1 != v2 {
			return v1 - v2
		}
	}
	if len(s1) < len(s2) {
		return -1
	}
	return 0
}

type isZeroer interface {
	IsZero() bool
}

var isZeroerType = reflect.TypeOf((*isZeroer)(nil)).Elem()

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []StdField) (StdField, bool) {
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level, either both tagged or neither tagged.
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return StdField{}, false
	}
	return fields[0], true
}


// foldName returns a folded string such that foldName(x) == foldName(y)
// is identical to bytes.EqualFold(x, y).
func foldName(in []byte) []byte {
	// This is inlinable to take advantage of "function outlining".
	var arr [32]byte // large enough for most JSON names
	return appendFoldedName(arr[:0], in)
}

func appendFoldedName(out, in []byte) []byte {
	for i := 0; i < len(in); {
		// Handle single-byte ASCII.
		if c := in[i]; c < utf8.RuneSelf {
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			out = append(out, c)
			i++
			continue
		}
		// Handle multi-byte Unicode.
		r, n := utf8.DecodeRune(in[i:])
		out = utf8.AppendRune(out, foldRune(r))
		i += n
	}
	return out
}

// foldRune is returns the smallest rune for all runes in the same fold set.
func foldRune(r rune) rune {
	for {
		r2 := unicode.SimpleFold(r)
		if r2 <= r {
			return r2
		}
		r = r2
	}
}
