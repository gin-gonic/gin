package optdec

import (
	"fmt"
	"reflect"
	"unsafe"

	"sync"

	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/utf8"
)


type ErrorCode int

const (
	SONIC_OK                = 0;
	SONIC_CONTROL_CHAR      = 1;
	SONIC_INVALID_ESCAPED   = 2;
	SONIC_INVALID_NUM       = 3;
	SONIC_FLOAT_INF         = 4;
	SONIC_EOF               = 5;
	SONIC_INVALID_CHAR   	= 6;
	SONIC_EXPECT_KEY        = 7;
	SONIC_EXPECT_COLON      = 8;
	SONIC_EXPECT_OBJ_COMMA_OR_END  = 9;
	SONIC_EXPECT_ARR_COMMA_OR_END  = 10;
	SONIC_VISIT_FAILED        = 11;
	SONIC_INVALID_ESCAPED_UTF = 12;
	SONIC_INVALID_LITERAL 	  = 13;
	SONIC_STACK_OVERFLOW	  = 14;
)

var ParsingErrors = []string{
    SONIC_OK                      	: "ok",
	SONIC_CONTROL_CHAR      	  	: "control chars in string",
	SONIC_INVALID_ESCAPED 		  	: "invalid escaped chars in string",
	SONIC_INVALID_NUM       		: "invalid number",
	SONIC_FLOAT_INF         		: "float infinity",
	SONIC_EOF               		: "eof",
	SONIC_INVALID_CHAR		  		: "invalid chars",
	SONIC_EXPECT_KEY        		: "expect a json key",
	SONIC_EXPECT_COLON      		: "expect a `:`",
	SONIC_EXPECT_OBJ_COMMA_OR_END	: "expect a `,` or `}`",
	SONIC_EXPECT_ARR_COMMA_OR_END	: "expect a `,` or `]`",
	SONIC_VISIT_FAILED     			: "failed in json visitor",
	SONIC_INVALID_ESCAPED_UTF		: "invalid escaped unicodes",
	SONIC_INVALID_LITERAL			: "invalid literal(true/false/null)",
	SONIC_STACK_OVERFLOW			: "json is exceeded max depth 4096, cause stack overflow",
}

func (code ErrorCode) Error() string {
	return ParsingErrors[code]
}

type node struct {
	typ uint64
	val uint64
}

// should consistent with native/parser.c
type _nospaceBlock struct {
	_ [8]byte
	_ [8]byte
}

// should consistent with native/parser.c
type nodeBuf struct {
	ncur    uintptr
	parent  int64
	depth   uint64
	nstart  uintptr
	nend    uintptr
	iskey   bool
	stat    jsonStat
}

func (self *nodeBuf) init(nodes []node) {
	self.ncur = uintptr(unsafe.Pointer(&nodes[0]))
	self.nstart = self.ncur
	self.nend = self.ncur + uintptr(cap(nodes)) * unsafe.Sizeof(node{})
	self.parent = -1
}

// should consistent with native/parser.c
type Parser struct {
	Json    string
	padded	[]byte
	nodes 	[]node
	dbuf 	[]byte
	backup  []node

	options uint64
	// JSON cursor
	start   uintptr
	cur     uintptr
	end     uintptr
	_nbk    _nospaceBlock

	// node buffer cursor
	nbuf   	nodeBuf
	Utf8Inv  	bool
	isEface    bool
}

// only when parse non-empty object/array are needed.
type jsonStat struct {
    object 		uint32
    array 		uint32
    str 		uint32
    number 		uint32
    array_elems uint32
    object_keys uint32
    max_depth	uint32
}


var (
	defaultJsonPaddedCap uintptr =  1 << 20  // 1 Mb
	defaultNodesCap      uintptr =  (1 << 20) / unsafe.Sizeof(node{})  // 1 Mb
)

var parsePool sync.Pool = sync.Pool {
	New: func () interface{} {
		return &Parser{
			options: 0,
			padded: make([]byte, 0, defaultJsonPaddedCap),
			nodes: make([]node, defaultNodesCap, defaultNodesCap),
			dbuf: make([]byte, types.MaxDigitNums, types.MaxDigitNums),
		}
	},
}

var padding string = "x\"x\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

