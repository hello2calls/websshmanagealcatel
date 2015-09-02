package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
)

// WriteFile encrypt data and write the file
func WriteFile(dataFile S.Data, file string) {
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	key := []byte("nit5i8drod0it1of0en8hylg9ov0out4") // 32 bytes!
	plaintext := []byte(jsonIndent)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	ioutil.WriteFile(file, ciphertext, 0777)
}

// ReadFile decrypt data and read the file
func ReadFile(file string) S.Data {
	var ciphertext, _ = ioutil.ReadFile(file)

	key := []byte("nit5i8drod0it1of0en8hylg9ov0out4") // 32 bytes!

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	var dataFile S.Data
	err = json.Unmarshal(ciphertext, &dataFile)
	if err != nil {
		fmt.Println("JSON Unmarshal file in indexHandler error:", err)
	}

	return dataFile
}
