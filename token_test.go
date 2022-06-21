package jsonwindow_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/peterstace/jsonwindow"
)

func TestTokenValidString(t *testing.T) {
	for i, tc := range []string{
		`""`, `"hello"`,
		`"\u1234"`, `"_\u1234"`, `"\u1234_"`,
		`"\""`, `"\\"`, `"\/"`, `"\b"`, `"\f"`, `"\n"`, `"\r"`, `"\t"`,
		`"_\""`, `"_\\"`, `"_\/"`, `"_\b"`, `"_\f"`, `"_\n"`, `"_\r"`, `"_\t"`,
		`"\"_"`, `"\\_"`, `"\/_"`, `"\b_"`, `"\f_"`, `"\n_"`, `"\r_"`, `"\t_"`,
		`"_\"_"`, `"_\\_"`, `"_\/_"`, `"_\b_"`, `"_\f_"`, `"_\n_"`, `"_\r_"`, `"_\t_"`,
		`£`, // 2 bytes
		`€`, // 3 bytes
		`𐍈`, // 4 bytes
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
			if got.Type != jsonwindow.StringToken {
				t.Errorf("expected string token but got: %v", got.Type)
			}
			if _, err := win.NextToken(); err != io.EOF {
				t.Errorf("expected io.EOF after error but got: %v", err)
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

//func TestTokenValid(t *testing.T) {
//	type tokenSpec struct {
//		typ jsonwindow.TokenType
//		txt string
//	}
//	for i, tc := range []struct {
//		input string
//		want  []tokenSpec
//	}{
//		{
//			input: "",
//			want:  nil,
//		},
//		{
//			input: `123`,
//			want:  []tokenSpec{{jsonwindow.NumberToken, `123`}},
//		},
//		{
//			input: " \t\n\r",
//			want:  []tokenSpec{{jsonwindow.WhitespaceToken, " \t\n\r"}},
//		},
//	} {
//		t.Run(strconv.Itoa(i), func(t *testing.T) {
//			win := jsonwindow.New([]byte(tc.input))
//			var got []tokenSpec
//			for {
//				tok, err := win.NextToken()
//				if err == io.EOF {
//					break
//				}
//				if err != nil {
//					t.Fatalf("unexpected error: %v", err)
//				}
//				got = append(got, tokenSpec{tok.Type, string(tok.Raw)})
//			}
//			if !reflect.DeepEqual(got, tc.want) {
//				t.Log("input:")
//				t.Log(hex.Dump([]byte(tc.input)))
//				t.Logf("want: %v", tc.want)
//				t.Logf("got: %v", got)
//				t.Errorf("mismatch")
//			}
//		})
//	}
//}
