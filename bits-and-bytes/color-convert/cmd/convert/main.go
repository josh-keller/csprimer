package main

import (
	"os"

	color_convert "github.com/josh-keller/csprimer/color-convert"
)

func main() {
	color_convert.ConvertCSS(os.Stdin, os.Stdout)
}
