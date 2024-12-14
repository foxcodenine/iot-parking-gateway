package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateJWTSecretKey(length int) (string, error) {
	// Create a byte slice to hold the random data
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode the random bytes to a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
