package common

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Add(a, b int) int {
	return a + b
}

func Dec(a, b int) int {
	return a - b
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func StrTrip(src string) string {
	str := strings.Replace(src, " ", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}


