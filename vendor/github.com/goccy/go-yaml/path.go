package yaml

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
)

// PathString create Path from string
//
// YAMLPath rule
// $     : the root object/element
// .     : child operator
// ..    : recursive descent
// [num] : object/element of array by number
// [*]   : all objects/elements for array.
//
// If you want to use reserved characters such as `.` and `*` as a key name,
// enclose them in single quotation as follows ( $.foo.'bar.baz-*'.hoge ).
// If you want to use a single quote with reserved characters, escape it with `\` ( $.foo.'bar.baz\'s value'.hoge ).
func PathString(s string) (*Path, error) {
	buf := []rune(s)
	length := len(buf)
	cursor := 0
	builder := &PathBuilder{}
	for cursor < length {
		c := buf[cursor]
		switch c {
		case '$':
			builder = builder.Root()
			cursor++
		case '.':
			b, buf, c, err := parsePathDot(builder, buf, cursor)
			if err != nil {
				return nil, err
			}
			length = len(buf)
			builder = b
			cursor = c
		case '[':
			b, buf, c, err := parsePathIndex(builder, buf, cursor)
			if err != nil {
				return nil, err
			}
			length = len(buf)
			builder = b
			cursor = c
		default:
			return nil, fmt.Errorf("invalid path at %d: %w", cursor, ErrInvalidPathString)
		}
	}
	return builder.Build(), nil
}

func parsePathRecursive(b *PathBuilder, buf []rune, cursor int) (*PathBuilder, []rune, int, error) {
	length := len(buf)
	cursor += 2 // skip .. characters
	start := cursor
	for ; cursor < length; cursor++ {
		c := buf[cursor]
		switch c {
		case '$':
			return nil, nil, 0, fmt.Errorf("specified '$' after '..' character: %w", ErrInvalidPathString)
		case '*':
			return nil, nil, 0, fmt.Errorf("specified '*' after '..' character: %w", ErrInvalidPathString)
		case '.', '[':
			goto end
		case ']':
			return nil, nil, 0, fmt.Errorf("specified ']' after '..' character: %w", ErrInvalidPathString)
		}
	}
end:
	if start == cursor {
		return nil, nil, 0, fmt.Errorf("not found recursive selector: %w", ErrInvalidPathString)
	}
	return b.Recursive(string(buf[start:cursor])), buf, cursor, nil
}

