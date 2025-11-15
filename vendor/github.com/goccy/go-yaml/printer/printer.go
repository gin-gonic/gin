package printer

import (
	"fmt"
	"math"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

// Property additional property set for each the token
type Property struct {
	Prefix string
	Suffix string
}

// PrintFunc returns property instance
type PrintFunc func() *Property

// Printer create text from token collection or ast
type Printer struct {
	LineNumber       bool
	LineNumberFormat func(num int) string
	MapKey           PrintFunc
	Anchor           PrintFunc
	Alias            PrintFunc
	Bool             PrintFunc
	String           PrintFunc
	Number           PrintFunc
	Comment          PrintFunc
}

func defaultLineNumberFormat(num int) string {
	return fmt.Sprintf("%2d | ", num)
}

func (p *Printer) property(tk *token.Token) *Property {
	prop := &Property{}
	switch tk.PreviousType() {
	case token.AnchorType:
		if p.Anchor != nil {
			return p.Anchor()
		}
		return prop
	case token.AliasType:
		if p.Alias != nil {
			return p.Alias()
		}
		return prop
	}
	switch tk.NextType() {
	case token.MappingValueType:
		if p.MapKey != nil {
			return p.MapKey()
		}
		return prop
	}
	switch tk.Type {
	case token.BoolType:
		if p.Bool != nil {
			return p.Bool()
		}
		return prop
	case token.AnchorType:
		if p.Anchor != nil {
			return p.Anchor()
		}
		return prop
	case token.AliasType:
		if p.Anchor != nil {
			return p.Alias()
		}
		return prop
	case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
		if p.String != nil {
			return p.String()
		}
		return prop
	case token.IntegerType, token.FloatType:
		if p.Number != nil {
			return p.Number()
		}
		return prop
	case token.CommentType:
		if p.Comment != nil {
			return p.Comment()
		}
		return prop
	default:
	}
	return prop
}

// PrintTokens create text from token collection
func (p *Printer) PrintTokens(tokens token.Tokens) string {
	if len(tokens) == 0 {
		return ""
	}
	if p.LineNumber {
		if p.LineNumberFormat == nil {
			p.LineNumberFormat = defaultLineNumberFormat
		}
	}
	texts := []string{}
	lineNumber := tokens[0].Position.Line
	for _, tk := range tokens {
		lines := strings.Split(tk.Origin, "\n")
		prop := p.property(tk)
		header := ""
		if p.LineNumber {
			header = p.LineNumberFormat(lineNumber)
		}
		if len(lines) == 1 {
			line := prop.Prefix + lines[0] + prop.Suffix
			if len(texts) == 0 {
				texts = append(texts, header+line)
				lineNumber++
			} else {
				text := texts[len(texts)-1]
				texts[len(texts)-1] = text + line
			}
		} else {
			for idx, src := range lines {
				if p.LineNumber {
					header = p.LineNumberFormat(lineNumber)
				}
				line := prop.Prefix + src + prop.Suffix
				if idx == 0 {
					if len(texts) == 0 {
						texts = append(texts, header+line)
						lineNumber++
					} else {
						text := texts[len(texts)-1]
						texts[len(texts)-1] = text + line
					}
				} else {
					texts = append(texts, fmt.Sprintf("%s%s", header, line))
					lineNumber++
				}
			}
		}
	}
	return strings.Join(texts, "\n")
}

// PrintNode create text from ast.Node
func (p *Printer) PrintNode(node ast.Node) []byte {
	return []byte(fmt.Sprintf("%+v\n", node))
}

func (p *Printer) setDefaultColorSet() {
	p.Bool = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiMagenta),
			Suffix: format(ColorReset),
		}
	}
	p.Number = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiMagenta),
			Suffix: format(ColorReset),
		}
	}
	p.MapKey = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiCyan),
			Suffix: format(ColorReset),
		}
	}
	p.Anchor = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiYellow),
			Suffix: format(ColorReset),
		}
	}
	p.Alias = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiYellow),
			Suffix: format(ColorReset),
		}
	}
	p.String = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiGreen),
			Suffix: format(ColorReset),
		}
	}
	p.Comment = func() *Property {
		return &Property{
			Prefix: format(ColorFgHiBlack),
			Suffix: format(ColorReset),
		}
	}
}

func (p *Printer) PrintErrorMessage(msg string, isColored bool) string {
	if isColored {
		return fmt.Sprintf("%s%s%s",
			format(ColorFgHiRed),
			msg,
			format(ColorReset),
		)
	}
	return msg
}

func (p *Printer) removeLeftSideNewLineChar(src string) string {
	return strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(src, "\r"), "\n"), "\r\n")
}

func (p *Printer) removeRightSideNewLineChar(src string) string {
	return strings.TrimRight(strings.TrimRight(strings.TrimRight(src, "\r"), "\n"), "\r\n")
}

