package main

import "fmt"

type Jug struct {
	capacity int
	Volume   int
}

func NewJug(size int) *Jug {
	return &Jug{
		capacity: size,
		Volume:   0,
	}
}

func (j Jug) String() string {
	return fmt.Sprintf("%d/%d gals", j.Volume, j.capacity)
}

func FillUp(j Jug) Jug {
	j.Volume = j.capacity
	return j
}

func DumpOut(j Jug) Jug {
	j.Volume = 0
	return j
}

func (j Jug) SpaceLeft() int {
	return j.capacity - j.Volume
}

func PourInto(dest, source Jug) (Jug, Jug) {
	var amt int
	if source.Volume >= dest.SpaceLeft() {
		amt = dest.SpaceLeft()
	} else {
		amt = source.Volume
	}
	source.Volume -= amt
	dest.Volume += amt
	return dest, source
}

type SearchInfo struct {
	History []string
	State   [2]Jug
}

// func nextStates(jugs [2]Jug) []SearchInfo {
// 	nextStates := []SearchInfo{}
// 	var tmpJug1, tmpJug2 Jug
// 	copy(&tmpJug1, jugs[0])
//
// 	// - Fill up first jug
//
//
// 	// - Fill up second jug
// 	// - Dump first jug
// 	// - Dump second jug
// 	// - Pour 1 into 2
// 	// - Pour 2 into 1
//
// }

// func Search(target int, jugs [2]Jug) []string {
// 	visited := map[[2]Jug]struct{}{jugs: struct{}{}}
// 	to_visit := []SearchInfo{{[]string{""}, jugs}}
//
// 	for len(to_visit) > 0 {
// 		curr := to_visit[0]
// 		to_visit = to_visit[1:]
// 		next := nextStates(curr.State)
//
// 		for _, jug := range(jugs) {
// 			if jug.Volume == target {
// 				return curr.History
// 			}
// 		}
// 		to_visit = append(to_visit, )
// 	}
//
//
// }

func main() {
	three := *NewJug(3)
	five := *NewJug(5)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 0, 0

	five = FillUp(five)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 0,5

	three, five = PourInto(three, five)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 3,2

	three = DumpOut(three)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 0,2

	three, five = PourInto(three, five)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 2,0

	five = FillUp(five)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 2,5

	three, five = PourInto(three, five)
	fmt.Println(three)
	fmt.Println(five)
	fmt.Println()
	// 3,4
}
