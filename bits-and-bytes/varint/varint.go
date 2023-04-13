package varint

import (
	"encoding/binary"
	"os"
)

const (
	MASK     = 0b01111111
	CONT_BIT = 0b10000000
)

func Encode(i uint64) []byte {
	result := []byte{}

	curr := byte(i & MASK)
	i = i >> 7

	for i != 0 {
		result = append(result, curr|CONT_BIT)
		curr = byte(i & MASK)
		i = i >> 7
	}

	result = append(result, curr)
	return result
}

// start with result = 0 (uint64)
// start with shift = 0 (will increment by 7)
// for each byte in varint:
// - convert to uint64
// - shift by shift value
// - add to result
func Decode(varint []byte) uint64 {
	var result uint64
	shift := 0

	for _, b := range varint {
		result = result | (uint64(b) << shift)
		shift += 7
	}
	return result
}

func Read(filename string) uint64 {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return binary.BigEndian.Uint64(data)
}
