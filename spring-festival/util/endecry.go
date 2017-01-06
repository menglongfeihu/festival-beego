package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	//"strconv"

	//"github.com/astaxie/beego/logs"
)

const MAX_ENCRYPT_BLOCK = 117

func ReadBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// RSA加密明文最大长度117字节，所以在加密的过程中需要分块进行
func RSAEncrypt(data []byte) ([]byte, error) {
	publicKey, err := ReadBytes(`conf/public.pem`)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	inputLen := len(data)
	if inputLen > MAX_ENCRYPT_BLOCK {
		var cache bytes.Buffer
		chunk := inputLen / MAX_ENCRYPT_BLOCK
		var lastPosition int
		for index := 0; index < chunk; index++ {
			begin := index * MAX_ENCRYPT_BLOCK
			end := (index + 1) * MAX_ENCRYPT_BLOCK
			temp, _ := rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data[begin:end])
			lastPosition = end
			cache.Write(temp)
		}
		temp, _ := rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data[lastPosition:inputLen])
		cache.Write(temp)
		return cache.Bytes(), nil
	} else {
		return rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data)
	}

}

func MD5Encrypt(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	result := h.Sum(nil)
	return result
}
