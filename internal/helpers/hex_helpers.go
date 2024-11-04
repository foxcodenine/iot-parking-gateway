package helpers

import (
	"strconv"
	"strings"
)

func HexSliceToBase10(hexSlice []string) (int64, error) {
	// Step 1: Join the hex values into a single string
	hexString := strings.Join(hexSlice, "")

	// Step 2: Convert the hex string to a base-10 integer
	base10Value, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return 0, err
	}

	return base10Value, nil
}
