package yaml

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
)

const (
	// DefaultIndentSpaces default number of space for indent
	DefaultIndentSpaces = 2
)

// Encoder writes YAML values to an output stream.
type Encoder struct {
	writer                     io.Writer
	opts                       []EncodeOption
	singleQuote                bool
	isFlowStyle                bool
	isJSONStyle                bool
	useJSONMarshaler           bool
	enableSmartAnchor          bool
	aliasRefToName             map[uintptr]string
	anchorRefToName            map[uintptr]string
	anchorNameMap              map[string]struct{}
	anchorCallback             func(*ast.AnchorNode, interface{}) error
	customMarshalerMap         map[reflect.Type]func(context.Context, interface{}) ([]byte, error)
	omitZero                   bool
	omitEmpty                  bool
	autoInt                    bool
	useLiteralStyleIfMultiline bool
	commentMap                 map[*Path][]*Comment
	written                    bool

	line           int
	column         int
	offset         int
	indentNum      int
	indentLevel    int
	indentSequence bool
}

// NewEncoder returns a new encoder that writes to w.
// The Encoder should be closed after use to flush all data to w.
func NewEncoder(w io.Writer, opts ...EncodeOption) *Encoder {
	return &Encoder{
		writer:             w,
		opts:               opts,
		customMarshalerMap: map[reflect.Type]func(context.Context, interface{}) ([]byte, error){},
		line:               1,
		column:             1,
		offset:             0,
		indentNum:          DefaultIndentSpaces,
		anchorRefToName:    make(map[uintptr]string),
		anchorNameMap:      make(map[string]struct{}),
		aliasRefToName:     make(map[uintptr]string),
	}
}

// Close closes the encoder by writing any remaining data.
// It does not write a stream terminating string "...".
func (e *Encoder) Close() error {
	return nil
}

// Encode writes the YAML encoding of v to the stream.
// If multiple items are encoded to the stream,
// the second and subsequent document will be preceded with a "---" document separator,
// but the first will not.
//
// See the documentation for Marshal for details about the conversion of Go values to YAML.
func (e *Encoder) Encode(v interface{}) error {
	return e.EncodeContext(context.Background(), v)
}

// EncodeContext writes the YAML encoding of v to the stream with context.Context.
func (e *Encoder) EncodeContext(ctx context.Context, v interface{}) error {
	node, err := e.EncodeToNodeContext(ctx, v)
	if err != nil {
		return err
	}
	if err := e.setCommentByCommentMap(node); err != nil {
		return err
	}
	if !e.written {
		e.written = true
	} else {
		// write document separator
		_, _ = e.writer.Write([]byte("---\n"))
	}
	var p printer.Printer
	_, _ = e.writer.Write(p.PrintNode(node))
	return nil
}

// EncodeToNode convert v to ast.Node.
func (e *Encoder) EncodeToNode(v interface{}) (ast.Node, error) {
	return e.EncodeToNodeContext(context.Background(), v)
}

