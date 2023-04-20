package color_convert

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertString(t *testing.T) {
	assert.Equal(t, "rgb(254 3 10)", string(ConvertColor([]byte("#fe030a"))))
	assert.Equal(t, "rgb(15 13 239)", string(ConvertColor([]byte("#0f0def"))))

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
