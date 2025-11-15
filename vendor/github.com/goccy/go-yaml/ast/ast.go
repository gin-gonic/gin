package ast

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/token"
)

var (
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrInvalidAnchorName = errors.New("invalid anchor name")
	ErrInvalidAliasName  = errors.New("invalid alias name")
)

// NodeType type identifier of node
type NodeType int

const (
	// UnknownNodeType type identifier for default
	UnknownNodeType NodeType = iota
	// DocumentType type identifier for document node
	DocumentType
	// NullType type identifier for null node
	NullType
	// BoolType type identifier for boolean node
	BoolType
	// IntegerType type identifier for integer node
	IntegerType
	// FloatType type identifier for float node
	FloatType
	// InfinityType type identifier for infinity node
	InfinityType
	// NanType type identifier for nan node
	NanType
	// StringType type identifier for string node
	StringType
	// MergeKeyType type identifier for merge key node
	MergeKeyType
	// LiteralType type identifier for literal node
	LiteralType
	// MappingType type identifier for mapping node
	MappingType
	// MappingKeyType type identifier for mapping key node
	MappingKeyType
	// MappingValueType type identifier for mapping value node
	MappingValueType
	// SequenceType type identifier for sequence node
	SequenceType
	// SequenceEntryType type identifier for sequence entry node
	SequenceEntryType
	// AnchorType type identifier for anchor node
	AnchorType
	// AliasType type identifier for alias node
	AliasType
	// DirectiveType type identifier for directive node
	DirectiveType
	// TagType type identifier for tag node
	TagType
	// CommentType type identifier for comment node
	CommentType
	// CommentGroupType type identifier for comment group node
	CommentGroupType
)

// String node type identifier to text
func (t NodeType) String() string {
	switch t {
	case UnknownNodeType:
		return "UnknownNode"
	case DocumentType:
		return "Document"
	case NullType:
		return "Null"
	case BoolType:
		return "Bool"
	case IntegerType:
		return "Integer"
	case FloatType:
		return "Float"
	case InfinityType:
		return "Infinity"
	case NanType:
		return "Nan"
	case StringType:
		return "String"
	case MergeKeyType:
		return "MergeKey"
	case LiteralType:
		return "Literal"
	case MappingType:
		return "Mapping"
	case MappingKeyType:
		return "MappingKey"
	case MappingValueType:
		return "MappingValue"
	case SequenceType:
		return "Sequence"
	case SequenceEntryType:
		return "SequenceEntry"
	case AnchorType:
		return "Anchor"
	case AliasType:
		return "Alias"
	case DirectiveType:
		return "Directive"
	case TagType:
		return "Tag"
	case CommentType:
		return "Comment"
	case CommentGroupType:
		return "CommentGroup"
	}
	return ""
}

// String node type identifier to YAML Structure name
// based on https://yaml.org/spec/1.2/spec.html
func (t NodeType) YAMLName() string {
	switch t {
	case UnknownNodeType:
		return "unknown"
	case DocumentType:
		return "document"
	case NullType:
		return "null"
	case BoolType:
		return "boolean"
	case IntegerType:
		return "int"
	case FloatType:
		return "float"
	case InfinityType:
		return "inf"
	case NanType:
		return "nan"
	case StringType:
		return "string"
	case MergeKeyType:
		return "merge key"
	case LiteralType:
		return "scalar"
	case MappingType:
		return "mapping"
	case MappingKeyType:
		return "key"
	case MappingValueType:
		return "value"
	case SequenceType:
		return "sequence"
	case SequenceEntryType:
		return "value"
	case AnchorType:
		return "anchor"
	case AliasType:
		return "alias"
	case DirectiveType:
		return "directive"
	case TagType:
		return "tag"
	case CommentType:
		return "comment"
	case CommentGroupType:
		return "comment"
	}
	return ""
}

// Node type of node
type Node interface {
	io.Reader
	// String node to text
	String() string
	// GetToken returns token instance
	GetToken() *token.Token
	// Type returns type of node
	Type() NodeType
	// AddColumn add column number to child nodes recursively
	AddColumn(int)
	// SetComment set comment token to node
	SetComment(*CommentGroupNode) error
	// Comment returns comment token instance
	GetComment() *CommentGroupNode
	// GetPath returns YAMLPath for the current node
	GetPath() string
	// SetPath set YAMLPath for the current node
	SetPath(string)
	// MarshalYAML
	MarshalYAML() ([]byte, error)
	// already read length
	readLen() int
	// append read length
	addReadLen(int)
	// clean read length
	clearLen()
}

// MapKeyNode type for map key node
type MapKeyNode interface {
	Node
	IsMergeKey() bool
	// String node to text without comment
	stringWithoutComment() string
}

// ScalarNode type for scalar node
type ScalarNode interface {
	MapKeyNode
	GetValue() interface{}
}

type BaseNode struct {
	Path    string
	Comment *CommentGroupNode
	read    int
}

func addCommentString(base string, node *CommentGroupNode) string {
	return fmt.Sprintf("%s %s", base, node.String())
}

func (n *BaseNode) readLen() int {
	return n.read
}

func (n *BaseNode) clearLen() {
	n.read = 0
}

func (n *BaseNode) addReadLen(len int) {
	n.read += len
}

// GetPath returns YAMLPath for the current node.
func (n *BaseNode) GetPath() string {
	if n == nil {
		return ""
	}
	return n.Path
}

// SetPath set YAMLPath for the current node.
func (n *BaseNode) SetPath(path string) {
	if n == nil {
		return
	}
	n.Path = path
}

// GetComment returns comment token instance
func (n *BaseNode) GetComment() *CommentGroupNode {
	return n.Comment
}

