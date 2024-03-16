package random

import (
	"math/rand"
	"time"
)

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var gen = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// String generates a random string
func String(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = runes[gen.Intn(len(runes))]
	}
	return string(b)
}
