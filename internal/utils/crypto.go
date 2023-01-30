package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, key string) (string, error) {
	//encrypt data.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Hash(text string) string {
	encryptKey := text + "wehavetouseenvfile"
	hasher := sha1.New()
	hasher.Write([]byte(encryptKey))
	hashKey := base32.HexEncoding.EncodeToString(hasher.Sum(nil))
	return hashKey
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, key string) (string, error) {
	//key is just 32bit hash key from database.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
