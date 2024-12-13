package main

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func main() {
	idx := slices.IndexFunc([]int{1, 2, 3, 4, 5}, func(i int) bool {
		return i == 3
	})
	fmt.Println(idx)
}
