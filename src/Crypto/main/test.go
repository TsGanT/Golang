package main

import (
	"fmt"
)

func permutations(iterable []rune, r int) {
	pool := iterable
	n := len(pool)
	if r > n {
		return
	}
	indices := make([]rune, n)
	for i := range indices {
		indices[i] = i
	}
	//fmt.Println("indicies:", indices)
	cycles := make([]int, r)
	for i := range cycles {
		cycles[i] = n - i

	}
	//fmt.Println("cylcles:", cycles)
	result := make([]rune, r)
	for i, el := range indices[:r] {
		result[i] = pool[el]
		//fmt.Println("i:", i, "   el:", el)
	}
	//fmt.Println("result:", result)

	for n > 0 {
		i := r - 1
		for ; i >= 0; i -= 1 { //i从后往前
			cycles[i] -= 1
			if cycles[i] == 0 {
				index := indices[i]
				for j := i; j < n-1; j += 1 {
					indices[j] = indices[j+1]
				}
				indices[n-1] = index
				cycles[i] = n - i
				//fmt.Println("cycles2:", cycles)
				//fmt.Println("indience2:", indices)
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]
				for k := i; k < r; k += 1 {
					result[k] = pool[indices[k]]
					//fmt.Println("result[k], k", result[k], k)
				}
				fmt.Println("result:", result)
				break
			}
		}
		if i < 0 {
			return
		}
	}
}

func main() {
	fmt.Println("Itertools permutations in Go:")
	// permutations('ABCD', 2) --> AB AC AD BA BC BD CA CB CD DA DB DC
	// permutations(range(3)) --> 012 021 102 120 201 210
	//fmt.Printf("iterable = %s, r = %d", "[]int{1, 2, 3, 4}", 3)
	fmt.Println()
	permutations(permutations('ABCD', 2))

}