func newParser(data string, pos int, opt uint64) *Parser {
	p := parsePool.Get().(*Parser)

	/* validate json if needed */
	if (opt & (1 << _F_validate_string)) != 0  && !utf8.ValidateString(data){
		dbuf := utf8.CorrectWith(nil, rt.Str2Mem(data[pos:]), "\ufffd")
		dbuf = append(dbuf, padding...)
		p.Json = rt.Mem2Str(dbuf[:len(dbuf) - len(padding)])
		p.Utf8Inv = true
		p.start = uintptr((*rt.GoString)(unsafe.Pointer(&p.Json)).Ptr)
	} else {
		p.Json = data
		// TODO: prevent too large JSON
		p.padded = append(p.padded, data[pos:]...)
		p.padded = append(p.padded, padding...)
		p.start = uintptr((*rt.GoSlice)(unsafe.Pointer(&p.padded)).Ptr)
	}

	p.cur 	= p.start
	p.end   = p.cur + uintptr(len(p.Json))
	p.options = opt
	p.nbuf.init(p.nodes)
	return p
}


func (p *Parser) Pos() int {
	return int(p.cur - p.start)
}

func (p *Parser) JsonBytes() []byte {
	if p.Utf8Inv {
		return (rt.Str2Mem(p.Json))
	} else {
		return p.padded
	}
}

var nodeType = rt.UnpackType(reflect.TypeOf(node{}))

//go:inline
func calMaxNodeCap(jsonSize int) int {
	return jsonSize / 2 + 2
}

func (p *Parser) parse() ErrorCode {
	// when decode into struct, we should decode number as possible
	old := p.options
	if !p.isEface {
		p.options &^= 1 << _F_use_number
	}

	// fast path with limited node buffer
	err := ErrorCode(native.ParseWithPadding(unsafe.Pointer(p)))
	if err != SONIC_VISIT_FAILED {
		p.options = old
		return err
	}

	// check OoB here
	offset := p.nbuf.ncur - p.nbuf.nstart
	curLen :=  int(offset / unsafe.Sizeof(node{}))
	if curLen != len(p.nodes) {
		panic(fmt.Sprintf("current len: %d, real len: %d cap: %d", curLen, len(p.nodes), cap(p.nodes)))
	}

	// node buf is not enough, continue parse
	// the maxCap is always meet all valid JSON
	maxCap := curLen + calMaxNodeCap(len(p.Json) - int(p.cur - p.start))
	slice := rt.GoSlice{
		Ptr: rt.Mallocgc(uintptr(maxCap) * nodeType.Size, nodeType, false),
		Len: maxCap,
		Cap: maxCap,
	}
	rt.Memmove(unsafe.Pointer(slice.Ptr), unsafe.Pointer(&p.nodes[0]), offset)
	p.backup = p.nodes
	p.nodes = *(*[]node)(unsafe.Pointer(&slice))

	// update node cursor
	p.nbuf.nstart = uintptr(unsafe.Pointer(&p.nodes[0]))
	p.nbuf.nend = p.nbuf.nstart + uintptr(cap(p.nodes)) * unsafe.Sizeof(node{})
	p.nbuf.ncur = p.nbuf.nstart + offset

	// continue parse json
	err = ErrorCode(native.ParseWithPadding(unsafe.Pointer(p)))
	p.options = old
	return err
}

func (p *Parser) reset() {
	p.options = 0
	p.padded = p.padded[:0]
	// nodes is too large here, we will not reset it and use small backup nodes buffer
	if p.backup != nil {
		p.nodes = p.backup
		p.backup = nil
	}
	p.start = 0
	p.cur = 0
	p.end = 0
	p.Json = ""
	p.nbuf = nodeBuf{}
	p._nbk = _nospaceBlock{}
	p.Utf8Inv = false
	p.isEface = false
}

func (p *Parser) free() {
	p.reset()
	parsePool.Put(p)
}

//go:noinline
func (p *Parser) fixError(code ErrorCode) error {
	if code == SONIC_OK {
		return nil
	}

	if p.Pos() == 0 {
		code = SONIC_EOF;
	}

	pos := p.Pos() - 1
	return error_syntax(pos, p.Json, ParsingErrors[code])
}

func Parse(data string, opt uint64) error {
	p := newParser(data, 0, opt)
	err := p.parse()
	p.free()
	return err
}
