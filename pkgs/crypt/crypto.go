package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Crypt struct {
	key []byte
}

var DefaultCrypt = &Crypt{
	key: []byte("x^)dixf&*1$free]"),
}

func (that *Crypt) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (that *Crypt) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (that *Crypt) AesEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(that.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = that.pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, that.key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (that *Crypt) AesDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(that.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, that.key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = that.pKCS7UnPadding(origData)
	return origData, nil
}

func DecodeBase64(str string) (res string) {
	count := len(str) % 4
	for i := 0; i < count; i++ {
		str += "="
	}
	if s, err := base64.StdEncoding.DecodeString(str); err == nil {
		res = string(s)
	}
	return
}