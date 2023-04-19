package color_convert

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertString(hexStr string) string {
	// Read as a hex digit
	hexInt, err := strconv.ParseInt(strings.TrimLeft(hexStr, "#"), 16, 32)
	if err != nil {
		panic(err)
	}

	rgb := HexToRGB(int(hexInt))

	// Split into three components
	// Incorporate into string
	return fmt.Sprintf("rgb(%d, %d, %d)", rgb[0], rgb[1], rgb[2])
}

func HexToRGB(hex int) [3]int {
	red := hex >> 0x10
	green := hex & 0x00ff00 >> 0x08
	blue := hex & 0x0000ff
	return [3]int{red, green, blue}
}
