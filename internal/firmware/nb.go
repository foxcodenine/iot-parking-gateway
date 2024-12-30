package firmware

import "github.com/foxcodenine/iot-parking-gateway/internal/helpers"

// Validates event ID
func isValidEventID(eventID int) bool {
	validEventIDs := []int{6, 10, 26}
	return helpers.Contains(validEventIDs, eventID)
}

const (
	Multiplier256 = 256
	Divider4      = 4
	Divider2      = 2
	Divider16     = 16
)
