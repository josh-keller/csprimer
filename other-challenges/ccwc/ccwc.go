package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse flags
	// Read file? Stat file?
	//
	countBytes := flag.Bool("c", false, "The number of bytes in each input file is written to the standard output")
	flag.Parse()
	fmt.Println(*countBytes)
}

// TODO: handle error
func CountBytes(filename string) int {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return len(file)
}
