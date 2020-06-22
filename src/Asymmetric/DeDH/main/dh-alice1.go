package main

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func BigOdd(b *big.Int) (res *big.Int) {
	im := big.NewInt(2)
	io := big.NewInt(1023)
	ip := big.NewInt(1)
	ip.Exp(im, io, nil).Add(ip, big.NewInt(1)).Add(ip, b.Mul(b, big.NewInt(2)))
	ip.Add(ip, big.NewInt(2))
	return ip
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

func MillerRabin(a *big.Int, p *big.Int) bool {
	tmp := big.NewInt(0)
	tmpa := big.NewInt(0)
	tmp.Set(p)
	tmpa.Set(a)
	if Fermat(a, p.Add(p, big.NewInt(-1)), tmp).Cmp(big.NewInt(1)) == 0 {
		tmp2 := big.NewInt(0)
		tmp2.Set(tmp)
		tmp21 := big.NewInt(0)
		tmp21.Set(tmp)
		tmp22 := big.NewInt(0)
		tmp22.Set(tmp)
		tmp211 := big.NewInt(0)
		tmp211.Set(tmp)
		tmp212 := big.NewInt(0)
		tmp212.Set(tmp)
		u := tmp2.Add(tmp2, big.NewInt(-1)).Div(tmp2, big.NewInt(2))
		tmp3 := big.NewInt(0)
		tmp4 := big.NewInt(0)
		tmp5 := big.NewInt(0)
		tmp3.Set(u)
		tmp4.Set(u)
		tmp5.Set(u)
		for {
			if tmp3.And(tmp3, big.NewInt(1)).Cmp(big.NewInt(0)) == 0 {
				tmp4x := big.NewInt(0)
				tmp4x.Set(tmp5)
				tmpa2 := big.NewInt(0)
				tmpa2.Set(tmpa)
				t := Fermat(tmpa2, tmp4x, tmp22)
				if t.Cmp(big.NewInt(1)) == 0 {
					tmp5 = tmp5.Div(tmp5, big.NewInt(2))
				} else {
					if t.Cmp(tmp21.Add(tmp21, big.NewInt(-1))) == 0 {
						return true
					}
					return false
				}
			} else {
				tmpa3 := big.NewInt(0)
				tmpa3.Set(tmpa)
				tmpu1 := big.NewInt(0)
				tmpu1.Set(tmp5)
				t := Fermat(tmpa3, tmpu1, tmp211)
				if t.Cmp(big.NewInt(1)) == 0 || t.Cmp(tmp212.Sub(tmp212, big.NewInt(1))) == 0 {

					return true
				}
				return false
			}
		}
	}
	return false
}

func TestMillarRabin(p *big.Int, root *big.Int) *big.Int {
	tmp := big.NewInt(0)
	tmp.Set(p)
	tmp2 := big.NewInt(0)
	tmp2.Set(p)
	if MillerRabin(root, tmp) == false {
		return big.NewInt(0)
	}
	if checkbit(tmp2) == false {
		return big.NewInt(0)
	}
	fmt.Print("Success!! you have found a prime number/0.000001")
	return p
}

func checkbit(p *big.Int) bool {
	count := 0
	for {
		if p.Cmp(big.NewInt(0)) != 0 {
			p.Div(p, big.NewInt(2))
			count = count + 1
		} else {
			break
		}
	}
	if count == 1024 {
		return true
	}
	return false
}

func GetPrime(g *big.Int) *big.Int {
	rand.Seed(time.Now().Unix())
	for i := int64(rand.Intn(1000)); i < (i + 1000); i++ {
		tmpg := big.NewInt(0)
		tmpg.Set(g)
		resp := BigOdd(big.NewInt(i))
		res := TestMillarRabin(resp, tmpg)
		if res.Cmp(big.NewInt(0)) != 0 {
			return res
		}
	}
	return big.NewInt(0)
}

func main() {
	var g int64
	g = 7
	seed := big.NewInt(g)
	pri := GetPrime(seed)
	tmtri := big.NewInt(0)
	tmtri.Set(pri)
	fmt.Print("This is the 1024 prime:", pri)
	AliceA := big.NewInt(1997654)
	ip := big.NewInt(1)
	gpowa := ip.Exp(big.NewInt(g), AliceA, nil)
	AlicePublicKey := gpowa.Mod(gpowa, tmtri)
	fmt.Print("\nThis is Alice's pubilc key:", AlicePublicKey)
	fmt.Print("\nThis is g:", g)
	for idx, args := range os.Args {
		if idx == 1 {
			d1 := pri.String()
			d2 := AlicePublicKey.String()
			d3 := big.NewInt(g).String()
			d := d1 + "\n" + d2 + "\n" + d3
			filename := args
			file, err := os.Create(filename)
			check(err)

			_, err = io.WriteString(file, d)
			check(err)
			file.Close()
		}
		if idx == 2 {
			d4 := pri.String()
			d5 := big.NewInt(g).String()
			d6 := big.NewInt(1997654).String()
			d7 := d4 + "\n" + d5 + "\n" + d6
			filename2 := args
			file2, err := os.Create(filename2)
			check(err)

			_, err = io.WriteString(file2, d7)
			check(err)
			file2.Close()
		}
	}
}
