package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func preDecrypt(iv []byte, ciphertext []byte) []byte {
	ciphertext = append(iv, ciphertext...)
	fmt.Println(ciphertext)
	buffer := make([]byte, len(ciphertext))
	blockcipher := make([]byte, 32)
	copy(buffer, ciphertext)
	blocknums := len(ciphertext) / 16
	for i := blocknums - 1; i > 0; i-- {
		copy(blockcipher, buffer[i*16-16:i*16+16])
		plaintext := Decrypt(blockcipher)
		copy(buffer[i*16:i*16+16], plaintext)
	}
	return buffer[16:]
}
func Decrypt(blockcipher []byte) []byte {
	pre_ciphertext := make([]byte, 16)
	copy(pre_ciphertext, blockcipher[len(blockcipher)-32:len(blockcipher)-16])
	recoveredText := make([]byte, 16)
	intermedia := make([]byte, 16)
	mod_ciphertext := blockcipher[len(blockcipher)-32 : len(blockcipher)-16]
	_, err := rand.Read(mod_ciphertext)
	if err != nil {
		panic(err)
	}
	for i := 15; i >= 0; i-- {
		pad := byte(16 - i)
		for j := i + 1; j < 16; j++ {
			mod_ciphertext[j] = pad ^ intermedia[j]
		}
		for k := 0x00; k < 0x100; {
			mod_ciphertext[i] = byte(k)
			fmt.Println(blockcipher)
			fmt.Println("we want pad:", pad)
			ioutil.WriteFile("temp", blockcipher, 0777)
			output, _ := exec.Command("./decrypt-test", "-i", "temp").CombinedOutput()
			fmt.Println(string(output))
			//if err1 != nil{
			// panic(err1)
			//}
			if !strings.Contains(string(output), "INVALID PADDING") {
				fmt.Println("1111111")
				break
			}
			k++
		}

		intermedia[i] = pad ^ mod_ciphertext[i]
	}
	for i := range intermedia {
		recoveredText[i] = intermedia[i] ^ pre_ciphertext[i]
	}
	fmt.Println(recoveredText)
	return recoveredText
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Invalid input argument")
		fmt.Println("Expected input is: decrypt-attack -i <ciphertext file>")
		os.Exit(1)
	}
	ciphertext_iv, _ := ioutil.ReadFile(os.Args[2])
	ciphertext := make([]byte, len(ciphertext_iv)-16)
	copy(ciphertext, ciphertext_iv[16:])
	iv := ciphertext_iv[:16]
	raw_plaintext := preDecrypt(iv, ciphertext)
	padLen := int(raw_plaintext[len(raw_plaintext)-1])
	rem := len(raw_plaintext) - padLen - 32
	fmt.Println(rem, len(raw_plaintext))
	plaintext := raw_plaintext[0:rem]
	fmt.Print(string(plaintext))
	return
}
