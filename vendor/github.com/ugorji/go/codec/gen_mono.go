// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build codec.build

package codec

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
)

// This tool will monomorphize types scoped to a specific format.
//
// This tool only monomorphized the type Name, and not a function Name.
// Explicitly, generic functions are not supported, as they cannot be monomorphized
// to a specific format without a corresponding name change.
//
// However, for types constrained to encWriter or decReader,
// which are shared across formats, there's no place to put them without duplication.

const genMonoParserMode = parser.AllErrors | parser.SkipObjectResolution

var genMonoSpecialFieldTypes = []string{"helperDecReader"}

// These functions should take the address of first param when monomorphized
var genMonoSpecialFunc4Addr = []string{} // {"decByteSlice"}

var genMonoImportsToSkip = []string{`"errors"`, `"fmt"`, `"net/rpc"`}

var genMonoRefImportsVia_ = [][2]string{
	// {"errors", "New"},
}

var genMonoCallsToSkip = []string{"callMake"}

type genMonoFieldState uint

const (
	genMonoFieldRecv genMonoFieldState = iota << 1
	genMonoFieldParamsResult
	genMonoFieldStruct
)

type genMonoImports struct {
	set   map[string]struct{}
	specs []*ast.ImportSpec
}

type genMono struct {
	files             map[string][]byte
	typParam          map[string]*ast.Field
	typParamTransient map[string]*ast.Field
}

func (x *genMono) init() {
	x.files = make(map[string][]byte)
	x.typParam = make(map[string]*ast.Field)
	x.typParamTransient = make(map[string]*ast.Field)
}

func (x *genMono) reset() {
	clear(x.typParam)
	clear(x.typParamTransient)
}

func (m *genMono) hdl(hname string) {
	m.reset()
	m.do(hname, []string{"encode.go", "decode.go", hname + ".go"}, []string{"base.notfastpath.go", "base.notfastpath.notmono.go"}, "", "")
	m.do(hname, []string{"base.notfastpath.notmono.go"}, nil, ".notfastpath", ` && (notfastpath || codec.notfastpath)`)
	m.do(hname, []string{"base.fastpath.notmono.generated.go"}, []string{"base.fastpath.generated.go"}, ".fastpath", ` && !notfastpath && !codec.notfastpath`)
}

func (m *genMono) do(hname string, fnames, tnames []string, fnameInfx string, buildTagsSfx string) {
	// keep m.typParams across whole call, as all others use it
	const fnameSfx = ".mono.generated.go"
	fname := hname + fnameInfx + fnameSfx

	var imports = genMonoImports{set: make(map[string]struct{})}

	r1, fset := m.merge(fnames, tnames, &imports)
	m.trFile(r1, hname, true)

	r2, fset := m.merge(fnames, tnames, &imports)
	m.trFile(r2, hname, false)

	r0 := genMonoOutInit(imports.specs, fname)
	r0.Decls = append(r0.Decls, r1.Decls...)
	r0.Decls = append(r0.Decls, r2.Decls...)

	// output r1 to a file
	f, err := os.Create(fname)
	halt.onerror(err)
	defer f.Close()

	var s genMonoStrBuilder
	s.s(`//go:build !notmono && !codec.notmono `).s(buildTagsSfx).s(`

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

`)
	_, err = f.Write(s.v)
	halt.onerror(err)
	err = format.Node(f, fset, r0)
	halt.onerror(err)
}

func (x *genMono) file(fname string) (b []byte) {
	b = x.files[fname]
	if b == nil {
		var err error
		b, err = os.ReadFile(fname)
		halt.onerror(err)
		x.files[fname] = b
	}
	return
}

