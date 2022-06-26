package jsonwindow_test

import (
	"strconv"
	"testing"

	"github.com/peterstace/jsonwindow"
)

func TestUnmarshalValidStringToken(t *testing.T) {
	for i, tc := range []struct {
		input string
		want  string
	}{
		{`""`, ""},
		{`"x"`, "x"},
		{`" "`, " "},
		{`"\n"`, "\n"},
		{`"\u0023"`, "#"},
		{`"한"`, "한"},
		{`"\uD55C"`, "한"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := jsonwindow.UnmarshalStringToken([]byte(tc.input))
			mustBeNoError(t, err)
			expectEq(t, tc.want, string(got))
		})
	}
}

func TestUnmarshalInvalidStringToken(t *testing.T) {
	for i, str := range []string{
		``, `"`, `"X`, `"\x"`, `"\u"`, `"\uXXXX"`, `"\u12"`,
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := jsonwindow.UnmarshalStringToken([]byte(str))
			expectErr(t, err)
		})
	}
}
