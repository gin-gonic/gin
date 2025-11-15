package jitdec

import (
    `unsafe`
    `encoding/json`
    `reflect`
    `runtime`

	`github.com/bytedance/sonic/internal/decoder/consts`
	`github.com/bytedance/sonic/internal/decoder/errors`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/utf8`
	`github.com/bytedance/sonic/option`
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
	_F_no_validate_json = consts.F_no_validate_json
	_F_validate_string = consts.F_validate_string
    _F_case_sensitive = consts.F_case_sensitive
)

var (
	error_wrap = errors.ErrorWrap
	error_type = errors.ErrorType
	error_field = errors.ErrorField
	error_value = errors.ErrorValue
	error_mismatch = errors.ErrorMismatch
	stackOverflow = errors.StackOverflow
)


// Decode parses the JSON-encoded data from current position and stores the result
// in the value pointed to by val.
func Decode(s *string, i *int, f uint64, val interface{}) error {
    /* validate json if needed */
    if (f & (1 << _F_validate_string)) != 0  && !utf8.ValidateString(*s){
        dbuf := utf8.CorrectWith(nil, rt.Str2Mem(*s), "\ufffd")
        *s = rt.Mem2Str(dbuf)
    }

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
        etp  = vv.Type
        vp   = unsafe.Pointer(&newp)
    }

    /* create a new stack, and call the decoder */
    sb := newStack()
    nb, err := decodeTypedPointer(*s, *i, etp, vp, sb, f)
    /* return the stack back */
    *i = nb
    freeStack(sb)

    /* avoid GC ahead */
    runtime.KeepAlive(vv)
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
        if pp, err := compiler.compile(_vt); err != nil {
            return nil, err
        } else {
            as := newAssembler(pp)
            as.name = _vt.String()
            return as.Load(), nil
        }
    }

    /* find or compile */
    vt := rt.UnpackType(_vt)
    if val := programCache.Get(vt); val != nil {
        return nil, nil
    } else if _, err := programCache.Compute(vt, decoder); err == nil {
        return compiler.rec, nil
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

