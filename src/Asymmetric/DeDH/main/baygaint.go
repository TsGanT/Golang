package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

// func powBig(a *big.Int, n *big.Int) *big.Int {
// 	tmp := big.NewInt(0)
// 	tmp.Set(a)
// 	res := big.NewInt(1)
// 	for n.Cmp(big.NewInt(0)) == 1 {
// 		temp := new(big.Int)
// 		tmpn := big.NewInt(0)
// 		tmpn.Set(n)
// 		if tmpn.Mod(tmpn, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
// 			temp.Mul(res, tmp)
// 			res = temp
// 		}
// 		temp = new(big.Int)
// 		temp.Mul(tmp, tmp)
// 		tmp = temp
// 		n = n.Div(n, big.NewInt(2))
// 	}
// 	return res
// }

func babyGaint(p *big.Int, g *big.Int, h *big.Int) int32 {
	s := big.NewInt(0)
	y := big.NewInt(0)
	var i, j int32

	tmpp := big.NewInt(0)
	tmpp.Set(p)
	tmpp2 := big.NewInt(0)
	tmpp2.Set(p)

	buffer := tmpp.Sub(tmpp, big.NewInt(1))
	buffer.Sqrt(buffer)
	n := int32(math.Ceil(float64(buffer.Int64())))
	m := make(map[string]int32)
	for i = 0; i < n; i++ {
		tmpg := big.NewInt(0)
		tmpg.Set(g)
		value := Fermat(tmpg, big.NewInt(int64(i)), p)
		m[value.String()] = i
	}
	p2 := tmpp2.Sub(tmpp2, big.NewInt(2))
	s.Mul(big.NewInt(int64(n)), p2)

	tmpg2 := big.NewInt(0)
	tmpg2.Set(g)
	tmpp3 := big.NewInt(0)
	tmpp3.Set(p)

	c := Fermat(tmpg2, s, tmpp3)
	for j = 0; j < n; j++ {
		tmpp4 := big.NewInt(0)
		tmpp4.Set(p)
		buffer2 := Fermat(c, big.NewInt(int64(j)), tmpp4)
		buffer3 := h.Mul(h, buffer2)
		y.Mod(buffer3, p)
		if val, ok := m[y.String()]; ok {
			r := int32(j * n)
			res := int32(r + val)
			fmt.Println("This is your result:", res)
			return res
		}
	}
	return int32(-1)
}

// func main() {	//This can work!!!
// 	a := big.NewInt(7)
// 	b := big.NewInt(33)
// 	c := big.NewInt(37)
// 	res := babyGaint(c, a, b)
// 	fmt.Println(res)
// }

func main() { //This can not work!!
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
			res := babyGaint(p, g, PAlice)
			fmt.Println("This is the secrect key:", res)

		}
	}
}
