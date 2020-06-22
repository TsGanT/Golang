package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

type PublicKey struct {
	G, P, Y *big.Int
}

// var public *PublicKey = new(PublicKey)

// func SetPublicKey(g *big.Int, p *big.Int, y *big.Int) {
// 	public.G = g
// 	public.P = p
// 	public.Y = y
// }

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

func GetSHA256HashCode(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

func GCMEncrypt(key string, data []byte) ([]byte, []byte, error) {
	keyuse, _ := hex.DecodeString(key)
	block, err := aes.NewCipher(keyuse)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return ciphertext, nonce, nil
}

func main() {
	var message string
	var k string

	gexpab := big.NewInt(0)
	gexpb := big.NewInt(0)
	gexpa := big.NewInt(0)

	p := big.NewInt(0)
	PElgamal := big.NewInt(0)
	b := big.NewInt(1968)
	g := big.NewInt(0)
	gb := big.NewInt(0)

	for idx, args := range os.Args {

		if idx == 1 {
			message = args
		}
		if idx == 2 {
			fContent, err := ioutil.ReadFile(args)
			if err != nil {
				panic(err)
			}
			p, _ = p.SetString(string(fContent[0:308]), 10)                 //This is the big prime number
			PElgamal, _ = PElgamal.SetString(string(fContent[309:617]), 10) //This is g^a mod p
			g, _ = g.SetString(string(fContent[618:619]), 10)               //This is g
			//fmt.Println("g", PElgamal)
			gexpa.Set(PElgamal)
			tmpPElgamal := big.NewInt(0)
			tmpPElgamal.Set(PElgamal)
			tmpPElgamal2 := big.NewInt(0)
			tmpPElgamal2.Set(PElgamal)
			tmpp := big.NewInt(0)
			tmpp.Set(p)
			tmpg := big.NewInt(0)
			tmpg.Set(g)
			// SetPublicKey(g, p, PElgamal)
			gexpb = Fermat(tmpg, b, tmpp)
			gb.Set(gexpb)
			gexpab = Fermat(tmpPElgamal, big.NewInt(1968), p)
			tmpgb := big.NewInt(0)
			tmpgab := big.NewInt(0)
			bytegb := []byte(tmpgb.Set(gexpb).String())
			bytegab := []byte(tmpgab.Set(gexpab).String())
			bytega := []byte(tmpPElgamal2.String())
			tt := append(bytega, bytegb...)
			tt = append(tt, bytegab...)
			k = GetSHA256HashCode(tt)
			//fmt.Println("k:", k)
		}

		if idx == 3 {
			ciphertext, nonce, _ := GCMEncrypt(k, []byte(message))
			gbstring := gb.String()
			noncestring := hex.EncodeToString(nonce)
			d := gbstring + "\n" + noncestring + hex.EncodeToString(ciphertext)
			filename2 := args
			file2, err := os.Create(filename2)
			check(err)
			_, err = io.WriteString(file2, d)
		}
	}
}