func (x *genMono) merge(fNames, tNames []string, imports *genMonoImports) (dst *ast.File, fset *token.FileSet) {
	// typParams used in fnLoadTyps
	var typParams map[string]*ast.Field
	var loadTyps bool
	fnLoadTyps := func(node ast.Node) bool {
		var ok bool
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Tok == token.TYPE {
				for _, v := range n.Specs {
					nn := v.(*ast.TypeSpec)
					ok = genMonoTypeParamsOk(nn.TypeParams)
					if ok {
						// each decl will have only 1 var/type
						typParams[nn.Name.Name] = nn.TypeParams.List[0]
						if loadTyps {
							dst.Decls = append(dst.Decls, &ast.GenDecl{Tok: n.Tok, Specs: []ast.Spec{v}})
						}
					}
				}
			}
			return false
		}
		return true
	}

	// we only merge top-level methods and types
	fnIdX := func(n *ast.FuncDecl, n2 *ast.IndexExpr) (ok bool) {
		n9, ok9 := n2.Index.(*ast.Ident)
		n3, ok := n2.X.(*ast.Ident) // n3 = type name
		ok = ok && ok9 && n9.Name == "T"
		if ok {
			_, ok = x.typParam[n3.Name]
		}
		return
	}

	fnLoadMethodsAndImports := func(node ast.Node) bool {
		var ok bool
		switch n := node.(type) {
		case *ast.FuncDecl:
			// TypeParams is nil for methods, as it is defined at the type node
			// instead, look at the name, and
			// if IndexExpr.Index=T, and IndexExpr.X matches a type name seen already
			//     then ok = true
			if n.Recv == nil || len(n.Recv.List) != 1 {
				return false
			}
			ok = false
			switch nn := n.Recv.List[0].Type.(type) {
			case *ast.IndexExpr:
				ok = fnIdX(n, nn)
			case *ast.StarExpr:
				switch nn2 := nn.X.(type) {
				case *ast.IndexExpr:
					ok = fnIdX(n, nn2)
				}
			}
			if ok {
				dst.Decls = append(dst.Decls, n)
			}
			return false
		case *ast.GenDecl:
			if n.Tok == token.IMPORT {
				for _, v := range n.Specs {
					nn := v.(*ast.ImportSpec)
					if slices.Contains(genMonoImportsToSkip, nn.Path.Value) {
						continue
					}
					if _, ok = imports.set[nn.Path.Value]; !ok {
						imports.specs = append(imports.specs, nn)
						imports.set[nn.Path.Value] = struct{}{}
					}
				}
			}
			return false
		}
		return true
	}

	fset = token.NewFileSet()
	fnLoadAsts := func(names []string) (asts []*ast.File) {
		for _, fname := range names {
			fsrc := x.file(fname)
			f, err := parser.ParseFile(fset, fname, fsrc, genMonoParserMode)
			halt.onerror(err)
			asts = append(asts, f)
		}
		return
	}

	clear(x.typParamTransient)

	dst = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}

	fs := fnLoadAsts(fNames)
	ts := fnLoadAsts(tNames)

	loadTyps = true
	typParams = x.typParam
	for _, v := range fs {
		ast.Inspect(v, fnLoadTyps)
	}
	loadTyps = false
	typParams = x.typParamTransient
	for _, v := range ts {
		ast.Inspect(v, fnLoadTyps)
	}
	typParams = nil
	for _, v := range fs {
		ast.Inspect(v, fnLoadMethodsAndImports)
	}

	return
}

func (x *genMono) trFile(r *ast.File, hname string, isbytes bool) {
	fn := func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.TypeSpec:
			// type x[T encDriver] struct { ... }
			if !genMonoTypeParamsOk(n.TypeParams) {
				return false
			}
			x.trType(n, hname, isbytes)
			return false
		case *ast.FuncDecl:
			if n.Recv == nil || len(n.Recv.List) != 1 {
				return false
			}
			if _, ok := n.Recv.List[0].Type.(*ast.Ident); ok {
				return false
			}
			tp := x.trMethodSign(n, hname, isbytes) // receiver, params, results
			// handle the body
			x.trMethodBody(n.Body, tp, hname, isbytes)
			return false
		}
		return true
	}
	ast.Inspect(r, fn)

	// set type params to nil, and Pos to NoPos
	fn = func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncType:
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		case *ast.TypeSpec: // for type ...
			if genMonoTypeParamsOk(n.TypeParams) {
				n.TypeParams = nil
			}
		}
		return true
	}
	ast.Inspect(r, fn)
}

func (x *genMono) trType(n *ast.TypeSpec, hname string, isbytes bool) {
	sfx, _, _, hnameUp := genMonoIsBytesVals(hname, isbytes)
	tp := n.TypeParams.List[0]
	switch tp.Type.(*ast.Ident).Name {
	case "encDriver", "decDriver":
		n.Name.Name += hnameUp + sfx
	case "encWriter", "decReader":
		n.Name.Name += sfx
	}

	// handle the Struct and Array types
	switch nn := n.Type.(type) {
	case *ast.StructType:
		x.trStruct(nn, tp, hname, isbytes)
	case *ast.ArrayType:
		x.trArray(nn, tp, hname, isbytes)
	}
}

func (x *genMono) trMethodSign(n *ast.FuncDecl, hname string, isbytes bool) (tp *ast.Field) {
	// check if recv type is not parameterized
	tp = x.trField(n.Recv.List[0], nil, hname, isbytes, genMonoFieldRecv)
	// handle params and results
	x.trMethodSignNonRecv(n.Type.Params, tp, hname, isbytes)
	x.trMethodSignNonRecv(n.Type.Results, tp, hname, isbytes)
	return
}

