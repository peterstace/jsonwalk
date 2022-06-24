package jsonwindow

import (
	"io"
)

type Token struct {
	Type TokenType
	Raw  []byte
}

func countWhitespace(raw []byte) int {
	i := 0
	for {
		if i >= len(raw) {
			return i
		}
		switch raw[i] {
		case ' ', '\t', '\n', '\r':
			i++
		default:
			return i
		}
	}
}

var peekLookup = map[byte]TokenType{
	',': CommaToken,
	':': ColonToken,
	'{': OpenObjectToken,
	'}': CloseObjectToken,
	'[': OpenArrayToken,
	']': CloseArrayToken,
	'"': StringToken,
	't': TrueToken,
	'f': FalseToken,
	'n': NullToken,
	'-': NumberToken,
	'0': NumberToken,
	'1': NumberToken,
	'2': NumberToken,
	'3': NumberToken,
	'4': NumberToken,
	'5': NumberToken,
	'6': NumberToken,
	'7': NumberToken,
	'8': NumberToken,
	'9': NumberToken,
}

func peekTokenType(raw []byte) (TokenType, error) {
	if len(raw) == 0 {
		return 0, io.EOF
	}
	c := raw[0]
	typ, ok := peekLookup[c]
	if !ok {
		return 0, unexpectedStartOfTokenError(c)
	}
	return typ, nil
}

func parseToken(raw []byte) (Token, error) {
	typ, err := peekTokenType(raw)
	if err != nil {
		return Token{}, err
	}
	switch typ {
	case StringToken:
		return parseNextStringToken(raw)
	case NumberToken:
		return parseNextNumberToken(raw), nil
	case CommaToken:
		return Token{CommaToken, raw[:1]}, nil
	case ColonToken:
		return Token{ColonToken, raw[:1]}, nil
	case OpenObjectToken:
		return Token{OpenObjectToken, raw[:1]}, nil
	case CloseObjectToken:
		return Token{CloseObjectToken, raw[:1]}, nil
	case OpenArrayToken:
		return Token{OpenArrayToken, raw[:1]}, nil
	case CloseArrayToken:
		return Token{CloseArrayToken, raw[:1]}, nil
	case TrueToken:
		return parseNextKeywordToken(raw, TrueToken)
	case FalseToken:
		return parseNextKeywordToken(raw, FalseToken)
	case NullToken:
		return parseNextKeywordToken(raw, NullToken)
	default:
		panic("unknown token type: " + typ.String())
	}
}

var (
	trueBytes  = []byte("true")
	falseBytes = []byte("false")
	nullBytes  = []byte("null")
)

func parseNextKeywordToken(raw []byte, typ TokenType) (Token, error) {
	var keyword []byte
	switch typ {
	case TrueToken:
		keyword = trueBytes
	case FalseToken:
		keyword = falseBytes
	case NullToken:
		keyword = nullBytes
	default:
		panic("unexpected token type: " + typ.String())
	}

	if len(raw) < len(keyword) {
		return Token{}, io.ErrUnexpectedEOF
	}
	for i, c := range keyword {
		if c != raw[i] {
			return Token{}, unexpectedCharWithinTokenError(raw[i])
		}
	}
	return Token{Type: typ, Raw: raw[:len(keyword)]}, nil
}

func parseNextStringToken(raw []byte) (Token, error) {
	i := 1 // Already consumed the start quote char.
	for {
		if i >= len(raw) {
			return Token{}, io.ErrUnexpectedEOF
		}
		c := raw[i]
		i++
		if c == '"' {
			return Token{Type: StringToken, Raw: raw[:i]}, nil
		}
		if c == '\\' {
			if i >= len(raw) {
				return Token{}, io.ErrUnexpectedEOF
			}
			if raw[i] == 'u' {
				i += 4 // Skip the next 4 hex digits.
			} else {
				i++ // Skip the single escaped character.
			}
		}
	}
}

func parseNextNumberToken(raw []byte) Token {
	i := 1 // Already checked the leading char.
	for {
		if i >= len(raw) || !isNumberChar(raw[i]) {
			return Token{NumberToken, raw[:i]}
		}
		i++
	}
}

func isStartNumberChar(c byte) bool {
	return false ||
		(c >= '0' && c <= '9') ||
		c == '-'
}

func isNumberChar(c byte) bool {
	return false ||
		(c >= '0' && c <= '9') ||
		c == '.' ||
		c == '-' ||
		c == 'E' ||
		c == 'e' ||
		c == '+'
}
