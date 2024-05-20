package helpers

import (
	"math/rand"
	"time"
)

func GenerateCouponCode(prefix string) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate the code by selecting random characters from the charset
	code := make([]byte, 10)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return prefix + string(code)
}
