// get file name with flag
// go run main.go -name=hello

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	file := flag.String("f", "", "a file name to read")
	flag.Parse()

	if *file == "" {
		fmt.Println("no file name")
		return
	}

	content, err := readFile(*file)

	if err != nil {
		fmt.Println(err)
		return
	}

	encrypt(content)
	decryptedText := decrypt()

	fmt.Println(decryptedText)

}

func readFile(in string) (string, error) {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	fmt.Println("Contents of file:", string(data))

	return string(data), nil
}

func encrypt(content string) {

	key := []byte("passphrasewhichneedstobe32bytes!")
	plaintext := []byte(content)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("myfile.data", gcm.Seal(nonce, nonce, plaintext, nil), 0777)
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}
}

func decrypt() string {
	key := []byte("passphrasewhichneedstobe32bytes!")
	ciphertext, err := ioutil.ReadFile("myfile.data")
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	// delete file
	err = os.Remove("myfile.data")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(plaintext)
}
