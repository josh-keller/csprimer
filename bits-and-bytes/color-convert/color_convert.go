package color_convert

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var css_color_re = regexp.MustCompile(`#[a-z0-9]{6}`)

func ConvertColor(hexStr []byte) []byte {
	hexInt, err := strconv.ParseInt(string(bytes.TrimLeft(hexStr, "#")), 16, 32)
	if err != nil {
		panic(err)
	}

	rgb := HexToRGB(int(hexInt))
	return []byte(fmt.Sprintf("rgb(%d %d %d)", rgb[0], rgb[1], rgb[2]))
}

func HexToRGB(hex int) [3]int {
	red := hex >> 16
	green := hex & 0x00ff00 >> 8
	blue := hex & 0x0000ff
	return [3]int{red, green, blue}
}

func ConvertCSS(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		output.Write(css_color_re.ReplaceAllFunc(scanner.Bytes(), ConvertColor))
		output.Write([]byte("\n"))
	}
}

func main() {
	ConvertCSS(os.Stdin, os.Stdout)
}
