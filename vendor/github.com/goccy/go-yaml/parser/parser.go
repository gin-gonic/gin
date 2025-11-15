package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

type Mode uint

const (
	ParseComments Mode = 1 << iota // parse comments and add them to AST
)

// ParseBytes parse from byte slice, and returns ast.File
func ParseBytes(bytes []byte, mode Mode, opts ...Option) (*ast.File, error) {
	tokens := lexer.Tokenize(string(bytes))
	f, err := Parse(tokens, mode, opts...)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Parse parse from token instances, and returns ast.File
func Parse(tokens token.Tokens, mode Mode, opts ...Option) (*ast.File, error) {
	if tk := tokens.InvalidToken(); tk != nil {
		return nil, errors.ErrSyntax(tk.Error, tk)
	}
	p, err := newParser(tokens, mode, opts)
	if err != nil {
		return nil, err
	}
	f, err := p.parse(newContext())
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Parse parse from filename, and returns ast.File
func ParseFile(filename string, mode Mode, opts ...Option) (*ast.File, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f, err := ParseBytes(file, mode, opts...)
	if err != nil {
		return nil, err
	}
	f.Name = filename
	return f, nil
}

type YAMLVersion string

const (
	YAML10 YAMLVersion = "1.0"
	YAML11 YAMLVersion = "1.1"
	YAML12 YAMLVersion = "1.2"
	YAML13 YAMLVersion = "1.3"
)

var yamlVersionMap = map[string]YAMLVersion{
	"1.0": YAML10,
	"1.1": YAML11,
	"1.2": YAML12,
	"1.3": YAML13,
}

type parser struct {
	tokens                []*Token
	pathMap               map[string]ast.Node
	yamlVersion           YAMLVersion
	allowDuplicateMapKey  bool
	secondaryTagDirective *ast.DirectiveNode
}

func newParser(tokens token.Tokens, mode Mode, opts []Option) (*parser, error) {
	filteredTokens := []*token.Token{}
	if mode&ParseComments != 0 {
		filteredTokens = tokens
	} else {
		for _, tk := range tokens {
			if tk.Type == token.CommentType {
				continue
			}
			// keep prev/next reference between tokens containing comments
			// https://github.com/goccy/go-yaml/issues/254
			filteredTokens = append(filteredTokens, tk)
		}
	}
	tks, err := CreateGroupedTokens(token.Tokens(filteredTokens))
	if err != nil {
		return nil, err
	}
	p := &parser{
		tokens:  tks,
		pathMap: make(map[string]ast.Node),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *parser) parse(ctx *context) (*ast.File, error) {
	file := &ast.File{Docs: []*ast.DocumentNode{}}
	for _, token := range p.tokens {
		doc, err := p.parseDocument(ctx, token.Group)
		if err != nil {
			return nil, err
		}
		file.Docs = append(file.Docs, doc)
	}
	return file, nil
}

func (p *parser) parseDocument(ctx *context, docGroup *TokenGroup) (*ast.DocumentNode, error) {
	if len(docGroup.Tokens) == 0 {
		return ast.Document(docGroup.RawToken(), nil), nil
	}

	p.pathMap = make(map[string]ast.Node)

	var (
		tokens = docGroup.Tokens
		start  *token.Token
		end    *token.Token
	)
	if docGroup.First().Type() == token.DocumentHeaderType {
		start = docGroup.First().RawToken()
		tokens = tokens[1:]
	}
	if docGroup.Last().Type() == token.DocumentEndType {
		end = docGroup.Last().RawToken()
		tokens = tokens[:len(tokens)-1]
		defer func() {
			// clear yaml version value if DocumentEnd token (...) is specified.
			p.yamlVersion = ""
		}()
	}

	if len(tokens) == 0 {
		return ast.Document(docGroup.RawToken(), nil), nil
	}

	body, err := p.parseDocumentBody(ctx.withGroup(&TokenGroup{
		Type:   TokenGroupDocumentBody,
		Tokens: tokens,
	}))
	if err != nil {
		return nil, err
	}
	node := ast.Document(start, body)
	node.End = end
	return node, nil
}

func (p *parser) parseDocumentBody(ctx *context) (ast.Node, error) {
	node, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if ctx.next() {
		return nil, errors.ErrSyntax("value is not allowed in this context", ctx.currentToken().RawToken())
	}
	return node, nil
}

func (p *parser) parseToken(ctx *context, tk *Token) (ast.Node, error) {
	switch tk.GroupType() {
	case TokenGroupMapKey, TokenGroupMapKeyValue:
		return p.parseMap(ctx)
	case TokenGroupDirective:
		node, err := p.parseDirective(ctx.withGroup(tk.Group), tk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupDirectiveName:
		node, err := p.parseDirectiveName(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupAnchor:
		node, err := p.parseAnchor(ctx.withGroup(tk.Group), tk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupAnchorName:
		anchor, err := p.parseAnchorName(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		if ctx.isTokenNotFound() {
			return nil, errors.ErrSyntax("could not find anchor value", tk.RawToken())
		}
		value, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		if _, ok := value.(*ast.AnchorNode); ok {
			return nil, errors.ErrSyntax("anchors cannot be used consecutively", value.GetToken())
		}
		anchor.Value = value
		return anchor, nil
	case TokenGroupAlias:
		node, err := p.parseAlias(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupLiteral, TokenGroupFolded:
		node, err := p.parseLiteral(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupScalarTag:
		node, err := p.parseTag(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	}
	switch tk.Type() {
	case token.CommentType:
		return p.parseComment(ctx)
	case token.TagType:
		return p.parseTag(ctx)
	case token.MappingStartType:
		return p.parseFlowMap(ctx.withFlow(true))
	case token.SequenceStartType:
		return p.parseFlowSequence(ctx.withFlow(true))
	case token.SequenceEntryType:
		return p.parseSequence(ctx)
	case token.SequenceEndType:
		// SequenceEndType is always validated in parseFlowSequence.
		// Therefore, if this is found in other cases, it is treated as a syntax error.
		return nil, errors.ErrSyntax("could not find '[' character corresponding to ']'", tk.RawToken())
	case token.MappingEndType:
		// MappingEndType is always validated in parseFlowMap.
		// Therefore, if this is found in other cases, it is treated as a syntax error.
		return nil, errors.ErrSyntax("could not find '{' character corresponding to '}'", tk.RawToken())
	case token.MappingValueType:
		return nil, errors.ErrSyntax("found an invalid key for this map", tk.RawToken())
	}
	node, err := p.parseScalarValue(ctx, tk)
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	return node, nil
}

func (p *parser) parseScalarValue(ctx *context, tk *Token) (ast.ScalarNode, error) {
	if tk.Group != nil {
		switch tk.GroupType() {
		case TokenGroupAnchor:
			return p.parseAnchor(ctx.withGroup(tk.Group), tk.Group)
		case TokenGroupAnchorName:
			anchor, err := p.parseAnchorName(ctx.withGroup(tk.Group))
			if err != nil {
				return nil, err
			}
			ctx.goNext()
			if ctx.isTokenNotFound() {
				return nil, errors.ErrSyntax("could not find anchor value", tk.RawToken())
			}
			value, err := p.parseToken(ctx, ctx.currentToken())
			if err != nil {
				return nil, err
			}
			if _, ok := value.(*ast.AnchorNode); ok {
				return nil, errors.ErrSyntax("anchors cannot be used consecutively", value.GetToken())
			}
			anchor.Value = value
			return anchor, nil
		case TokenGroupAlias:
			return p.parseAlias(ctx.withGroup(tk.Group))
		case TokenGroupLiteral, TokenGroupFolded:
			return p.parseLiteral(ctx.withGroup(tk.Group))
		case TokenGroupScalarTag:
			return p.parseTag(ctx.withGroup(tk.Group))
		default:
			return nil, errors.ErrSyntax("unexpected scalar value", tk.RawToken())
		}
	}
	switch tk.Type() {
	case token.MergeKeyType:
		return newMergeKeyNode(ctx, tk)
	case token.NullType, token.ImplicitNullType:
		return newNullNode(ctx, tk)
	case token.BoolType:
		return newBoolNode(ctx, tk)
	case token.IntegerType, token.BinaryIntegerType, token.OctetIntegerType, token.HexIntegerType:
		return newIntegerNode(ctx, tk)
	case token.FloatType:
		return newFloatNode(ctx, tk)
	case token.InfinityType:
		return newInfinityNode(ctx, tk)
	case token.NanType:
		return newNanNode(ctx, tk)
	case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
		return newStringNode(ctx, tk)
	case token.TagType:
		// this case applies when it is a scalar tag and its value does not exist.
		// Examples of cases where the value does not exist include cases like `key: !!str,` or `!!str : value`.
		return p.parseScalarTag(ctx)
	}
	return nil, errors.ErrSyntax("unexpected scalar value type", tk.RawToken())
}

func (p *parser) parseFlowMap(ctx *context) (*ast.MappingNode, error) {
	node, err := newMappingNode(ctx, ctx.currentToken(), true)
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip MappingStart token

	isFirst := true
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type() == token.MappingEndType {
			node.End = tk.RawToken()
			break
		}

		var entryTk *Token
		if tk.Type() == token.CollectEntryType {
			entryTk = tk
			ctx.goNext()
		} else if !isFirst {
			return nil, errors.ErrSyntax("',' or '}' must be specified", tk.RawToken())
		}

		if tk := ctx.currentToken(); tk.Type() == token.MappingEndType {
			// this case is here: "{ elem, }".
			// In this case, ignore the last element and break mapping parsing.
			node.End = tk.RawToken()
			break
		}

		mapKeyTk := ctx.currentToken()
		switch mapKeyTk.GroupType() {
		case TokenGroupMapKeyValue:
			value, err := p.parseMapKeyValue(ctx.withGroup(mapKeyTk.Group), mapKeyTk.Group, entryTk)
			if err != nil {
				return nil, err
			}
			node.Values = append(node.Values, value)
			ctx.goNext()
		case TokenGroupMapKey:
			key, err := p.parseMapKey(ctx.withGroup(mapKeyTk.Group), mapKeyTk.Group)
			if err != nil {
				return nil, err
			}
			ctx := ctx.withChild(p.mapKeyText(key))
			colonTk := mapKeyTk.Group.Last()
			if p.isFlowMapDelim(ctx.nextToken()) {
				value, err := newNullNode(ctx, ctx.insertNullToken(colonTk))
				if err != nil {
					return nil, err
				}
				mapValue, err := newMappingValueNode(ctx, colonTk, entryTk, key, value)
				if err != nil {
					return nil, err
				}
				node.Values = append(node.Values, mapValue)
				ctx.goNext()
			} else {
				ctx.goNext()
				if ctx.isTokenNotFound() {
					return nil, errors.ErrSyntax("could not find map value", colonTk.RawToken())
				}
				value, err := p.parseToken(ctx, ctx.currentToken())
				if err != nil {
					return nil, err
				}
				mapValue, err := newMappingValueNode(ctx, colonTk, entryTk, key, value)
				if err != nil {
					return nil, err
				}
				node.Values = append(node.Values, mapValue)
			}
		default:
			if !p.isFlowMapDelim(ctx.nextToken()) {
				errTk := mapKeyTk
				if errTk == nil {
					errTk = tk
				}
				return nil, errors.ErrSyntax("could not find flow map content", errTk.RawToken())
			}
			key, err := p.parseScalarValue(ctx, mapKeyTk)
			if err != nil {
				return nil, err
			}
			value, err := newNullNode(ctx, ctx.insertNullToken(mapKeyTk))
			if err != nil {
				return nil, err
			}
			mapValue, err := newMappingValueNode(ctx, mapKeyTk, entryTk, key, value)
			if err != nil {
				return nil, err
			}
			node.Values = append(node.Values, mapValue)
			ctx.goNext()
		}
		isFirst = false
	}
	if node.End == nil {
		return nil, errors.ErrSyntax("could not find flow mapping end token '}'", node.Start)
	}
	ctx.goNext() // skip mapping end token.
	return node, nil
}

func (p *parser) isFlowMapDelim(tk *Token) bool {
	return tk.Type() == token.MappingEndType || tk.Type() == token.CollectEntryType
}

func (p *parser) parseMap(ctx *context) (*ast.MappingNode, error) {
	keyTk := ctx.currentToken()
	if keyTk.Group == nil {
		return nil, errors.ErrSyntax("unexpected map key", keyTk.RawToken())
	}
	var keyValueNode *ast.MappingValueNode
	if keyTk.GroupType() == TokenGroupMapKeyValue {
		node, err := p.parseMapKeyValue(ctx.withGroup(keyTk.Group), keyTk.Group, nil)
		if err != nil {
			return nil, err
		}
		keyValueNode = node
		ctx.goNext()
		if err := p.validateMapKeyValueNextToken(ctx, keyTk, ctx.currentToken()); err != nil {
			return nil, err
		}
	} else {
		key, err := p.parseMapKey(ctx.withGroup(keyTk.Group), keyTk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()

		valueTk := ctx.currentToken()
		if keyTk.Line() == valueTk.Line() && valueTk.Type() == token.SequenceEntryType {
			return nil, errors.ErrSyntax("block sequence entries are not allowed in this context", valueTk.RawToken())
		}
		ctx := ctx.withChild(p.mapKeyText(key))
		value, err := p.parseMapValue(ctx, key, keyTk.Group.Last())
		if err != nil {
			return nil, err
		}
		node, err := newMappingValueNode(ctx, keyTk.Group.Last(), nil, key, value)
		if err != nil {
			return nil, err
		}
		keyValueNode = node
	}
	mapNode, err := newMappingNode(ctx, &Token{Token: keyValueNode.GetToken()}, false, keyValueNode)
	if err != nil {
		return nil, err
	}
	var tk *Token
	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	} else {
		tk = ctx.currentToken()
	}
	for tk.Column() == keyTk.Column() {
		typ := tk.Type()
		if ctx.isFlow && typ == token.SequenceEndType {
			// [
			// key: value
			// ] <=
			break
		}
		if !p.isMapToken(tk) {
			return nil, errors.ErrSyntax("non-map value is specified", tk.RawToken())
		}
		cm := p.parseHeadComment(ctx)
		if typ == token.MappingEndType {
			// a: {
			//  b: c
			// } <=
			ctx.goNext()
			break
		}
		node, err := p.parseMap(ctx)
		if err != nil {
			return nil, err
		}
		if len(node.Values) != 0 {
			if err := setHeadComment(cm, node.Values[0]); err != nil {
				return nil, err
			}
		}
		mapNode.Values = append(mapNode.Values, node.Values...)
		if node.FootComment != nil {
			mapNode.Values[len(mapNode.Values)-1].FootComment = node.FootComment
		}
		tk = ctx.currentToken()
	}
	if ctx.isComment() {
		if keyTk.Column() <= ctx.currentToken().Column() {
			// If the comment is in the same or deeper column as the last element column in map value,
			// treat it as a footer comment for the last element.
			if len(mapNode.Values) == 1 {
				mapNode.Values[0].FootComment = p.parseFootComment(ctx, keyTk.Column())
				mapNode.Values[0].FootComment.SetPath(mapNode.Values[0].Key.GetPath())
			} else {
				mapNode.FootComment = p.parseFootComment(ctx, keyTk.Column())
				mapNode.FootComment.SetPath(mapNode.GetPath())
			}
		}
	}
	return mapNode, nil
}

func (p *parser) validateMapKeyValueNextToken(ctx *context, keyTk, tk *Token) error {
	if tk == nil {
		return nil
	}
	if tk.Column() <= keyTk.Column() {
		return nil
	}
	if ctx.isComment() {
		return nil
	}
	if ctx.isFlow && (tk.Type() == token.CollectEntryType || tk.Type() == token.SequenceEndType) {
		return nil
	}
	// a: b
	//  c <= this token is invalid.
	return errors.ErrSyntax("value is not allowed in this context. map key-value is pre-defined", tk.RawToken())
}

func (p *parser) isMapToken(tk *Token) bool {
	if tk.Group == nil {
		return tk.Type() == token.MappingStartType || tk.Type() == token.MappingEndType
	}
	g := tk.Group
	return g.Type == TokenGroupMapKey || g.Type == TokenGroupMapKeyValue
}

func (p *parser) parseMapKeyValue(ctx *context, g *TokenGroup, entryTk *Token) (*ast.MappingValueNode, error) {
	if g.Type != TokenGroupMapKeyValue {
		return nil, errors.ErrSyntax("unexpected map key-value pair", g.RawToken())
	}
	if g.First().Group == nil {
		return nil, errors.ErrSyntax("unexpected map key", g.RawToken())
	}
	keyGroup := g.First().Group
	key, err := p.parseMapKey(ctx.withGroup(keyGroup), keyGroup)
	if err != nil {
		return nil, err
	}

	c := ctx.withChild(p.mapKeyText(key))
	value, err := p.parseToken(c, g.Last())
	if err != nil {
		return nil, err
	}
	return newMappingValueNode(c, keyGroup.Last(), entryTk, key, value)
}

func (p *parser) parseMapKey(ctx *context, g *TokenGroup) (ast.MapKeyNode, error) {
	if g.Type != TokenGroupMapKey {
		return nil, errors.ErrSyntax("unexpected map key", g.RawToken())
	}
	if g.First().Type() == token.MappingKeyType {
		mapKeyTk := g.First()
		if mapKeyTk.Group != nil {
			ctx = ctx.withGroup(mapKeyTk.Group)
		}
		key, err := newMappingKeyNode(ctx, mapKeyTk)
		if err != nil {
			return nil, err
		}
		ctx.goNext() // skip mapping key token
		if ctx.isTokenNotFound() {
			return nil, errors.ErrSyntax("could not find value for mapping key", mapKeyTk.RawToken())
		}

		scalar, err := p.parseScalarValue(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		key.Value = scalar
		keyText := p.mapKeyText(scalar)
		keyPath := ctx.withChild(keyText).path
		key.SetPath(keyPath)
		if err := p.validateMapKey(ctx, key.GetToken(), keyPath, g.Last()); err != nil {
			return nil, err
		}
		p.pathMap[keyPath] = key
		return key, nil
	}
	if g.Last().Type() != token.MappingValueType {
		return nil, errors.ErrSyntax("expected map key-value delimiter ':'", g.Last().RawToken())
	}

	scalar, err := p.parseScalarValue(ctx, g.First())
	if err != nil {
		return nil, err
	}
	key, ok := scalar.(ast.MapKeyNode)
	if !ok {
		return nil, errors.ErrSyntax("cannot take map-key node", scalar.GetToken())
	}
	keyText := p.mapKeyText(key)
	keyPath := ctx.withChild(keyText).path
	key.SetPath(keyPath)
	if err := p.validateMapKey(ctx, key.GetToken(), keyPath, g.Last()); err != nil {
		return nil, err
	}
	p.pathMap[keyPath] = key
	return key, nil
}

func (p *parser) validateMapKey(ctx *context, tk *token.Token, keyPath string, colonTk *Token) error {
	if !p.allowDuplicateMapKey {
		if n, exists := p.pathMap[keyPath]; exists {
			pos := n.GetToken().Position
			return errors.ErrSyntax(
				fmt.Sprintf("mapping key %q already defined at [%d:%d]", tk.Value, pos.Line, pos.Column),
				tk,
			)
		}
	}
	origin := p.removeLeftWhiteSpace(tk.Origin)
	if ctx.isFlow {
		if tk.Type == token.StringType {
			origin = p.removeRightWhiteSpace(origin)
			if tk.Position.Line+p.newLineCharacterNum(origin) != colonTk.Line() {
				return errors.ErrSyntax("map key definition includes an implicit line break", tk)
			}
		}
		return nil
	}
	if tk.Type != token.StringType && tk.Type != token.SingleQuoteType && tk.Type != token.DoubleQuoteType {
		return nil
	}
	if p.existsNewLineCharacter(origin) {
		return errors.ErrSyntax("unexpected key name", tk)
	}
	return nil
}

func (p *parser) removeLeftWhiteSpace(src string) string {
	// CR or LF or CRLF
	return strings.TrimLeftFunc(src, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})
}

func (p *parser) removeRightWhiteSpace(src string) string {
	// CR or LF or CRLF
	return strings.TrimRightFunc(src, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})
}

func (p *parser) existsNewLineCharacter(src string) bool {
	return p.newLineCharacterNum(src) > 0
}

func (p *parser) newLineCharacterNum(src string) int {
	var num int
	for i := 0; i < len(src); i++ {
		switch src[i] {
		case '\r':
			if len(src) > i+1 && src[i+1] == '\n' {
				i++
			}
			num++
		case '\n':
			num++
		}
	}
	return num
}

func (p *parser) mapKeyText(n ast.Node) string {
	if n == nil {
		return ""
	}
	switch nn := n.(type) {
	case *ast.MappingKeyNode:
		return p.mapKeyText(nn.Value)
	case *ast.TagNode:
		return p.mapKeyText(nn.Value)
	case *ast.AnchorNode:
		return p.mapKeyText(nn.Value)
	case *ast.AliasNode:
		return ""
	}
	return n.GetToken().Value
}

func (p *parser) parseMapValue(ctx *context, key ast.MapKeyNode, colonTk *Token) (ast.Node, error) {
	tk := ctx.currentToken()
	if tk == nil {
		return newNullNode(ctx, ctx.addNullValueToken(colonTk))
	}

	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	}
	keyCol := key.GetToken().Position.Column
	keyLine := key.GetToken().Position.Line

	if tk.Column() != keyCol && tk.Line() == keyLine && (tk.GroupType() == TokenGroupMapKey || tk.GroupType() == TokenGroupMapKeyValue) {
		// a: b:
		//    ^
		//
		// a: b: c
		//    ^
		return nil, errors.ErrSyntax("mapping value is not allowed in this context", tk.RawToken())
	}

	if tk.Column() == keyCol && p.isMapToken(tk) {
		// in this case,
		// ----
		// key: <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(colonTk))
	}

	if tk.Line() == keyLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() == keyCol && p.isMapToken(ctx.nextToken()) {
		// in this case,
		// ----
		// key: &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createImplicitNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	if tk.Column() <= keyCol && tk.GroupType() == TokenGroupAnchorName {
		// key: <value does not defined>
		// &anchor
		return nil, errors.ErrSyntax("anchor is not allowed in this context", tk.RawToken())
	}
	if tk.Column() <= keyCol && tk.Type() == token.TagType {
		// key: <value does not defined>
		// !!tag
		return nil, errors.ErrSyntax("tag is not allowed in this context", tk.RawToken())
	}

	if tk.Column() < keyCol {
		// in this case,
		// ----
		//   key: <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(colonTk))
	}

	if tk.Line() == keyLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() < keyCol {
		// in this case,
		// ----
		//   key: &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createImplicitNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if err := p.validateAnchorValueInMapOrSeq(value, keyCol); err != nil {
		return nil, err
	}
	return value, nil
}

func (p *parser) validateAnchorValueInMapOrSeq(value ast.Node, col int) error {
	anchor, ok := value.(*ast.AnchorNode)
	if !ok {
		return nil
	}
	tag, ok := anchor.Value.(*ast.TagNode)
	if !ok {
		return nil
	}
	anchorTk := anchor.GetToken()
	tagTk := tag.GetToken()

	if anchorTk.Position.Line == tagTk.Position.Line {
		// key:
		//   &anchor !!tag
		//
		// - &anchor !!tag
		return nil
	}

	if tagTk.Position.Column <= col {
		// key: &anchor
		// !!tag
		//
		// - &anchor
		// !!tag
		return errors.ErrSyntax("tag is not allowed in this context", tagTk)
	}
	return nil
}

func (p *parser) parseAnchor(ctx *context, g *TokenGroup) (*ast.AnchorNode, error) {
	anchorNameGroup := g.First().Group
	anchor, err := p.parseAnchorName(ctx.withGroup(anchorNameGroup))
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find anchor value", anchor.GetToken())
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if _, ok := value.(*ast.AnchorNode); ok {
		return nil, errors.ErrSyntax("anchors cannot be used consecutively", value.GetToken())
	}
	anchor.Value = value
	return anchor, nil
}

func (p *parser) parseAnchorName(ctx *context) (*ast.AnchorNode, error) {
	anchor, err := newAnchorNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find anchor value", anchor.GetToken())
	}

	anchorName, err := p.parseScalarValue(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if anchorName == nil {
		return nil, errors.ErrSyntax("unexpected anchor. anchor name is not scalar value", ctx.currentToken().RawToken())
	}
	anchor.Name = anchorName
	return anchor, nil
}

func (p *parser) parseAlias(ctx *context) (*ast.AliasNode, error) {
	alias, err := newAliasNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find alias value", alias.GetToken())
	}

	aliasName, err := p.parseScalarValue(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if aliasName == nil {
		return nil, errors.ErrSyntax("unexpected alias. alias name is not scalar value", ctx.currentToken().RawToken())
	}
	alias.Value = aliasName
	return alias, nil
}

func (p *parser) parseLiteral(ctx *context) (*ast.LiteralNode, error) {
	node, err := newLiteralNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip literal/folded token

	tk := ctx.currentToken()
	if tk == nil {
		value, err := newStringNode(ctx, &Token{Token: token.New("", "", node.Start.Position)})
		if err != nil {
			return nil, err
		}
		node.Value = value
		return node, nil
	}
	value, err := p.parseToken(ctx, tk)
	if err != nil {
		return nil, err
	}
	str, ok := value.(*ast.StringNode)
	if !ok {
		return nil, errors.ErrSyntax("unexpected token. required string token", value.GetToken())
	}
	node.Value = str
	return node, nil
}

func (p *parser) parseScalarTag(ctx *context) (*ast.TagNode, error) {
	tag, err := p.parseTag(ctx)
	if err != nil {
		return nil, err
	}
	if tag.Value == nil {
		return nil, errors.ErrSyntax("specified not scalar tag", tag.GetToken())
	}
	if _, ok := tag.Value.(ast.ScalarNode); !ok {
		return nil, errors.ErrSyntax("specified not scalar tag", tag.GetToken())
	}
	return tag, nil
}

func (p *parser) parseTag(ctx *context) (*ast.TagNode, error) {
	tagTk := ctx.currentToken()
	tagRawTk := tagTk.RawToken()
	node, err := newTagNode(ctx, tagTk)
	if err != nil {
		return nil, err
	}
	ctx.goNext()

	comment := p.parseHeadComment(ctx)

	var tagValue ast.Node
	if p.secondaryTagDirective != nil {
		value, err := newStringNode(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		tagValue = value
		node.Directive = p.secondaryTagDirective
	} else {
		value, err := p.parseTagValue(ctx, tagRawTk, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		tagValue = value
	}
	if err := setHeadComment(comment, tagValue); err != nil {
		return nil, err
	}
	node.Value = tagValue
	return node, nil
}

func (p *parser) parseTagValue(ctx *context, tagRawTk *token.Token, tk *Token) (ast.Node, error) {
	if tk == nil {
		return newNullNode(ctx, ctx.createImplicitNullToken(&Token{Token: tagRawTk}))
	}
	switch token.ReservedTagKeyword(tagRawTk.Value) {
	case token.MappingTag, token.SetTag:
		if !p.isMapToken(tk) {
			return nil, errors.ErrSyntax("could not find map", tk.RawToken())
		}
		if tk.Type() == token.MappingStartType {
			return p.parseFlowMap(ctx.withFlow(true))
		}
		return p.parseMap(ctx)
	case token.IntegerTag, token.FloatTag, token.StringTag, token.BinaryTag, token.TimestampTag, token.BooleanTag, token.NullTag:
		if tk.GroupType() == TokenGroupLiteral || tk.GroupType() == TokenGroupFolded {
			return p.parseLiteral(ctx.withGroup(tk.Group))
		} else if tk.Type() == token.CollectEntryType || tk.Type() == token.MappingValueType {
			return newTagDefaultScalarValueNode(ctx, tagRawTk)
		}
		scalar, err := p.parseScalarValue(ctx, tk)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return scalar, nil
	case token.SequenceTag, token.OrderedMapTag:
		if tk.Type() == token.SequenceStartType {
			return p.parseFlowSequence(ctx.withFlow(true))
		}
		return p.parseSequence(ctx)
	}
	return p.parseToken(ctx, tk)
}

func (p *parser) parseFlowSequence(ctx *context) (*ast.SequenceNode, error) {
	node, err := newSequenceNode(ctx, ctx.currentToken(), true)
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip SequenceStart token

	isFirst := true
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type() == token.SequenceEndType {
			node.End = tk.RawToken()
			break
		}

		var entryTk *Token
		if tk.Type() == token.CollectEntryType {
			if isFirst {
				return nil, errors.ErrSyntax("expected sequence element, but found ','", tk.RawToken())
			}
			entryTk = tk
			ctx.goNext()
		} else if !isFirst {
			return nil, errors.ErrSyntax("',' or ']' must be specified", tk.RawToken())
		}

		if tk := ctx.currentToken(); tk.Type() == token.SequenceEndType {
			// this case is here: "[ elem, ]".
			// In this case, ignore the last element and break sequence parsing.
			node.End = tk.RawToken()
			break
		}

		if ctx.isTokenNotFound() {
			break
		}

		ctx := ctx.withIndex(uint(len(node.Values)))
		value, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		node.Values = append(node.Values, value)
		seqEntry := ast.SequenceEntry(entryTk.RawToken(), value, nil)
		if err := setLineComment(ctx, seqEntry, entryTk); err != nil {
			return nil, err
		}
		seqEntry.SetPath(ctx.path)
		node.Entries = append(node.Entries, seqEntry)

		isFirst = false
	}
	if node.End == nil {
		return nil, errors.ErrSyntax("sequence end token ']' not found", node.Start)
	}
	ctx.goNext() // skip sequence end token.
	return node, nil
}

func (p *parser) parseSequence(ctx *context) (*ast.SequenceNode, error) {
	seqTk := ctx.currentToken()
	seqNode, err := newSequenceNode(ctx, seqTk, false)
	if err != nil {
		return nil, err
	}

	tk := seqTk
	for tk.Type() == token.SequenceEntryType && tk.Column() == seqTk.Column() {
		seqTk := tk
		headComment := p.parseHeadComment(ctx)
		ctx.goNext() // skip sequence entry token

		ctx := ctx.withIndex(uint(len(seqNode.Values)))
		value, err := p.parseSequenceValue(ctx, seqTk)
		if err != nil {
			return nil, err
		}
		seqEntry := ast.SequenceEntry(seqTk.RawToken(), value, headComment)
		if err := setLineComment(ctx, seqEntry, seqTk); err != nil {
			return nil, err
		}
		seqEntry.SetPath(ctx.path)
		seqNode.ValueHeadComments = append(seqNode.ValueHeadComments, headComment)
		seqNode.Values = append(seqNode.Values, value)
		seqNode.Entries = append(seqNode.Entries, seqEntry)

		if ctx.isComment() {
			tk = ctx.nextNotCommentToken()
		} else {
			tk = ctx.currentToken()
		}
	}
	if ctx.isComment() {
		if seqTk.Column() <= ctx.currentToken().Column() {
			// If the comment is in the same or deeper column as the last element column in sequence value,
			// treat it as a footer comment for the last element.
			seqNode.FootComment = p.parseFootComment(ctx, seqTk.Column())
			if len(seqNode.Values) != 0 {
				seqNode.FootComment.SetPath(seqNode.Values[len(seqNode.Values)-1].GetPath())
			}
		}
	}
	return seqNode, nil
}

func (p *parser) parseSequenceValue(ctx *context, seqTk *Token) (ast.Node, error) {
	tk := ctx.currentToken()
	if tk == nil {
		return newNullNode(ctx, ctx.addNullValueToken(seqTk))
	}

	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	}
	seqCol := seqTk.Column()
	seqLine := seqTk.Line()

	if tk.Column() == seqCol && tk.Type() == token.SequenceEntryType {
		// in this case,
		// ----
		// - <value does not defined>
		// -
		return newNullNode(ctx, ctx.insertNullToken(seqTk))
	}

	if tk.Line() == seqLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() == seqCol && ctx.nextToken().Type() == token.SequenceEntryType {
		// in this case,
		// ----
		// - &anchor
		// -
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createImplicitNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	if tk.Column() <= seqCol && tk.GroupType() == TokenGroupAnchorName {
		// - <value does not defined>
		// &anchor
		return nil, errors.ErrSyntax("anchor is not allowed in this sequence context", tk.RawToken())
	}
	if tk.Column() <= seqCol && tk.Type() == token.TagType {
		// - <value does not defined>
		// !!tag
		return nil, errors.ErrSyntax("tag is not allowed in this sequence context", tk.RawToken())
	}

	if tk.Column() < seqCol {
		// in this case,
		// ----
		//   - <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(seqTk))
	}

	if tk.Line() == seqLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() < seqCol {
		// in this case,
		// ----
		//   - &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createImplicitNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if err := p.validateAnchorValueInMapOrSeq(value, seqCol); err != nil {
		return nil, err
	}
	return value, nil
}

func (p *parser) parseDirective(ctx *context, g *TokenGroup) (*ast.DirectiveNode, error) {
	directiveNameGroup := g.First().Group
	directive, err := p.parseDirectiveName(ctx.withGroup(directiveNameGroup))
	if err != nil {
		return nil, err
	}

	switch directive.Name.String() {
	case "YAML":
		if len(g.Tokens) != 2 {
			return nil, errors.ErrSyntax("unexpected format YAML directive", g.First().RawToken())
		}
		valueTk := g.Tokens[1]
		valueRawTk := valueTk.RawToken()
		value := valueRawTk.Value
		ver, exists := yamlVersionMap[value]
		if !exists {
			return nil, errors.ErrSyntax(fmt.Sprintf("unknown YAML version %q", value), valueRawTk)
		}
		if p.yamlVersion != "" {
			return nil, errors.ErrSyntax("YAML version has already been specified", valueRawTk)
		}
		p.yamlVersion = ver
		versionNode, err := newStringNode(ctx, valueTk)
		if err != nil {
			return nil, err
		}
		directive.Values = append(directive.Values, versionNode)
	case "TAG":
		if len(g.Tokens) != 3 {
			return nil, errors.ErrSyntax("unexpected format TAG directive", g.First().RawToken())
		}
		tagKey, err := newStringNode(ctx, g.Tokens[1])
		if err != nil {
			return nil, err
		}
		if tagKey.Value == "!!" {
			p.secondaryTagDirective = directive
		}
		tagValue, err := newStringNode(ctx, g.Tokens[2])
		if err != nil {
			return nil, err
		}
		directive.Values = append(directive.Values, tagKey, tagValue)
	default:
		if len(g.Tokens) > 1 {
			for _, tk := range g.Tokens[1:] {
				value, err := newStringNode(ctx, tk)
				if err != nil {
					return nil, err
				}
				directive.Values = append(directive.Values, value)
			}
		}
	}
	return directive, nil
}

func (p *parser) parseDirectiveName(ctx *context) (*ast.DirectiveNode, error) {
	directive, err := newDirectiveNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find directive value", directive.GetToken())
	}

	directiveName, err := p.parseScalarValue(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if directiveName == nil {
		return nil, errors.ErrSyntax("unexpected directive. directive name is not scalar value", ctx.currentToken().RawToken())
	}
	directive.Name = directiveName
	return directive, nil
}

func (p *parser) parseComment(ctx *context) (ast.Node, error) {
	cm := p.parseHeadComment(ctx)
	if ctx.isTokenNotFound() {
		return cm, nil
	}
	node, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if err := setHeadComment(cm, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) parseHeadComment(ctx *context) *ast.CommentGroupNode {
	tks := []*token.Token{}
	for ctx.isComment() {
		tks = append(tks, ctx.currentToken().RawToken())
		ctx.goNext()
	}
	if len(tks) == 0 {
		return nil
	}
	return ast.CommentGroup(tks)
}

func (p *parser) parseFootComment(ctx *context, col int) *ast.CommentGroupNode {
	tks := []*token.Token{}
	for ctx.isComment() && col <= ctx.currentToken().Column() {
		tks = append(tks, ctx.currentToken().RawToken())
		ctx.goNext()
	}
	if len(tks) == 0 {
		return nil
	}
	return ast.CommentGroup(tks)
}
