package jsonwindow

type Window struct {
	cur int
	buf []byte
}

func New(raw []byte) *Window {
	return &Window{0, raw}
}

func (w *Window) NextToken() (Token, error) {
	w.cur += countWhitespace(w.buf[w.cur:])
	tok, err := parseToken(w.buf[w.cur:])
	if err != nil {
		return Token{}, err
	}
	w.cur += len(tok.Raw)
	return tok, nil
}

func (w *Window) PeekToken() (Token, error) {
	idx := w.cur + countWhitespace(w.buf[w.cur:])
	tok, err := parseToken(w.buf[idx:])
	if err != nil {
		return Token{}, err
	}
	return tok, nil
}

func (w *Window) NextValue() (Value, error) {
	tok, err := w.NextToken()
	if err != nil {
		return Value{}, err
	}
	switch tok.Type {
	case StringToken:
		return Value{Type: StringValue, Raw: tok.Raw}, nil
	case NumberToken:
		return Value{Type: NumberValue, Raw: tok.Raw}, nil
	case TrueToken:
		return Value{Type: BooleanValue, Raw: tok.Raw}, nil
	case FalseToken:
		return Value{Type: BooleanValue, Raw: tok.Raw}, nil
	case NullToken:
		return Value{Type: NullValue, Raw: tok.Raw}, nil
	case OpenObjectToken:
		start := w.cur - len(tok.Raw)
		for {
			// Consume the key. It must be a string.
			key, err := w.NextToken()
			if err != nil {
				return Value{}, err
			}
			if key.Type != StringToken {
				return Value{}, unexpectedTokenTypeError(key.Type)
			}

			// Consume the colon separating the key from the value.
			colon, err := w.NextToken()
			if err != nil {
				return Value{}, err
			}
			if colon.Type != ColonToken {
				return Value{}, unexpectedTokenTypeError(colon.Type)
			}

			// Consume the value. It doesn't matter what type of value it is.
			if _, err := w.NextValue(); err != nil {
				return Value{}, err
			}

			// Check to see if we're at the end of the object, of if there are
			// more key/value pairs.
			commaOrClose, err := w.NextToken()
			if err != nil {
				return Value{}, err
			}
			if commaOrClose.Type == CloseObjectToken {
				return Value{Type: ObjectValue, Raw: w.buf[start:w.cur]}, nil
			}
			if commaOrClose.Type != CommaToken {
				return Value{}, unexpectedTokenTypeError(commaOrClose.Type)
			}
		}
	case OpenArrayToken:
		start := w.cur - len(tok.Raw)
		for {
			// Consume the value. It doesn't matter what type of value it is.
			if _, err := w.NextValue(); err != nil {
				return Value{}, err
			}

			// Check to see if we're at the end of the object, or if there are
			// more values.
			commaOrClose, err := w.NextToken()
			if err != nil {
				return Value{}, err
			}
			if commaOrClose.Type == CloseArrayToken {
				return Value{Type: ArrayValue, Raw: w.buf[start:w.cur]}, nil
			}
			if commaOrClose.Type != CommaToken {
				return Value{}, unexpectedTokenTypeError(commaOrClose.Type)
			}
		}
	case CommaToken, ColonToken, CloseObjectToken, CloseArrayToken:
		return Value{}, unexpectedTokenTypeError(tok.Type)
	default:
		panic("unknown token type: " + tok.Type.String())
	}
}
