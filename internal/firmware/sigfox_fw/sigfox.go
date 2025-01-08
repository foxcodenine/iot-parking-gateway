package sigfoxfw

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// Validates event ID
func isValidEventID(eventID int) bool {
	validEventIDs := []int{6, 10, 26, 31}
	return helpers.Contains(validEventIDs, eventID)
}
