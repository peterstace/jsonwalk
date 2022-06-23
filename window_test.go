package jsonwindow_test

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/peterstace/jsonwindow"
)

var (
	startObj = jsonwindow.Token{
		Type: jsonwindow.OpenObjectToken,
		Raw:  []byte("{"),
	}
	closeObj = jsonwindow.Token{
		Type: jsonwindow.CloseObjectToken,
		Raw:  []byte("}"),
	}
	startArr = jsonwindow.Token{
		Type: jsonwindow.OpenArrayToken,
		Raw:  []byte("["),
	}
	closeArr = jsonwindow.Token{
		Type: jsonwindow.CloseArrayToken,
		Raw:  []byte("]"),
	}
	null = jsonwindow.Token{
		Type: jsonwindow.NullToken,
		Raw:  []byte("null"),
	}
	colon = jsonwindow.Token{
		Type: jsonwindow.ColonToken,
		Raw:  []byte{':'},
	}
	comma = jsonwindow.Token{
		Type: jsonwindow.CommaToken,
		Raw:  []byte{','},
	}
)

func num(raw string) jsonwindow.Token {
	return jsonwindow.Token{
		Type: jsonwindow.NumberToken,
		Raw:  []byte(raw),
	}
}

func str(raw string) jsonwindow.Token {
	return jsonwindow.Token{
		Type: jsonwindow.StringToken,
		Raw:  []byte(raw),
	}
}

func boolean(b bool) jsonwindow.Token {
	raw := []byte(strconv.FormatBool(b))
	typ := jsonwindow.FalseToken
	if b {
		typ = jsonwindow.TrueToken
	}
	return jsonwindow.Token{
		Type: typ,
		Raw:  raw,
	}
}

func TestNextToken(t *testing.T) {
	for i, tc := range []struct {
		input string
		want  []jsonwindow.Token
	}{
		{`1`, []jsonwindow.Token{num(`1`)}},
		{`"x"`, []jsonwindow.Token{str(`"x"`)}},
		{`true`, []jsonwindow.Token{boolean(true)}},
		{`false`, []jsonwindow.Token{boolean(false)}},
		{`null`, []jsonwindow.Token{null}},
		{`{}`, []jsonwindow.Token{startObj, closeObj}},
		{`{"x":1}`, []jsonwindow.Token{startObj, str(`"x"`), colon, num(`1`), closeObj}},
		{`{"x":1,"y":"two"}`, []jsonwindow.Token{startObj, str(`"x"`), colon, num(`1`), comma, str(`"y"`), colon, str(`"two"`), closeObj}},
		{`[]`, []jsonwindow.Token{startArr, closeArr}},
		{`[true]`, []jsonwindow.Token{startArr, boolean(true), closeArr}},
		{`[false,true]`, []jsonwindow.Token{startArr, boolean(false), comma, boolean(true), closeArr}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log("input:")
			logBytes(t, []byte(tc.input))

			// Make sure the parsed tokens match those expected.
			got := getAllTokens(t, []byte(tc.input))
			expectTokenSliceEq(t, got, tc.want)

			// Make sure that when whitespace is added between any two tokens
			// in the stream, that we get the original stream back when
			// re-parsed.
			for j := 0; j < len(got)+1; j++ {
				t.Run(fmt.Sprintf("WS at %d", j), func(t *testing.T) {
					var rawWithWS []byte
					for k, tok := range got {
						if k == j {
							rawWithWS = append(rawWithWS, ' ')
						}
						rawWithWS = append(rawWithWS, tok.Raw...)
					}
					if j == len(got) {
						rawWithWS = append(rawWithWS, ' ')
					}

					t.Log("input:")
					logBytes(t, rawWithWS)

					reparsed := getAllTokens(t, rawWithWS)
					expectTokenSliceEq(t, reparsed, tc.want)
				})
			}
		})
	}
}

func getAllTokens(t *testing.T, input []byte) []jsonwindow.Token {
	t.Helper()
	win := jsonwindow.New([]byte(input))
	var all []jsonwindow.Token
	for {
		tok, err := win.NextToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		all = append(all, tok)
	}
	return all
}

func expectTokenSliceEq(t *testing.T, got, want []jsonwindow.Token) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Logf("got: %v", got)
		t.Logf("want: %v", want)
		t.Fatalf("mismatch")
	}
}

func logBytes(t *testing.T, byts []byte) {
	t.Helper()
	dump := strings.TrimSpace(hex.Dump(byts))
	for _, line := range strings.Split(dump, "\n") {
		t.Log(line)
	}
}
