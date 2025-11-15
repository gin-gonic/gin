// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build codec.build

package codec

import (
	"bytes"
	"encoding/base32"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"text/template"
	// "ugorji.net/zz"
)

// ---------------------------------------------------

const (
	genTopLevelVarName = "x"

	// genFastpathCanonical configures whether we support Canonical in fast path. Low savings.
	//
	// MARKER: This MUST ALWAYS BE TRUE. fastpath.go.tmpl doesn't handle it being false.
	genFastpathCanonical = true

	// genFastpathTrimTypes configures whether we trim uncommon fastpath types.
	genFastpathTrimTypes = true
)

var genFormats = []string{"Json", "Cbor", "Msgpack", "Binc", "Simple"}

var (
	errGenAllTypesSamePkg        = errors.New("All types must be in the same package")
	errGenExpectArrayOrMap       = errors.New("unexpected type - expecting array/map/slice")
	errGenUnexpectedTypeFastpath = errors.New("fastpath: unexpected type - requires map or slice")

	// don't use base64, only 63 characters allowed in valid go identifiers
	// ie ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_
	//
	// don't use numbers, as a valid go identifer must start with a letter.
	genTypenameEnc = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef")
	genQNameRegex  = regexp.MustCompile(`[A-Za-z_.]+`)
)

// --------

func genCheckErr(err error) {
	halt.onerror(err)
}

func genTitleCaseName(s string) string {
	switch s {
	case "interface{}", "interface {}":
		return "Intf"
	case "[]byte", "[]uint8", "bytes":
		return "Bytes"
	default:
		return strings.ToUpper(s[0:1]) + s[1:]
	}
}

// --------

type genFastpathV struct {
	// genFastpathV is either a primitive (Primitive != "") or a map (MapKey != "") or a slice
	MapKey      string
	Elem        string
	Primitive   string
	Size        int
	NoCanonical bool
}

func (x *genFastpathV) MethodNamePfx(prefix string, prim bool) string {
	var name []byte
	if prefix != "" {
		name = append(name, prefix...)
	}
	if prim {
		name = append(name, genTitleCaseName(x.Primitive)...)
	} else {
		if x.MapKey == "" {
			name = append(name, "Slice"...)
		} else {
			name = append(name, "Map"...)
			name = append(name, genTitleCaseName(x.MapKey)...)
		}
		name = append(name, genTitleCaseName(x.Elem)...)
	}
	return string(name)
}

// --------

type genTmpl struct {
	Values  []genFastpathV
	Formats []string
}

func (x genTmpl) FastpathLen() (l int) {
	for _, v := range x.Values {
		// if v.Primitive == "" && !(v.MapKey == "" && v.Elem == "uint8") {
		if v.Primitive == "" {
			l++
		}
	}
	return
}

func genTmplZeroValue(s string) string {
	switch s {
	case "interface{}", "interface {}":
		return "nil"
	case "[]byte", "[]uint8", "bytes":
		return "nil"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return "0"
	}
}

var genTmplNonZeroValueIdx [6]uint64
var genTmplNonZeroValueStrs = [...][6]string{
	{`"string-is-an-interface-1"`, "true", `"some-string-1"`, `[]byte("some-string-1")`, "11.1", "111"},
	{`"string-is-an-interface-2"`, "false", `"some-string-2"`, `[]byte("some-string-2")`, "22.2", "77"},
	{`"string-is-an-interface-3"`, "true", `"some-string-3"`, `[]byte("some-string-3")`, "33.3e3", "127"},
}

// Note: last numbers must be in range: 0-127 (as they may be put into a int8, uint8, etc)

func genTmplNonZeroValue(s string) string {
	var i int
	switch s {
	case "interface{}", "interface {}":
		i = 0
	case "bool":
		i = 1
	case "string":
		i = 2
	case "bytes", "[]byte", "[]uint8":
		i = 3
	case "float32", "float64", "float", "double", "complex", "complex64", "complex128":
		i = 4
	default:
		i = 5
	}
	genTmplNonZeroValueIdx[i]++
	idx := genTmplNonZeroValueIdx[i]
	slen := uint64(len(genTmplNonZeroValueStrs))
	return genTmplNonZeroValueStrs[idx%slen][i] // return string, to remove ambiguity
}

