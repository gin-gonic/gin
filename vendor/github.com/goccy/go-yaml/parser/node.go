package parser

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/token"
)

func newMappingNode(ctx *context, tk *Token, isFlow bool, values ...*ast.MappingValueNode) (*ast.MappingNode, error) {
	node := ast.Mapping(tk.RawToken(), isFlow, values...)
	node.SetPath(ctx.path)
	return node, nil
}

func newMappingValueNode(ctx *context, colonTk, entryTk *Token, key ast.MapKeyNode, value ast.Node) (*ast.MappingValueNode, error) {
	node := ast.MappingValue(colonTk.RawToken(), key, value)
	node.SetPath(ctx.path)
	node.CollectEntry = entryTk.RawToken()
	if key.GetToken().Position.Line == value.GetToken().Position.Line {
		// originally key was commented, but now that null value has been added, value must be commented.
		if err := setLineComment(ctx, value, colonTk); err != nil {
			return nil, err
		}
		// set line comment by colonTk or entryTk.
		if err := setLineComment(ctx, value, entryTk); err != nil {
			return nil, err
		}
	} else {
		if err := setLineComment(ctx, key, colonTk); err != nil {
			return nil, err
		}
		// set line comment by colonTk or entryTk.
		if err := setLineComment(ctx, key, entryTk); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func newMappingKeyNode(ctx *context, tk *Token) (*ast.MappingKeyNode, error) {
	node := ast.MappingKey(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newAnchorNode(ctx *context, tk *Token) (*ast.AnchorNode, error) {
	node := ast.Anchor(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newAliasNode(ctx *context, tk *Token) (*ast.AliasNode, error) {
	node := ast.Alias(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newDirectiveNode(ctx *context, tk *Token) (*ast.DirectiveNode, error) {
	node := ast.Directive(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newMergeKeyNode(ctx *context, tk *Token) (*ast.MergeKeyNode, error) {
	node := ast.MergeKey(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newNullNode(ctx *context, tk *Token) (*ast.NullNode, error) {
	node := ast.Null(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newBoolNode(ctx *context, tk *Token) (*ast.BoolNode, error) {
	node := ast.Bool(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newIntegerNode(ctx *context, tk *Token) (*ast.IntegerNode, error) {
	node := ast.Integer(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newFloatNode(ctx *context, tk *Token) (*ast.FloatNode, error) {
	node := ast.Float(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newInfinityNode(ctx *context, tk *Token) (*ast.InfinityNode, error) {
	node := ast.Infinity(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newNanNode(ctx *context, tk *Token) (*ast.NanNode, error) {
	node := ast.Nan(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newStringNode(ctx *context, tk *Token) (*ast.StringNode, error) {
	node := ast.String(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newLiteralNode(ctx *context, tk *Token) (*ast.LiteralNode, error) {
	node := ast.Literal(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newTagNode(ctx *context, tk *Token) (*ast.TagNode, error) {
	node := ast.Tag(tk.RawToken())
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newSequenceNode(ctx *context, tk *Token, isFlow bool) (*ast.SequenceNode, error) {
	node := ast.Sequence(tk.RawToken(), isFlow)
	node.SetPath(ctx.path)
	if err := setLineComment(ctx, node, tk); err != nil {
		return nil, err
	}
	return node, nil
}

func newTagDefaultScalarValueNode(ctx *context, tag *token.Token) (ast.ScalarNode, error) {
	pos := *(tag.Position)
	pos.Column++

	var (
		tk   *Token
		node ast.ScalarNode
	)
	switch token.ReservedTagKeyword(tag.Value) {
	case token.IntegerTag:
		tk = &Token{Token: token.New("0", "0", &pos)}
		n, err := newIntegerNode(ctx, tk)
		if err != nil {
			return nil, err
		}
		node = n
	case token.FloatTag:
		tk = &Token{Token: token.New("0", "0", &pos)}
		n, err := newFloatNode(ctx, tk)
		if err != nil {
			return nil, err
		}
		node = n
	case token.StringTag, token.BinaryTag, token.TimestampTag:
		tk = &Token{Token: token.New("", "", &pos)}
		n, err := newStringNode(ctx, tk)
		if err != nil {
			return nil, err
		}
		node = n
	case token.BooleanTag:
		tk = &Token{Token: token.New("false", "false", &pos)}
		n, err := newBoolNode(ctx, tk)
		if err != nil {
			return nil, err
		}
		node = n
	case token.NullTag:
		tk = &Token{Token: token.New("null", "null", &pos)}
		n, err := newNullNode(ctx, tk)
		if err != nil {
			return nil, err
		}
		node = n
	default:
		return nil, errors.ErrSyntax(fmt.Sprintf("cannot assign default value for %q tag", tag.Value), tag)
	}
	ctx.insertToken(tk)
	ctx.goNext()
	return node, nil
}

func setLineComment(ctx *context, node ast.Node, tk *Token) error {
	if tk == nil || tk.LineComment == nil {
		return nil
	}
	comment := ast.CommentGroup([]*token.Token{tk.LineComment})
	comment.SetPath(ctx.path)
	if err := node.SetComment(comment); err != nil {
		return err
	}
	return nil
}

func setHeadComment(cm *ast.CommentGroupNode, value ast.Node) error {
	if cm == nil {
		return nil
	}
	switch n := value.(type) {
	case *ast.MappingNode:
		if len(n.Values) != 0 && value.GetComment() == nil {
			cm.SetPath(n.Values[0].GetPath())
			return n.Values[0].SetComment(cm)
		}
	case *ast.MappingValueNode:
		cm.SetPath(n.GetPath())
		return n.SetComment(cm)
	}
	cm.SetPath(value.GetPath())
	return value.SetComment(cm)
}
