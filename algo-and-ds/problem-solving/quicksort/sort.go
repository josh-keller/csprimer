package main

import "math/rand"

func MSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	first := MSort(arr[:len(arr)/2])
	second := MSort(arr[len(arr)/2:])

	i, j := 0, 0
	merged := make([]int, 0, len(arr))

	for i < len(first) && j < len(second) {
		if first[i] < second[j] {
			merged = append(merged, first[i])
			i++
		} else {
			merged = append(merged, second[j])
			j++
		}
	}

	if i < len(first) {
		merged = append(merged, first[i:]...)
	} else {
		merged = append(merged, second[j:]...)
	}

	return merged
}

func QSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	// Choose a random pivot point and swap with first element
	pivot := rand.Intn(len(arr))
	arr[pivot], arr[0] = arr[0], arr[pivot]

	// Partition the slice
	m := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[0] {
			arr[m+1], arr[i] = arr[i], arr[m+1]
			m++
		}
	}

	// Swap the pivot value to its correct place
	arr[0], arr[m] = arr[m], arr[0]

	// Sort the two halves
	QSort(arr[:m])
	QSort(arr[m+1:])
}
