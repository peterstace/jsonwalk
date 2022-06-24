package jsonwindow_test

import (
	"strconv"
	"testing"

	"github.com/peterstace/jsonwindow"
)

func TestWalkObject(t *testing.T) {
	for i, tc := range []struct {
		raw      string
		wantKeys []string
		wantVals []string
	}{
		{
			`{}`,
			nil,
			nil,
		},
		{
			`{"x":"y"}`,
			[]string{`"x"`},
			[]string{`"y"`},
		},
		{
			`{"x":"y","z":"w"}`,
			[]string{`"x"`, `"z"`},
			[]string{`"y"`, `"w"`},
		},
		{
			`{"k1":[1,2,3],"k2":{"x":"y","z":"w"}}`,
			[]string{`"k1"`, `"k2"`},
			[]string{`[1,2,3]`, `{"x":"y","z":"w"}`},
		},
		{
			// "python -m json.tool --no-indent" style whitespace.
			`{"k1": [1, 2, 3], "k2": {"x": "y", "z": "w"}}`,
			[]string{`"k1"`, `"k2"`},
			[]string{`[1, 2, 3]`, `{"x": "y", "z": "w"}`},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log("input:")
			logBytes(t, []byte(tc.raw))
			var gotKeys, gotVals []string
			err := jsonwindow.WalkObject([]byte(tc.raw), func(key, val []byte) error {
				gotKeys = append(gotKeys, string(key))
				gotVals = append(gotVals, string(val))
				return nil
			})
			assertNoError(t, err)
			expectDeepEq(t, tc.wantKeys, gotKeys)
			expectDeepEq(t, tc.wantVals, gotVals)
		})
	}
}
