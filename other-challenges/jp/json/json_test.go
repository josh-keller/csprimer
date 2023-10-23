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
		{`{"key1": "value1"}`, nil},
	}

	for _, tc := range testCases {
		_, err := json.Parse(strings.NewReader(tc.Str))
		if err != tc.Err {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}
