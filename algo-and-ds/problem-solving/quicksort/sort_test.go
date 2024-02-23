package main

import (
	"reflect"
	"testing"
)

func TestMergeSort(t *testing.T) {
	testCases := []struct {
		arr  []int
		want []int
	}{
		{arr: []int{4, 8, 1, 2, 7, 3}, want: []int{1, 2, 3, 4, 7, 8}},
		{arr: []int{8, 7, 4, 3, 2, 1}, want: []int{1, 2, 3, 4, 7, 8}},
		{arr: []int{8, 8, 1, 8, 1, 1}, want: []int{1, 1, 1, 8, 8, 8}},
	}

	for _, tc := range testCases {
		got := MSort(tc.arr)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("Got %v, want %v", got, tc.want)
		}

	}
}
func TestQuickSort(t *testing.T) {
	testCases := []struct {
		arr  []int
		want []int
	}{
		{arr: []int{4, 8, 1, 2, 7, 3}, want: []int{1, 2, 3, 4, 7, 8}},
		{arr: []int{8, 7, 4, 3, 2, 1}, want: []int{1, 2, 3, 4, 7, 8}},
		{arr: []int{8, 8, 1, 8, 1, 1}, want: []int{1, 1, 1, 8, 8, 8}},
	}

	for _, tc := range testCases {
		QSort(tc.arr)
		if !reflect.DeepEqual(tc.arr, tc.want) {
			t.Fatalf("Got %v, want %v", tc.arr, tc.want)
		}

	}
}
