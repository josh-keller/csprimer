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
	curr := byte(0)

	for {
		curr = byte(i & MASK)
		i = i >> 7

		if i == 0 {
			result = append(result, curr)
			break
		}

		result = append(result, curr|CONT_BIT)
	}

	return result
}

func Decode(varint []byte) uint64 {
	var result uint64
	shift := 0
	i := 0

	for {
		result |= uint64(varint[i]&MASK) << shift
		shift += 7
		if varint[i]&CONT_BIT == 0 {
			return result
		}
		i++
	}
}

func Read(filename string) uint64 {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return binary.BigEndian.Uint64(data)
}
