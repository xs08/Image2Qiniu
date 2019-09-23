package utils

import (
	"bytes"
	"math/rand"
	"time"
)

// JoinStrs join strs
func JoinStrs(strs ...string) string {
	var buf bytes.Buffer
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf.String()
}

// RandomStr get a random string with length
func RandomStr(len int) string {
	if len <= 0 {
		return ""
	}
	chars := "1234567890qwertyuiopasdfghjklzxcvbnm"
	var strBuffer bytes.Buffer

	for len > 0 {
		len--
		strBuffer.WriteByte(chars[rand.Intn(len)])
	}
	return strBuffer.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
