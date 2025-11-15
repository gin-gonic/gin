package format

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

func FormatNodeWithResolvedAlias(n ast.Node, anchorNodeMap map[string]ast.Node) string {
	tk := getFirstToken(n)
	if tk == nil {
		return ""
	}
	formatter := newFormatter(tk, hasComment(n))
	formatter.anchorNodeMap = anchorNodeMap
	return formatter.format(n)
}

func FormatNode(n ast.Node) string {
	tk := getFirstToken(n)
	if tk == nil {
		return ""
	}
	return newFormatter(tk, hasComment(n)).format(n)
}

func FormatFile(file *ast.File) string {
	if len(file.Docs) == 0 {
		return ""
	}
	tk := getFirstToken(file.Docs[0])
	if tk == nil {
		return ""
	}
	return newFormatter(tk, hasCommentFile(file)).formatFile(file)
}

func hasCommentFile(f *ast.File) bool {
	for _, doc := range f.Docs {
		if hasComment(doc.Body) {
			return true
		}
	}
	return false
}

func hasComment(n ast.Node) bool {
	if n == nil {
		return false
	}
	switch nn := n.(type) {
	case *ast.DocumentNode:
		return hasComment(nn.Body)
	case *ast.NullNode:
		return nn.Comment != nil
	case *ast.BoolNode:
		return nn.Comment != nil
	case *ast.IntegerNode:
		return nn.Comment != nil
	case *ast.FloatNode:
		return nn.Comment != nil
	case *ast.StringNode:
		return nn.Comment != nil
	case *ast.InfinityNode:
		return nn.Comment != nil
	case *ast.NanNode:
		return nn.Comment != nil
	case *ast.LiteralNode:
		return nn.Comment != nil
	case *ast.DirectiveNode:
		if nn.Comment != nil {
			return true
		}
		for _, value := range nn.Values {
			if hasComment(value) {
				return true
			}
		}
	case *ast.TagNode:
		if nn.Comment != nil {
			return true
		}
		return hasComment(nn.Value)
	case *ast.MappingNode:
		if nn.Comment != nil || nn.FootComment != nil {
			return true
		}
		for _, value := range nn.Values {
			if value.Comment != nil || value.FootComment != nil {
				return true
			}
			if hasComment(value.Key) {
				return true
			}
			if hasComment(value.Value) {
				return true
			}
		}
	case *ast.MappingKeyNode:
		return nn.Comment != nil
	case *ast.MergeKeyNode:
		return nn.Comment != nil
	case *ast.SequenceNode:
		if nn.Comment != nil || nn.FootComment != nil {
			return true
		}
		for _, entry := range nn.Entries {
			if entry.Comment != nil || entry.HeadComment != nil || entry.LineComment != nil {
				return true
			}
			if hasComment(entry.Value) {
				return true
			}
		}
	case *ast.AnchorNode:
		if nn.Comment != nil {
			return true
		}
		if hasComment(nn.Name) || hasComment(nn.Value) {
			return true
		}
	case *ast.AliasNode:
		if nn.Comment != nil {
			return true
		}
		if hasComment(nn.Value) {
			return true
		}
	}
	return false
}

func getFirstToken(n ast.Node) *token.Token {
	if n == nil {
		return nil
	}
	switch nn := n.(type) {
	case *ast.DocumentNode:
		if nn.Start != nil {
			return nn.Start
		}
		return getFirstToken(nn.Body)
	case *ast.NullNode:
		return nn.Token
	case *ast.BoolNode:
		return nn.Token
	case *ast.IntegerNode:
		return nn.Token
	case *ast.FloatNode:
		return nn.Token
	case *ast.StringNode:
		return nn.Token
	case *ast.InfinityNode:
		return nn.Token
	case *ast.NanNode:
		return nn.Token
	case *ast.LiteralNode:
		return nn.Start
	case *ast.DirectiveNode:
		return nn.Start
	case *ast.TagNode:
		return nn.Start
	case *ast.MappingNode:
		if nn.IsFlowStyle {
			return nn.Start
		}
		if len(nn.Values) == 0 {
			return nn.Start
		}
		return getFirstToken(nn.Values[0].Key)
	case *ast.MappingKeyNode:
		return nn.Start
	case *ast.MergeKeyNode:
		return nn.Token
	case *ast.SequenceNode:
		return nn.Start
	case *ast.AnchorNode:
		return nn.Start
	case *ast.AliasNode:
		return nn.Start
	}
	return nil
}

type Formatter struct {
	existsComment    bool
	tokenToOriginMap map[*token.Token]string
	anchorNodeMap    map[string]ast.Node
}

