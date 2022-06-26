package jsonwindow_test

import (
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

func logBytes(t testing.TB, byts []byte) {
	t.Helper()
	dump := strings.TrimSpace(hex.Dump(byts))
	for _, line := range strings.Split(dump, "\n") {
		t.Log(line)
	}
}

func mustBeNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}
}

func expectErr(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func mustBeTrue(t testing.TB, b bool) {
	if !b {
		t.Fatalf("expected true but got false")
	}
}

func expectDeepEq[T any](t testing.TB, want, got T) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Logf("want: %v", want)
		t.Logf("got: %v", got)
		t.Errorf("not deep equal")
	}
}

func mustBeDeepEq[T any](t testing.TB, want, got T) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Logf("want: %v", want)
		t.Logf("got: %v", got)
		t.Fatalf("not deep equal")
	}
}

func expectEq[T comparable](t testing.TB, want, got T) {
	t.Helper()
	if want != got {
		t.Logf("want: %v", want)
		t.Logf("got: %v", got)
		t.Errorf("not equal")
	}
}

func mustBeEq[T comparable](t testing.TB, want, got T) {
	t.Helper()
	if want != got {
		t.Logf("want: %v", want)
		t.Logf("got: %v", got)
		t.Fatalf("not equal")
	}
}
