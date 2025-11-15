package token

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Character type for character
type Character byte

const (
	// SequenceEntryCharacter character for sequence entry
	SequenceEntryCharacter Character = '-'
	// MappingKeyCharacter character for mapping key
	MappingKeyCharacter Character = '?'
	// MappingValueCharacter character for mapping value
	MappingValueCharacter Character = ':'
	// CollectEntryCharacter character for collect entry
	CollectEntryCharacter Character = ','
	// SequenceStartCharacter character for sequence start
	SequenceStartCharacter Character = '['
	// SequenceEndCharacter character for sequence end
	SequenceEndCharacter Character = ']'
	// MappingStartCharacter character for mapping start
	MappingStartCharacter Character = '{'
	// MappingEndCharacter character for mapping end
	MappingEndCharacter Character = '}'
	// CommentCharacter character for comment
	CommentCharacter Character = '#'
	// AnchorCharacter character for anchor
	AnchorCharacter Character = '&'
	// AliasCharacter character for alias
	AliasCharacter Character = '*'
	// TagCharacter character for tag
	TagCharacter Character = '!'
	// LiteralCharacter character for literal
	LiteralCharacter Character = '|'
	// FoldedCharacter character for folded
	FoldedCharacter Character = '>'
	// SingleQuoteCharacter character for single quote
	SingleQuoteCharacter Character = '\''
	// DoubleQuoteCharacter character for double quote
	DoubleQuoteCharacter Character = '"'
	// DirectiveCharacter character for directive
	DirectiveCharacter Character = '%'
	// SpaceCharacter character for space
	SpaceCharacter Character = ' '
	// LineBreakCharacter character for line break
	LineBreakCharacter Character = '\n'
)

// Type type identifier for token
type Type int

const (
	// UnknownType reserve for invalid type
	UnknownType Type = iota
	// DocumentHeaderType type for DocumentHeader token
	DocumentHeaderType
	// DocumentEndType type for DocumentEnd token
	DocumentEndType
	// SequenceEntryType type for SequenceEntry token
	SequenceEntryType
	// MappingKeyType type for MappingKey token
	MappingKeyType
	// MappingValueType type for MappingValue token
	MappingValueType
	// MergeKeyType type for MergeKey token
	MergeKeyType
	// CollectEntryType type for CollectEntry token
	CollectEntryType
	// SequenceStartType type for SequenceStart token
	SequenceStartType
	// SequenceEndType type for SequenceEnd token
	SequenceEndType
	// MappingStartType type for MappingStart token
	MappingStartType
	// MappingEndType type for MappingEnd token
	MappingEndType
	// CommentType type for Comment token
	CommentType
	// AnchorType type for Anchor token
	AnchorType
	// AliasType type for Alias token
	AliasType
	// TagType type for Tag token
	TagType
	// LiteralType type for Literal token
	LiteralType
	// FoldedType type for Folded token
	FoldedType
	// SingleQuoteType type for SingleQuote token
	SingleQuoteType
	// DoubleQuoteType type for DoubleQuote token
	DoubleQuoteType
	// DirectiveType type for Directive token
	DirectiveType
	// SpaceType type for Space token
	SpaceType
	// NullType type for Null token
	NullType
	// ImplicitNullType type for implicit Null token.
	// This is used when explicit keywords such as null or ~ are not specified.
	// It is distinguished during encoding and output as an empty string.
	ImplicitNullType
	// InfinityType type for Infinity token
	InfinityType
	// NanType type for Nan token
	NanType
	// IntegerType type for Integer token
	IntegerType
	// BinaryIntegerType type for BinaryInteger token
	BinaryIntegerType
	// OctetIntegerType type for OctetInteger token
	OctetIntegerType
	// HexIntegerType type for HexInteger token
	HexIntegerType
	// FloatType type for Float token
	FloatType
	// StringType type for String token
	StringType
	// BoolType type for Bool token
	BoolType
	// InvalidType type for invalid token
	InvalidType
)