func parsePathDot(b *PathBuilder, buf []rune, cursor int) (*PathBuilder, []rune, int, error) {
	if b.root == nil || b.node == nil {
		return nil, nil, 0, fmt.Errorf("required '$' character at first: %w", ErrInvalidPathString)
	}
	length := len(buf)
	if cursor+1 < length && buf[cursor+1] == '.' {
		b, buf, c, err := parsePathRecursive(b, buf, cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		return b, buf, c, nil
	}
	cursor++ // skip . character
	start := cursor

	// if started single quote, looking for end single quote char
	if cursor < length && buf[cursor] == '\'' {
		return parseQuotedKey(b, buf, cursor)
	}
	for ; cursor < length; cursor++ {
		c := buf[cursor]
		switch c {
		case '$':
			return nil, nil, 0, fmt.Errorf("specified '$' after '.' character: %w", ErrInvalidPathString)
		case '*':
			return nil, nil, 0, fmt.Errorf("specified '*' after '.' character: %w", ErrInvalidPathString)
		case '.', '[':
			goto end
		case ']':
			return nil, nil, 0, fmt.Errorf("specified ']' after '.' character: %w", ErrInvalidPathString)
		}
	}
end:
	if start == cursor {
		return nil, nil, 0, fmt.Errorf("could not find by empty key: %w", ErrInvalidPathString)
	}
	return b.child(string(buf[start:cursor])), buf, cursor, nil
}

func parseQuotedKey(b *PathBuilder, buf []rune, cursor int) (*PathBuilder, []rune, int, error) {
	if b.root == nil || b.node == nil {
		return nil, nil, 0, fmt.Errorf("required '$' character at first: %w", ErrInvalidPathString)
	}

	cursor++ // skip single quote
	start := cursor
	length := len(buf)
	var foundEndDelim bool
	for ; cursor < length; cursor++ {
		switch buf[cursor] {
		case '\\':
			buf = append(append([]rune{}, buf[:cursor]...), buf[cursor+1:]...)
			length = len(buf)
		case '\'':
			foundEndDelim = true
			goto end
		}
	}
end:
	if !foundEndDelim {
		return nil, nil, 0, fmt.Errorf("could not find end delimiter for key: %w", ErrInvalidPathString)
	}
	if start == cursor {
		return nil, nil, 0, fmt.Errorf("could not find by empty key: %w", ErrInvalidPathString)
	}
	selector := buf[start:cursor]
	cursor++
	if cursor < length {
		switch buf[cursor] {
		case '$':
			return nil, nil, 0, fmt.Errorf("specified '$' after '.' character: %w", ErrInvalidPathString)
		case '*':
			return nil, nil, 0, fmt.Errorf("specified '*' after '.' character: %w", ErrInvalidPathString)
		case ']':
			return nil, nil, 0, fmt.Errorf("specified ']' after '.' character: %w", ErrInvalidPathString)
		}
	}
	return b.child(string(selector)), buf, cursor, nil
}

func parsePathIndex(b *PathBuilder, buf []rune, cursor int) (*PathBuilder, []rune, int, error) {
	if b.root == nil || b.node == nil {
		return nil, nil, 0, fmt.Errorf("required '$' character at first: %w", ErrInvalidPathString)
	}

	length := len(buf)
	cursor++ // skip '[' character
	if length <= cursor {
		return nil, nil, 0, fmt.Errorf("unexpected end of YAML Path: %w", ErrInvalidPathString)
	}
	c := buf[cursor]
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '*':
		start := cursor
		cursor++
		for ; cursor < length; cursor++ {
			c := buf[cursor]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				continue
			}
			break
		}
		if buf[cursor] != ']' {
			return nil, nil, 0, fmt.Errorf("invalid character %s at %d: %w", string(buf[cursor]), cursor, ErrInvalidPathString)
		}
		numOrAll := string(buf[start:cursor])
		if numOrAll == "*" {
			return b.IndexAll(), buf, cursor + 1, nil
		}
		num, err := strconv.ParseInt(numOrAll, 10, 64)
		if err != nil {
			return nil, nil, 0, err
		}
		return b.Index(uint(num)), buf, cursor + 1, nil
	}
	return nil, nil, 0, fmt.Errorf("invalid character %q at %d: %w", c, cursor, ErrInvalidPathString)
}

// Path represent YAMLPath ( like a JSONPath ).
type Path struct {
	node pathNode
}

// String path to text.
func (p *Path) String() string {
	return p.node.String()
}

// Read decode from r and set extracted value by YAMLPath to v.
func (p *Path) Read(r io.Reader, v interface{}) error {
	node, err := p.ReadNode(r)
	if err != nil {
		return err
	}
	if err := Unmarshal([]byte(node.String()), v); err != nil {
		return err
	}
	return nil
}

// ReadNode create AST from r and extract node by YAMLPath.
func (p *Path) ReadNode(r io.Reader) (ast.Node, error) {
	if p.node == nil {
		return nil, ErrInvalidPath
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}
	f, err := parser.ParseBytes(buf.Bytes(), 0)
	if err != nil {
		return nil, err
	}
	node, err := p.FilterFile(f)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Filter filter from target by YAMLPath and set it to v.
func (p *Path) Filter(target, v interface{}) error {
	b, err := Marshal(target)
	if err != nil {
		return err
	}
	if err := p.Read(bytes.NewBuffer(b), v); err != nil {
		return err
	}
	return nil
}

// FilterFile filter from ast.File by YAMLPath.
func (p *Path) FilterFile(f *ast.File) (ast.Node, error) {
	for _, doc := range f.Docs {
		node, err := p.FilterNode(doc.Body)
		if err != nil {
			return nil, err
		}
		if node != nil {
			return node, nil
		}
	}
	return nil, fmt.Errorf("failed to find path ( %s ): %w", p.node, ErrNotFoundNode)
}

// FilterNode filter from node by YAMLPath.
func (p *Path) FilterNode(node ast.Node) (ast.Node, error) {
	if node == nil {
		return nil, nil
	}
	n, err := p.node.filter(node)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// MergeFromReader merge YAML text into ast.File.
func (p *Path) MergeFromReader(dst *ast.File, src io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return err
	}
	file, err := parser.ParseBytes(buf.Bytes(), 0)
	if err != nil {
		return err
	}
	if err := p.MergeFromFile(dst, file); err != nil {
		return err
	}
	return nil
}

// MergeFromFile merge ast.File into ast.File.
func (p *Path) MergeFromFile(dst *ast.File, src *ast.File) error {
	base, err := p.FilterFile(dst)
	if err != nil {
		return err
	}
	for _, doc := range src.Docs {
		if err := ast.Merge(base, doc); err != nil {
			return err
		}
	}
	return nil
}

// MergeFromNode merge ast.Node into ast.File.
func (p *Path) MergeFromNode(dst *ast.File, src ast.Node) error {
	base, err := p.FilterFile(dst)
	if err != nil {
		return err
	}
	if err := ast.Merge(base, src); err != nil {
		return err
	}
	return nil
}

// ReplaceWithReader replace ast.File with io.Reader.
func (p *Path) ReplaceWithReader(dst *ast.File, src io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return err
	}
	file, err := parser.ParseBytes(buf.Bytes(), 0)
	if err != nil {
		return err
	}
	if err := p.ReplaceWithFile(dst, file); err != nil {
		return err
	}
	return nil
}