func newFormatter(tk *token.Token, existsComment bool) *Formatter {
	tokenToOriginMap := make(map[*token.Token]string)
	for tk.Prev != nil {
		tk = tk.Prev
	}
	tokenToOriginMap[tk] = tk.Origin

	var origin string
	for tk.Next != nil {
		tk = tk.Next
		if tk.Type == token.CommentType {
			origin += strings.Repeat("\n", strings.Count(normalizeNewLineChars(tk.Origin), "\n"))
			continue
		}
		origin += tk.Origin
		tokenToOriginMap[tk] = origin
		origin = ""
	}
	return &Formatter{
		existsComment:    existsComment,
		tokenToOriginMap: tokenToOriginMap,
	}
}

func getIndentNumByFirstLineToken(tk *token.Token) int {
	defaultIndent := tk.Position.Column - 1

	// key: value
	//    ^
	//   next
	if tk.Type == token.SequenceEntryType {
		// If the current token is the sequence entry.
		// the indent is calculated from the column value of the current token.
		return defaultIndent
	}

	// key: value
	//    ^
	//   next
	if tk.Next != nil && tk.Next.Type == token.MappingValueType {
		// If the current token is the key in the mapping-value,
		// the indent is calculated from the column value of the current token.
		return defaultIndent
	}

	if tk.Prev == nil {
		return defaultIndent
	}
	prev := tk.Prev

	// key: value
	//    ^
	//   prev
	if prev.Type == token.MappingValueType {
		// If the current token is the value in the mapping-value,
		// the indent is calculated from the column value of the key two steps back.
		if prev.Prev == nil {
			return defaultIndent
		}
		return prev.Prev.Position.Column - 1
	}

	// - value
	// ^
	// prev
	if prev.Type == token.SequenceEntryType {
		// If the value is not a mapping-value and the previous token was a sequence entry,
		// the indent is calculated using the column value of the sequence entry token.
		return prev.Position.Column - 1
	}

	return defaultIndent
}

func (f *Formatter) format(n ast.Node) string {
	return f.trimSpacePrefix(
		f.trimIndentSpace(
			getIndentNumByFirstLineToken(getFirstToken(n)),
			f.trimNewLineCharPrefix(f.formatNode(n)),
		),
	)
}

func (f *Formatter) formatFile(file *ast.File) string {
	if len(file.Docs) == 0 {
		return ""
	}
	var ret string
	for _, doc := range file.Docs {
		ret += f.formatDocument(doc)
	}
	return ret
}

func (f *Formatter) origin(tk *token.Token) string {
	if tk == nil {
		return ""
	}
	if f.existsComment {
		return tk.Origin
	}
	return f.tokenToOriginMap[tk]
}

func (f *Formatter) formatDocument(n *ast.DocumentNode) string {
	return f.origin(n.Start) + f.formatNode(n.Body) + f.origin(n.End)
}

