package main

import (
	"fmt"
	"math/big"
)

func main() {
	x := big.NewInt(5)
	p := big.NewInt(15)
	fmt.Print("lalal:", x.Mul(x, x).Mod(x, p))
	fmt.Print("x:", x)
}
