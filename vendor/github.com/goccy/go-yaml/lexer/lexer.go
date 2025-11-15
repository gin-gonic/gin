package lexer

import (
	"io"

	"github.com/goccy/go-yaml/scanner"
	"github.com/goccy/go-yaml/token"
)

// Tokenize split to token instances from string
func Tokenize(src string) token.Tokens {
	var s scanner.Scanner
	s.Init(src)
	var tokens token.Tokens
	for {
		subTokens, err := s.Scan()
		if err == io.EOF {
			break
		}
		tokens.Add(subTokens...)
	}
	return tokens
}