// Note: used for fastpath only
func genTmplEncCommandAsString(s string, vname string) string {
	switch s {
	case "uint64":
		return "e.e.EncodeUint(" + vname + ")"
	case "uint", "uint8", "uint16", "uint32":
		return "e.e.EncodeUint(uint64(" + vname + "))"
	case "int64":
		return "e.e.EncodeInt(" + vname + ")"
	case "int", "int8", "int16", "int32":
		return "e.e.EncodeInt(int64(" + vname + "))"
	case "[]byte", "[]uint8", "bytes":
		// return fmt.Sprintf(
		// 	"if %s != nil { e.e.EncodeStringBytesRaw(%s) } "+
		// 		"else if e.h.NilCollectionToZeroLength { e.e.WriteArrayEmpty() } "+
		// 		"else { e.e.EncodeNil() }", vname, vname)
		// return "e.e.EncodeStringBytesRaw(" + vname + ")"
		return "e.e.EncodeBytes(" + vname + ")"
	case "string":
		return "e.e.EncodeString(" + vname + ")"
	case "float32":
		return "e.e.EncodeFloat32(" + vname + ")"
	case "float64":
		return "e.e.EncodeFloat64(" + vname + ")"
	case "bool":
		return "e.e.EncodeBool(" + vname + ")"
	// case "symbol":
	// 	return "e.e.EncodeSymbol(" + vname + ")"
	default:
		return fmt.Sprintf("if !e.encodeBuiltin(%s) { e.encodeR(reflect.ValueOf(%s)) }", vname, vname)
		// return "e.encodeI(" + vname + ")"
	}
}

// Note: used for fastpath only
func genTmplDecCommandAsString(s string, mapkey bool) string {
	switch s {
	case "uint":
		return "uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "uint8":
		return "uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))"
	case "uint16":
		return "uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))"
	case "uint32":
		return "uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))"
	case "uint64":
		return "d.d.DecodeUint64()"
	case "uintptr":
		return "uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "int":
		return "int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))"
	case "int8":
		return "int8(chkOvf.IntV(d.d.DecodeInt64(), 8))"
	case "int16":
		return "int16(chkOvf.IntV(d.d.DecodeInt64(), 16))"
	case "int32":
		return "int32(chkOvf.IntV(d.d.DecodeInt64(), 32))"
	case "int64":
		return "d.d.DecodeInt64()"

	case "string":
		// if mapkey {
		// 	return "d.stringZC(d.d.DecodeStringAsBytes())"
		// }
		// return "string(d.d.DecodeStringAsBytes())"
		return "d.detach2Str(d.d.DecodeStringAsBytes())"
	case "[]byte", "[]uint8", "bytes":
		// return "bytesOk(d.d.DecodeBytes())"
		return "bytesOKdbi(d.decodeBytesInto(v[uint(j)], false))"
	case "float32":
		return "float32(d.d.DecodeFloat32())"
	case "float64":
		return "d.d.DecodeFloat64()"
	case "complex64":
		return "complex(d.d.DecodeFloat32(), 0)"
	case "complex128":
		return "complex(d.d.DecodeFloat64(), 0)"
	case "bool":
		return "d.d.DecodeBool()"
	default:
		halt.error(errors.New("gen internal: unknown type for decode: " + s))
	}
	return ""
}

func genTmplSortType(s string, elem bool) string {
	if elem {
		return s
	}
	return s + "Slice"
}

// var genTmplMu sync.Mutex
var genTmplV = genTmpl{}
var genTmplFuncs template.FuncMap
var genTmplOnce sync.Once