// String type identifier to text
func (t Type) String() string {
	switch t {
	case UnknownType:
		return "Unknown"
	case DocumentHeaderType:
		return "DocumentHeader"
	case DocumentEndType:
		return "DocumentEnd"
	case SequenceEntryType:
		return "SequenceEntry"
	case MappingKeyType:
		return "MappingKey"
	case MappingValueType:
		return "MappingValue"
	case MergeKeyType:
		return "MergeKey"
	case CollectEntryType:
		return "CollectEntry"
	case SequenceStartType:
		return "SequenceStart"
	case SequenceEndType:
		return "SequenceEnd"
	case MappingStartType:
		return "MappingStart"
	case MappingEndType:
		return "MappingEnd"
	case CommentType:
		return "Comment"
	case AnchorType:
		return "Anchor"
	case AliasType:
		return "Alias"
	case TagType:
		return "Tag"
	case LiteralType:
		return "Literal"
	case FoldedType:
		return "Folded"
	case SingleQuoteType:
		return "SingleQuote"
	case DoubleQuoteType:
		return "DoubleQuote"
	case DirectiveType:
		return "Directive"
	case SpaceType:
		return "Space"
	case StringType:
		return "String"
	case BoolType:
		return "Bool"
	case IntegerType:
		return "Integer"
	case BinaryIntegerType:
		return "BinaryInteger"
	case OctetIntegerType:
		return "OctetInteger"
	case HexIntegerType:
		return "HexInteger"
	case FloatType:
		return "Float"
	case NullType:
		return "Null"
	case ImplicitNullType:
		return "ImplicitNull"
	case InfinityType:
		return "Infinity"
	case NanType:
		return "Nan"
	case InvalidType:
		return "Invalid"
	}
	return ""
}

// CharacterType type for character category
type CharacterType int

const (
	// CharacterTypeIndicator type of indicator character
	CharacterTypeIndicator CharacterType = iota
	// CharacterTypeWhiteSpace type of white space character
	CharacterTypeWhiteSpace
	// CharacterTypeMiscellaneous type of miscellaneous character
	CharacterTypeMiscellaneous
	// CharacterTypeEscaped type of escaped character
	CharacterTypeEscaped
	// CharacterTypeInvalid type for a invalid token.
	CharacterTypeInvalid
)

// String character type identifier to text
func (c CharacterType) String() string {
	switch c {
	case CharacterTypeIndicator:
		return "Indicator"
	case CharacterTypeWhiteSpace:
		return "WhiteSpace"
	case CharacterTypeMiscellaneous:
		return "Miscellaneous"
	case CharacterTypeEscaped:
		return "Escaped"
	}
	return ""
}

// Indicator type for indicator
type Indicator int

const (
	// NotIndicator not indicator
	NotIndicator Indicator = iota
	// BlockStructureIndicator indicator for block structure ( '-', '?', ':' )
	BlockStructureIndicator
	// FlowCollectionIndicator indicator for flow collection ( '[', ']', '{', '}', ',' )
	FlowCollectionIndicator
	// CommentIndicator indicator for comment ( '#' )
	CommentIndicator
	// NodePropertyIndicator indicator for node property ( '!', '&', '*' )
	NodePropertyIndicator
	// BlockScalarIndicator indicator for block scalar ( '|', '>' )
	BlockScalarIndicator
	// QuotedScalarIndicator indicator for quoted scalar ( ''', '"' )
	QuotedScalarIndicator
	// DirectiveIndicator indicator for directive ( '%' )
	DirectiveIndicator
	// InvalidUseOfReservedIndicator indicator for invalid use of reserved keyword ( '@', '`' )
	InvalidUseOfReservedIndicator
)

// String indicator to text
func (i Indicator) String() string {
	switch i {
	case NotIndicator:
		return "NotIndicator"
	case BlockStructureIndicator:
		return "BlockStructure"
	case FlowCollectionIndicator:
		return "FlowCollection"
	case CommentIndicator:
		return "Comment"
	case NodePropertyIndicator:
		return "NodeProperty"
	case BlockScalarIndicator:
		return "BlockScalar"
	case QuotedScalarIndicator:
		return "QuotedScalar"
	case DirectiveIndicator:
		return "Directive"
	case InvalidUseOfReservedIndicator:
		return "InvalidUseOfReserved"
	}
	return ""
}

