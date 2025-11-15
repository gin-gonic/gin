package decoder

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type PathString string

func (s PathString) Build() (*Path, error) {
	builder := new(PathBuilder)
	return builder.Build([]rune(s))
}

type PathBuilder struct {
	root                    PathNode
	node                    PathNode
	singleQuotePathSelector bool
	doubleQuotePathSelector bool
}

func (b *PathBuilder) Build(buf []rune) (*Path, error) {
	node, err := b.build(buf)
	if err != nil {
		return nil, err
	}
	return &Path{
		node:                    node,
		RootSelectorOnly:        node == nil,
		SingleQuotePathSelector: b.singleQuotePathSelector,
		DoubleQuotePathSelector: b.doubleQuotePathSelector,
	}, nil
}

func (b *PathBuilder) build(buf []rune) (PathNode, error) {
	if len(buf) == 0 {
		return nil, errors.ErrEmptyPath()
	}
	if buf[0] != '$' {
		return nil, errors.ErrInvalidPath("JSON Path must start with a $ character")
	}
	if len(buf) == 1 {
		return nil, nil
	}
	buf = buf[1:]
	offset, err := b.buildNext(buf)
	if err != nil {
		return nil, err
	}
	if len(buf) > offset {
		return nil, errors.ErrInvalidPath("remain invalid path %q", buf[offset:])
	}
	return b.root, nil
}

func (b *PathBuilder) buildNextCharIfExists(buf []rune, cursor int) (int, error) {
	if len(buf) > cursor {
		offset, err := b.buildNext(buf[cursor:])
		if err != nil {
			return 0, err
		}
		return cursor + 1 + offset, nil
	}
	return cursor, nil
}

func (b *PathBuilder) buildNext(buf []rune) (int, error) {
	switch buf[0] {
	case '.':
		if len(buf) == 1 {
			return 0, errors.ErrInvalidPath("JSON Path ends with dot character")
		}
		offset, err := b.buildSelector(buf[1:])
		if err != nil {
			return 0, err
		}
		return offset + 1, nil
	case '[':
		if len(buf) == 1 {
			return 0, errors.ErrInvalidPath("JSON Path ends with left bracket character")
		}
		offset, err := b.buildIndex(buf[1:])
		if err != nil {
			return 0, err
		}
		return offset + 1, nil
	default:
		return 0, errors.ErrInvalidPath("expect dot or left bracket character. but found %c character", buf[0])
	}
}

func (b *PathBuilder) buildSelector(buf []rune) (int, error) {
	switch buf[0] {
	case '.':
		if len(buf) == 1 {
			return 0, errors.ErrInvalidPath("JSON Path ends with double dot character")
		}
		offset, err := b.buildPathRecursive(buf[1:])
		if err != nil {
			return 0, err
		}
		return 1 + offset, nil
	case '[', ']', '$', '*':
		return 0, errors.ErrInvalidPath("found invalid path character %c after dot", buf[0])
	}
	for cursor := 0; cursor < len(buf); cursor++ {
		switch buf[cursor] {
		case '$', '*', ']':
			return 0, errors.ErrInvalidPath("found %c character in field selector context", buf[cursor])
		case '.':
			if cursor+1 >= len(buf) {
				return 0, errors.ErrInvalidPath("JSON Path ends with dot character")
			}
			selector := buf[:cursor]
			b.addSelectorNode(string(selector))
			offset, err := b.buildSelector(buf[cursor+1:])
			if err != nil {
				return 0, err
			}
			return cursor + 1 + offset, nil
		case '[':
			if cursor+1 >= len(buf) {
				return 0, errors.ErrInvalidPath("JSON Path ends with left bracket character")
			}
			selector := buf[:cursor]
			b.addSelectorNode(string(selector))
			offset, err := b.buildIndex(buf[cursor+1:])
			if err != nil {
				return 0, err
			}
			return cursor + 1 + offset, nil
		case '"':
			if cursor+1 >= len(buf) {
				return 0, errors.ErrInvalidPath("JSON Path ends with double quote character")
			}
			offset, err := b.buildQuoteSelector(buf[cursor+1:], DoubleQuotePathSelector)
			if err != nil {
				return 0, err
			}
			return cursor + 1 + offset, nil
		}
	}
	b.addSelectorNode(string(buf))
	return len(buf), nil
}