// SetComment set comment token
func (n *BaseNode) SetComment(node *CommentGroupNode) error {
	n.Comment = node
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readNode(p []byte, node Node) (int, error) {
	s := node.String()
	readLen := node.readLen()
	remain := len(s) - readLen
	if remain == 0 {
		node.clearLen()
		return 0, io.EOF
	}
	size := min(remain, len(p))
	for idx, b := range []byte(s[readLen : readLen+size]) {
		p[idx] = byte(b)
	}
	node.addReadLen(size)
	return size, nil
}

func checkLineBreak(t *token.Token) bool {
	if t.Prev != nil {
		lbc := "\n"
		prev := t.Prev
		var adjustment int
		// if the previous type is sequence entry use the previous type for that
		if prev.Type == token.SequenceEntryType {
			// as well as switching to previous type count any new lines in origin to account for:
			// -
			//   b: c
			adjustment = strings.Count(strings.TrimRight(t.Origin, lbc), lbc)
			if prev.Prev != nil {
				prev = prev.Prev
			}
		}
		lineDiff := t.Position.Line - prev.Position.Line - 1
		if lineDiff > 0 {
			if prev.Type == token.StringType {
				// Remove any line breaks included in multiline string
				adjustment += strings.Count(strings.TrimRight(strings.TrimSpace(prev.Origin), lbc), lbc)
			}
			// Due to the way that comment parsing works its assumed that when a null value does not have new line in origin
			// it was squashed therefore difference is ignored.
			// foo:
			//  bar:
			//  # comment
			//  baz: 1
			// becomes
			// foo:
			//  bar: null # comment
			//
			//  baz: 1
			if prev.Type == token.NullType || prev.Type == token.ImplicitNullType {
				return strings.Count(prev.Origin, lbc) > 0
			}
			if lineDiff-adjustment > 0 {
				return true
			}
		}
	}
	return false
}

// Null create node for null value
func Null(tk *token.Token) *NullNode {
	return &NullNode{
		BaseNode: &BaseNode{},
		Token:    tk,
	}
}

// Bool create node for boolean value
func Bool(tk *token.Token) *BoolNode {
	b, _ := strconv.ParseBool(tk.Value)
	return &BoolNode{
		BaseNode: &BaseNode{},
		Token:    tk,
		Value:    b,
	}
}

// Integer create node for integer value
func Integer(tk *token.Token) *IntegerNode {
	var v any
	if num := token.ToNumber(tk.Value); num != nil {
		v = num.Value
	}
	return &IntegerNode{
		BaseNode: &BaseNode{},
		Token:    tk,
		Value:    v,
	}
}

// Float create node for float value
func Float(tk *token.Token) *FloatNode {
	var v float64
	if num := token.ToNumber(tk.Value); num != nil && num.Type == token.NumberTypeFloat {
		value, ok := num.Value.(float64)
		if ok {
			v = value
		}
	}
	return &FloatNode{
		BaseNode: &BaseNode{},
		Token:    tk,
		Value:    v,
	}
}

// Infinity create node for .inf or -.inf value
func Infinity(tk *token.Token) *InfinityNode {
	node := &InfinityNode{
		BaseNode: &BaseNode{},
		Token:    tk,
	}
	switch tk.Value {
	case ".inf", ".Inf", ".INF":
		node.Value = math.Inf(0)
	case "-.inf", "-.Inf", "-.INF":
		node.Value = math.Inf(-1)
	}
	return node
}

// Nan create node for .nan value
func Nan(tk *token.Token) *NanNode {
	return &NanNode{
		BaseNode: &BaseNode{},
		Token:    tk,
	}
}

// String create node for string value
func String(tk *token.Token) *StringNode {
	return &StringNode{
		BaseNode: &BaseNode{},
		Token:    tk,
		Value:    tk.Value,
	}
}

// Comment create node for comment
func Comment(tk *token.Token) *CommentNode {
	return &CommentNode{
		BaseNode: &BaseNode{},
		Token:    tk,
	}
}

func CommentGroup(comments []*token.Token) *CommentGroupNode {
	nodes := []*CommentNode{}
	for _, comment := range comments {
		nodes = append(nodes, Comment(comment))
	}
	return &CommentGroupNode{
		BaseNode: &BaseNode{},
		Comments: nodes,
	}
}

// MergeKey create node for merge key ( << )
func MergeKey(tk *token.Token) *MergeKeyNode {
	return &MergeKeyNode{
		BaseNode: &BaseNode{},
		Token:    tk,
	}
}

// Mapping create node for map
func Mapping(tk *token.Token, isFlowStyle bool, values ...*MappingValueNode) *MappingNode {
	node := &MappingNode{
		BaseNode:    &BaseNode{},
		Start:       tk,
		IsFlowStyle: isFlowStyle,
		Values:      []*MappingValueNode{},
	}
	node.Values = append(node.Values, values...)
	return node
}

// MappingValue create node for mapping value
func MappingValue(tk *token.Token, key MapKeyNode, value Node) *MappingValueNode {
	return &MappingValueNode{
		BaseNode: &BaseNode{},
		Start:    tk,
		Key:      key,
		Value:    value,
	}
}

// MappingKey create node for map key ( '?' ).
func MappingKey(tk *token.Token) *MappingKeyNode {
	return &MappingKeyNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

// Sequence create node for sequence
func Sequence(tk *token.Token, isFlowStyle bool) *SequenceNode {
	return &SequenceNode{
		BaseNode:    &BaseNode{},
		Start:       tk,
		IsFlowStyle: isFlowStyle,
		Values:      []Node{},
	}
}

func Anchor(tk *token.Token) *AnchorNode {
	return &AnchorNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

func Alias(tk *token.Token) *AliasNode {
	return &AliasNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

func Document(tk *token.Token, body Node) *DocumentNode {
	return &DocumentNode{
		BaseNode: &BaseNode{},
		Start:    tk,
		Body:     body,
	}
}

func Directive(tk *token.Token) *DirectiveNode {
	return &DirectiveNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

func Literal(tk *token.Token) *LiteralNode {
	return &LiteralNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

func Tag(tk *token.Token) *TagNode {
	return &TagNode{
		BaseNode: &BaseNode{},
		Start:    tk,
	}
}

// File contains all documents in YAML file
type File struct {
	Name string
	Docs []*DocumentNode
}

// Read implements (io.Reader).Read
func (f *File) Read(p []byte) (int, error) {
	for _, doc := range f.Docs {
		n, err := doc.Read(p)
		if err == io.EOF {
			continue
		}
		return n, nil
	}
	return 0, io.EOF
}

// String all documents to text
func (f *File) String() string {
	docs := []string{}
	for _, doc := range f.Docs {
		docs = append(docs, doc.String())
	}
	if len(docs) > 0 {
		return strings.Join(docs, "\n") + "\n"
	} else {
		return ""
	}
}

// DocumentNode type of Document
type DocumentNode struct {
	*BaseNode
	Start *token.Token // position of DocumentHeader ( `---` )
	End   *token.Token // position of DocumentEnd ( `...` )
	Body  Node
}

// Read implements (io.Reader).Read
func (d *DocumentNode) Read(p []byte) (int, error) {
	return readNode(p, d)
}

// Type returns DocumentNodeType
func (d *DocumentNode) Type() NodeType { return DocumentType }

// GetToken returns token instance
func (d *DocumentNode) GetToken() *token.Token {
	return d.Body.GetToken()
}

// AddColumn add column number to child nodes recursively
func (d *DocumentNode) AddColumn(col int) {
	if d.Body != nil {
		d.Body.AddColumn(col)
	}
}

// String document to text
func (d *DocumentNode) String() string {
	doc := []string{}
	if d.Start != nil {
		doc = append(doc, d.Start.Value)
	}
	if d.Body != nil {
		doc = append(doc, d.Body.String())
	}
	if d.End != nil {
		doc = append(doc, d.End.Value)
	}
	return strings.Join(doc, "\n")
}

// MarshalYAML encodes to a YAML text
func (d *DocumentNode) MarshalYAML() ([]byte, error) {
	return []byte(d.String()), nil
}

// NullNode type of null node
type NullNode struct {
	*BaseNode
	Token *token.Token
}

// Read implements (io.Reader).Read
func (n *NullNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns NullType
func (n *NullNode) Type() NodeType { return NullType }

// GetToken returns token instance
func (n *NullNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *NullNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns nil value
func (n *NullNode) GetValue() interface{} {
	return nil
}

// String returns `null` text
func (n *NullNode) String() string {
	if n.Token.Type == token.ImplicitNullType {
		if n.Comment != nil {
			return n.Comment.String()
		}
		return ""
	}
	if n.Comment != nil {
		return addCommentString("null", n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *NullNode) stringWithoutComment() string {
	return "null"
}

// MarshalYAML encodes to a YAML text
func (n *NullNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *NullNode) IsMergeKey() bool {
	return false
}

// IntegerNode type of integer node
type IntegerNode struct {
	*BaseNode
	Token *token.Token
	Value interface{} // int64 or uint64 value
}

// Read implements (io.Reader).Read
func (n *IntegerNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns IntegerType
func (n *IntegerNode) Type() NodeType { return IntegerType }

// GetToken returns token instance
func (n *IntegerNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *IntegerNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns int64 value
func (n *IntegerNode) GetValue() interface{} {
	return n.Value
}

// String int64 to text
func (n *IntegerNode) String() string {
	if n.Comment != nil {
		return addCommentString(n.Token.Value, n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *IntegerNode) stringWithoutComment() string {
	return n.Token.Value
}

// MarshalYAML encodes to a YAML text
func (n *IntegerNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *IntegerNode) IsMergeKey() bool {
	return false
}

// FloatNode type of float node
type FloatNode struct {
	*BaseNode
	Token     *token.Token
	Precision int
	Value     float64
}

// Read implements (io.Reader).Read
func (n *FloatNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns FloatType
func (n *FloatNode) Type() NodeType { return FloatType }

// GetToken returns token instance
func (n *FloatNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *FloatNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns float64 value
func (n *FloatNode) GetValue() interface{} {
	return n.Value
}

// String float64 to text
func (n *FloatNode) String() string {
	if n.Comment != nil {
		return addCommentString(n.Token.Value, n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *FloatNode) stringWithoutComment() string {
	return n.Token.Value
}

// MarshalYAML encodes to a YAML text
func (n *FloatNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *FloatNode) IsMergeKey() bool {
	return false
}

// StringNode type of string node
type StringNode struct {
	*BaseNode
	Token *token.Token
	Value string
}

// Read implements (io.Reader).Read
func (n *StringNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns StringType
func (n *StringNode) Type() NodeType { return StringType }

// GetToken returns token instance
func (n *StringNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *StringNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns string value
func (n *StringNode) GetValue() interface{} {
	return n.Value
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *StringNode) IsMergeKey() bool {
	return false
}

// escapeSingleQuote escapes s to a single quoted scalar.
// https://yaml.org/spec/1.2.2/#732-single-quoted-style
func escapeSingleQuote(s string) string {
	var sb strings.Builder
	growLen := len(s) + // s includes also one ' from the doubled pair
		2 + // opening and closing '
		strings.Count(s, "'") // ' added by ReplaceAll
	sb.Grow(growLen)
	sb.WriteString("'")
	sb.WriteString(strings.ReplaceAll(s, "'", "''"))
	sb.WriteString("'")
	return sb.String()
}

// String string value to text with quote or literal header if required
func (n *StringNode) String() string {
	switch n.Token.Type {
	case token.SingleQuoteType:
		quoted := escapeSingleQuote(n.Value)
		if n.Comment != nil {
			return addCommentString(quoted, n.Comment)
		}
		return quoted
	case token.DoubleQuoteType:
		quoted := strconv.Quote(n.Value)
		if n.Comment != nil {
			return addCommentString(quoted, n.Comment)
		}
		return quoted
	}

	lbc := token.DetectLineBreakCharacter(n.Value)
	if strings.Contains(n.Value, lbc) {
		// This block assumes that the line breaks in this inside scalar content and the Outside scalar content are the same.
		// It works mostly, but inconsistencies occur if line break characters are mixed.
		header := token.LiteralBlockHeader(n.Value)
		space := strings.Repeat(" ", n.Token.Position.Column-1)
		indent := strings.Repeat(" ", n.Token.Position.IndentNum)
		values := []string{}
		for _, v := range strings.Split(n.Value, lbc) {
			values = append(values, fmt.Sprintf("%s%s%s", space, indent, v))
		}
		block := strings.TrimSuffix(strings.TrimSuffix(strings.Join(values, lbc), fmt.Sprintf("%s%s%s", lbc, indent, space)), fmt.Sprintf("%s%s", indent, space))
		return fmt.Sprintf("%s%s%s", header, lbc, block)
	} else if len(n.Value) > 0 && (n.Value[0] == '{' || n.Value[0] == '[') {
		return fmt.Sprintf(`'%s'`, n.Value)
	}
	if n.Comment != nil {
		return addCommentString(n.Value, n.Comment)
	}
	return n.Value
}

func (n *StringNode) stringWithoutComment() string {
	switch n.Token.Type {
	case token.SingleQuoteType:
		quoted := fmt.Sprintf(`'%s'`, n.Value)
		return quoted
	case token.DoubleQuoteType:
		quoted := strconv.Quote(n.Value)
		return quoted
	}

	lbc := token.DetectLineBreakCharacter(n.Value)
	if strings.Contains(n.Value, lbc) {
		// This block assumes that the line breaks in this inside scalar content and the Outside scalar content are the same.
		// It works mostly, but inconsistencies occur if line break characters are mixed.
		header := token.LiteralBlockHeader(n.Value)
		space := strings.Repeat(" ", n.Token.Position.Column-1)
		indent := strings.Repeat(" ", n.Token.Position.IndentNum)
		values := []string{}
		for _, v := range strings.Split(n.Value, lbc) {
			values = append(values, fmt.Sprintf("%s%s%s", space, indent, v))
		}
		block := strings.TrimSuffix(strings.TrimSuffix(strings.Join(values, lbc), fmt.Sprintf("%s%s%s", lbc, indent, space)), fmt.Sprintf("  %s", space))
		return fmt.Sprintf("%s%s%s", header, lbc, block)
	} else if len(n.Value) > 0 && (n.Value[0] == '{' || n.Value[0] == '[') {
		return fmt.Sprintf(`'%s'`, n.Value)
	}
	return n.Value
}

// MarshalYAML encodes to a YAML text
func (n *StringNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// LiteralNode type of literal node
type LiteralNode struct {
	*BaseNode
	Start *token.Token
	Value *StringNode
}

// Read implements (io.Reader).Read
func (n *LiteralNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns LiteralType
func (n *LiteralNode) Type() NodeType { return LiteralType }

// GetToken returns token instance
func (n *LiteralNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *LiteralNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// GetValue returns string value
func (n *LiteralNode) GetValue() interface{} {
	return n.String()
}

// String literal to text
func (n *LiteralNode) String() string {
	origin := n.Value.GetToken().Origin
	lit := strings.TrimRight(strings.TrimRight(origin, " "), "\n")
	if n.Comment != nil {
		return fmt.Sprintf("%s %s\n%s", n.Start.Value, n.Comment.String(), lit)
	}
	return fmt.Sprintf("%s\n%s", n.Start.Value, lit)
}

func (n *LiteralNode) stringWithoutComment() string {
	return n.String()
}

// MarshalYAML encodes to a YAML text
func (n *LiteralNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *LiteralNode) IsMergeKey() bool {
	return false
}

// MergeKeyNode type of merge key node
type MergeKeyNode struct {
	*BaseNode
	Token *token.Token
}

// Read implements (io.Reader).Read
func (n *MergeKeyNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns MergeKeyType
func (n *MergeKeyNode) Type() NodeType { return MergeKeyType }

// GetToken returns token instance
func (n *MergeKeyNode) GetToken() *token.Token {
	return n.Token
}

// GetValue returns '<<' value
func (n *MergeKeyNode) GetValue() interface{} {
	return n.Token.Value
}

// String returns '<<' value
func (n *MergeKeyNode) String() string {
	return n.stringWithoutComment()
}

func (n *MergeKeyNode) stringWithoutComment() string {
	return n.Token.Value
}

// AddColumn add column number to child nodes recursively
func (n *MergeKeyNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// MarshalYAML encodes to a YAML text
func (n *MergeKeyNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *MergeKeyNode) IsMergeKey() bool {
	return true
}

// BoolNode type of boolean node
type BoolNode struct {
	*BaseNode
	Token *token.Token
	Value bool
}

// Read implements (io.Reader).Read
func (n *BoolNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns BoolType
func (n *BoolNode) Type() NodeType { return BoolType }

// GetToken returns token instance
func (n *BoolNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *BoolNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns boolean value
func (n *BoolNode) GetValue() interface{} {
	return n.Value
}

// String boolean to text
func (n *BoolNode) String() string {
	if n.Comment != nil {
		return addCommentString(n.Token.Value, n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *BoolNode) stringWithoutComment() string {
	return n.Token.Value
}

// MarshalYAML encodes to a YAML text
func (n *BoolNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *BoolNode) IsMergeKey() bool {
	return false
}

// InfinityNode type of infinity node
type InfinityNode struct {
	*BaseNode
	Token *token.Token
	Value float64
}

// Read implements (io.Reader).Read
func (n *InfinityNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns InfinityType
func (n *InfinityNode) Type() NodeType { return InfinityType }

// GetToken returns token instance
func (n *InfinityNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *InfinityNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns math.Inf(0) or math.Inf(-1)
func (n *InfinityNode) GetValue() interface{} {
	return n.Value
}

// String infinity to text
func (n *InfinityNode) String() string {
	if n.Comment != nil {
		return addCommentString(n.Token.Value, n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *InfinityNode) stringWithoutComment() string {
	return n.Token.Value
}

// MarshalYAML encodes to a YAML text
func (n *InfinityNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *InfinityNode) IsMergeKey() bool {
	return false
}

// NanNode type of nan node
type NanNode struct {
	*BaseNode
	Token *token.Token
}

// Read implements (io.Reader).Read
func (n *NanNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns NanType
func (n *NanNode) Type() NodeType { return NanType }

// GetToken returns token instance
func (n *NanNode) GetToken() *token.Token {
	return n.Token
}

// AddColumn add column number to child nodes recursively
func (n *NanNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// GetValue returns math.NaN()
func (n *NanNode) GetValue() interface{} {
	return math.NaN()
}

// String returns .nan
func (n *NanNode) String() string {
	if n.Comment != nil {
		return addCommentString(n.Token.Value, n.Comment)
	}
	return n.stringWithoutComment()
}

func (n *NanNode) stringWithoutComment() string {
	return n.Token.Value
}

// MarshalYAML encodes to a YAML text
func (n *NanNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *NanNode) IsMergeKey() bool {
	return false
}

// MapNode interface of MappingValueNode / MappingNode
type MapNode interface {
	MapRange() *MapNodeIter
}

// MapNodeIter is an iterator for ranging over a MapNode
type MapNodeIter struct {
	values []*MappingValueNode
	idx    int
}

const (
	startRangeIndex = -1
)

// Next advances the map iterator and reports whether there is another entry.
// It returns false when the iterator is exhausted.
func (m *MapNodeIter) Next() bool {
	m.idx++
	next := m.idx < len(m.values)
	return next
}

// Key returns the key of the iterator's current map node entry.
func (m *MapNodeIter) Key() MapKeyNode {
	return m.values[m.idx].Key
}

// Value returns the value of the iterator's current map node entry.
func (m *MapNodeIter) Value() Node {
	return m.values[m.idx].Value
}

// KeyValue returns the MappingValueNode of the iterator's current map node entry.
func (m *MapNodeIter) KeyValue() *MappingValueNode {
	return m.values[m.idx]
}

// MappingNode type of mapping node
type MappingNode struct {
	*BaseNode
	Start       *token.Token
	End         *token.Token
	IsFlowStyle bool
	Values      []*MappingValueNode
	FootComment *CommentGroupNode
}

func (n *MappingNode) startPos() *token.Position {
	if len(n.Values) == 0 {
		return n.Start.Position
	}
	return n.Values[0].Key.GetToken().Position
}

// Merge merge key/value of map.
func (n *MappingNode) Merge(target *MappingNode) {
	keyToMapValueMap := map[string]*MappingValueNode{}
	for _, value := range n.Values {
		key := value.Key.String()
		keyToMapValueMap[key] = value
	}
	column := n.startPos().Column - target.startPos().Column
	target.AddColumn(column)
	for _, value := range target.Values {
		mapValue, exists := keyToMapValueMap[value.Key.String()]
		if exists {
			mapValue.Value = value.Value
		} else {
			n.Values = append(n.Values, value)
		}
	}
}

// SetIsFlowStyle set value to IsFlowStyle field recursively.
func (n *MappingNode) SetIsFlowStyle(isFlow bool) {
	n.IsFlowStyle = isFlow
	for _, value := range n.Values {
		value.SetIsFlowStyle(isFlow)
	}
}

// Read implements (io.Reader).Read
func (n *MappingNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns MappingType
func (n *MappingNode) Type() NodeType { return MappingType }

// GetToken returns token instance
func (n *MappingNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *MappingNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	n.End.AddColumn(col)
	for _, value := range n.Values {
		value.AddColumn(col)
	}
}

func (n *MappingNode) flowStyleString(commentMode bool) string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, strings.TrimLeft(value.String(), " "))
	}
	mapText := fmt.Sprintf("{%s}", strings.Join(values, ", "))
	if commentMode && n.Comment != nil {
		return addCommentString(mapText, n.Comment)
	}
	return mapText
}

func (n *MappingNode) blockStyleString(commentMode bool) string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	mapText := strings.Join(values, "\n")
	if commentMode && n.Comment != nil {
		value := values[0]
		var spaceNum int
		for i := 0; i < len(value); i++ {
			if value[i] != ' ' {
				break
			}
			spaceNum++
		}
		comment := n.Comment.StringWithSpace(spaceNum)
		return fmt.Sprintf("%s\n%s", comment, mapText)
	}
	return mapText
}

// String mapping values to text
func (n *MappingNode) String() string {
	if len(n.Values) == 0 {
		if n.Comment != nil {
			return addCommentString("{}", n.Comment)
		}
		return "{}"
	}

	commentMode := true
	if n.IsFlowStyle || len(n.Values) == 0 {
		return n.flowStyleString(commentMode)
	}
	return n.blockStyleString(commentMode)
}

// MapRange implements MapNode protocol
func (n *MappingNode) MapRange() *MapNodeIter {
	return &MapNodeIter{
		idx:    startRangeIndex,
		values: n.Values,
	}
}

// MarshalYAML encodes to a YAML text
func (n *MappingNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// MappingKeyNode type of tag node
type MappingKeyNode struct {
	*BaseNode
	Start *token.Token
	Value Node
}

// Read implements (io.Reader).Read
func (n *MappingKeyNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns MappingKeyType
func (n *MappingKeyNode) Type() NodeType { return MappingKeyType }

// GetToken returns token instance
func (n *MappingKeyNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *MappingKeyNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// String tag to text
func (n *MappingKeyNode) String() string {
	return n.stringWithoutComment()
}

func (n *MappingKeyNode) stringWithoutComment() string {
	return fmt.Sprintf("%s %s", n.Start.Value, n.Value.String())
}

// MarshalYAML encodes to a YAML text
func (n *MappingKeyNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *MappingKeyNode) IsMergeKey() bool {
	if n.Value == nil {
		return false
	}
	key, ok := n.Value.(MapKeyNode)
	if !ok {
		return false
	}
	return key.IsMergeKey()
}

// MappingValueNode type of mapping value
type MappingValueNode struct {
	*BaseNode
	Start        *token.Token // delimiter token ':'.
	CollectEntry *token.Token // collect entry token ','.
	Key          MapKeyNode
	Value        Node
	FootComment  *CommentGroupNode
	IsFlowStyle  bool
}

// Replace replace value node.
func (n *MappingValueNode) Replace(value Node) error {
	column := n.Value.GetToken().Position.Column - value.GetToken().Position.Column
	value.AddColumn(column)
	n.Value = value
	return nil
}

// Read implements (io.Reader).Read
func (n *MappingValueNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns MappingValueType
func (n *MappingValueNode) Type() NodeType { return MappingValueType }

// GetToken returns token instance
func (n *MappingValueNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *MappingValueNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Key != nil {
		n.Key.AddColumn(col)
	}
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetIsFlowStyle set value to IsFlowStyle field recursively.
func (n *MappingValueNode) SetIsFlowStyle(isFlow bool) {
	n.IsFlowStyle = isFlow
	switch value := n.Value.(type) {
	case *MappingNode:
		value.SetIsFlowStyle(isFlow)
	case *MappingValueNode:
		value.SetIsFlowStyle(isFlow)
	case *SequenceNode:
		value.SetIsFlowStyle(isFlow)
	}
}

// String mapping value to text
func (n *MappingValueNode) String() string {
	var text string
	if n.Comment != nil {
		text = fmt.Sprintf(
			"%s\n%s",
			n.Comment.StringWithSpace(n.Key.GetToken().Position.Column-1),
			n.toString(),
		)
	} else {
		text = n.toString()
	}
	if n.FootComment != nil {
		text += fmt.Sprintf("\n%s", n.FootComment.StringWithSpace(n.Key.GetToken().Position.Column-1))
	}
	return text
}

func (n *MappingValueNode) toString() string {
	space := strings.Repeat(" ", n.Key.GetToken().Position.Column-1)
	if checkLineBreak(n.Key.GetToken()) {
		space = fmt.Sprintf("%s%s", "\n", space)
	}
	keyIndentLevel := n.Key.GetToken().Position.IndentLevel
	valueIndentLevel := n.Value.GetToken().Position.IndentLevel
	keyComment := n.Key.GetComment()
	if _, ok := n.Value.(ScalarNode); ok {
		value := n.Value.String()
		if value == "" {
			// implicit null value.
			return fmt.Sprintf("%s%s:", space, n.Key.String())
		}
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), value)
	} else if keyIndentLevel < valueIndentLevel && !n.IsFlowStyle {
		if keyComment != nil {
			return fmt.Sprintf(
				"%s%s: %s\n%s",
				space,
				n.Key.stringWithoutComment(),
				keyComment.String(),
				n.Value.String(),
			)
		}
		return fmt.Sprintf("%s%s:\n%s", space, n.Key.String(), n.Value.String())
	} else if m, ok := n.Value.(*MappingNode); ok && (m.IsFlowStyle || len(m.Values) == 0) {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if s, ok := n.Value.(*SequenceNode); ok && (s.IsFlowStyle || len(s.Values) == 0) {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AnchorNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AliasNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*TagNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	}

	if keyComment != nil {
		return fmt.Sprintf(
			"%s%s: %s\n%s",
			space,
			n.Key.stringWithoutComment(),
			keyComment.String(),
			n.Value.String(),
		)
	}
	if m, ok := n.Value.(*MappingNode); ok && m.Comment != nil {
		return fmt.Sprintf(
			"%s%s: %s",
			space,
			n.Key.String(),
			strings.TrimLeft(n.Value.String(), " "),
		)
	}
	return fmt.Sprintf("%s%s:\n%s", space, n.Key.String(), n.Value.String())
}

// MapRange implements MapNode protocol
func (n *MappingValueNode) MapRange() *MapNodeIter {
	return &MapNodeIter{
		idx:    startRangeIndex,
		values: []*MappingValueNode{n},
	}
}

// MarshalYAML encodes to a YAML text
func (n *MappingValueNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// ArrayNode interface of SequenceNode
type ArrayNode interface {
	ArrayRange() *ArrayNodeIter
}

// ArrayNodeIter is an iterator for ranging over a ArrayNode
type ArrayNodeIter struct {
	values []Node
	idx    int
}

// Next advances the array iterator and reports whether there is another entry.
// It returns false when the iterator is exhausted.
func (m *ArrayNodeIter) Next() bool {
	m.idx++
	next := m.idx < len(m.values)
	return next
}

// Value returns the value of the iterator's current array entry.
func (m *ArrayNodeIter) Value() Node {
	return m.values[m.idx]
}

// Len returns length of array
func (m *ArrayNodeIter) Len() int {
	return len(m.values)
}

// SequenceNode type of sequence node
type SequenceNode struct {
	*BaseNode
	Start             *token.Token
	End               *token.Token
	IsFlowStyle       bool
	Values            []Node
	ValueHeadComments []*CommentGroupNode
	Entries           []*SequenceEntryNode
	FootComment       *CommentGroupNode
}

// Replace replace value node.
func (n *SequenceNode) Replace(idx int, value Node) error {
	if len(n.Values) <= idx {
		return fmt.Errorf(
			"invalid index for sequence: sequence length is %d, but specified %d index",
			len(n.Values), idx,
		)
	}
	column := n.Values[idx].GetToken().Position.Column - value.GetToken().Position.Column
	value.AddColumn(column)
	n.Values[idx] = value
	return nil
}

// Merge merge sequence value.
func (n *SequenceNode) Merge(target *SequenceNode) {
	column := n.Start.Position.Column - target.Start.Position.Column
	target.AddColumn(column)
	n.Values = append(n.Values, target.Values...)
	if len(target.ValueHeadComments) == 0 {
		n.ValueHeadComments = append(n.ValueHeadComments, make([]*CommentGroupNode, len(target.Values))...)
		return
	}
	n.ValueHeadComments = append(n.ValueHeadComments, target.ValueHeadComments...)
}

// SetIsFlowStyle set value to IsFlowStyle field recursively.
func (n *SequenceNode) SetIsFlowStyle(isFlow bool) {
	n.IsFlowStyle = isFlow
	for _, value := range n.Values {
		switch value := value.(type) {
		case *MappingNode:
			value.SetIsFlowStyle(isFlow)
		case *MappingValueNode:
			value.SetIsFlowStyle(isFlow)
		case *SequenceNode:
			value.SetIsFlowStyle(isFlow)
		}
	}
}

// Read implements (io.Reader).Read
func (n *SequenceNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns SequenceType
func (n *SequenceNode) Type() NodeType { return SequenceType }

// GetToken returns token instance
func (n *SequenceNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *SequenceNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	n.End.AddColumn(col)
	for _, value := range n.Values {
		value.AddColumn(col)
	}
}

func (n *SequenceNode) flowStyleString() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

func (n *SequenceNode) blockStyleString() string {
	space := strings.Repeat(" ", n.Start.Position.Column-1)
	values := []string{}
	if n.Comment != nil {
		values = append(values, n.Comment.StringWithSpace(n.Start.Position.Column-1))
	}

	for idx, value := range n.Values {
		if value == nil {
			continue
		}
		valueStr := value.String()
		newLinePrefix := ""
		if strings.HasPrefix(valueStr, "\n") {
			valueStr = valueStr[1:]
			newLinePrefix = "\n"
		}
		splittedValues := strings.Split(valueStr, "\n")
		trimmedFirstValue := strings.TrimLeft(splittedValues[0], " ")
		diffLength := len(splittedValues[0]) - len(trimmedFirstValue)
		if len(splittedValues) > 1 && value.Type() == StringType || value.Type() == LiteralType {
			// If multi-line string, the space characters for indent have already been added, so delete them.
			prefix := space + "  "
			for i := 1; i < len(splittedValues); i++ {
				splittedValues[i] = strings.TrimPrefix(splittedValues[i], prefix)
			}
		}
		newValues := []string{trimmedFirstValue}
		for i := 1; i < len(splittedValues); i++ {
			if len(splittedValues[i]) <= diffLength {
				// this line is \n or white space only
				newValues = append(newValues, "")
				continue
			}
			trimmed := splittedValues[i][diffLength:]
			newValues = append(newValues, fmt.Sprintf("%s  %s", space, trimmed))
		}
		newValue := strings.Join(newValues, "\n")
		if len(n.ValueHeadComments) == len(n.Values) && n.ValueHeadComments[idx] != nil {
			values = append(values, fmt.Sprintf("%s%s", newLinePrefix, n.ValueHeadComments[idx].StringWithSpace(n.Start.Position.Column-1)))
			newLinePrefix = ""
		}
		values = append(values, fmt.Sprintf("%s%s- %s", newLinePrefix, space, newValue))
	}
	if n.FootComment != nil {
		values = append(values, n.FootComment.StringWithSpace(n.Start.Position.Column-1))
	}
	return strings.Join(values, "\n")
}

// String sequence to text
func (n *SequenceNode) String() string {
	if n.IsFlowStyle || len(n.Values) == 0 {
		return n.flowStyleString()
	}
	return n.blockStyleString()
}

// ArrayRange implements ArrayNode protocol
func (n *SequenceNode) ArrayRange() *ArrayNodeIter {
	return &ArrayNodeIter{
		idx:    startRangeIndex,
		values: n.Values,
	}
}

// MarshalYAML encodes to a YAML text
func (n *SequenceNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// SequenceEntryNode is the sequence entry.
type SequenceEntryNode struct {
	*BaseNode
	HeadComment *CommentGroupNode // head comment.
	LineComment *CommentGroupNode // line comment e.g.) - # comment.
	Start       *token.Token      // entry token.
	Value       Node              // value node.
}

// String node to text
func (n *SequenceEntryNode) String() string {
	return "" // TODO
}

// GetToken returns token instance
func (n *SequenceEntryNode) GetToken() *token.Token {
	return n.Start
}

// Type returns type of node
func (n *SequenceEntryNode) Type() NodeType {
	return SequenceEntryType
}

// AddColumn add column number to child nodes recursively
func (n *SequenceEntryNode) AddColumn(col int) {
	n.Start.AddColumn(col)
}

// SetComment set line comment.
func (n *SequenceEntryNode) SetComment(cm *CommentGroupNode) error {
	n.LineComment = cm
	return nil
}

// Comment returns comment token instance
func (n *SequenceEntryNode) GetComment() *CommentGroupNode {
	return n.LineComment
}

// MarshalYAML
func (n *SequenceEntryNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *SequenceEntryNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// SequenceEntry creates SequenceEntryNode instance.
func SequenceEntry(start *token.Token, value Node, headComment *CommentGroupNode) *SequenceEntryNode {
	return &SequenceEntryNode{
		BaseNode:    &BaseNode{},
		HeadComment: headComment,
		Start:       start,
		Value:       value,
	}
}

// SequenceMergeValue creates SequenceMergeValueNode instance.
func SequenceMergeValue(values ...MapNode) *SequenceMergeValueNode {
	return &SequenceMergeValueNode{
		values: values,
	}
}

// SequenceMergeValueNode is used to convert the Sequence node specified for the merge key into a MapNode format.
type SequenceMergeValueNode struct {
	values []MapNode
}

// MapRange returns MapNodeIter instance.
func (n *SequenceMergeValueNode) MapRange() *MapNodeIter {
	ret := &MapNodeIter{idx: startRangeIndex}
	for _, value := range n.values {
		iter := value.MapRange()
		ret.values = append(ret.values, iter.values...)
	}
	return ret
}

// AnchorNode type of anchor node
type AnchorNode struct {
	*BaseNode
	Start *token.Token
	Name  Node
	Value Node
}

func (n *AnchorNode) stringWithoutComment() string {
	return n.Value.String()
}

func (n *AnchorNode) SetName(name string) error {
	if n.Name == nil {
		return ErrInvalidAnchorName
	}
	s, ok := n.Name.(*StringNode)
	if !ok {
		return ErrInvalidAnchorName
	}
	s.Value = name
	return nil
}

// Read implements (io.Reader).Read
func (n *AnchorNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns AnchorType
func (n *AnchorNode) Type() NodeType { return AnchorType }

// GetToken returns token instance
func (n *AnchorNode) GetToken() *token.Token {
	return n.Start
}

func (n *AnchorNode) GetValue() any {
	return n.Value.GetToken().Value
}

// AddColumn add column number to child nodes recursively
func (n *AnchorNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Name != nil {
		n.Name.AddColumn(col)
	}
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// String anchor to text
func (n *AnchorNode) String() string {
	anchor := "&" + n.Name.String()
	value := n.Value.String()
	if s, ok := n.Value.(*SequenceNode); ok && !s.IsFlowStyle {
		return fmt.Sprintf("%s\n%s", anchor, value)
	} else if m, ok := n.Value.(*MappingNode); ok && !m.IsFlowStyle {
		return fmt.Sprintf("%s\n%s", anchor, value)
	}
	if value == "" {
		// implicit null value.
		return anchor
	}
	return fmt.Sprintf("%s %s", anchor, value)
}

// MarshalYAML encodes to a YAML text
func (n *AnchorNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *AnchorNode) IsMergeKey() bool {
	if n.Value == nil {
		return false
	}
	key, ok := n.Value.(MapKeyNode)
	if !ok {
		return false
	}
	return key.IsMergeKey()
}

// AliasNode type of alias node
type AliasNode struct {
	*BaseNode
	Start *token.Token
	Value Node
}

func (n *AliasNode) stringWithoutComment() string {
	return n.Value.String()
}

func (n *AliasNode) SetName(name string) error {
	if n.Value == nil {
		return ErrInvalidAliasName
	}
	s, ok := n.Value.(*StringNode)
	if !ok {
		return ErrInvalidAliasName
	}
	s.Value = name
	return nil
}

// Read implements (io.Reader).Read
func (n *AliasNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns AliasType
func (n *AliasNode) Type() NodeType { return AliasType }

// GetToken returns token instance
func (n *AliasNode) GetToken() *token.Token {
	return n.Start
}

func (n *AliasNode) GetValue() any {
	return n.Value.GetToken().Value
}

// AddColumn add column number to child nodes recursively
func (n *AliasNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// String alias to text
func (n *AliasNode) String() string {
	return fmt.Sprintf("*%s", n.Value.String())
}

// MarshalYAML encodes to a YAML text
func (n *AliasNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *AliasNode) IsMergeKey() bool {
	return false
}

// DirectiveNode type of directive node
type DirectiveNode struct {
	*BaseNode
	// Start is '%' token.
	Start *token.Token
	// Name is directive name e.g.) "YAML" or "TAG".
	Name Node
	// Values is directive values e.g.) "1.2" or "!!" and "tag:clarkevans.com,2002:app/".
	Values []Node
}

// Read implements (io.Reader).Read
func (n *DirectiveNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns DirectiveType
func (n *DirectiveNode) Type() NodeType { return DirectiveType }

// GetToken returns token instance
func (n *DirectiveNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *DirectiveNode) AddColumn(col int) {
	if n.Name != nil {
		n.Name.AddColumn(col)
	}
	for _, value := range n.Values {
		value.AddColumn(col)
	}
}

// String directive to text
func (n *DirectiveNode) String() string {
	values := make([]string, 0, len(n.Values))
	for _, val := range n.Values {
		values = append(values, val.String())
	}
	return strings.Join(append([]string{"%" + n.Name.String()}, values...), " ")
}

// MarshalYAML encodes to a YAML text
func (n *DirectiveNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// TagNode type of tag node
type TagNode struct {
	*BaseNode
	Directive *DirectiveNode
	Start     *token.Token
	Value     Node
}

func (n *TagNode) GetValue() any {
	scalar, ok := n.Value.(ScalarNode)
	if !ok {
		return nil
	}
	return scalar.GetValue()
}

func (n *TagNode) stringWithoutComment() string {
	return n.Value.String()
}

// Read implements (io.Reader).Read
func (n *TagNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns TagType
func (n *TagNode) Type() NodeType { return TagType }

// GetToken returns token instance
func (n *TagNode) GetToken() *token.Token {
	return n.Start
}

// AddColumn add column number to child nodes recursively
func (n *TagNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// String tag to text
func (n *TagNode) String() string {
	value := n.Value.String()
	if s, ok := n.Value.(*SequenceNode); ok && !s.IsFlowStyle {
		return fmt.Sprintf("%s\n%s", n.Start.Value, value)
	} else if m, ok := n.Value.(*MappingNode); ok && !m.IsFlowStyle {
		return fmt.Sprintf("%s\n%s", n.Start.Value, value)
	}

	return fmt.Sprintf("%s %s", n.Start.Value, value)
}

// MarshalYAML encodes to a YAML text
func (n *TagNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// IsMergeKey returns whether it is a MergeKey node.
func (n *TagNode) IsMergeKey() bool {
	if n.Value == nil {
		return false
	}
	key, ok := n.Value.(MapKeyNode)
	if !ok {
		return false
	}
	return key.IsMergeKey()
}

func (n *TagNode) ArrayRange() *ArrayNodeIter {
	arr, ok := n.Value.(ArrayNode)
	if !ok {
		return nil
	}
	return arr.ArrayRange()
}

// CommentNode type of comment node
type CommentNode struct {
	*BaseNode
	Token *token.Token
}

// Read implements (io.Reader).Read
func (n *CommentNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns TagType
func (n *CommentNode) Type() NodeType { return CommentType }

// GetToken returns token instance
func (n *CommentNode) GetToken() *token.Token { return n.Token }

// AddColumn add column number to child nodes recursively
func (n *CommentNode) AddColumn(col int) {
	if n.Token == nil {
		return
	}
	n.Token.AddColumn(col)
}

// String comment to text
func (n *CommentNode) String() string {
	return fmt.Sprintf("#%s", n.Token.Value)
}

// MarshalYAML encodes to a YAML text
func (n *CommentNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// CommentGroupNode type of comment node
type CommentGroupNode struct {
	*BaseNode
	Comments []*CommentNode
}

// Read implements (io.Reader).Read
func (n *CommentGroupNode) Read(p []byte) (int, error) {
	return readNode(p, n)
}

// Type returns TagType
func (n *CommentGroupNode) Type() NodeType { return CommentType }

// GetToken returns token instance
func (n *CommentGroupNode) GetToken() *token.Token {
	if len(n.Comments) > 0 {
		return n.Comments[0].Token
	}
	return nil
}

// AddColumn add column number to child nodes recursively
func (n *CommentGroupNode) AddColumn(col int) {
	for _, comment := range n.Comments {
		comment.AddColumn(col)
	}
}

// String comment to text
func (n *CommentGroupNode) String() string {
	values := []string{}
	for _, comment := range n.Comments {
		values = append(values, comment.String())
	}
	return strings.Join(values, "\n")
}

func (n *CommentGroupNode) StringWithSpace(col int) string {
	values := []string{}
	space := strings.Repeat(" ", col)
	for _, comment := range n.Comments {
		space := space
		if checkLineBreak(comment.Token) {
			space = fmt.Sprintf("%s%s", "\n", space)
		}
		values = append(values, space+comment.String())
	}
	return strings.Join(values, "\n")
}

// MarshalYAML encodes to a YAML text
func (n *CommentGroupNode) MarshalYAML() ([]byte, error) {
	return []byte(n.String()), nil
}

// Visitor has Visit method that is invokded for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children of node with the visitor w,
// followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(Node) Visitor
}

// Walk traverses an AST in depth-first order: It starts by calling v.Visit(node); node must not be nil.
// If the visitor w returned by v.Visit(node) is not nil,
// Walk is invoked recursively with visitor w for each of the non-nil children of node,
// followed by a call of w.Visit(nil).
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *CommentNode:
	case *NullNode:
		walkComment(v, n.BaseNode)
	case *IntegerNode:
		walkComment(v, n.BaseNode)
	case *FloatNode:
		walkComment(v, n.BaseNode)
	case *StringNode:
		walkComment(v, n.BaseNode)
	case *MergeKeyNode:
		walkComment(v, n.BaseNode)
	case *BoolNode:
		walkComment(v, n.BaseNode)
	case *InfinityNode:
		walkComment(v, n.BaseNode)
	case *NanNode:
		walkComment(v, n.BaseNode)
	case *LiteralNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *DirectiveNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Name)
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *TagNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *DocumentNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Body)
	case *MappingNode:
		walkComment(v, n.BaseNode)
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *MappingKeyNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *MappingValueNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *SequenceNode:
		walkComment(v, n.BaseNode)
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *AnchorNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Name)
		Walk(v, n.Value)
	case *AliasNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	}
}

func walkComment(v Visitor, base *BaseNode) {
	if base == nil {
		return
	}
	if base.Comment == nil {
		return
	}
	Walk(v, base.Comment)
}

type filterWalker struct {
	typ     NodeType
	results []Node
}

func (v *filterWalker) Visit(n Node) Visitor {
	if v.typ == n.Type() {
		v.results = append(v.results, n)
	}
	return v
}

type parentFinder struct {
	target Node
}

func (f *parentFinder) walk(parent, node Node) Node {
	if f.target == node {
		return parent
	}
	switch n := node.(type) {
	case *CommentNode:
		return nil
	case *NullNode:
		return nil
	case *IntegerNode:
		return nil
	case *FloatNode:
		return nil
	case *StringNode:
		return nil
	case *MergeKeyNode:
		return nil
	case *BoolNode:
		return nil
	case *InfinityNode:
		return nil
	case *NanNode:
		return nil
	case *LiteralNode:
		return f.walk(node, n.Value)
	case *DirectiveNode:
		if found := f.walk(node, n.Name); found != nil {
			return found
		}
		for _, value := range n.Values {
			if found := f.walk(node, value); found != nil {
				return found
			}
		}
	case *TagNode:
		return f.walk(node, n.Value)
	case *DocumentNode:
		return f.walk(node, n.Body)
	case *MappingNode:
		for _, value := range n.Values {
			if found := f.walk(node, value); found != nil {
				return found
			}
		}
	case *MappingKeyNode:
		return f.walk(node, n.Value)
	case *MappingValueNode:
		if found := f.walk(node, n.Key); found != nil {
			return found
		}
		return f.walk(node, n.Value)
	case *SequenceNode:
		for _, value := range n.Values {
			if found := f.walk(node, value); found != nil {
				return found
			}
		}
	case *AnchorNode:
		if found := f.walk(node, n.Name); found != nil {
			return found
		}
		return f.walk(node, n.Value)
	case *AliasNode:
		return f.walk(node, n.Value)
	}
	return nil
}

// Parent get parent node from child node.
func Parent(root, child Node) Node {
	finder := &parentFinder{target: child}
	return finder.walk(root, root)
}

// Filter returns a list of nodes that match the given type.
func Filter(typ NodeType, node Node) []Node {
	walker := &filterWalker{typ: typ}
	Walk(walker, node)
	return walker.results
}

// FilterFile returns a list of nodes that match the given type.
func FilterFile(typ NodeType, file *File) []Node {
	results := []Node{}
	for _, doc := range file.Docs {
		walker := &filterWalker{typ: typ}
		Walk(walker, doc)
		results = append(results, walker.results...)
	}
	return results
}

type ErrInvalidMergeType struct {
	dst Node
	src Node
}

func (e *ErrInvalidMergeType) Error() string {
	return fmt.Sprintf("cannot merge %s into %s", e.src.Type(), e.dst.Type())
}

// Merge merge document, map, sequence node.
func Merge(dst Node, src Node) error {
	if doc, ok := src.(*DocumentNode); ok {
		src = doc.Body
	}
	err := &ErrInvalidMergeType{dst: dst, src: src}
	switch dst.Type() {
	case DocumentType:
		node, _ := dst.(*DocumentNode)
		return Merge(node.Body, src)
	case MappingType:
		node, _ := dst.(*MappingNode)
		target, ok := src.(*MappingNode)
		if !ok {
			return err
		}
		node.Merge(target)
		return nil
	case SequenceType:
		node, _ := dst.(*SequenceNode)
		target, ok := src.(*SequenceNode)
		if !ok {
			return err
		}
		node.Merge(target)
		return nil
	}
	return err
}
