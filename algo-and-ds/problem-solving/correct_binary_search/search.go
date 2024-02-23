package main

import "fmt"

func assert(cond bool) {
	if !cond {
		panic("Cond not met")
	}
}

func search(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	// fmt.Println("---Target:", target, "---")

	for left <= right {
		mid := (right-left)/2 + left
		// fmt.Println(left, mid, right)
		if arr[mid] == target {
			return mid
		} else if target > arr[mid] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

func main() {
	assert(search([]int{1, 3, 4, 5, 8, 9}, 4) == 2)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 1) == 0)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 9) == 5)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 0) == -1)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 10) == -1)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 6) == -1)
	assert(search([]int{1, 3, 4, 5, 8}, 9) == -1)
	assert(search([]int{1, 3, 4, 5, 8}, 8) == 4)
	assert(search([]int{1, 3, 4, 5, 8}, 3) == 1)
	assert(search([]int{1, 3, 4, 5, 8, 9}, 3) == 1)
	assert(search([]int{1}, 6) == -1)
	assert(search([]int{1}, 1) == 0)
	fmt.Println("ok")
}
