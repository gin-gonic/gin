package scanner

import "github.com/goccy/go-yaml/token"

type InvalidTokenError struct {
	Token *token.Token
}

func (e *InvalidTokenError) Error() string {
	return e.Token.Error
}

func ErrInvalidToken(tk *token.Token) *InvalidTokenError {
	return &InvalidTokenError{
		Token: tk,
	}
}
