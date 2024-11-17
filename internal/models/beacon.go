package models

// Beacon struct represents the JSON structure within the beacons JSONB column.
type Beacon struct {
	BeaconNumber int `json:"beacon_number"` // Unique number for each beacon within an activity
	Major        int `json:"major"`         // Major identifier of the beacon
	Minor        int `json:"minor"`         // Minor identifier of the beacon
	RSSI         int `json:"rssi"`          // Received Signal Strength Indicator (RSSI) value
}
