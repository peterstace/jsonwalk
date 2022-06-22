package jsonwindow_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/peterstace/jsonwindow"
)

// TODO: test keyword tokens: true, false, null.

var (
	validStringTokens = []string{
		`""`, `"hello"`,
		`"\u1234"`, `"_\u1234"`, `"\u1234_"`,
		`"\""`, `"\\"`, `"\/"`, `"\b"`, `"\f"`, `"\n"`, `"\r"`, `"\t"`,
		`"_\""`, `"_\\"`, `"_\/"`, `"_\b"`, `"_\f"`, `"_\n"`, `"_\r"`, `"_\t"`,
		`"\"_"`, `"\\_"`, `"\/_"`, `"\b_"`, `"\f_"`, `"\n_"`, `"\r_"`, `"\t_"`,
		`"_\"_"`, `"_\\_"`, `"_\/_"`, `"_\b_"`, `"_\f_"`, `"_\n_"`, `"_\r_"`, `"_\t_"`,
		`"£"`, // 2 bytes
		`"€"`, // 3 bytes
		`"🍣"`, // 4 bytes
		`"_£"`, `"_€"`, `"_🍣"`,
		`"£_"`, `"€_"`, `"🍣_"`,
	}
)

func TestTokenValidString(t *testing.T) {
	for i, tc := range validStringTokens {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			win := jsonwindow.New([]byte(tc))
			got, err := win.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got.Raw) != tc {
				t.Errorf("got: %v want: %v", string(got.Raw), tc)
			}
			if got.Type != jsonwindow.StringToken {
				t.Errorf("expected string token but got: %v", got.Type)
			}
			if _, err := win.NextToken(); err != io.EOF {
				t.Errorf("expected io.EOF after error but got: %v", err)
			}
		})
	}
}

func TestTokenTruncatedString(t *testing.T) {
	for i, goodTok := range validStringTokens {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var badTokens []string
			for j := 1; j < len(goodTok); j++ {
				badTokens = append(badTokens, goodTok[:j])
			}
			for j, badTok := range badTokens {
				t.Run(strconv.Itoa(j), func(t *testing.T) {
					win := jsonwindow.New([]byte(badTok))
					_, err := win.NextToken()
					if err == nil {
						t.Error("expected error but got nil")
					}
				})
			}
		})
	}
}

func TestTokenValidNumber(t *testing.T) {
	for i, tc := range []string{
		`0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`,
		`10`, `11`, `12`, `13`, `14`, `15`, `16`, `17`, `18`, `19`,
		`-0`, `-1`,
		`0.5`, `0.55`,
		`0e-1`, `0E-1`, `0e+1`, `0e-12`,
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			win := jsonwindow.New([]byte(tc))
			got, err := win.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got.Raw) != tc {
				t.Errorf("got: %v want: %v", string(got.Raw), tc)
			}
			if got.Type != jsonwindow.NumberToken {
				t.Errorf("expected number token but got: %v", got.Type)
			}
			if _, err := win.NextToken(); err != io.EOF {
				t.Errorf("expected io.EOF after error but got: %v", err)
			}
		})
	}
}

func TestTokenWhitespace(t *testing.T) {
	for i, tc := range []string{
		" ", "\t", "\r", "\n",
		"  ", "\t\t", "\r\r", "\n\n",
		" \t", " \r", " \n",
		"\t ", "\r ", "\n ",
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			win := jsonwindow.New([]byte(tc))
			got, err := win.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got.Raw) != tc {
				t.Errorf("got: %v want: %v", string(got.Raw), tc)
			}
			if got.Type != jsonwindow.WhitespaceToken {
				t.Errorf("expected whitespace token but got: %v", got.Type)
			}
			if _, err := win.NextToken(); err != io.EOF {
				t.Errorf("expected io.EOF after error but got: %v", err)
			}
		})
	}
}

func TestTokenDelim(t *testing.T) {
	for c, typ := range map[byte]jsonwindow.TokenType{
		':': jsonwindow.ColonToken,
		',': jsonwindow.CommaToken,
		'{': jsonwindow.OpenObjectToken,
		'[': jsonwindow.OpenArrayToken,
		'}': jsonwindow.CloseObjectToken,
		']': jsonwindow.CloseArrayToken,
	} {
		t.Run(string([]byte{c}), func(t *testing.T) {
			win := jsonwindow.New([]byte{c})
			got, err := win.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got.Raw) != 1 || got.Raw[0] != c {
				t.Errorf("got: %v want: %v", string(got.Raw), string([]byte{c}))
			}
			if got.Type != typ {
				t.Errorf("expected %s but got %v", typ, got.Type)
			}
			if _, err := win.NextToken(); err != io.EOF {
				t.Errorf("expected io.EOF after error but got: %v", err)
			}
		})
	}
}
