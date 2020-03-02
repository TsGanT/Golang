package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/mkideal/cli"
)

const DescriptionTemplate = `
usage: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -o <output file>

Where <mode> is one of either encrypt or decrypt, and the input/output ﬁles contain raw binary data. 
You should parse the ﬁrst 16 bytes of the key as the encryption key k enc , 
and the second 16 bytes as the MAC key k mac .
YOu can run the program as:
go run test.go --Auth "decrypt" -K 123tangshitangsh87654321abcdefgh -I "test.txt" -O "result.txt"

Enjoy!
`

type CLIOpts struct {
	Help      bool   `cli:"!h,help" usage:"Show help."`
	Condensed bool   `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`
	Auth      string `cli:"Auth" usage:"Chose encrypt or decrypt mode."`
	K         string `cli:"K" usage:"enter the Ksec and Kmac, the ﬁrst 16 bytes of the key as the encryption key kenc , and the second 16 bytes as the MAC key kmac."`
	I         string `cli:"I" usage:"The path and name of input file."`
	O         string `cli:"O" usage:"The path and name of output file."`
	IV        string `cli:"IV" usage:"The will be only use when decryption."`
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

func Padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}
func UnPadding(cipherText []byte) []byte {
	end := cipherText[len(cipherText)-1]
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
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
func ComputeM3(m []byte) []byte {
	n := 16 - len(m)%16
	var PS []byte
	if n == 0 {
		PS = bytes.Repeat(IntToBytes(16), 16)
	}
	if n != 0 {
		PS = bytes.Repeat(IntToBytes(16-n), 16-n)
	}
	M3 := append(m, PS...)
	return M3
}

func deM3(m []byte) []byte {
	a := m[len(m)-1]
	// b := m[len(m)-8 : len(m)]
	// c := bytes.Repeat(b, int(a))
	if int(a) > 16 {
		return nil
	}
	n := int(a) * 8
	res := m[:len(m)-n]
	return res
}
func deM2(m []byte) []byte {
	res := m[:len(m)-32]
	return res
}

func AES_CBC_Encrypt(plainText []byte, iv []byte, key []byte) []byte {
	block, err := aes.NewCipher(key) //feed back a block interface(something like class in other language)
	if err != nil {
		panic(err)
	}
	plainText = Padding(plainText, block.BlockSize()) //padding here
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText
}
func AES_CBC_Decrypt(cipherText []byte, iv []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = UnPadding(plainText)
	return plainText
}

func VerifyHamc(m []byte, Kmac []byte) int {
	ver := m[len(m)-32:]
	mes := m[:len(m)-32]
	Tag := HashMac(mes, Kmac)
	return bytes.Compare(Tag[:], ver)
}

func EnAESCBCHmac(plainText []byte, iv []byte, Ksec []byte, Kmac []byte) []byte {
	Tag := HashMac(plainText, Kmac)
	M2 := ComputeM2(plainText, Tag)
	M3 := ComputeM3(M2)
	cipherText := AES_CBC_Encrypt(M3, iv, Ksec)
	return cipherText
}

func DeAESCBCHmac(cipherText []byte, iv []byte, Ksec []byte, Kmac []byte) []byte {
	plainText := AES_CBC_Decrypt(cipherText, iv, Ksec)
	dem3 := deM3(plainText)
	if dem3 == nil {
		fmt.Printf("You used a wrong key!!")
		return nil
	}
	if VerifyHamc(dem3, Kmac) != 0 {
		fmt.Printf("Hmac Verify Failed!!!")
		return nil
	}
	dem2 := deM2(dem3)
	return dem2
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {

		argv := ctx.Argv().(*CLIOpts)
		path := argv.I
		message, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Print(err)
		}
		str := string(message)
		m := []byte(str)
		if argv.Help || len(argv.K) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}
		key := []byte(argv.K)
		Ksec := key[:16]
		Kmac := key[16:]
		ivd := []byte(argv.IV)
		path2 := argv.O

		if argv.Auth == "decrypt" {
			Plaintext := DeAESCBCHmac(m, ivd, Ksec, Kmac)
			if Plaintext == nil {
				fmt.Println("You may under hacking...")
			} else {
				fmt.Println("Decrypt success!")
				err := ioutil.WriteFile(path2, Plaintext, 0644)
				check(err)
			}
		} else {
			i := RandStringRunes(16)
			ive := []byte(i)
			fmt.Println("Encrypt success!")
			fmt.Println("This is your random IV, Please remember it!!", string(ive))
			Ciptertext := EnAESCBCHmac(m, ive, Ksec, Kmac)
			err := ioutil.WriteFile(path2, Ciptertext, 0644)
			check(err)
		}
		return nil
	})
}