func (f *Formatter) formatNull(n *ast.NullNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatString(n *ast.StringNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatInteger(n *ast.IntegerNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatFloat(n *ast.FloatNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatBool(n *ast.BoolNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatInfinity(n *ast.InfinityNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatNan(n *ast.NanNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatLiteral(n *ast.LiteralNode) string {
	return f.origin(n.Start) + f.formatCommentGroup(n.Comment) + f.origin(n.Value.Token)
}

func (f *Formatter) formatMergeKey(n *ast.MergeKeyNode) string {
	return f.origin(n.Token)
}

func (f *Formatter) formatMappingValue(n *ast.MappingValueNode) string {
	return f.formatCommentGroup(n.Comment) +
		f.origin(n.Key.GetToken()) + ":" + f.formatCommentGroup(n.Key.GetComment()) + f.formatNode(n.Value) +
		f.formatCommentGroup(n.FootComment)
}

func (f *Formatter) formatDirective(n *ast.DirectiveNode) string {
	ret := f.origin(n.Start) + f.formatNode(n.Name)
	for _, val := range n.Values {
		ret += f.formatNode(val)
	}
	return ret
}

func (f *Formatter) formatMapping(n *ast.MappingNode) string {
	var ret string
	if n.IsFlowStyle {
		ret = f.origin(n.Start)
	}
	ret += f.formatCommentGroup(n.Comment)
	for _, value := range n.Values {
		if value.CollectEntry != nil {
			ret += f.origin(value.CollectEntry)
		}
		ret += f.formatMappingValue(value)
	}
	if n.IsFlowStyle {
		ret += f.origin(n.End)
	}
	return ret
}

func (f *Formatter) formatTag(n *ast.TagNode) string {
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatMappingKey(n *ast.MappingKeyNode) string {
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatSequence(n *ast.SequenceNode) string {
	var ret string
	if n.IsFlowStyle {
		ret = f.origin(n.Start)
	}
	if n.Comment != nil {
		// add head comment.
		ret += f.formatCommentGroup(n.Comment)
	}
	for _, entry := range n.Entries {
		ret += f.formatNode(entry)
	}
	if n.IsFlowStyle {
		ret += f.origin(n.End)
	}
	ret += f.formatCommentGroup(n.FootComment)
	return ret
}

func (f *Formatter) formatSequenceEntry(n *ast.SequenceEntryNode) string {
	return f.formatCommentGroup(n.HeadComment) + f.origin(n.Start) + f.formatCommentGroup(n.LineComment) + f.formatNode(n.Value)
}

func (f *Formatter) formatAnchor(n *ast.AnchorNode) string {
	return f.origin(n.Start) + f.formatNode(n.Name) + f.formatNode(n.Value)
}

func (f *Formatter) formatAlias(n *ast.AliasNode) string {
	if f.anchorNodeMap != nil {
		anchorName := n.Value.GetToken().Value
		node := f.anchorNodeMap[anchorName]
		if node != nil {
			formatted := f.formatNode(node)
			// If formatted text contains newline characters, indentation needs to be considered.
			if strings.Contains(formatted, "\n") {
				// If the first character is not a newline, the first line should be output without indentation.
				isIgnoredFirstLine := !strings.HasPrefix(formatted, "\n")
				formatted = f.addIndentSpace(n.GetToken().Position.IndentNum, formatted, isIgnoredFirstLine)
			}
			return formatted
		}
	}
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatNode(n ast.Node) string {
	switch nn := n.(type) {
	case *ast.DocumentNode:
		return f.formatDocument(nn)
	case *ast.NullNode:
		return f.formatNull(nn)
	case *ast.BoolNode:
		return f.formatBool(nn)
	case *ast.IntegerNode:
		return f.formatInteger(nn)
	case *ast.FloatNode:
		return f.formatFloat(nn)
	case *ast.StringNode:
		return f.formatString(nn)
	case *ast.InfinityNode:
		return f.formatInfinity(nn)
	case *ast.NanNode:
		return f.formatNan(nn)
	case *ast.LiteralNode:
		return f.formatLiteral(nn)
	case *ast.DirectiveNode:
		return f.formatDirective(nn)
	case *ast.TagNode:
		return f.formatTag(nn)
	case *ast.MappingNode:
		return f.formatMapping(nn)
	case *ast.MappingKeyNode:
		return f.formatMappingKey(nn)
	case *ast.MappingValueNode:
		return f.formatMappingValue(nn)
	case *ast.MergeKeyNode:
		return f.formatMergeKey(nn)
	case *ast.SequenceNode:
		return f.formatSequence(nn)
	case *ast.SequenceEntryNode:
		return f.formatSequenceEntry(nn)
	case *ast.AnchorNode:
		return f.formatAnchor(nn)
	case *ast.AliasNode:
		return f.formatAlias(nn)
	}
	return ""
}

func (f *Formatter) formatCommentGroup(g *ast.CommentGroupNode) string {
	if g == nil {
		return ""
	}
	var ret string
	for _, cm := range g.Comments {
		ret += f.formatComment(cm)
	}
	return ret
}

func (f *Formatter) formatComment(n *ast.CommentNode) string {
	if n == nil {
		return ""
	}
	return n.Token.Origin
}

// nolint: unused
func (f *Formatter) formatIndent(col int) string {
	if col <= 1 {
		return ""
	}
	return strings.Repeat(" ", col-1)
}

func (f *Formatter) trimNewLineCharPrefix(v string) string {
	return strings.TrimLeftFunc(v, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
}

func (f *Formatter) trimSpacePrefix(v string) string {
	return strings.TrimLeftFunc(v, func(r rune) bool {
		return r == ' '
	})
}

func (f *Formatter) trimIndentSpace(trimIndentNum int, v string) string {
	if trimIndentNum == 0 {
		return v
	}
	lines := strings.Split(normalizeNewLineChars(v), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		var cnt int
		out = append(out, strings.TrimLeftFunc(line, func(r rune) bool {
			cnt++
			return r == ' ' && cnt <= trimIndentNum
		}))
	}
	return strings.Join(out, "\n")
}

func (f *Formatter) addIndentSpace(indentNum int, v string, isIgnoredFirstLine bool) string {
	if indentNum == 0 {
		return v
	}
	indent := strings.Repeat(" ", indentNum)
	lines := strings.Split(normalizeNewLineChars(v), "\n")
	out := make([]string, 0, len(lines))
	for idx, line := range lines {
		if line == "" || (isIgnoredFirstLine && idx == 0) {
			out = append(out, line)
			continue
		}
		out = append(out, indent+line)
	}
	return strings.Join(out, "\n")
}

// normalizeNewLineChars normalize CRLF and CR to LF.
func normalizeNewLineChars(v string) string {
	return strings.ReplaceAll(strings.ReplaceAll(v, "\r\n", "\n"), "\r", "\n")
}
