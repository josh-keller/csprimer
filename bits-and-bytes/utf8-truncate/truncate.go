package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cases, err := os.Open("./cases")
	must(err)
	expected, err := os.Open("./expected")
	caseScanner := bufio.NewScanner(cases)
	expScanner := bufio.NewScanner(expected)

	for caseScanner.Scan() && expScanner.Scan() {
		currCase := caseScanner.Bytes()
		currExp := expScanner.Text()
		trunc := uint8(currCase[0])
		text := string(currCase[1:])
		truncText := Truncate(currCase[1:], int(trunc))
		fmt.Printf("Input: '%s' will be trucated at %d bytes\n", text, trunc)
		fmt.Printf("Expect: '%s'\n", currExp)
		fmt.Printf("Output: '%s'\n", truncText)
		fmt.Printf("Matches: %t\n-----------\n", truncText == currExp)
	}
}

func Truncate(b []byte, truncIdx int) string {
	if int(truncIdx) >= len(b) {
		return string(b)
	}

	for truncIdx >= 0 {
		// If this is not the start of UTF-8 or this is a zero-width-joiner, keep going backwards
		if !isStartingByte(b[truncIdx]) || (truncIdx < len(b)-3 && isZWJ(b[truncIdx:truncIdx+3])) {
			truncIdx--
			continue
		}

		// If the codepoint before this is a ZWJ, skip it
		if truncIdx >= 3 && isZWJ(b[truncIdx-3:truncIdx]) {
			truncIdx -= 3
			continue
		}

		// This is a starting byte that is not a ZWJ, and not preceded by a ZWJ
		break
	}

	return string(b[:truncIdx])
}

func isZWJ(b []byte) bool {
	return b[0] == 0xe2 && b[1] == 0x80 && b[2] == 0x8d
}

// Starting bytes begin with 00, 01, or 11
func isStartingByte(b byte) bool {
	return b&0xc0 != 0x80
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}

func Encode(number uint64) []byte {
	if number&0xffffff80 == 0 {
		return []byte{byte(number)}
	} else if number&0xfffff800 == 0 {
		return []byte{
			byte(number>>6) | 0xc0,
			byte(number&0x3f) | 0x80,
		}
	} else if number&0xffff0000 == 0 {
		return []byte{
			byte(number>>12) | 0xe0,
			byte(number>>6&0x3f) | 0x80,
			byte(number&0x3f) | 0x80,
		}
	} else if number&0xffe00000 == 0 {
		return []byte{
			byte(number>>18) | 0xf0,
			byte(number>>12&0x3f) | 0x80,
			byte(number>>6&0x3f) | 0x80,
			byte(number&0x3f) | 0x80,
		}
	}

	return []byte{}
}
