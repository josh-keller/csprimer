package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/josh-keller/coding-challenges/json/json"
)

func main() {
	_, err := json.Parse(bufio.NewReader(os.Stdin))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("ok!")
	os.Exit(0)
}
