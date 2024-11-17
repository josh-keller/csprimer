package main

import (
	"fmt"
	"math"
)

var (
	memo    = map[int]int{1: 1}
	squares = []int{1}
)

func main() {
	tc := map[int]int{1: 1, 2: 2, 3: 3, 4: 1, 5: 2, 6: 3, 7: 4, 8: 2, 9: 1, 18: 2, 23: 4, 25: 1, 33: 3, 34: 2}
	for n, want := range tc {
		got := minSquares(n)
		fmt.Printf("Min square sum for %d: %d - %t\n", n, got, want == got)
	}
}

func minSquares(target int) int {
	for i := 1; i <= target; i++ {
		if _, exists := memo[i]; exists {
			continue
		}

		if isPerfectSquare(i) {
			squares = append(squares, i)
			memo[i] = 1
			continue
		}

		min := i

		for _, sq := range squares {
			curr := 1 + memo[i-sq]
			if curr < min {
				min = curr
			}
		}

		memo[i] = min
	}

	return memo[target]
}

func isPerfectSquare(n int) bool {
	sqrt := math.Sqrt(float64(n))
	return math.Trunc(sqrt) == sqrt
}
