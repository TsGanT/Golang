package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/mkideal/cli"
)

const DescriptionTemplate = `
usage: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -O <output file>

Where <mode> is one of either encrypt or decrypt, and the input/output ﬁles contain raw binary data. 
You should parse the ﬁrst 16 bytes of the key as the encryption key k enc , 
and the second 16 bytes as the MAC key k mac .
YOu can run the program as:
go run aescbc.go --Auth "decrypt" -K 123tangshitangsh87654321abcdefgh -I "test.txt" -O "result.txt"

Enjoy!
`

type CLIOpts struct {
	Help      bool   `cli:"!h,help" usage:"Show help."`
	Condensed bool   `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`
	Auth      string `cli:"Auth" usage:"Chose encrypt or decrypt mode."`
	K         string `cli:"K" usage:"enter the 16 byte Ksec."`
	I         string `cli:"I" usage:"The path and name of input file."`
	O         string `cli:"O" usage:"The path and name of output file."`
	Checksum  string `cli:"Checksum" usage:"You can input any checksum you want."`
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

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Addchecksum(checksum []byte, message []byte) []byte {
	r := BytesToInt(checksum)
	res := r % 256
	M := append(IntToBytes(res), message...)
	return M
}

func Dechecksum(checksum []byte, message []byte) []byte {
	checkstring := message[:8]
	r := BytesToInt(checksum)
	res := r % 256
	//fmt.Println(res)
	if BytesToInt(checkstring) != res {
		fmt.Println("Checksum verify falied")
	}
	M := message[8:]

	return M
}

func AEC_CRT_DeCrypt(text []byte, key []byte, iv []byte, checksum []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCTR(block, iv)
	message := make([]byte, len(text))
	blockMode.XORKeyStream(message, text)
	rres := Dechecksum(checksum, message)
	return rres
}

func AEC_CRT_EnCrypt(text []byte, key []byte, iv []byte, checksum []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	M := Addchecksum(checksum, text)
	blockMode := cipher.NewCTR(block, iv)
	message := make([]byte, len(M))
	blockMode.XORKeyStream(message, M)
	message = append(iv, message...)
	return message
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
		checksum := []byte(argv.Checksum)

		path2 := argv.O
		if argv.Auth == "decrypt" {
			ivd := m[:16]
			m = m[16:]
			plaintext := AEC_CRT_DeCrypt(m, key, ivd, checksum)
			if plaintext == nil {
				fmt.Println("You may under hacking...")
			} else {
				fmt.Println("Decrypt success!")
				err := ioutil.WriteFile(path2, plaintext, 0644)
				check(err)
			}
		} else {
			i := RandStringRunes(16)
			ive := []byte(i)
			fmt.Println("Encrypt success!")
			fmt.Println("This is your random IV, Please remember it!!", string(ive))
			ciphertext := AEC_CRT_EnCrypt(m, key, ive, checksum)
			err := ioutil.WriteFile(path2, ciphertext, 0644)
			check(err)
		}
		return nil
	})
	// key := []byte("aaaaaaaaaaaaaaaa")
	// message := []byte("aaaaaaaaa")
	// checksum := []byte("asdsdasd")
	// ive := []byte("aaaaaaaaaaaaaaaa")
	// fmt.Println("checksum firsttime:", checksum)
	// ciphertext := AEC_CRT_EnCrypt(message, key, ive, checksum)
	// fmt.Println("ciphertext", ciphertext)
	// plaintext := AEC_CRT_DeCrypt(ciphertext, key, ive, checksum)
	// fmt.Println("plaintext:", plaintext)
}
