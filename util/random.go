package util

import (
	"math/rand"
	"strings"
	"time"
)

const alpha = "1234567890abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alpha)

	for i := 0; i < n; i++ {
		c := alpha[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
