package jsonwindow

type Window struct {
	cur int
	buf []byte
}

func New(raw []byte) *Window {
	return &Window{0, raw}
}

func (w *Window) NextToken() (Token, error) {
	tok, err := parseNextToken(w.buf[w.cur:])
	if err != nil {
		return Token{}, err
	}
	w.cur += len(tok.Raw)
	return tok, nil
}
