package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertString(t *testing.T) {
	assert.Equal(t, "rgb(254 3 10)", string(HexColorToRGB([]byte("#fe030a"))))
	assert.Equal(t, "rgb(15 13 239)", string(HexColorToRGB([]byte("#0f0def"))))

}

func TestHexDigitToInt(t *testing.T) {
	assert.Equal(t, 0, HexDigitToInt(byte('0')))
	assert.Equal(t, 1, HexDigitToInt(byte('1')))
	assert.Equal(t, 9, HexDigitToInt(byte('9')))
	assert.Equal(t, 10, HexDigitToInt(byte('a')))
	assert.Equal(t, 10, HexDigitToInt(byte('A')))
	assert.Equal(t, 15, HexDigitToInt(byte('f')))
}

func TestHexToRGB(t *testing.T) {
	assert.Equal(t, [3]int{254, 3, 10}, HexToRGB(0xfe030a))
	assert.Equal(t, [3]int{15, 13, 239}, HexToRGB(0x0f0def))
}

func TestConvertCSS(t *testing.T) {
	result := ConvertTestString("    color: #fe030a;\n")
	assert.Equal(t, "    color: rgb(254 3 10);\n", result)
	result = ConvertTestString("not a color line\n")
	assert.Equal(t, "not a color line\n", result)
}

func TestConvertCSSFile(t *testing.T) {
	FilePairTest(t, "./input_files/simple.css", "./input_files/simple_expected.css")
	FilePairTest(t, "./input_files/advanced.css", "./input_files/advanced_expected.css")
}

func TestNormalize(t *testing.T) {
	assert.Equal(t, []byte{1, 1}, normalize([]byte{1}))
	assert.Equal(t, []byte{1, 1, 2, 2}, normalize([]byte{1, 2}))
}

func FilePairTest(t *testing.T, inputFile, outputFile string) {
	t.Helper()

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	expected, err := ioutil.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}

	var b strings.Builder

	ConvertCSS(file, &b)
	assert.Equal(t, string(expected), b.String())
}

func ConvertTestString(input string) string {
	r := strings.NewReader(input)
	var b strings.Builder
	ConvertCSS(r, &b)
	return b.String()
}
