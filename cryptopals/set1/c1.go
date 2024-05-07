package set1

import (
	"encoding/base64"
	"encoding/hex"
	"math"
	"strings"
)

func HexToBase64(h string) string {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

func FixedXor(in1, in2 string) string {
	b1, err := hex.DecodeString(in1)
	if err != nil {
		panic(err)
	}
	b2, err := hex.DecodeString(in2)
	if err != nil {
		panic(err)
	}
	for i := range b1 {
		b1[i] ^= b2[i]
	}

	return hex.EncodeToString(b1)
}

func calculateWeight(bs []byte) float64 {
	freq := make(map[byte]int)
	ignored := 0
	for _, b := range bs {
		if b >= 32 && b <= 126 {
			freq[b]++
		} else if b == 9 || b == 10 || b == 13 {
			freq[b]++
		} else {
			return math.Inf(1)
		}
	}

	chi2 := 0.0
	len := len(bs) - ignored
	for i := byte(0); i < 128; i++ {
		observed := float64(freq[i])
		exFreq, exists := englishFreq[i]
		if !exists {
			ignored++
			continue
		}
		expected := float64(len) * exFreq
		difference := observed - expected
		chi2 += difference * difference / expected
	}

	return chi2
}

func bestByteAndScore(h string) (byte, float64, string) {
	bs, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	xored := make([]byte, len(bs))
	copy(xored, bs)
	bestWeight := 1000.0
	currBest := make([]byte, len(bs))
	bestByte := 0

	for xorByte := 0; xorByte < 256; xorByte++ {
		for i := range bs {
			xored[i] = bs[i] ^ byte(xorByte)
		}
		weight := calculateWeight(xored)
		if weight < bestWeight {
			bestWeight = weight
			copy(currBest, xored)
			bestByte = xorByte
		}
	}

	return byte(bestByte), bestWeight, string(currBest)
}

func CrackSingleByteXor(h string) string {
	_, _, cracked := bestByteAndScore(h)
	return cracked
}

func FindSingleXor(input string) string {
	// Read through all lines
	tmp := strings.Fields(input)
	// Filter out any that are not 60-characters
	// lines := tmp
	lines := make([]string, 0, len(tmp))

	for _, l := range tmp {
		if len(l) == 60 {
			lines = append(lines, l)
		}
	}

	bestScore := 1000.0
	output := ""
	// Get the score of the highest single xor
	for _, l := range lines {
		_, score, out := bestByteAndScore(l)
		if score < bestScore {
			bestScore = score
			output = out
		}
	}

	return output
}
