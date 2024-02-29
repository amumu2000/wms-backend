package utils

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	salt string
)

func InitCrypto(md5Salt string) {
	salt = md5Salt
}

func MD5Salt(str string) string {
	b := []byte(str)
	s := []byte(salt)
	h := md5.New()
	h.Write(s)
	h.Write(b)
	var res []byte
	res = h.Sum(nil)
	for i := 0; i < 2; i++ {
		h.Reset()
		h.Write(res)
		res = h.Sum(nil)
	}
	return hex.EncodeToString(res)
}
