package main

import (
	"fmt"
	"math/big"
	//"io/ioutil"
)

func BigOdd(b *big.Int) (res *big.Int) {
	im := big.NewInt(2)
	io := big.NewInt(1023)
	ip := big.NewInt(1)
	ip.Exp(im, io, nil).Add(ip, big.NewInt(1)).Add(ip, b.Mul(b, big.NewInt(2)))
	ip.Add(ip, big.NewInt(2))
	return ip
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	// a := 12222
	// d1 := []byte(a)
	// d2 := []byte("ahahhahah")
	// d2 = append(d2,d1...)
	// err := ioutil.WriteFile("test.txt", d2, 0644)
	a := BigOdd(big.NewInt(1))
	a2 := BigOdd(big.NewInt(2))
	s := a.String()
	s2 := a2.String()
	s3 := s + "\n" +s2
    // check(err)
	// h := "hello"
	// w := "world"
	// r := h+w
	fmt.Print(s3)
	a = a.Add(a, big.NewInt(10))
	fmt.Print(a)
}

// package main

// import (
//     "fmt"
//     "os"
//     "strconv"
// )

// func main () {
//     for idx, args := range os.Args {
// 		if idx == 1{
// 			fmt.Println("参数" + strconv.Itoa(idx) + ":", args)
// 			fmt.Print(idx)
// 		}
		
//     }
// }

if idx == 1 {
	d1 := pri.Bytes()
	d2 := AlicePublicKey.Bytes()
	d3 := big.NewInt(g).Bytes()
	n := []byte("\n")
	d := append(d1, n...)
	d = append(d, d2...)
	d = append(d, n...)
	d = append(d, d3...)
	err := ioutil.WriteFile(args, d, 0644)
	check(err)
}
if idx == 2 {
	d4 := pri.Bytes()
	d5 := big.NewInt(g).Bytes()
	d6 := big.NewInt(1997).Bytes()
	n2 := []byte("\n")
	d7 := append(d4, n2...)
	d7 = append(d7, d5...)
	d7 = append(d7, n2...)
	d7 = append(d7, d6...)
	err := ioutil.WriteFile(args, d7, 0644)
	check(err)
}