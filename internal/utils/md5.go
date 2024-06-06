package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Crypt(plainText string) string {
	hash := md5.New()
	hash.Write([]byte(plainText))
	cypher := hash.Sum([]byte(plainText))
	return hex.EncodeToString(cypher)
}