var (
	reservedNullKeywords = []string{
		"null",
		"Null",
		"NULL",
		"~",
	}
	reservedBoolKeywords = []string{
		"true",
		"True",
		"TRUE",
		"false",
		"False",
		"FALSE",
	}
	// For compatibility with other YAML 1.1 parsers
	// Note that we use these solely for encoding the bool value with quotes.
	// go-yaml should not treat these as reserved keywords at parsing time.
	// as go-yaml is supposed to be compliant only with YAML 1.2.
	reservedLegacyBoolKeywords = []string{
		"y",
		"Y",
		"yes",
		"Yes",
		"YES",
		"n",
		"N",
		"no",
		"No",
		"NO",
		"on",
		"On",
		"ON",
		"off",
		"Off",
		"OFF",
	}
	reservedInfKeywords = []string{
		".inf",
		".Inf",
		".INF",
		"-.inf",
		"-.Inf",
		"-.INF",
	}
	reservedNanKeywords = []string{
		".nan",
		".NaN",
		".NAN",
	}
	reservedKeywordMap = map[string]func(string, string, *Position) *Token{}
	// reservedEncKeywordMap contains is the keyword map used at encoding time.
	// This is supposed to be a superset of reservedKeywordMap,
	// and used to quote legacy keywords present in YAML 1.1 or lesser for compatibility reasons,
	// even though this library is supposed to be YAML 1.2-compliant.
	reservedEncKeywordMap = map[string]func(string, string, *Position) *Token{}
)

func reservedKeywordToken(typ Type, value, org string, pos *Position) *Token {
	return &Token{
		Type:          typ,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

func init() {
	for _, keyword := range reservedNullKeywords {
		f := func(value, org string, pos *Position) *Token {
			return reservedKeywordToken(NullType, value, org, pos)
		}

		reservedKeywordMap[keyword] = f
		reservedEncKeywordMap[keyword] = f
	}
	for _, keyword := range reservedBoolKeywords {
		f := func(value, org string, pos *Position) *Token {
			return reservedKeywordToken(BoolType, value, org, pos)
		}
		reservedKeywordMap[keyword] = f
		reservedEncKeywordMap[keyword] = f
	}
	for _, keyword := range reservedLegacyBoolKeywords {
		reservedEncKeywordMap[keyword] = func(value, org string, pos *Position) *Token {
			return reservedKeywordToken(BoolType, value, org, pos)
		}
	}
	for _, keyword := range reservedInfKeywords {
		reservedKeywordMap[keyword] = func(value, org string, pos *Position) *Token {
			return reservedKeywordToken(InfinityType, value, org, pos)
		}
	}
	for _, keyword := range reservedNanKeywords {
		reservedKeywordMap[keyword] = func(value, org string, pos *Position) *Token {
			return reservedKeywordToken(NanType, value, org, pos)
		}
	}
}

// ReservedTagKeyword type of reserved tag keyword
type ReservedTagKeyword string

const (
	// IntegerTag `!!int` tag
	IntegerTag ReservedTagKeyword = "!!int"
	// FloatTag `!!float` tag
	FloatTag ReservedTagKeyword = "!!float"
	// NullTag `!!null` tag
	NullTag ReservedTagKeyword = "!!null"
	// SequenceTag `!!seq` tag
	SequenceTag ReservedTagKeyword = "!!seq"
	// MappingTag `!!map` tag
	MappingTag ReservedTagKeyword = "!!map"
	// StringTag `!!str` tag
	StringTag ReservedTagKeyword = "!!str"
	// BinaryTag `!!binary` tag
	BinaryTag ReservedTagKeyword = "!!binary"
	// OrderedMapTag `!!omap` tag
	OrderedMapTag ReservedTagKeyword = "!!omap"
	// SetTag `!!set` tag
	SetTag ReservedTagKeyword = "!!set"
	// TimestampTag `!!timestamp` tag
	TimestampTag ReservedTagKeyword = "!!timestamp"
	// BooleanTag `!!bool` tag
	BooleanTag ReservedTagKeyword = "!!bool"
	// MergeTag `!!merge` tag
	MergeTag ReservedTagKeyword = "!!merge"
)

var (
	// ReservedTagKeywordMap map for reserved tag keywords
	ReservedTagKeywordMap = map[ReservedTagKeyword]func(string, string, *Position) *Token{
		IntegerTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		FloatTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		NullTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		SequenceTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		MappingTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		StringTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		BinaryTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		OrderedMapTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		SetTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		TimestampTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		BooleanTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		MergeTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
	}
)

type NumberType string

const (
	NumberTypeDecimal NumberType = "decimal"
	NumberTypeBinary  NumberType = "binary"
	NumberTypeOctet   NumberType = "octet"
	NumberTypeHex     NumberType = "hex"
	NumberTypeFloat   NumberType = "float"
)

type NumberValue struct {
	Type  NumberType
	Value any
	Text  string
}

func ToNumber(value string) *NumberValue {
	num, err := toNumber(value)
	if err != nil {
		return nil
	}
	return num
}

func isNumber(value string) bool {
	num, err := toNumber(value)
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) && errors.Is(numErr.Err, strconv.ErrRange) {
			return true
		}
		return false
	}
	return num != nil
}

