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
decrypt-attack -i <ciphertext file> -i <input file> -O <output file>

You can use this to decypt the AESCTR mode. I just embed the test function in this program.
The name of the test program is Test2_AEC_CRT_DeCrypt, Test1_AEC_CRT_DeCrypt. You can just change password in this teo function
YOu can run the program as:
go run decrypt-attack-crc.go -I "test.txt" -O "result.txt"

Enjoy!
`

type CLIOpts struct {
	Help      bool   `cli:"!h,help" usage:"Show help."`
	Condensed bool   `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`
	I         string `cli:"I" usage:"The path and name of input file."`
	O         string `cli:"O" usage:"The path and name of output file."`
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
	fmt.Println("***real checksum:", res)
	M := append(IntToBytes(res), message...)
	return M
}

func Dechecksum(checksum []byte, message []byte) []byte {
	checkstring := message[:8]
	r := BytesToInt(checksum)
	res := r % 256
	if BytesToInt(checkstring) != res {
		fmt.Println("Checksum verify falied")
	}
	M := message[8:]

	return M
}

func Verichecksum(message []byte) int {
	checksum := []byte("okjiuxxxx")
	padd := []byte{0, 0, 0, 0, 0, 0, 0}
	checkstring := []byte{message[7]}
	checkstring = append(padd, checkstring...)
	r := BytesToInt(checksum)
	res := r % 256
	if BytesToInt(checkstring) != res {
		//fmt.Println("Checksum verify falied")
		return 0
	}
	return res
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
	return message
}

func Test_AEC_CRT_DeCrypt(text []byte, iv []byte) int {
	key := []byte("aaaaaaaaaaaaaaaa")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCTR(block, iv)
	message := make([]byte, len(text))
	blockMode.XORKeyStream(message, text)
	r := Verichecksum(message)
	if r == 0 {
		return 0
	}
	return r
}

func Verichecksum2(message []byte, checksum int, place int) []byte {
	//fmt.Println("*********message[place]", message[place])
	checkstring := []byte{message[place]}
	padd := []byte{0, 0, 0, 0, 0, 0, 0}
	checkstring = append(padd, checkstring...)
	//fmt.Println("checkstring:", BytesToInt(checkstring))
	if BytesToInt(checkstring) != checksum {
		return nil
	}
	return []byte{message[place]}
}

func Test2_AEC_CRT_DeCrypt(text []byte, checksum int, place int, iv []byte) []byte {
	key := []byte("aaaaaaaaaaaaaaaa")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCTR(block, iv)
	message := make([]byte, len(text))
	blockMode.XORKeyStream(message, text)
	r := Verichecksum2(message, checksum, place)
	if r == nil {
		return nil
	}
	return []byte{text[place]}
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

func findchecksum(ciphertext []byte, iv []byte) int {
	var v int
	//usebyte := ciphertext[7]
	for i := 0; i < 255; i = i + 1 {
		ciphertext[7] = byte(i)
		res := Test_AEC_CRT_DeCrypt(ciphertext, iv)
		if res == 0 {
			continue
		}
		v = res
		break
	}
	return v
}

func Attack_plaintext(ciphertext []byte, checksum int, iv []byte) []byte {
	var usebyte byte
	var plaintext []byte
	n := len(ciphertext)
	for i := 0; i < n; i = i + 1 {
		usebyte = ciphertext[i]
		//fmt.Println("usebyte:", usebyte)
		for j := 0; j < 255; j = j + 1 {
			ciphertext[i] = byte(j)
			//print("	i:", i)
			r := Test2_AEC_CRT_DeCrypt(ciphertext, checksum, i, iv)
			if r == nil {
				continue
			}
			//fmt.Println("r:", r)
			p := r[0] ^ usebyte ^ byte(checksum)
			plaintext = append(plaintext, []byte{p}...)
			j = j + 255
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
		iv := m[:16]
		m = m[16:]

		getchecksum := findchecksum(m, iv)
		attackplain := Attack_plaintext(m, getchecksum, iv)
		attackplain = attackplain[8:]

		fmt.Printf("%s", string(attackplain))
		err = ioutil.WriteFile(path2, attackplain, 0644)
		check(err)
		return nil
	})

}
