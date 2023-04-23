package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

var hexColorRegex = regexp.MustCompile(`#[a-fA-F0-9]+`)

func normalize(hexBytes []byte) []byte {
	result := make([]byte, 0, 2*len(hexBytes))

	for _, b := range hexBytes {
		result = append(result, b, b)
	}

	return result
}

func HexColorToRGB(hexBytes []byte) []byte {
	hexBytes = hexBytes[1:]
	intsFromBytes := make([]int, 8)
	if len(hexBytes) <= 4 {
		hexBytes = normalize(hexBytes)
	}

	for i := 0; i < len(hexBytes); i++ {
		intsFromBytes[i] = HexDigitToInt(hexBytes[i])
	}

	red := intsFromBytes[0]<<4 | intsFromBytes[1]
	green := intsFromBytes[2]<<4 | intsFromBytes[3]
	blue := intsFromBytes[4]<<4 | intsFromBytes[5]

	if len(hexBytes) == 6 {
		return []byte(fmt.Sprintf("rgb(%d %d %d)", red, green, blue))
	}

	alpha := float64(intsFromBytes[6]<<4|intsFromBytes[7]) / 255.0
	return []byte(fmt.Sprintf("rgba(%d %d %d / %.5f)", red, green, blue, alpha))
}

func HexDigitToInt(hex byte) int {
	bit7 := hex & 0b01000000
	addend := (bit7 >> 3) | (bit7 >> 6)
	return (int(hex) & 0x0f) + int(addend)
}

func ConvertCSS(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		output.Write(hexColorRegex.ReplaceAllFunc(scanner.Bytes(), HexColorToRGB))
		output.Write([]byte("\n"))
	}
}

// Old implentation when I was using the standard library to parse the hex to an single integer
func HexToRGB(hex int) [3]int {
	red := hex >> 16
	green := hex & 0x00ff00 >> 8
	blue := hex & 0x0000ff
	return [3]int{red, green, blue}
}

func main() {
	ConvertCSS(os.Stdin, os.Stdout)
}
