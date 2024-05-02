package random

import (
	"math/rand"
	"time"
)

const letterChars = "0123456789"

// RandNumberToString 随机返回指定位数字符串
func RandNumberToString(n int) string {
	randInstance := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for index := range b {
		b[index] = letterChars[randInstance.Intn(len(letterChars))]
	}
	return string(b)
}
