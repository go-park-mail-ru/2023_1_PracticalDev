package shortener

import (
	"math/rand"
	"time"
)

func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	rn := rand.New(rand.NewSource(time.Now().Unix()))

	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[rn.Intn(len(charset))]
	}

	return string(shortCode)
}
