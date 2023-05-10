package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/josh-keller/csprimer/image-rotate/bitmap"
)

func main() {
	// Define output file flag
	outputFile := flag.String("output", "", "output file name")
	flag.StringVar(outputFile, "o", "", "output file name (shorthand)")
	flag.Parse()

	// Determine input file
	var input io.Reader
	if len(flag.Args()) > 0 {
		inputFile := flag.Arg(0)
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening input file:", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	} else if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		input = os.Stdin
	} else {
		fmt.Fprintln(os.Stderr, "No input file specified and no input piped to stdin.")
		os.Exit(1)
	}

	// Determine output file
	var output io.Writer
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output file:", err)
			os.Exit(1)
		}
		defer f.Close()
		output = f
	} else if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		output = os.Stdout
	} else {
		var defaultOutputName string
		if inputFile, ok := input.(*os.File); ok {
			if inputFile.Name() == os.Stdin.Name() {
				defaultOutputName = "out.bmp"
			} else {
				baseName := filepath.Base(inputFile.Name())
				ext := filepath.Ext(inputFile.Name())
				defaultOutputName = fmt.Sprintf("%s_out%s", baseName[:len(baseName)-len(ext)], ext)
			}
		} else {
			fmt.Fprintln(os.Stderr, "No output file specified and no output piped to stdout.")
			os.Exit(1)
		}

		f, err := os.Create(defaultOutputName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output file:", err)
			os.Exit(1)
		}
		defer f.Close()
		output = f
	}

	bmp := bitmap.NewFromReader(input)
	rotated := bitmap.Rotate(bmp)
	binary.Write(output, binary.LittleEndian, &rotated.Header)
	output.Write(rotated.RawInfo)
	output.Write(rotated.Image)

	// Print output file name
	if outputFile == nil {
		fmt.Fprintln(os.Stderr, "Output written to", output)
	}
}
