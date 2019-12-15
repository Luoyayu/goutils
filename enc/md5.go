package enc

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5HexStr return hexString of string input
func MD5HexStr(text string) string {
	md5Bytes := md5.Sum([]byte(text))
	return hex.EncodeToString(md5Bytes[:])
}
