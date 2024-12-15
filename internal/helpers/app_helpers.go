package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

// StructToMap converts a struct to a map[string]any using JSON marshaling and unmarshaling.
func StructToMap(v any) (map[string]any, error) {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %v", err)
	}

	// Unmarshal JSON into a map
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into map: %v", err)
	}

	return result, nil
}