// ReplaceWithFile replace ast.File with ast.File.
func (p *Path) ReplaceWithFile(dst *ast.File, src *ast.File) error {
	for _, doc := range src.Docs {
		if err := p.ReplaceWithNode(dst, doc); err != nil {
			return err
		}
	}
	return nil
}

// ReplaceNode replace ast.File with ast.Node.
func (p *Path) ReplaceWithNode(dst *ast.File, node ast.Node) error {
	for _, doc := range dst.Docs {
		if node.Type() == ast.DocumentType {
			node = node.(*ast.DocumentNode).Body
		}
		if err := p.node.replace(doc.Body, node); err != nil {
			return err
		}
	}
	return nil
}

// AnnotateSource add annotation to passed source ( see section 5.1 in README.md ).
func (p *Path) AnnotateSource(source []byte, colored bool) ([]byte, error) {
	file, err := parser.ParseBytes([]byte(source), 0)
	if err != nil {
		return nil, err
	}
	node, err := p.FilterFile(file)
	if err != nil {
		return nil, err
	}
	var pp printer.Printer
	return []byte(pp.PrintErrorToken(node.GetToken(), colored)), nil
}

// PathBuilder represent builder for YAMLPath.
type PathBuilder struct {
	root *rootNode
	node pathNode
}

// Root add '$' to current path.
func (b *PathBuilder) Root() *PathBuilder {
	root := newRootNode()
	return &PathBuilder{root: root, node: root}
}

// IndexAll add '[*]' to current path.
func (b *PathBuilder) IndexAll() *PathBuilder {
	b.node = b.node.chain(newIndexAllNode())
	return b
}

// Recursive add '..selector' to current path.
func (b *PathBuilder) Recursive(selector string) *PathBuilder {
	b.node = b.node.chain(newRecursiveNode(selector))
	return b
}

func (b *PathBuilder) containsReservedPathCharacters(path string) bool {
	if strings.Contains(path, ".") {
		return true
	}
	if strings.Contains(path, "*") {
		return true
	}
	return false
}

func (b *PathBuilder) enclosedSingleQuote(name string) bool {
	return strings.HasPrefix(name, "'") && strings.HasSuffix(name, "'")
}

func (b *PathBuilder) normalizeSelectorName(name string) string {
	if b.enclosedSingleQuote(name) {
		// already escaped name
		return name
	}
	if b.containsReservedPathCharacters(name) {
		escapedName := strings.ReplaceAll(name, `'`, `\'`)
		return "'" + escapedName + "'"
	}
	return name
}

func (b *PathBuilder) child(name string) *PathBuilder {
	b.node = b.node.chain(newSelectorNode(name))
	return b
}

// Child add '.name' to current path.
func (b *PathBuilder) Child(name string) *PathBuilder {
	return b.child(b.normalizeSelectorName(name))
}

// Index add '[idx]' to current path.
func (b *PathBuilder) Index(idx uint) *PathBuilder {
	b.node = b.node.chain(newIndexNode(idx))
	return b
}

// Build build YAMLPath.
func (b *PathBuilder) Build() *Path {
	return &Path{node: b.root}
}

type pathNode interface {
	fmt.Stringer
	chain(pathNode) pathNode
	filter(ast.Node) (ast.Node, error)
	replace(ast.Node, ast.Node) error
}

type basePathNode struct {
	child pathNode
}

func (n *basePathNode) chain(node pathNode) pathNode {
	n.child = node
	return node
}

