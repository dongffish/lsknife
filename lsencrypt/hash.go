package lsencrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Hash calc md5 code for origData
// 返回 计算所得MD5码
func MD5Hash(origData string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(origData))
	if err != nil {
		return "", err
	}
	s := h.Sum(nil)
	return hex.EncodeToString(s), nil
}
