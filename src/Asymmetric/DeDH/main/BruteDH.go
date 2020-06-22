package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

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

func bruteforce(pkey *big.Int, g *big.Int, p *big.Int) *big.Int {
	var i int64
	for i = 0; i < 2233720368547; i++ {
		tmpp := big.NewInt(0)
		tmpg := big.NewInt(0)
		tmpkey := big.NewInt(0)
		tmpp.Set(p)
		tmpg.Set(g)
		tmpkey.Set(pkey)
		res := Fermat(tmpg, big.NewInt(i), tmpp)
		if res.Cmp(tmpkey) == 0 {
			fmt.Println("Bruteforce success!")
			return big.NewInt(i)
		}
	}
	return big.NewInt(-1)
}

func main() {
	p := big.NewInt(0)
	PAlice := big.NewInt(0)
	g := big.NewInt(0)
	for idx, args := range os.Args {
		if idx == 1 {
			fContent, err := ioutil.ReadFile(args)
			if err != nil {
				panic(err)
			}
			p, _ = p.SetString(string(fContent[0:308]), 10)
			PAlice, _ = PAlice.SetString(string(fContent[309:617]), 10)
			g, _ = g.SetString(string(fContent[618:619]), 10)
			res := bruteforce(PAlice, g, p)
			fmt.Println("This is the secrect key:", res)

		}
	}
}