type rootNode struct {
	*basePathNode
}

func newRootNode() *rootNode {
	return &rootNode{basePathNode: &basePathNode{}}
}

func (n *rootNode) String() string {
	s := "$"
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

func (n *rootNode) filter(node ast.Node) (ast.Node, error) {
	if n.child == nil {
		return node, nil
	}
	filtered, err := n.child.filter(node)
	if err != nil {
		return nil, err
	}
	return filtered, nil
}

func (n *rootNode) replace(node ast.Node, target ast.Node) error {
	if n.child == nil {
		return nil
	}
	if err := n.child.replace(node, target); err != nil {
		return err
	}
	return nil
}

type selectorNode struct {
	*basePathNode
	selector string
}

func newSelectorNode(selector string) *selectorNode {
	return &selectorNode{
		basePathNode: &basePathNode{},
		selector:     selector,
	}
}

func (n *selectorNode) filter(node ast.Node) (ast.Node, error) {
	selector := n.selector
	if len(selector) > 1 && selector[0] == '\'' && selector[len(selector)-1] == '\'' {
		selector = selector[1 : len(selector)-1]
	}
	switch node.Type() {
	case ast.MappingType:
		for _, value := range node.(*ast.MappingNode).Values {
			key := value.Key.GetToken().Value
			if len(key) > 0 {
				switch key[0] {
				case '"':
					var err error
					key, err = strconv.Unquote(key)
					if err != nil {
						return nil, err
					}
				case '\'':
					if len(key) > 1 && key[len(key)-1] == '\'' {
						key = key[1 : len(key)-1]
					}
				}
			}
			if key == selector {
				if n.child == nil {
					return value.Value, nil
				}
				filtered, err := n.child.filter(value.Value)
				if err != nil {
					return nil, err
				}
				return filtered, nil
			}
		}
	case ast.MappingValueType:
		value, _ := node.(*ast.MappingValueNode)
		key := value.Key.GetToken().Value
		if key == selector {
			if n.child == nil {
				return value.Value, nil
			}
			filtered, err := n.child.filter(value.Value)
			if err != nil {
				return nil, err
			}
			return filtered, nil
		}
	default:
		return nil, fmt.Errorf("expected node type is map or map value. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	return nil, nil
}

func (n *selectorNode) replaceMapValue(value *ast.MappingValueNode, target ast.Node) error {
	key := value.Key.GetToken().Value
	if key != n.selector {
		return nil
	}
	if n.child == nil {
		if err := value.Replace(target); err != nil {
			return err
		}
	} else {
		if err := n.child.replace(value.Value, target); err != nil {
			return err
		}
	}
	return nil
}

func (n *selectorNode) replace(node ast.Node, target ast.Node) error {
	switch node.Type() {
	case ast.MappingType:
		for _, value := range node.(*ast.MappingNode).Values {
			if err := n.replaceMapValue(value, target); err != nil {
				return err
			}
		}
	case ast.MappingValueType:
		value, _ := node.(*ast.MappingValueNode)
		if err := n.replaceMapValue(value, target); err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected node type is map or map value. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	return nil
}

func (n *selectorNode) String() string {
	var builder PathBuilder
	selector := builder.normalizeSelectorName(n.selector)
	s := fmt.Sprintf(".%s", selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type indexNode struct {
	*basePathNode
	selector uint
}

func newIndexNode(selector uint) *indexNode {
	return &indexNode{
		basePathNode: &basePathNode{},
		selector:     selector,
	}
}

func (n *indexNode) filter(node ast.Node) (ast.Node, error) {
	if node.Type() != ast.SequenceType {
		return nil, fmt.Errorf("expected sequence type node. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	sequence, _ := node.(*ast.SequenceNode)
	if n.selector >= uint(len(sequence.Values)) {
		return nil, fmt.Errorf("expected index is %d. but got sequences has %d items: %w", n.selector, len(sequence.Values), ErrInvalidQuery)
	}
	value := sequence.Values[n.selector]
	if n.child == nil {
		return value, nil
	}
	filtered, err := n.child.filter(value)
	if err != nil {
		return nil, err
	}
	return filtered, nil
}

func (n *indexNode) replace(node ast.Node, target ast.Node) error {
	if node.Type() != ast.SequenceType {
		return fmt.Errorf("expected sequence type node. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	sequence, _ := node.(*ast.SequenceNode)
	if n.selector >= uint(len(sequence.Values)) {
		return fmt.Errorf("expected index is %d. but got sequences has %d items: %w", n.selector, len(sequence.Values), ErrInvalidQuery)
	}
	if n.child == nil {
		if err := sequence.Replace(int(n.selector), target); err != nil {
			return err
		}
		return nil
	}
	if err := n.child.replace(sequence.Values[n.selector], target); err != nil {
		return err
	}
	return nil
}

func (n *indexNode) String() string {
	s := fmt.Sprintf("[%d]", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type indexAllNode struct {
	*basePathNode
}

func newIndexAllNode() *indexAllNode {
	return &indexAllNode{
		basePathNode: &basePathNode{},
	}
}

func (n *indexAllNode) String() string {
	s := "[*]"
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

func (n *indexAllNode) filter(node ast.Node) (ast.Node, error) {
	if node.Type() != ast.SequenceType {
		return nil, fmt.Errorf("expected sequence type node. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	sequence, _ := node.(*ast.SequenceNode)
	if n.child == nil {
		return sequence, nil
	}
	out := *sequence
	out.Values = []ast.Node{}
	for _, value := range sequence.Values {
		filtered, err := n.child.filter(value)
		if err != nil {
			return nil, err
		}
		out.Values = append(out.Values, filtered)
	}
	return &out, nil
}

func (n *indexAllNode) replace(node ast.Node, target ast.Node) error {
	if node.Type() != ast.SequenceType {
		return fmt.Errorf("expected sequence type node. but got %s: %w", node.Type(), ErrInvalidQuery)
	}
	sequence, _ := node.(*ast.SequenceNode)
	if n.child == nil {
		for idx := range sequence.Values {
			if err := sequence.Replace(idx, target); err != nil {
				return err
			}
		}
		return nil
	}
	for _, value := range sequence.Values {
		if err := n.child.replace(value, target); err != nil {
			return err
		}
	}
	return nil
}

type recursiveNode struct {
	*basePathNode
	selector string
}

func newRecursiveNode(selector string) *recursiveNode {
	return &recursiveNode{
		basePathNode: &basePathNode{},
		selector:     selector,
	}
}

func (n *recursiveNode) String() string {
	s := fmt.Sprintf("..%s", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

func (n *recursiveNode) filterNode(node ast.Node) (*ast.SequenceNode, error) {
	sequence := &ast.SequenceNode{BaseNode: &ast.BaseNode{}}
	switch typedNode := node.(type) {
	case *ast.MappingNode:
		for _, value := range typedNode.Values {
			seq, err := n.filterNode(value)
			if err != nil {
				return nil, err
			}
			sequence.Values = append(sequence.Values, seq.Values...)
		}
	case *ast.MappingValueNode:
		key := typedNode.Key.GetToken().Value
		if n.selector == key {
			sequence.Values = append(sequence.Values, typedNode.Value)
		}
		seq, err := n.filterNode(typedNode.Value)
		if err != nil {
			return nil, err
		}
		sequence.Values = append(sequence.Values, seq.Values...)
	case *ast.SequenceNode:
		for _, value := range typedNode.Values {
			seq, err := n.filterNode(value)
			if err != nil {
				return nil, err
			}
			sequence.Values = append(sequence.Values, seq.Values...)
		}
	}
	return sequence, nil
}

func (n *recursiveNode) filter(node ast.Node) (ast.Node, error) {
	sequence, err := n.filterNode(node)
	if err != nil {
		return nil, err
	}
	sequence.Start = node.GetToken()
	return sequence, nil
}

func (n *recursiveNode) replaceNode(node ast.Node, target ast.Node) error {
	switch typedNode := node.(type) {
	case *ast.MappingNode:
		for _, value := range typedNode.Values {
			if err := n.replaceNode(value, target); err != nil {
				return err
			}
		}
	case *ast.MappingValueNode:
		key := typedNode.Key.GetToken().Value
		if n.selector == key {
			if err := typedNode.Replace(target); err != nil {
				return err
			}
		}
		if err := n.replaceNode(typedNode.Value, target); err != nil {
			return err
		}
	case *ast.SequenceNode:
		for _, value := range typedNode.Values {
			if err := n.replaceNode(value, target); err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *recursiveNode) replace(node ast.Node, target ast.Node) error {
	if err := n.replaceNode(node, target); err != nil {
		return err
	}
	return nil
}
