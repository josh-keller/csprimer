package json_test

import (
	"testing"

	"github.com/josh-keller/coding-challenges/json/json"
)

func TestParseObject(t *testing.T) {
	testCases := []struct {
		Str []byte
		Err error
	}{
		{[]byte(`{}`), nil},
		{[]byte(`{\`), json.ErrMissingClosingBrace},
	}

	for _, tc := range testCases {
		_, err := json.Parse(tc.Str)
		if err != tc.Err {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}
