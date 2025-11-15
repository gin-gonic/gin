package optdec

import (
	"reflect"
	"unsafe"

	"encoding/json"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/option"
	"github.com/bytedance/sonic/internal/decoder/errors"
	"github.com/bytedance/sonic/internal/decoder/consts"
)


type (
	MismatchTypeError = errors.MismatchTypeError
	SyntaxError = errors.SyntaxError
)

const (
	_F_allow_control = consts.F_allow_control
	_F_copy_string = consts.F_copy_string
	_F_disable_unknown = consts.F_disable_unknown
	_F_disable_urc = consts.F_disable_urc
	_F_use_int64 = consts.F_use_int64
	_F_use_number = consts.F_use_number
	_F_validate_string = consts.F_validate_string
)

type Options = consts.Options

const (
	OptionUseInt64     = consts.OptionUseInt64
	OptionUseNumber    = consts.OptionUseNumber
	OptionUseUnicodeErrors = consts.OptionUseUnicodeErrors
	OptionDisableUnknown = consts.OptionDisableUnknown
	OptionCopyString = consts.OptionCopyString
	OptionValidateString = consts.OptionValidateString
)


func Decode(s *string, i *int, f uint64, val interface{}) error {
	vv := rt.UnpackEface(val)
	vp := vv.Value

	/* check for nil type */
	if vv.Type == nil {
		return &json.InvalidUnmarshalError{}
	}

	/* must be a non-nil pointer */
	if vp == nil || vv.Type.Kind() != reflect.Ptr {
		return &json.InvalidUnmarshalError{Type: vv.Type.Pack()}
	}

	etp := rt.PtrElem(vv.Type)

	/* check the defined pointer type for issue 379 */
	if vv.Type.IsNamed() {
		newp := vp
		etp = vv.Type
		vp = unsafe.Pointer(&newp)
	}

	dec, err := findOrCompile(etp)
	if err != nil {
		return err
	}

	/* parse into document */
	ctx, err := NewContext(*s, *i, uint64(f), etp)
	defer ctx.Delete()
	if ctx.Parser.Utf8Inv {
		*s = ctx.Parser.Json
	}
	if err != nil {
		goto fix_error;
	}
	err = dec.FromDom(vp, ctx.Root(), &ctx)

fix_error:
	err = fix_error(*s, *i, err)

	// update position at last
	*i += ctx.Parser.Pos()
	return err
}

func fix_error(json string, pos int, err error) error {
	if e, ok := err.(SyntaxError); ok {
		return SyntaxError{
			Pos: int(e.Pos) + pos,
			Src: json,
			Msg: e.Msg,
		}
	}

	if e, ok := err.(MismatchTypeError); ok {
		return &MismatchTypeError {
			Pos: int(e.Pos) + pos,
			Src: json,
			Type: e.Type,
		}
	}

	return err
}

// Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
// order to reduce the first-hit latency.
//
// Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
// a compile option to set the depth of recursive compile for the nested struct type.
func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    cfg := option.DefaultCompileOptions()
    for _, opt := range opts {
        opt(&cfg)
    }
    return pretouchRec(map[reflect.Type]bool{vt:true}, cfg)
}

func pretouchType(_vt reflect.Type, opts option.CompileOptions) (map[reflect.Type]bool, error) {
    /* compile function */
    compiler := newCompiler().apply(opts)
    decoder := func(vt *rt.GoType, _ ...interface{}) (interface{}, error) {
        if f, err := compiler.compileType(_vt); err != nil {
            return nil, err
        } else {
            return f, nil
        }
    }

    /* find or compile */
    vt := rt.UnpackType(_vt)
    if val := programCache.Get(vt); val != nil {
        return nil, nil
    } else if _, err := programCache.Compute(vt, decoder); err == nil {
        return compiler.visited, nil
    } else {
        return nil, err
    }
}

func pretouchRec(vtm map[reflect.Type]bool, opts option.CompileOptions) error {
    if opts.RecursiveDepth < 0 || len(vtm) == 0 {
        return nil
    }
    next := make(map[reflect.Type]bool)
    for vt := range(vtm) {
        sub, err := pretouchType(vt, opts)
        if err != nil {
            return err
        }
        for svt := range(sub) {
            next[svt] = true
        }
    }
    opts.RecursiveDepth -= 1
    return pretouchRec(next, opts)
}
