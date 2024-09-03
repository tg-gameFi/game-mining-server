package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs: invalid data")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func aesEncryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	encrypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	if len(encryptBytes)%blockMode.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return encrypted, nil
}

func aesDecryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	if len(data)%blockMode.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	decrypted := make([]byte, len(data))
	blockMode.CryptBlocks(decrypted, data)
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func EncryptByAes(data string, key string) (string, error) {
	res, err := aesEncryptBytes([]byte(data), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptByAes(data string, key string) (string, error) {
	dataByte, e0 := base64.StdEncoding.DecodeString(data)
	if e0 != nil {
		return "", e0
	}
	resByte, e1 := aesDecryptBytes(dataByte, []byte(key))
	if e1 != nil {
		return "", e1
	} else {
		return string(resByte), nil
	}
}
