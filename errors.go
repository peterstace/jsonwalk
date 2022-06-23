package jsonwindow

import "fmt"

type unexpectedStartOfTokenError byte

func (u unexpectedStartOfTokenError) Error() string {
	return fmt.Sprintf("unexpected char at start of token: %d", u)
}

type unexpectedCharWithinTokenError byte

func (u unexpectedCharWithinTokenError) Error() string {
	return fmt.Sprintf("unexpected char within token: %d", u)
}

type unexpectedTokenTypeError TokenType

func (u unexpectedTokenTypeError) Error() string {
	return fmt.Sprintf("unexpected token of type %s", TokenType(u))
}
