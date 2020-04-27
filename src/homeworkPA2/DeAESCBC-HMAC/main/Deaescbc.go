package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/mkideal/cli"
)

const DescriptionTemplate = `
usage: decrypt-attack -i <ciphertext file> -O <output file>

Where <mode> is one of either encrypt or decrypt, and the input/output ﬁles contain raw binary data. 
You should parse the ﬁrst 16 bytes of the key as the encryption key k enc , 
and the second 16 bytes as the MAC key k mac .
YOu can run the program as:
go run Deaescbc.go -I "test.txt" -O "result.txt"

Enjoy!
`

type CLIOpts struct {
	Help bool   `cli:"!h,help" usage:"Show help."`
	I    string `cli:"I" usage:"The path and name of input file."`
	O    string `cli:"O" usage:"The path and name of output file."`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IntToBytes(i int) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
func BytesToInt(buf []byte) int {
	return int(binary.BigEndian.Uint64(buf))
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func HashMac(Kmac []byte, M []byte) [32]byte {
	ipad := 0b0011011000110110001101100011011000110110001101100011011000110110
	opad := 0b0101110001011100010111000101110001011100010111000101110001011100

	var s1Int int = BytesToInt(Kmac) ^ ipad
	s1 := IntToBytes(s1Int)
	plainText2 := append(s1, M...)
	hashP1 := sha256.Sum256(plainText2)
	var s2Int int = BytesToInt(Kmac) ^ opad
	s2 := IntToBytes(s2Int)
	plainText3 := append(s2, hashP1[:]...)
	hmac := sha256.Sum256(plainText3)
	return hmac

}

func ComputeM2(m []byte, Tag [32]byte) []byte {
	M2 := append(m, Tag[:]...)
	return M2
}

func deM2(m []byte) []byte {
	res := m[:len(m)-32]
	return res
}

func AES_CBC_Decrypt(cipherText []byte, iv []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = PKCS5UnPadding(plainText)
	return plainText
}

func VerifyHamc(m []byte, Kmac []byte) int {
	ver := m[len(m)-32:]
	mes := m[:len(m)-32]
	Tag := HashMac(mes, Kmac)
	return bytes.Compare(Tag[:], ver)
}

func DeIv(plaintext []byte, iv []byte) []byte {
	First := plaintext[len(iv):]
	//fmt.Println(First)
	return First
}

func DeAESCBCHmac(cipherText []byte, iv []byte, K []byte) []byte {
	p := DeIv(cipherText, iv)
	Ksec := K
	Kmac := []byte("aaaaaaaaaaaaaaaa")
	plainText := AES_CBC_Decrypt(p, iv, Ksec)
	if VerifyHamc(plainText, Kmac) != 0 {
		//fmt.Printf("Hmac Verify Failed!!!")
		return nil
	}
	dem2 := deM2(plainText)
	return dem2
}

func GetIv(message []byte) []byte {
	iv := message[:16]
	return iv
}

func main() {
	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		path := argv.I
		if argv.Help || len(path) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}
		message, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Print(err)
		}
		str := string(message)
		m := []byte(str)
		path2 := argv.O
		posiblelist := []byte("abcdefghijklmnopqrstuvwxyz1234567890")
		s := len(posiblelist)
		//key := make([]byte, 32)
		//firstkey := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		firstkey := bytes.Repeat([]byte("a"), 16)
		iv := GetIv(m)
		return nil
	})
}
