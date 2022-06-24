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

func (w *Window) PeekTokenType() (TokenType, error) {
	idx := w.cur + countWhitespace(w.buf[w.cur:])
	return peekTokenType(w.buf[idx:])
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
		raw, err := w.completeObject(nil)
		return Value{Type: ObjectValue, Raw: raw}, err
	case OpenArrayToken:
		raw, err := w.completeArray(nil)
		return Value{Type: ArrayValue, Raw: raw}, err
	case CommaToken, ColonToken, CloseObjectToken, CloseArrayToken:
		return Value{}, unexpectedTokenTypeError(tok.Type)
	default:
		panic("unknown token type: " + tok.Type.String())
	}
}

func (w *Window) completeObject(fn func(keyToken, value []byte) error) ([]byte, error) {
	start := w.cur - 1

	// Check for empty object.
	cls, err := w.PeekTokenType()
	if err != nil {
		return nil, err
	}
	if cls == CloseObjectToken {
		if _, err := w.NextToken(); err != nil {
			return nil, err
		}
		return w.buf[start:w.cur], nil
	}

	for {
		// Consume the key. It must be a string.
		key, err := w.NextToken()
		if err != nil {
			return nil, err
		}
		if key.Type != StringToken {
			return nil, unexpectedTokenTypeError(key.Type)
		}

		// Consume the colon separating the key from the value.
		colon, err := w.NextToken()
		if err != nil {
			return nil, err
		}
		if colon.Type != ColonToken {
			return nil, unexpectedTokenTypeError(colon.Type)
		}

		// Consume the value. It doesn't matter what type of value it is.
		val, err := w.NextValue()
		if err != nil {
			return nil, err
		}

		// Use callback if provided.
		if fn != nil {
			if err := fn(key.Raw, val.Raw); err != nil {
				return nil, err
			}
		}

		// Check to see if we're at the end of the object, of if there are
		// more key/value pairs.
		commaOrClose, err := w.NextToken()
		if err != nil {
			return nil, err
		}
		if commaOrClose.Type == CloseObjectToken {
			return w.buf[start:w.cur], nil
		}
		if commaOrClose.Type != CommaToken {
			return nil, unexpectedTokenTypeError(commaOrClose.Type)
		}
	}
}

func (w *Window) WalkNextObject(fn func(keyToken, value []byte) error) error {
	tok, err := w.NextToken()
	if err != nil {
		return err
	}
	if tok.Type != OpenObjectToken {
		return unexpectedTokenTypeError(tok.Type)
	}

	_, err = w.completeObject(fn)
	return err
}

func (w *Window) completeArray(fn func(idx int, value []byte) error) ([]byte, error) {
	start := w.cur - 1

	// Check for empty array.
	cls, err := w.PeekTokenType()
	if err != nil {
		return nil, err
	}
	if cls == CloseArrayToken {
		if _, err := w.NextToken(); err != nil {
			return nil, err
		}
		return w.buf[start:w.cur], nil
	}

	for i := 0; true; i++ {
		// Consume the value. It doesn't matter what type of value it is.
		val, err := w.NextValue()
		if err != nil {
			return nil, err
		}

		// Use callback if provided.
		if fn != nil {
			if err := fn(i, val.Raw); err != nil {
				return nil, err
			}
		}

		// Check to see if we're at the end of the object, or if there are
		// more values.
		commaOrClose, err := w.NextToken()
		if err != nil {
			return nil, err
		}
		if commaOrClose.Type == CloseArrayToken {
			return w.buf[start:w.cur], nil
		}
		if commaOrClose.Type != CommaToken {
			return nil, unexpectedTokenTypeError(commaOrClose.Type)
		}
	}
	panic("unreachable")
}

func (w *Window) WalkNextArray(fn func(idx int, value []byte) error) error {
	tok, err := w.NextToken()
	if err != nil {
		return err
	}
	if tok.Type != OpenArrayToken {
		return unexpectedTokenTypeError(tok.Type)
	}

	_, err = w.completeArray(fn)
	return err
}
