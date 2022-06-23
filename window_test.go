package jsonwindow_test

import (
	"encoding/hex"
	"io"
	"reflect"
	"strconv"
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
)

func TestNextToken(t *testing.T) {
	for i, tc := range []struct {
		input string
		want  []jsonwindow.Token
	}{
		{`{}`, []jsonwindow.Token{startObj, closeObj}},
		{`{ }`, []jsonwindow.Token{startObj, closeObj}},
		{` {}`, []jsonwindow.Token{startObj, closeObj}},
		{`{} `, []jsonwindow.Token{startObj, closeObj}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log("input:")
			t.Log(hex.Dump([]byte(tc.input)))
			win := jsonwindow.New([]byte(tc.input))
			var got []jsonwindow.Token
			for {
				tok, err := win.NextToken()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				got = append(got, tok)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Logf("got: %v", got)
				t.Logf("want: %v", tc.want)
				t.Errorf("mismatch")
			}
		})
	}
}
