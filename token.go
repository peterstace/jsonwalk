package jsonwindow

import (
	"fmt"
	"io"
)

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
	WhitespaceToken
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
	case WhitespaceToken:
		return "WhitespaceToken"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}

type Token struct {
	Type TokenType
	Raw  []byte
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
		n, err := parseNextStringToken(raw)
		if err != nil {
			return Token{}, err
		}
		return Token{StringToken, raw[:n]}, nil
	case 't':
		n, err := parseNextKeywordToken(raw, trueBytes)
		if err != nil {
			return Token{}, err
		}
		return Token{TrueToken, raw[:n]}, nil
	case 'f':
		n, err := parseNextKeywordToken(raw, falseBytes)
		if err != nil {
			return Token{}, err
		}
		return Token{FalseToken, raw[:n]}, nil
	case 'n':
		n, err := parseNextKeywordToken(raw, nullBytes)
		if err != nil {
			return Token{}, err
		}
		return Token{NullToken, raw[:n]}, nil
	default:
		c := raw[0]
		switch {
		case isStartNumberChar(c):
			n := parseNextNumberToken(raw)
			return Token{NumberToken, raw[:n]}, nil
		case isWhitespaceChar(c):
			n := parseNextWhitespaceToken(raw)
			return Token{WhitespaceToken, raw[:n]}, nil
		default:
			return Token{}, unexpectedStartOfTokenError(c)
		}
	}
}

func parseNextKeywordToken(raw, keyword []byte) (int, error) {
	if len(raw) < len(keyword) {
		return 0, io.ErrUnexpectedEOF
	}
	for i, c := range keyword {
		if c != raw[i] {
			return 0, unexpectedCharWithinTokenError(raw[i])
		}
	}
	return len(keyword), nil
}

func parseNextStringToken(raw []byte) (int, error) {
	i := 1 // Already consumed the start quote char.
	for {
		if i >= len(raw) {
			return 0, io.ErrUnexpectedEOF
		}
		c := raw[i]
		i++
		if c == '"' {
			return i, nil
		}
		if c == '\\' {
			if i >= len(raw) {
				return 0, io.ErrUnexpectedEOF
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

func parseNextWhitespaceToken(raw []byte) int {
	i := 1 // Already checked the leading char.
	for {
		if i >= len(raw) || !isWhitespaceChar(raw[i]) {
			return i
		}
		i++
	}
}

func isWhitespaceChar(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
