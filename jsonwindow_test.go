package jsonwindow_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/peterstace/jsonwindow"
)

type objectCallSpec struct {
	key, value string
}

func TestWalkObject(t *testing.T) {
	for i, tc := range []struct {
		input     string
		wantCalls []objectCallSpec
	}{
		{
			input:     `{}`,
			wantCalls: nil,
		},
		{
			input: `{"foo":"bar"}`,
			wantCalls: []objectCallSpec{
				{`"foo"`, `"bar"`},
			},
		},
		{
			input: `{"foo":"bar","baz":["a",2,"c"]}`,
			wantCalls: []objectCallSpec{
				{`"foo"`, `"bar"`},
				{`"baz"`, `["a",2,"c"]`},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Logf("input: %v", tc.input)
			var gotCalls []objectCallSpec
			if err := jsonwindow.Object(
				[]byte(tc.input),
				func(key, value []byte) error {
					gotCalls = append(gotCalls, objectCallSpec{
						string(key), string(value),
					})
					return nil
				},
			); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(gotCalls, tc.wantCalls) {
				t.Logf("gotCalls: %v", gotCalls)
				t.Logf("wantCalls: %v", tc.wantCalls)
				t.Errorf("not equal")
			}
		})
	}
}
