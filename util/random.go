package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvxyz"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomString(length int) string {
	var sb strings.Builder
	var alphLength = len(alphabet)

	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rand.Intn(alphLength)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(-1000, 500000000)
}

func RandomCurrency() string {
	currs := []string{"EUR", "USD", "CAD"}
	n := len(currs)
	return currs[rand.Intn(n)]
}
