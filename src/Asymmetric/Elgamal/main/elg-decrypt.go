package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
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

func GCMDecryption(k string, ciphertext []byte, nonce []byte) ([]byte, error) {
	key, _ := hex.DecodeString(k)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext, nil
}

func GetSHA256HashCode(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
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

	gb := big.NewInt(0)
	ga := big.NewInt(0)
	gab := big.NewInt(0)

	p := big.NewInt(0)
	//PElgamal := big.NewInt(0)
	a := big.NewInt(0)
	g := big.NewInt(0)

	var ciphertext []byte
	var nonce []byte
	for idx, args := range os.Args {
		if idx == 1 {
			fContent, err := ioutil.ReadFile(args)
			if err != nil {
				panic(err)
			}
			gb, _ = gb.SetString(string(fContent[0:308]), 10)
			nonce, _ = hex.DecodeString(string(fContent[309:333]))
			ciphertext, _ = hex.DecodeString(string(fContent[333:]))
		}

		if idx == 2 {
			fContent2, err2 := ioutil.ReadFile(args)
			if err2 != nil {
				panic(err2)
			}
			p, _ = p.SetString(string(fContent2[0:308]), 10)
			g, _ = g.SetString(string(fContent2[309:310]), 10)
			a, _ = a.SetString(string(fContent2[311:315]), 10)
			tmpp := big.NewInt(0)
			tmpp.Set(p)
			tmpg := big.NewInt(0)
			tmpg.Set(g)
			tmpa := big.NewInt(0)
			tmpa.Set(a)
			tmpgb := big.NewInt(0)
			tmpgb.Set(gb)
			ga = Fermat(tmpg, tmpa, tmpp)
			gab = Fermat(tmpgb, a, p)
			bytegb := []byte(gb.String())
			bytegab := []byte(gab.String())
			bytega := []byte(ga.String())
			tt := append(bytega, bytegb...)
			tt = append(tt, bytegab...)
			k := GetSHA256HashCode(tt)
			//fmt.Println("here is your key:", k)
			plaintext, e := GCMDecryption(k, ciphertext, nonce)
			if e != nil {
				fmt.Println("Your decyption failed!!")
			}
			fmt.Printf("%s\n", plaintext)
		}

	}
}
