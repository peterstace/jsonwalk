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

type invalidStringEscapeError byte

func (i invalidStringEscapeError) Error() string {
	return fmt.Sprintf("invalid char after string escape: %d", i)
}

type invalidHexDigitError byte

func (i invalidHexDigitError) Error() string {
	return fmt.Sprintf("invalid hex digit: %d", i)
}

type outOfRangeStringCharError byte

func (o outOfRangeStringCharError) Error() string {
	return fmt.Sprintf("out of range string char: %d", o)
}
