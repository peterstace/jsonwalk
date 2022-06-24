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
		{
			`""`,
			"",
		},
		{
			`"x"`,
			"x",
		},
		{
			`" "`,
			" ",
		},
		{
			`"\n"`,
			"\n",
		},
		{
			`"\u0023"`,
			"#",
		},
		{
			`"\u0023"`,
			"#",
		},
		{
			`"한"`,
			"한",
		},
		{
			`"\uD55C"`,
			"한",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := jsonwindow.UnmarshalStringToken([]byte(tc.input))
			assertNoError(t, err)
			expectEq(t, tc.want, string(got))
		})
	}
}
