package varint

import (
	"encoding/binary"
	"os"
)

const MASK = 0b01111111

func Encode(i uint64) []byte {
	result := []byte{}

	curr := byte(i & MASK)
	i = i >> 7

	for i != 0 {
		result = append(result, 0b10000000|curr)
		curr = byte(i & MASK)
		i = i >> 7
	}

	result = append(result, curr)
	return result
}

func Read(filename string) uint64 {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return binary.BigEndian.Uint64(data)
}