func toNumber(value string) (*NumberValue, error) {
	if len(value) == 0 {
		return nil, nil
	}
	if strings.HasPrefix(value, "_") {
		return nil, nil
	}
	dotCount := strings.Count(value, ".")
	if dotCount > 1 {
		return nil, nil
	}

	isNegative := strings.HasPrefix(value, "-")
	normalized := strings.ReplaceAll(strings.TrimPrefix(strings.TrimPrefix(value, "+"), "-"), "_", "")

	var (
		typ  NumberType
		base int
	)
	switch {
	case strings.HasPrefix(normalized, "0x"):
		normalized = strings.TrimPrefix(normalized, "0x")
		base = 16
		typ = NumberTypeHex
	case strings.HasPrefix(normalized, "0o"):
		normalized = strings.TrimPrefix(normalized, "0o")
		base = 8
		typ = NumberTypeOctet
	case strings.HasPrefix(normalized, "0b"):
		normalized = strings.TrimPrefix(normalized, "0b")
		base = 2
		typ = NumberTypeBinary
	case strings.HasPrefix(normalized, "0") && len(normalized) > 1 && dotCount == 0:
		base = 8
		typ = NumberTypeOctet
	case dotCount == 1:
		typ = NumberTypeFloat
	default:
		typ = NumberTypeDecimal
		base = 10
	}

	text := normalized
	if isNegative {
		text = "-" + text
	}

	var v any
	if typ == NumberTypeFloat {
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, err
		}
		v = f
	} else if isNegative {
		i, err := strconv.ParseInt(text, base, 64)
		if err != nil {
			return nil, err
		}
		v = i
	} else {
		u, err := strconv.ParseUint(text, base, 64)
		if err != nil {
			return nil, err
		}
		v = u
	}

	return &NumberValue{
		Type:  typ,
		Value: v,
		Text:  text,
	}, nil
}

// This is a subset of the formats permitted by the regular expression
// defined at http://yaml.org/type/timestamp.html. Note that time.Parse
// cannot handle: "2001-12-14 21:59:43.10 -5" from the examples.
var timestampFormats = []string{
	time.RFC3339Nano,
	"2006-01-02t15:04:05.999999999Z07:00", // RFC3339Nano with lower-case "t".
	time.DateTime,
	time.DateOnly,

	// Not in examples, but to preserve backward compatibility by quoting time values.
	"15:4",
}

