package main

import (
	"fmt"
	"io"
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
	PAlice := big.NewInt(0)
	b := big.NewInt(1968)
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

		}

		if idx == 2 {
			tmpA := big.NewInt(0)
			tmpA.Set(PAlice)
			tmpp := big.NewInt(0)
			tmpp.Set(p)
			tmpp2 := big.NewInt(0)
			tmpp2.Set(p)
			Shared := big.NewInt(0)
			Shared = Fermat(tmpA, big.NewInt(1968), tmpp)
			//fmt.Println("ex:", Shared)
			Bpublic := big.NewInt(0)
			Bpublic = Fermat(g, b, tmpp2)
			d := Bpublic.String()
			filename := args
			file, err := os.Create(filename)
			check(err)
			_, err = io.WriteString(file, d)
			check(err)
			file.Close()
			fmt.Println("Here is the shared key calculated by Bob:", Shared)
		}

	}

}
