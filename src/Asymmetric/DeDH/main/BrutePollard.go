package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

func ext_euclid(a *big.Int, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	if b.Cmp(big.NewInt(0)) == 0 {
		return a, big.NewInt(1), big.NewInt(0)
	}
	tmpb := big.NewInt(0)
	tmpb2 := big.NewInt(0)
	tmpb.Set(b)
	tmpb2.Set(b)
	tmpa := big.NewInt(0)
	tmpa2 := big.NewInt(0)
	tmpa.Set(b)
	tmpa2.Set(b)
	d, xx, yy := ext_euclid(tmpb, tmpa.Mod(tmpa, tmpb))
	x := yy
	adivb := tmpa2.Div(tmpa2, tmpb2).Mul(tmpa2, yy)
	y := xx.Sub(xx, adivb)
	return d, x, y
}

func inverse(a *big.Int, n *big.Int) *big.Int {
	_, res, _ := ext_euclid(a, n)
	return res
}

func NewXab(x *big.Int, a *big.Int, b *big.Int, p *big.Int, g *big.Int, h *big.Int, Q *big.Int) {
	tmpx := big.NewInt(0)
	tmpx.Set(x)
	c := tmpx.Mod(tmpx, big.NewInt(3))

	if c.Cmp(big.NewInt(0)) == 0 {
		x = x.Mul(x, g).Mod(x, p)
		a = a.Add(a, big.NewInt(1)).Mod(a, Q)
	}
	if c.Cmp(big.NewInt(1)) == 0 {
		x = x.Mul(x, h).Mod(x, p)
		b = b.Add(b, big.NewInt(1)).Mod(b, Q)
	}
	if c.Cmp(big.NewInt(2)) == 0 {
		x = x.Mul(x, x).Mod(x, p)
		a = a.Mul(a, big.NewInt(2)).Mod(a, Q)
		b = b.Mul(b, big.NewInt(2)).Mod(b, Q)
	}
}

func pollard(p *big.Int, g *big.Int, h *big.Int) *big.Int {
	tmpp := big.NewInt(0)
	tmpp.Set(p)
	Q := tmpp.Sub(tmpp, big.NewInt(1)).Div(tmpp, big.NewInt(2))
	tmpg := big.NewInt(0)
	tmpg.Set(g)
	tmph := big.NewInt(0)
	tmph.Set(h)
	x := g.Mul(g, h)
	a := big.NewInt(1)
	b := big.NewInt(1)

	X := big.NewInt(0)
	X.Set(x)
	A := big.NewInt(1)
	B := big.NewInt(1)
	fmt.Println("test2")
	for {
		NewXab(x, a, b, p, g, h, Q)
		NewXab(X, A, B, p, g, h, Q)
		NewXab(X, A, B, p, g, h, Q)
		if x.Cmp(X) == 0 {
			break
		}
	}
	nom := a.Sub(a, A)
	denom := B.Sub(B, b)
	fmt.Println("nom:", nom)
	fmt.Println("denom:", denom)
	if denom.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("Falied!!")
	}
	tmpQ := big.NewInt(0)
	tmpQ.Set(Q)
	rev := inverse(denom, tmpQ.Mul(tmpQ, nom))
	res := rev.Mod(rev, Q)
	return res.Add(res, Q)
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
			//fmt.Println("test")
			res := pollard(p, g, PAlice)
			fmt.Println("This is the secrect key:", res)

		}
	}
}
