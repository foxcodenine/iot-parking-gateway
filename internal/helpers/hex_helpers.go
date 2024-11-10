package helpers

import (
	"fmt"
	"math/big"
	"strconv"
)

// parseHexSubstring extracts an integer value from a specified hex substring
// within a larger hex string, based on the offset and length in bytes.
func ParseHexSubstring(hexStr string, byteOffset, byteLength int) (int, int, error) {
	byteOffset *= 2 // Convert byte offset to hex characters
	byteLength *= 2 // Convert byte length to hex characters
	endIndex := byteOffset + byteLength

	// Check if endIndex is within the hexStr bounds
	if endIndex > len(hexStr) {
		return 0, byteOffset / 2, fmt.Errorf("substring out of range")
	}

	hexSubStr := hexStr[byteOffset:endIndex] // Extract hex substring

	value64, err := strconv.ParseInt(hexSubStr, 16, 64) // Parse hex to int64
	if err != nil {
		return 0, byteOffset / 2, err
	}

	// Convert int64 to int and adjust endIndex back to bytes
	value := int(value64)
	nextOffset := endIndex / 2

	return value, nextOffset, nil
}

// ParseHexSubstringBigInt extracts a large integer value from a specified hex substring
// within a larger hex string, based on the offset and length in bytes, returning *big.Int.
func ParseHexSubstringBigInt(hexStr string, byteOffset, byteLength int) (*big.Int, int, error) {
	byteOffset *= 2 // Convert byte offset to hex characters
	byteLength *= 2 // Convert byte length to hex characters
	endIndex := byteOffset + byteLength

	// Check if endIndex is within the hexStr bounds
	if endIndex > len(hexStr) {
		return nil, byteOffset / 2, fmt.Errorf("substring out of range")
	}

	hexSubStr := hexStr[byteOffset:endIndex] // Extract hex substring

	// Parse hex substring to big.Int
	value := new(big.Int)
	_, success := value.SetString(hexSubStr, 16)
	if !success {
		return nil, byteOffset / 2, fmt.Errorf("failed to parse hex substring as big integer")
	}

	// Adjust endIndex back to bytes
	nextOffset := endIndex / 2

	return value, nextOffset, nil
}