func isTimestamp(value string) bool {
	for _, format := range timestampFormats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

// IsNeedQuoted checks whether the value needs quote for passed string or not
func IsNeedQuoted(value string) bool {
	if value == "" {
		return true
	}
	if _, exists := reservedEncKeywordMap[value]; exists {
		return true
	}
	if isNumber(value) {
		return true
	}
	if value == "-" {
		return true
	}
	first := value[0]
	switch first {
	case '*', '&', '[', '{', '}', ']', ',', '!', '|', '>', '%', '\'', '"', '@', ' ', '`':
		return true
	}
	last := value[len(value)-1]
	switch last {
	case ':', ' ':
		return true
	}
	if isTimestamp(value) {
		return true
	}
	for i, c := range value {
		switch c {
		case '#', '\\':
			return true
		case ':', '-':
			if i+1 < len(value) && value[i+1] == ' ' {
				return true
			}
		}
	}
	return false
}

// LiteralBlockHeader detect literal block scalar header
func LiteralBlockHeader(value string) string {
	lbc := DetectLineBreakCharacter(value)

	switch {
	case !strings.Contains(value, lbc):
		return ""
	case strings.HasSuffix(value, fmt.Sprintf("%s%s", lbc, lbc)):
		return "|+"
	case strings.HasSuffix(value, lbc):
		return "|"
	default:
		return "|-"
	}
}

// New create reserved keyword token or number token and other string token.
func New(value string, org string, pos *Position) *Token {
	fn := reservedKeywordMap[value]
	if fn != nil {
		return fn(value, org, pos)
	}
	if num := ToNumber(value); num != nil {
		tk := &Token{
			Type:          IntegerType,
			CharacterType: CharacterTypeMiscellaneous,
			Indicator:     NotIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
		switch num.Type {
		case NumberTypeFloat:
			tk.Type = FloatType
		case NumberTypeBinary:
			tk.Type = BinaryIntegerType
		case NumberTypeOctet:
			tk.Type = OctetIntegerType
		case NumberTypeHex:
			tk.Type = HexIntegerType
		}
		return tk
	}
	return String(value, org, pos)
}

// Position type for position in YAML document
type Position struct {
	Line        int
	Column      int
	Offset      int
	IndentNum   int
	IndentLevel int
}

// String position to text
func (p *Position) String() string {
	return fmt.Sprintf("[level:%d,line:%d,column:%d,offset:%d]", p.IndentLevel, p.Line, p.Column, p.Offset)
}

// Token type for token
type Token struct {
	// Type is a token type.
	Type Type
	// CharacterType is a character type.
	CharacterType CharacterType
	// Indicator is a indicator type.
	Indicator Indicator
	// Value is a string extracted with only meaningful characters, with spaces and such removed.
	Value string
	// Origin is a string that stores the original text as-is.
	Origin string
	// Error keeps error message for InvalidToken.
	Error string
	// Position is a token position.
	Position *Position
	// Next is a next token reference.
	Next *Token
	// Prev is a previous token reference.
	Prev *Token
}

// PreviousType previous token type
func (t *Token) PreviousType() Type {
	if t.Prev != nil {
		return t.Prev.Type
	}
	return UnknownType
}

// NextType next token type
func (t *Token) NextType() Type {
	if t.Next != nil {
		return t.Next.Type
	}
	return UnknownType
}

// AddColumn append column number to current position of column
func (t *Token) AddColumn(col int) {
	if t == nil {
		return
	}
	t.Position.Column += col
}

// Clone copy token ( preserve Prev/Next reference )
func (t *Token) Clone() *Token {
	if t == nil {
		return nil
	}
	copied := *t
	if t.Position != nil {
		pos := *(t.Position)
		copied.Position = &pos
	}
	return &copied
}

// Dump outputs token information to stdout for debugging.
func (t *Token) Dump() {
	fmt.Printf(
		"[TYPE]:%q [CHARTYPE]:%q [INDICATOR]:%q [VALUE]:%q [ORG]:%q [POS(line:column:level:offset)]: %d:%d:%d:%d\n",
		t.Type, t.CharacterType, t.Indicator, t.Value, t.Origin, t.Position.Line, t.Position.Column, t.Position.IndentLevel, t.Position.Offset,
	)
}

// Tokens type of token collection
type Tokens []*Token

func (t Tokens) InvalidToken() *Token {
	for _, tt := range t {
		if tt.Type == InvalidType {
			return tt
		}
	}
	return nil
}

func (t *Tokens) add(tk *Token) {
	tokens := *t
	if len(tokens) == 0 {
		tokens = append(tokens, tk)
	} else {
		last := tokens[len(tokens)-1]
		last.Next = tk
		tk.Prev = last
		tokens = append(tokens, tk)
	}
	*t = tokens
}

// Add append new some tokens
func (t *Tokens) Add(tks ...*Token) {
	for _, tk := range tks {
		t.add(tk)
	}
}

// Dump dump all token structures for debugging
func (t Tokens) Dump() {
	for _, tk := range t {
		fmt.Print("- ")
		tk.Dump()
	}
}

// String create token for String
func String(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          StringType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// SequenceEntry create token for SequenceEntry
func SequenceEntry(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceEntryType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(SequenceEntryCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingKey create token for MappingKey
func MappingKey(pos *Position) *Token {
	return &Token{
		Type:          MappingKeyType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(MappingKeyCharacter),
		Origin:        string(MappingKeyCharacter),
		Position:      pos,
	}
}

// MappingValue create token for MappingValue
func MappingValue(pos *Position) *Token {
	return &Token{
		Type:          MappingValueType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(MappingValueCharacter),
		Origin:        string(MappingValueCharacter),
		Position:      pos,
	}
}

// CollectEntry create token for CollectEntry
func CollectEntry(org string, pos *Position) *Token {
	return &Token{
		Type:          CollectEntryType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(CollectEntryCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// SequenceStart create token for SequenceStart
func SequenceStart(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceStartType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(SequenceStartCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// SequenceEnd create token for SequenceEnd
func SequenceEnd(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceEndType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(SequenceEndCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingStart create token for MappingStart
func MappingStart(org string, pos *Position) *Token {
	return &Token{
		Type:          MappingStartType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(MappingStartCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingEnd create token for MappingEnd
func MappingEnd(org string, pos *Position) *Token {
	return &Token{
		Type:          MappingEndType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(MappingEndCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Comment create token for Comment
func Comment(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          CommentType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     CommentIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Anchor create token for Anchor
func Anchor(org string, pos *Position) *Token {
	return &Token{
		Type:          AnchorType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         string(AnchorCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Alias create token for Alias
func Alias(org string, pos *Position) *Token {
	return &Token{
		Type:          AliasType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         string(AliasCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Tag create token for Tag
func Tag(value string, org string, pos *Position) *Token {
	fn := ReservedTagKeywordMap[ReservedTagKeyword(value)]
	if fn != nil {
		return fn(value, org, pos)
	}
	return &Token{
		Type:          TagType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Literal create token for Literal
func Literal(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          LiteralType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Folded create token for Folded
func Folded(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          FoldedType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// SingleQuote create token for SingleQuote
func SingleQuote(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          SingleQuoteType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     QuotedScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// DoubleQuote create token for DoubleQuote
func DoubleQuote(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          DoubleQuoteType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     QuotedScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Directive create token for Directive
func Directive(org string, pos *Position) *Token {
	return &Token{
		Type:          DirectiveType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     DirectiveIndicator,
		Value:         string(DirectiveCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Space create token for Space
func Space(pos *Position) *Token {
	return &Token{
		Type:          SpaceType,
		CharacterType: CharacterTypeWhiteSpace,
		Indicator:     NotIndicator,
		Value:         string(SpaceCharacter),
		Origin:        string(SpaceCharacter),
		Position:      pos,
	}
}

// MergeKey create token for MergeKey
func MergeKey(org string, pos *Position) *Token {
	return &Token{
		Type:          MergeKeyType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "<<",
		Origin:        org,
		Position:      pos,
	}
}

// DocumentHeader create token for DocumentHeader
func DocumentHeader(org string, pos *Position) *Token {
	return &Token{
		Type:          DocumentHeaderType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "---",
		Origin:        org,
		Position:      pos,
	}
}

// DocumentEnd create token for DocumentEnd
func DocumentEnd(org string, pos *Position) *Token {
	return &Token{
		Type:          DocumentEndType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "...",
		Origin:        org,
		Position:      pos,
	}
}

func Invalid(err string, org string, pos *Position) *Token {
	return &Token{
		Type:          InvalidType,
		CharacterType: CharacterTypeInvalid,
		Indicator:     NotIndicator,
		Value:         org,
		Origin:        org,
		Error:         err,
		Position:      pos,
	}
}

// DetectLineBreakCharacter detect line break character in only one inside scalar content scope.
func DetectLineBreakCharacter(src string) string {
	nc := strings.Count(src, "\n")
	rc := strings.Count(src, "\r")
	rnc := strings.Count(src, "\r\n")
	switch {
	case nc == rnc && rc == rnc:
		return "\r\n"
	case rc > nc:
		return "\r"
	default:
		return "\n"
	}
}
