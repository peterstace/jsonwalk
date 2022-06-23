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

func parseNextToken(raw []byte) (Token, error) {
	if len(raw) == 0 {
		return Token{}, io.EOF
	}
	switch raw[0] {
	case ',':
		return Token{CommaToken, raw[:1]}, nil
	case ':':
		return Token{ColonToken, raw[:1]}, nil
	case '{':
		return Token{OpenObjectToken, raw[:1]}, nil
	case '}':
		return Token{CloseObjectToken, raw[:1]}, nil
	case '[':
		return Token{OpenArrayToken, raw[:1]}, nil
	case ']':
		return Token{CloseArrayToken, raw[:1]}, nil
	case '"':
		return parseNextStringToken(raw)
	case 't':
		return parseNextKeywordToken(raw, TrueToken)
	case 'f':
		return parseNextKeywordToken(raw, FalseToken)
	case 'n':
		return parseNextKeywordToken(raw, NullToken)
	default:
		c := raw[0]
		if !isStartNumberChar(c) {
			return Token{}, unexpectedStartOfTokenError(c)
		}
		n := parseNextNumberToken(raw)
		return Token{NumberToken, raw[:n]}, nil
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

func parseNextNumberToken(raw []byte) int {
	i := 1 // Already checked the leading char.
	for {
		if i >= len(raw) || !isNumberChar(raw[i]) {
			return i
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
