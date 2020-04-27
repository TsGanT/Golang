package main

import (
	"fmt"
	"math/big"
)

var zero *big.Int = big.NewInt(0)
var one *big.Int = big.NewInt(1)
var two *big.Int = big.NewInt(2)

func BigOdd(b *big.Int) (res *big.Int) {
	im := big.NewInt(2)
	io := big.NewInt(1023)
	ip := one
	ip2 := one
	res = ip.Exp(im, io, nil).Add(ip, ip2.Mul(ip2, b))
	return res
}

func Fermat(x *big.Int, n *big.Int, p *big.Int) (res *big.Int) {
	if n.Cmp(zero) == 0 {
		return one
	}
	res = Fermat(x.Mul(x, x).Mod(x, p), n.Div(n, two), p)
	if n.And(n, one).Cmp(zero) != 0 {
		res = res.Mul(res, x).Mod(res, p)
	}
	return res
}

func MillerRabin(a *big.Int, p *big.Int) bool {
	fmt.Print("    p:", Fermat(a, p.Sub(p, one), p))
	if Fermat(a, p.Sub(p, one), p).Cmp(one) == 0 {
		fmt.Print("-------------------test1", p)
		u := p.Sub(p, one).Div(p, two)
		fmt.Print("-------------------test2", u)
		for {
			if u.And(u, one).Cmp(zero) == 0 {
				t := Fermat(a, u, p)
				if t.Cmp(one) == 0 {
					u = u.Div(u, two)
				} else {
					if t.Cmp(p.Sub(p, one)) == 0 {
						return true
					}
					return false
				}
			} else {
				t := Fermat(a, u, p)
				if t.Cmp(one) == 0 || t.Cmp(p.Sub(p, one)) == 0 {
					return true
				}
				return false
			}
		}
	}
	return false
}

func TestMillarRabin(p *big.Int, root *big.Int) *big.Int {
	for k := 0; k < 7; k++ {
		if MillerRabin(root, p) == false {
			return zero
		}
	}
	fmt.Print("Success!! you have found a prime number/0.000001")
	return p
}

func GetPrime(g *big.Int) *big.Int {
	var i int64
	for i = 100; i < (i + 1000); i++ {
		iuse := big.NewInt(i)
		resp := BigOdd(iuse)
		fmt.Print("		resp:", resp)
		res := TestMillarRabin(resp, g)
		if res.Cmp(zero) == 0 {
			// if checkbit(res) {
			return res
			// }
		}
	}
	return zero
}

func checkbit(p int) bool {
	count := 0
	for {
		if p < 0 {
			break
		}
		p = p >> 1
		count = count + 1
	}
	if count == 1024 {
		return true
	}
	return false
}

func main() {
	seed := big.NewInt(7)
	pri := GetPrime(seed)
	fmt.Print("This is the 1024 prime:", pri)
}

// func main() {
// 	//a := BigOdd(5)
// 	// for i:=1; i<100; i++{
// 	// 	a := BigOdd(i)

// 	// }
// 	im := big.NewInt(2)
// 	io := big.NewInt(10)
// 	ip := big.NewInt(1)
// 	ip2 := big.NewInt(1)
// 	ip3 := big.NewInt((2))
// 	ip.Exp(im, io, nil).Add(ip, ip2.Mul(ip2, ip3))
// 	//sip.Mul(im, io) //.Add(ip, im).Div(ip, io)
// 	fmt.Printf("Big Int: %v\n", ip)

// }
