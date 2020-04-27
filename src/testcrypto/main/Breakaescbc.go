package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/mkideal/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const DescriptionTemplate = `
decrypt-attack -i <ciphertext file> -i <input file> -O <output file>

Where <mode> is one of either encrypt or decrypt, and the input/output ﬁles contain raw binary data. 
You should parse the ﬁrst 16 bytes of the key as the encryption key k enc , 
and the second 16 bytes as the MAC key k mac .
YOu can run the program as:
go run test2.go -I "test.txt" -O "result.txt"

Enjoy!
`

type CLIOpts struct {
	Help      bool   `cli:"!h,help" usage:"Show help."`
	Condensed bool   `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`
	I         string `cli:"I" usage:"The path and name of input file."`
	O         string `cli:"O" usage:"The path and name of output file."`
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

func Validpadding(origData []byte, wanted int, place int, blocknumber int) int {
	length := len(origData)
	unpadding := int(origData[length-place-(blocknumber-1)*16])
	if unpadding > 16 || unpadding == 0 {
		fmt.Println("Invalid Padding")
		return 0
	}
	if unpadding == wanted {
		return wanted
	}
	return 0
}

func AES_CBC_Decrypt(cipherText []byte, iv []byte, key []byte, place int, blocknumber int, want int) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	n := len(cipherText)
	//fmt.Println("cipher text:", cipherText)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText) //depands on program2
	//fmt.Println("decypher plaintext:", plainText)
	wanted := want
	res := Validpadding(plainText, wanted, place, blocknumber)
	if res == wanted {
		needbyte := cipherText[n-(blocknumber)*16-place]
		return []byte{needbyte}
	}
	return nil
}

func AddIv(cipherText []byte, iv []byte) []byte {
	Final := append(iv, cipherText...)
	//fmt.Println(Final)
	return Final
}

func DeIv(plaintext []byte, iv []byte) []byte {
	First := plaintext[len(iv):]
	//fmt.Println(First)
	return First
}

func Changebyte(cipher []byte, iv []byte, Ksec []byte, Kmac []byte) []byte {
	n := len(cipher)
	blocknumbers := int(n / 16)
	cipheruse := make([]byte, n)
	var plaintext []byte
	for i := 1; i < blocknumbers; i = i + 1 {
		copy(cipheruse, cipher)
		ciphertext2 := cipheruse[n-16*i-16 : n-16*i]
		for j := 15; j >= 0; j-- {
			for w := 1; w <= 16; w = w + 1 {
				cipherbyte := ciphertext2[j]
				for k := 0; k < 255; k = k + 1 {
					ciphertext2[j] = byte(k)
					_ = append(append(cipheruse[:n-16*i-16], ciphertext2...), cipheruse[n-i*16:]...)
					getbyte := DeAESCBCHmac(cipheruse, iv, Ksec, Kmac, 16-j, i, w)
					if getbyte == nil {
						continue
					} else {
						if w == 16-j {
							plainbyte := getbyte[0] ^ cipherbyte ^ byte(w)
							plaintext = append([]byte{plainbyte}, plaintext...)
							k = k + 255
						} else {
							ciphertext2[j] = getbyte[0]
							_ = append(append(cipheruse[:n-16*i-16], ciphertext2...), cipheruse[n-i*16:]...) //still need something
							j = j - 1
							w = w - 1
						}
						k = k + 255
					}
				}
			}
		}
	}
	return plaintext
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
		if argv.Help || len(argv.I) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}
		path2 := argv.O
		iv := []byte("aaaaaaaaaaaaaaaa")
		Ksec := []byte("aaaaaaaaaaaaaaaa")
		Kmac := []byte("aaaaaaaaaaaaaaaa")
		plaintext := Changebyte(m, iv, Ksec, Kmac)

		plainf := PKCS5UnPadding(plaintext)
		plainf = plainf[:len(plainf)-32]
		fmt.Printf("%s", string(plainf))
		err = ioutil.WriteFile(path2, plainf, 0644)
		check(err)
		return nil
	})

}
