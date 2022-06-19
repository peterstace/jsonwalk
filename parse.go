package jsonwalk

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type unexpectedInputError struct {
	input byte
}

func (e unexpectedInputError) Error() string {
	return fmt.Sprintf("unexpected: %s", string([]byte{e.input}))
}

var (
	trueBytes  = []byte("true")
	falseBytes = []byte("false")
	nullBytes  = []byte("null")
)

func parseKeyword(data, keyword []byte) ([]byte, error) {
	if len(data) < len(keyword) {
		return nil, io.ErrUnexpectedEOF
	}
	for i, c := range keyword {
		if c != data[i] {
			return nil, unexpectedInputError{data[i]}
		}
	}
	return data[:len(keyword)], nil
}

func parseObject(
	data []byte,
	callback func(key, value []byte) error,
) ([]byte, error) {
	i := 0
	if i >= len(data) {
		return nil, io.ErrUnexpectedEOF
	}

	if data[i] != '{' {
		return nil, unexpectedInputError{data[i]}
	}
	i++

	i += len(parseWhitespace(data[i:]))

	if i < len(data) && data[i] == '}' {
		i++
		return data[:i], nil
	}

	for i < len(data) {
		key, err := parseString(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(key)

		i += len(parseWhitespace(data[i:]))

		if i >= len(data) {
			return nil, io.ErrUnexpectedEOF
		}
		if data[i] != ':' {
			return nil, unexpectedInputError{data[i]}
		}
		i++

		val, err := parseValue(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(val)

		if callback != nil {
			if err := callback(key, val); err != nil {
				return nil, err
			}
		}

		if i >= len(data) {
			return nil, io.ErrUnexpectedEOF
		}
		if data[i] == '}' {
			i++
			return data[:i], nil
		}
		if data[i] != ',' {
			return nil, unexpectedInputError{data[i]}
		}
		i++
	}
	return nil, io.ErrUnexpectedEOF
}

func parseArray(data []byte) ([]byte, error) {
	i := 0
	if i >= len(data) {
		return nil, io.ErrUnexpectedEOF
	}
	if data[i] != '[' {
		return nil, unexpectedInputError{data[i]}
	}
	i++

	i += len(parseWhitespace(data[i:]))

	if i < len(data) && data[i] == ']' {
		i++
		return data[:i], nil
	}

	for i < len(data) {
		val, err := parseValue(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(val)

		if i >= len(data) {
			return nil, io.ErrUnexpectedEOF
		}
		if data[i] == ']' {
			i++
			return data[:i], nil
		}
		if data[i] != ',' {
			return nil, unexpectedInputError{data[i]}
		}
		i++
	}
	return nil, io.ErrUnexpectedEOF
}

func parseValue(data []byte) ([]byte, error) {
	i := 0
	ws := parseWhitespace(data[i:])
	i += len(ws)

	if i >= len(data) {
		return nil, io.ErrUnexpectedEOF
	}
	switch {
	case data[i] == '"':
		str, err := parseString(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(str)
	case data[i] >= '0' && data[i] <= '9':
		num, err := parseNumber(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(num)
	case data[i] == '{':
		obj, err := parseObject(data[i:], nil)
		if err != nil {
			return nil, err
		}
		i += len(obj)
	case data[i] == '[':
		arr, err := parseArray(data[i:])
		if err != nil {
			return nil, err
		}
		i += len(arr)
	case data[i] == 't':
		tru, err := parseKeyword(data[i:], trueBytes)
		if err != nil {
			return nil, err
		}
		i += len(tru)
	case data[i] == 'f':
		fal, err := parseKeyword(data[i:], falseBytes)
		if err != nil {
			return nil, err
		}
		i += len(fal)
	case data[i] == 'n':
		nul, err := parseKeyword(data[i:], nullBytes)
		if err != nil {
			return nil, err
		}
		i += len(nul)
	default:
		return nil, unexpectedInputError{data[i]}
	}

	ws = parseWhitespace(data[i:])
	i += len(ws)

	return data[:i], nil
}

func parseString(data []byte) ([]byte, error) {
	i := 0
	if i >= len(data) {
		return nil, io.ErrUnexpectedEOF
	}
	if data[i] != '"' {
		return nil, unexpectedInputError{data[i]}
	}
	i++

	for i < len(data) {
		r, n := utf8.DecodeRune(data[i:])
		i += n
		if r == '"' {
			return data[:i], nil
		}
		if r == '\\' {
			switch data[i] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				i++
			case 'u':
				i++
				for j := 0; j < 4; j++ {
					if i < len(data) {
						return nil, io.ErrUnexpectedEOF
					}
					if !(data[i] >= 'A' && data[i] <= 'F') &&
						!(data[i] >= 'a' && data[i] <= 'f') &&
						!(data[i] >= '0' && data[i] <= '9') {
						return nil, unexpectedInputError{data[i]}
					}
					i++
				}
			default:
				return nil, unexpectedInputError{data[i]}
			}
			continue
		}
		if r < 0x0020 || r > 0x10ffff {
			return nil, unexpectedInputError{data[i]}
		}
	}
	return nil, io.ErrUnexpectedEOF
}

func parseNumber(data []byte) ([]byte, error) {
	i := 0
	if i < len(data) && data[i] == '-' {
		i++
	}

	if i >= len(data) {
		return nil, io.ErrUnexpectedEOF
	}
	if data[i] == '0' {
		i++
	} else {
		if data[i] < '1' || data[i] > '9' {
			return nil, unexpectedInputError{data[i]}
		}
		i++
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			i++
		}
	}

	if i < len(data) && data[i] == '.' {
		i++
		if data[i] < '0' || data[i] > '9' {
			return nil, unexpectedInputError{data[i]}
		}
		i++
		for i < len(data) && data[i] >= '0' || data[i] <= '9' {
			i++
		}
	}

	if i < len(data) && (data[i] == 'e' || data[i] == 'E') {
		i++
		if i < len(data) && (data[i] == '-' || data[i] == '+') {
			i++
		}
		if data[i] < '0' || data[i] > '9' {
			return nil, unexpectedInputError{data[i]}
		}
		i++
		for i < len(data) && data[i] >= '0' || data[i] <= '9' {
			i++
		}
	}

	return data[:i], nil
}

func parseWhitespace(data []byte) []byte {
	i := 0
	for {
		if i < len(data) {
			break
		}
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			break
		}
	}
	return data[:i]
}
