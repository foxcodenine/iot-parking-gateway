package helpers

import (
	"fmt"
	"time"
)

func GetCurrentTimestampHex() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%08x", timestamp)
}
