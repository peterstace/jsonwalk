package jsonwindow

import "fmt"

type Window struct {
	cur int
	buf []byte
}

func New(raw []byte) *Window {
	return &Window{0, raw}
}

func (w *Window) NextToken() (Token, error) {
	return Token{}, fmt.Errorf("not implemented yet")
}
