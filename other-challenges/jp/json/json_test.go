package json_test

import (
	"strings"
	"testing"

	"github.com/josh-keller/coding-challenges/json/json"
)

func TestParseObject(t *testing.T) {
	testCases := []struct {
		Str string
		Err error
	}{
		{`{}`, nil},
		{`{\`, json.ErrMissingClosingBrace},
	}

	for _, tc := range testCases {
		_, err := json.Parse(strings.NewReader(tc.Str))
		if err != tc.Err {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}