func (x *genMono) trMethodSignNonRecv(r *ast.FieldList, tp *ast.Field, hname string, isbytes bool) {
	if r == nil || len(r.List) == 0 {
		return
	}
	for _, v := range r.List {
		x.trField(v, tp, hname, isbytes, genMonoFieldParamsResult)
	}
}

func (x *genMono) trStruct(r *ast.StructType, tp *ast.Field, hname string, isbytes bool) {
	// search for fields, and update accordingly
	//   type x[T encDriver] struct { w T }
	//   var x *A[T]
	//   A[T]
	if r == nil || r.Fields == nil || len(r.Fields.List) == 0 {
		return
	}
	for _, v := range r.Fields.List {
		x.trField(v, tp, hname, isbytes, genMonoFieldStruct)
	}
}

func (x *genMono) trArray(n *ast.ArrayType, tp *ast.Field, hname string, isbytes bool) {
	sfx, _, _, hnameUp := genMonoIsBytesVals(hname, isbytes)
	// type fastpathEs[T encDriver] [56]fastpathE[T]
	// p := tp.Names[0].Name
	switch elt := n.Elt.(type) {
	// case *ast.InterfaceType:
	case *ast.IndexExpr:
		if elt.Index.(*ast.Ident).Name == "T" { // generic
			n.Elt = ast.NewIdent(elt.X.(*ast.Ident).Name + hnameUp + sfx)
		}
	}
}

func (x *genMono) trMethodBody(r *ast.BlockStmt, tp *ast.Field, hname string, isbytes bool) {
	// find the parent node for an indexExpr, or a T/*T, and set the value back in there

	fn := func(pnode ast.Node) bool {
		var pn *ast.Ident
		fnUp := func() {
			x.updateIdentForT(pn, hname, tp, isbytes, false)
		}
		switch n := pnode.(type) {
		// case *ast.SelectorExpr:
		// case *ast.TypeAssertExpr:
		// case *ast.IndexExpr:
		case *ast.StarExpr:
			if genMonoUpdateIndexExprT(&pn, n.X) {
				n.X = pn
				fnUp()
			}
		case *ast.CallExpr:
			for i4, n4 := range n.Args {
				if genMonoUpdateIndexExprT(&pn, n4) {
					n.Args[i4] = pn
					fnUp()
				}
			}
			if n4, ok4 := n.Fun.(*ast.Ident); ok4 && slices.Contains(genMonoSpecialFunc4Addr, n4.Name) {
				n.Args[0] = &ast.UnaryExpr{Op: token.AND, X: n.Args[0].(*ast.SelectorExpr)}
			}
		case *ast.CompositeLit:
			if genMonoUpdateIndexExprT(&pn, n.Type) {
				n.Type = pn
				fnUp()
			}
		case *ast.ArrayType:
			if genMonoUpdateIndexExprT(&pn, n.Elt) {
				n.Elt = pn
				fnUp()
			}
		case *ast.ValueSpec:
			for i2, n2 := range n.Values {
				if genMonoUpdateIndexExprT(&pn, n2) {
					n.Values[i2] = pn
					fnUp()
				}
			}
			if genMonoUpdateIndexExprT(&pn, n.Type) {
				n.Type = pn
				fnUp()
			}
		case *ast.BinaryExpr:
			// early return here, since the 2 things can apply
			if genMonoUpdateIndexExprT(&pn, n.X) {
				n.X = pn
				fnUp()
			}
			if genMonoUpdateIndexExprT(&pn, n.Y) {
				n.Y = pn
				fnUp()
			}
			return true
		}
		return true
	}
	ast.Inspect(r, fn)
}

func (x *genMono) trField(f *ast.Field, tpt *ast.Field, hname string, isbytes bool, state genMonoFieldState) (tp *ast.Field) {
	var pn *ast.Ident
	switch nn := f.Type.(type) {
	case *ast.IndexExpr:
		if genMonoUpdateIndexExprT(&pn, nn) {
			f.Type = pn
		}
	case *ast.StarExpr:
		if genMonoUpdateIndexExprT(&pn, nn.X) {
			nn.X = pn
		}
	case *ast.FuncType:
		x.trMethodSignNonRecv(nn.Params, tpt, hname, isbytes)
		x.trMethodSignNonRecv(nn.Results, tpt, hname, isbytes)
		return
	case *ast.ArrayType:
		x.trArray(nn, tpt, hname, isbytes)
		return
	case *ast.Ident:
		if state == genMonoFieldRecv || nn.Name != "T" {
			return
		}
		pn = nn // "T"
		if state == genMonoFieldParamsResult {
			f.Type = &ast.StarExpr{X: pn}
		}
	}
	if pn == nil {
		return
	}

	tp = x.updateIdentForT(pn, hname, tpt, isbytes, true)
	return
}

