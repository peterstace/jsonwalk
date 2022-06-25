package jsonwindow

import (
	"encoding/json"
	"fmt"
)

func UnmarshalStringToken(strTok []byte) ([]byte, error) {
	n := len(strTok)
	if strTok[0] != '"' || strTok[n-1] != '"' {
		return nil, fmt.Errorf("invalid string token: must start and end with quote char")
	}

	// TODO: outline slow path?

	// If the string token doesn't use any escapes, we can just strip off the
	// quotes.
	fastPath := true
	for i := 1; i < n-1; i++ {
		if strTok[i] == '\\' {
			fastPath = false
			break
		}
	}
	if fastPath {
		return strTok[1 : n-1], nil
	}

	// The slow path is complicated, so we just rely on the standard
	// encoding/json library.
	var str string
	err := json.Unmarshal(strTok, &str)
	return []byte(str), err
}
