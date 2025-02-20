package models

import (
	"encoding/json"
	"errors"
)

// Beacon struct represents the JSON structure within the beacons JSONB column.
type Beacon struct {
	BeaconNumber int `json:"beacon_number"` // Unique number for each beacon within an activity
	Major        int `json:"major"`         // Major identifier of the beacon
	Minor        int `json:"minor"`         // Minor identifier of the beacon
	RSSI         int `json:"rssi"`          // Received Signal Strength Indicator (RSSI) value
}

// ---------------------------------------------------------------------

// BeaconSlice is used to handle JSONB data scanning
type BeaconSlice []Beacon

// Scan implements the sql.Scanner interface for BeaconSlice.
// This method is automatically called by Go’s database/sql package when scanning
// a database query result into a BeaconSlice field within a struct.
func (bs *BeaconSlice) Scan(src interface{}) error {
	// Case 1: If the database column contains a NULL value, `src` will be nil.
	// We explicitly set `bs` to nil to match the database's representation.
	if src == nil {
		*bs = nil // This ensures that the BeaconSlice field is correctly represented as nil in Go.
		return nil
	}

	// Case 2: The database stores JSONB data as a binary format ([]uint8 in Go).
	// We need to confirm that `src` is indeed a byte slice ([]byte).
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed") // Handle unexpected types gracefully.
	}

	// Case 3: If the JSONB column is empty (i.e., contains an empty JSON object `{}` or an empty array `[]`),
	// the length of `b` would be zero. We treat this as a `nil` BeaconSlice to maintain consistency.
	if len(b) == 0 {
		*bs = nil // Ensure the field is treated as nil in Go when there's no actual data.
		return nil
	}

	// Case 4: If the data exists and is valid JSON, we attempt to unmarshal it into the BeaconSlice.
	// json.Unmarshal automatically converts the JSON byte slice into a slice of Beacon structs.
	return json.Unmarshal(b, bs)
}
