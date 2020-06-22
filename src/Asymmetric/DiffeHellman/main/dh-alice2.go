package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Fermat(x *big.Int, n *big.Int, p *big.Int) (res *big.Int) {

	if n.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1)
	}
	tmpx := big.NewInt(0)
	tmpx.Set(x)
	tmpx2 := big.NewInt(0)
	tmpx2.Set(x)
	tmpn := big.NewInt(0)
	tmpn.Set(n)
	tmpn2 := big.NewInt(0)
	tmpn2.Set(n)
	tx := tmpx.Mul(tmpx, tmpx).Mod(tmpx, p)
	tn := tmpn.Div(tmpn, big.NewInt(2))
	res = Fermat(tx, tn, p)
	if tmpn2.And(tmpn2, big.NewInt(1)).Cmp(big.NewInt(0)) != 0 {
		res = res.Mul(res, tmpx2).Mod(res, p)
	}
	return res
}

func main() {

	p := big.NewInt(0)
	PBob := big.NewInt(0)
	a := big.NewInt(0)
	g := big.NewInt(0)
	for idx, args := range os.Args {
		if idx == 1 {
			fContent, err := ioutil.ReadFile(args)
			if err != nil {
				panic(err)
			}
			PBob, _ = PBob.SetString(string(fContent[0:308]), 10)
		}

		if idx == 2 {
			fContent2, err2 := ioutil.ReadFile(args)
			if err2 != nil {
				panic(err2)
			}
			p, _ = p.SetString(string(fContent2[0:308]), 10)
			g, _ = g.SetString(string(fContent2[309:310]), 10)
			a, _ = a.SetString(string(fContent2[311:315]), 10)
			Shared := big.NewInt(0)
			Shared = Fermat(PBob, a, p)
			fmt.Println("This is Alice's shared key:", Shared)
		}

	}

}
