package jsonwindow_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/peterstace/jsonwindow"
)

func BenchmarkWalkObjectCustom(b *testing.B) {
	fname := os.Getenv("JSON_TEST_FILE")
	if fname == "" {
		b.Skip("JSON_TEST_FILE not set")
	}
	zipped, err := os.ReadFile(fname)
	mustBeNoError(b, err)
	zr, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	mustBeNoError(b, err)
	mustBeEq(b, 1, len(zr.File))
	zf, err := zr.File[0].Open()
	mustBeNoError(b, err)
	raw, err := io.ReadAll(zf)
	zf.Close()

	b.SetBytes(int64(len(raw)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := jsonwindow.WalkObject(raw, nil)
		mustBeNoError(b, err)
	}
}

func BenchmarkWhitespace(b *testing.B) {
	const n = 10_000
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		strs[i] = strings.Repeat("abc", i%10)
	}
	raw, err := json.Marshal(map[string][]string{"": strs})
	mustBeNoError(b, err)

	for desc, fn := range map[string]func(t testing.TB, buf []byte) []byte{
		"compact":          func(_ testing.TB, buf []byte) []byte { return buf },
		"python_no_indent": pythonNoIndent,
		"python_indent":    pythonIndent,
	} {
		b.Run(desc, func(b *testing.B) {
			alt := fn(b, raw)
			b.SetBytes(int64(len(alt)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := jsonwindow.WalkObject(alt, nil)
				mustBeNoError(b, err)
			}
		})
	}

}

func pythonNoIndent(t testing.TB, buf []byte) []byte {
	return pythonFilter(t, buf, "--no-indent")
}

func pythonIndent(t testing.TB, buf []byte) []byte {
	return pythonFilter(t, buf)
}

func pythonFilter(t testing.TB, buf []byte, args ...string) []byte {
	args = append([]string{"-m", "json.tool"}, args...)
	cmd := exec.Command("python", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stdin = bytes.NewReader(buf)
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not run python command: %v %v", err, stdout.String())
	}
	return stdout.Bytes()
}

func BenchmarkWalkObjectKeyword(b *testing.B) {
	const n = 10_000
	vals := make([]any, n)
	for i := 0; i < n; i++ {
		switch i % 3 {
		case 0:
			vals[i] = true
		case 1:
			vals[i] = false
		}
	}
	raw, err := json.Marshal(map[string]any{"": vals})
	mustBeNoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := jsonwindow.WalkObject(raw, nil)
		mustBeNoError(b, err)
	}
}