func genTmplInit() {
	wordSizeBytes := int(intBitsize) / 8

	typesizes := map[string]int{
		"interface{}": 2 * wordSizeBytes,
		"string":      2 * wordSizeBytes,
		"[]byte":      3 * wordSizeBytes,
		"uint":        1 * wordSizeBytes,
		"uint8":       1,
		"uint16":      2,
		"uint32":      4,
		"uint64":      8,
		"uintptr":     1 * wordSizeBytes,
		"int":         1 * wordSizeBytes,
		"int8":        1,
		"int16":       2,
		"int32":       4,
		"int64":       8,
		"float32":     4,
		"float64":     8,
		"complex64":   8,
		"complex128":  16,
		"bool":        1,
	}

	// keep as slice, so it is in specific iteration order.
	// Initial order was uint64, string, interface{}, int, int64, ...

	var types = [...]string{
		"interface{}",
		"string",
		"[]byte",
		"float32",
		"float64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"bool",
	}

	var primitivetypes, slicetypes, mapkeytypes, mapvaltypes []string

	primitivetypes = types[:]

	slicetypes = types[:]
	mapkeytypes = types[:]
	mapvaltypes = types[:]

	if genFastpathTrimTypes {
		// Note: we only create fastpaths for commonly used types.
		// Consequently, things like int8, uint16, uint, etc are commented out.
		slicetypes = []string{
			"interface{}",
			"string",
			"[]byte",
			"float32",
			"float64",
			"uint8", // keep fastpath, so it doesn't have to go through reflection
			"uint64",
			"int",
			"int32", // rune
			"int64",
			"bool",
		}
		mapkeytypes = []string{
			"string",
			"uint8",  // byte
			"uint64", // used for keys
			"int",    // default number key
			"int32",  // rune
		}
		mapvaltypes = []string{
			"interface{}",
			"string",
			"[]byte",
			"uint8",  // byte
			"uint64", // used for keys, etc
			"int",    // default number
			"int32",  // rune (mostly used for unicode)
			"float64",
			"bool",
		}
	}

	var gt = genTmpl{Formats: genFormats}

	// For each slice or map type, there must be a (symmetrical) Encode and Decode fastpath function

	for _, s := range primitivetypes {
		gt.Values = append(gt.Values,
			genFastpathV{Primitive: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical})
	}
	for _, s := range slicetypes {
		gt.Values = append(gt.Values,
			genFastpathV{Elem: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical})
	}
	for _, s := range mapkeytypes {
		for _, ms := range mapvaltypes {
			gt.Values = append(gt.Values,
				genFastpathV{MapKey: s, Elem: ms, Size: typesizes[s] + typesizes[ms], NoCanonical: !genFastpathCanonical})
		}
	}

	funcs := make(template.FuncMap)
	// funcs["haspfx"] = strings.HasPrefix
	funcs["encmd"] = genTmplEncCommandAsString
	funcs["decmd"] = genTmplDecCommandAsString
	funcs["zerocmd"] = genTmplZeroValue
	funcs["nonzerocmd"] = genTmplNonZeroValue
	funcs["hasprefix"] = strings.HasPrefix
	funcs["sorttype"] = genTmplSortType

	genTmplV = gt
	genTmplFuncs = funcs
}

// genTmplGoFile is used to generate source files from templates.
func genTmplGoFile(r io.Reader, w io.Writer) (err error) {
	genTmplOnce.Do(genTmplInit)

	gt := genTmplV

	t := template.New("").Funcs(genTmplFuncs)

	tmplstr, err := io.ReadAll(r)
	if err != nil {
		return
	}

	if t, err = t.Parse(string(tmplstr)); err != nil {
		return
	}

	var out bytes.Buffer
	err = t.Execute(&out, gt)
	if err != nil {
		return
	}

	bout, err := format.Source(out.Bytes())
	if err != nil {
		w.Write(out.Bytes()) // write out if error, so we can still see.
		// w.Write(bout) // write out if error, as much as possible, so we can still see.
		return
	}
	w.Write(bout)
	return
}

func genTmplRun2Go(fnameIn, fnameOut string) {
	// println("____ " + fnameIn + " --> " + fnameOut + " ______")
	fin, err := os.Open(fnameIn)
	genCheckErr(err)
	defer fin.Close()
	fout, err := os.Create(fnameOut)
	genCheckErr(err)
	defer fout.Close()
	err = genTmplGoFile(fin, fout)
	genCheckErr(err)
}
