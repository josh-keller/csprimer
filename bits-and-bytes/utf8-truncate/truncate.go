package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// sandwich := 0x1F96A
	num, _ := strconv.ParseInt(os.Args[1], 0, 32)

	bytes := Encode(uint64(num))
	for _, b := range bytes {
		fmt.Printf("%x ", b)
	}
	fmt.Print("\n")
}

func Encode(number uint64) []byte {
	if number&0xffffff80 == 0 {
		return []byte{byte(number)}
	} else if number&0xfffff800 == 0 {
		return []byte{
			byte(number>>6) | 0xc0,
			byte(number&0x3f) | 0x80,
		}
	} else if number&0xffff0000 == 0 {
		return []byte{
			byte(number>>12) | 0xe0,
			byte(number>>6&0x3f) | 0x80,
			byte(number&0x3f) | 0x80,
		}
	} else if number&0xffe00000 == 0 {
		return []byte{
			byte(number>>18) | 0xf0,
			byte(number>>12&0x3f) | 0x80,
			byte(number>>6&0x3f) | 0x80,
			byte(number&0x3f&0x3f) | 0x80,
		}
	}

	return []byte{}
}
