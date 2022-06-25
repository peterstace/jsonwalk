package jsonwindow_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
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

// TestWalkObjectSingleValue tests that each possible type of value is captured
// correctly when walking. It's done with a single key to make the test as
// simple as possible and focus just on the values.
func TestWalkObjectSingleValue(t *testing.T) {
	for i, val := range []string{
		// Strings:
		`""`,
		`"X"`, `"XY"`, `"XYZ"`,
		`"\u1234"`, `"_\u1234"`, `"\u1234_"`,
		`"\""`, `"\\"`, `"\/"`, `"\b"`, `"\f"`, `"\n"`, `"\r"`, `"\t"`,
		`"_\""`, `"_\\"`, `"_\/"`, `"_\b"`, `"_\f"`, `"_\n"`, `"_\r"`, `"_\t"`,
		`"\"_"`, `"\\_"`, `"\/_"`, `"\b_"`, `"\f_"`, `"\n_"`, `"\r_"`, `"\t_"`,
		`"_\"_"`, `"_\\_"`, `"_\/_"`, `"_\b_"`, `"_\f_"`, `"_\n_"`, `"_\r_"`, `"_\t_"`,
		`"£"`, // 2 byte utf-8 codepoint
		`"€"`, // 3 byte utf-8 codepoint
		`"🍣"`, // 4 byte utf-8 codepoint
		`"_£"`, `"_€"`, `"_🍣"`,
		`"£_"`, `"€_"`, `"🍣_"`,

		// Numbers:
		`0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`,
		`10`, `11`, `12`, `13`, `14`, `15`, `16`, `17`, `18`, `19`,
		`-0`, `-1`,
		`0.5`, `0.55`,
		`0e-1`, `0E-1`, `0e+1`, `0e-12`,

		// Keywords:
		`true`,
		`false`,
		`null`,

		// Objects:
		`{}`,
		`{"K":"V"}`,
		`{"K1":"V1","K2":"V2"}`,
		`{"K":{}}`,
		`{"K":{"K":"V"}}`,
		`{"K":{"K1":"V1","K2":"V2"}}`,

		// Arrays:
		`[]`,
		`[1]`,
		`[1,2]`,
		`[[1]]`,
		`[[1,2]]`,
		`[[[1]]]`,
		`[[[1,2]]]`,
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			input := []byte(`{"":` + val + `}`)
			t.Log("input:")
			logBytes(t, input)
			var roundtrip string
			err := jsonwindow.WalkObject(input, func(_, val []byte) error {
				roundtrip = string(val)
				return nil
			})
			assertNoError(t, err)
			expectEq(t, val, roundtrip)
		})
	}
}

// TestWalkObjectWhitespace exhaustively tests scenarios where whitespace is
// inserted into each possible location in an input string, comparing it
// against the standard `encoding/json` library. This ensures that whitespace
// stripping between tokens works correctly.
func TestWalkObjectWhitespace(t *testing.T) {
	const withoutWS = `{"A":"V","B":{},"C":{"K":"V"},"D":{"1":1,"2":2},"E":[],"F":[1],"G":[1,2]}`
	for i := 0; i < len(withoutWS)+1; i++ {
		t.Run(fmt.Sprintf("idx_%d", i), func(t *testing.T) {
			withWS := []byte(withoutWS[:i] + " " + withoutWS[i:])
			t.Log("input:")
			logBytes(t, withWS)
			var gotKeys, gotVals []string
			err := jsonwindow.WalkObject(withWS, func(keyTok, val []byte) error {
				key, err := jsonwindow.UnmarshalStringToken(keyTok)
				assertNoError(t, err)
				gotKeys = append(gotKeys, string(key))
				gotVals = append(gotVals, string(val))
				return nil
			})

			var wantMap map[string]json.RawMessage
			err = json.Unmarshal(withWS, &wantMap)
			assertNoError(t, err)
			var wantPairs [][2]string
			for k, v := range wantMap {
				wantPairs = append(wantPairs, [2]string{k, string(v)})
			}
			sort.Slice(wantPairs, func(i, j int) bool {
				lhs := strings.TrimSpace(wantPairs[i][0])
				rhs := strings.TrimSpace(wantPairs[j][0])
				return lhs < rhs
			})
			var wantKeys, wantVals []string
			for _, kv := range wantPairs {
				wantKeys = append(wantKeys, kv[0])
				wantVals = append(wantVals, kv[1])
			}

			expectDeepEq(t, wantKeys, gotKeys)
			expectDeepEq(t, wantVals, gotVals)
		})
	}
}