func (b *PathBuilder) buildQuoteSelector(buf []rune, sel QuotePathSelector) (int, error) {
	switch buf[0] {
	case '[', ']', '$', '.', '*', '\'', '"':
		return 0, errors.ErrInvalidPath("found invalid path character %c after quote", buf[0])
	}
	for cursor := 0; cursor < len(buf); cursor++ {
		switch buf[cursor] {
		case '\'':
			if sel != SingleQuotePathSelector {
				return 0, errors.ErrInvalidPath("found double quote character in field selector with single quote context")
			}
			if len(buf) <= cursor+1 {
				return 0, errors.ErrInvalidPath("JSON Path ends with single quote character in field selector context")
			}
			if buf[cursor+1] != ']' {
				return 0, errors.ErrInvalidPath("expect right bracket for field selector with single quote but found %c", buf[cursor+1])
			}
			selector := buf[:cursor]
			b.addSelectorNode(string(selector))
			b.singleQuotePathSelector = true
			return b.buildNextCharIfExists(buf, cursor+2)
		case '"':
			if sel != DoubleQuotePathSelector {
				return 0, errors.ErrInvalidPath("found single quote character in field selector with double quote context")
			}
			selector := buf[:cursor]
			b.addSelectorNode(string(selector))
			b.doubleQuotePathSelector = true
			return b.buildNextCharIfExists(buf, cursor+1)
		}
	}
	return 0, errors.ErrInvalidPath("couldn't find quote character in selector quote path context")
}

func (b *PathBuilder) buildPathRecursive(buf []rune) (int, error) {
	switch buf[0] {
	case '.', '[', ']', '$', '*':
		return 0, errors.ErrInvalidPath("found invalid path character %c after double dot", buf[0])
	}
	for cursor := 0; cursor < len(buf); cursor++ {
		switch buf[cursor] {
		case '$', '*', ']':
			return 0, errors.ErrInvalidPath("found %c character in field selector context", buf[cursor])
		case '.':
			if cursor+1 >= len(buf) {
				return 0, errors.ErrInvalidPath("JSON Path ends with dot character")
			}
			selector := buf[:cursor]
			b.addRecursiveNode(string(selector))
			offset, err := b.buildSelector(buf[cursor+1:])
			if err != nil {
				return 0, err
			}
			return cursor + 1 + offset, nil
		case '[':
			if cursor+1 >= len(buf) {
				return 0, errors.ErrInvalidPath("JSON Path ends with left bracket character")
			}
			selector := buf[:cursor]
			b.addRecursiveNode(string(selector))
			offset, err := b.buildIndex(buf[cursor+1:])
			if err != nil {
				return 0, err
			}
			return cursor + 1 + offset, nil
		}
	}
	b.addRecursiveNode(string(buf))
	return len(buf), nil
}

func (b *PathBuilder) buildIndex(buf []rune) (int, error) {
	switch buf[0] {
	case '.', '[', ']', '$':
		return 0, errors.ErrInvalidPath("found invalid path character %c after left bracket", buf[0])
	case '\'':
		if len(buf) == 1 {
			return 0, errors.ErrInvalidPath("JSON Path ends with single quote character")
		}
		offset, err := b.buildQuoteSelector(buf[1:], SingleQuotePathSelector)
		if err != nil {
			return 0, err
		}
		return 1 + offset, nil
	case '*':
		if len(buf) == 1 {
			return 0, errors.ErrInvalidPath("JSON Path ends with star character")
		}
		if buf[1] != ']' {
			return 0, errors.ErrInvalidPath("expect right bracket character for index all path but found %c character", buf[1])
		}
		b.addIndexAllNode()
		offset := len("*]")
		if len(buf) > 2 {
			buildOffset, err := b.buildNext(buf[2:])
			if err != nil {
				return 0, err
			}
			return offset + buildOffset, nil
		}
		return offset, nil
	}

	for cursor := 0; cursor < len(buf); cursor++ {
		switch buf[cursor] {
		case ']':
			index, err := strconv.ParseInt(string(buf[:cursor]), 10, 64)
			if err != nil {
				return 0, errors.ErrInvalidPath("%q is unexpected index path", buf[:cursor])
			}
			b.addIndexNode(int(index))
			return b.buildNextCharIfExists(buf, cursor+1)
		}
	}
	return 0, errors.ErrInvalidPath("couldn't find right bracket character in index path context")
}

func (b *PathBuilder) addIndexAllNode() {
	node := newPathIndexAllNode()
	if b.root == nil {
		b.root = node
		b.node = node
	} else {
		b.node = b.node.chain(node)
	}
}

func (b *PathBuilder) addRecursiveNode(selector string) {
	node := newPathRecursiveNode(selector)
	if b.root == nil {
		b.root = node
		b.node = node
	} else {
		b.node = b.node.chain(node)
	}
}

func (b *PathBuilder) addSelectorNode(name string) {
	node := newPathSelectorNode(name)
	if b.root == nil {
		b.root = node
		b.node = node
	} else {
		b.node = b.node.chain(node)
	}
}

func (b *PathBuilder) addIndexNode(idx int) {
	node := newPathIndexNode(idx)
	if b.root == nil {
		b.root = node
		b.node = node
	} else {
		b.node = b.node.chain(node)
	}
}

