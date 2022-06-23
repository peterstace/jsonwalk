package jsonwindow

import "fmt"

type TokenType int

const (
	StringToken TokenType = iota + 1
	NumberToken
	CommaToken
	ColonToken
	OpenObjectToken
	CloseObjectToken
	OpenArrayToken
	CloseArrayToken
	TrueToken
	FalseToken
	NullToken
)

func (t TokenType) String() string {
	switch t {
	case StringToken:
		return "StringToken"
	case NumberToken:
		return "NumberToken"
	case CommaToken:
		return "CommaToken"
	case ColonToken:
		return "ColonToken"
	case OpenObjectToken:
		return "OpenObjectToken"
	case CloseObjectToken:
		return "CloseObjectToken"
	case OpenArrayToken:
		return "OpenArrayToken"
	case CloseArrayToken:
		return "CloseArrayToken"
	case TrueToken:
		return "TrueToken"
	case FalseToken:
		return "FalseToken"
	case NullToken:
		return "NullToken"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