// EncodeToNodeContext convert v to ast.Node with context.Context.
func (e *Encoder) EncodeToNodeContext(ctx context.Context, v interface{}) (ast.Node, error) {
	for _, opt := range e.opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	if e.enableSmartAnchor {
		// during the first encoding, store all mappings between alias addresses and their names.
		if _, err := e.encodeValue(ctx, reflect.ValueOf(v), 1); err != nil {
			return nil, err
		}
		e.clearSmartAnchorRef()
	}
	node, err := e.encodeValue(ctx, reflect.ValueOf(v), 1)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (e *Encoder) setCommentByCommentMap(node ast.Node) error {
	if e.commentMap == nil {
		return nil
	}
	for path, comments := range e.commentMap {
		n, err := path.FilterNode(node)
		if err != nil {
			return err
		}
		if n == nil {
			continue
		}
		for _, comment := range comments {
			commentTokens := []*token.Token{}
			for _, text := range comment.Texts {
				commentTokens = append(commentTokens, token.New(text, text, nil))
			}
			commentGroup := ast.CommentGroup(commentTokens)
			switch comment.Position {
			case CommentHeadPosition:
				if err := e.setHeadComment(node, n, commentGroup); err != nil {
					return err
				}
			case CommentLinePosition:
				if err := e.setLineComment(node, n, commentGroup); err != nil {
					return err
				}
			case CommentFootPosition:
				if err := e.setFootComment(node, n, commentGroup); err != nil {
					return err
				}
			default:
				return ErrUnknownCommentPositionType
			}
		}
	}
	return nil
}

func (e *Encoder) setHeadComment(node ast.Node, filtered ast.Node, comment *ast.CommentGroupNode) error {
	parent := ast.Parent(node, filtered)
	if parent == nil {
		return ErrUnsupportedHeadPositionType(node)
	}
	switch p := parent.(type) {
	case *ast.MappingValueNode:
		if err := p.SetComment(comment); err != nil {
			return err
		}
	case *ast.MappingNode:
		if err := p.SetComment(comment); err != nil {
			return err
		}
	case *ast.SequenceNode:
		if len(p.ValueHeadComments) == 0 {
			p.ValueHeadComments = make([]*ast.CommentGroupNode, len(p.Values))
		}
		var foundIdx int
		for idx, v := range p.Values {
			if v == filtered {
				foundIdx = idx
				break
			}
		}
		p.ValueHeadComments[foundIdx] = comment
	default:
		return ErrUnsupportedHeadPositionType(node)
	}
	return nil
}

func (e *Encoder) setLineComment(node ast.Node, filtered ast.Node, comment *ast.CommentGroupNode) error {
	switch filtered.(type) {
	case *ast.MappingValueNode, *ast.SequenceNode:
		// Line comment cannot be set for mapping value node.
		// It should probably be set for the parent map node
		if err := e.setLineCommentToParentMapNode(node, filtered, comment); err != nil {
			return err
		}
	default:
		if err := filtered.SetComment(comment); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) setLineCommentToParentMapNode(node ast.Node, filtered ast.Node, comment *ast.CommentGroupNode) error {
	parent := ast.Parent(node, filtered)
	if parent == nil {
		return ErrUnsupportedLinePositionType(node)
	}
	switch p := parent.(type) {
	case *ast.MappingValueNode:
		if err := p.Key.SetComment(comment); err != nil {
			return err
		}
	case *ast.MappingNode:
		if err := p.SetComment(comment); err != nil {
			return err
		}
	default:
		return ErrUnsupportedLinePositionType(parent)
	}
	return nil
}

func (e *Encoder) setFootComment(node ast.Node, filtered ast.Node, comment *ast.CommentGroupNode) error {
	parent := ast.Parent(node, filtered)
	if parent == nil {
		return ErrUnsupportedFootPositionType(node)
	}
	switch n := parent.(type) {
	case *ast.MappingValueNode:
		n.FootComment = comment
	case *ast.MappingNode:
		n.FootComment = comment
	case *ast.SequenceNode:
		n.FootComment = comment
	default:
		return ErrUnsupportedFootPositionType(n)
	}
	return nil
}

func (e *Encoder) encodeDocument(doc []byte) (ast.Node, error) {
	f, err := parser.ParseBytes(doc, 0)
	if err != nil {
		return nil, err
	}
	for _, docNode := range f.Docs {
		if docNode.Body != nil {
			return docNode.Body, nil
		}
	}
	return nil, nil
}

func (e *Encoder) isInvalidValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	kind := v.Type().Kind()
	if kind == reflect.Ptr && v.IsNil() {
		return true
	}
	if kind == reflect.Interface && v.IsNil() {
		return true
	}
	return false
}

type jsonMarshaler interface {
	MarshalJSON() ([]byte, error)
}

func (e *Encoder) existsTypeInCustomMarshalerMap(t reflect.Type) bool {
	if _, exists := e.customMarshalerMap[t]; exists {
		return true
	}

	globalCustomMarshalerMu.Lock()
	defer globalCustomMarshalerMu.Unlock()
	if _, exists := globalCustomMarshalerMap[t]; exists {
		return true
	}
	return false
}

func (e *Encoder) marshalerFromCustomMarshalerMap(t reflect.Type) (func(context.Context, interface{}) ([]byte, error), bool) {
	if marshaler, exists := e.customMarshalerMap[t]; exists {
		return marshaler, exists
	}

	globalCustomMarshalerMu.Lock()
	defer globalCustomMarshalerMu.Unlock()
	if marshaler, exists := globalCustomMarshalerMap[t]; exists {
		return marshaler, exists
	}
	return nil, false
}

func (e *Encoder) canEncodeByMarshaler(v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	if e.existsTypeInCustomMarshalerMap(v.Type()) {
		return true
	}
	iface := v.Interface()
	switch iface.(type) {
	case BytesMarshalerContext:
		return true
	case BytesMarshaler:
		return true
	case InterfaceMarshalerContext:
		return true
	case InterfaceMarshaler:
		return true
	case time.Time, *time.Time:
		return true
	case time.Duration:
		return true
	case encoding.TextMarshaler:
		return true
	case jsonMarshaler:
		return e.useJSONMarshaler
	}
	return false
}

func (e *Encoder) encodeByMarshaler(ctx context.Context, v reflect.Value, column int) (ast.Node, error) {
	iface := v.Interface()

	if marshaler, exists := e.marshalerFromCustomMarshalerMap(v.Type()); exists {
		doc, err := marshaler(ctx, iface)
		if err != nil {
			return nil, err
		}
		node, err := e.encodeDocument(doc)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	if marshaler, ok := iface.(BytesMarshalerContext); ok {
		doc, err := marshaler.MarshalYAML(ctx)
		if err != nil {
			return nil, err
		}
		node, err := e.encodeDocument(doc)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	if marshaler, ok := iface.(BytesMarshaler); ok {
		doc, err := marshaler.MarshalYAML()
		if err != nil {
			return nil, err
		}
		node, err := e.encodeDocument(doc)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	if marshaler, ok := iface.(InterfaceMarshalerContext); ok {
		marshalV, err := marshaler.MarshalYAML(ctx)
		if err != nil {
			return nil, err
		}
		return e.encodeValue(ctx, reflect.ValueOf(marshalV), column)
	}

	if marshaler, ok := iface.(InterfaceMarshaler); ok {
		marshalV, err := marshaler.MarshalYAML()
		if err != nil {
			return nil, err
		}
		return e.encodeValue(ctx, reflect.ValueOf(marshalV), column)
	}

	if t, ok := iface.(time.Time); ok {
		return e.encodeTime(t, column), nil
	}
	// Handle *time.Time explicitly since it implements TextMarshaler and shouldn't be treated as plain text
	if t, ok := iface.(*time.Time); ok && t != nil {
		return e.encodeTime(*t, column), nil
	}

	if t, ok := iface.(time.Duration); ok {
		return e.encodeDuration(t, column), nil
	}

	if marshaler, ok := iface.(encoding.TextMarshaler); ok {
		text, err := marshaler.MarshalText()
		if err != nil {
			return nil, err
		}
		node := e.encodeString(string(text), column)
		return node, nil
	}

	if e.useJSONMarshaler {
		if marshaler, ok := iface.(jsonMarshaler); ok {
			jsonBytes, err := marshaler.MarshalJSON()
			if err != nil {
				return nil, err
			}
			doc, err := JSONToYAML(jsonBytes)
			if err != nil {
				return nil, err
			}
			node, err := e.encodeDocument(doc)
			if err != nil {
				return nil, err
			}
			return node, nil
		}
	}

	return nil, errors.New("does not implemented Marshaler")
}

func (e *Encoder) encodeValue(ctx context.Context, v reflect.Value, column int) (ast.Node, error) {
	if e.isInvalidValue(v) {
		return e.encodeNil(), nil
	}
	if e.canEncodeByMarshaler(v) {
		node, err := e.encodeByMarshaler(ctx, v, column)
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeUint(v.Uint()), nil
	case reflect.Float32:
		return e.encodeFloat(v.Float(), 32), nil
	case reflect.Float64:
		return e.encodeFloat(v.Float(), 64), nil
	case reflect.Ptr:
		if value := e.encodePtrAnchor(v, column); value != nil {
			return value, nil
		}
		return e.encodeValue(ctx, v.Elem(), column)
	case reflect.Interface:
		return e.encodeValue(ctx, v.Elem(), column)
	case reflect.String:
		return e.encodeString(v.String(), column), nil
	case reflect.Bool:
		return e.encodeBool(v.Bool()), nil
	case reflect.Slice:
		if mapSlice, ok := v.Interface().(MapSlice); ok {
			return e.encodeMapSlice(ctx, mapSlice, column)
		}
		if value := e.encodePtrAnchor(v, column); value != nil {
			return value, nil
		}
		return e.encodeSlice(ctx, v)
	case reflect.Array:
		return e.encodeArray(ctx, v)
	case reflect.Struct:
		if v.CanInterface() {
			if mapItem, ok := v.Interface().(MapItem); ok {
				return e.encodeMapItem(ctx, mapItem, column)
			}
			if t, ok := v.Interface().(time.Time); ok {
				return e.encodeTime(t, column), nil
			}
		}
		return e.encodeStruct(ctx, v, column)
	case reflect.Map:
		if value := e.encodePtrAnchor(v, column); value != nil {
			return value, nil
		}
		return e.encodeMap(ctx, v, column)
	default:
		return nil, fmt.Errorf("unknown value type %s", v.Type().String())
	}
}

func (e *Encoder) encodePtrAnchor(v reflect.Value, column int) ast.Node {
	anchorName, exists := e.getAnchor(v.Pointer())
	if !exists {
		return nil
	}
	aliasName := anchorName
	alias := ast.Alias(token.New("*", "*", e.pos(column)))
	alias.Value = ast.String(token.New(aliasName, aliasName, e.pos(column)))
	e.setSmartAlias(aliasName, v.Pointer())
	return alias
}

func (e *Encoder) pos(column int) *token.Position {
	return &token.Position{
		Line:        e.line,
		Column:      column,
		Offset:      e.offset,
		IndentNum:   e.indentNum,
		IndentLevel: e.indentLevel,
	}
}

func (e *Encoder) encodeNil() *ast.NullNode {
	value := "null"
	return ast.Null(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeInt(v int64) *ast.IntegerNode {
	value := strconv.FormatInt(v, 10)
	return ast.Integer(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeUint(v uint64) *ast.IntegerNode {
	value := strconv.FormatUint(v, 10)
	return ast.Integer(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeFloat(v float64, bitSize int) ast.Node {
	if v == math.Inf(0) {
		value := ".inf"
		return ast.Infinity(token.New(value, value, e.pos(e.column)))
	} else if v == math.Inf(-1) {
		value := "-.inf"
		return ast.Infinity(token.New(value, value, e.pos(e.column)))
	} else if math.IsNaN(v) {
		value := ".nan"
		return ast.Nan(token.New(value, value, e.pos(e.column)))
	}
	value := strconv.FormatFloat(v, 'g', -1, bitSize)
	if !strings.Contains(value, ".") && !strings.Contains(value, "e") {
		if e.autoInt {
			return ast.Integer(token.New(value, value, e.pos(e.column)))
		}
		// append x.0 suffix to keep float value context
		value = fmt.Sprintf("%s.0", value)
	}
	return ast.Float(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) isNeedQuoted(v string) bool {
	if e.isJSONStyle {
		return true
	}
	if e.useLiteralStyleIfMultiline && strings.ContainsAny(v, "\n\r") {
		return false
	}
	if e.isFlowStyle && strings.ContainsAny(v, `]},'"`) {
		return true
	}
	if e.isFlowStyle {
		for i := 0; i < len(v); i++ {
			if v[i] != ':' {
				continue
			}
			if i+1 < len(v) && v[i+1] == '/' {
				continue
			}
			return true
		}
	}
	if token.IsNeedQuoted(v) {
		return true
	}
	return false
}

func (e *Encoder) encodeString(v string, column int) *ast.StringNode {
	if e.isNeedQuoted(v) {
		if e.singleQuote {
			v = quoteWith(v, '\'')
		} else {
			v = strconv.Quote(v)
		}
	}
	return ast.String(token.New(v, v, e.pos(column)))
}

func (e *Encoder) encodeBool(v bool) *ast.BoolNode {
	value := strconv.FormatBool(v)
	return ast.Bool(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeSlice(ctx context.Context, value reflect.Value) (*ast.SequenceNode, error) {
	if e.indentSequence {
		e.column += e.indentNum
		defer func() { e.column -= e.indentNum }()
	}
	column := e.column
	sequence := ast.Sequence(token.New("-", "-", e.pos(column)), e.isFlowStyle)
	for i := 0; i < value.Len(); i++ {
		node, err := e.encodeValue(ctx, value.Index(i), column)
		if err != nil {
			return nil, err
		}
		sequence.Values = append(sequence.Values, node)
	}
	return sequence, nil
}

func (e *Encoder) encodeArray(ctx context.Context, value reflect.Value) (*ast.SequenceNode, error) {
	if e.indentSequence {
		e.column += e.indentNum
		defer func() { e.column -= e.indentNum }()
	}
	column := e.column
	sequence := ast.Sequence(token.New("-", "-", e.pos(column)), e.isFlowStyle)
	for i := 0; i < value.Len(); i++ {
		node, err := e.encodeValue(ctx, value.Index(i), column)
		if err != nil {
			return nil, err
		}
		sequence.Values = append(sequence.Values, node)
	}
	return sequence, nil
}

func (e *Encoder) encodeMapItem(ctx context.Context, item MapItem, column int) (*ast.MappingValueNode, error) {
	k := reflect.ValueOf(item.Key)
	v := reflect.ValueOf(item.Value)
	value, err := e.encodeValue(ctx, v, column)
	if err != nil {
		return nil, err
	}
	if e.isMapNode(value) {
		value.AddColumn(e.indentNum)
	}
	if e.isTagAndMapNode(value) {
		value.AddColumn(e.indentNum)
	}
	return ast.MappingValue(
		token.New("", "", e.pos(column)),
		e.encodeString(k.Interface().(string), column),
		value,
	), nil
}

func (e *Encoder) encodeMapSlice(ctx context.Context, value MapSlice, column int) (*ast.MappingNode, error) {
	node := ast.Mapping(token.New("", "", e.pos(column)), e.isFlowStyle)
	for _, item := range value {
		encoded, err := e.encodeMapItem(ctx, item, column)
		if err != nil {
			return nil, err
		}
		node.Values = append(node.Values, encoded)
	}
	return node, nil
}

func (e *Encoder) isMapNode(node ast.Node) bool {
	_, ok := node.(ast.MapNode)
	return ok
}

func (e *Encoder) isTagAndMapNode(node ast.Node) bool {
	tn, ok := node.(*ast.TagNode)
	return ok && e.isMapNode(tn.Value)
}

func (e *Encoder) encodeMap(ctx context.Context, value reflect.Value, column int) (ast.Node, error) {
	node := ast.Mapping(token.New("", "", e.pos(column)), e.isFlowStyle)
	keys := make([]interface{}, len(value.MapKeys()))
	for i, k := range value.MapKeys() {
		keys[i] = k.Interface()
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})
	for _, key := range keys {
		k := reflect.ValueOf(key)
		v := value.MapIndex(k)
		encoded, err := e.encodeValue(ctx, v, column)
		if err != nil {
			return nil, err
		}
		if e.isMapNode(encoded) {
			encoded.AddColumn(e.indentNum)
		}
		if e.isTagAndMapNode(encoded) {
			encoded.AddColumn(e.indentNum)
		}
		keyText := fmt.Sprint(key)
		vRef := e.toPointer(v)

		// during the second encoding, an anchor is assigned if it is found to be used by an alias.
		if aliasName, exists := e.getSmartAlias(vRef); exists {
			anchorName := aliasName
			anchorNode := ast.Anchor(token.New("&", "&", e.pos(column)))
			anchorNode.Name = ast.String(token.New(anchorName, anchorName, e.pos(column)))
			anchorNode.Value = encoded
			encoded = anchorNode
		}
		node.Values = append(node.Values, ast.MappingValue(
			nil,
			e.encodeString(keyText, column),
			encoded,
		))
		e.setSmartAnchor(vRef, keyText)
	}
	return node, nil
}

// IsZeroer is used to check whether an object is zero to determine
// whether it should be omitted when marshaling with the omitempty flag.
// One notable implementation is time.Time.
type IsZeroer interface {
	IsZero() bool
}

func (e *Encoder) isOmittedByOmitZero(v reflect.Value) bool {
	kind := v.Kind()
	if z, ok := v.Interface().(IsZeroer); ok {
		if (kind == reflect.Ptr || kind == reflect.Interface) && v.IsNil() {
			return true
		}
		return z.IsZero()
	}
	switch kind {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		for i := v.NumField() - 1; i >= 0; i-- {
			if vt.Field(i).PkgPath != "" {
				continue // private field
			}
			if !e.isOmittedByOmitZero(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

func (e *Encoder) isOmittedByOmitEmptyOption(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	}
	return false
}

// The current implementation of the omitempty tag combines the functionality of encoding/json's omitempty and omitzero tags.
// This stems from a historical decision to respect the implementation of gopkg.in/yaml.v2, but it has caused confusion,
// so we are working to integrate it into the functionality of encoding/json. (However, this will take some time.)
// In the current implementation, in addition to the exclusion conditions of omitempty,
// if a type implements IsZero, that implementation will be used.
// Furthermore, for non-pointer structs, if all fields are eligible for exclusion,
// the struct itself will also be excluded. These behaviors are originally the functionality of omitzero.
func (e *Encoder) isOmittedByOmitEmptyTag(v reflect.Value) bool {
	kind := v.Kind()
	if z, ok := v.Interface().(IsZeroer); ok {
		if (kind == reflect.Ptr || kind == reflect.Interface) && v.IsNil() {
			return true
		}
		return z.IsZero()
	}
	switch kind {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		for i := v.NumField() - 1; i >= 0; i-- {
			if vt.Field(i).PkgPath != "" {
				continue // private field
			}
			if !e.isOmittedByOmitEmptyTag(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

func (e *Encoder) encodeTime(v time.Time, column int) *ast.StringNode {
	value := v.Format(time.RFC3339Nano)
	if e.isJSONStyle {
		value = strconv.Quote(value)
	}
	return ast.String(token.New(value, value, e.pos(column)))
}

func (e *Encoder) encodeDuration(v time.Duration, column int) *ast.StringNode {
	value := v.String()
	if e.isJSONStyle {
		value = strconv.Quote(value)
	}
	return ast.String(token.New(value, value, e.pos(column)))
}

func (e *Encoder) encodeAnchor(anchorName string, value ast.Node, fieldValue reflect.Value, column int) (*ast.AnchorNode, error) {
	anchorNode := ast.Anchor(token.New("&", "&", e.pos(column)))
	anchorNode.Name = ast.String(token.New(anchorName, anchorName, e.pos(column)))
	anchorNode.Value = value
	if e.anchorCallback != nil {
		if err := e.anchorCallback(anchorNode, fieldValue.Interface()); err != nil {
			return nil, err
		}
		if snode, ok := anchorNode.Name.(*ast.StringNode); ok {
			anchorName = snode.Value
		}
	}
	if fieldValue.Kind() == reflect.Ptr {
		e.setAnchor(fieldValue.Pointer(), anchorName)
	}
	return anchorNode, nil
}

func (e *Encoder) encodeStruct(ctx context.Context, value reflect.Value, column int) (ast.Node, error) {
	node := ast.Mapping(token.New("", "", e.pos(column)), e.isFlowStyle)
	structType := value.Type()
	fieldMap, err := structFieldMap(structType)
	if err != nil {
		return nil, err
	}
	hasInlineAnchorField := false
	var inlineAnchorValue reflect.Value
	for i := 0; i < value.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		fieldValue := value.FieldByName(field.Name)
		sf := fieldMap[field.Name]
		if (e.omitZero || sf.IsOmitZero) && e.isOmittedByOmitZero(fieldValue) {
			// omit encoding by omitzero tag or OmitZero option.
			continue
		}
		if e.omitEmpty && e.isOmittedByOmitEmptyOption(fieldValue) {
			// omit encoding by OmitEmpty option.
			continue
		}
		if sf.IsOmitEmpty && e.isOmittedByOmitEmptyTag(fieldValue) {
			// omit encoding by omitempty tag.
			continue
		}
		ve := e
		if !e.isFlowStyle && sf.IsFlow {
			ve = &Encoder{}
			*ve = *e
			ve.isFlowStyle = true
		}
		encoded, err := ve.encodeValue(ctx, fieldValue, column)
		if err != nil {
			return nil, err
		}
		if e.isMapNode(encoded) {
			encoded.AddColumn(e.indentNum)
		}
		var key ast.MapKeyNode = e.encodeString(sf.RenderName, column)
		switch {
		case encoded.Type() == ast.AliasType:
			if aliasName := sf.AliasName; aliasName != "" {
				alias, ok := encoded.(*ast.AliasNode)
				if !ok {
					return nil, errors.ErrUnexpectedNodeType(encoded.Type(), ast.AliasType, encoded.GetToken())
				}
				got := alias.Value.String()
				if aliasName != got {
					return nil, fmt.Errorf("expected alias name is %q but got %q", aliasName, got)
				}
			}
			if sf.IsInline {
				// if both used alias and inline, output `<<: *alias`
				key = ast.MergeKey(token.New("<<", "<<", e.pos(column)))
			}
		case sf.AnchorName != "":
			anchorNode, err := e.encodeAnchor(sf.AnchorName, encoded, fieldValue, column)
			if err != nil {
				return nil, err
			}
			encoded = anchorNode
		case sf.IsInline:
			isAutoAnchor := sf.IsAutoAnchor
			if !hasInlineAnchorField {
				hasInlineAnchorField = isAutoAnchor
			}
			if isAutoAnchor {
				inlineAnchorValue = fieldValue
			}
			mapNode, ok := encoded.(ast.MapNode)
			if !ok {
				// if an inline field is null, skip encoding it
				if _, ok := encoded.(*ast.NullNode); ok {
					continue
				}
				return nil, errors.New("inline value is must be map or struct type")
			}
			mapIter := mapNode.MapRange()
			for mapIter.Next() {
				mapKey := mapIter.Key()
				mapValue := mapIter.Value()
				keyName := mapKey.GetToken().Value
				if fieldMap.isIncludedRenderName(keyName) {
					// if declared the same key name, skip encoding this field
					continue
				}
				mapKey.AddColumn(-e.indentNum)
				mapValue.AddColumn(-e.indentNum)
				node.Values = append(node.Values, ast.MappingValue(nil, mapKey, mapValue))
			}
			continue
		case sf.IsAutoAnchor:
			anchorNode, err := e.encodeAnchor(sf.RenderName, encoded, fieldValue, column)
			if err != nil {
				return nil, err
			}
			encoded = anchorNode
		}
		node.Values = append(node.Values, ast.MappingValue(nil, key, encoded))
	}
	if hasInlineAnchorField {
		node.AddColumn(e.indentNum)
		anchorName := "anchor"
		anchorNode := ast.Anchor(token.New("&", "&", e.pos(column)))
		anchorNode.Name = ast.String(token.New(anchorName, anchorName, e.pos(column)))
		anchorNode.Value = node
		if e.anchorCallback != nil {
			if err := e.anchorCallback(anchorNode, value.Addr().Interface()); err != nil {
				return nil, err
			}
			if snode, ok := anchorNode.Name.(*ast.StringNode); ok {
				anchorName = snode.Value
			}
		}
		if inlineAnchorValue.Kind() == reflect.Ptr {
			e.setAnchor(inlineAnchorValue.Pointer(), anchorName)
		}
		return anchorNode, nil
	}
	return node, nil
}

func (e *Encoder) toPointer(v reflect.Value) uintptr {
	if e.isInvalidValue(v) {
		return 0
	}

	switch v.Type().Kind() {
	case reflect.Ptr:
		return v.Pointer()
	case reflect.Interface:
		return e.toPointer(v.Elem())
	case reflect.Slice:
		return v.Pointer()
	case reflect.Map:
		return v.Pointer()
	}
	return 0
}

func (e *Encoder) clearSmartAnchorRef() {
	if !e.enableSmartAnchor {
		return
	}
	e.anchorRefToName = make(map[uintptr]string)
	e.anchorNameMap = make(map[string]struct{})
}

func (e *Encoder) setSmartAnchor(ptr uintptr, name string) {
	if !e.enableSmartAnchor {
		return
	}
	e.setAnchor(ptr, e.generateAnchorName(name))
}

func (e *Encoder) setAnchor(ptr uintptr, name string) {
	if ptr == 0 {
		return
	}
	if name == "" {
		return
	}
	e.anchorRefToName[ptr] = name
	e.anchorNameMap[name] = struct{}{}
}

func (e *Encoder) generateAnchorName(base string) string {
	if _, exists := e.anchorNameMap[base]; !exists {
		return base
	}
	for i := 1; i < 100; i++ {
		name := base + strconv.Itoa(i)
		if _, exists := e.anchorNameMap[name]; exists {
			continue
		}
		return name
	}
	return ""
}

func (e *Encoder) getAnchor(ref uintptr) (string, bool) {
	anchorName, exists := e.anchorRefToName[ref]
	return anchorName, exists
}

func (e *Encoder) setSmartAlias(name string, ref uintptr) {
	if !e.enableSmartAnchor {
		return
	}
	e.aliasRefToName[ref] = name
}

func (e *Encoder) getSmartAlias(ref uintptr) (string, bool) {
	if !e.enableSmartAnchor {
		return "", false
	}
	aliasName, exists := e.aliasRefToName[ref]
	return aliasName, exists
}