type QuotePathSelector int

const (
	SingleQuotePathSelector QuotePathSelector = 1
	DoubleQuotePathSelector QuotePathSelector = 2
)

type Path struct {
	node                    PathNode
	RootSelectorOnly        bool
	SingleQuotePathSelector bool
	DoubleQuotePathSelector bool
}

func (p *Path) Field(sel string) (PathNode, bool, error) {
	if p.node == nil {
		return nil, false, nil
	}
	return p.node.Field(sel)
}

func (p *Path) Get(src, dst reflect.Value) error {
	if p.node == nil {
		return nil
	}
	return p.node.Get(src, dst)
}

func (p *Path) String() string {
	if p.node == nil {
		return "$"
	}
	return p.node.String()
}

type PathNode interface {
	fmt.Stringer
	Index(idx int) (PathNode, bool, error)
	Field(fieldName string) (PathNode, bool, error)
	Get(src, dst reflect.Value) error
	chain(PathNode) PathNode
	target() bool
	single() bool
}

type BasePathNode struct {
	child PathNode
}

func (n *BasePathNode) chain(node PathNode) PathNode {
	n.child = node
	return node
}

func (n *BasePathNode) target() bool {
	return n.child == nil
}

func (n *BasePathNode) single() bool {
	return true
}

type PathSelectorNode struct {
	*BasePathNode
	selector string
}

func newPathSelectorNode(selector string) *PathSelectorNode {
	return &PathSelectorNode{
		BasePathNode: &BasePathNode{},
		selector:     selector,
	}
}

func (n *PathSelectorNode) Index(idx int) (PathNode, bool, error) {
	return nil, false, &errors.PathError{}
}

func (n *PathSelectorNode) Field(fieldName string) (PathNode, bool, error) {
	if n.selector == fieldName {
		return n.child, true, nil
	}
	return nil, false, nil
}

func (n *PathSelectorNode) Get(src, dst reflect.Value) error {
	switch src.Type().Kind() {
	case reflect.Map:
		iter := src.MapRange()
		for iter.Next() {
			key, ok := iter.Key().Interface().(string)
			if !ok {
				return fmt.Errorf("invalid map key type %T", src.Type().Key())
			}
			child, found, err := n.Field(key)
			if err != nil {
				return err
			}
			if found {
				if child != nil {
					return child.Get(iter.Value(), dst)
				}
				return AssignValue(iter.Value(), dst)
			}
		}
	case reflect.Struct:
		typ := src.Type()
		for i := 0; i < typ.Len(); i++ {
			tag := runtime.StructTagFromField(typ.Field(i))
			child, found, err := n.Field(tag.Key)
			if err != nil {
				return err
			}
			if found {
				if child != nil {
					return child.Get(src.Field(i), dst)
				}
				return AssignValue(src.Field(i), dst)
			}
		}
	case reflect.Ptr:
		return n.Get(src.Elem(), dst)
	case reflect.Interface:
		return n.Get(reflect.ValueOf(src.Interface()), dst)
	case reflect.Float64, reflect.String, reflect.Bool:
		return AssignValue(src, dst)
	}
	return fmt.Errorf("failed to get %s value from %s", n.selector, src.Type())
}

