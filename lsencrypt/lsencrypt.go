package lsencrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

var (
	// Keys 默认密钥Key,也是用来补充Key长度用的。
	// 当用户调用加、解密时，一般会输入密钥，此时这个keys就会失效。所以在这里明文也不会造成密钥泄露
	Keys = "HzLv&ShGao@#HIC^(2).$China*202412"
)

func normalizeKey(key []byte, keyLen int) []byte {
	strKey := string(key)
	retStr := ""
	s := len(strKey)
	if s == 0 {
		retStr = ""
	} else if s >= keyLen {
		retStr = strKey[:keyLen]
	} else {
		retStr = strKey + Keys[s:keyLen]
	}
	return []byte(retStr)
}

// DesEncrypt DES加密
func DesEncrypt(origData, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	key = normalizeKey(key, 8)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	//origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return []byte(base64.StdEncoding.EncodeToString(crypted)), nil
}

// DesDecrypt DES解密
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	key = normalizeKey(key, 8)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	rawstr, _ := base64.StdEncoding.DecodeString(string(crypted))
	blockMode.CryptBlocks(origData, rawstr)
	origData = PKCS5UnPadding(origData)
	//origData = ZeroUnPadding(origData)
	return origData, nil
}

// TripleDesEncrypt 3DES加密,返回加密后的数据,base64编码
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	key = normalizeKey(key, 24)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	//origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return []byte(base64.StdEncoding.EncodeToString(crypted)), nil
}

// TripleDesDecrypt 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	key = normalizeKey(key, 24)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	// origData := crypted
	rawstr, _ := base64.StdEncoding.DecodeString(string(crypted))
	origData := make([]byte, len(rawstr))
	blockMode.CryptBlocks(origData, rawstr)
	origData = PKCS5UnPadding(origData)
	//origData = ZeroUnPadding(origData)
	return origData, nil
}

// ZeroPadding zero padding
// 返回ZeroPadding后的字节数组
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// ZeroUnPadding zero unpadding
// 返回ZeroUnPadding后的数据
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// PKCS5Padding pkcs5 padding
// 返回PKCS5Padding后的数据
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding pkcs5 unpadding
// 返回PKCS5UnPadding后的数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