func (p *Printer) removeRightSideWhiteSpaceChar(src string) string {
	return p.removeRightSideNewLineChar(strings.TrimRight(src, " "))
}

func (p *Printer) newLineCount(s string) int {
	src := []rune(s)
	size := len(src)
	cnt := 0
	for i := 0; i < size; i++ {
		c := src[i]
		switch c {
		case '\r':
			if i+1 < size && src[i+1] == '\n' {
				i++
			}
			cnt++
		case '\n':
			cnt++
		}
	}
	return cnt
}

func (p *Printer) isNewLineLastChar(s string) bool {
	for i := len(s) - 1; i > 0; i-- {
		c := s[i]
		switch c {
		case ' ':
			continue
		case '\n', '\r':
			return true
		}
		break
	}
	return false
}

func (p *Printer) printBeforeTokens(tk *token.Token, minLine, extLine int) token.Tokens {
	for tk.Prev != nil {
		if tk.Prev.Position.Line < minLine {
			break
		}
		tk = tk.Prev
	}
	minTk := tk.Clone()
	if minTk.Prev != nil {
		// add white spaces to minTk by prev token
		prev := minTk.Prev
		whiteSpaceLen := len(prev.Origin) - len(strings.TrimRight(prev.Origin, " "))
		minTk.Origin = strings.Repeat(" ", whiteSpaceLen) + minTk.Origin
	}
	minTk.Origin = p.removeLeftSideNewLineChar(minTk.Origin)
	tokens := token.Tokens{minTk}
	tk = minTk.Next
	for tk != nil && tk.Position.Line <= extLine {
		clonedTk := tk.Clone()
		tokens.Add(clonedTk)
		tk = clonedTk.Next
	}
	lastTk := tokens[len(tokens)-1]
	trimmedOrigin := p.removeRightSideWhiteSpaceChar(lastTk.Origin)
	suffix := lastTk.Origin[len(trimmedOrigin):]
	lastTk.Origin = trimmedOrigin

	if lastTk.Next != nil && len(suffix) > 1 {
		next := lastTk.Next.Clone()
		// add suffix to header of next token
		if suffix[0] == '\n' || suffix[0] == '\r' {
			suffix = suffix[1:]
		}
		next.Origin = suffix + next.Origin
		lastTk.Next = next
	}
	return tokens
}

func (p *Printer) printAfterTokens(tk *token.Token, maxLine int) token.Tokens {
	tokens := token.Tokens{}
	if tk == nil {
		return tokens
	}
	if tk.Position.Line > maxLine {
		return tokens
	}
	minTk := tk.Clone()
	minTk.Origin = p.removeLeftSideNewLineChar(minTk.Origin)
	tokens.Add(minTk)
	tk = minTk.Next
	for tk != nil && tk.Position.Line <= maxLine {
		clonedTk := tk.Clone()
		tokens.Add(clonedTk)
		tk = clonedTk.Next
	}
	return tokens
}

func (p *Printer) setupErrorTokenFormat(annotateLine int, isColored bool) {
	prefix := func(annotateLine, num int) string {
		if annotateLine == num {
			return fmt.Sprintf("> %2d | ", num)
		}
		return fmt.Sprintf("  %2d | ", num)
	}
	p.LineNumber = true
	p.LineNumberFormat = func(num int) string {
		if isColored {
			return colorize(prefix(annotateLine, num), ColorBold, ColorFgHiWhite)
		}
		return prefix(annotateLine, num)
	}
	if isColored {
		p.setDefaultColorSet()
	}
}

func (p *Printer) PrintErrorToken(tk *token.Token, isColored bool) string {
	errToken := tk
	curLine := tk.Position.Line
	curExtLine := curLine + p.newLineCount(p.removeLeftSideNewLineChar(tk.Origin))
	if p.isNewLineLastChar(tk.Origin) {
		// if last character ( exclude white space ) is new line character, ignore it.
		curExtLine--
	}

	minLine := int(math.Max(float64(curLine-3), 1))
	maxLine := curExtLine + 3
	p.setupErrorTokenFormat(curLine, isColored)

	beforeTokens := p.printBeforeTokens(tk, minLine, curExtLine)
	lastTk := beforeTokens[len(beforeTokens)-1]
	afterTokens := p.printAfterTokens(lastTk.Next, maxLine)

	beforeSource := p.PrintTokens(beforeTokens)
	prefixSpaceNum := len(fmt.Sprintf("  %2d | ", curLine))
	annotateLine := strings.Repeat(" ", prefixSpaceNum+errToken.Position.Column-1) + "^"
	afterSource := p.PrintTokens(afterTokens)
	return fmt.Sprintf("%s\n%s\n%s", beforeSource, annotateLine, afterSource)
}
