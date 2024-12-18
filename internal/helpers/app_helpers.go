package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
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

// StructSliceToMapSlice converts a slice of structs into a slice of map[string]any.
func StructSliceToMapSlice(slice any) ([]map[string]any, error) {
	// Ensure input is a slice
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice of structs")
	}

	// Resultant slice of maps
	var result []map[string]any

	// Iterate over the slice
	for i := 0; i < sliceValue.Len(); i++ {
		// Convert each struct to map
		item := sliceValue.Index(i).Interface()
		structMap, err := StructToMap(item)
		if err != nil {
			return nil, fmt.Errorf("failed to convert struct at index %d: %v", i, err)
		}
		result = append(result, structMap)
	}

	return result, nil
}
