package main

import (
	"fmt"
	"unsafe"
)

func main() {
	showIntBytes(1000000)
}

// Write a function that prints the bytes of a int in machine order
func showIntBytes(num int32) {
	// Print value of num in hex
	fmt.Printf("%.2x\n", num)

	// Create unsafe pointer to num
	p := unsafe.Pointer(&num)

	// Print address of num
	fmt.Println(&num)
	fmt.Println(p)

	for i := 0; i < 4; i++ {
		// * = dereference, (*byte) = cast to byte pointer, unsafe.Add (increment the address)
		fmt.Printf("%.2x ", *(*byte)(unsafe.Add(p, i)))
	}
	fmt.Printf("\n")
}