func (x *genMono) updateIdentForT(pn *ast.Ident, hname string, tp *ast.Field,
	isbytes bool, lookupTP bool) (tp2 *ast.Field) {
	sfx, writer, reader, hnameUp := genMonoIsBytesVals(hname, isbytes)
	// handle special ones e.g. helperDecReader et al
	if slices.Contains(genMonoSpecialFieldTypes, pn.Name) {
		pn.Name += sfx
		return
	}

	if pn.Name != "T" && lookupTP {
		tp = x.typParam[pn.Name]
		if tp == nil {
			tp = x.typParamTransient[pn.Name]
		}
	}

	paramtyp := tp.Type.(*ast.Ident).Name
	if pn.Name == "T" {
		switch paramtyp {
		case "encDriver", "decDriver":
			pn.Name = hname + genMonoTitleCase(paramtyp) + sfx
		case "encWriter":
			pn.Name = writer
		case "decReader":
			pn.Name = reader
		}
	} else {
		switch paramtyp {
		case "encDriver", "decDriver":
			pn.Name += hnameUp + sfx
		case "encWriter", "decReader":
			pn.Name += sfx
		}
	}
	return tp
}

func genMonoUpdateIndexExprT(pn **ast.Ident, node ast.Node) (pnok bool) {
	*pn = nil
	if n2, ok := node.(*ast.IndexExpr); ok {
		n9, ok9 := n2.Index.(*ast.Ident)
		n3, ok := n2.X.(*ast.Ident)
		if ok && ok9 && n9.Name == "T" {
			*pn, pnok = ast.NewIdent(n3.Name), true
		}
	}
	return
}

func genMonoTitleCase(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func genMonoIsBytesVals(hName string, isbytes bool) (suffix, writer, reader, hNameUp string) {
	hNameUp = genMonoTitleCase(hName)
	if isbytes {
		return "Bytes", "bytesEncAppender", "bytesDecReader", hNameUp
	}
	return "IO", "bufioEncWriter", "ioDecReader", hNameUp
}

func genMonoTypeParamsOk(v *ast.FieldList) (ok bool) {
	if v == nil || v.List == nil || len(v.List) != 1 {
		return false
	}
	pn := v.List[0]
	if len(pn.Names) != 1 {
		return false
	}
	pnName := pn.Names[0].Name
	if pnName != "T" {
		return false
	}
	// ignore any nodes which are not idents e.g. cmp.orderedRv
	vv, ok := pn.Type.(*ast.Ident)
	if !ok {
		return false
	}
	switch vv.Name {
	case "encDriver", "decDriver", "encWriter", "decReader":
		return true
	}
	return false
}

func genMonoCopy(src *ast.File) (dst *ast.File) {
	dst = &ast.File{
		Name: &ast.Ident{Name: "codec"},
	}
	dst.Decls = append(dst.Decls, src.Decls...)
	return
}

type genMonoStrBuilder struct {
	v []byte
}

func (x *genMonoStrBuilder) s(v string) *genMonoStrBuilder {
	x.v = append(x.v, v...)
	return x
}

func genMonoOutInit(importSpecs []*ast.ImportSpec, fname string) (f *ast.File) {
	// ParseFile seems to skip the //go:build stanza
	// it should be written directly into the file
	var s genMonoStrBuilder
	s.s(`
package codec

import (
`)
	for _, v := range importSpecs {
		s.s("\t").s(v.Path.Value).s("\n")
	}
	s.s(")\n")
	for _, v := range genMonoRefImportsVia_ {
		s.s("var _ = ").s(v[0]).s(".").s(v[1]).s("\n")
	}
	f, err := parser.ParseFile(token.NewFileSet(), fname, s.v, genMonoParserMode)
	halt.onerror(err)
	return
}

func genMonoAll() {
	// hdls := []Handle{
	// 	(*SimpleHandle)(nil),
	// 	(*JsonHandle)(nil),
	// 	(*CborHandle)(nil),
	// 	(*BincHandle)(nil),
	// 	(*MsgpackHandle)(nil),
	// }
	hdls := []string{"simple", "json", "cbor", "binc", "msgpack"}
	var m genMono
	m.init()
	for _, v := range hdls {
		m.hdl(v)
	}
}
