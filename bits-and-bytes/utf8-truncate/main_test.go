package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	cases, err := os.Open("./cases")
	must(err)
	expected, err := os.Open("./expected")
	must(err)
	caseScanner := bufio.NewScanner(cases)
	expScanner := bufio.NewScanner(expected)

	for caseScanner.Scan() && expScanner.Scan() {
		currCase := caseScanner.Bytes()
		expectedText := expScanner.Text()
		trunc := uint8(currCase[0])
		truncText := Truncate(currCase[1:], int(trunc))
		assert.Equal(t, expectedText, truncText)
	}
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		unicode uint64
		utf8    []byte
	}{
		{0x0053, []byte{0x53}},
		{0x007f, []byte{0x7f}},
		{0x0080, []byte{0xc2, 0x80}},
		{0x00bf, []byte{0xc2, 0xbf}},
		{0x00c0, []byte{0xc3, 0x80}},
		{0x07ff, []byte{0xdf, 0xbf}},
		{0x0800, []byte{0xe0, 0xa0, 0x80}},
		{0x085b, []byte{0xe0, 0xa1, 0x9b}},
		{0x0c00, []byte{0xe0, 0xb0, 0x80}},
		{0x1f308, []byte{0xf0, 0x9f, 0x8c, 0x88}},
		{0x1f96a, []byte{0xf0, 0x9f, 0xa5, 0xaa}},
		{0x1000ff, []byte{0xf4, 0x80, 0x83, 0xbf}},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.utf8, Encode(tc.unicode))
	}
}
