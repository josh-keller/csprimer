package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	countBytes := flag.Bool("c", false, "The number of bytes in each input file is written to the standard output")
	countLines := flag.Bool("l", false, "The number of lines in each input file is written to the standard output")
	flag.Parse()

	// TODO: Take more than file
	filename := flag.CommandLine.Arg(0)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	totalBytes := 0
	totalLines := 0
	totalWords := 0
	buffer := make([]byte, 11)
	index := 0

	for {
		n, err := f.Read(buffer[index:])
		totalBytes += n
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			i := bytes.IndexAny(buffer[index:], " \n\r\t\v\f")
			if i >= 0 {
				if i > 0 {
					totalWords++
				}
				if buffer[index+i] == '\n' {
					totalLines++
				}
				index += i + 1
			} else {
				index = 0
				break
			}
		}
	}

	fmt.Printf("%8d%8d%8d %s\n", totalLines, totalWords, totalBytes, filename)

	if *countBytes {
		fmt.Fprintf(os.Stdout, "  %d %s\n", CountBytes(filename), filename)
	}

	if *countLines {
		fmt.Fprintf(os.Stdout, "  %d %s\n", CountLines(filename), filename)
	}

}

// TODO: handle error opening file
func CountBytes(filename string) int {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	totalBytes := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		totalBytes++
	}

	return totalBytes
}

func CountLines(filename string) int {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	totalLines := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		totalLines++
	}

	return totalLines
}
