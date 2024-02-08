package main

import "fmt"

func main() {
	fmt.Println(exp(2, 2) == 4)
	fmt.Println(exp(2, 3) == 8)
	fmt.Println(exp(2, 4) == 16)
	fmt.Println(exp(2, 5) == 32)
	fmt.Println(exp(2, 10) == 1024)
	fmt.Println(exp(2, 31) == 2147483648)
	fmt.Println(exp(3, 10) == 59049)
}

func exp(b, n int) int {
	if n == 0 {
		return 1
	}

	result := exp(b, n/2)
	return result*result*(n&1)*(b-1) + 1
}