func (n *PathSelectorNode) String() string {
	s := fmt.Sprintf(".%s", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type PathIndexNode struct {
	*BasePathNode
	selector int
}

func newPathIndexNode(selector int) *PathIndexNode {
	return &PathIndexNode{
		BasePathNode: &BasePathNode{},
		selector:     selector,
	}
}

func (n *PathIndexNode) Index(idx int) (PathNode, bool, error) {
	if n.selector == idx {
		return n.child, true, nil
	}
	return nil, false, nil
}

func (n *PathIndexNode) Field(fieldName string) (PathNode, bool, error) {
	return nil, false, &errors.PathError{}
}

func (n *PathIndexNode) Get(src, dst reflect.Value) error {
	switch src.Type().Kind() {
	case reflect.Array, reflect.Slice:
		if src.Len() > n.selector {
			if n.child != nil {
				return n.child.Get(src.Index(n.selector), dst)
			}
			return AssignValue(src.Index(n.selector), dst)
		}
	case reflect.Ptr:
		return n.Get(src.Elem(), dst)
	case reflect.Interface:
		return n.Get(reflect.ValueOf(src.Interface()), dst)
	}
	return fmt.Errorf("failed to get [%d] value from %s", n.selector, src.Type())
}

func (n *PathIndexNode) String() string {
	s := fmt.Sprintf("[%d]", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type PathIndexAllNode struct {
	*BasePathNode
}

func newPathIndexAllNode() *PathIndexAllNode {
	return &PathIndexAllNode{
		BasePathNode: &BasePathNode{},
	}
}

func (n *PathIndexAllNode) Index(idx int) (PathNode, bool, error) {
	return n.child, true, nil
}

func (n *PathIndexAllNode) Field(fieldName string) (PathNode, bool, error) {
	return nil, false, &errors.PathError{}
}

func (n *PathIndexAllNode) Get(src, dst reflect.Value) error {
	switch src.Type().Kind() {
	case reflect.Array, reflect.Slice:
		var arr []interface{}
		for i := 0; i < src.Len(); i++ {
			var v interface{}
			rv := reflect.ValueOf(&v)
			if n.child != nil {
				if err := n.child.Get(src.Index(i), rv); err != nil {
					return err
				}
			} else {
				if err := AssignValue(src.Index(i), rv); err != nil {
					return err
				}
			}
			arr = append(arr, v)
		}
		if err := AssignValue(reflect.ValueOf(arr), dst); err != nil {
			return err
		}
		return nil
	case reflect.Ptr:
		return n.Get(src.Elem(), dst)
	case reflect.Interface:
		return n.Get(reflect.ValueOf(src.Interface()), dst)
	}
	return fmt.Errorf("failed to get all value from %s", src.Type())
}

func (n *PathIndexAllNode) String() string {
	s := "[*]"
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type PathRecursiveNode struct {
	*BasePathNode
	selector string
}

func newPathRecursiveNode(selector string) *PathRecursiveNode {
	node := newPathSelectorNode(selector)
	return &PathRecursiveNode{
		BasePathNode: &BasePathNode{
			child: node,
		},
		selector: selector,
	}
}

func (n *PathRecursiveNode) Field(fieldName string) (PathNode, bool, error) {
	if n.selector == fieldName {
		return n.child, true, nil
	}
	return nil, false, nil
}

func (n *PathRecursiveNode) Index(_ int) (PathNode, bool, error) {
	return n, true, nil
}

func valueToSliceValue(v interface{}) []interface{} {
	rv := reflect.ValueOf(v)
	ret := []interface{}{}
	if rv.Type().Kind() == reflect.Slice || rv.Type().Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			ret = append(ret, rv.Index(i).Interface())
		}
		return ret
	}
	return []interface{}{v}
}

func (n *PathRecursiveNode) Get(src, dst reflect.Value) error {
	if n.child == nil {
		return fmt.Errorf("failed to get by recursive path ..%s", n.selector)
	}
	var arr []interface{}
	switch src.Type().Kind() {
	case reflect.Map:
		iter := src.MapRange()
		for iter.Next() {
			key, ok := iter.Key().Interface().(string)
			if !ok {
				return fmt.Errorf("invalid map key type %T", src.Type().Key())
			}
			child, found, err := n.Field(key)
			if err != nil {
				return err
			}
			if found {
				var v interface{}
				rv := reflect.ValueOf(&v)
				_ = child.Get(iter.Value(), rv)
				arr = append(arr, valueToSliceValue(v)...)
			} else {
				var v interface{}
				rv := reflect.ValueOf(&v)
				_ = n.Get(iter.Value(), rv)
				if v != nil {
					arr = append(arr, valueToSliceValue(v)...)
				}
			}
		}
		_ = AssignValue(reflect.ValueOf(arr), dst)
		return nil
	case reflect.Struct:
		typ := src.Type()
		for i := 0; i < typ.Len(); i++ {
			tag := runtime.StructTagFromField(typ.Field(i))
			child, found, err := n.Field(tag.Key)
			if err != nil {
				return err
			}
			if found {
				var v interface{}
				rv := reflect.ValueOf(&v)
				_ = child.Get(src.Field(i), rv)
				arr = append(arr, valueToSliceValue(v)...)
			} else {
				var v interface{}
				rv := reflect.ValueOf(&v)
				_ = n.Get(src.Field(i), rv)
				if v != nil {
					arr = append(arr, valueToSliceValue(v)...)
				}
			}
		}
		_ = AssignValue(reflect.ValueOf(arr), dst)
		return nil
	case reflect.Array, reflect.Slice:
		for i := 0; i < src.Len(); i++ {
			var v interface{}
			rv := reflect.ValueOf(&v)
			_ = n.Get(src.Index(i), rv)
			if v != nil {
				arr = append(arr, valueToSliceValue(v)...)
			}
		}
		_ = AssignValue(reflect.ValueOf(arr), dst)
		return nil
	case reflect.Ptr:
		return n.Get(src.Elem(), dst)
	case reflect.Interface:
		return n.Get(reflect.ValueOf(src.Interface()), dst)
	}
	return fmt.Errorf("failed to get %s value from %s", n.selector, src.Type())
}

func (n *PathRecursiveNode) String() string {
	s := fmt.Sprintf("..%s", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}
